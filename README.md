# AI CLI
A command-line interface for interacting with various AI providers and models.

## Features
- Support for multiple AI providers:
  - Ollama (local LLM runtime)
  - LocalAI (self-hosted AI model server)
- Interactive chat mode
- System prompts for controlling AI behavior
- Command history management
- Pre-defined persona presets
- Streaming responses
- Model selection



## Installation
```bash
git clone https://github.com/ahr9n/ai-cli.git
cd ai-cli
go mod download
make build
make install
```

## Requirements
- Go 1.24 or later installed
- For Ollama provider: [Ollama](https://ollama.ai/) installed and running
- For LocalAI provider: [LocalAI](https://github.com/go-skynet/LocalAI) installed and running

## Usage

### Basic Examples
```bash
ai-cli ollama "What is the capital of Palestine?" # use Ollama with a single prompt 
ai-cli localai "What is the old capital of Egypt?" # use LocalAI with a single prompt 
ai-cli ollama -i # start interactive chat mode with Ollama
ai-cli ollama --list-models # list available models with Ollama
ai-cli ollama --model llama3 "What is AI CLI?" # use a specific model
```

### Advanced Features
```bash
ai-cli ollama -s "You are a math tutor" "Explain calculus" # use a specific system prompt
ai-cli ollama -p creative "Tell me a story about a robot" # use a preset prompt style
ai-cli default set ollama # set a default provider
ai-cli ollama -i --max-history 10 # limit conversation history (in interactive mode)
ai-cli ollama -u http://localhost:8080 "Hello" # use a custom API URL
```

### Interactive Mode
In interactive mode, you can:
- Type `exit` or `quit` to end the session
- Type `clear` to reset the conversation history
- Have a continuous conversation with context

## Available Commands
```
ai-cli
  ├── version        - Print version information
  ├── providers      - List available AI providers
  ├── ollama         - Use Ollama provider
  ├── localai        - Use LocalAI provider
  └── default        - Manage default provider settings
      ├── set        - Set default provider
      ├── show       - Show current default provider
      └── clear      - Clear default provider setting
```

## System Prompt Presets
- `creative` - For imaginative and engaging responses
- `concise` - For brief, direct answers
- `code` - For coding assistance with clean, well-documented examples

## Development
```bash
make build    # Build the binary
make run      # Run the CLI
make test     # Run tests
make format   # Format code
```
