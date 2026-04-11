package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ahr9n/ai-cli/internal/loader"
	"github.com/ahr9n/ai-cli/internal/provider"
)

func runChat(ctx context.Context, p provider.Provider, opts *ChatOptions, args []string) error {
	if opts.Interactive {
		return runInteractiveMode(ctx, p, opts)
	}

	if len(args) == 0 {
		return fmt.Errorf("please provide a prompt or use -i for interactive mode")
	}

	userPrompt := strings.Join(args, " ")
	return handleSinglePrompt(ctx, p, userPrompt, opts)
}

func handleSinglePrompt(ctx context.Context, p provider.Provider, userPrompt string, opts *ChatOptions) error {
	var messages []provider.Message

	if opts.SystemPrompt != "" {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
			Content: opts.SystemPrompt,
		})
	}
	messages = append(messages, provider.Message{
		Role:    provider.RoleUser,
		Content: userPrompt,
	})

	l := loader.New(loader.Dots)
	l.Start()
	defer l.Stop()

	err := p.StreamCompletion(ctx, messages, &provider.CompletionOptions{
		Model:       opts.Model,
		Temperature: opts.Temperature,
	}, func(chunk string) {
		l.Stop()
		fmt.Print(chunk)
	})
	if err != nil {
		return fmt.Errorf("chat completion failed: %w", err)
	}

	fmt.Println()
	return nil
}

func runInteractiveMode(ctx context.Context, p provider.Provider, opts *ChatOptions) error {
	if opts.MaxHistory > 0 && opts.MaxHistory < 3 && opts.SystemPrompt != "" {
		opts.MaxHistory = 3
	}

	fmt.Printf("Starting interactive chat mode with %s (type 'exit' or 'quit' to quit, 'clear' to reset history)\n", p.Name())
	fmt.Printf("Model: %s\n", opts.Model)
	if opts.MaxHistory > 0 {
		fmt.Printf("Message history limit: %d messages\n", opts.MaxHistory)
	}

	var messages []provider.Message
	if opts.SystemPrompt != "" {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
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
					Role:    provider.RoleSystem,
					Content: opts.SystemPrompt,
				}}
			} else {
				messages = nil
			}
			fmt.Println("Conversation history cleared")
			continue
		}

		messages = append(messages, provider.Message{
			Role:    provider.RoleUser,
			Content: input,
		})

		l := loader.New(loader.Dots)
		l.Start()

		var response strings.Builder
		err := p.StreamCompletion(ctx, messages, &provider.CompletionOptions{
			Model:       opts.Model,
			Temperature: opts.Temperature,
		}, func(chunk string) {
			if response.Len() == 0 {
				l.Stop()
				fmt.Print("\nAssistant: ")
			}
			fmt.Print(chunk)
			response.WriteString(chunk)
		})
		l.Stop()

		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			continue
		}
		fmt.Println()

		messages = append(messages, provider.Message{
			Role:    provider.RoleAssistant,
			Content: response.String(),
		})

		if opts.MaxHistory > 0 && len(messages) > opts.MaxHistory {
			if messages[0].Role == provider.RoleSystem {
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
