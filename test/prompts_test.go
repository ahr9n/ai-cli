package test

import (
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/prompts"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSystem(t *testing.T) {
	prompt := prompts.DefaultSystem()
	assert.NotEmpty(t, prompt)
}
