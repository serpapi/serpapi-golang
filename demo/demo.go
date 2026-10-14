package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	serpapi "github.com/serpapi/serpapi-golang"
)

/***
 * Demonstrate how to run a search on Google using SerpApi and paginate through the results.
 *
 * Each response includes `serpapi_pagination.next` when more results are available.
 * The query parameters of that link (e.g. `start`) are merged into the next request.
 *
 * go get -u github.com/serpapi/serpapi-golang
 *
 * The SERPAPI_KEY environment variable must be set to your secret SerpApi API key.
 */
func main() {
	// Read SERPAPI key from environment variable
	api_key := os.Getenv("SERPAPI_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable SERPAPI_KEY\n $ export SERPAPI_KEY='secret api key'")
	}
	// Initialize the SerpApi client
	setting := serpapi.NewSerpApiClientSetting(api_key)
	setting.Engine = "google"          // Set the search engine to Google
	setting.Persistent = false         // Close the HTTP connection after the request to avoid keeping it open
	setting.Asynchronous = false       // Block search query until results are returned
	setting.Timeout = 60 * time.Second // Set timeout for HTTP requests
	client := serpapi.NewClient(setting)
	// define search parameters
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Austin,Texas,United States",
		"hl":       "en",
		"gl":       "us",
	}
	maxPages := 3
	for page := 1; page <= maxPages; page++ {
		fmt.Printf("search page %d is running\n", page)
		data, err := client.Search(parameter)
		if err != nil {
			panic(err)
		}
		// decode data and display the organic result titles
		results, _ := data["organic_results"].([]interface{})
		for _, result := range results {
			fmt.Println(" -", result.(map[string]interface{})["title"])
		}

		pagination, _ := data["serpapi_pagination"].(map[string]interface{})
		next, _ := pagination["next"].(string)
		if next == "" {
			break
		}
		nextURL, err := url.Parse(next)
		if err != nil {
			panic(err)
		}
		for name, values := range nextURL.Query() {
			parameter[name] = values[0]
		}
	}
	fmt.Println("done")
}
