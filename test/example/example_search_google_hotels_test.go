package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_hotels engine
// doc: https://serpapi.com/google-hotels-api
//
func TestGoogleHotels(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_hotels", 
    "q": "Bali Resorts", 
    "check_in_date": "2025-05-26", 
    "check_out_date": "2025-05-27", 
    "adults": "2", 
    "currency": "USD", 
    "gl": "us", 
    "hl": "en",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["properties"] == nil {
    t.Error("key is not found: properties")
    return 
  }

  if len(rsp["properties"].([]interface{})) < 5 {
    t.Error("expect more than 5 properties") 
    return
  }
}  
