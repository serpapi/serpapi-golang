package main

import (
	"fmt"
	serpapi "github.com/serpapi/serpapi-golang"
	"os"
	"time"
)

/***
 * The code snippet aims to improve the efficiency of searching using the SerpApi client using `async` mode.
 * The request are non-blocking which allows batching a large amount of query, and wait before fetching the result back.
 *
 * **Process:**
 * 1. **Request Queue:** The company list is iterated over, and each company is queried using the SerpApi client. Requests
 * are stored in a queue to avoid blocking the main thread.
 *
 * 2. **Client Retrieval:** After each request, the code checks the status of the search result. If it's cached or
 * successful, the company name is printed, and the request is skipped. Otherwise, the result is added to the queue for
 * further processing.
 *
 * 3. **Queue Processing:** The queue is processed until it's empty. In each iteration, the last result is retrieved and
 * its client ID is extracted.
 *
 * 4. **Archived Client Retrieval:** Using the client ID, the code retrieves the archived client and checks its status. If
 * it's cached or successful, the company name is printed, and the client is skipped. Otherwise, the result is added back
 * to the queue for further processing.
 *
 * 5. **Completion:** The queue is closed, and a message is printed indicating that the process is complete.
 *
 * * **Asynchronous Requests:** The `async: true` option ensures that search requests are processed in parallel, improving
 * efficiency.
 * * **Queue Management:** The queue allows requests to be processed asynchronously without blocking the main thread.
 * * **Status Checking:** The code checks the status of each search result before processing it, avoiding unnecessary work.
 * * **Queue Processing:** The queue ensures that all requests are processed in the order they were submitted.
 *
 * **Overall, the code snippet demonstrates a well-structured approach to improve the efficiency of searching for company
 * information using SerpApi.**
 *
 * go get -u github.com/serpapi/serpapi-golang
 */
func main() {
	api_key := os.Getenv("SERPAPI_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable API_KEY\n $ export API_KEY='secret api key'")
	}
	setting := serpapi.NewSerpApiClientSetting(api_key)
	setting.Persistent = false                     // Enable persistent search
	setting.Asynchronous = true                    // Enable asynchronous search
	setting.Timeout = 60 * time.Second             // Set timeout for HTTP requests
	setting.MaxIdleConnection = 10                 // Set maximum idle connections
	setting.KeepAlive = 60 * time.Second           // Set keep-alive duration
	setting.TLSHandshakeTimeout = 10 * time.Second // Set TLS handshake timeout

	client := serpapi.NewClient(setting)

	// Target MAANG companies
	companyList := []string{"meta", "amazon", "apple", "netflix", "google"}
	scheduleSearch := make(chan string, len(companyList))

	var lastSearchMetadata map[string]interface{}

	for _, company := range companyList {
		// Store request into scheduleSearch - non-blocking
		fmt.Printf("Schedule search for: %s\n", company)
		result, err := client.Search(map[string]string{"q": company})
		if err != nil {
			panic(err)
		}

		searchMetadata := result["search_metadata"].(map[string]interface{})
		if status, ok := searchMetadata["status"].(string); ok && status == "Cached" {
			fmt.Printf("%s: search results found in cache for: %s\n", company, company)
		}

		// Add results to the client queue
		scheduleSearch <- searchMetadata["id"].(string)
		lastSearchMetadata = searchMetadata
	}

	fmt.Printf("Last search submitted at: %s\n", lastSearchMetadata["created_at"].(string))

	fmt.Println("Wait 5s for all requests to be completed")
	time.Sleep(5 * time.Second)

	fmt.Println("Wait until all searches are cached or successful")
	for len(scheduleSearch) > 0 {
		// Extract client ID
		searchID := <-scheduleSearch

		// Retrieve client from the archive - blocking
		searchArchived, err := client.SearchArchive(searchID)
		if err != nil {
			panic(err)
		}

		searchParameters := searchArchived["search_parameters"].(map[string]interface{})
		company := searchParameters["q"].(string)

		searchMetadata := searchArchived["search_metadata"].(map[string]interface{})
		if status, ok := searchMetadata["status"].(string); ok && (status == "Cached" || status == "Success") {
			fmt.Printf("search results found in archive for: %s\n", company)
			continue
		}

		// Add results back to the client queue if the search is still in progress
		scheduleSearch <- searchID
	}

	close(scheduleSearch)
	fmt.Println("done")
}
