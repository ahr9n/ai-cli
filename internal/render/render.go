package render

import (
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

// Renderer applies Glamour markdown formatting when stdout is a terminal.
type Renderer struct {
	glamour    *glamour.TermRenderer
	isTerminal bool
}

// New creates a Renderer. Glamour is only initialized when stdout is a TTY.
func New() *Renderer {
	isTerm := term.IsTerminal(int(os.Stdout.Fd()))
	r := &Renderer{isTerminal: isTerm}

	if isTerm {
		tr, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(termWidth()),
		)
		if err == nil {
			r.glamour = tr
		}
	}
	return r
}

// RenderBlock renders a complete markdown string.
// Falls back to raw text if glamour is unavailable or rendering fails.
func (r *Renderer) RenderBlock(markdown string) string {
	if r.glamour == nil {
		return markdown
	}
	rendered, err := r.glamour.Render(markdown)
	if err != nil {
		return markdown
	}
	return strings.TrimSpace(rendered)
}

// IsTerminal reports whether stdout is a terminal.
func (r *Renderer) IsTerminal() bool {
	return r.isTerminal
}

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
