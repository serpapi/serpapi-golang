
package serpapi

import (
  "testing"
)

// basic use case
func TestGoogleReverseImage(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  client_parameter := map[string]string{
    "engine": "google_reverse_image",
    "api_key": *getApiKey(),
  }
  client := NewClient(client_parameter)

  parameter := map[string]string{
    "image_url": "https://i.imgur.com/5bGzZi7.jpg",
  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if len(rsp["image_sizes"].([]interface{})) < 1 {
    t.Error("expect more than 5 image_sizes")
    return
  }
}  
