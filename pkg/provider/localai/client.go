package localai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahr9n/ai-cli/pkg/api"
	"github.com/ahr9n/ai-cli/pkg/provider"
)

type Client struct {
	*api.BaseClient
}

type completionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type modelInfo struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
type listModelsResponse struct {
	Data []modelInfo `json:"data"`
}

func NewClient(baseURL string) provider.Provider {
	return &Client{
		BaseClient: &api.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		},
	}
}

func (c *Client) CreateCompletion(messages []provider.Message, opts *provider.CompletionOptions) (string, error) {
	var result string
	err := c.StreamCompletion(messages, opts, func(response string) {
		result += response
	})
	return result, err
}

func (c *Client) StreamCompletion(messages []provider.Message, opts *provider.CompletionOptions, onResponse func(string)) error {
	if len(messages) == 0 {
		return fmt.Errorf("no messages provided")
	}

	localMessages := make([]Message, len(messages))
	for i, msg := range messages {
		localMessages[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	reqBody := completionRequest{
		Model:       opts.Model,
		Messages:    localMessages,
		Temperature: opts.Temperature,
		Stream:      true,
	}

	resp, err := c.DoPost("v1/chat/completions", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.HandleError(resp); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var response completionResponse
		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			continue
		}

		if len(response.Choices) > 0 {
			content := response.Choices[0].Message.Content
			if content != "" {
				onResponse(content)
			}
		}
	}

	return scanner.Err()
}

func (c *Client) ListModels() ([]provider.ModelInfo, error) {
	resp, err := c.DoGet("v1/models")
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if err := c.HandleError(resp); err != nil {
		return nil, err
	}

	var response listModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]provider.ModelInfo, len(response.Data))
	for i, m := range response.Data {
		models[i] = provider.ModelInfo{
			Name:        m.ID,
			Family:      m.Object,
			Description: m.Description,
		}
	}

	return models, nil
}

func (c *Client) GetDefaultModel() string {
	return "gpt-3.5-turbo"
}

func (c *Client) Name() string {
	return "LocalAI"
}

func (c *Client) Description() string {
	return "Self-hosted AI model server compatible with OpenAI's API"
}
