package serpapi

import (
	"testing"
)

// basic use case
func TestBing(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	client_parameter := map[string]string{
		"engine":  "bing",
		"api_key": *getApiKey(),
	}
	client := NewClient(client_parameter)

	parameter := map[string]string{
		"q": "coffee",
	}
	rsp, err := client.Search(parameter)

	if err != nil {
		t.Error("unexpected error", err)
		return
	}

	if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
		t.Error("bad status")
		return
	}

	if len(rsp["organic_results"].([]interface{})) < 1 {
		t.Error("expect at least one organic_results")
		return
	}
}
