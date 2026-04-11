package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ahr9n/ai-cli/internal/provider"
)

func setupTestDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
}

func newTestSession() *Session {
	return &Session{
		Meta: Metadata{
			ID:        "test1234",
			Title:     "What is Go?",
			Provider:  "ollama",
			Model:     "deepseek-r1:1.5b",
			CreatedAt: time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 4, 11, 12, 5, 0, 0, time.UTC),
			Messages:  2,
		},
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "What is Go?"},
			{Role: provider.RoleAssistant, Content: "Go is a statically typed language."},
		},
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestDir(t)
	sess := newTestSession()

	if err := Save(sess); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	loaded, err := Load("test1234")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if loaded.Meta.ID != "test1234" {
		t.Errorf("ID = %q, want %q", loaded.Meta.ID, "test1234")
	}
	if loaded.Meta.Title != "What is Go?" {
		t.Errorf("Title = %q, want %q", loaded.Meta.Title, "What is Go?")
	}
	if loaded.Meta.Provider != "ollama" {
		t.Errorf("Provider = %q, want %q", loaded.Meta.Provider, "ollama")
	}
	if len(loaded.Messages) != 2 {
		t.Fatalf("Messages = %d, want 2", len(loaded.Messages))
	}
	if loaded.Messages[0].Content != "What is Go?" {
		t.Errorf("Messages[0].Content = %q, want %q", loaded.Messages[0].Content, "What is Go?")
	}
	if loaded.Messages[1].Role != provider.RoleAssistant {
		t.Errorf("Messages[1].Role = %q, want %q", loaded.Messages[1].Role, provider.RoleAssistant)
	}
}

func TestLoadNonexistent(t *testing.T) {
	setupTestDir(t)
	_, err := Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent session")
	}
}

func TestList(t *testing.T) {
	setupTestDir(t)

	s1 := newTestSession()
	s1.Meta.ID = "aaa11111"
	s1.Meta.UpdatedAt = time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)

	s2 := newTestSession()
	s2.Meta.ID = "bbb22222"
	s2.Meta.Title = "How does concurrency work?"
	s2.Meta.UpdatedAt = time.Date(2026, 4, 11, 0, 0, 0, 0, time.UTC)

	if err := Save(s1); err != nil {
		t.Fatal(err)
	}
	if err := Save(s2); err != nil {
		t.Fatal(err)
	}

	sessions, err := List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("List returned %d, want 2", len(sessions))
	}
	// Sorted by UpdatedAt descending — bbb22222 should be first.
	if sessions[0].ID != "bbb22222" {
		t.Errorf("sessions[0].ID = %q, want %q", sessions[0].ID, "bbb22222")
	}
	if sessions[1].ID != "aaa11111" {
		t.Errorf("sessions[1].ID = %q, want %q", sessions[1].ID, "aaa11111")
	}
}

func TestListEmpty(t *testing.T) {
	setupTestDir(t)
	sessions, err := List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("List returned %d, want 0", len(sessions))
	}
}

func TestDelete(t *testing.T) {
	setupTestDir(t)
	sess := newTestSession()
	if err := Save(sess); err != nil {
		t.Fatal(err)
	}

	if err := Delete("test1234"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	path := filepath.Join(DataDir(), "test1234.jsonl")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("session file should have been deleted")
	}
}

func TestDeleteNonexistent(t *testing.T) {
	setupTestDir(t)
	err := Delete("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent session")
	}
}

func TestAppendMessage(t *testing.T) {
	setupTestDir(t)
	sess := newTestSession()
	if err := Save(sess); err != nil {
		t.Fatal(err)
	}

	newMsg := provider.Message{Role: provider.RoleUser, Content: "Tell me more"}
	if err := AppendMessage(sess, newMsg); err != nil {
		t.Fatalf("AppendMessage error: %v", err)
	}

	if sess.Meta.Messages != 3 {
		t.Errorf("Messages count = %d, want 3", sess.Meta.Messages)
	}

	// Reload and verify.
	loaded, err := Load("test1234")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(loaded.Messages) != 3 {
		t.Fatalf("loaded Messages = %d, want 3", len(loaded.Messages))
	}
	if loaded.Messages[2].Content != "Tell me more" {
		t.Errorf("Messages[2].Content = %q, want %q", loaded.Messages[2].Content, "Tell me more")
	}
}

func TestNewID(t *testing.T) {
	id := NewID()
	if len(id) != 8 {
		t.Errorf("NewID() length = %d, want 8", len(id))
	}
	// Should be hex characters.
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("NewID() contains non-hex character: %c", c)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"hi", 2, "hi"},
		{"hello", 3, "hel"},
		{"ab", 1, "a"},
	}
	for _, tt := range tests {
		got := Truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}
