package cli

import (
	"bytes"
	"testing"
)

func TestRootCommandHelp(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("--help failed: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Fatal("expected help output, got empty")
	}
}

func TestVersionCommand(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"version"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}
}

func TestProvidersCommand(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{"providers"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("providers command failed: %v", err)
	}
}

func TestUnknownPresetReturnsError(t *testing.T) {
	opts := &ChatOptions{PresetPrompt: "nonexistent"}
	err := resolveSystemPrompt(opts)
	if err == nil {
		t.Fatal("expected error for unknown preset")
	}
}

func TestResolveSystemPrompt(t *testing.T) {
	tests := []struct {
		name     string
		preset   string
		explicit string
		wantErr  bool
	}{
		{name: "default", preset: "", wantErr: false},
		{name: "creative", preset: "creative", wantErr: false},
		{name: "concise", preset: "concise", wantErr: false},
		{name: "code", preset: "code", wantErr: false},
		{name: "explicit takes precedence", preset: "creative", explicit: "custom prompt", wantErr: false},
		{name: "unknown preset", preset: "bogus", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &ChatOptions{
				PresetPrompt: tt.preset,
				SystemPrompt: tt.explicit,
			}
			err := resolveSystemPrompt(opts)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && opts.SystemPrompt == "" {
				t.Error("SystemPrompt should be set after resolve")
			}
			if tt.explicit != "" && opts.SystemPrompt != tt.explicit {
				t.Errorf("explicit prompt not preserved: got %q, want %q", opts.SystemPrompt, tt.explicit)
			}
		})
	}
}

func TestDefaultSubcommands(t *testing.T) {
	for _, sub := range []string{"default set --help", "default show --help", "default clear --help"} {
		t.Run(sub, func(t *testing.T) {
			cmd := NewRootCommand()
			cmd.SetArgs([]string{sub})

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			// --help always succeeds
			_ = cmd.Execute()
		})
	}
}
