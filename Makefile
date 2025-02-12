.PHONY: build run test clean install format

build:
	go build -o bin/ollama-cli cmd/ollama-cli/main.go

run: build
	./bin/ollama-cli

test:
	go test -v ./test/...

clean:
	rm -f bin/ollama-cli
	go clean

install: build
	mv bin/ollama-cli $(GOPATH)/bin/ollama-cli

format:
	go fmt ./...
	go vet ./...
