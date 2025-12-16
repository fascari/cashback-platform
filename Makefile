.PHONY: all build clean proto deps

# Default target
all: deps build

# Download dependencies for all services
deps:
	@echo "Downloading dependencies..."
	cd services/cashback-service-api && go mod tidy
	cd services/mint-consumer && go mod tidy
	cd services/blockchain-adapter && go mod tidy

# Build all services
build: build-cashback-service build-mint-consumer build-blockchain-adapter

build-cashback-service:
	@echo "Building cashback-service-api..."
	cd services/cashback-service-api && go build -o ../../bin/cashback-service-api ./cmd/api/main.go

build-mint-consumer:
	@echo "Building mint-consumer..."
	cd services/mint-consumer && go build -o ../../bin/mint-consumer ./cmd/main.go

build-blockchain-adapter:
	@echo "Building blockchain-adapter..."
	cd services/blockchain-adapter && go build -o ../../bin/blockchain-adapter ./cmd/main.go

# Run services (use in separate terminals)
run-cashback-service:
	cd services/cashback-service-api && go run ./cmd/api/main.go

run-mint-consumer:
	cd services/mint-consumer && go run ./cmd/main.go

run-blockchain-adapter:
	cd services/blockchain-adapter && go run ./cmd/main.go

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go-grpc_out=. proto/token.proto

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf services/*/vendor

# Test all services
test:
	@echo "Running tests..."
	cd services/cashback-service-api && go test ./...
	cd services/mint-consumer && go test ./...
	cd services/blockchain-adapter && go test ./...

# Format code
fmt:
	@echo "Formatting code..."
	cd services/cashback-service-api && go fmt ./...
	cd services/mint-consumer && go fmt ./...
	cd services/blockchain-adapter && go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	cd services/cashback-service-api && golangci-lint run
	cd services/mint-consumer && golangci-lint run
	cd services/blockchain-adapter && golangci-lint run

# Docker targets
docker-build:
	docker build -t cashback-service-api:latest -f services/cashback-service-api/Dockerfile .
	docker build -t mint-consumer:latest -f services/mint-consumer/Dockerfile .
	docker build -t blockchain-adapter:latest -f services/blockchain-adapter/Dockerfile .

# Database setup (requires PostgreSQL running)
db-setup:
	@echo "Setting up databases..."
	psql -U postgres -c "CREATE DATABASE cashback_service_db;" || true
	psql -U postgres -c "CREATE DATABASE mint_consumer_db;" || true
	psql -U postgres -c "CREATE DATABASE blockchain_adapter_db;" || true

# Help
help:
	@echo "Available targets:"
	@echo "  all                  - Download deps and build all services"
	@echo "  deps                 - Download dependencies"
	@echo "  build                - Build all services"
	@echo "  run-cashback-service - Run cashback service API"
	@echo "  run-mint-consumer    - Run mint consumer"
	@echo "  run-blockchain-adapter - Run blockchain adapter"
	@echo "  proto                - Generate protobuf code"
	@echo "  clean                - Clean build artifacts"
	@echo "  test                 - Run tests"
	@echo "  fmt                  - Format code"
	@echo "  lint                 - Lint code"
	@echo "  docker-build         - Build Docker images"
	@echo "  db-setup             - Create databases"

