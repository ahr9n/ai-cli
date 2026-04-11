package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ahr9n/ai-cli/internal/loader"
	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/render"
)

// readStdin reads all data from r if it is not a terminal (i.e. piped input).
// Returns empty string when r is a terminal or on read error.
func readStdin(r *os.File) string {
	info, err := r.Stat()
	if err != nil {
		return ""
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return "" // interactive terminal, not a pipe
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func runChat(ctx context.Context, p provider.Provider, opts *ChatOptions, args []string) error {
	if opts.Interactive {
		return runInteractiveMode(ctx, p, opts)
	}

	piped := readStdin(os.Stdin)
	prompt := strings.Join(args, " ")

	switch {
	case prompt != "" && piped != "":
		prompt = prompt + "\n\n" + piped
	case piped != "":
		prompt = piped
	}

	if prompt == "" {
		return fmt.Errorf("please provide a prompt or use -i for interactive mode")
	}

	return handleSinglePrompt(ctx, p, prompt, opts)
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

	jsonMode := opts.Output == "json"

	l := loader.New(loader.Dots)
	if !jsonMode {
		l.Start()
	}
	defer l.Stop()

	r := render.New()
	var buf strings.Builder
	var printed bool

	err := p.StreamCompletion(ctx, messages, &provider.CompletionOptions{
		Model:       opts.Model,
		Temperature: opts.Temperature,
	}, func(chunk string) {
		l.Stop()
		if !jsonMode {
			fmt.Print(chunk)
		}
		buf.WriteString(chunk)
		printed = true
	})

	if jsonMode {
		return emitJSON(os.Stdout, p.Name(), opts.Model, userPrompt, buf.String(), err)
	}

	if err != nil {
		return fmt.Errorf("chat completion failed: %w", err)
	}

	if printed {
		if r.IsTerminal() && buf.Len() > 0 {
			// Clear raw streaming output, reprint with Glamour rendering.
			lines := strings.Count(buf.String(), "\n") + 1
			for i := 0; i < lines; i++ {
				fmt.Print("\033[A\033[2K")
			}
			fmt.Print("\r")
			fmt.Println(r.RenderBlock(buf.String()))
		} else {
			fmt.Println()
		}
	}
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

	r := render.New()

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

		raw := response.String()
		if r.IsTerminal() && raw != "" {
			// Clear raw streaming output, reprint rendered.
			lines := strings.Count(raw, "\n") + 1
			// +1 for the "Assistant: " prefix line.
			lines++
			for i := 0; i < lines; i++ {
				fmt.Print("\033[A\033[2K")
			}
			fmt.Print("\r")
			fmt.Printf("\nAssistant: %s\n", r.RenderBlock(raw))
		} else {
			fmt.Println()
		}

		messages = append(messages, provider.Message{
			Role:    provider.RoleAssistant,
			Content: raw,
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

// emitJSON writes a structured JSON response to w.
// When completionErr is non-nil, the error field is populated and the function returns nil
// (the error has been reported in the JSON output).
func emitJSON(w io.Writer, providerName, model, prompt, response string, completionErr error) error {
	result := struct {
		Provider string  `json:"provider"`
		Model    string  `json:"model"`
		Prompt   string  `json:"prompt"`
		Response string  `json:"response"`
		Error    *string `json:"error"`
	}{
		Provider: providerName,
		Model:    model,
		Prompt:   prompt,
		Response: response,
	}

	if completionErr != nil {
		errMsg := completionErr.Error()
		result.Error = &errMsg
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
