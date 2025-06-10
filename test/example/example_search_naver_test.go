package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for naver engine
// doc: https://serpapi.com/naver-search-api
//
func TestNaver(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "naver", 
    "query": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["ads_results"] == nil {
    t.Error("key is not found: ads_results")
    return 
  }

  if len(rsp["ads_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 ads_results") 
    return
  }
}  
