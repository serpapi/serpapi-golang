
package serpapi

import (
  "testing"
)

// basic use case
func TestYoutube(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  client_parameter := map[string]string{
    "engine": "youtube",
    "api_key": *getApiKey(),
  }
  client := NewClient(client_parameter)

  parameter := map[string]string{
    "search_query": "coffee",
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

  if len(rsp["video_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 video_results")
    return
  }
}  
