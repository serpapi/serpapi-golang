package serpapi

import (
	"testing"
)

// basic use case
func TestGoogleProduct(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	client_parameter := map[string]string{
		"engine":  "google_product",
		"api_key": *getApiKey(),
	}
	client := NewClient(client_parameter)

	parameter := map[string]string{
		"product_id": "1547596720935658503",
	}
	rsp, err := client.Search(parameter)

	if err != nil {
		t.Error("search error", err)
		return
	}

	if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
		t.Error("bad search status")
		return
	}

	result, ok := rsp["product_results"]
	if !ok {
		t.Error("product_results is not found in result")
		return
	}
	if len(result.(map[string]interface{})) < 5 {
		t.Error("expect more than 5 product_results")
		return
	}
}
