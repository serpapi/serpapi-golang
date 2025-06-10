package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_maps engine
// doc: https://serpapi.com/google-maps-api
//
func TestGoogleMaps(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_maps", 
    "q": "Coffee", 
    "ll": "@40.7455096,-74.0083012,14z", 
    "type": "search",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["local_results"] == nil {
    t.Error("key is not found: local_results")
    return 
  }

  if len(rsp["local_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 local_results") 
    return
  }
}  
