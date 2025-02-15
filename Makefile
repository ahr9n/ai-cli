.PHONY: build run test clean install format

build:
	go build -o bin/ai-cli cmd/main.go

run: build
	./bin/ai-cli

test:
	go test -v ./test/...

clean:
	rm -rf bin/
	go clean

install: build
	mv bin/ai-cli $(GOPATH)/bin/ai-cli

format:
	go fmt ./...
	go vet ./...
