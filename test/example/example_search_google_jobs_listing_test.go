package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_jobs_listing engine
// doc: http://serpapi.com/google-jobs-listing-api
//
func TestGoogleJobsListing(t *testing.T) {
  // skip test if no SERPAPI_KEY provided
  //  and account is required. see: https://serpapi.com/
  api_key := os.Getenv("SERPAPI_KEY")
  if api_key == "" {
    t.Skip("SERPAPI_KEY required")
    return
  }

  // Initialize the client with custom setting
	setting := serpapi.NewSerpApiClientSetting(api_key) // Replace with your SerpApi key
	setting.Persistent = false                     // Enable persistent search
	setting.Asynchronous = false                   // Enable asynchronous search
	setting.Timeout = 60 * time.Second             // Set timeout for HTTP requests
	setting.MaxIdleConnection = 10                 // Set maximum idle connections
	setting.KeepAlive = 60 * time.Second           // Set keep-alive duration
	setting.TLSHandshakeTimeout = 10 * time.Second // Set TLS handshake timeout

  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_jobs_listing", 
    "q": "eyJqb2JfdGl0bGUiOiJCYXJpc3RhIiwiaHRpZG9jaWQiOiJ5Vy1laV9FQ3Y3Z0FBQUFBQUFBQUFBPT0iLCJnbCI6InVzIiwiaGwiOiJlbiIsImZjIjoiRXZjQkNyY0JRVUYwVm14aVJETmtXVmxsYm5SNVNqZFVNM3BEVkd0d1drcFdZVXRzTTNOQmFIaHVPVEpXWWsxbGVsRldiMGxYVjBWdUxVdzNYMlF5V0VKTVpEaDRMVkZ6Umtwek5qSklaRkJtVTJReU5FbGxZa0ZDWnpCemVUY3lYemc1UkU5blNIWlpRVnBRU1doMFJHMXljRk50VkhCemJsOUxjbUprYURKNU4ybE5hMmt5Vmpkc2RuUmpORnB3VkcwemEzUmFTV3RZYWxGcmFHRjJkek0yTVcxeGNGbGliM2xCWmtveVl6ZDJRMTlrYTB0alYzQkpjbVZ2RWhkSVNHVnNXVFpFY2toTU1tOXhkSE5RYms1MVIzRkJaeG9pUVVSVmVVVkhaV2xpVmxaaVgxRnRkbXRrVmpaVWQxVnVhbWsxYW5KT2QyaE9adyIsImZjdiI6IjMiLCJmY19pZCI6ImZjXzEiLCJhcHBseV9saW5rIjp7InRpdGxlIjoiLm5GZzJlYntmb250LXdlaWdodDo1MDB9LkJpNkRkY3tmb250LXdlaWdodDo1MDB9QXBwbHkgZGlyZWN0bHkgb24gSW5kZWVkIiwibGluayI6Imh0dHBzOi8vd3d3LmluZGVlZC5jb20vdmlld2pvYj9qaz03ZTA0YWYyNmIyZGE2NjljXHUwMDI2dXRtX2NhbXBhaWduPWdvb2dsZV9qb2JzX2FwcGx5XHUwMDI2dXRtX3NvdXJjZT1nb29nbGVfam9ic19hcHBseVx1MDAyNnV0bV9tZWRpdW09b3JnYW5pYyJ9fQ",
  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error: ", err)
    return
  }

  status := rsp["search_metadata"].(map[string]interface{})["status"]
  if status != "Success" {
    t.Error("unexpected status: ", status)
    return
  }

  apply_options, ok := rsp["apply_options"].([]interface{})
  if !ok || len(apply_options) < 1 {
    t.Error("expect non-empty apply_options")
    t.Errorf("results: %v", rsp["apply_options"])
    return
  }
}  
