package serpapi

import (
	"os"
	"testing"
	"time"

	"github.com/serpapi/serpapi-golang"
)

// example test for google_light_search engine
// doc: https://serpapi.com/google_light_search
func TestGoogleLightSearch(t *testing.T) {
	api_key := os.Getenv("SERPAPI_KEY")
	if api_key == "" {
		t.Skip("SERPAPI_KEY required")
		return
	}

	setting := serpapi.NewSerpApiClientSetting(api_key)
	setting.Timeout = 60 * time.Second
	client := serpapi.NewClient(setting)

	parameter := map[string]string{
		"engine": "google_light_search",
		"q":      "coffee",
	}
	rsp, err := client.Search(parameter)
	if err != nil {
		t.Error("unexpected error: ", err)
		return
	}

	status := rsp["search_metadata"].(map[string]interface{})["status"]
	if status != "Success" {
		t.Error("unexpected status: ", status)
		return
	}

	organic_results, ok := rsp["organic_results"].([]interface{})
	if !ok || len(organic_results) < 1 {
		t.Error("expect non-empty organic_results")
		t.Errorf("results: %v", rsp["organic_results"])
		return
	}
}
