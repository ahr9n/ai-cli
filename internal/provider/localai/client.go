package localai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ahr9n/ai-cli/internal/api"
	"github.com/ahr9n/ai-cli/internal/provider"
)

type Client struct {
	base *api.BaseClient
}

type completionRequest struct {
	Model       string             `json:"model"`
	Messages    []provider.Message `json:"messages"`
	Temperature float32            `json:"temperature"`
	Stream      bool               `json:"stream"`
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

	reqBody := completionRequest{
		Model:       opts.Model,
		Messages:    messages,
		Temperature: opts.Temperature,
		Stream:      true,
	}

	resp, err := c.base.DoPost(ctx, "v1/chat/completions", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.base.HandleError(resp); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	const maxScanTokenSize = 10 * 1024 * 1024
	scanBuf := make([]byte, maxScanTokenSize)
	scanner.Buffer(scanBuf, maxScanTokenSize)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		}

		if line == "[DONE]" {
			continue
		}

		var streamResp streamingResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			var response completionResponse
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				continue
			}
			if len(response.Choices) > 0 {
				content := response.Choices[0].Message.Content
				if content != "" {
					onResponse(content)
				}
			}
			continue
		}

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

func (c *Client) ListModels(ctx context.Context) ([]provider.ModelInfo, error) {
	resp, err := c.base.DoGet(ctx, "v1/models")
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if err := c.base.HandleError(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response listModelsResponse
	if err := json.Unmarshal(body, &response); err != nil {
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
