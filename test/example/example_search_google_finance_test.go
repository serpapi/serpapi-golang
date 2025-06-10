package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_finance engine
// doc: https://serpapi.com/google-finance-api
//
func TestGoogleFinance(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_finance", 
    "q": "GOOG:NASDAQ",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["markets"] == nil {
    t.Error("key is not found: markets")
    return 
  }

  if len(rsp["markets"].([]interface{})) < 5 {
    t.Error("expect more than 5 markets") 
    return
  }
}  
