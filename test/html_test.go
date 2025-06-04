package serpapi

import (
	"strings"
	"testing"

	"github.com/serpapi/serpapi-golang"
)

func TestHtml(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	setting = serpapi.NewSerpApiClientSetting(*getApiKey())
	setting.Engine = "google" // Set the search engine to Google
	client := serpapi.NewClient(client_parameter)

	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	data, err := client.Html(parameter)
	if err != nil {
		t.Error("err must be nil")
		return
	}
	if !strings.Contains(*data, "</html>") {
		t.Error("data does not contains <html> tag")
	}
}
