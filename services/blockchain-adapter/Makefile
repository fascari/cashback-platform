.PHONY: build test lint run clean mocks fmt deps

build:
	@echo "Building blockchain-adapter..."
	@mkdir -p ../../bin
	go build -o ../../bin/blockchain-adapter ./cmd/main.go

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Linting..."
	golangci-lint run

run:
	@echo "Running blockchain-adapter..."
	go run ./cmd/main.go

clean:
	@echo "Cleaning..."
	rm -f ../../bin/blockchain-adapter
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

