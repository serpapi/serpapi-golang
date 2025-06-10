package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_play engine
// doc: https://serpapi.com/google-play-api
//
func TestGooglePlay(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_play", 
    "q": "kite", 
    "store": "apps",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    t.Error("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 1 {
    t.Error("expect more than 1 organic_results") 
    return
  }
}  
