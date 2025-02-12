package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/ahr9n/ollama-cli/pkg/ollama"
	"github.com/ahr9n/ollama-cli/pkg/prompts"
)

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
			Role:    prompts.RoleUser,
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
			Role:    prompts.RoleUser,
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
			Role:    prompts.RoleAssistant,
			Content: response.String(),
		})
	}

	return nil
}
