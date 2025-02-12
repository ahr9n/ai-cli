package test

import (
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/ollama"
	"github.com/ahr9n/ollama-cli/test/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	server := utils.SetupTestServer()
	defer server.Close()

	client := utils.SetupTestClient(server.URL)
	assert.NotNil(t, client)
}

func TestCreateChatCompletion(t *testing.T) {
	server := utils.SetupTestServer()
	defer server.Close()

	client := utils.SetupTestClient(server.URL)

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
