.PHONY: fmt build build-api build-consumer lint swagger

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

swagger:
	@echo "Generating swagger documentation..."
	@swag init -g cmd/api/main.go -o docs