package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for yahoo_shopping engine
// doc: http://serpapi.com/yahoo-shopping-api
//
func TestYahooShopping(t *testing.T) {
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
    "engine": "yahoo_shopping", 
    "p": "coffee", 
    "merchants": "3cf3c2e4-90aa-4398-a970-83eb69413f9b", 
    "min_price": "500", 
    "max_price": "2000", 
    "sort_by": "discountPercentage",
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

  shopping_results, ok := rsp["shopping_results"].([]interface{})
  if !ok || len(shopping_results) < 1 {
    t.Error("expect non-empty shopping_results")
    t.Errorf("results: %v", rsp["shopping_results"])
    return
  }
}  
