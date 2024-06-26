all: tidy format build test

build:
	go build ./...

format:
	go fmt ./...

test:
	go test ./... -race

tidy:
	go mod tidy
