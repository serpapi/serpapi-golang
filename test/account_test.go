package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

// doc: https://serpapi.com/account-api
func TestAccount(t *testing.T) {
	// Skip this test
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	var rsp map[string]interface{}
	var err error
	setting := serpapi.NewSerpApiClientSetting(*getApiKey())
	client := serpapi.NewClient(setting)
	rsp, err = client.Account()

	if err != nil {
		t.Error("fail to fetch data")
		t.Error(err)
		return
	}

	if _, exists := rsp["account_id"]; !exists {
		t.Error("key account_id does not exist in response")
		return
	}
}
