// Package openaicompat implements the OpenAI-compatible chat completions API.
// Used as the foundation for OpenAI, Groq, LocalAI, and any compatible backend.
package openaicompat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ahr9n/ai-cli/internal/api"
	"github.com/ahr9n/ai-cli/internal/provider"
)

// Client implements provider.Provider for any OpenAI-compatible API.
type Client struct {
	base *api.BaseClient
	name string
	desc string
}

// NewClient creates a client for an OpenAI-compatible endpoint.
func NewClient(baseURL, apiKey string) *Client {
	var opts []api.Option
	if apiKey != "" {
		opts = append(opts, api.WithHeader("Authorization", "Bearer "+apiKey))
	}
	return &Client{
		base: api.NewBaseClient(baseURL, opts...),
		name: "OpenAI Compatible",
		desc: "OpenAI-compatible API",
	}
}

// WithInfo sets the display name and description. Returns the client for chaining.
func (c *Client) WithInfo(name, desc string) *Client {
	c.name = name
	c.desc = desc
	return c
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

func (c *Client) CreateCompletion(ctx context.Context, messages []provider.Message, opts *provider.CompletionOptions) (string, error) {
	var result strings.Builder
	err := c.StreamCompletion(ctx, messages, opts, func(chunk string) {
		result.WriteString(chunk)
	})
	return result.String(), err
}

func (c *Client) StreamCompletion(ctx context.Context, messages []provider.Message, opts *provider.CompletionOptions, onResponse func(string)) error {
	if len(messages) == 0 {
		return provider.ErrNoMessages
	}

	reqBody := completionRequest{
		Model:       opts.Model,
		Messages:    messages,
		Temperature: opts.Temperature,
		Stream:      true,
	}

	resp, err := c.base.DoPost(ctx, "chat/completions", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.base.CheckStatus(resp); err != nil {
		return err
	}

	scanner := bufio.NewScanner(resp.Body)
	const maxScanTokenSize = 10 * 1024 * 1024
	scanBuf := make([]byte, maxScanTokenSize)
	scanner.Buffer(scanBuf, maxScanTokenSize)

	var parseFailures int
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
				parseFailures++
				if parseFailures > 5 {
					return fmt.Errorf("too many parse failures in stream, last line: %q", line)
				}
				continue
			}
			parseFailures = 0
			if len(response.Choices) > 0 {
				content := response.Choices[0].Message.Content
				if content != "" {
					onResponse(content)
				}
			}
			continue
		}
		parseFailures = 0

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
	resp, err := c.base.DoGet(ctx, "models")
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if err := c.base.CheckStatus(resp); err != nil {
		return nil, err
	}

	// ReadAll because we may need to retry parsing as array vs object.
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

func (c *Client) Name() string        { return c.name }
func (c *Client) Description() string { return c.desc }
