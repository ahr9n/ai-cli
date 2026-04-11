package render

import (
	"strings"
	"testing"

	"github.com/charmbracelet/glamour"
)

func TestRenderBlockFallback(t *testing.T) {
	// When glamour is nil (non-terminal), raw text is returned unchanged.
	r := &Renderer{glamour: nil, isTerminal: false}
	input := "# Hello\n\nSome **bold** text\n\n```go\nfmt.Println(\"hi\")\n```"
	got := r.RenderBlock(input)
	if got != input {
		t.Errorf("RenderBlock without glamour should return raw text\ngot:  %q\nwant: %q", got, input)
	}
}

func TestRenderBlockWithGlamour(t *testing.T) {
	// Construct a renderer with glamour enabled (regardless of actual terminal).
	r := newTestRenderer()
	if r.glamour == nil {
		t.Skip("glamour init failed in test environment")
	}

	input := "# Heading\n\nHello world\n\n```go\nfmt.Println(\"hi\")\n```"
	got := r.RenderBlock(input)

	// Glamour-rendered output differs from raw input.
	if got == input {
		t.Error("expected glamour to transform the markdown")
	}
	if !strings.Contains(got, "Heading") {
		t.Errorf("rendered output missing 'Heading': %q", got)
	}
	if !strings.Contains(got, "Hello world") {
		t.Errorf("rendered output missing 'Hello world': %q", got)
	}
}

func TestIsTerminal(t *testing.T) {
	r := &Renderer{isTerminal: false}
	if r.IsTerminal() {
		t.Error("expected IsTerminal() = false")
	}
	r.isTerminal = true
	if !r.IsTerminal() {
		t.Error("expected IsTerminal() = true")
	}
}

func TestTermWidth(t *testing.T) {
	w := termWidth()
	if w <= 0 {
		t.Errorf("termWidth() = %d, want > 0", w)
	}
}

// newTestRenderer creates a Renderer with glamour enabled, bypassing terminal detection.
func newTestRenderer() *Renderer {
	r := &Renderer{isTerminal: true}

	tr, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err == nil {
		r.glamour = tr
	}
	return r
}
