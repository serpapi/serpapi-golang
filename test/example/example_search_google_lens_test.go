package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_lens engine
// doc: https://serpapi.com/google-lens-api
//
func TestGoogleLens(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_lens", 
    "url": "https://i.imgur.com/HBrB8p0.png",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["visual_matches"] == nil {
    t.Error("key is not found: visual_matches")
    return 
  }

  if len(rsp["visual_matches"].([]interface{})) < 5 {
    t.Error("expect more than 5 visual_matches") 
    return
  }
}  
