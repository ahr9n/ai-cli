package test

import (
	"testing"

	"github.com/ahr9n/ai-cli/pkg/cli"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	cmd := cli.NewRootCommand()

	assert.Equal(t, "ai-cli", cmd.Use)
	assert.Equal(t, "A CLI tool to interact with various AI models", cmd.Short)
	assert.Contains(t, cmd.Long, "command-line interface")

	subCommands := cmd.Commands()
	subCommandNames := make([]string, len(subCommands))
	for i, subCmd := range subCommands {
		subCommandNames[i] = subCmd.Name()
	}

	expectedCommands := []string{
		"version",
		"providers",
		"ollama",
		"localai",
		"default",
	}

	for _, expected := range expectedCommands {
		assert.Contains(t, subCommandNames, expected, "Subcommand %s should exist", expected)
	}

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRootCommandHelp(t *testing.T) {
	cmd := cli.NewRootCommand()
	cmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRootCommandUnknownCommand(t *testing.T) {
	cmd := cli.NewRootCommand()
	cmd.SetArgs([]string{"unknown"})

	err := cmd.Execute()
	assert.Error(t, err)
}
