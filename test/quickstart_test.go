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
	setting.Timeout = 30 * time.Second
	setting.Engine = "google" // Set the search engine to Google

	client := serpapi.NewClient(setting)

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
