# AI CLI

A command-line interface for interacting with various AI providers and models.

## Features

- Multiple AI provider support:
  - Ollama (local LLM runtime)
  - LocalAI (self-hosted AI model server)
- Interactive chat mode
- Single prompt mode
- Model management
- Provider-specific configurations
- Stream responses in real-time

## Prerequisites

Before using this CLI, make sure you have:

1. Go 1.24 or later installed
2. Install & run at least one AI provider:
  - [Ollama](https://ollama.ai/) or
  - [LocalAI](https://github.com/go-skynet/LocalAI)

## Installation

```bash
# Clone the repository
git clone https://github.com/ahr9n/ai-cli
cd ai-cli

# Install dependencies
go mod download

# Build and install
make install
```

## Usage

```bash
# List available providers
ai-cli providers

# Use Ollama
ai-cli ollama "What is the capital of Palestine?"
ai-cli ollama -i  # Interactive mode [with the default model]
ai-cli ollama --model mistral "Explain quantum computing"

# Use LocalAI
ai-cli localai "What is the capital of Palestine?"
ai-cli localai -i --model gpt-3.5-turbo # Interactive mode [with a certain model]

# Set default provider and base URL
ai-cli default set ollama
ai-cli default set localai --url http://custom:8080

# Show current default
ai-cli default show

# Clear default
ai-cli default clear
```

### Available Commands
```
Commands:
  version     Print version information
  providers   List available providers
  ollama      Use Ollama provider
  localai     Use LocalAI provider
  default     Manage default provider settings
  help        Help about any command

Common Flags:
  -i, --interactive       Start interactive chat mode
  -m, --model string      Model to use (provider-specific)
  -t, --temperature float Temperature for response generation (default 0.7)
  -u, --url string        Provider API URL (optional)
```

## Development

```bash
make build    # Build the binary
make run      # Run the CLI
make test     # Run tests
make format   # Format code
```
