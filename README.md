# AI CLI

A command-line interface for interacting with various AI providers and models.

<!-- ![demo](demo.gif) -->

[![Go Tests](https://github.com/ahr9n/ai-cli/actions/workflows/go.yml/badge.svg)](https://github.com/ahr9n/ai-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ahr9n/ai-cli)](https://goreportcard.com/report/github.com/ahr9n/ai-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Installation

### Homebrew (macOS/Linux)

```bash
brew install ahr9n/tap/ai-cli
```

### Go install

```bash
go install github.com/ahr9n/ai-cli/cmd/ai-cli@latest
```

### Binary download

Download from [Releases](https://github.com/ahr9n/ai-cli/releases).

### Build from source

```bash
git clone https://github.com/ahr9n/ai-cli.git
cd ai-cli
make build
make install
```

Requires Go 1.24 or later.

## Supported Providers

| Provider | Type | Default URL | API Key Env Var |
|----------|------|-------------|-----------------|
| [Ollama](https://ollama.ai/) | `ollama` | `http://localhost:11434` | — |
| [OpenAI](https://platform.openai.com/) | `openai` | `https://api.openai.com/v1` | `OPENAI_API_KEY` |
| [Groq](https://groq.com/) | `groq` | `https://api.groq.com/openai/v1` | `GROQ_API_KEY` |
| [LocalAI](https://github.com/mudler/LocalAI) | `localai` | `http://localhost:8080/v1` | — |
| Custom | `custom` | `http://localhost:8080/v1` | `AI_CLI_API_KEY` |

Any OpenAI-compatible API works via the `custom` provider with `--url`.

## Usage

### Single prompt

```bash
ai-cli ollama "What is the capital of Palestine?"
ai-cli openai "Explain quantum computing"
ai-cli groq --model llama-3.3-70b-versatile "What is AI CLI?"
```

### Piped input

```bash
echo "Explain this code" | ai-cli ollama
cat error.log | ai-cli openai "What went wrong?"
```

### Interactive mode

```bash
ai-cli ollama -i
ai-cli openai -i --max-history 10
```

Type `exit`/`quit` to end, `clear` to reset history. Sessions auto-save by default.

### System prompts

```bash
ai-cli ollama -s "You are a math tutor" "Explain calculus"
ai-cli ollama -p creative "Tell me a story"
ai-cli ollama -p code "Explain binary search"
```

Presets: `creative`, `concise`, `code`.

### Default provider

```bash
ai-cli default set ollama
ai-cli default set openai --model gpt-4o
ai-cli default show
ai-cli default clear
```

With a default set, prompts work directly:

```bash
ai-cli "Hello, world"
```

### Conversation management

```bash
ai-cli chat list
ai-cli chat show <session-id>
ai-cli chat continue <session-id>
ai-cli chat delete <session-id>
ai-cli chat export <session-id>
ai-cli chat export <session-id> -f json
```

Sessions are stored as JSONL in `~/.local/share/ai-cli/sessions/`.

### JSON output

```bash
ai-cli ollama --output json "What is Go?"
```

### Other commands

```bash
ai-cli version
ai-cli providers
ai-cli ollama --list-models
ai-cli completion bash    # also: zsh, fish, powershell
```

## Shell Completions

### Bash

```bash
source <(ai-cli completion bash)
# Or persist:
ai-cli completion bash > /etc/bash_completion.d/ai-cli
```

### Zsh

```bash
ai-cli completion zsh > "${fpath[1]}/_ai-cli"
```

### Fish

```bash
ai-cli completion fish | source
ai-cli completion fish > ~/.config/fish/completions/ai-cli.fish
```

### PowerShell

```powershell
ai-cli completion powershell | Out-String | Invoke-Expression
```

## Project Structure

```
cmd/ai-cli/main.go          Entrypoint: signal handling, context setup
internal/
├── cli/                     Cobra command wiring, chat loop, session management
├── api/                     Shared HTTP client (BaseClient, CheckStatus, HTTPError)
├── config/                  YAML config via Viper (XDG_CONFIG_HOME)
├── provider/                Provider interface, registry, error types
│   ├── ollama/              Ollama (native /api/chat)
│   ├── openai/              OpenAI
│   ├── groq/                Groq
│   ├── localai/             LocalAI
│   ├── custom/              Custom OpenAI-compatible
│   └── openaicompat/        Shared OpenAI-compatible client (SSE streaming)
├── session/                 JSONL conversation persistence (XDG_DATA_HOME)
├── render/                  Glamour markdown rendering
├── prompt/                  System prompt constants
└── loader/                  Terminal spinner (stderr, TTY-aware)
```

Providers self-register via `init()`. Adding a new provider: create `internal/provider/<name>/` with an `init()` calling `provider.Register(...)`, then add a blank import in `internal/cli/provider.go`.

## Development

```bash
make build     # Build to bin/ai-cli (injects version via ldflags)
make test      # Run all tests with race detector
make format    # go fmt + go vet
make clean     # Remove build artifacts
```

## License

See [LICENSE](LICENSE).
