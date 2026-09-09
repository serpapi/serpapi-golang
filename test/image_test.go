package serpapi

import (
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
	imagePath := t.TempDir() + "/image.png"
	if err := os.WriteFile(imagePath, []byte("fake-image-data"), 0600); err != nil {
		t.Fatal(err)
	}

	for name, image := range map[string]interface{}{
		"path":   imagePath,
		"reader": strings.NewReader("fake-image-data"),
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
				file, _, err := request.FormFile("image")
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()
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

			result, err := client.UploadImage(image)
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
