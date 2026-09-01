package serpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	VERSION        = "1.2.0"
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
	KeepAlive           time.Duration     // Time between keep-alive probes (default: 60 seconds)
	TLSHandshakeTimeout time.Duration     // Timeout for TLS handshake (default: 10 seconds)
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
		TLSHandshakeTimeout: 10 * time.Second, // Default TLS handshake timeout
	}
}

// NewClient initializes a new SerpApiClient client
func NewClient(setting SerpApiClientSetting) SerpApiClient {
	transport := &http.Transport{
		TLSHandshakeTimeout: setting.TLSHandshakeTimeout, // Adjust as needed for HTTPS
		DisableKeepAlives:   !setting.Persistent,         // Keep-alives are enabled by default in Go >=1.5
		MaxIdleConns:        setting.MaxIdleConnection,   // Set maximum idle connections
	}
	httpSearch := &http.Client{
		Timeout:   setting.Timeout, // Use the timeout from the setting
		Transport: transport,
	}
	return SerpApiClient{Setting: setting, HttpSearch: httpSearch}
}

// Search returns search result as a map
func (client *SerpApiClient) Search(parameter map[string]string) (map[string]interface{}, error) {
	rsp, err := client.execute("/search", "json", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// Html returns raw HTML search result
func (client *SerpApiClient) Html(parameter map[string]string) (*string, error) {
	rsp, err := client.execute("/search", "html", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeText(rsp.Body)
}

// Markdown returns search results formatted as Markdown
func (client *SerpApiClient) Markdown(parameter map[string]string) (*string, error) {
	rsp, err := client.execute("/search", "md", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeTextResponse(rsp)
}

// Location returns standardized location data
func (client *SerpApiClient) Location(location string, limit int) ([]interface{}, error) {
	parameter := map[string]string{
		"q":     location,
		"limit": fmt.Sprint(limit),
	}
	rsp, err := client.execute("/locations.json", "json", parameter)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSONArray(rsp.Body)
}

// Account returns account information
func (client *SerpApiClient) Account() (map[string]interface{}, error) {
	rsp, err := client.execute("/account", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// SearchArchive retrieves previous search results from the archive
func (client *SerpApiClient) SearchArchive(id string) (map[string]interface{}, error) {
	rsp, err := client.execute("/searches/"+id+".json", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeJSON(rsp.Body)
}

// SearchArchiveMarkdown retrieves previous search results formatted as Markdown
func (client *SerpApiClient) SearchArchiveMarkdown(id string) (*string, error) {
	rsp, err := client.execute("/searches/"+id+".md", "md", map[string]string{})
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return client.decodeTextResponse(rsp)
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

// decodeText decodes a response body to a string.
func (client *SerpApiClient) decodeText(body io.Reader) (*string, error) {
	buffer, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	text := string(buffer)
	return &text, nil
}

// decodeTextResponse decodes Markdown responses to a string. SerpApi returns
// API errors as JSON even when Markdown output was requested.
func (client *SerpApiClient) decodeTextResponse(rsp *http.Response) (*string, error) {
	buffer, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	mediaType, _, _ := mime.ParseMediaType(rsp.Header.Get("Content-Type"))
	mediaType = strings.ToLower(mediaType)
	isJSON := mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
	if isJSON {
		var data map[string]interface{}
		if err := json.Unmarshal(buffer, &data); err == nil {
			if errorMessage, exists := data["error"].(string); exists {
				return nil, errors.New(errorMessage)
			}
		}
	}

	if rsp.StatusCode < http.StatusOK || rsp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("HTTP request failed with status: %d", rsp.StatusCode)
	}
	if isJSON {
		return nil, errors.New("expected a Markdown response, received JSON")
	}

	text := string(buffer)
	return &text, nil
}

// execute sends an HTTP GET request and returns the response
func (client *SerpApiClient) execute(path string, output string, parameter map[string]string) (*http.Response, error) {
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
	query.Set("output", output)

	endpoint := BaseURL + path + "?" + query.Encode()
	rsp, err := client.HttpSearch.Get(endpoint)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
