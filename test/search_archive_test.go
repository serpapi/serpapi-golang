package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

// Search archive API
// doc: https://serpapi.com/search-archive-api
func TestSearchArchive(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	setting := serpapi.NewSerpApiClientSetting(*getApiKey())
	setting.Engine = "google" // Set the search engine to Google
	client := serpapi.NewClient(setting)
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	rsp, err := client.Search(parameter)

	if err != nil {
		t.Error("unexpected error", err)
		return
	}

	searchID := rsp["search_metadata"].(map[string]interface{})["id"].(string)
	if len(searchID) == 0 {
		t.Error("search_metadata.id must be defined")
		return
	}

	searchArchive, err := client.SearchArchive(searchID)
	if err != nil {
		t.Error(err)
		return
	}

	searchIDArchive := searchArchive["search_metadata"].(map[string]interface{})["id"].(string)
	if searchIDArchive != searchID {
		t.Error("search_metadata.id do not match", searchIDArchive, searchID)
	}
}
