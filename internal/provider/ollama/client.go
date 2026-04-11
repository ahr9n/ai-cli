package ollama

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ahr9n/ai-cli/internal/api"
	"github.com/ahr9n/ai-cli/internal/provider"
)

type Client struct {
	base *api.BaseClient
}

type chatRequest struct {
	Model       string             `json:"model"`
	Messages    []provider.Message `json:"messages"`
	Temperature float32            `json:"temperature"`
	Stream      bool               `json:"stream"`
}

type chatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

type modelInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	Details  struct {
		Family  string `json:"family"`
		License string `json:"license"`
	} `json:"details"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		base: &api.BaseClient{
			BaseURL:    baseURL,
			HTTPClient: &http.Client{},
		},
	}
}

func (c *Client) CreateCompletion(ctx context.Context, messages []provider.Message, opts *provider.CompletionOptions) (string, error) {
	var result strings.Builder
	err := c.StreamCompletion(ctx, messages, opts, func(chunk string) {
		result.WriteString(chunk)
	})
	return result.String(), err
}

func (c *Client) StreamCompletion(ctx context.Context, messages []provider.Message, opts *provider.CompletionOptions, onResponse func(string)) error {
	if len(messages) == 0 {
		return fmt.Errorf("no messages provided")
	}

	reqBody := chatRequest{
		Model:       opts.Model,
		Messages:    messages,
		Temperature: opts.Temperature,
		Stream:      true,
	}

	resp, err := c.base.DoPost(ctx, "api/chat", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("model '%s' not found - try running: ollama pull %s", opts.Model, opts.Model)
	}

	if err := c.base.HandleError(resp); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var response chatResponse
		if err := json.Unmarshal(line, &response); err != nil {
			continue
		}

		if response.Message.Content != "" {
			onResponse(response.Message.Content)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

func (c *Client) ListModels(ctx context.Context) ([]provider.ModelInfo, error) {
	resp, err := c.base.DoGet(ctx, "api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.base.HandleError(resp); err != nil {
		return nil, err
	}

	var response struct {
		Models []modelInfo `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]provider.ModelInfo, len(response.Models))
	for i, m := range response.Models {
		models[i] = provider.ModelInfo{
			Name:     m.Name,
			Size:     m.Size,
			Modified: m.Modified,
			Family:   m.Details.Family,
		}
	}

	return models, nil
}

func (c *Client) GetDefaultModel() string {
	return "deepseek-r1:1.5b"
}

func (c *Client) Name() string {
	return "Ollama"
}

func (c *Client) Description() string {
	return "Run large language models locally"
}
