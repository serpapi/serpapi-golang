package serpapi

import (
	"testing"
	"time"

	"github.com/serpapi/serpapi-golang"
)

func TestQuickStart(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}
	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	setting.Persistent = false // Close the HTTP connection after the request to avoid keeping it open
	setting.Asynchronous = false // Block search query until results are returned
	setting.Timeout = 30 * time.Second
	setting.Engine = "google" // Set the search engine to Google
	// set default parameters
	setting.Parameter = map[string]string{
		"start": "0",
		"hl":    "en",
		"google_domain": "google.com",
	}
	client := serpapi.NewClient(setting)

	parameter := map[string]string{
		"q":             "Coffee",
		"location":      "Portland, Oregon, United States",
		"hl":            "en",
		"gl":            "us",
		"safe":          "active",
		"start":         "10",
		"num":           "10",
		"device":        "desktop",
		"engine":        "google",  // override engine default
	}
	results, err := client.Search(parameter)

	if err != nil {
		t.Error(err)
		return
	}
	result := results["organic_results"].([]interface{})[0].(map[string]interface{})
	if len(result["title"].(string)) == 0 {
		t.Error("empty title in local results")
		return
	}
}
