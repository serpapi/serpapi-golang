package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_reverse_image engine
// doc: https://serpapi.com/google-reverse-image
//
func TestGoogleReverseImage(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_reverse_image", 
    "image_url": "https://i.imgur.com/5bGzZi7.jpg",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["image_sizes"] == nil {
    t.Error("key is not found: image_sizes")
    return 
  }

  if len(rsp["image_sizes"].([]interface{})) < 1 {
    t.Error("expect more than 1 image_sizes") 
    return
  }
}  
