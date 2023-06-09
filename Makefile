.PHONY: build run lint

build:
	go build ./cmd/server

run:
	go run ./...

lint:
	golangci-lint run

docker-up:
	docker-compose up

integration:
	go test -tags=integration ./internal/integration -v -count=1 -p=1
