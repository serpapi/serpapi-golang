package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)


func TestAccount(t *testing.T) {
	// Skip this test
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	var rsp map[string]interface{}
	var err error
	parameter := map[string]string{
		"api_key": *getApiKey(),
	}
	client := serpapi.NewClient(parameter)
	rsp, err = client.Account()

	if err != nil {
		t.Error("fail to fetch data")
		t.Error(err)
		return
	}

	if rsp["account_id"] == nil {
		t.Error("no account_id found")
		return
	}
}