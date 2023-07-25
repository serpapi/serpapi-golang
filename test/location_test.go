package serpapi

import (
	"testing"

	"github.com/serpapi/serpapi-golang"
)

func TestLocation(t *testing.T) {
	var rsp []interface{}
	var err error
	client := serpapi.NewClient(map[string]string{})
	rsp, err = client.Location("Austin", 3)

	if err != nil {
		t.Error(err)
	}

	//log.Println(rsp[0])
	first := rsp[0].(map[string]interface{})
	googleID := first["google_id"].(float64)
	if googleID != float64(200635) {
		t.Error(googleID)
		return
	}
}
