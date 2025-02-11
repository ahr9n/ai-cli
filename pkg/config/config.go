package config

type Config struct {
	OllamaURL    string `yaml:"ollama_url"`
	DefaultModel string `yaml:"default_model"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		OllamaURL:    "http://localhost:11434", // default Ollama URL
		DefaultModel: "deepseek-r1:1.5b",       // default model
	}

	return cfg, nil
}
