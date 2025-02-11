.PHONY: build run test

build:
	go build -o ollama-cli cmd/ollama-cli/main.go

run:
	./ollama-cli

test:
	go test ./test/...

format:
	go fmt ./...