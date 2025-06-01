package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for home_depot engine
// doc: https://serpapi.com/home-depot-search-api
//
func TestHomeDepot(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "home_depot", 
    "q": "table",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["products"] == nil {
    t.Error("key is not found: products")
    return 
  }

  if len(rsp["products"].([]interface{})) < 5 {
    t.Error("expect more than 5 products") 
    return
  }
}  
