package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/ahr9n/ollama-cli/pkg/ollama"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return runChat(opts, args)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.Interactive, "interactive", "i", false, "Start interactive chat mode")
	flags.StringVarP(&opts.Model, "model", "m", "deepseek-r1:1.5b", "Model to use (deepseek-r1:1.5b, llama2, mistral, etc.)")
	flags.Float32VarP(&opts.Temperature, "temperature", "t", 0.7, "Sampling temperature (0.0-2.0)")

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

func handleSinglePrompt(client *ollama.Client, prompt string, opts *ChatOptions) error {
	messages := []ollama.Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	fmt.Print("thinking...")

	response, err := client.CreateChatCompletion(messages, &ollama.ChatOptions{
		Model:       opts.Model,
		Temperature: opts.Temperature,
	})

	if err != nil {
		return fmt.Errorf("chat completion failed: %w", err)
	}

	fmt.Print("\r\033[K")
	fmt.Println(response)
	return nil
}

func runInteractiveMode(client *ollama.Client, opts *ChatOptions) error {
	fmt.Println("Starting interactive chat mode (type 'exit' to quit)")
	fmt.Println("Model:", opts.Model)

	scanner := bufio.NewScanner(os.Stdin)
	messages := []ollama.Message{}

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		messages = append(messages, ollama.Message{
			Role:    "user",
			Content: input,
		})

		fmt.Print("\nthinking...")

		var response strings.Builder
		err := client.StreamChatCompletion(messages, &ollama.ChatOptions{
			Model:       opts.Model,
			Temperature: opts.Temperature,
		}, func(chunk string) {
			if response.Len() == 0 {
				fmt.Print("\r\033[K")
				fmt.Print("\nAssistant: ")
			}
			fmt.Print(chunk)
			response.WriteString(chunk)
		})

		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			continue
		}
		fmt.Println()

		messages = append(messages, ollama.Message{
			Role:    "assistant",
			Content: response.String(),
		})
	}

	return nil
}
