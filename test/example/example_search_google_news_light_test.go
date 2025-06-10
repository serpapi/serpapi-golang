package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_news_light engine
// doc: https://serpapi.com/google-news-light-api
//
func TestGoogleNewsLight(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_news_light", 
    "q": "pizza",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["news_results"] == nil {
    t.Error("key is not found: news_results")
    return 
  }

  if len(rsp["news_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 news_results") 
    return
  }
}  
