package provider

import "errors"

var (
	ErrModelNotFound = errors.New("model not found")
	ErrNoMessages    = errors.New("no messages provided")
)
