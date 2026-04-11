package cli

import (
	"testing"
)

func TestCompletionBash(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"completion", "bash"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("completion bash error: %v", err)
	}
}

func TestCompletionZsh(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"completion", "zsh"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("completion zsh error: %v", err)
	}
}

func TestCompletionFish(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"completion", "fish"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("completion fish error: %v", err)
	}
}

func TestCompletionInvalidShell(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"completion", "invalid"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid shell")
	}
}
