package custom

import (
	"os"

	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/provider/openaicompat"
)

func init() {
	provider.Register(provider.Custom, provider.Registration{
		Factory: func(baseURL string) provider.Provider {
			return openaicompat.NewClient(baseURL, os.Getenv("AI_CLI_API_KEY")).
				WithInfo("Custom", "Custom OpenAI-compatible endpoint")
		},
		DefaultURL:   "http://localhost:8080/v1",
		DefaultModel: "",
		Name:         "Custom",
		Description:  "Custom OpenAI-compatible endpoint (set AI_CLI_API_KEY env var for auth)",
	})
}
