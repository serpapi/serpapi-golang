package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_maps_photos engine
// doc: http://serpapi.com/google-maps-photos-api
//
func TestGoogleMapsPhotos(t *testing.T) {
  // skip test if no SERPAPI_KEY provided
  //  and account is required. see: https://serpapi.com/
  api_key := os.Getenv("SERPAPI_KEY")
  if api_key == "" {
    t.Skip("SERPAPI_KEY required")
    return
  }

  // Initialize the client with custom setting
	setting := serpapi.NewSerpApiClientSetting(api_key) // Replace with your SerpApi key
	setting.Persistent = false                     // Enable persistent search
	setting.Asynchronous = false                   // Enable asynchronous search
	setting.Timeout = 60 * time.Second             // Set timeout for HTTP requests
	setting.MaxIdleConnection = 10                 // Set maximum idle connections
	setting.KeepAlive = 60 * time.Second           // Set keep-alive duration
	setting.TLSHandshakeTimeout = 10 * time.Second // Set TLS handshake timeout

  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_maps_photos", 
    "data_id": "0x89c25998e556791b:0xbdcd67f46e37b16d",
  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error: ", err)
    return
  }

  status := rsp["search_metadata"].(map[string]interface{})["status"]
  if status != "Success" {
    t.Error("unexpected status: ", status)
    return
  }

  photos, ok := rsp["photos"].([]interface{})
  if !ok || len(photos) < 1 {
    t.Error("expect non-empty photos")
    t.Errorf("results: %v", rsp["photos"])
    return
  }
}  
