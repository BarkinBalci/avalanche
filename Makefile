.PHONY: fmt build build-api build-consumer lint

fmt:
	@echo "Formatting code..."
	@go fmt ./...

build: build-api build-consumer

build-api:
	@echo "Building API..."
	@go build -o bin/api ./cmd/api

build-consumer:
	@echo "Building consumer..."
	@go build -o bin/consumer ./cmd/consumer

test:
	@echo "Running tests..."
	@go test -v -race ./...

lint:
	@echo "Running linter..."
	@golangci-lint run