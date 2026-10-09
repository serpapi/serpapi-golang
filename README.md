## SerpApi Go Library

[![serpapi-go](https://github.com/serpapi/serpapi-golang/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/serpapi/serpapi-golang/actions/workflows/ci.yml)

Integrate search data into your Go application. This library is the official wrapper for [SerpApi](https://serpapi.com).

[SerpApi](https://serpapi.com) supports Google, Google Maps, Google Shopping, Baidu, Yandex, Yahoo, eBay, App Stores, and more.

## Installation

Go 1.17+ is required.

```bash
go get -u github.com/serpapi/serpapi-golang
```

## Quick start

```golang
import (
  "fmt"
  "net/url"

  "github.com/serpapi/serpapi-golang"
)

setting := serpapi.NewSerpApiClientSetting("<SERPAPI_KEY>") // Replace with your SerpApi key
setting.Engine = "google" // Set the search engine to Google
client := serpapi.NewClient(setting)
parameter := map[string]string{
  "q":             "Coffee",
  "location":      "Austin, Texas, United States",
}
for page := 1; page <= 3; page++ {
  results, err := client.Search(parameter)
  if err != nil {
    panic(err)
  }
  fmt.Println(results["organic_results"])

  // follow serpapi_pagination.next until there are no more pages
  pagination, _ := results["serpapi_pagination"].(map[string]interface{})
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
 ```

This example runs a search for "coffee" on Google and walks through the first 3 pages of results, each returned as a Go map.
See [Pagination](#pagination) for details and the [playground](https://serpapi.com/playground) to generate your own code.

### Context-aware requests

The existing methods remain available and continue to use the client's configured
timeout. For request-scoped cancellation or deadlines, use the corresponding
`Context` method:

```golang
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

results, err := client.SearchContext(ctx, map[string]string{
  "q": "Coffee",
})
if err != nil {
  // err may be context.Canceled or context.DeadlineExceeded.
  panic(err)
}
fmt.Println(results)
```

Context-aware variants are available for `Search`, `Html`, `Markdown`, `Location`,
`Account`, and `SearchArchive`. The original methods are backward-compatible
convenience wrappers.

### HTTP errors

Non-successful HTTP responses return `*serpapi.HTTPError`. It exposes the HTTP
status code, request URL (API key removed), and response body so callers can handle rate limits
and authentication failures explicitly:

```golang
var httpErr *serpapi.HTTPError
if errors.As(err, &httpErr) {
  fmt.Println(httpErr.StatusCode, httpErr.Body)
}
```

## Advanced usage

### Search API
```golang

func main() {
  // Initialize the client with custom setting
	setting := serpapi.NewSerpApiClientSetting("<SERPAPI_KEY>") // Replace with your SerpApi key
	setting.Persistent = false                     // Close the HTTP connection after each request
	setting.Asynchronous = true                    // Enable asynchronous search
	setting.Timeout = 60 * time.Second             // Set timeout for HTTP requests
	setting.MaxIdleConnection = 10                 // Set maximum idle connections
	setting.KeepAlive = 60 * time.Second           // Set keep-alive duration
	setting.TLSHandshakeTimeout = 10 * time.Second // Set TLS handshake timeout

	client := serpapi.NewClient(setting)

  // search query overview (more fields available depending on search engine)
  parameter := map[string]string{
    "q":             "Coffee",
    "location":      "Austin, Texas, United States",
    "hl":            "en",
    "gl":            "us",
    "google_domain": "google.com",
    "safe":          "active",
    "start":         "10",
    "device":        "desktop",
  }

  // formatted search results as a map
  // serpapi.com converts HTML -> JSON
  rsp, err := client.Search(parameter)

  if err != nil {
    panic(err)
  }
  fmt.Println(rsp)

  // raw search engine html as a String
  // serpapi.com acts as a proxy to provide high throughput, no search limit and more.
  raw_html, err := client.Html(parameter)
  if err != nil {
    panic(err)
  }
  fmt.Println(raw_html)

  // search results as markdown optimized for LLMs and AI agents
  // serpapi.com converts HTML -> Markdown
  markdown, err := client.Markdown(parameter)
  if err != nil {
    panic(err)
  }
  fmt.Println(markdown)
}
```

[Google search documentation](https://serpapi.com/search-api).
More hands on examples are available below.

#### Documentation

 * [Full documentation on SerpApi.com](https://serpapi.com)
 * [Library Github page](https://github.com/serpapi/serpapi-golang)
 * [API health status](https://serpapi.com/status)

### Markdown API

`Markdown` returns the search results as markdown instead of JSON. The format is
optimized for LLMs and AI agents: it carries most of the information found in the
JSON response, but is far more token efficient thanks to tables, markdown links
and YAML frontmatter.

It is the `output=md` variant of the Search API, so it accepts exactly the same
parameters as `Search` and `Html`.

```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
  "strings"
)

 func main() {

	// Initialize the SerpApi client with the API key
	// and set the search engine to Google
	setting := serpapi.NewSerpApiClientSetting("secret_api_key")
	setting.Engine = "google" // Set the search engine to Google
	client := serpapi.NewClient(setting)

	// Define the search parameters
	parameter := map[string]string{
		"q":        "Coffee",
		"location": "Portland"}

	// Perform the search and get the markdown response
	data, err := client.Markdown(parameter)
	if err != nil {
		fmt.Println("err must be nil")
		return
	}
	if !strings.Contains(*data, "#") {
		fmt.Println("data does not contain any markdown heading")
	}
}

```

 * source code: [test/markdown_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/markdown_test.go)

 * [Markdown output documentation](https://serpapi.com/blog/turning-search-results-into-markdown-for-llms/)

### Location API

```golang
import (
	"fmt"
	"github.com/serpapi/serpapi-golang"
)

setting := serpapi.NewSerpApiClientSetting("<SERPAPI_KEY>") // Replace with your SerpApi secret key
client := serpapi.NewClient(setting)
locationList, err := client.Location("Austin", 5)

if err != nil {
  panic(err)
}
fmt.Println(locationList)
```

It prints the first 5 locations matching Austin (Texas, Texas, Rochester)
```
[map[canonical_name:Austin,TX,Texas,United States country_code:US google_id:200635 google_parent_id:21176 gps:[-97.7430608 30.267153]...
```

 * source code: [test/location_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/location_test.go)

### Search Archive API

This API allows retrieving previous search results.
To fetch earlier results from the search_id.

First, you need to run a search and save the search id.
```golang
setting := serpapi.NewSerpApiClientSetting("<SERPAPI_KEY>") // Replace with your SerpApi key
client := serpapi.NewClient(setting)
parameter := map[string]string{
  "q":        "Coffee",
  "location": "Portland"}

rsp, err := client.Search(parameter)
if err != nil {
  panic(err)
}

// Now let's retrieve the previous search results from the archive.
searchID := rsp["search_metadata"].(map[string]interface{})["id"].(string)
searchArchive, err := client.SearchArchive(searchID)
if err != nil {
  panic(err)
}
fmt.Println(searchArchive["search_metadata"])
```

This code prints the search results from the archive. :)

### Account API
```golang
setting := serpapi.NewSerpApiClientSetting("<SERPAPI_KEY>") // Replace with your SerpApi key
client := serpapi.NewClient(setting)
rsp, err := client.Account()
if err != nil {
  panic(err)
}
fmt.Println(rsp)
```

It prints your account information.

## Basic examples in Go

### Search Google
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_test.go)
* see: [serpapi.com/search-api](https://serpapi.com/search-api)

### Search Google Light
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_light", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_light_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_light_test.go)
* see: [serpapi.com/google-light-api](https://serpapi.com/google-light-api)

### Search Google Scholar
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_scholar", 
    "q": "biology",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_scholar_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_scholar_test.go)
* see: [serpapi.com/google-scholar-api](https://serpapi.com/google-scholar-api)

### Search Google Autocomplete
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_autocomplete", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["suggestions"] == nil {
    fmt.Println("key is not found: suggestions")
    return 
  }

  if len(rsp["suggestions"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 suggestions") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_autocomplete_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_autocomplete_test.go)
* see: [serpapi.com/google-autocomplete-api](https://serpapi.com/google-autocomplete-api)

### Search Google Product
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_product", 
    "q": "coffee", 
    "product_id": "4887235756540435899",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["product_results"] == nil {
    fmt.Println("key is not found: product_results")
    return 
  }

  if len(rsp["product_results"].(map[string]interface{})) < 5 {
    fmt.Println("expect more than  5 product_results")
    return
  }
}  

```

 * source code: [test/example/example_search_google_product_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_product_test.go)
* see: [serpapi.com/google-product-api](https://serpapi.com/google-product-api)

### Search Google Reverse Image
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_reverse_image", 
    "image_url": "https://i.imgur.com/5bGzZi7.jpg",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["image_sizes"] == nil {
    fmt.Println("key is not found: image_sizes")
    return 
  }

  if len(rsp["image_sizes"].([]interface{})) < 1 {
    fmt.Println("expect more than 1 image_sizes") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_reverse_image_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_reverse_image_test.go)
* see: [serpapi.com/google-reverse-image](https://serpapi.com/google-reverse-image)

### Search Google Events
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_events", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["events_results"] == nil {
    fmt.Println("key is not found: events_results")
    return 
  }

  if len(rsp["events_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 events_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_events_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_events_test.go)
* see: [serpapi.com/google-events-api](https://serpapi.com/google-events-api)

### Search Google Local Services
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_local_services", 
    "q": "electrician", 
    "data_cid": "6745062158417646970",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["local_ads"] == nil {
    fmt.Println("key is not found: local_ads")
    return 
  }

  if len(rsp["local_ads"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 local_ads") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_local_services_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_local_services_test.go)
* see: [serpapi.com/google-local-services-api](https://serpapi.com/google-local-services-api)

### Search Google Maps
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_maps", 
    "q": "Coffee", 
    "ll": "@40.7455096,-74.0083012,14z", 
    "type": "search",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["local_results"] == nil {
    fmt.Println("key is not found: local_results")
    return 
  }

  if len(rsp["local_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 local_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_maps_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_maps_test.go)
* see: [serpapi.com/google-maps-api](https://serpapi.com/google-maps-api)

### Search Google Jobs
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_jobs", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["jobs_results"] == nil {
    fmt.Println("key is not found: jobs_results")
    return 
  }

  if len(rsp["jobs_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 jobs_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_jobs_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_jobs_test.go)
* see: [serpapi.com/google-jobs-api](https://serpapi.com/google-jobs-api)

### Search Google Play
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_play", 
    "q": "kite", 
    "store": "apps",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 1 {
    fmt.Println("expect more than 1 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_play_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_play_test.go)
* see: [serpapi.com/google-play-api](https://serpapi.com/google-play-api)

### Search Google Images
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_images", 
    "tbm": "isch", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["images_results"] == nil {
    fmt.Println("key is not found: images_results")
    return 
  }

  if len(rsp["images_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 images_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_images_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_images_test.go)
* see: [serpapi.com/images-results](https://serpapi.com/images-results)

### Search Google Lens
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_lens", 
    "url": "https://i.imgur.com/HBrB8p0.png",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["visual_matches"] == nil {
    fmt.Println("key is not found: visual_matches")
    return 
  }

  if len(rsp["visual_matches"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 visual_matches") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_lens_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_lens_test.go)
* see: [serpapi.com/google-lens-api](https://serpapi.com/google-lens-api)

### Search Google Images Light
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_images_light", 
    "q": "Coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["images_results"] == nil {
    fmt.Println("key is not found: images_results")
    return 
  }

  if len(rsp["images_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 images_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_images_light_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_images_light_test.go)
* see: [serpapi.com/google-images-light-api](https://serpapi.com/google-images-light-api)

### Search Google Hotels
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_hotels", 
    "q": "Bali Resorts", 
    "check_in_date": "2025-05-26", 
    "check_out_date": "2025-05-27", 
    "adults": "2", 
    "currency": "USD", 
    "gl": "us", 
    "hl": "en",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["properties"] == nil {
    fmt.Println("key is not found: properties")
    return 
  }

  if len(rsp["properties"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 properties") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_hotels_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_hotels_test.go)
* see: [serpapi.com/google-hotels-api](https://serpapi.com/google-hotels-api)

### Search Google Flights
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_flights", 
    "departure_id": "PEK", 
    "arrival_id": "AUS", 
    "outbound_date": "2025-05-26", 
    "return_date": "2025-06-01", 
    "currency": "USD", 
    "hl": "en",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["best_flights"] == nil {
    fmt.Println("key is not found: best_flights")
    return 
  }

  if len(rsp["best_flights"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 best_flights") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_flights_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_flights_test.go)
* see: [serpapi.com/google-flights-api](https://serpapi.com/google-flights-api)

### Search Google Finance
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_finance", 
    "q": "GOOG:NASDAQ",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["markets"] == nil {
    fmt.Println("key is not found: markets")
    return 
  }

  if len(rsp["markets"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 markets") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_finance_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_finance_test.go)
* see: [serpapi.com/google-finance-api](https://serpapi.com/google-finance-api)

### Search Google News
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_news", 
    "q": "pizza", 
    "gl": "us", 
    "hl": "en",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["news_results"] == nil {
    fmt.Println("key is not found: news_results")
    return 
  }

  if len(rsp["news_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 news_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_news_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_news_test.go)
* see: [serpapi.com/google-news-api](https://serpapi.com/google-news-api)

### Search Google News Light
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_news_light", 
    "q": "pizza",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["news_results"] == nil {
    fmt.Println("key is not found: news_results")
    return 
  }

  if len(rsp["news_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 news_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_news_light_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_news_light_test.go)
* see: [serpapi.com/google-news-light-api](https://serpapi.com/google-news-light-api)

### Search Google Patents
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_patents", 
    "q": "(Coffee)",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_patents_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_patents_test.go)
* see: [serpapi.com/google-patents-api](https://serpapi.com/google-patents-api)

### Search Google Trends
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_trends", 
    "q": "coffee,milk,bread,pasta,steak", 
    "data_type": "TIMESERIES",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["interest_over_time"] == nil {
    fmt.Println("key is not found: interest_over_time")
    return 
  }

  if len(rsp["interest_over_time"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 interest_over_time") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_trends_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_trends_test.go)
* see: [serpapi.com/google-trends-api](https://serpapi.com/google-trends-api)

### Search Google Shopping
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_shopping", 
    "q": "Macbook M4",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["shopping_results"] == nil {
    fmt.Println("key is not found: shopping_results")
    return 
  }

  if len(rsp["shopping_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 shopping_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_shopping_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_shopping_test.go)
* see: [serpapi.com/google-shopping-api](https://serpapi.com/google-shopping-api)

### Search Google Immersive Product
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_immersive_product", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["immersive_product_results"] == nil {
    fmt.Println("key is not found: immersive_product_results")
    return 
  }

  if len(rsp["immersive_product_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 immersive_product_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_immersive_product_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_immersive_product_test.go)
* see: [serpapi.com/google-immersive-product-api](https://serpapi.com/google-immersive-product-api)

### Search Google Videos
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "google_videos", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_google_videos_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_google_videos_test.go)
* see: [serpapi.com/google-videos-api](https://serpapi.com/google-videos-api)

### Search Amazon
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "amazon", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_amazon_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_amazon_test.go)
* see: [serpapi.com/amazon-search-api](https://serpapi.com/amazon-search-api)

### Search Baidu
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "baidu", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_baidu_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_baidu_test.go)
* see: [serpapi.com/baidu-search-api](https://serpapi.com/baidu-search-api)

### Search Yahoo
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "yahoo", 
    "p": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_yahoo_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_yahoo_test.go)
* see: [serpapi.com/yahoo-search-api](https://serpapi.com/yahoo-search-api)

### Search Youtube
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "youtube", 
    "search_query": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["video_results"] == nil {
    fmt.Println("key is not found: video_results")
    return 
  }

  if len(rsp["video_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 video_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_youtube_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_youtube_test.go)
* see: [serpapi.com/youtube-search-api](https://serpapi.com/youtube-search-api)

### Search Walmart
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "walmart", 
    "query": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_walmart_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_walmart_test.go)
* see: [serpapi.com/walmart-search-api](https://serpapi.com/walmart-search-api)

### Search eBay
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "ebay", 
    "_nkw": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_ebay_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_ebay_test.go)
* see: [serpapi.com/ebay-search-api](https://serpapi.com/ebay-search-api)

### Search Naver
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "naver", 
    "query": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["ads_results"] == nil {
    fmt.Println("key is not found: ads_results")
    return 
  }

  if len(rsp["ads_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 ads_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_naver_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_naver_test.go)
* see: [serpapi.com/naver-search-api](https://serpapi.com/naver-search-api)

### Search Home Depot
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "home_depot", 
    "q": "table",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["products"] == nil {
    fmt.Println("key is not found: products")
    return 
  }

  if len(rsp["products"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 products") 
    return
  }
}  

```

 * source code: [test/example/example_search_home_depot_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_home_depot_test.go)
* see: [serpapi.com/home-depot-search-api](https://serpapi.com/home-depot-search-api)

### Search Apple App Store
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "apple_app_store", 
    "term": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_apple_app_store_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_apple_app_store_test.go)
* see: [serpapi.com/apple-app-store](https://serpapi.com/apple-app-store)

### Search DuckDuckGo
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "duckduckgo", 
    "q": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_duckduckgo_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_duckduckgo_test.go)
* see: [serpapi.com/duckduckgo-search-api](https://serpapi.com/duckduckgo-search-api)

### Search Yandex
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "yandex", 
    "text": "coffee",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_yandex_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_yandex_test.go)
* see: [serpapi.com/yandex-search-api](https://serpapi.com/yandex-search-api)

### Search Yelp
```golang
 import (	
  "github.com/serpapi/serpapi-golang" 
  "fmt"
)

 func main() {

  setting := serpapi.NewSerpApiClientSetting("secret_api_key")
  client := serpapi.NewClient(setting)

  parameter := map[string]string{
    "engine": "yelp", 
    "find_desc": "Coffee", 
    "find_loc": "New York, NY, USA",  }
  rsp, err := client.Search(parameter)

  if err != nil {
    fmt.Println("unexpected error", err)
    return
  }

  if rsp["search_metadata"].(map[string]interface{})["status"] != "Success" {
    fmt.Println("bad status")
    return
  }

  if rsp["organic_results"] == nil {
    fmt.Println("key is not found: organic_results")
    return 
  }

  if len(rsp["organic_results"].([]interface{})) < 5 {
    fmt.Println("expect more than 5 organic_results") 
    return
  }
}  

```

 * source code: [test/example/example_search_yelp_test.go](https://github.com/serpapi/serpapi-golang/blob/master/test/example/example_search_yelp_test.go)
* see: [serpapi.com/yelp-search-api](https://serpapi.com/yelp-search-api)

## Advanced search API usage
### Highly scalable batching

Search API features non-blocking search using the option: `async=true`.
 - Non-blocking - async=true - a single parent process can handle unlimited concurrent searches.
 - Blocking - async=false - many processes must be forked and synchronized to handle concurrent searches. This strategy is I/O intensive because each client would hold a network connection.

Search API enables `async` search.
 - Non-blocking (`async=true`) : the development is more complex, but this allows handling many simultaneous connections.
 - Blocking (`async=false`) : it's easy to write the code but more compute-intensive when the parent process needs to hold many connections.

Here is an example of asynchronous searches using Go 
```golang
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
 *
 * The SERPAPI_KEY environment variable must be set to your secret SerpApi API key.
 */
func main() {
	// Read SERPAPI key from environment variable
	api_key := os.Getenv("SERPAPI_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable SERPAPI_KEY\n $ export SERPAPI_KEY='secret api key'")
	}
	setting := serpapi.NewSerpApiClientSetting(api_key)
	setting.Persistent = false                     // Close the HTTP connection after each request
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

```

 * source code: [demo/demo_async.go](https://github.com/serpapi/serpapi-golang/blob/master/demo/demo_async.go)

This code shows a simple solution to batch searches asynchronously into a [queue](https://en.wikipedia.org/wiki/Queue_(abstract_data_type)). 
Each search takes a few seconds before completion by SerpApi service and the search engine. By the time the first element pops out of the queue. The search result might be already available in the archive. If not, the `search_archive` method blocks until the search results are available. 

### Pagination

Each response includes `serpapi_pagination.next` when more results are available.
Merge the query parameters of that link into the next request until it is missing.

```golang
parameter := map[string]string{"q": "Coffee", "location": "Austin,Texas,United States"}
for page := 1; page <= 3; page++ {
	data, err := client.Search(parameter)
	if err != nil {
		panic(err)
	}
	results, _ := data["organic_results"].([]interface{})
	fmt.Printf("page %d: %d results\n", page, len(results))

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
```

 * source code: [demo/demo.go](https://github.com/serpapi/serpapi-golang/blob/master/demo/demo.go)

## Supported Go versions
Go versions validated by Github Actions:
 - 1.17+
 * see: [Github Actions.](https://github.com/serpapi/serpapi-golang/actions/workflows/ci.yml)

## Changelog
 * [2026-09-24] 1.3.0 Context-aware requests
  - Added backward-compatible `Context` variants for all API methods
  - Added request cancellation and deadline support
 * [2026-08-16] 1.2.0 Markdown output support
  - New `Markdown()` method returning `output=md` results optimized for LLMs and AI agents
 * [2026-01-26] 1.1.0 Asynchronous & Persistent Mode Support
  - Major features (async/persistent mode, API key handling, client configuration)
  - New test examples
  - Documentation and code quality improvements
  - Build system and CI/CD enhancements
 * [2024-10-01] 1.0.0 Full API support

## Developer guide
### Key goals
 - Brand centric instead of search engine based
   - No hard-coded logic per search engine
 - Simple HTTP client (lightweight, reduced dependency)
   - No magic default values
   - Thread safe
 - Easy extension
 - Defensive code style (raise a custom exception)
 - TDD
 - Best API coding practice per platform
 - KiSS principles

### Inspirations
The source code and coding style of this project are inspired by Go.
The Go programming language provides native recommendations for building excellent software.

### Code quality expectations
 - 0 lint offense: `make lint`
 - 100% tests passing: `make test`
 - Code coverage report: `make coverage`

## Design : UML diagram
### Class diagram
```mermaid
classDiagram
  Application *-- serpapi 
  serpapi *-- Client
  class Client {
    engine String
    api_key String
    params Map
    search() Map
    html() String
    markdown() String
    location() String
    search_archive() Map
    account() Map
  }
  net/http <.. Client
  json <.. Client
  Go <.. net/http
  Go <.. json
```
### search() : Sequence diagram
```mermaid
sequenceDiagram
    Client->>SerpApi.com: search() : http request 
    SerpApi.com-->>SerpApi.com: query search engine
    SerpApi.com-->>SerpApi.com: parse HTML into JSON
    SerpApi.com-->>Client: JSON string payload
    Client-->>Client: decode JSON into map
```
where:
  - The end user implements the application.
  - Client refers to serpapi.Client.
  - SerpApi.com is the backend HTTP / REST service.
  - Engine refers to Google, Baidu, Bing, and more.

The SerpApi.com service (backend)
 - executes a scalable search on `engine: "google"` using the search query: `q: "coffee"`.
 - parses the messy HTML responses from Google on the backend.
 - returns a standardized JSON response.
The class serpapi.Client (client side / golang):
 - Format the request to SerpApi.com server.
 - Execute HTTP Get request.
 - Parse JSON into Go map using a standard JSON library.
Et voila!

## Continuous integration
We love "true open source" and "continuous integration", and Test Driven Development (TDD).
 We are using Go test to test [our infrastructure around the clock]) using Github Action to achieve the best QoS (Quality Of Service).

The directory test/ includes specification which serves the dual purposes of examples and functional tests.

Set your secret API key in your shell before running a test.
```bash
export SERPAPI_KEY="your_secret_key"
```
Install testing dependency
```bash
$ make test
```
Contributions are welcome. Feel free to submit a pull request!

## License

MIT License.
