package prompts

func DefaultSystem() string {
	return "You are a helpful assistant. Provide clear, direct & short responses."
}

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)
