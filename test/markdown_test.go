package serpapi

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/serpapi/serpapi-golang"
)

func TestMd(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}

	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	setting.Engine = "google"
	client := serpapi.NewClient(setting)

	markdown, err := client.Md(map[string]string{
		"q":        "Coffee",
		"location": "Austin, Texas, United States",
	})
	if err != nil {
		t.Fatalf("Md returned an error: %v", err)
	}
	if !strings.HasPrefix(*markdown, "---") {
		t.Error("Markdown response does not contain YAML frontmatter")
	}
	if !strings.Contains(*markdown, "## Organic Results") {
		t.Error("Markdown response does not contain organic results")
	}
}

func TestMarkdownContext(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}

	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	client := serpapi.NewClient(setting)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	markdown, err := client.MarkdownContext(ctx, map[string]string{
		"q":        "Coffee",
		"location": "Portland",
	})
	if err != nil {
		t.Fatalf("MarkdownContext returned an error: %v", err)
	}
	if !strings.HasPrefix(*markdown, "---") {
		t.Error("Markdown response does not contain YAML frontmatter")
	}
}

func TestMdInvalidAPIKey(t *testing.T) {
	client := serpapi.NewClient(serpapi.NewSerpApiClientSetting("invalid_api_key"))

	_, err := client.Md(map[string]string{"q": "Coffee"})

	var httpErr *serpapi.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *serpapi.HTTPError, got %v", err)
	}
	if httpErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", httpErr.StatusCode)
	}
	if strings.Contains(httpErr.URL, "invalid_api_key") {
		t.Error("HTTPError URL leaks the API key")
	}
}
