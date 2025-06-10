package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_shopping engine
// doc: https://serpapi.com/google-shopping-api
//
func TestGoogleShopping(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_shopping", 
    "q": "Macbook M4",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["shopping_results"] == nil {
    t.Error("key is not found: shopping_results")
    return 
  }

  if len(rsp["shopping_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 shopping_results") 
    return
  }
}  
