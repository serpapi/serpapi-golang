package serpapi

import (
	"os"
)

func getApiKey() *string {
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		return nil
	}
	return &apiKey
}

func shoulSkip() bool {
	return len(*getApiKey()) == 0
}
