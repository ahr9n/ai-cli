package cli

import (
	"os"
	"testing"
	"time"

	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/session"
)

func setupSessionTestDir(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dir)
}

func seedSession(t *testing.T) {
	t.Helper()
	sess := &session.Session{
		Meta: session.Metadata{
			ID:        "abc12345",
			Title:     "Test conversation",
			Provider:  "ollama",
			Model:     "deepseek-r1:1.5b",
			CreatedAt: time.Date(2026, 4, 11, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 4, 11, 12, 5, 0, 0, time.UTC),
			Messages:  2,
		},
		Messages: []provider.Message{
			{Role: provider.RoleUser, Content: "Hello"},
			{Role: provider.RoleAssistant, Content: "Hi there!"},
		},
	}
	if err := session.Save(sess); err != nil {
		t.Fatal(err)
	}
}

func TestChatListCommand(t *testing.T) {
	setupSessionTestDir(t)
	seedSession(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("chat list error: %v", err)
	}
}

func TestChatListEmpty(t *testing.T) {
	setupSessionTestDir(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("chat list error: %v", err)
	}
}

func TestChatShowCommand(t *testing.T) {
	setupSessionTestDir(t)
	seedSession(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "show", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("chat show error: %v", err)
	}
}

func TestChatShowNonexistent(t *testing.T) {
	setupSessionTestDir(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "show", "nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent session")
	}
}

func TestChatDeleteCommand(t *testing.T) {
	setupSessionTestDir(t)
	seedSession(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "delete", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("chat delete error: %v", err)
	}

	// Verify deletion.
	_, err := session.Load("abc12345")
	if err == nil {
		t.Fatal("session should have been deleted")
	}
}

func TestChatExportMarkdown(t *testing.T) {
	setupSessionTestDir(t)
	seedSession(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "export", "abc12345"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("chat export error: %v", err)
	}
}

func TestChatExportJSON(t *testing.T) {
	setupSessionTestDir(t)
	seedSession(t)

	// Redirect stdout to verify JSON output.
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "export", "abc12345", "-f", "json"})
	err := cmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("chat export json error: %v", err)
	}

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if output == "" {
		t.Fatal("expected JSON output")
	}
}

func TestChatExportInvalidFormat(t *testing.T) {
	setupSessionTestDir(t)
	seedSession(t)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{"chat", "export", "abc12345", "-f", "xml"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}
