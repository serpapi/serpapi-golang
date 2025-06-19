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
	// run a non blocking search
	setting.Asynchronous = true
	// set default search engine
	setting.Engine = "google"
	client := serpapi.NewClient(setting)
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	scheduleSearch, err := client.Search(parameter)

	if err != nil {
		t.Error("unexpected error", err)
		return
	}
	metadata, ok := scheduleSearch["search_metadata"].(map[string]interface{})
	if !ok {
		t.Error("search_metadata is missing or invalid")
		return
	}

	searchID, ok := metadata["id"].(string)
	if !ok || len(searchID) == 0 {
		t.Error("search_metadata.id is missing or invalid")
		return
	}

	// wait for the search to complete
	t.Logf("Waiting for search %s to complete...", searchID)
	var archiveMetadata map[string]interface{}
	for {
		searchArchive, err := client.SearchArchive(searchID)
		if err != nil {
			t.Error(err)
			return
		}

		archiveMetadata, ok = searchArchive["search_metadata"].(map[string]interface{})
		if !ok {
			t.Error("search_metadata in search archive is missing or invalid")
			return
		}

		status, ok := archiveMetadata["status"].(string)
		if !ok {
			t.Error("search_metadata.status is missing or invalid")
			return
		}

		if status == "Cached" || status == "Success" {
			break
		}
	}

	searchIDArchive, ok := archiveMetadata["id"].(string)
	if !ok || searchIDArchive != searchID {
		t.Errorf("search_metadata.id mismatch: got %v, expected %v", searchIDArchive, searchID)
	}
	// print search results from search archive
}
