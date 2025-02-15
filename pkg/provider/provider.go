package provider

type Message struct {
	Role    string
	Content string
}

type CompletionOptions struct {
	Model       string
	Temperature float32
}

type ModelInfo struct {
	Name        string
	Size        int64
	Modified    string
	Family      string
	Description string
}

type Provider interface {
	CreateCompletion(messages []Message, opts *CompletionOptions) (string, error)
	StreamCompletion(messages []Message, opts *CompletionOptions, onResponse func(string)) error

	ListModels() ([]ModelInfo, error)
	GetDefaultModel() string

	Name() string
	Description() string
}

type ProviderType string

const (
	Ollama  ProviderType = "ollama"
	LocalAI ProviderType = "localai"
)

var DefaultURLs = map[ProviderType]string{
	Ollama:  "http://localhost:11434",
	LocalAI: "http://localhost:8080",
}
