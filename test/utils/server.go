package utils

import (
	"fmt"
	"github.com/ahr9n/ai-cli/pkg/provider"
	"net/http"
	"net/http/httptest"

	"github.com/ahr9n/ai-cli/pkg/provider/localai"
	"github.com/ahr9n/ai-cli/pkg/provider/ollama"
)

var MockResponses = struct {
	SuccessfulChat  string
	ModelsList      string
	ModelNotFound   string
	InvalidJSON     string
	InternalError   string
	EmptyModelsList string
}{
	SuccessfulChat: `{"response": "Test response", "done": true}`,
	ModelsList: `{
		"models": [
			{
				"name": "model1",
				"size": 1234567,
				"modified": "2024-02-12",
				"details": {
					"format": "gguf",
					"family": "llama",
					"license": "mit"
				}
			},
			{
				"name": "model2",
				"size": 7654321,
				"modified": "2024-02-13",
				"details": {
					"format": "gguf",
					"family": "mistral",
					"license": "apache"
				}
			}
		]
	}`,
	ModelNotFound:   `{"error": "model not found"}`,
	InvalidJSON:     `{"invalid": json`,
	InternalError:   `{"error": "internal server error"}`,
	EmptyModelsList: `{"models": []}`,
}

func NewTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func NewTestClient(providerType provider.ProviderType, serverURL string) (provider.Provider, error) {
	switch providerType {
	case provider.Ollama:
		return ollama.NewClient(serverURL), nil
	case provider.LocalAI:
		return localai.NewClient(serverURL), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
}

func CreateSuccessHandler(response string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

func CreateErrorHandler(statusCode int, response string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	}
}

func CreateStreamHandler(responses []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		for _, resp := range responses {
			w.Write([]byte(resp + "\n"))
			w.(http.Flusher).Flush()
		}
	}
}
