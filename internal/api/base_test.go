package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]string
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("invalid JSON payload: %v", err)
		}
		if payload["key"] != "value" {
			t.Errorf("payload[key] = %q, want %q", payload["key"], "value")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewBaseClient(srv.URL)
	resp, err := c.DoPost(context.Background(), "test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("DoPost failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestDoGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := NewBaseClient(srv.URL)
	resp, err := c.DoGet(context.Background(), "test")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestCheckStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
		wantCode   int
	}{
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			body:       "ok",
			wantErr:    false,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			body:       "not found",
			wantErr:    true,
			wantCode:   404,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			body:       "server error",
			wantErr:    true,
			wantCode:   500,
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       "unauthorized",
			wantErr:    true,
			wantCode:   401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := NewBaseClient(srv.URL)
			resp, err := c.DoGet(context.Background(), "test")
			if err != nil {
				t.Fatalf("DoGet failed: %v", err)
			}
			defer resp.Body.Close()

			err = c.CheckStatus(resp)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var httpErr *HTTPError
				if !errors.As(err, &httpErr) {
					t.Fatalf("error type = %T, want *HTTPError", err)
				}
				if httpErr.StatusCode != tt.wantCode {
					t.Errorf("HTTPError.StatusCode = %d, want %d", httpErr.StatusCode, tt.wantCode)
				}
				if string(httpErr.Body) != tt.body {
					t.Errorf("HTTPError.Body = %q, want %q", httpErr.Body, tt.body)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCheckStatusDrainsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error body content"))
	}))
	defer srv.Close()

	c := NewBaseClient(srv.URL)
	resp, err := c.DoGet(context.Background(), "test")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}
	defer resp.Body.Close()

	_ = c.CheckStatus(resp)

	// Body should be drained — reading again returns empty.
	remaining, _ := io.ReadAll(resp.Body)
	if len(remaining) != 0 {
		t.Errorf("body not drained: got %d bytes remaining", len(remaining))
	}
}

func TestDoPostContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	c := NewBaseClient(srv.URL)
	_, err := c.DoPost(ctx, "test", map[string]string{})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
