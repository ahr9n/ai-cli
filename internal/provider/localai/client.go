package localai

import (
	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/provider/openaicompat"
)

func init() {
	provider.Register(provider.LocalAI, provider.Registration{
		Factory:      func(baseURL string) provider.Provider { return NewClient(baseURL) },
		DefaultURL:   "http://localhost:8080/v1",
		DefaultModel: "gpt-3.5-turbo",
		Name:         "LocalAI",
		Description:  "Self-hosted AI model server compatible with OpenAI's API",
	})
}

// Client wraps openaicompat.Client with LocalAI-specific defaults.
type Client struct {
	*openaicompat.Client
}

func NewClient(baseURL string) *Client {
	return &Client{Client: openaicompat.NewClient(baseURL, "")}
}

func (c *Client) Name() string        { return "LocalAI" }
func (c *Client) Description() string { return "Self-hosted AI model server compatible with OpenAI's API" }
