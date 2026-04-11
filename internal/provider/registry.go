package provider

import "fmt"

// Registration holds metadata and a factory for a provider type.
type Registration struct {
	Factory      func(baseURL string) Provider
	DefaultURL   string
	DefaultModel string
	Name         string
	Description  string
}

var (
	registry          = make(map[ProviderType]Registration)
	registrationOrder []ProviderType
)

// Register adds a provider to the global registry.
func Register(pt ProviderType, r Registration) {
	if _, exists := registry[pt]; !exists {
		registrationOrder = append(registrationOrder, pt)
	}
	registry[pt] = r
}

// New creates a provider instance by type, using the default URL if baseURL is empty.
func New(pt ProviderType, baseURL string) (Provider, error) {
	r, ok := registry[pt]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", pt)
	}
	if baseURL == "" {
		baseURL = r.DefaultURL
	}
	return r.Factory(baseURL), nil
}

// DefaultURL returns the default base URL for a provider type.
func DefaultURL(pt ProviderType) string {
	if r, ok := registry[pt]; ok {
		return r.DefaultURL
	}
	return ""
}

// DefaultModel returns the default model for a provider type.
func DefaultModel(pt ProviderType) string {
	if r, ok := registry[pt]; ok {
		return r.DefaultModel
	}
	return ""
}

// All returns info for all registered providers in registration order.
func All() []ProviderInfo {
	result := make([]ProviderInfo, 0, len(registrationOrder))
	for _, pt := range registrationOrder {
		r := registry[pt]
		result = append(result, ProviderInfo{
			Type:        pt,
			Name:        r.Name,
			Description: r.Description,
		})
	}
	return result
}
