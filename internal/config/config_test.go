package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidYAML(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `active_provider: openai
providers:
  openai:
    base_url: https://api.openai.com/v1
    model: gpt-4o
`
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(dir))
	// Dir() appends "ai-cli" — rename our temp dir to match.
	aiDir := filepath.Join(filepath.Dir(dir), "ai-cli")
	if err := os.Rename(dir, aiDir); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.ActiveProvider != "openai" {
		t.Errorf("ActiveProvider = %q, want %q", cfg.ActiveProvider, "openai")
	}
	pc, ok := cfg.Providers["openai"]
	if !ok {
		t.Fatal("expected openai in Providers map")
	}
	if pc.Model != "gpt-4o" {
		t.Errorf("Model = %q, want %q", pc.Model, "gpt-4o")
	}
	if pc.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("BaseURL = %q, want %q", pc.BaseURL, "https://api.openai.com/v1")
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	aiDir := filepath.Join(dir, "ai-cli")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error on missing file: %v", err)
	}
	if cfg.ActiveProvider != "" {
		t.Errorf("ActiveProvider = %q, want empty", cfg.ActiveProvider)
	}
}

func TestLoadMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	aiDir := filepath.Join(dir, "ai-cli")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(aiDir, "config.yaml"), []byte(":::bad yaml"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	aiDir := filepath.Join(dir, "ai-cli")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg := &Config{
		ActiveProvider: "groq",
		Providers: map[string]ProviderConfig{
			"groq": {Model: "llama-3.3-70b-versatile"},
		},
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.ActiveProvider != "groq" {
		t.Errorf("ActiveProvider = %q, want %q", loaded.ActiveProvider, "groq")
	}
	if loaded.Providers["groq"].Model != "llama-3.3-70b-versatile" {
		t.Errorf("Model = %q, want %q", loaded.Providers["groq"].Model, "llama-3.3-70b-versatile")
	}
}

func TestEnvVarOverride(t *testing.T) {
	dir := t.TempDir()
	aiDir := filepath.Join(dir, "ai-cli")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		t.Fatal(err)
	}
	yamlContent := `active_provider: ollama
`
	if err := os.WriteFile(filepath.Join(aiDir, "config.yaml"), []byte(yamlContent), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("AI_CLI_ACTIVE_PROVIDER", "openai")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.ActiveProvider != "openai" {
		t.Errorf("ActiveProvider = %q, want %q (env override)", cfg.ActiveProvider, "openai")
	}
}

func TestMigrateFromJSON(t *testing.T) {
	dir := t.TempDir()
	aiDir := filepath.Join(dir, "ai-cli")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		t.Fatal(err)
	}

	legacy := `{"provider":"localai","provider_url":"http://myhost:9000/v1"}`
	if err := os.WriteFile(filepath.Join(aiDir, "default.json"), []byte(legacy), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.ActiveProvider != "localai" {
		t.Errorf("ActiveProvider = %q, want %q", cfg.ActiveProvider, "localai")
	}
	if pc := cfg.Providers["localai"]; pc.BaseURL != "http://myhost:9000/v1" {
		t.Errorf("BaseURL = %q, want %q", pc.BaseURL, "http://myhost:9000/v1")
	}

	// Legacy file should be removed after migration.
	if _, err := os.Stat(filepath.Join(aiDir, "default.json")); !os.IsNotExist(err) {
		t.Error("legacy default.json should have been removed after migration")
	}
}
