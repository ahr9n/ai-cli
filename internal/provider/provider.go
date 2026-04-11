package provider

import "context"

// Role represents a message role in a conversation.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message represents a single message in a conversation.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// CompletionOptions configures a completion request.
type CompletionOptions struct {
	Model       string
	Temperature float32
}

// ModelInfo describes an available model.
type ModelInfo struct {
	Name        string
	Size        int64
	Modified    string
	Family      string
	Description string
}

// Provider defines the interface for AI providers.
type Provider interface {
	CreateCompletion(ctx context.Context, messages []Message, opts *CompletionOptions) (string, error)
	StreamCompletion(ctx context.Context, messages []Message, opts *CompletionOptions, onResponse func(string)) error
	ListModels(ctx context.Context) ([]ModelInfo, error)
	GetDefaultModel() string
	Name() string
	Description() string
}

// ProviderType identifies a provider implementation.
type ProviderType string

const (
	Ollama  ProviderType = "ollama"
	LocalAI ProviderType = "localai"
)

// ProviderInfo describes an available provider.
type ProviderInfo struct {
	Type        ProviderType
	Name        string
	Description string
}
