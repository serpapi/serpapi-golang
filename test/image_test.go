package serpapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/serpapi/serpapi-golang"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func imageClient(handler roundTripFunc) serpapi.SerpApiClient {
	setting := serpapi.NewSerpApiClientSetting("client-api-key")
	client := serpapi.NewClient(setting)
	client.HttpSearch = &http.Client{Transport: handler}
	return client
}

func TestUploadImage(t *testing.T) {
	imagePath := t.TempDir() + string(os.PathSeparator) + "image.png"
	if err := os.WriteFile(imagePath, []byte("fake-image-data"), 0600); err != nil {
		t.Fatal(err)
	}

	imageFile, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer imageFile.Close()

	for name, test := range map[string]struct {
		image    interface{}
		filename string
	}{
		"path":   {imagePath, "image.png"},
		"file":   {imageFile, "image.png"},
		"reader": {strings.NewReader("fake-image-data"), "image"},
	} {
		t.Run(name, func(t *testing.T) {
			client := imageClient(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodPost || request.URL.String() != "https://serpapi.com/image" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				if err := request.ParseMultipartForm(1 << 20); err != nil {
					t.Fatal(err)
				}
				if value := request.FormValue("api_key"); value != "client-api-key" {
					t.Errorf("api_key = %q, want client-api-key", value)
				}
				file, header, err := request.FormFile("image")
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()
				if header.Filename != test.filename {
					t.Errorf("filename = %q, want %q", header.Filename, test.filename)
				}
				contents, err := io.ReadAll(file)
				if err != nil {
					t.Fatal(err)
				}
				if string(contents) != "fake-image-data" {
					t.Errorf("image = %q, want fake-image-data", contents)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"image_id":"image-123"}`)),
				}, nil
			})

			result, err := client.UploadImage(test.image)
			if err != nil {
				t.Fatal(err)
			}
			if result["image_id"] != "image-123" {
				t.Errorf("image_id = %v, want image-123", result["image_id"])
			}
		})
	}
}

func TestUploadImageRejectsUnsupportedInput(t *testing.T) {
	client := imageClient(func(request *http.Request) (*http.Response, error) {
		t.Fatal("request should not be sent")
		return nil, nil
	})

	_, err := client.UploadImage(42)
	if err == nil || err.Error() != "image must be a file path or io.Reader" {
		t.Errorf("error = %v", err)
	}
}

func TestUploadImageReturnsHTTPError(t *testing.T) {
	const body = `{"error":"Invalid API key."}`
	client := imageClient(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Status:     "401 Unauthorized",
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})

	_, err := client.UploadImage(strings.NewReader("fake-image-data"))
	var httpErr *serpapi.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *serpapi.HTTPError, got %v", err)
	}
	if httpErr.StatusCode != http.StatusUnauthorized || httpErr.Body != body {
		t.Errorf("HTTPError = %d %q", httpErr.StatusCode, httpErr.Body)
	}
	if httpErr.URL != "https://serpapi.com/image" {
		t.Errorf("URL = %q, want https://serpapi.com/image", httpErr.URL)
	}
	if strings.Contains(err.Error(), "client-api-key") {
		t.Errorf("error leaks the API key: %v", err)
	}
}

func TestUploadImageContextCancellation(t *testing.T) {
	client := imageClient(func(request *http.Request) (*http.Response, error) {
		<-request.Context().Done()
		return nil, request.Context().Err()
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.UploadImageContext(ctx, strings.NewReader("fake-image-data"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
