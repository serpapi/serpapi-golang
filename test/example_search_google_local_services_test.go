package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// basic use case
func TestGoogleLocalServices(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_local_services", 
    "q": "electrician", 
    "data_cid": "6745062158417646970",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["local_ads"] == nil {
    t.Error("key is not found: local_ads")
    return 
  }

  if len(rsp["local_ads"].([]interface{})) < 5 {
    t.Error("expect more than 5 local_ads") 
    return
  }
}  
