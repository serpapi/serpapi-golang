package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for youtube engine
// doc: https://serpapi.com/youtube-search-api
//
func TestYoutube(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "youtube", 
    "search_query": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["video_results"] == nil {
    t.Error("key is not found: video_results")
    return 
  }

  if len(rsp["video_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 video_results") 
    return
  }
}  
