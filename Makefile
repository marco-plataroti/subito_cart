# Project variables
BINARY_NAME=subito-cart
PKG=./...
BUILD_DIR=bin/server

# Go commands
GO=go
GOTEST=$(GO) test -v
GOBUILD=$(GO) build
GORUN=$(GO) run
GOTIDY=$(GO) mod tidy
GOCLEAN=$(GO) clean
GOLINT=golangci-lint run

# Default target: build
all: build

## Build the binary
build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

## Run the application
run:
	$(GORUN) ./cmd/main.go

## Run tests
test:
	$(GOTEST) $(PKG)

## Format code
fmt:
	$(GO) fmt $(PKG)
	gofmt -s -w .
	goimports -w .
## Lint code
lint:
	$(GOLINT)

## Clean up compiled binaries
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

## Update dependencies
tidy:
	$(GOTIDY)

## Show available make commands
help:
	@echo "Available commands:"
	@echo "  make build     - Build the project"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Run linter"
	@echo "  make clean     - Clean up built files"
	@echo "  make tidy      - Update dependencies"

.PHONY: all init build run test fmt lint clean tidy help
