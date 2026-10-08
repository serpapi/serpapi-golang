package serpapi

import (
  "testing"
  "os"
  "time"
  "github.com/serpapi/serpapi-golang"
)

// example test for google_related_questions engine
// doc: http://serpapi.com/google-related-questions-api
//
func TestGoogleRelatedQuestions(t *testing.T) {
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
    "engine": "google_related_questions", 
    "next_page_token": "eyJvbnMiOiIxMDA0MSIsImZjIjoiRW9vQkNreEJTa2M1U210T1JsQjZhMjAyTVRCMVVsRmhlV05HVnpSdUxTMTRWemhVUm5OVk4yaGxjRXRCWlRKR1RUbHlRWGswU2sxWWVXSTVYemhIWkhWTFNscExNRUpPV25WNWFtcFFUa3d3RWhaT1RIWmFXbk54WDBSTVV6RjNUalJRYnpkcVdFOUJHaUpCUmxoeVJXTnZZMHBFWmxacU5GQnljR3h0Um5KV1lsOVdOVmh1WkcxWFZGRjMiLCJmY3YiOiIzIiwiZWkiOiJOTHZaWnNxX0RMUzF3TjRQbzdqWE9BIiwicWMiOiJDZ1pqYjJabVpXVVFBSDA1VEQ4XyIsInF1ZXN0aW9uIjoiSXMgY29mZmVlIGdvb2Qgb3IgYmFkIGZvciBoZWFsdGg/IiwibGsiOiJHaUJwY3lCamIyWm1aV1VnWjI5dlpDQnZjaUJpWVdRZ1ptOXlJR2hsWVd4MGFBIiwiYnMiOiJjLVB5NGxMMExGWkl6azlMUzAxVlNNX1BUMUhJTDFKSVNreFJTQVBTR2FtSk9TVVo5aEtiZUl5VXBCUXlDYWpqY3VHU0M4LW9SRkpYbkE5UkNsSlRtVjlxTDdGWnpraGVTclljbnlJdVZ5NzU4SXpFRW9YRW9sUUZRd09GcE5TODFMVE1rbUtGX0RTb0Ruc0pFeU1GS2JseXZJcTQ0cmgwd01hazVLY0NoZk5TRlpKTEN4RFNRR0dGa255UWZVQmY1S2RVMmtzc3JqWFNsdElzaDJzeHhLOUJnQkVBIiwiaWQiOiJmY19OTHZaWnNxX0RMUzF3TjRQbzdqWE9BXzQifQ==",
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

  related_questions, ok := rsp["related_questions"].([]interface{})
  if !ok || len(related_questions) < 1 {
    t.Error("expect non-empty related_questions")
    t.Errorf("results: %v", rsp["related_questions"])
    return
  }
}  
