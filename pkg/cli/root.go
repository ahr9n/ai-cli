package cli

import (
	"github.com/spf13/cobra"
)

type ChatOptions struct {
	Interactive bool
	Model       string
	Temperature float32
}

func NewRootCommand() *cobra.Command {
	opts := &ChatOptions{}

	cmd := &cobra.Command{
		Use:   "ollama-cli [prompt]",
		Short: "A CLI tool to interact with Ollama LLMs",
		Long: `A command-line interface for interacting with Ollama's language models.
Supports both single prompts and interactive conversations.

Before using, make sure to:
1. Install Ollama (https://ollama.ai)
2. Run the Ollama service
3. Pull your desired model (e.g., 'ollama pull deepseek-r1:1.5b')`,
		Example: `  ollama-cli "What is the capital of France?"
  ollama-cli -i  # Start interactive mode
  ollama-cli --model mistral "Explain quantum computing"`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChat(opts, args)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Interactive, "interactive", "i", false, "Start interactive chat mode")
	flags.StringVarP(&opts.Model, "model", "m", "deepseek-r1:1.5b", "Model to use (deepseek-r1:1.5b, llama2, mistral, etc.)")
	flags.Float32VarP(&opts.Temperature, "temperature", "t", 0.7, "Sampling temperature (0.0-2.0)")

	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newModelsCommand())

	return cmd
}
