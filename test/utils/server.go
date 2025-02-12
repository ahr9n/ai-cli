package utils

import (
	"net/http"
	"net/http/httptest"

	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/ahr9n/ollama-cli/pkg/ollama"
)

func SetupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"response": "Test response", "done": true}`))
	}))
}

func SetupTestClient(serverURL string) *ollama.Client {
	cfg := &config.Config{
		OllamaURL: serverURL,
	}
	return ollama.NewClient(cfg)
}
