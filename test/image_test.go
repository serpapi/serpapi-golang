package serpapi

import (
	"bytes"
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
				if value := request.FormValue("source"); value != "go:"+serpapi.VERSION {
					t.Errorf("source = %q, want go:%s", value, serpapi.VERSION)
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

	_, err = client.UploadImage(t.TempDir() + string(os.PathSeparator) + "missing.png")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
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

// uploadForm is the multipart form received by a fake Image API.
type uploadForm struct {
	sent     bool
	fields   map[string][]string
	filename string
	image    []byte
}

// uploadRecorder returns a client whose Image API records the request in form.
func uploadRecorder(t *testing.T, apiKey string, form *uploadForm) serpapi.SerpApiClient {
	t.Helper()
	client := imageClient(func(request *http.Request) (*http.Response, error) {
		form.sent = true
		if err := request.ParseMultipartForm(1 << 20); err != nil {
			t.Fatal(err)
		}
		form.fields = request.MultipartForm.Value
		file, header, err := request.FormFile("image")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		form.filename = header.Filename
		if form.image, err = io.ReadAll(file); err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"image_id":"image-123"}`)),
		}, nil
	})
	client.Setting.SerpApiKey = apiKey
	return client
}

func TestUploadImageReader(t *testing.T) {
	for name, test := range map[string]struct {
		filename string
		want     string
	}{
		"filename":  {"photo.png", "photo.png"},
		"directory": {"local" + string(os.PathSeparator) + "photo.png", "photo.png"},
		"empty":     {"", "image"},
	} {
		t.Run(name, func(t *testing.T) {
			var form uploadForm
			client := uploadRecorder(t, "client-api-key", &form)

			result, err := client.UploadImageReader(strings.NewReader("fake-image-data"), test.filename, nil)
			if err != nil {
				t.Fatal(err)
			}
			if result["image_id"] != "image-123" {
				t.Errorf("image_id = %v, want image-123", result["image_id"])
			}
			if form.filename != test.want {
				t.Errorf("filename = %q, want %q", form.filename, test.want)
			}
			if string(form.image) != "fake-image-data" {
				t.Errorf("image = %q, want fake-image-data", form.image)
			}
			if got := form.fields["api_key"]; len(got) != 1 || got[0] != "client-api-key" {
				t.Errorf("api_key = %q, want [client-api-key]", got)
			}
		})
	}
}

func TestUploadImageReaderOverridesAPIKey(t *testing.T) {
	var form uploadForm
	client := uploadRecorder(t, "client-api-key", &form)

	_, err := client.UploadImageReaderContext(context.Background(), strings.NewReader("fake-image-data"), "photo.png", map[string]string{
		"api_key": "call-api-key",
		"engine":  "google_lens",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := form.fields["api_key"]; len(got) != 1 || got[0] != "call-api-key" {
		t.Errorf("api_key = %q, want [call-api-key]", got)
	}
	// the Image API only takes the api_key and the source
	if _, ok := form.fields["engine"]; ok {
		t.Errorf("engine must not be sent to the Image API: %v", form.fields)
	}
	if got := form.fields["source"]; len(got) != 1 || got[0] != "go:"+serpapi.VERSION {
		t.Errorf("source = %q, want [go:%s]", got, serpapi.VERSION)
	}
	if client.Setting.SerpApiKey != "client-api-key" {
		t.Errorf("client api key was modified: %q", client.Setting.SerpApiKey)
	}
}

func TestUploadImageAPIKeyPrecedence(t *testing.T) {
	for name, test := range map[string]struct {
		serpApiKey string
		setting    map[string]string
		parameter  map[string]string
		want       string
	}{
		"setting parameter only":      {"", map[string]string{"api_key": "setting-api-key"}, nil, "setting-api-key"},
		"setting parameter over key":  {"client-api-key", map[string]string{"api_key": "setting-api-key"}, nil, "setting-api-key"},
		"call parameter over setting": {"client-api-key", map[string]string{"api_key": "setting-api-key"}, map[string]string{"api_key": "call-api-key"}, "call-api-key"},
	} {
		t.Run(name, func(t *testing.T) {
			var form uploadForm
			client := uploadRecorder(t, test.serpApiKey, &form)
			client.Setting.Parameter = test.setting

			if _, err := client.UploadImageReader(strings.NewReader("fake-image-data"), "photo.png", test.parameter); err != nil {
				t.Fatal(err)
			}
			if got := form.fields["api_key"]; len(got) != 1 || got[0] != test.want {
				t.Errorf("api_key = %q, want [%s]", got, test.want)
			}
		})
	}
}

func TestUploadImageReaderRejectsNil(t *testing.T) {
	var form uploadForm
	client := uploadRecorder(t, "client-api-key", &form)

	_, err := client.UploadImageReader(nil, "photo.png", nil)
	if err == nil || err.Error() != "image must be a non-nil io.Reader" {
		t.Errorf("error = %v", err)
	}
	if form.sent {
		t.Error("the request must not be sent without an image")
	}
}

func TestUploadImageWithoutAPIKey(t *testing.T) {
	var form uploadForm
	client := uploadRecorder(t, "", &form)

	if _, err := client.UploadImage(strings.NewReader("fake-image-data")); err != nil {
		t.Fatal(err)
	}
	if _, ok := form.fields["api_key"]; ok {
		t.Errorf("api_key must not be sent when empty: %v", form.fields)
	}
	if got := form.fields["source"]; len(got) != 1 || got[0] != "go:"+serpapi.VERSION {
		t.Errorf("source = %q, want [go:%s]", got, serpapi.VERSION)
	}
}

func TestUploadImageSizeLimit(t *testing.T) {
	var form uploadForm
	client := uploadRecorder(t, "client-api-key", &form)

	_, err := client.UploadImage(bytes.NewReader(make([]byte, serpapi.MaxImageSize+1)))
	if !errors.Is(err, serpapi.ErrImageTooLarge) {
		t.Fatalf("expected serpapi.ErrImageTooLarge, got %v", err)
	}
	if want := "image is too large: the Image API limit is 512000 bytes"; err.Error() != want {
		t.Errorf("error = %q, want %q", err, want)
	}
	if form.sent {
		t.Error("an oversized image must not be sent")
	}

	// an image of exactly MaxImageSize bytes is accepted
	if _, err := client.UploadImage(bytes.NewReader(make([]byte, serpapi.MaxImageSize))); err != nil {
		t.Fatal(err)
	}
	if len(form.image) != serpapi.MaxImageSize {
		t.Errorf("uploaded %d bytes, want %d", len(form.image), serpapi.MaxImageSize)
	}
}

// imageFixture is a 14 KB PNG, well under the Image API limit.
const imageFixture = "testdata/serpapi-logo.png"

type failingReader struct{}

var errDiskRead = errors.New("disk read failed")

func (failingReader) Read([]byte) (int, error) {
	return 0, errDiskRead
}

func TestUploadImageReadError(t *testing.T) {
	var form uploadForm
	client := uploadRecorder(t, "client-api-key", &form)

	_, err := client.UploadImageReader(failingReader{}, "photo.png", nil)
	if !errors.Is(err, errDiskRead) {
		t.Fatalf("expected the read error, got %v", err)
	}
	if form.sent {
		t.Error("the request must not be sent after a read error")
	}
}

// doc: https://serpapi.com/image-api
// The Image API is free: this test spends no search credit.
func TestUploadImageLive(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}
	client := serpapi.NewClient(serpapi.NewSerpApiClientSetting(getApiKey()))

	upload, err := client.UploadImage(imageFixture)
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := upload["image_id"].(string); id == "" {
		t.Fatalf("expected a non-empty image_id, got %v", upload)
	}

	// the api_key can be given per call instead of on the client
	image, err := os.Open(imageFixture)
	if err != nil {
		t.Fatal(err)
	}
	defer image.Close()
	keyless := serpapi.NewClient(serpapi.NewSerpApiClientSetting(""))
	readerUpload, err := keyless.UploadImageReader(image, "serpapi-logo.png", map[string]string{"api_key": getApiKey()})
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := readerUpload["image_id"].(string); id == "" {
		t.Fatalf("expected a non-empty image_id, got %v", readerUpload)
	}
}

// doc: https://serpapi.com/google-lens-api
// This test spends one search credit.
func TestUploadImageGoogleLensLive(t *testing.T) {
	if shoulSkip() {
		t.Skip("SERPAPI_KEY required")
		return
	}
	client := serpapi.NewClient(serpapi.NewSerpApiClientSetting(getApiKey()))

	upload, err := client.UploadImage(imageFixture)
	if err != nil {
		t.Fatal(err)
	}
	imageID, _ := upload["image_id"].(string)
	if imageID == "" {
		t.Fatalf("expected a non-empty image_id, got %v", upload)
	}

	results, err := client.Search(map[string]string{
		"engine":   "google_lens",
		"image_id": imageID,
	})
	if err != nil {
		t.Fatal(err)
	}
	metadata, _ := results["search_metadata"].(map[string]interface{})
	t.Logf("google_lens search_metadata.id = %v", metadata["id"])
	if metadata["status"] != "Success" {
		t.Errorf("search_metadata.status = %v, want Success", metadata["status"])
	}
	parameters, _ := results["search_parameters"].(map[string]interface{})
	if parameters["engine"] != "google_lens" {
		t.Errorf("search_parameters.engine = %v, want google_lens", parameters["engine"])
	}
}
