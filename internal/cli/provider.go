package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/spf13/cobra"

	// Register providers via init().
	_ "github.com/ahr9n/ai-cli/internal/provider/localai"
	_ "github.com/ahr9n/ai-cli/internal/provider/ollama"
)

func listProvidersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "providers",
		Short: "List available AI providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			providers := provider.All()

			fmt.Println("Available Providers:")
			fmt.Println("--------------------")
			for _, p := range providers {
				fmt.Printf("\n%s\n", p.Name)
				fmt.Printf("Description: %s\n", p.Description)
				fmt.Printf("Default URL: %s\n", provider.DefaultURL(p.Type))
				fmt.Printf("Type: %s\n", p.Type)
			}

			return nil
		},
	}
}

func addCommonFlags(cmd *cobra.Command, opts *ChatOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&opts.Interactive, "interactive", "i", false, "Start interactive chat mode")
	flags.Float32VarP(&opts.Temperature, "temperature", "t", 0.7, "Sampling temperature (0.0-2.0)")
	flags.BoolVar(&opts.ListModels, "list-models", false, "List available models")
	flags.StringVarP(&opts.SystemPrompt, "system", "s", "", "System prompt to set the assistant's behavior")
	flags.IntVar(&opts.MaxHistory, "max-history", 20, "Maximum conversation history to keep (0 = unlimited)")
	flags.StringVarP(&opts.PresetPrompt, "preset", "p", "", "Use a preset system prompt (creative, concise, code)")
}

func newProviderCommand(providerType provider.ProviderType, defaultModel string) *cobra.Command {
	opts := &ChatOptions{}

	cmd := &cobra.Command{
		Use:   string(providerType) + " [prompt]",
		Short: fmt.Sprintf("Use %s provider", providerType),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := provider.New(providerType, opts.ProviderURL)
			if err != nil {
				return err
			}
			if opts.ListModels {
				return displayModels(cmd.Context(), p)
			}
			if err := resolveSystemPrompt(opts); err != nil {
				return err
			}
			return runChat(cmd.Context(), p, opts, args)
		},
	}

	addCommonFlags(cmd, opts)
	cmd.Flags().StringVarP(&opts.Model, "model", "m", defaultModel, "Model to use")
	cmd.Flags().StringVarP(&opts.ProviderURL, "url", "u", provider.DefaultURL(providerType), "Provider API URL")

	return cmd
}

func displayModels(ctx context.Context, p provider.Provider) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	models, err := p.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	if len(models) == 0 {
		fmt.Printf("No models found for %s\n", p.Name())
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "Available models for %s:\n\n", p.Name())
	fmt.Fprintln(w, "NAME\tSIZE\tFAMILY\tMODIFIED")
	fmt.Fprintln(w, "----\t----\t------\t--------")

	for _, model := range models {
		size := formatSize(model.Size)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			model.Name,
			size,
			model.Family,
			model.Modified,
		)
	}
	return w.Flush()
}

func formatSize(bytes int64) string {
	if bytes == 0 {
		return "-"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
