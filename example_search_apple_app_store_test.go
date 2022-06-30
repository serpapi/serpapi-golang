
package serpapi

import (
  "testing"
)

// basic use case
func TestAppleAppStore(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  client_parameter := map[string]string{
    "engine": "apple_app_store",
    "api_key": *getApiKey(),
  }
  client := NewClient(client_parameter)

  parameter := map[string]string{
    "term": "coffee",
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

  if len(rsp["organic_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 organic_results")
    return
  }
}  