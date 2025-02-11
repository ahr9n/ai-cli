.PHONY: build run test

build:
	go build -o bin/ollama-cli cmd/ollama-cli/main.go

run:
	./bin/ollama-cli

test:
	go test ./test/...

format:
	go fmt ./...