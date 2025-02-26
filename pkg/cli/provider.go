package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ahr9n/ai-cli/pkg/provider"
	"github.com/ahr9n/ai-cli/pkg/provider/localai"
	"github.com/ahr9n/ai-cli/pkg/provider/ollama"
	"github.com/spf13/cobra"
)

func AvailableProvidersList() []struct {
	Type        provider.ProviderType
	Name        string
	Description string
} {
	return []struct {
		Type        provider.ProviderType
		Name        string
		Description string
	}{
		{
			Type:        provider.Ollama,
			Name:        "Ollama",
			Description: "Run large language models locally",
		},
		{
			Type:        provider.LocalAI,
			Name:        "LocalAI",
			Description: "Self-hosted AI model server compatible with OpenAI's API",
		},
	}
}

func listProvidersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "providers",
		Short: "List available AI providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			providers := AvailableProvidersList()

			fmt.Println("Available Providers:")
			fmt.Println("--------------------")
			for _, p := range providers {
				fmt.Printf("\n%s\n", p.Name)
				fmt.Printf("Description: %s\n", p.Description)
				fmt.Printf("Default URL: %s\n", provider.DefaultURLs[p.Type])
				fmt.Printf("Type: %s\n", p.Type)
			}

			return nil
		},
	}
}

func NewProvider(providerType provider.ProviderType, baseURL string) (provider.Provider, error) {
	if baseURL == "" {
		baseURL = provider.DefaultURLs[providerType]
	}

	switch providerType {
	case provider.Ollama:
		return ollama.NewClient(baseURL), nil
	case provider.LocalAI:
		return localai.NewClient(baseURL), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
}

func addCommonFlags(cmd *cobra.Command, opts *ChatOptions) {
	flags := cmd.Flags()
	flags.BoolVarP(&opts.Interactive, "interactive", "i", false, "Start interactive chat mode")
	flags.Float32VarP(&opts.Temperature, "temperature", "t", 0.7, "Sampling temperature (0.0-2.0)")
	flags.BoolVar(&opts.ListModels, "list-models", false, "List available models")
	flags.StringVarP(&opts.SystemPrompt, "system", "s", "", "System prompt to set the assistant's behavior")
	flags.IntVarP(&opts.MaxHistory, "max-history", "", 20, "Maximum conversation history to keep (0 = unlimited)")
	flags.StringVarP(&opts.PresetPrompt, "preset", "p", "", "Use a preset system prompt (creative, concise, code)")
}

func newOllamaCommand() *cobra.Command {
	opts := &ChatOptions{}

	cmd := &cobra.Command{
		Use:   "ollama [prompt]",
		Short: "Use Ollama provider",
		Long:  `Use Ollama to run large language models locally`,
		Example: `  ai-cli ollama "What is the capital of Palestine?"
  ai-cli ollama -i  # Start interactive mode
  ai-cli ollama --model mistral "Write a story"
  ai-cli ollama -p creative "Tell me a short story"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := NewProvider(provider.Ollama, opts.ProviderURL)
			if err != nil {
				return err
			}
			if opts.ListModels {
				return displayModels(p)
			}
			resolveSystemPrompt(opts)

			return runChat(p, opts, args)
		},
	}

	addCommonFlags(cmd, opts)
	cmd.Flags().StringVarP(&opts.Model, "model", "m", "deepseek-r1:1.5b", "Model to use")
	cmd.Flags().StringVarP(&opts.ProviderURL, "url", "u", provider.DefaultURLs[provider.Ollama], "Provider API URL (optional)")

	return cmd
}

func newLocalAICommand() *cobra.Command {
	opts := &ChatOptions{}

	cmd := &cobra.Command{
		Use:   "localai [prompt]",
		Short: "Use LocalAI provider",
		Long:  `Use LocalAI self-hosted model server`,
		Example: `  ai-cli localai "What is the capital of Palestine?"
  ai-cli localai -i  # Start interactive mode
  ai-cli localai --model gpt-3.5-turbo "Write a story"
  ai-cli localai -p code "Explain binary search"`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := NewProvider(provider.LocalAI, opts.ProviderURL)
			if err != nil {
				return err
			}
			if opts.ListModels {
				return displayModels(p)
			}
			resolveSystemPrompt(opts)

			return runChat(p, opts, args)
		},
	}

	addCommonFlags(cmd, opts)
	cmd.Flags().StringVarP(&opts.Model, "model", "m", "gpt-3.5-turbo", "Model to use")
	cmd.Flags().StringVarP(&opts.ProviderURL, "url", "u", provider.DefaultURLs[provider.LocalAI], "Provider API URL (optional)")

	return cmd
}

func displayModels(p provider.Provider) error {
	models, err := p.ListModels()
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
