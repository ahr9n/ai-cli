package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ahr9n/ai-cli/pkg/prompts"
	"github.com/ahr9n/ai-cli/pkg/provider"
)

func runChat(p provider.Provider, opts *ChatOptions, args []string) error {
	if opts.Interactive {
		return runInteractiveMode(p, opts)
	}

	if len(args) == 0 {
		return fmt.Errorf("please provide a prompt or use -i for interactive mode")
	}

	prompt := strings.Join(args, " ")
	return handleSinglePrompt(p, prompt, opts)
}

func handleSinglePrompt(p provider.Provider, prompt string, opts *ChatOptions) error {
	messages := []provider.Message{
		{
			Role:    prompts.RoleUser,
			Content: prompt,
		},
	}

	fmt.Print("thinking...")

	response, err := p.CreateCompletion(messages, &provider.CompletionOptions{
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

func runInteractiveMode(p provider.Provider, opts *ChatOptions) error {
	fmt.Printf("Starting interactive chat mode with %s (type 'exit' to quit)\n", p.Name())
	fmt.Printf("Model: %s\n", opts.Model)

	var messages []provider.Message
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		messages = append(messages, provider.Message{
			Role:    prompts.RoleUser,
			Content: input,
		})

		fmt.Print("\nthinking...")

		var response strings.Builder
		err := p.StreamCompletion(messages, &provider.CompletionOptions{
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

		messages = append(messages, provider.Message{
			Role:    prompts.RoleAssistant,
			Content: response.String(),
		})
	}

	return nil
}
