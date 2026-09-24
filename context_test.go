package serpapi

import (
	"context"
	"errors"
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
