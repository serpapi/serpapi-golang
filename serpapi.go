package serpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	VERSION        = "1.3.0"
	BaseURL        = "https://serpapi.com"
	DefaultTimeout = 60 * time.Second
)

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

// Markdown returns the search result as markdown optimized for LLMs and AI agents
func (client *SerpApiClient) Markdown(parameter map[string]string) (*string, error) {
	return client.MarkdownContext(context.Background(), parameter)
}

// MarkdownContext returns markdown search results and supports cancellation.
func (client *SerpApiClient) MarkdownContext(ctx context.Context, parameter map[string]string) (*string, error) {
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

// decodeJSON decodes response body to a map
func (client *SerpApiClient) decodeJSON(body io.ReadCloser) (map[string]interface{}, error) {
	defer body.Close()
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
func (client *SerpApiClient) decodeJSONArray(body io.ReadCloser) ([]interface{}, error) {
	defer body.Close()
	decoder := json.NewDecoder(body)
	var rsp []interface{}
	if err := decoder.Decode(&rsp); err != nil {
		return nil, errors.New("failed to decode JSON array")
	}
	return rsp, nil
}

// decodeText decodes response body to a raw string (HTML or markdown)
func (client *SerpApiClient) decodeText(body io.ReadCloser) (*string, error) {
	defer body.Close()
	buffer, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	text := string(buffer)
	return &text, nil
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
	if _, ok := query["api_key"]; !ok {
		if client.Setting.SerpApiKey != "" {
			query.Add("api_key", client.Setting.SerpApiKey)
		}
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
	query.Add("output", output)

	endpoint := BaseURL + path + "?" + query.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	rsp, err := client.HttpSearch.Do(request)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode < http.StatusOK || rsp.StatusCode >= http.StatusMultipleChoices {
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
	return rsp, nil
}
