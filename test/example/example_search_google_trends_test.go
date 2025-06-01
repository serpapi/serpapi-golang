package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_trends engine
// doc: https://serpapi.com/google-trends-api
//
func TestGoogleTrends(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_trends", 
    "q": "coffee,milk,bread,pasta,steak", 
    "data_type": "TIMESERIES",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["interest_over_time"] == nil {
    t.Error("key is not found: interest_over_time")
    return 
  }

  if len(rsp["interest_over_time"].([]interface{})) < 5 {
    t.Error("expect more than 5 interest_over_time") 
    return
  }
}  
