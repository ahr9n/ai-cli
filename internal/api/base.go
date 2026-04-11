package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Option configures a BaseClient.
type Option func(*BaseClient)

// WithHeader adds a default header to all requests.
func WithHeader(key, value string) Option {
	return func(c *BaseClient) {
		if value != "" {
			c.headers[key] = value
		}
	}
}

type BaseClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

func NewBaseClient(baseURL string, opts ...Option) *BaseClient {
	c := &BaseClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		headers:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *BaseClient) DoPost(ctx context.Context, path string, payload any) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", c.baseURL, path), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

func (c *BaseClient) DoGet(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.baseURL, path), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

// CheckStatus returns nil for 200 OK responses. For all other status codes,
// it drains the response body (for TCP connection reuse) and returns an *HTTPError.
func (c *BaseClient) CheckStatus(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return &HTTPError{StatusCode: resp.StatusCode, Body: body}
}
