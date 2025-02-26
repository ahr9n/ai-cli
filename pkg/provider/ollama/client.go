package ollama

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ahr9n/ai-cli/pkg/api"
	"github.com/ahr9n/ai-cli/pkg/prompts"
	"github.com/ahr9n/ai-cli/pkg/provider"
)

type Client struct {
	*api.BaseClient
}

type generateRequest struct {
	Model       string    `json:"model"`
	Prompt      string    `json:"prompt"`
	Temperature float32   `json:"temperature"`
	Stream      bool      `json:"stream"`
	System      string    `json:"system,omitempty"`
	Messages    []Message `json:"messages,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
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

func NewClient(baseURL string) provider.Provider {
	return &Client{
		BaseClient: &api.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 60 * time.Second,
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

	var systemPrompt string
	var ollamaMessages []Message
	var userPrompt string
	for _, msg := range messages {
		if msg.Role == prompts.RoleSystem {
			systemPrompt = msg.Content
		} else {
			ollamaMessages = append(ollamaMessages, Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == prompts.RoleUser {
			userPrompt = messages[i].Content
			break
		}
	}

	reqBody := generateRequest{
		Model:       opts.Model,
		Prompt:      userPrompt,
		System:      systemPrompt,
		Messages:    ollamaMessages,
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
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var response generateResponse
		if err := json.Unmarshal(line, &response); err != nil {
			continue
		}

		if response.Response != "" {
			onResponse(response.Response)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

func (c *Client) ListModels() ([]provider.ModelInfo, error) {
	resp, err := c.DoGet("api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := c.HandleError(resp); err != nil {
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
