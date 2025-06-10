package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_immersive_product engine
// doc: https://serpapi.com/google-immersive-product-api
//
func TestGoogleImmersiveProduct(t *testing.T) {
  if shoulSkip() {
    t.Skip("SERPSERPAPI_KEY required")
    return
  }

  auth := map[string]string{
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(auth)

  parameter := map[string]string{
    "engine": "google_immersive_product", 
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

  if rsp["immersive_product_results"] == nil {
    t.Error("key is not found: immersive_product_results")
    return 
  }

  if len(rsp["immersive_product_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 immersive_product_results") 
    return
  }
}  
