package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/ahr9n/ai-cli/internal/session"
	"github.com/spf13/cobra"
)

func newChatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Manage saved conversations",
	}
	cmd.AddCommand(
		chatListCommand(),
		chatShowCommand(),
		chatContinueCommand(),
		chatDeleteCommand(),
		chatExportCommand(),
	)
	return cmd
}

func chatListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List saved conversations",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, err := session.List()
			if err != nil {
				return err
			}
			if len(sessions) == 0 {
				fmt.Println("No saved conversations")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tPROVIDER\tMODEL\tMSGS\tUPDATED")
			fmt.Fprintln(w, "--\t-----\t--------\t-----\t----\t-------")
			for _, s := range sessions {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n",
					s.ID,
					session.Truncate(s.Title, 40),
					s.Provider,
					s.Model,
					s.Messages,
					s.UpdatedAt.Format("2006-01-02 15:04"),
				)
			}
			return w.Flush()
		},
	}
}

func chatShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <session-id>",
		Short: "Display a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, err := session.Load(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Session: %s\n", sess.Meta.ID)
			providerDisplay := sess.Meta.ProviderName
			if providerDisplay == "" {
				providerDisplay = sess.Meta.Provider
			}
			fmt.Printf("Provider: %s | Model: %s\n", providerDisplay, sess.Meta.Model)
			fmt.Printf("Created: %s | Updated: %s\n",
				sess.Meta.CreatedAt.Format("2006-01-02 15:04"),
				sess.Meta.UpdatedAt.Format("2006-01-02 15:04"),
			)
			fmt.Println(strings.Repeat("-", 60))

			for _, msg := range sess.Messages {
				switch msg.Role {
				case provider.RoleSystem:
					fmt.Printf("\n[System]: %s\n", msg.Content)
				case provider.RoleUser:
					fmt.Printf("\nYou: %s\n", msg.Content)
				case provider.RoleAssistant:
					fmt.Printf("\nAssistant: %s\n", msg.Content)
				}
			}
			return nil
		},
	}
}

func chatContinueCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "continue <session-id>",
		Short: "Continue a saved conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, err := session.Load(args[0])
			if err != nil {
				return err
			}

			providerType := provider.ProviderType(sess.Meta.Provider)
			p, err := provider.New(providerType, "")
			if err != nil {
				return fmt.Errorf("cannot create provider %q: %w", sess.Meta.Provider, err)
			}

			opts := &ChatOptions{
				Interactive: true,
				Model:       sess.Meta.Model,
				MaxHistory:  20,
			}

			fmt.Printf("Resuming session %s (%d messages)\n", sess.Meta.ID, len(sess.Messages))
			return runInteractiveWithSession(cmd.Context(), p, opts, sess)
		},
	}
}

func chatDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <session-id>",
		Short: "Delete a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := session.Delete(args[0]); err != nil {
				return err
			}
			fmt.Printf("Session %s deleted\n", args[0])
			return nil
		},
	}
}

func chatExportCommand() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export <session-id>",
		Short: "Export a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, err := session.Load(args[0])
			if err != nil {
				return err
			}
			switch format {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(sess)
			case "markdown":
				return exportMarkdown(sess)
			default:
				return fmt.Errorf("unknown format: %q (available: markdown, json)", format)
			}
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "Export format (markdown, json)")
	cmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"markdown", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

func exportMarkdown(sess *session.Session) error {
	fmt.Printf("# %s\n\n", sess.Meta.Title)
	providerDisplay := sess.Meta.ProviderName
	if providerDisplay == "" {
		providerDisplay = sess.Meta.Provider
	}
	fmt.Printf("**Provider:** %s | **Model:** %s\n\n", providerDisplay, sess.Meta.Model)
	fmt.Println("---")
	fmt.Println()

	for _, msg := range sess.Messages {
		switch msg.Role {
		case provider.RoleSystem:
			fmt.Printf("*System: %s*\n\n", msg.Content)
		case provider.RoleUser:
			fmt.Printf("## User\n\n%s\n\n", msg.Content)
		case provider.RoleAssistant:
			fmt.Printf("## Assistant\n\n%s\n\n", msg.Content)
		}
	}
	return nil
}
