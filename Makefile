.PHONY: build run test clean install format

build:
	go build -o ollama-cli cmd/main.go

run: build
	./ollama-cli

test:
	go test -v ./test/...

clean:
	rm -f ollama-cli
	go clean

install: build
	mv ollama-cli $(GOPATH)/ollama-cli

format:
	go fmt ./...
	go vet ./...
