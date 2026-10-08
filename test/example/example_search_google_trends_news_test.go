package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_trends_news engine
// doc: http://serpapi.com/google-trends-news-api
//
func TestGoogleTrendsNews(t *testing.T) {
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
    "engine": "google_trends_news", 
    "page_token": "GOpbs3ica1xTlFpYmlpcEp-SWJI4bfK5BzeuRi5KzVsUGjz53MObP3Rg7LvN0xXg4jUfT8HZd3nPwthgvXk5QDYAHy0v_g",
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

  trending_searches, ok := rsp["trending_searches"].([]interface{})
  if !ok || len(trending_searches) < 1 {
    t.Error("expect non-empty trending_searches")
    t.Errorf("results: %v", rsp["trending_searches"])
    return
  }
}  
