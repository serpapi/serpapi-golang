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

	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	setting.Engine = "google"
	client := serpapi.NewClient(setting)

	markdown, err := client.Markdown(map[string]string{
		"q":        "Coffee",
		"location": "Austin, Texas, United States",
	})
	if err != nil {
		t.Fatalf("Markdown returned an error: %v", err)
	}
	if !strings.HasPrefix(*markdown, "---") {
		t.Error("Markdown response does not contain YAML frontmatter")
	}
	if !strings.Contains(*markdown, "## Organic Results") {
		t.Error("Markdown response does not contain organic results")
	}
}
