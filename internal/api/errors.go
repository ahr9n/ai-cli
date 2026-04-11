package api

import "fmt"

// HTTPError represents a non-OK HTTP response.
type HTTPError struct {
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("request failed (status %d): %s", e.StatusCode, e.Body)
}
