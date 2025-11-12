.DEFAULT_GOAL := build

lint:
	golangci-lint run ./...

build: lint
	go build .

run:
	go run .
