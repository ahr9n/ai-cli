package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ahr9n/ai-cli/pkg/provider"
	"github.com/spf13/cobra"
)

type defaultConfig struct {
	Provider    string `json:"provider"`
	ProviderURL string `json:"provider_url,omitempty"`
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
	var providerURL string

	cmd := &cobra.Command{
		Use:   "set [provider]",
		Short: "Set default provider",
		Example: `  ai-cli default set ollama
  ai-cli default set localai --url http://custom:8080`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			providerType := provider.ProviderType(args[0])

			if _, err := NewProvider(providerType, ""); err != nil {
				return fmt.Errorf("invalid provider: %s", args[0])
			}

			config := defaultConfig{
				Provider:    args[0],
				ProviderURL: providerURL,
			}

			if err := saveDefaultConfig(config); err != nil {
				return fmt.Errorf("failed to save default config: %w", err)
			}

			fmt.Printf("Default provider set to: %s\n", args[0])
			if providerURL != "" {
				fmt.Printf("Provider URL: %s\n", providerURL)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&providerURL, "url", "u", "", "Provider API URL (optional)")

	return cmd
}

func newDefaultShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current default provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDefaultConfig()
			if err != nil {
				return fmt.Errorf("failed to load default config: %w", err)
			}

			if config.Provider == "" {
				fmt.Println("No default provider set")
				return nil
			}

			fmt.Printf("Default provider: %s\n", config.Provider)
			if config.ProviderURL != "" {
				fmt.Printf("Provider URL: %s\n", config.ProviderURL)
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
			if err := removeDefaultConfig(); err != nil {
				return fmt.Errorf("failed to clear default config: %w", err)
			}
			fmt.Println("Default provider cleared")
			return nil
		},
	}
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "ai-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "default.json"), nil
}

func saveDefaultConfig(config defaultConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func loadDefaultConfig() (defaultConfig, error) {
	var config defaultConfig
	path, err := getConfigPath()
	if err != nil {
		return config, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func removeDefaultConfig() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); os.IsNotExist(err) {
		return nil
	} else {
		return err
	}
}
