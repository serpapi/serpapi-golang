package serpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	VERSION        = "1.5.0"
	BaseURL        = "https://serpapi.com"
	DefaultTimeout = 60 * time.Second
	// MaxImageSize is the largest image, in bytes, accepted by the Image API.
	MaxImageSize = 500 << 10
)

// ErrImageTooLarge is wrapped by the error returned when an image to upload
// exceeds MaxImageSize.
var ErrImageTooLarge = errors.New("image is too large")

// SerpApiClient holds configuration settings and an HTTP client for making requests
type SerpApiClient struct {
	Setting    SerpApiClientSetting // Configuration settings for the client
	HttpSearch *http.Client         // HTTP client for making requests
}

// SerpApiClientSetting holds configuration settings for the SerpApiClient
type SerpApiClientSetting struct {
	Persistent          bool              // Enable persistent search (default: false)
	Asynchronous        bool              // Enable asynchronous search (default: false)
	Timeout             time.Duration     // Timeout for HTTP requests
	SerpApiKey          string            // SerpAPI Key for authentication
	Engine              string            // Search engine to use [default: "google"]
	Parameter           map[string]string // Additional default parameters for the search
	MaxIdleConnection   int               // Maximum number of idle connections to keep
	KeepAlive           time.Duration     // Time between keep-alive probes
	TLSHandshakeTimeout time.Duration     // Timeout for TLS handshake (default: 10 seconds)
}

// HTTPError describes a non-successful response from the SerpApi service.
type HTTPError struct {
	StatusCode int
	Status     string
	URL        string
	Body       string
}

// Error returns a human-readable description of the HTTP failure.
func (err *HTTPError) Error() string {
	if err.Body == "" {
		return fmt.Sprintf("serpapi request failed: %s", err.Status)
	}
	return fmt.Sprintf("serpapi request failed: %s: %s", err.Status, err.Body)
}

// NewSerpApiClientSetting initializes a new SerpApiClientSetting with default values
func NewSerpApiClientSetting(serpApiKey string) SerpApiClientSetting {
	return SerpApiClientSetting{
		Persistent:          false,
		Asynchronous:        false,
		Timeout:             DefaultTimeout, // Default timeout of 60 seconds
		Parameter:           make(map[string]string),
		SerpApiKey:          serpApiKey,
		Engine:              "google",         // Default search engine
		KeepAlive:           60 * time.Second, // Default TCP keep-alive interval
		TLSHandshakeTimeout: 10 * time.Second, // Default TLS handshake timeout
	}
}

// NewClient initializes a new SerpApiClient client
func NewClient(setting SerpApiClientSetting) SerpApiClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSHandshakeTimeout = setting.TLSHandshakeTimeout
	transport.DisableKeepAlives = !setting.Persistent
	transport.MaxIdleConns = setting.MaxIdleConnection
	transport.MaxIdleConnsPerHost = setting.MaxIdleConnection
	if setting.KeepAlive > 0 {
		transport.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: setting.KeepAlive,
		}).DialContext
	}
	httpSearch := &http.Client{
		Timeout:   setting.Timeout, // Use the timeout from the setting
		Transport: transport,
	}
	return SerpApiClient{Setting: setting, HttpSearch: httpSearch}
}

// Search returns search result as a map
func (client *SerpApiClient) Search(parameter map[string]string) (map[string]interface{}, error) {
	return client.SearchContext(context.Background(), parameter)
}

// SearchContext returns search results as a map and supports cancellation.
func (client *SerpApiClient) SearchContext(ctx context.Context, parameter map[string]string) (map[string]interface{}, error) {
	rsp, err := client.execute(ctx, "/search", "json", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// UploadImage uploads an image to SerpApi's Image API and returns the JSON
// response. Its image_id can be passed to engines such as google_lens.
// See https://serpapi.com/image-api for the supported formats.
//
// image is a file path (string) or an io.Reader. The uploaded filename is the
// base name of the path, or of the reader's Name() method when it has one (as
// *os.File does). Any other reader is uploaded as "image". Use
// UploadImageReader to choose the filename or to override the api_key.
//
// An image larger than MaxImageSize is not sent, and the returned error wraps
// ErrImageTooLarge. A reader is then left partly consumed.
func (client *SerpApiClient) UploadImage(image interface{}) (map[string]interface{}, error) {
	return client.UploadImageContext(context.Background(), image)
}

// UploadImageContext uploads an image like UploadImage and supports cancellation.
func (client *SerpApiClient) UploadImageContext(ctx context.Context, image interface{}) (map[string]interface{}, error) {
	switch value := image.(type) {
	case string:
		file, err := os.Open(value)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return client.UploadImageReaderContext(ctx, file, value, nil)
	case io.Reader:
		filename := ""
		if named, ok := value.(interface{ Name() string }); ok {
			filename = named.Name()
		}
		return client.UploadImageReaderContext(ctx, value, filename, nil)
	default:
		return nil, errors.New("image must be a file path or io.Reader")
	}
}

// UploadImageReader uploads the image read from r to SerpApi's Image API, like
// UploadImage, under the base name of filename (for example "photo.png"). An
// empty filename is uploaded as "image".
//
// parameter may set "api_key" to override the client's API key for this call.
// The key is resolved like Search: parameter, then Setting.Parameter, then
// Setting.SerpApiKey. The Image API takes no other parameter, so other entries
// are ignored. parameter may be nil.
//
// An image larger than MaxImageSize is not sent, and the returned error wraps
// ErrImageTooLarge. r is then left partly consumed.
func (client *SerpApiClient) UploadImageReader(r io.Reader, filename string, parameter map[string]string) (map[string]interface{}, error) {
	return client.UploadImageReaderContext(context.Background(), r, filename, parameter)
}

// UploadImageReaderContext uploads an image like UploadImageReader and supports cancellation.
func (client *SerpApiClient) UploadImageReaderContext(ctx context.Context, r io.Reader, filename string, parameter map[string]string) (map[string]interface{}, error) {
	if r == nil {
		return nil, errors.New("image must be a non-nil io.Reader")
	}
	rsp, err := client.executeImageUpload(ctx, r, filename, parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// executeImageUpload sends a multipart POST request to the Image API.
func (client *SerpApiClient) executeImageUpload(ctx context.Context, reader io.Reader, filename string, parameter map[string]string) (*http.Response, error) {
	if filename == "" {
		filename = "image"
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if apiKey, ok := client.apiKey(parameter); ok {
		if err := writer.WriteField("api_key", apiKey); err != nil {
			return nil, err
		}
	}
	if err := writer.WriteField("source", "go:"+VERSION); err != nil {
		return nil, err
	}
	part, err := writer.CreateFormFile("image", filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	// read one byte past the limit to detect an oversized image without buffering all of it
	size, err := io.Copy(part, io.LimitReader(reader, MaxImageSize+1))
	if err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}
	if size > MaxImageSize {
		return nil, fmt.Errorf("%w: the Image API limit is %d bytes", ErrImageTooLarge, MaxImageSize)
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	// The API key is in the body, so the endpoint is safe to show in errors.
	endpoint := BaseURL + "/image"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	rsp, err := client.HttpSearch.Do(request)
	if err != nil {
		return nil, err
	}
	return checkStatus(rsp, endpoint)
}

// Html returns raw HTML search result
func (client *SerpApiClient) Html(parameter map[string]string) (*string, error) {
	return client.HtmlContext(context.Background(), parameter)
}

// HtmlContext returns raw HTML search results and supports cancellation.
func (client *SerpApiClient) HtmlContext(ctx context.Context, parameter map[string]string) (*string, error) {
	rsp, err := client.execute(ctx, "/search", "html", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeText(rsp.Body)
}

// Md returns the search result as markdown optimized for LLMs and AI agents
func (client *SerpApiClient) Md(parameter map[string]string) (*string, error) {
	return client.MdContext(context.Background(), parameter)
}

// MdContext returns markdown search results and supports cancellation.
func (client *SerpApiClient) MdContext(ctx context.Context, parameter map[string]string) (*string, error) {
	rsp, err := client.execute(ctx, "/search", "md", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeText(rsp.Body)
}

// Location returns standardized location data
func (client *SerpApiClient) Location(location string, limit int) ([]interface{}, error) {
	return client.LocationContext(context.Background(), location, limit)
}

// LocationContext returns standardized location data and supports cancellation.
func (client *SerpApiClient) LocationContext(ctx context.Context, location string, limit int) ([]interface{}, error) {
	parameter := map[string]string{
		"q":     location,
		"limit": fmt.Sprint(limit),
	}
	rsp, err := client.execute(ctx, "/locations.json", "json", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSONArray(rsp.Body)
}

// Account returns account information
func (client *SerpApiClient) Account() (map[string]interface{}, error) {
	return client.AccountContext(context.Background())
}

// AccountContext returns account information and supports cancellation.
func (client *SerpApiClient) AccountContext(ctx context.Context) (map[string]interface{}, error) {
	rsp, err := client.execute(ctx, "/account", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// SearchArchive retrieves previous search results from the archive
func (client *SerpApiClient) SearchArchive(id string) (map[string]interface{}, error) {
	return client.SearchArchiveContext(context.Background(), id)
}

// SearchArchiveContext retrieves a previous search result and supports cancellation.
func (client *SerpApiClient) SearchArchiveContext(ctx context.Context, id string) (map[string]interface{}, error) {
	rsp, err := client.execute(ctx, "/searches/"+url.PathEscape(id)+".json", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// SearchArchiveMd retrieves previous search results formatted as Markdown
func (client *SerpApiClient) SearchArchiveMd(id string) (*string, error) {
	return client.SearchArchiveMdContext(context.Background(), id)
}

// SearchArchiveMdContext retrieves a previous search result as Markdown and supports cancellation.
func (client *SerpApiClient) SearchArchiveMdContext(ctx context.Context, id string) (*string, error) {
	rsp, err := client.execute(ctx, "/searches/"+url.PathEscape(id)+".md", "md", map[string]string{})
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeText(rsp.Body)
}

// decodeJSON decodes response body to a map
func (client *SerpApiClient) decodeJSON(body io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(body)
	var rsp map[string]interface{}
	if err := decoder.Decode(&rsp); err != nil {
		return nil, errors.New("failed to decode JSON")
	}
	if errorMessage, exists := rsp["error"].(string); exists {
		return nil, errors.New(errorMessage)
	}
	return rsp, nil
}

// decodeJSONArray decodes response body to a slice
func (client *SerpApiClient) decodeJSONArray(body io.Reader) ([]interface{}, error) {
	decoder := json.NewDecoder(body)
	var rsp []interface{}
	if err := decoder.Decode(&rsp); err != nil {
		return nil, errors.New("failed to decode JSON array")
	}
	return rsp, nil
}

// decodeText decodes response body to a raw string (HTML or markdown)
func (client *SerpApiClient) decodeText(body io.Reader) (*string, error) {
	buffer, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	text := string(buffer)
	return &text, nil
}

// apiKey returns the API key of a request and whether one is set. The key is
// taken from parameter, then from the client Setting.Parameter, then from
// Setting.SerpApiKey when it is not empty.
func (client *SerpApiClient) apiKey(parameter map[string]string) (string, bool) {
	if apiKey, ok := parameter["api_key"]; ok {
		return apiKey, true
	}
	if apiKey, ok := client.Setting.Parameter["api_key"]; ok {
		return apiKey, true
	}
	return client.Setting.SerpApiKey, client.Setting.SerpApiKey != ""
}

// execute sends an HTTP GET request and returns the response
func (client *SerpApiClient) execute(ctx context.Context, path string, output string, parameter map[string]string) (*http.Response, error) {
	query := url.Values{}
	for name, value := range parameter {
		query.Add(name, value)
	}
	for name, value := range client.Setting.Parameter {
		if _, ok := query[name]; !ok {
			query.Add(name, value)
		}
	}
	if apiKey, ok := client.apiKey(parameter); ok {
		query.Set("api_key", apiKey)
	}
	if _, ok := query["engine"]; !ok {
		if client.Setting.Engine != "" {
			query.Add("engine", client.Setting.Engine)
		}
	}
	if client.Setting.Asynchronous {
		query.Add("async", "true")
	}
	query.Add("source", "go:"+VERSION)
	query.Set("output", output)

	endpoint := BaseURL + path + "?" + query.Encode()
	// errors expose the URL, so never include the secret API key
	query.Del("api_key")
	redactedEndpoint := BaseURL + path + "?" + query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	rsp, err := client.HttpSearch.Do(request)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			urlErr.URL = redactedEndpoint
		}
		return nil, err
	}
	return checkStatus(rsp, redactedEndpoint)
}

// checkStatus returns rsp if it is successful, otherwise it closes the body
// and returns an *HTTPError for endpoint, which must not contain the API key.
func checkStatus(rsp *http.Response, endpoint string) (*http.Response, error) {
	if rsp.StatusCode >= http.StatusOK && rsp.StatusCode < http.StatusMultipleChoices {
		return rsp, nil
	}
	body, readErr := io.ReadAll(rsp.Body)
	closeErr := rsp.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return nil, &HTTPError{
		StatusCode: rsp.StatusCode,
		Status:     rsp.Status,
		URL:        endpoint,
		Body:       string(body),
	}
}
