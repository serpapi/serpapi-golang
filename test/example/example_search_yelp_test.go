package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for yelp engine
// doc: https://serpapi.com/yelp-search-api
//
func TestYelp(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "yelp", 
    "find_desc": "Coffee", 
    "find_loc": "New York, NY, USA",  }
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
