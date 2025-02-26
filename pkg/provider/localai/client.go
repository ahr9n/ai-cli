package localai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

type streamingResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
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
				Timeout: 60 * time.Second, // Increased timeout
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
	// Increase scanner buffer to 10MB (default is 64KB)
	const maxScanTokenSize = 10 * 1024 * 1024
	scanBuf := make([]byte, maxScanTokenSize)
	scanner.Buffer(scanBuf, maxScanTokenSize)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip "data: " prefix if present (SSE format)
		if len(line) > 6 && line[0:6] == "data: " {
			line = line[6:]
		}

		// Skip "[DONE]" message
		if line == "[DONE]" {
			continue
		}

		var streamResp streamingResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			// Try parsing as non-streaming response
			var response completionResponse
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				continue // Skip lines we can't parse
			}

			// Handle non-streaming response
			if len(response.Choices) > 0 {
				content := response.Choices[0].Message.Content
				if content != "" {
					onResponse(content)
				}
			}
			continue
		}

		// Handle streaming response
		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				onResponse(content)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response listModelsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Try to parse as array if object parsing fails
		var altResponse []modelInfo
		if altErr := json.Unmarshal(body, &altResponse); altErr != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		response.Data = altResponse
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
