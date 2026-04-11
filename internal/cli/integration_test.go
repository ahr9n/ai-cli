package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// mockProvider starts an httptest.Server implementing /v1/chat/completions and /v1/models.
func mockProvider(t *testing.T, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/chat/completions") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "text/event-stream")
			words := strings.Fields(response)
			for _, word := range words {
				chunk := map[string]any{
					"choices": []map[string]any{
						{"delta": map[string]string{"content": word + " "}},
					},
				}
				data, _ := json.Marshal(chunk)
				fmt.Fprintf(w, "data: %s\n\n", data)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
			fmt.Fprint(w, "data: [DONE]\n\n")

		case strings.HasSuffix(r.URL.Path, "/models") && r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]string{
					{"id": "test-model", "object": "model"},
				},
			})

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}
	}))
}

func TestSinglePromptEndToEnd(t *testing.T) {
	srv := mockProvider(t, "Hello from the test server")
	defer srv.Close()

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"custom", "--url", srv.URL + "/v1", "--model", "test", "say hello"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}
}

func TestJSONOutputEndToEnd(t *testing.T) {
	srv := mockProvider(t, "test response")
	defer srv.Close()

	// JSON mode writes to os.Stdout, so we capture it via a pipe.
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"custom", "--url", srv.URL + "/v1", "--model", "test", "--output", "json", "hello"})
	err := cmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	buf := make([]byte, 8192)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	var result map[string]any
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, output)
	}
	if result["error"] != nil {
		t.Errorf("unexpected error: %v", result["error"])
	}
	if result["response"] == "" {
		t.Error("expected non-empty response")
	}
}

func TestListModelsEndToEnd(t *testing.T) {
	srv := mockProvider(t, "")
	defer srv.Close()

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"custom", "--url", srv.URL + "/v1", "--list-models"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}
}

func TestUnicodeStreaming(t *testing.T) {
	// Response contains emoji, CJK, and Arabic text.
	srv := mockProvider(t, "Hello 🌍 你好 مرحبا")
	defer srv.Close()

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"custom", "--url", srv.URL + "/v1", "--model", "test", "unicode test"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}
}

func TestMaxTokensAndTopPFlags(t *testing.T) {
	srv := mockProvider(t, "short")
	defer srv.Close()

	cmd := NewRootCommand()
	cmd.SetArgs([]string{
		"custom", "--url", srv.URL + "/v1", "--model", "test",
		"--max-tokens", "100", "--top-p", "0.9",
		"hello",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command failed: %v", err)
	}
}

func TestConnectionRefused(t *testing.T) {
	cmd := NewRootCommand()
	// Use a port that's almost certainly not listening.
	cmd.SetArgs([]string{"custom", "--url", "http://127.0.0.1:1/v1", "--model", "test", "hello"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for connection refused")
	}
	if !strings.Contains(err.Error(), "connect") {
		t.Errorf("error = %v, expected connection-related message", err)
	}
}
