package serpapi

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestQuickStart(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	client_parameter := map[string]string{
		"timeout": "30",
		"api_key": *getApiKey(),
		"engine":  "google",
	}
	client := NewClient(client_parameter)

	parameter := map[string]string{
		"q":             "Coffee",
		"location":      "Portland, Oregon, United States",
		"hl":            "en",
		"gl":            "us",
		"google_domain": "google.com",
		"safe":          "active",
		"start":         "10",
		"num":           "10",
		"device":        "desktop",
	}
	rsp, err := client.Search(parameter)

	if err != nil {
		t.Error(err)
		return
	}
	result := rsp["organic_results"].([]interface{})[0].(map[string]interface{})
	if len(result["title"].(string)) == 0 {
		t.Error("empty title in local results")
		return
	}
}

func TestHtml(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	client_parameter := map[string]string{
		"engine":  "google",
		"api_key": *getApiKey(),
	}
	client := NewClient(client_parameter)

	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	data, err := client.Html(parameter)
	if err != nil {
		t.Error("err must be nil")
		return
	}
	if !strings.Contains(*data, "</html>") {
		t.Error("data does not contains <html> tag")
	}
}

func TestDecodeJson(t *testing.T) {
	reader, err := os.Open("./data/search_coffee_sample.json")
	if err != nil {
		panic(err)
	}
	var client SerpApiClient
	rsp, err := client.decodeJSON(reader)
	if err != nil {
		t.Error("error should be nil", err)
		return
	}

	results := rsp["organic_results"].([]interface{})
	ref := results[0].(map[string]interface{})
	if ref["title"] != "Portland Roasting Coffee" {
		t.Error("empty title in local results")
		return
	}
}

func TestDecodeJsonPage20(t *testing.T) {
	t.Log("run test")
	reader, err := os.Open("./data/search_coffee_sample_page20.json")
	if err != nil {
		panic(err)
	}
	var client SerpApiClient
	rsp, err := client.decodeJSON(reader)
	if err != nil {
		t.Error("error should be nil")
		t.Error(err)
	}
	t.Log(reflect.ValueOf(rsp).MapKeys())
	results := rsp["organic_results"].([]interface{})
	ref := results[0].(map[string]interface{})
	t.Log(ref["title"].(string))
	if ref["title"].(string) != "Coffee | HuffPost" {
		t.Error("fail decoding the title ")
	}
}

func TestDecodeJsonError(t *testing.T) {
	reader, err := os.Open("./data/error_sample.json")
	if err != nil {
		panic(err)
	}
	var client SerpApiClient
	rsp, err := client.decodeJSON(reader)
	if rsp != nil {
		t.Error("response should not be nil")
		return
	}

	if err == nil {
		t.Error("unexcepted err is nil")
	} else if strings.Compare(err.Error(), "Your account credit is too low, plesae add more credits.") == 0 {
		t.Error("empty title in local results")
		return
	}
}

func TestLocation(t *testing.T) {
	var rsp []interface{}
	var err error
	client := NewClient(map[string]string{})
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
	client := NewClient(parameter)
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

// Search archive API
func TestSearchArchive(t *testing.T) {
	if shoulSkip() {
		t.Skip("API_KEY required")
		return
	}

	client_parameter := map[string]string{
		"engine":  "google",
		"api_key": *getApiKey(),
	}
	client := NewClient(client_parameter)
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	rsp, err := client.Search(parameter)

	if err != nil {
		t.Error("unexpected error", err)
		return
	}

	searchID := rsp["search_metadata"].(map[string]interface{})["id"].(string)
	if len(searchID) == 0 {
		t.Error("search_metadata.id must be defined")
		return
	}

	searchArchive, err := client.SearchArchive(searchID)
	if err != nil {
		t.Error(err)
		return
	}

	searchIDArchive := searchArchive["search_metadata"].(map[string]interface{})["id"].(string)
	if searchIDArchive != searchID {
		t.Error("search_metadata.id do not match", searchIDArchive, searchID)
	}
}
