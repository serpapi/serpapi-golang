package serpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// Current version
	VERSION = "1.0.0"
)

// Search holds query
type SerpApiClient struct {
	Parameter  map[string]string
	HttpSearch *http.Client
}

// NewClient initiliaze a new SerpApiClient client
func NewClient(parameter map[string]string) SerpApiClient {
	var timeout time.Duration
	var err error
	timeout = time.Second * 60
	if v, ok := parameter["timeout"]; ok {
		timeout, err = time.ParseDuration(v + "s")
		if err != nil {
			log.Fatal("timeout must be an integer value indicating the number of seconds until HTTP client timeout")
			panic(err)
		}
	}
	// Create the http search
	httpSearch := &http.Client{
		Timeout: timeout,
	}
	return SerpApiClient{Parameter: parameter, HttpSearch: httpSearch}
}

// Search returns search result as a Map
func (client *SerpApiClient) Search(parameter map[string]string) (map[string]interface{}, error) {
	rsp, err := client.execute("/search", "json", parameter)
	if err != nil {
		return nil, err
	}
	return client.decodeJSON(rsp.Body)
}

// Html returns raw html search result
func (client *SerpApiClient) Html(parameter map[string]string) (*string, error) {
	rsp, err := client.execute("/search", "html", parameter)
	if err != nil {
		return nil, err
	}
	return client.decodeHTML(rsp.Body)
}

// Location returns the standardize set of location takes location.
func (client *SerpApiClient) Location(location string, limit int) ([]interface{}, error) {
	client.Parameter = map[string]string{
		"q":     location,
		"limit": fmt.Sprint(limit),
	}
	rsp, err := client.execute("/locations.json", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	return client.decodeJSONArray(rsp.Body)
}

// Account return account information
func (client *SerpApiClient) Account() (map[string]interface{}, error) {
	rsp, err := client.execute("/account", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	return client.decodeJSON(rsp.Body)
}

// SearchArchive retrieves previous search results from the SerpApiClient archive (no credit charge)
func (client *SerpApiClient) SearchArchive(id string) (map[string]interface{}, error) {
	rsp, err := client.execute("/searches/"+id+".json", "json", map[string]string{})
	if err != nil {
		return nil, err
	}
	return client.decodeJSON(rsp.Body)
}

// decodeJson response
func (client *SerpApiClient) decodeJSON(body io.ReadCloser) (map[string]interface{}, error) {
	// Decode JSON from response body
	decoder := json.NewDecoder(body)

	// Response data
	var rsp map[string]interface{}
	err := decoder.Decode(&rsp)
	if err != nil {
		return nil, errors.New("fail to decode")
	}

	// check error message
	errorMessage, derror := rsp["error"].(string)
	if derror {
		return nil, errors.New(errorMessage)
	}
	return rsp, nil
}

// decodeJSONArray decodes response body to SearchResultArray
func (client *SerpApiClient) decodeJSONArray(body io.ReadCloser) ([]interface{}, error) {
	decoder := json.NewDecoder(body)
	var rsp []interface{}
	err := decoder.Decode(&rsp)
	if err != nil {
		return nil, errors.New("fail to decode array")
	}
	return rsp, nil
}

// decodeHTML decodes response body to html string
func (client *SerpApiClient) decodeHTML(body io.ReadCloser) (*string, error) {
	buffer, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	text := string(buffer)
	return &text, nil
}

// execute HTTP get reuqest and returns http response
func (client *SerpApiClient) execute(path string, output string, parameter map[string]string) (*http.Response, error) {
	query := url.Values{}
	for k, v := range client.Parameter {
		query.Add(k, v)
	}
	for k, v := range parameter {
		query.Add(k, v)
	}

	// source programming language
	query.Add("source", "go")

	// set output
	query.Add("output", output)

	endpoint := "https://serpapi.com" + path + "?" + query.Encode()
	rsp, err := client.HttpSearch.Get(endpoint)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}
