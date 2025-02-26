package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ahr9n/ai-cli/pkg/prompts"
	"github.com/ahr9n/ai-cli/pkg/provider"
	"github.com/ahr9n/ai-cli/pkg/utils"
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
	messages := []provider.Message{}

	if opts.SystemPrompt != "" {
		messages = append(messages, provider.Message{
			Role:    prompts.RoleSystem,
			Content: opts.SystemPrompt,
		})
	}
	messages = append(messages, provider.Message{
		Role:    prompts.RoleUser,
		Content: prompt,
	})

	loader := utils.InitLoader(utils.Dots)
	response, err := p.CreateCompletion(messages, &provider.CompletionOptions{
		Model:       opts.Model,
		Temperature: opts.Temperature,
	})

	loader.Stop()
	if err != nil {
		return fmt.Errorf("chat completion failed: %w", err)
	}

	fmt.Println(response)

	return nil
}

func runInteractiveMode(p provider.Provider, opts *ChatOptions) error {
	fmt.Printf("Starting interactive chat mode with %s (type 'exit' or 'quit' to quit, 'clear' to reset history)\n", p.Name())
	fmt.Printf("Model: %s\n", opts.Model)
	if opts.MaxHistory > 0 {
		fmt.Printf("Message history limit: %d messages\n", opts.MaxHistory)
	}

	var messages []provider.Message
	if opts.SystemPrompt != "" {
		messages = append(messages, provider.Message{
			Role:    prompts.RoleSystem,
			Content: opts.SystemPrompt,
		})
		fmt.Printf("System prompt: %s\n", opts.SystemPrompt)
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		switch strings.ToLower(input) {
		case "exit", "quit":
			return nil
		case "clear":
			if opts.SystemPrompt != "" {
				messages = []provider.Message{{
					Role:    prompts.RoleSystem,
					Content: opts.SystemPrompt,
				}}
			} else {
				messages = []provider.Message{}
			}
			fmt.Println("Conversation history cleared")
			continue
		}

		messages = append(messages, provider.Message{
			Role:    prompts.RoleUser,
			Content: input,
		})
		loader := utils.InitLoader(utils.Dots)

		var response strings.Builder
		err := p.StreamCompletion(messages, &provider.CompletionOptions{
			Model:       opts.Model,
			Temperature: opts.Temperature,
		}, func(chunk string) {
			if response.Len() == 0 {
				loader.Stop()
				fmt.Print("\nAssistant: ")
			}
			fmt.Print(chunk)
			response.WriteString(chunk)
		})
		loader.Stop()

		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			continue
		}
		fmt.Println()

		messages = append(messages, provider.Message{
			Role:    prompts.RoleAssistant,
			Content: response.String(),
		})

		if opts.MaxHistory > 0 && len(messages) > opts.MaxHistory {
			if messages[0].Role == prompts.RoleSystem {
				messages = append([]provider.Message{messages[0]}, messages[len(messages)-opts.MaxHistory+1:]...)
			} else {
				messages = messages[len(messages)-opts.MaxHistory:]
			}
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("input error: %w", scanner.Err())
	}

	return nil
}
