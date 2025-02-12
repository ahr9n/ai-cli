package test

import (
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/ollama"
	"github.com/ahr9n/ollama-cli/test/utils"
)

func BenchmarkCreateChatCompletion(b *testing.B) {
	server := utils.NewTestServer(utils.CreateSuccessHandler(utils.MockResponses.ModelsList))
	defer server.Close()

	client := utils.NewTestClient(server.URL)
	messages := []ollama.Message{{Role: "user", Content: "test"}}
	opts := &ollama.ChatOptions{Model: "test", Temperature: 0.7}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateChatCompletion(messages, opts)
	}
}
