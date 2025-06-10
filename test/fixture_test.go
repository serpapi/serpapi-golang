package serpapi

import (
	"os"
)

func getApiKey() string {
	return os.Getenv("SERPAPI_KEY")
}

func shoulSkip() bool {
	return getApiKey() == ""
}
