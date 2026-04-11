# Changelog

## v1.0.0

### Added
- `--max-tokens` flag to limit response length
- `--top-p` flag for nucleus sampling control
- Friendly network error messages (connection refused, DNS failure)
- Integration tests with mock SSE server (6 end-to-end tests)
- `CONTRIBUTING.md` with architecture guide and provider addition instructions
- `.golangci.yml` lint configuration

### Changed
- Provider interface frozen: `CompletionOptions` now includes `MaxTokens` and `TopP`
- `api.BaseClient` produces user-readable errors for common network failures

## v0.5.0

### Added
- GoReleaser config for cross-platform builds (linux/darwin/windows x amd64/arm64)
- GitHub Actions release workflow triggered on tag push
- Homebrew tap support (`brew install ahr9n/tap/ai-cli`)
- `go install` support with `debug.ReadBuildInfo` version fallback
- VHS demo tape for terminal recording
- CHANGELOG.md

### Changed
- `ai-cli version` now resolves version from ldflags, build info, or "dev" fallback
- README rewritten with install methods, provider table, badges

## v0.4.0

### Added
- Conversation persistence with JSONL session storage (`~/.local/share/ai-cli/sessions/`)
- `ai-cli chat list` — list saved conversations
- `ai-cli chat show <id>` — display a conversation
- `ai-cli chat continue <id>` — resume a saved conversation
- `ai-cli chat delete <id>` — delete a conversation
- `ai-cli chat export <id> [-f markdown|json]` — export conversations
- Auto-save in interactive mode (disable with `--no-save`)
- `--session <name>` flag for custom session names
- `--continue <id>` flag to resume sessions from provider commands
- Tab completion for session IDs

## v0.3.0

### Added
- Glamour markdown rendering package (`internal/render/`)
- Shell completion command (`ai-cli completion bash/zsh/fish/powershell`)
- Dynamic model name tab completion (`--model <TAB>`)
- `--output json` flag for structured JSON output
- Tab completion for `--preset` and `--output` flags

## v0.2.0

### Added
- OpenAI-compatible shared client (`internal/provider/openaicompat/`)
- OpenAI, Groq, and Custom provider registrations
- YAML config via Viper (`~/.config/ai-cli/config.yaml`)
- Legacy JSON config migration
- Environment variable overrides (`AI_CLI_` prefix)
- `--model` and `--url` flags on `default set`
- Config merge into provider commands (respects explicit flags)
- Stdin pipe support (`echo "prompt" | ai-cli ollama`)
- Auth header support via `api.WithHeader` functional option

### Changed
- LocalAI refactored to thin wrapper over openaicompat
- Config stored as YAML instead of JSON

## v0.1.0

### Added
- Initial release
- Ollama and LocalAI providers with self-registering registry
- Interactive chat mode with conversation history
- Streaming responses (SSE parsing)
- System prompts and preset personas (creative, concise, code)
- Default provider management (`default set/show/clear`)
- Terminal spinner on stderr with TTY detection
- Signal handling (Ctrl+C / SIGTERM)
- `--list-models` flag
- `--max-history` flag for conversation history limits
