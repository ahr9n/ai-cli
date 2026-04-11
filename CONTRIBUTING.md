# Contributing to AI CLI

## Architecture

```
cmd/ai-cli/main.go          Entrypoint: signal handling, context setup
internal/
├── api/                     Shared HTTP client (BaseClient, functional options)
├── cli/                     Cobra command wiring, chat loop, session management
├── config/                  YAML config via Viper (XDG_CONFIG_HOME)
├── loader/                  Terminal spinner (stderr, TTY-aware)
├── prompt/                  System prompt constants
├── provider/                Provider interface, self-registering registry
│   ├── ollama/              Ollama (native /api/chat)
│   ├── openaicompat/        Shared OpenAI-compatible SSE client
│   ├── openai/              OpenAI (thin wrapper over openaicompat)
│   ├── groq/                Groq (thin wrapper over openaicompat)
│   ├── localai/             LocalAI (thin wrapper over openaicompat)
│   └── custom/              Custom OpenAI-compatible endpoint
├── render/                  Glamour markdown rendering
└── session/                 JSONL conversation persistence (XDG_DATA_HOME)
```

### Provider Registry

Providers self-register via `init()` and blank imports in `internal/cli/provider.go`. The `provider.Provider` interface is frozen as of v1.0.0:

```go
type Provider interface {
    Name() string
    StreamCompletion(ctx context.Context, messages []Message, opts *CompletionOptions, onChunk func(string)) error
    ListModels(ctx context.Context) ([]Model, error)
}
```

### Adding a New Provider

1. Create `internal/provider/<name>/client.go`
2. Implement the `Provider` interface
3. Call `provider.Register(...)` in `init()`
4. Add a blank import in `internal/cli/provider.go`

If your provider uses the OpenAI-compatible API, wrap `openaicompat.Client` instead of reimplementing SSE parsing.

## Development

### Prerequisites

- Go 1.24+
- [golangci-lint](https://golangci-lint.run/) (optional, for linting)

### Build & Test

```bash
make build          # Build to bin/ai-cli
make test           # Run tests with race detector
make format         # go fmt + go vet
make lint           # Run golangci-lint (if installed)
```

### Running Tests

```bash
go test ./... -race -count=1         # All tests
go test ./internal/cli/ -run Test    # Integration tests only
go test ./... -coverprofile=c.out    # Coverage report
go tool cover -func=c.out           # Coverage summary
```

### Code Style

- Follow standard Go conventions (`go fmt`, `go vet`)
- Use `context.Context` for all I/O operations
- Propagate errors with `fmt.Errorf("...: %w", err)`
- Keep provider implementations thin — shared logic belongs in `openaicompat` or `api`

### Commit Messages

Use conventional-style prefixes:

```
feat: add support for Anthropic provider
fix: handle empty SSE chunks in streaming parser
docs: update provider table in README
test: add integration tests for JSON output mode
```

## Releasing

Releases are automated via GoReleaser on tag push:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow builds binaries, creates a GitHub release, and updates the Homebrew tap.
