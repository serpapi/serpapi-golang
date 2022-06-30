package main

import (
	"fmt"
	"os"

	serpapi "github.com/serpapi/serpapi-golang"
)

/***
 * Demonstrate how to search on Google
 *
 * go get -u github.com/serpapi/serpapi-golang
 */
func main() {
	api_key := os.Getenv("API_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable API_KEY\n $ export API_KEY='secret api key'")
	}
	client_parameter := map[string]string{
		"api_key": api_key,
		"engine":  "google",
	}
	client := serpapi.NewClient(client_parameter)
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Austin,Texas",
	}
	data, err := client.Search(parameter)
	if err != nil {
		panic(err)
	}
	// decode data and display the first organic result title
	results := data["organic_results"].([]interface{})
	firstResult := results[0].(map[string]interface{})
	fmt.Println(firstResult["title"].(string))
}
