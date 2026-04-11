package cli

import (
	"fmt"

	"github.com/ahr9n/ai-cli/internal/prompt"
	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/spf13/cobra"
)

type ChatOptions struct {
	Interactive  bool
	Model        string
	Temperature  float32
	ProviderURL  string
	ListModels   bool
	SystemPrompt string
	MaxHistory   int
	PresetPrompt string
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai-cli",
		Short: "A CLI tool to interact with various AI models",
		Long: `A command-line interface for interacting with various AI providers and models.
Supports both single prompt and interactive conversations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleDefaultProvider(cmd, args)
		},
	}

	cmd.AddCommand(
		versionCommand(),
		listProvidersCommand(),
		newProviderCommand(provider.Ollama, "deepseek-r1:1.5b"),
		newProviderCommand(provider.LocalAI, "gpt-3.5-turbo"),
		newDefaultCommand(),
	)

	return cmd
}

func resolveSystemPrompt(opts *ChatOptions) error {
	if opts.SystemPrompt != "" {
		return nil
	}
	switch opts.PresetPrompt {
	case "creative":
		opts.SystemPrompt = prompt.Creative
	case "concise":
		opts.SystemPrompt = prompt.Concise
	case "code":
		opts.SystemPrompt = prompt.Code
	case "":
		opts.SystemPrompt = prompt.Default
	default:
		return fmt.Errorf("unknown preset prompt: %q (available: creative, concise, code)", opts.PresetPrompt)
	}
	return nil
}
