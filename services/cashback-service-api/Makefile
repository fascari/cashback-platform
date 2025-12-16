.PHONY: build test lint run clean mocks fmt deps

build:
	@echo "Building cashback-service-api..."
	@mkdir -p ../../bin
	go build -o ../../bin/cashback-service-api ./cmd/api/main.go

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Linting..."
	golangci-lint run

run:
	@echo "Running cashback-service-api..."
	go run ./cmd/api/main.go

clean:
	@echo "Cleaning..."
	rm -f ../../bin/cashback-service-api
	rm -f coverage.out
	rm -rf mocks/

mocks:
	@echo "Generating mocks..."
	mockery

fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
