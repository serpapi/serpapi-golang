package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	serpapi "github.com/serpapi/serpapi-golang"
)

/***
 * Demonstrate how to run a search on Google using a thread pool
 * in order to distribute the workload across multiple threads.
 *
 * go get -u github.com/serpapi/serpapi-golang
 */
func main() {
	api_key := os.Getenv("SERPAPI_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable API_KEY\n $ export API_KEY='secret api key'")
		return
	}

	numWorkers := 4
	queries := make(chan map[string]string, numWorkers)
	results := make(chan map[string]interface{}, numWorkers)

	var wg sync.WaitGroup

	// Worker function
	worker := func(queries <-chan map[string]string, results chan<- map[string]interface{}) {
		defer wg.Done()

		setting := serpapi.NewSerpApiClientSetting(api_key)
		setting.Persistent = true // Enable persistent search
		client := serpapi.NewClient(setting)
		for query := range queries {
			fmt.Printf("Search query: %s\n", query["q"])
			result, err := client.Search(query)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			results <- result
		}
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(queries, results)
	}

	// Goroutine to process results
	go func() {
		for result := range results {
			fmt.Printf("Query Result: %v\n", result["search_metadata"].(map[string]interface{})["id"])
			fmt.Println("Organic Results:")
			for i, result := range result["organic_results"].([]interface{}) {
				jsonResult, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fmt.Printf("Error converting result %d to JSON: %v\n", i+1, err)
					continue
				}
				fmt.Printf("Result %d:\n%s\n", i+1, string(jsonResult))
			}
			fmt.Println("Total Organic Results:", len(result["organic_results"].([]interface{})))
			if organicResults, ok := result["organic_results"].([]interface{}); ok && len(organicResults) > 0 {
				if firstResult, ok := organicResults[0].(map[string]interface{}); ok {
					if title, ok := firstResult["title"].(string); ok {
						fmt.Println("Title of the first organic result:", title)
					} else {
						fmt.Println("Error: Unable to extract title from the first organic result")
					}
				} else {
					fmt.Println("Error: First organic result is not in the expected format")
				}
			} else {
				fmt.Println("Error: No organic results found or unable to parse organic results")
			}
		}
	}()

	// Schedule queries
	coffees := []string{"latte", "espresso", "cappuccino", "americano", "mocha", "macchiato", "frappuccino", "cold_brew"}
	for _, query := range coffees {
		queries <- map[string]string{
			"engine":   "google",
			"q":        query,
			"location": "Austin, Texas, United States",
			"hl":       "en",
			"gl":       "us",
			"num":      "10",
		}
		fmt.Printf("Scheduled query for: %s\n", query)
	}
	close(queries) // Close queries channel to signal workers to stop

	// Wait for all workers to finish
	wg.Wait()
	close(results) // Close results channel after all workers are done
}
