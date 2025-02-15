package cli

import (
	"fmt"
	"github.com/ahr9n/ollama-cli/pkg/provider"
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
}

func newOllamaCommand() *cobra.Command {
	opts := &ChatOptions{}

	cmd := &cobra.Command{
		Use:   "ollama [prompt]",
		Short: "Use Ollama provider",
		Long:  `Use Ollama to run large language models locally`,
		Example: `  ai-cli ollama "What is the capital of France?"
  ai-cli ollama -i  # Start interactive mode
  ai-cli ollama --model mistral "Write a story"`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := NewProvider(provider.Ollama, opts.ProviderURL)
			if err != nil {
				return err
			}
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
		Example: `  ai-cli localai "What is the capital of France?"
  ai-cli localai -i  # Start interactive mode
  ai-cli localai --model gpt-3.5-turbo "Write a story"`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := NewProvider(provider.LocalAI, opts.ProviderURL)
			if err != nil {
				return err
			}
			return runChat(p, opts, args)
		},
	}

	addCommonFlags(cmd, opts)
	cmd.Flags().StringVarP(&opts.Model, "model", "m", "gpt-3.5-turbo", "Model to use")
	cmd.Flags().StringVarP(&opts.ProviderURL, "url", "u", provider.DefaultURLs[provider.Ollama], "Provider API URL (optional)")

	return cmd
}
