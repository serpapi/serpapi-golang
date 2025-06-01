package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

// doc: https://serpapi.com/locations-api
func TestLocation(t *testing.T) {
	var locations []interface{}
	var err error
	client := serpapi.NewClient(map[string]string{})
	locations, err = client.Location("Austin", 5)

	if err != nil {
		t.Error(err)
	}

	//log.Println(locations[0])
	if len(locations) < 1 {
		t.Error("expect more than 1 location")
		return
	}
	first := locations[0].(map[string]interface{})
	if _, ok := first["google_id"]; !ok {
		t.Error("key 'google_id' does not exist in the first location")
		return
	}
}
