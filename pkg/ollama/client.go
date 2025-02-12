package ollama

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahr9n/ollama-cli/pkg/api"
	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/ahr9n/ollama-cli/pkg/prompts"
)

type Message struct {
	Role    string
	Content string
}

type ChatOptions struct {
	Model       string
	Temperature float32
}

type Client struct {
	*api.BaseClient
}

type generateRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	Stream      bool    `json:"stream"`
	System      string  `json:"system"`
}

type generateResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseClient: &api.BaseClient{
			BaseURL: cfg.OllamaURL,
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

func (c *Client) StreamChatCompletion(messages []Message, opts *ChatOptions, onResponse func(string)) error {
	if len(messages) == 0 {
		return fmt.Errorf("no messages provided")
	}

	prompt := messages[len(messages)-1].Content

	reqBody := generateRequest{
		Model:       opts.Model,
		Prompt:      prompt,
		System:      prompts.DefaultSystem(),
		Temperature: opts.Temperature,
		Stream:      true,
	}

	resp, err := c.DoPost("api/generate", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("model '%s' not found - try running: ollama pull %s", opts.Model, opts.Model)
	}

	if err := c.HandleError(resp); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var response generateResponse
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		if response.Response != "" {
			onResponse(response.Response)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading response stream: %w", err)
	}

	return nil
}

func (c *Client) CreateChatCompletion(messages []Message, opts *ChatOptions) (string, error) {
	var fullResponse string
	err := c.StreamChatCompletion(messages, opts, func(response string) {
		fullResponse += response
	})
	return fullResponse, err
}
