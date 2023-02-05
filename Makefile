.DEFAULT_GOAL := build

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

vet:
	go run vet ./...

build: lint
	go build .

run:
	go run .

.PHONY:bench