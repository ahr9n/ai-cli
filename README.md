# Ollama CLI

A simple command-line interface for interacting with [Ollama](https://ollama.ai/) language models.

## Features

- Interactive chat mode
- Single prompt mode
- Support for different models (deepseek-r1:1.5b, llama2, mistral, etc.)
- Configurable temperature for response generation
- Simple and clean interface

## Prerequisites

Before using this CLI, make sure you have:

1. Go 1.21 or later installed
2. [Ollama](https://ollama.ai/) installed and running
3. At least one model pulled (e.g., `ollama pull deepseek-r1:1.5b`)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/ollama-cli
cd ollama-cli

# Install dependencies
go mod download

# Build the binary
go build -o ollama-cli cmd/ollama-cli/main.go
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

### Available Options

```
Flags:
  -i, --interactive        Start interactive chat mode
  -m, --model string      Model to use (default "deepseek-r1:1.5b")
  -t, --temperature float Temperature for response generation (default 0.7)
  -h, --help             Help for ollama-cli
```

## Project Structure

```
ollama-cli/
├── cmd/
│   └── ollama-cli/
│       └── main.go          # Entry point
├── pkg/
│   ├── cli/
│   │   └── cli.go          # CLI implementation
│   ├── config/
│   │   └── config.go       # Configuration handling
│   ├── ollama/
│   │   └── client.go       # Ollama API client
│   └── prompts/
│       └── prompts.go      # Prompt management
├── go.mod
└── README.md
```
