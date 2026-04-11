package openaicompat

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahr9n/ai-cli/internal/provider"
)

func sseLines(chunks ...string) string {
	var b strings.Builder
	for _, c := range chunks {
		b.WriteString(fmt.Sprintf("data: {\"choices\":[{\"delta\":{\"content\":%q}}]}\n\n", c))
	}
	b.WriteString("data: [DONE]\n\n")
	return b.String()
}

func TestStreamCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, sseLines("Hello", " World"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	var got strings.Builder
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	err := c.StreamCompletion(context.Background(), msgs, &provider.CompletionOptions{
		Model: "test", Temperature: 0.7,
	}, func(chunk string) {
		got.WriteString(chunk)
	})
	if err != nil {
		t.Fatalf("StreamCompletion error: %v", err)
	}
	if got.String() != "Hello World" {
		t.Errorf("got %q, want %q", got.String(), "Hello World")
	}
}

func TestStreamCompletionNoMessages(t *testing.T) {
	c := NewClient("http://unused", "")
	err := c.StreamCompletion(context.Background(), nil, &provider.CompletionOptions{
		Model: "test",
	}, func(string) {})
	if err != provider.ErrNoMessages {
		t.Errorf("error = %v, want ErrNoMessages", err)
	}
}

func TestAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-key")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-key")
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	_ = c.StreamCompletion(context.Background(), msgs, &provider.CompletionOptions{
		Model: "test",
	}, func(string) {})
}

func TestNoAuthHeaderWhenEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("Authorization = %q, want empty", auth)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	_ = c.StreamCompletion(context.Background(), msgs, &provider.CompletionOptions{
		Model: "test",
	}, func(string) {})
}

func TestHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	err := c.StreamCompletion(context.Background(), msgs, &provider.CompletionOptions{
		Model: "test",
	}, func(string) {})
	if err == nil {
		t.Fatal("expected error for 404")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %v, expected to contain 404", err)
	}
}

func TestContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, sseLines("Hello"))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := NewClient(srv.URL, "")
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	err := c.StreamCompletion(ctx, msgs, &provider.CompletionOptions{
		Model: "test",
	}, func(string) {})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestParseFailureThreshold(t *testing.T) {
	// Send 6 unparseable lines to exceed the threshold of 5.
	var body strings.Builder
	for i := 0; i < 6; i++ {
		body.WriteString("data: {invalid json}\n\n")
	}
	body.WriteString("data: [DONE]\n\n")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, body.String())
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	err := c.StreamCompletion(context.Background(), msgs, &provider.CompletionOptions{
		Model: "test",
	}, func(string) {})
	if err == nil {
		t.Fatal("expected error for too many parse failures")
	}
	if !strings.Contains(err.Error(), "too many parse failures") {
		t.Errorf("error = %v, expected 'too many parse failures'", err)
	}
}

func TestListModels(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id":"gpt-4o","object":"model"},{"id":"gpt-3.5-turbo","object":"model"}]}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels error: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("got %d models, want 2", len(models))
	}
	if models[0].Name != "gpt-4o" {
		t.Errorf("models[0].Name = %q, want %q", models[0].Name, "gpt-4o")
	}
}

func TestListModelsArrayResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":"model-a"},{"id":"model-b"}]`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	models, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels error: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("got %d models, want 2", len(models))
	}
	if models[1].Name != "model-b" {
		t.Errorf("models[1].Name = %q, want %q", models[1].Name, "model-b")
	}
}

func TestWithInfo(t *testing.T) {
	c := NewClient("http://unused", "").WithInfo("MyProvider", "My description")
	if c.Name() != "MyProvider" {
		t.Errorf("Name() = %q, want %q", c.Name(), "MyProvider")
	}
	if c.Description() != "My description" {
		t.Errorf("Description() = %q, want %q", c.Description(), "My description")
	}
}

func TestNonStreamingFallback(t *testing.T) {
	body := `data: {"choices":[{"message":{"content":"hello"}}]}` + "\n\n" +
		"data: [DONE]\n\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, body)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	var got strings.Builder
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	err := c.StreamCompletion(context.Background(), msgs, &provider.CompletionOptions{Model: "test"}, func(chunk string) {
		got.WriteString(chunk)
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got.String() != "hello" {
		t.Errorf("got %q, want %q", got.String(), "hello")
	}
}

func TestCreateCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, sseLines("Hello", " World"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "")
	msgs := []provider.Message{{Role: provider.RoleUser, Content: "hi"}}
	got, err := c.CreateCompletion(context.Background(), msgs, &provider.CompletionOptions{
		Model: "test",
	})
	if err != nil {
		t.Fatalf("CreateCompletion error: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("got %q, want %q", got, "Hello World")
	}
}
