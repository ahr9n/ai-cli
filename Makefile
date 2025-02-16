.PHONY: build run test clean install format

build:
	go build -o ai-cli cmd/main.go

run: build
	./ai-cli

test:
	go test -v ./test/...

clean:
	rm -rf bin/
	go clean

install: build
	mv ai-cli $(GOPATH)/ai-cli

format:
	go fmt ./...
	go vet ./...
