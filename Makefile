.PHONY: build run test clean install format

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-X github.com/ahr9n/ai-cli/internal/cli.Version=$(VERSION)" -o bin/ai-cli ./cmd/ai-cli

run: build
	./bin/ai-cli

test:
	go test -race -v ./...

clean:
	rm -rf bin/
	go clean

install: build
	install -d $(HOME)/go/bin
	install bin/ai-cli $(HOME)/go/bin/ai-cli

format:
	go fmt ./...
	go vet ./...
