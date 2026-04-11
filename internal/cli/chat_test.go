package cli

import (
	"os"
	"testing"
)

func TestReadStdinFromPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}

	input := "hello from pipe"
	if _, err := w.WriteString(input); err != nil {
		t.Fatal(err)
	}
	w.Close()

	got := readStdin(r)
	if got != input {
		t.Errorf("readStdin() = %q, want %q", got, input)
	}
}

func TestReadStdinTrimsWhitespace(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}

	if _, err := w.WriteString("  hello  \n\n"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	got := readStdin(r)
	if got != "hello" {
		t.Errorf("readStdin() = %q, want %q", got, "hello")
	}
}

func TestReadStdinEmptyPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}
	w.Close()

	got := readStdin(r)
	if got != "" {
		t.Errorf("readStdin() = %q, want empty", got)
	}
}

func TestReadStdinTerminal(t *testing.T) {
	// os.Stdin when attached to a TTY has ModeCharDevice set.
	// Open /dev/null as a regular file — it won't have ModeCharDevice,
	// so we test the terminal path by opening an actual device if available.
	f, err := os.Open("/dev/tty")
	if err != nil {
		t.Skip("no /dev/tty available")
	}
	defer f.Close()

	got := readStdin(f)
	if got != "" {
		t.Errorf("readStdin(tty) = %q, want empty", got)
	}
}
