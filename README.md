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

1. Go 1.24 or later installed
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

### Single Prompt Mode

```bash
# Using the default model (deepseek-r1:1.5b)
./ollama-cli "What is the capital of France?"

# Using a different model
./ollama-cli --model mistral "What is the capital of France?"

# Adjusting temperature
./ollama-cli --temperature 0.9 "Write a creative story"
```

### Interactive Mode

```bash
# Start interactive chat
./ollama-cli -i

# Start interactive chat with specific model
./ollama-cli -i --model llama2
```

### Model Management

```bash
# List all available models
./ollama-cli models
```

### Available Options

```
Commands:
  version     Print version information
  models      List available models
  help        Help about any command

Flags:
  -i, --interactive        Start interactive chat mode
  -m, --model string      Model to use (default "deepseek-r1:1.5b")
  -t, --temperature float Temperature for response generation (default 0.7)
  -h, --help             Help for ollama-cli
```

## Project Structure

```
ollama-cli/
├── Makefile
├── README.md
├── cmd/
│   └── ollama-cli/
│       └── main.go          # Entry point
├── pkg/
│   ├── api/
│   │   └── base.go         # Base HTTP client
│   ├── cli/
│   │   ├── cli.go          # CLI implementation
│   │   └── models.go       # Models command
│   ├── config/
│   │   └── config.go       # Configuration handling
│   ├── ollama/
│   │   ├── client.go       # Ollama client
│   │   └── models.go       # Models API
│   ├── prompts/
│   │   └── prompts.go      # System prompts
│   └── utils/
│       └── loader.go       # Loading animation
└── test/
    ├── utils/
    │   └── server.go       # Test utilities and mocks
    ├── benchmark_test.go   # Performance benchmarks
    ├── cli_test.go         # CLI tests
    ├── config_test.go      # Configuration tests
    ├── ollama_test.go      # Client tests
    └── prompts_test.go     # Prompts tests
```

## Development

### Available Make Commands

```bash
make build    # Build the binary
make run      # Run the CLI
make test     # Run tests
make format   # Format code
```
