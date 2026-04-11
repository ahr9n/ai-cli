// Package session provides persistent conversation storage using JSONL files.
// Each session is a file in XDG_DATA_HOME/ai-cli/sessions/<id>.jsonl.
// The first line is metadata; subsequent lines are messages.
package session

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ahr9n/ai-cli/internal/provider"
)

// Metadata describes a saved session.
type Metadata struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Provider     string    `json:"provider"`      // registry key: "ollama"
	ProviderName string    `json:"provider_name"` // display name: "Ollama"
	Model        string    `json:"model"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Messages     int       `json:"messages"`
}

// Session holds metadata and the full message history.
type Session struct {
	Meta     Metadata           `json:"meta"`
	Messages []provider.Message `json:"messages"`
}

// NewID returns an 8-character hex string derived from the current time.
func NewID() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(h[:])[:8]
}

// DataDir returns the session storage directory, respecting XDG_DATA_HOME.
func DataDir() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "ai-cli", "sessions")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "ai-cli", "sessions")
}

// Save writes a session as a JSONL file (metadata line + message lines).
func Save(sess *Session) error {
	dir := DataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating session dir: %w", err)
	}

	f, err := os.Create(filepath.Join(dir, sess.Meta.ID+".jsonl"))
	if err != nil {
		return fmt.Errorf("creating session file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(struct {
		Meta Metadata `json:"meta"`
	}{Meta: sess.Meta}); err != nil {
		return fmt.Errorf("writing metadata: %w", err)
	}

	for _, msg := range sess.Messages {
		if err := enc.Encode(msg); err != nil {
			return fmt.Errorf("writing message: %w", err)
		}
	}

	return nil
}

// Load reads a session from its JSONL file.
func Load(id string) (*Session, error) {
	f, err := os.Open(filepath.Join(DataDir(), id+".jsonl"))
	if err != nil {
		return nil, fmt.Errorf("opening session %s: %w", id, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// First line: metadata.
	if !scanner.Scan() {
		return nil, fmt.Errorf("session %s: empty file", id)
	}
	var metaLine struct {
		Meta Metadata `json:"meta"`
	}
	if err := json.Unmarshal(scanner.Bytes(), &metaLine); err != nil {
		return nil, fmt.Errorf("session %s: invalid metadata: %w", id, err)
	}

	sess := &Session{Meta: metaLine.Meta}

	// Remaining lines: messages.
	for scanner.Scan() {
		var msg provider.Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue // skip corrupt lines
		}
		sess.Messages = append(sess.Messages, msg)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("session %s: read error: %w", id, err)
	}

	return sess, nil
}

// List returns metadata for all saved sessions, sorted by UpdatedAt descending.
func List() ([]Metadata, error) {
	dir := DataDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading session dir: %w", err)
	}

	var sessions []Metadata
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".jsonl")
		f, err := os.Open(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			var metaLine struct {
				Meta Metadata `json:"meta"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &metaLine); err == nil {
				metaLine.Meta.ID = id // ensure ID matches filename
				sessions = append(sessions, metaLine.Meta)
			}
		}
		f.Close()
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	return sessions, nil
}

// Delete removes a session file.
func Delete(id string) error {
	path := filepath.Join(DataDir(), id+".jsonl")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleting session %s: %w", id, err)
	}
	return nil
}

// AppendMessage adds a message to an existing session file and updates metadata.
func AppendMessage(sess *Session, msg provider.Message) error {
	sess.Messages = append(sess.Messages, msg)
	sess.Meta.Messages = len(sess.Messages)
	sess.Meta.UpdatedAt = time.Now()
	return Save(sess)
}

// Truncate shortens s to maxLen characters, adding "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
