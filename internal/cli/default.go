package cli

import (
	"fmt"

	"github.com/ahr9n/ai-cli/internal/config"
	"github.com/ahr9n/ai-cli/internal/provider"
	"github.com/spf13/cobra"
)

func handleDefaultProvider(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	cfg, err := config.Load()
	if err != nil || cfg.ActiveProvider == "" {
		return cmd.Help()
	}

	providerType := provider.ProviderType(cfg.ActiveProvider)
	if provider.DefaultURL(providerType) == "" {
		return fmt.Errorf("unknown default provider type: %s", cfg.ActiveProvider)
	}

	defaultModel := provider.DefaultModel(providerType)
	if pc, ok := cfg.Providers[cfg.ActiveProvider]; ok && pc.Model != "" {
		defaultModel = pc.Model
	}

	providerCmd := newProviderCommand(providerType, defaultModel)
	providerCmd.SetArgs(args)
	return providerCmd.ExecuteContext(cmd.Context())
}

func newDefaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default",
		Short: "Manage default provider settings",
	}

	cmd.AddCommand(
		newDefaultSetCommand(),
		newDefaultShowCommand(),
		newDefaultClearCommand(),
	)

	return cmd
}

func newDefaultSetCommand() *cobra.Command {
	var (
		providerURL string
		model       string
	)

	cmd := &cobra.Command{
		Use:   "set [provider]",
		Short: "Set default provider",
		Example: `  ai-cli default set ollama
  ai-cli default set openai --model gpt-4o
  ai-cli default set localai --url http://custom:8080/v1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			providerType := provider.ProviderType(args[0])
			if _, err := provider.New(providerType, ""); err != nil {
				return fmt.Errorf("invalid provider: %s", args[0])
			}

			cfg, _ := config.Load()
			cfg.ActiveProvider = args[0]

			if providerURL != "" || model != "" {
				if cfg.Providers == nil {
					cfg.Providers = make(map[string]config.ProviderConfig)
				}
				pc := cfg.Providers[args[0]]
				if providerURL != "" {
					pc.BaseURL = providerURL
				}
				if model != "" {
					pc.Model = model
				}
				cfg.Providers[args[0]] = pc
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Default provider set to: %s\n", args[0])
			if providerURL != "" {
				fmt.Printf("Provider URL: %s\n", providerURL)
			}
			if model != "" {
				fmt.Printf("Model: %s\n", model)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&providerURL, "url", "u", "", "Provider API URL")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Default model")

	return cmd
}

func newDefaultShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current default provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if cfg.ActiveProvider == "" {
				fmt.Println("No default provider set")
				return nil
			}

			fmt.Printf("Default provider: %s\n", cfg.ActiveProvider)
			if pc, ok := cfg.Providers[cfg.ActiveProvider]; ok {
				if pc.BaseURL != "" {
					fmt.Printf("Base URL: %s\n", pc.BaseURL)
				}
				if pc.Model != "" {
					fmt.Printf("Model: %s\n", pc.Model)
				}
			}
			return nil
		},
	}
}

func newDefaultClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clear default provider setting",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.Load()
			cfg.ActiveProvider = ""
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Println("Default provider cleared")
			return nil
		},
	}
}
