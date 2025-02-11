package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/ahr9n/ollama-cli/pkg/ollama"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		OllamaURL: "http://localhost:11434",
	}

	client := ollama.NewClient(cfg)
	assert.NotNil(t, client)
}

func TestCreateChatCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/generate", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response": "Test response", "done": true}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		OllamaURL: server.URL,
	}
	client := ollama.NewClient(cfg)

	messages := []ollama.Message{
		{
			Role:    "user",
			Content: "Test prompt",
		},
	}

	opts := &ollama.ChatOptions{
		Model:       "test-model",
		Temperature: 0.7,
	}

	response, err := client.CreateChatCompletion(messages, opts)
	assert.NoError(t, err)
	assert.Equal(t, "Test response", response)
}
