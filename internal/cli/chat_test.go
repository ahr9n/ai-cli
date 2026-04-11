package cli

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestEmitJSON(t *testing.T) {
	var buf bytes.Buffer
	err := emitJSON(&buf, "ollama", "deepseek-r1:1.5b", "What is Go?", "Go is a language", nil)
	if err != nil {
		t.Fatalf("emitJSON error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if result["provider"] != "ollama" {
		t.Errorf("provider = %v, want %q", result["provider"], "ollama")
	}
	if result["model"] != "deepseek-r1:1.5b" {
		t.Errorf("model = %v, want %q", result["model"], "deepseek-r1:1.5b")
	}
	if result["prompt"] != "What is Go?" {
		t.Errorf("prompt = %v, want %q", result["prompt"], "What is Go?")
	}
	if result["response"] != "Go is a language" {
		t.Errorf("response = %v, want %q", result["response"], "Go is a language")
	}
	if result["error"] != nil {
		t.Errorf("error = %v, want nil", result["error"])
	}
}

func TestEmitJSONWithError(t *testing.T) {
	var buf bytes.Buffer
	err := emitJSON(&buf, "ollama", "bad-model", "hello", "", errors.New("model not found"))
	if err != nil {
		t.Fatalf("emitJSON error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result["response"] != "" {
		t.Errorf("response = %v, want empty", result["response"])
	}
	errVal, ok := result["error"].(string)
	if !ok || errVal != "model not found" {
		t.Errorf("error = %v, want %q", result["error"], "model not found")
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
