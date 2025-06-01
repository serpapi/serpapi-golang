package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_flights engine
// doc: https://serpapi.com/google-flights-api
//
func TestGoogleFlights(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_flights", 
    "departure_id": "PEK", 
    "arrival_id": "AUS", 
    "outbound_date": "2025-05-26", 
    "return_date": "2025-06-01", 
    "currency": "USD", 
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

  if rsp["best_flights"] == nil {
    t.Error("key is not found: best_flights")
    return 
  }

  if len(rsp["best_flights"].([]interface{})) < 5 {
    t.Error("expect more than 5 best_flights") 
    return
  }
}  
