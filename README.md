# AI CLI

A command-line interface for interacting with various AI providers and models.

## Features

- Multiple AI providers (Ollama, LocalAI) with a pluggable registry
- Interactive chat mode with conversation history
- Streaming responses for both single prompts and interactive mode
- System prompts and preset personas (creative, concise, code)
- Configurable default provider
- Graceful cancellation via Ctrl+C / SIGTERM

## Installation

```bash
git clone https://github.com/ahr9n/ai-cli.git
cd ai-cli
make build
make install
```

## Requirements

- Go 1.24 or later
- At least one provider running locally:
  - [Ollama](https://ollama.ai/) — local LLM runtime
  - [LocalAI](https://github.com/mudler/LocalAI) — self-hosted OpenAI-compatible server

## Usage

### Single prompt

```bash
ai-cli ollama "What is the capital of Palestine?"
ai-cli localai "What is the old capital of Egypt?"
ai-cli ollama --model llama3 "What is AI CLI?"
```

### Interactive mode

```bash
ai-cli ollama -i
ai-cli ollama -i --max-history 10
```

In interactive mode, type `exit` or `quit` to end the session, and `clear` to reset conversation history.

### System prompts

```bash
ai-cli ollama -s "You are a math tutor" "Explain calculus"
ai-cli ollama -p creative "Tell me a story about a robot"
ai-cli ollama -p code "Explain binary search"
```

Available presets: `creative`, `concise`, `code`.

### Default provider

```bash
ai-cli default set ollama
ai-cli default set localai --url http://custom:8080
ai-cli default show
ai-cli default clear
```

With a default provider set, prompts work directly:

```bash
ai-cli "Hello, world"
```

### Other commands

```bash
ai-cli version
ai-cli providers
ai-cli ollama --list-models
ai-cli ollama -u http://localhost:8080 "Hello"
```

## Project Structure

```
cmd/ai-cli/main.go          Entrypoint: signal handling, context setup
internal/
├── cli/                     Cobra command wiring, chat loop, config persistence
├── api/                     Shared HTTP client (BaseClient, CheckStatus, HTTPError)
├── provider/                Provider interface, registry, error types
│   ├── ollama/              Ollama implementation (self-registers via init)
│   └── localai/             LocalAI implementation (self-registers via init)
├── prompt/                  System prompt constants
└── loader/                  Terminal spinner (stderr, TTY-aware)
```

Providers register themselves at import time. Adding a new provider means creating a package under `internal/provider/<name>/` with an `init()` that calls `provider.Register(...)` and a blank import in `internal/cli/provider.go`.

## Development

```bash
make build     # Build to bin/ai-cli (injects version via ldflags)
make test      # Run all tests with race detector
make format    # go fmt + go vet
make clean     # Remove build artifacts
```

## License

See [LICENSE](LICENSE).
