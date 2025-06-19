package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

// doc: https://serpapi.com/locationList-api
func TestLocation(t *testing.T) {
	setting := serpapi.NewSerpApiClientSetting(getApiKey())
	client := serpapi.NewClient(setting)
	locationList, err := client.Location("Austin", 5)

	if err != nil {
		t.Error(err)
	}

	//log.Println(locationList[0])
	if len(locationList) < 1 {
		t.Error("expect more than 1 location")
		return
	}
	first := locationList[0].(map[string]interface{})
	if _, ok := first["google_id"]; !ok {
		t.Error("key 'google_id' does not exist in the first location")
		return
	}
}
