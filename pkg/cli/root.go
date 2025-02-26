package cli

import (
	"github.com/ahr9n/ai-cli/pkg/prompts"
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
Supports both single prompts and interactive conversations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return HandleDefaultProvider(cmd, args)
		},
	}

	cmd.AddCommand(
		versionCommand(),
		listProvidersCommand(),
		newOllamaCommand(),
		newLocalAICommand(),
		newDefaultCommand(),
	)

	return cmd
}

func resolveSystemPrompt(opts *ChatOptions) {
	if opts.SystemPrompt != "" {
		return
	}

	switch opts.PresetPrompt {
	case "creative":
		opts.SystemPrompt = prompts.CreativeSystem()
	case "concise":
		opts.SystemPrompt = prompts.ConciseSystem()
	case "code":
		opts.SystemPrompt = prompts.CodeSystem()
	case "":
		opts.SystemPrompt = prompts.DefaultSystem()
	}
}
