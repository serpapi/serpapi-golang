package main

import (
	"fmt"
	"os"

	serpapi "github.com/serpapi/serpapi-golang"
)

/***
 * Demonstrate how to upload an image with the Image API and search it with Google Lens.
 *
 * The upload returns an `image_id`, valid for 10 minutes, that replaces the `url`
 * parameter of the google_lens engine. Images must be JPG, PNG or WebP up to 500 KB.
 *
 * go run demo/demo_upload_image.go [path/to/image.png]
 *
 * Without an argument, the SerpApi logo from the test fixtures is uploaded.
 * The SERPAPI_KEY environment variable must be set to your secret SerpApi API key.
 */
func main() {
	api_key := os.Getenv("SERPAPI_KEY")
	if len(api_key) == 0 {
		println("you must obtain an api_key from serpapi\n and set the environment variable SERPAPI_KEY\n $ export SERPAPI_KEY='secret api key'")
		return
	}
	client := serpapi.NewClient(serpapi.NewSerpApiClientSetting(api_key))

	path := "test/testdata/serpapi-logo.png"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	fmt.Printf("uploading %s\n", path)
	upload, err := client.UploadImage(path)
	if err != nil {
		panic(err)
	}
	imageID := upload["image_id"].(string)
	fmt.Println(" - image uploaded with id:", imageID)

	results, err := client.Search(map[string]string{
		"engine":   "google_lens",
		"image_id": imageID,
	})
	if err != nil {
		panic(err)
	}
	// visual_matches is missing when Google Lens finds no match
	matches, _ := results["visual_matches"].([]interface{})
	fmt.Printf(" - %d visual matches\n", len(matches))
	for i, match := range matches {
		if i == 5 {
			break
		}
		fmt.Println("   -", match.(map[string]interface{})["title"])
	}
	fmt.Println("done")
}
