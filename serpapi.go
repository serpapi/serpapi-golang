package serpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	VERSION       = "1.0.0"
	BaseURL       = "https://serpapi.com"
	DefaultTimeout = 60 * time.Second
)

// SerpApiClient holds query parameters and HTTP client
type SerpApiClient struct {
	Parameter  map[string]string
	HttpSearch *http.Client
}

// NewClient initializes a new SerpApiClient client
func NewClient(parameter map[string]string) SerpApiClient {
	timeout := DefaultTimeout
	if v, ok := parameter["timeout"]; ok {
		parsedTimeout, err := time.ParseDuration(v + "s")
		if err != nil {
			log.Fatalf("Invalid timeout value: %v", err)
		}
		timeout = parsedTimeout

		delete(parameter, "timeout")
	}

	httpSearch := &http.Client{Timeout: timeout}
	return SerpApiClient{Parameter: parameter, HttpSearch: httpSearch}
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
	return client.decodeHTML(rsp.Body)
}

// Location returns standardized location data
func (client *SerpApiClient) Location(location string, limit int) ([]interface{}, error) {
	client.Parameter = map[string]string{
		"q":     location,
		"limit": fmt.Sprint(limit),
	}
	rsp, err := client.execute("/locations.json", "json", map[string]string{})
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

// decodeHTML decodes response body to an HTML string
func (client *SerpApiClient) decodeHTML(body io.ReadCloser) (*string, error) {
	defer body.Close()
	buffer, err := io.ReadAll(body)
	if err != nil {
		return nil, err
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
	for name, value := range client.Parameter {
		if _, exists := parameter[name]; !exists {
			query.Add(name, value)
		}
	}

	query.Add("source", "go")
	query.Add("output", output)

	endpoint := BaseURL + path + "?" + query.Encode()
	rsp, err := client.HttpSearch.Get(endpoint)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
