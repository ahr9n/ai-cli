package provider

import (
	"context"
	"testing"
)

// These tests mutate package-level registry state and must not use t.Parallel().
func TestRegisterAndNew(t *testing.T) {
	// Save and restore global state.
	origRegistry := registry
	origOrder := registrationOrder
	t.Cleanup(func() {
		registry = origRegistry
		registrationOrder = origOrder
	})
	registry = make(map[ProviderType]Registration)
	registrationOrder = nil

	reg := Registration{
		Factory: func(baseURL string) Provider {
			return &stubProvider{baseURL: baseURL}
		},
		DefaultURL:   "http://localhost:9999",
		DefaultModel: "test-model",
		Name:         "TestProvider",
		Description:  "A test provider",
	}

	Register("test", reg)

	t.Run("New with empty baseURL uses default", func(t *testing.T) {
		p, err := New("test", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		sp := p.(*stubProvider)
		if sp.baseURL != "http://localhost:9999" {
			t.Errorf("got baseURL %q, want %q", sp.baseURL, "http://localhost:9999")
		}
	})

	t.Run("New with custom baseURL", func(t *testing.T) {
		p, err := New("test", "http://custom:1234")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		sp := p.(*stubProvider)
		if sp.baseURL != "http://custom:1234" {
			t.Errorf("got baseURL %q, want %q", sp.baseURL, "http://custom:1234")
		}
	})

	t.Run("New with unknown type returns error", func(t *testing.T) {
		_, err := New("nonexistent", "")
		if err == nil {
			t.Fatal("expected error for unknown provider type")
		}
	})
}

func TestDefaultURL(t *testing.T) {
	origRegistry := registry
	origOrder := registrationOrder
	t.Cleanup(func() {
		registry = origRegistry
		registrationOrder = origOrder
	})
	registry = make(map[ProviderType]Registration)
	registrationOrder = nil

	Register("test", Registration{
		Factory:    func(baseURL string) Provider { return &stubProvider{baseURL: baseURL} },
		DefaultURL: "http://localhost:8080",
	})

	if got := DefaultURL("test"); got != "http://localhost:8080" {
		t.Errorf("DefaultURL = %q, want %q", got, "http://localhost:8080")
	}
	if got := DefaultURL("missing"); got != "" {
		t.Errorf("DefaultURL(missing) = %q, want empty", got)
	}
}

func TestDefaultModel(t *testing.T) {
	origRegistry := registry
	origOrder := registrationOrder
	t.Cleanup(func() {
		registry = origRegistry
		registrationOrder = origOrder
	})
	registry = make(map[ProviderType]Registration)
	registrationOrder = nil

	Register("test", Registration{
		Factory:      func(baseURL string) Provider { return &stubProvider{baseURL: baseURL} },
		DefaultModel: "llama3",
	})

	if got := DefaultModel("test"); got != "llama3" {
		t.Errorf("DefaultModel = %q, want %q", got, "llama3")
	}
	if got := DefaultModel("missing"); got != "" {
		t.Errorf("DefaultModel(missing) = %q, want empty", got)
	}
}

func TestAll(t *testing.T) {
	origRegistry := registry
	origOrder := registrationOrder
	t.Cleanup(func() {
		registry = origRegistry
		registrationOrder = origOrder
	})
	registry = make(map[ProviderType]Registration)
	registrationOrder = nil

	Register("alpha", Registration{
		Factory:     func(baseURL string) Provider { return &stubProvider{baseURL: baseURL} },
		Name:        "Alpha",
		Description: "First",
	})
	Register("beta", Registration{
		Factory:     func(baseURL string) Provider { return &stubProvider{baseURL: baseURL} },
		Name:        "Beta",
		Description: "Second",
	})

	all := All()
	if len(all) != 2 {
		t.Fatalf("got %d providers, want 2", len(all))
	}

	// Verify registration order is preserved.
	if all[0].Type != "alpha" {
		t.Errorf("first provider type = %q, want %q", all[0].Type, "alpha")
	}
	if all[1].Type != "beta" {
		t.Errorf("second provider type = %q, want %q", all[1].Type, "beta")
	}
}

func TestRegisterIdempotent(t *testing.T) {
	origRegistry := registry
	origOrder := registrationOrder
	t.Cleanup(func() {
		registry = origRegistry
		registrationOrder = origOrder
	})
	registry = make(map[ProviderType]Registration)
	registrationOrder = nil

	Register("dup", Registration{
		Factory: func(baseURL string) Provider { return &stubProvider{baseURL: baseURL} },
		Name:    "V1",
	})
	Register("dup", Registration{
		Factory: func(baseURL string) Provider { return &stubProvider{baseURL: baseURL} },
		Name:    "V2",
	})

	all := All()
	if len(all) != 1 {
		t.Fatalf("got %d providers, want 1 (no duplicate entry)", len(all))
	}
	if all[0].Name != "V2" {
		t.Errorf("name = %q, want %q (should be overwritten)", all[0].Name, "V2")
	}
}

// stubProvider is a minimal Provider implementation for testing.
type stubProvider struct {
	baseURL string
}

func (s *stubProvider) CreateCompletion(_ context.Context, _ []Message, _ *CompletionOptions) (string, error) {
	return "", nil
}
func (s *stubProvider) StreamCompletion(_ context.Context, _ []Message, _ *CompletionOptions, _ func(string)) error {
	return nil
}
func (s *stubProvider) ListModels(_ context.Context) ([]ModelInfo, error) {
	return nil, nil
}
func (s *stubProvider) Name() string        { return "stub" }
func (s *stubProvider) Description() string { return "stub" }
