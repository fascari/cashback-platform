.PHONY: build test lint run clean mocks fmt deps

build:
	@echo "Building mint-consumer..."
	@mkdir -p ../../bin
	go build -o ../../bin/mint-consumer ./cmd/main.go

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Linting..."
	golangci-lint run

run:
	@echo "Running mint-consumer..."
	go run ./cmd/main.go

clean:
	@echo "Cleaning..."
	rm -f ../../bin/mint-consumer
	rm -f coverage.out
	rm -rf mocks/

mocks:
	@echo "Generating mocks..."
	mockery

fmt:
	@echo "Formatting code..."
	goimports -w .
	gofmt -s -w .

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

