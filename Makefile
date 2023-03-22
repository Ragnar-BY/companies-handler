.PHONY build lint

build:
	go build ./cmd/server

lint:
	golangci-lint run