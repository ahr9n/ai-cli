package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the persistent CLI configuration.
type Config struct {
	ActiveProvider string                    `mapstructure:"active_provider" yaml:"active_provider,omitempty"`
	Providers      map[string]ProviderConfig `mapstructure:"providers" yaml:"providers,omitempty"`
}

// ProviderConfig stores per-provider settings.
type ProviderConfig struct {
	BaseURL string `mapstructure:"base_url" yaml:"base_url,omitempty"`
	Model   string `mapstructure:"model" yaml:"model,omitempty"`
}

// Load reads config from ~/.config/ai-cli/config.yaml with env var overrides.
// Returns a zero-value Config (not an error) if the file doesn't exist.
func Load() (*Config, error) {
	dir := Dir()

	// Migrate from legacy default.json if config.yaml doesn't exist yet.
	yamlPath := filepath.Join(dir, "config.yaml")
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		migrateFromJSON(filepath.Join(dir, "default.json"), yamlPath)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)
	v.SetEnvPrefix("AI_CLI")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// Save writes config to ~/.config/ai-cli/config.yaml.
func Save(cfg *Config) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	return os.WriteFile(filepath.Join(dir, "config.yaml"), data, 0644)
}

// Dir returns the config directory path, respecting XDG_CONFIG_HOME.
func Dir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "ai-cli")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "ai-cli")
}

// migrateFromJSON reads the legacy default.json and writes config.yaml.
func migrateFromJSON(jsonPath, yamlPath string) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return
	}

	var legacy struct {
		Provider    string `json:"provider"`
		ProviderURL string `json:"provider_url"`
	}
	if err := json.Unmarshal(data, &legacy); err != nil || legacy.Provider == "" {
		return
	}

	cfg := &Config{ActiveProvider: legacy.Provider}
	if legacy.ProviderURL != "" {
		cfg.Providers = map[string]ProviderConfig{
			legacy.Provider: {BaseURL: legacy.ProviderURL},
		}
	}

	if err := Save(cfg); err == nil {
		os.Remove(jsonPath)
	}
}
