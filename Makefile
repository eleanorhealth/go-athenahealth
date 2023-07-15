all: tidy format build lint test

build:
	go build ./...

format:
	go fmt ./...

lint:
	golangci-lint run --timeout=5m --tests=false

test:
	go test ./...

tidy:
	go mod tidy
