all: tidy format build test

build:
	go build ./...

format:
	go fmt ./...

replace:
	go mod edit -replace="github.com/eleanorhealth/go-athenahealth=github.com/copilotiq/go-athenahealth@@$(commit)"

test:
	go test ./...

tidy:
	go mod tidy
