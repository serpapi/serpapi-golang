package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_product engine
// doc: https://serpapi.com/google-product-api
//
func TestGoogleProduct(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_product", 
    "q": "coffee", 
    "product_id": "4887235756540435899",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if rsp["product_results"] == nil {
    t.Error("key is not found: product_results")
    return 
  }

  if len(rsp["product_results"].(map[string]interface{})) < 5 {
    t.Error("expect more than  5 product_results")
    return
  }
}  
