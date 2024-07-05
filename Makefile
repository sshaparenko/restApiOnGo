.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golangci-lint run ./...
.PHONY:lint

vet: lint
	go vet ./...
.PHONY:vet

run:
	go run ./cmd/main.go
integration-test:vet
	go test -cover ./...
build:integration-test
	docker compose up 