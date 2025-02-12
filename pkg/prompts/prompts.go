package prompts

func DefaultSystem() string {
	return "You are a helpful assistant. Provide clear, direct & short responses."
}

func CreativeSystem() string {
	return "You are a creative assistant. Generate imaginative and engaging responses."
}

func ConciseSystem() string {
	return "You are a precise assistant. Provide brief, direct answers without elaboration."
}

func CodeSystem() string {
	return "You are a coding assistant. Focus on providing clean, well-documented code examples."
}

const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)
