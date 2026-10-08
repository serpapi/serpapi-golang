package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for bing_product engine
// doc: http://serpapi.com/bing-product-api
//
func TestBingProduct(t *testing.T) {
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
    "engine": "bing_product", 
    "product_token": "pc13NHicdY3BbsIwEETVjyES8oUaEXrJAWgJqBQiwg-YeEOsWF53vRblV_s1jRIOPcBlR_tmRvP7Mv7O9qYFsbSqasWnsTY0yEKKHShugETpQLVAIamN5U6zUIFTZHA0VSMpX-fdEZfTzcMdyB5s9fDKVM5nMp1M32bp3dio0AzmpCe5xbOyh7oG2urwuLaI1Qodww_n0eh_7Zww-g_Hhm_PJvfoSo8uIIHuV4bYiSJ0knhCHSv26gIZdyxZH45fWbkpivciORNew8D_AJZIW7o",
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

  product_results, ok := rsp["product_results"].([]interface{})
  if !ok || len(product_results) < 1 {
    t.Error("expect non-empty product_results")
    t.Errorf("results: %v", rsp["product_results"])
    return
  }
}  
