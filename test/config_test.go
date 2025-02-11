package test

import (
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := config.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://localhost:11434", cfg.OllamaURL)
	assert.Equal(t, "deepseek-r1:1.5b", cfg.DefaultModel)
}
