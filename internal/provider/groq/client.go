package groq

import (
	"os"

	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/provider/openaicompat"
)

func init() {
	provider.Register(provider.Groq, provider.Registration{
		Factory: func(baseURL string) provider.Provider {
			return openaicompat.NewClient(baseURL, os.Getenv("GROQ_API_KEY")).
				WithInfo("Groq", "Groq cloud inference (Llama, Mixtral, Gemma)")
		},
		DefaultURL:   "https://api.groq.com/openai/v1",
		DefaultModel: "llama-3.3-70b-versatile",
		Name:         "Groq",
		Description:  "Groq cloud inference (Llama, Mixtral, Gemma)",
	})
}
