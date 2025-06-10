package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

// Search archive API
// doc: https://serpapi.com/search-archive-api
func TestSearchArchive(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}

	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	setting.Engine = "google" // Set the search engine to Google
	client := serpapi.NewClient(setting)
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	scheduleSearch, err := client.Search(parameter)

	if err != nil {
		t.Error("unexpected error", err)
		return
	}
	searchMetadata, ok := scheduleSearch["search_metadata"].(map[string]interface{})
	if !ok {
		t.Error("search_metadata is missing or invalid")
		return
	}

	searchID, ok := searchMetadata["id"].(string)
	if !ok || len(searchID) == 0 {
		t.Error("search_metadata.id is missing or invalid")
		return
	}

	searchArchive, err := client.SearchArchive(searchID)
	if err != nil {
		t.Error(err)
		return
	}

	searchMetadataArchive, ok := searchArchive["search_metadata"].(map[string]interface{})
	if !ok {
		t.Error("search_metadata in search archive is missing or invalid")
		return
	}

	searchIDArchive, ok := searchMetadataArchive["id"].(string)
	if !ok || searchIDArchive != searchID {
		t.Errorf("search_metadata.id mismatch: got %v, expected %v", searchIDArchive, searchID)
	}
	// print search results from search archive
}
