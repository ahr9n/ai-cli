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

- Go 1.24 or later
- At least one AI provider:
  - [Ollama](https://ollama.ai/) or
  - [LocalAI](https://github.com/go-skynet/LocalAI)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/ai-cli
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
ai-cli ollama "What is the capital of France?"
ai-cli ollama -i  # Interactive mode
ai-cli ollama --model mistral "Explain quantum computing"

# Use LocalAI
ai-cli localai "What is the capital of France?"
ai-cli localai -i --model gpt-3.5-turbo

# Set default provider
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
  ollama      Use Ollama provider
  localai     Use LocalAI provider
  default     Manage default provider settings
  providers   List available providers
  version     Print version information
  help        Help about any command

Common Flags:
  -i, --interactive        Start interactive chat mode
  -m, --model string      Model to use (provider-specific)
  -t, --temperature float Temperature for response generation (default 0.7)
      --url string        Provider API URL (optional)
```

## Development

```bash
make build    # Build the binary
make run      # Run the CLI
make test     # Run tests
make format   # Format code
```

## Project Structure

```
ai-cli/
├── cmd/
│   └── main.go          # Entry point
├── pkg/
│   ├── api/            
│   │   └── base.go      # Base HTTP client
│   ├── cli/             # CLI commands
│   │   ├── chat.go
│   │   ├── provider.go
│   │   ├── root.go
│   │   └── version.go
│   ├── provider/        # Provider implementations
│   │   ├── provider.go
│   │   ├── localai/
│   │   └── ollama/
│   ├── prompts/         # System prompts
│   └── utils/           # Utilities
└── test/                # Tests
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Run tests and format code
4. Push your changes
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
