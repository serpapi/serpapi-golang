package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_patents engine
// doc: https://serpapi.com/google-patents-api
//
func TestGooglePatents(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_patents", 
    "q": "(Coffee)",  }
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

  if len(rsp["organic_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 organic_results") 
    return
  }
}  
