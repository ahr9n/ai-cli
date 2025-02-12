package ollama

import (
	"encoding/json"
	"fmt"
)

type ModelInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	Details  struct {
		Format  string `json:"format"`
		Family  string `json:"family"`
		License string `json:"license"`
	} `json:"details"`
}

type ListModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

func (c *Client) ListModels() ([]ModelInfo, error) {
	resp, err := c.DoGet("api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if err := c.HandleError(resp); err != nil {
		return nil, err
	}

	var response ListModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Models, nil
}
