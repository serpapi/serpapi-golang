package serpapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type contextTransport struct{}

func (contextTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	<-request.Context().Done()
	return nil, request.Context().Err()
}

func TestSearchContextCancellation(t *testing.T) {
	client := NewClient(NewSerpApiClientSetting(""))
	client.HttpSearch = &http.Client{Transport: contextTransport{}}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.SearchContext(ctx, map[string]string{"q": "coffee"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

type responseTransport struct{}

func (responseTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&stringReader{value: `{"ok":true}`}),
		Header:     make(http.Header),
	}, nil
}

type stringReader struct {
	value string
	read  bool
}

func (reader *stringReader) Read(buffer []byte) (int, error) {
	if reader.read {
		return 0, io.EOF
	}
	reader.read = true
	return copy(buffer, reader.value), nil
}

func TestSearchRemainsUsableWithoutContext(t *testing.T) {
	client := NewClient(NewSerpApiClientSetting(""))
	client.HttpSearch = &http.Client{Transport: responseTransport{}}

	result, err := client.Search(map[string]string{"q": "coffee"})
	if err != nil {
		t.Fatalf("Search returned an error: %v", err)
	}
	if result["ok"] != true {
		t.Fatalf("expected ok response, got %#v", result)
	}
}

type errorTransport struct{}

func (errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Status:     "429 Too Many Requests",
		Body:       io.NopCloser(&stringReader{value: `{"error":"rate limit exceeded"}`}),
		Header:     make(http.Header),
	}, nil
}

func TestSearchReturnsHTTPError(t *testing.T) {
	client := NewClient(NewSerpApiClientSetting(""))
	client.HttpSearch = &http.Client{Transport: errorTransport{}}

	_, err := client.Search(map[string]string{"q": "coffee"})
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %v", err)
	}
	if httpErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status code %d, got %d", http.StatusTooManyRequests, httpErr.StatusCode)
	}
	if httpErr.Body != `{"error":"rate limit exceeded"}` {
		t.Fatalf("unexpected response body: %q", httpErr.Body)
	}
	if httpErr.URL == "" {
		t.Fatal("expected request URL in HTTPError")
	}
	if got := httpErr.Error(); got != fmt.Sprintf("serpapi request failed: 429 Too Many Requests: %s", httpErr.Body) {
		t.Fatalf("unexpected error string: %q", got)
	}
}
