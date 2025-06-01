package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_jobs engine
// doc: https://serpapi.com/google-jobs-api
//
func TestGoogleJobs(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_jobs", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["jobs_results"] == nil {
    t.Error("key is not found: jobs_results")
    return 
  }

  if len(rsp["jobs_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 jobs_results") 
    return
  }
}  
