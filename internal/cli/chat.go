package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ahr9n/ai-cli/internal/loader"
	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/session"
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
	// --continue: load an existing session and resume interactive mode.
	if opts.ContinueID != "" {
		sess, err := session.Load(opts.ContinueID)
		if err != nil {
			return err
		}
		opts.Interactive = true
		opts.Model = sess.Meta.Model
		fmt.Printf("Resuming session %s (%d messages)\n", sess.Meta.ID, len(sess.Messages))
		return runInteractiveWithSession(ctx, p, opts, sess)
	}

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

	var buf strings.Builder
	var printed bool

	err := p.StreamCompletion(ctx, messages, &provider.CompletionOptions{
		Model:       opts.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		TopP:        opts.TopP,
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
		fmt.Println()
	}
	return nil
}

func runInteractiveMode(ctx context.Context, p provider.Provider, opts *ChatOptions) error {
	// Create a new session for auto-save.
	var sess *session.Session
	if !opts.NoSave {
		id := opts.SessionName
		if id == "" {
			id = session.NewID()
		}
		sess = &session.Session{
			Meta: session.Metadata{
				ID:           id,
				Provider:     string(opts.ProviderType),
				ProviderName: p.Name(),
				Model:        opts.Model,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}
	}
	return runInteractiveWithSession(ctx, p, opts, sess)
}

// runInteractiveWithSession runs the interactive chat loop.
// If sess is non-nil, messages are persisted after each exchange.
func runInteractiveWithSession(ctx context.Context, p provider.Provider, opts *ChatOptions, sess *session.Session) error {
	if opts.MaxHistory > 0 && opts.MaxHistory < 3 && opts.SystemPrompt != "" {
		opts.MaxHistory = 3
	}

	fmt.Printf("Starting interactive chat mode with %s (type 'exit' or 'quit' to quit, 'clear' to reset history)\n", p.Name())
	fmt.Printf("Model: %s\n", opts.Model)
	if opts.MaxHistory > 0 {
		fmt.Printf("Message history limit: %d messages\n", opts.MaxHistory)
	}
	if sess != nil {
		fmt.Printf("Session: %s\n", sess.Meta.ID)
	}

	// Initialize messages — from session if resuming, otherwise fresh.
	var messages []provider.Message
	if sess != nil && len(sess.Messages) > 0 {
		messages = append(messages, sess.Messages...)
	} else if opts.SystemPrompt != "" {
		messages = append(messages, provider.Message{
			Role:    provider.RoleSystem,
			Content: opts.SystemPrompt,
		})
		fmt.Printf("System prompt: %s\n", opts.SystemPrompt)
	}

	// Persist the initial state (system prompt) for new sessions.
	if sess != nil && len(sess.Messages) == 0 && len(messages) > 0 {
		sess.Messages = append(sess.Messages, messages...)
		sess.Meta.Messages = len(sess.Messages)
		_ = session.Save(sess)
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
			MaxTokens:   opts.MaxTokens,
			TopP:        opts.TopP,
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

		// Persist after each exchange.
		if sess != nil {
			if sess.Meta.Title == "" {
				sess.Meta.Title = session.Truncate(input, 60)
			}
			sess.Messages = messages
			sess.Meta.Messages = len(messages)
			sess.Meta.UpdatedAt = time.Now()
			if err := session.Save(sess); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to save session: %v\n", err)
			}
		}

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
