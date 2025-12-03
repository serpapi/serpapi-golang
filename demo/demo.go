package main

import (
	"fmt"
	"os"

	serpapi "github.com/serpapi/serpapi-golang"
)

/***
 * Demonstrate how to run a simple search on Google using SerpApi.
 *
 * go get -u github.com/serpapi/serpapi-golang
 *
 * The SERPAPI_KEY environment variable must be set to your secret SerpApi API key.
 */
func main() {
	// Read SERPAPI key from environment variable
	api_key := os.Getenv("SERPAPI_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable API_KEY\n $ export API_KEY='secret api key'")
	}
	// Initialize the SerpApi client
	setting := serpapi.NewSerpApiClientSetting(api_key)
	setting.Engine = "google"    // Set the search engine to Google
	setting.Persistent = false   // Close the HTTP connection after the request to avoid keeping it open
	setting.Asynchronous = false // Block search query until results are returned
	setting.Timeout = 60         // Set timeout for HTTP requests in seconds (more timeout granularity available)
	client := serpapi.NewClient(setting)
	// define search parameters
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Austin,Texas,United States",
		"hl":       "en",
		"gl":       "us"
	}
	fmt.Println("search is running")
	data, err := client.Search(parameter)
	if err != nil {
		panic(err)
	}
	// decode data and display the first organic result title
	results := data["organic_results"].([]interface{})
	fmt.Println("search first result:")
	firstResult := results[0].(map[string]interface{})
	fmt.Println(firstResult["title"].(string))
	fmt.Println("done")
}
