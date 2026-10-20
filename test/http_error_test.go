package serpapi

import (
	"strings"
	"testing"

	"github.com/serpapi/serpapi-golang"
)

func TestHTTPErrorMessage(t *testing.T) {
	client := serpapi.NewClient(serpapi.NewSerpApiClientSetting("invalid_api_key"))

	_, err := client.Search(map[string]string{"q": "Coffee"})
	if err == nil {
		t.Fatal("expected an error for an invalid API key")
	}

	message := err.Error()
	if !strings.HasPrefix(message, "serpapi request failed: 401 Unauthorized") {
		t.Errorf("unexpected message prefix: %q", message)
	}
	if !strings.Contains(message, "Invalid API key") {
		t.Errorf("message does not include the response body: %q", message)
	}
	if strings.Contains(message, "invalid_api_key") {
		t.Error("message leaks the API key")
	}
}
