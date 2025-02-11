package test

import (
	"testing"

	"github.com/ahr9n/ollama-cli/pkg/cli"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	cmd := cli.NewRootCommand()

	assert.Equal(t, "ollama-cli [prompt]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	flags := cmd.Flags()
	interactiveFlag, _ := flags.GetBool("interactive")
	assert.False(t, interactiveFlag)

	modelFlag, _ := flags.GetString("model")
	assert.Equal(t, "deepseek-r1:1.5b", modelFlag)

	tempFlag, _ := flags.GetFloat32("temperature")
	assert.Equal(t, float32(0.7), tempFlag)
}
