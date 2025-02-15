package cli

import (
	"github.com/spf13/cobra"
)

type ChatOptions struct {
	Interactive bool
	Model       string
	Temperature float32
	ProviderURL string
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai-cli",
		Short: "A CLI tool to interact with various AI models",
		Long: `A command-line interface for interacting with various AI providers and models.
Supports both single prompts and interactive conversations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
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
