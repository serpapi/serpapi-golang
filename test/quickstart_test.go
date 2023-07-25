package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

func TestQuickStart(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	client_parameter := map[string]string{
		"timeout": "30",
		"api_key": *getApiKey(),
		"engine":  "google",
	}
	client := serpapi.NewClient(client_parameter)

	parameter := map[string]string{
		"q":             "Coffee",
		"location":      "Portland, Oregon, United States",
		"hl":            "en",
		"gl":            "us",
		"google_domain": "google.com",
		"safe":          "active",
		"start":         "10",
		"num":           "10",
		"device":        "desktop",
	}
	rsp, err := client.Search(parameter)

	if err != nil {
		t.Error(err)
		return
	}
	result := rsp["organic_results"].([]interface{})[0].(map[string]interface{})
	if len(result["title"].(string)) == 0 {
		t.Error("empty title in local results")
		return
	}
}
