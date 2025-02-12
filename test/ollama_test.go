package test

import (
	"net/http"
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/ollama"
	"github.com/ahr9n/ollama-cli/test/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.ModelsList))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	assert.NotNil(t, client)
}


func TestListModels_Success(t *testing.T) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.ModelsList))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	models, err := client.ListModels()

	assert.NoError(t, err)
	assert.Len(t, models, 2)
	assert.Equal(t, "model1", models[0].Name)
	assert.Equal(t, "model2", models[1].Name)
	assert.Equal(t, "llama", models[0].Details.Family)
	assert.Equal(t, "mistral", models[1].Details.Family)
}

func TestListModels_EmptyList(t *testing.T) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.EmptyModelsList))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	models, err := client.ListModels()

	assert.NoError(t, err)
	assert.Empty(t, models)
}

func TestListModels_ServerError(t *testing.T) {
	server := utils.NewTestServer(utils.CreateErrorHandler(http.StatusInternalServerError, utils.MockResponses.InternalError))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	models, err := client.ListModels()

	assert.Error(t, err)
	assert.Nil(t, models)
	assert.Contains(t, err.Error(), "request failed (status 500)")
}

func TestListModels_InvalidJSON(t *testing.T) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.InvalidJSON))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	models, err := client.ListModels()

	assert.Error(t, err)
	assert.Nil(t, models)
	assert.Contains(t, err.Error(), "failed to decode")
}

func TestCreateChatCompletion_Success(t *testing.T) {
	server := utils.NewTestServer(utils.CreateStreamHandler([]string{
		`{"response": "Test ", "done": false}`,
		`{"response": "response", "done": true}`,
	}))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	messages := []ollama.Message{{Role: "user", Content: "test"}}
	opts := &ollama.ChatOptions{Model: "test-model", Temperature: 0.7}

	response, err := client.CreateChatCompletion(messages, opts)

	assert.NoError(t, err)
	assert.Equal(t, "Test response", response)
}

func TestCreateChatCompletion_ModelNotFound(t *testing.T) {
	server := utils.NewTestServer(utils.CreateErrorHandler(http.StatusNotFound, utils.MockResponses.ModelNotFound))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	messages := []ollama.Message{{Role: "user", Content: "test"}}
	opts := &ollama.ChatOptions{Model: "nonexistent-model", Temperature: 0.7}

	_, err := client.CreateChatCompletion(messages, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "model 'nonexistent-model' not found")
}

func TestCreateChatCompletion_ServerError(t *testing.T) {
	server := utils.NewTestServer(utils.CreateErrorHandler(http.StatusInternalServerError, utils.MockResponses.InternalError))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	messages := []ollama.Message{{Role: "user", Content: "test"}}
	opts := &ollama.ChatOptions{Model: "test-model", Temperature: 0.7}

	response, err := client.CreateChatCompletion(messages, opts)

	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "request failed (status 500)")
}

func TestCreateChatCompletion_EmptyMessages(t *testing.T) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.SuccessfulChat))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	messages := []ollama.Message{}
	opts := &ollama.ChatOptions{Model: "test-model", Temperature: 0.7}

	_, err := client.CreateChatCompletion(messages, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no messages provided")
}

func TestCreateChatCompletion_InvalidOptions(t *testing.T) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.SuccessfulChat))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	messages := []ollama.Message{} // Empty messages
	opts := &ollama.ChatOptions{Model: "test-model", Temperature: 0.7}

	response, err := client.CreateChatCompletion(messages, opts)

	assert.Error(t, err)
	assert.Empty(t, response)
}
