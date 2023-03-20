package serpapi

import (
  "testing"
  "github.com/serpapi/serpapi-golang"
)

// basic use case
func TestGoogleProduct(t *testing.T) {
  if shoulSkip() {
    t.Skip("API_KEY required")
    return
  }

  client_parameter := map[string]string{
    "engine": "google_product",
    "api_key": *getApiKey(),
  }
  client := serpapi.NewClient(client_parameter)

  parameter := map[string]string{ 
    "q": "coffee",
 
    "product_id": "4172129135583325756",
  }
  rsp, err := client.Search(parameter)

  if err != nil {
    t.Error("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    t.Error("bad status")
    return
  }

  if len(rsp["product_results"].([]interface{})) < 5 {
    t.Error("expect more than 5 product_results") 
    return
  }
}  
