package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// basic use case
func TestGoogleAutocomplete(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_autocomplete", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["suggestions"] == nil {
    t.Error("key is not found: suggestions")
    return 
  }

  if len(rsp["suggestions"].([]interface{})) < 5 {
    t.Error("expect more than 5 suggestions") 
    return
  }
}  
