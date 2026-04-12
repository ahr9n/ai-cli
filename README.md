<div align="center">
  <h1>AI CLI</h1>
  <p>Talk to any LLM from your terminal. One tool, many providers.</p>
  <p>
    <a href="https://github.com/ahr9n/ai-cli/actions/workflows/go.yml"><img src="https://github.com/ahr9n/ai-cli/actions/workflows/go.yml/badge.svg" alt="Go Tests"></a>
    <a href="https://goreportcard.com/report/github.com/ahr9n/ai-cli"><img src="https://goreportcard.com/badge/github.com/ahr9n/ai-cli" alt="Go Report Card"></a>
    <a href="https://github.com/ahr9n/ai-cli/releases"><img src="https://img.shields.io/github/v/release/ahr9n/ai-cli?color=blue" alt="Release"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  </p>
</div>

[![demo](https://asciinema.org/a/uUOFVezOecXNJyGS.svg)](https://asciinema.org/a/uUOFVezOecXNJyGS)

---

**ai-cli** is a multi-provider AI chat client for the command line. It supports local models via [Ollama](https://ollama.ai/), cloud APIs like [OpenAI](https://platform.openai.com/) and [Groq](https://groq.com/), and any [OpenAI-compatible](https://platform.openai.com/docs/api-reference/chat) endpoint — all through a single binary with streaming, piping, session persistence, and shell completions.

## Highlights

- **5 providers out of the box** — Ollama, OpenAI, Groq, LocalAI, and any custom OpenAI-compatible API
- **Streaming responses** — tokens appear as they arrive, with a spinner during latency
- **Pipe-friendly** — `cat file.go | ai-cli openai "review this code"` works as expected
- **Conversation persistence** — sessions auto-save as JSONL; list, continue, export, or delete them
- **Shell completions** — dynamic model name completion via `--model <TAB>`
- **JSON output** — `--output json` for scripting and automation
- **Zero config to start** — `ai-cli ollama "hello"` works immediately if Ollama is running

## Install

### Homebrew

```bash
brew install ahr9n/tap/ai-cli
```

### Go

```bash
go install github.com/ahr9n/ai-cli/cmd/ai-cli@latest
```

### Binary

Download a prebuilt binary from [Releases](https://github.com/ahr9n/ai-cli/releases).

### Source

```bash
git clone https://github.com/ahr9n/ai-cli.git && cd ai-cli
make build && make install
```

Requires Go 1.24+.

## Providers

| Provider | Type | Default URL | Auth |
|----------|------|-------------|------|
| [Ollama](https://ollama.ai/) | `ollama` | `localhost:11434` | — |
| [OpenAI](https://platform.openai.com/) | `openai` | `api.openai.com/v1` | `OPENAI_API_KEY` |
| [Groq](https://groq.com/) | `groq` | `api.groq.com/openai/v1` | `GROQ_API_KEY` |
| [LocalAI](https://github.com/mudler/LocalAI) | `localai` | `localhost:8080/v1` | — |
| Any OpenAI-compatible | `custom` | `localhost:8080/v1` | `AI_CLI_API_KEY` |

> Use `--url` to point any provider at a custom endpoint (vLLM, LM Studio, OpenRouter, etc.)

## Quick Start

```bash
# Single prompt
ai-cli ollama "What is the capital of Palestine?"

# Cloud provider
export OPENAI_API_KEY=sk-...
ai-cli openai "Explain quantum computing in 3 sentences"

# Pipe input
cat error.log | ai-cli openai "What went wrong?"
echo "package main" | ai-cli groq "Complete this Go program"

# Interactive chat
ai-cli ollama -i
```

## Usage

```
ai-cli <provider> [flags] [prompt]
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--interactive` | `-i` | Start interactive chat mode |
| `--model` | `-m` | Model to use |
| `--system` | `-s` | Custom system prompt |
| `--preset` | `-p` | Preset persona (`creative`, `concise`, `code`) |
| `--temperature` | `-t` | Sampling temperature (0.0–2.0) |
| `--max-tokens` | | Max response tokens (0 = default) |
| `--top-p` | | Nucleus sampling (0 = default) |
| `--output` | `-o` | Output format (`text`, `json`) |
| `--url` | `-u` | Custom provider URL |
| `--max-history` | | Message history limit (default: 20) |
| `--continue` | | Resume a saved session by ID |
| `--session` | | Custom session name |
| `--no-save` | | Disable auto-save in interactive mode |
| `--list-models` | | List available models |

### Default Provider

```bash
ai-cli default set openai --model gpt-4o
ai-cli "Hello, world"    # uses OpenAI with gpt-4o
ai-cli default show
ai-cli default clear
```

### Conversations

Sessions auto-save in interactive mode. Manage them with:

```bash
ai-cli chat list                        # List saved sessions
ai-cli chat show <id>                   # Display a conversation
ai-cli chat continue <id>              # Resume where you left off
ai-cli chat export <id>                # Export as markdown
ai-cli chat export <id> -f json        # Export as JSON
ai-cli chat delete <id>                # Delete a session
```

Sessions are stored as JSONL in `~/.local/share/ai-cli/sessions/`.

### JSON Output

```bash
ai-cli ollama --output json "What is Go?" | jq .response
```

```json
{
  "provider": "Ollama",
  "model": "deepseek-r1:1.5b",
  "prompt": "What is Go?",
  "response": "Go is a statically typed...",
  "error": null
}
```

### Shell Completions

```bash
# Bash
source <(ai-cli completion bash)

# Zsh
ai-cli completion zsh > "${fpath[1]}/_ai-cli"

# Fish
ai-cli completion fish | source
```

Model names complete dynamically: `ai-cli ollama --model <TAB>`.

## Configuration

Config is stored in `~/.config/ai-cli/config.yaml`:

```yaml
active_provider: openai
providers:
  openai:
    base_url: https://api.openai.com/v1
    model: gpt-4o
  ollama:
    model: llama3.1
```

Settings merge with CLI flags (explicit flags always win). Environment variables with the `AI_CLI_` prefix override config values.

## Architecture

```
cmd/ai-cli/main.go            Entrypoint — signal handling, context setup
internal/
├── cli/                       Cobra commands, chat loop, session management
├── api/                       HTTP client with auth headers, friendly errors
├── config/                    YAML config via Viper, XDG paths, JSON migration
├── provider/                  Interface, self-registering registry, error types
│   ├── ollama/                Ollama native API (/api/chat)
│   ├── openaicompat/          Shared OpenAI-compatible SSE streaming client
│   ├── openai/                OpenAI (thin wrapper)
│   ├── groq/                  Groq (thin wrapper)
│   ├── localai/               LocalAI (thin wrapper)
│   └── custom/                Custom endpoint (thin wrapper)
├── session/                   JSONL conversation persistence
├── render/                    Glamour markdown rendering
├── prompt/                    System prompt constants
└── loader/                    TTY-aware terminal spinner (stderr)
```

Providers self-register via `init()`. Adding a new provider is four steps — see [CONTRIBUTING.md](CONTRIBUTING.md).

## Development

```bash
make build      # Build binary (version injected via ldflags)
make test       # Tests with race detector
make format     # go fmt + go vet
make clean      # Remove artifacts
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full guide.

## License

[MIT](LICENSE)
