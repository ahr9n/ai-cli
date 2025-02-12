package cli

import (
	"fmt"
	"strings"

	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/ahr9n/ollama-cli/pkg/ollama"
	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", Version)
	},
}

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

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(newModelsCommand())

	return cmd
}

func runChat(opts *ChatOptions, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := ollama.NewClient(cfg)

	if opts.Interactive {
		return runInteractiveMode(client, opts)
	}

	if len(args) == 0 {
		return fmt.Errorf("please provide a prompt or use -i for interactive mode")
	}

	prompt := strings.Join(args, " ")
	return handleSinglePrompt(client, prompt, opts)
}
