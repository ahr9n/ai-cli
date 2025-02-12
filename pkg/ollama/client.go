package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	baseURL    string
	httpClient *http.Client
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
		baseURL: cfg.OllamaURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) StreamChatCompletion(messages []Message, opts *ChatOptions, onResponse func(string)) error {
	prompt := messages[len(messages)-1].Content

	reqBody := generateRequest{
		Model:       opts.Model,
		Prompt:      prompt,
		System:      prompts.DefaultSystem(),
		Temperature: opts.Temperature,
		Stream:      true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("model '%s' not found - try running: ollama pull %s", opts.Model, opts.Model)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed (status %d): %s", resp.StatusCode, string(body))
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
