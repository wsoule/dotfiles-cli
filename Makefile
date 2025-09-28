# Dotfiles Makefile

.PHONY: build clean test install dev release-test help

# Build variables
BINARY_NAME=dotfiles
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse HEAD)
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target
all: build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## install: Install the binary to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .

## dev: Build and run in development mode
dev: build
	./$(BINARY_NAME)

## release-test: Test the release process
release-test:
	@echo "Testing release process..."
	goreleaser release --snapshot --clean

## deps: Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

## lint: Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying go modules..."
	go mod tidy

## check: Run all checks (test, lint, fmt)
check: test lint fmt

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'