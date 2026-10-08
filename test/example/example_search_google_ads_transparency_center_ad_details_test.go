package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_ads_transparency_center_ad_details engine
// doc: http://serpapi.com/google-ads-transparency-center-ad-details-api
//
func TestGoogleAdsTransparencyCenterAdDetails(t *testing.T) {
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
    "engine": "google_ads_transparency_center_ad_details", 
    "advertiser_id": "AR12039459359856525313", 
    "creative_id": "CR11382217505440268289",
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

  ad_creatives, ok := rsp["ad_creatives"].([]interface{})
  if !ok || len(ad_creatives) < 1 {
    t.Error("expect non-empty ad_creatives")
    t.Errorf("results: %v", rsp["ad_creatives"])
    return
  }
}  
