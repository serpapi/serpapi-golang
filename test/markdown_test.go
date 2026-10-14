package serpapi

import (
	"strings"
	"testing"

	"github.com/serpapi/serpapi-golang"
)

func TestMarkdown(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}

	// Initialize the SerpApi client with the API key
	// and set the search engine to Google
	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	setting.Engine = "google" // Set the search engine to Google
	client := serpapi.NewClient(setting)

	// Define the search parameters
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	// Perform the search and get the markdown response
	data, err := client.Markdown(parameter)
	if err != nil {
		t.Error("err must be nil")
		return
	}
	if !strings.Contains(*data, "#") {
		t.Error("data does not contain any markdown heading")
	}
}
