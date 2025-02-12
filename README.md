# Ollama CLI

A simple command-line interface for interacting with [Ollama](https://ollama.ai/) language models.

## Features

- Interactive chat mode
- Single prompt mode
- Model management (list available models)
- Support for different models (deepseek-r1:1.5b, llama2, mistral, etc.)
- Configurable temperature for response generation
- Simple and clean interface

## Prerequisites

Before using this CLI, make sure you have:

1. Go 1.23.4 or later installed
2. [Ollama](https://ollama.ai/) installed and running
3. At least one model pulled (e.g., `ollama pull deepseek-r1:1.5b`)

## Installation

```bash
# Clone the repository
git clone https://github.com/ahr9n/ollama-cli
cd ollama-cli

# Install dependencies
go mod download

# Build the binary
make build
```

## Usage

```
Usage:
  ollama-cli [prompt] [flags]
  ollama-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  models      List available models
  version     Print version information

Flags:
  -h, --help                  help for ollama-cli
  -i, --interactive           Start interactive chat mode
  -m, --model string          Model to use (deepseek-r1:1.5b, llama2, mistral, etc.) (default "deepseek-r1:1.5b")
  -t, --temperature float32   Sampling temperature (0.0-2.0) (default 0.7)

Examples:
  -> Single Prompt Mode
  ./ollama-cli "What is the capital of France?" # Using the default model (deepseek-r1:1.5b)
  ./ollama-cli --model mistral "Explain quantum computing" # Using a different model
  ./ollama-cli --temperature 0.9 "Write a creative story" # Adjusting temperature

  -> Start interactive mode
  ./ollama-cli -i # Start interactive chat
  ./ollama-cli -i --model llama2 # Start interactive chat with specific model

  -> Model Management
  ./ollama-cli models # List all available models

Use "./ollama-cli [command] --help" for more information about a command.
```

## Development

```bash
make build    # Build the binary
make run      # Run the CLI
make test     # Run tests
make format   # Format code
```
