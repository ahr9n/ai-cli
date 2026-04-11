package openai

import (
	"os"

	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/provider/openaicompat"
)

func init() {
	provider.Register(provider.OpenAI, provider.Registration{
		Factory: func(baseURL string) provider.Provider {
			return openaicompat.NewClient(baseURL, os.Getenv("OPENAI_API_KEY")).
				WithInfo("OpenAI", "OpenAI API (GPT-4o, GPT-4o-mini, o1, etc.)")
		},
		DefaultURL:   "https://api.openai.com/v1",
		DefaultModel: "gpt-4o-mini",
		Name:         "OpenAI",
		Description:  "OpenAI API (GPT-4o, GPT-4o-mini, o1, etc.)",
	})
}
