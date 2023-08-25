package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// basic use case
func TestGoogle(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  client_parameter := map[string]string{
    "engine": "google",
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(client_parameter)

  parameter := map[string]string{ 
    "q": "coffee", 
    "engine": "google",
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
