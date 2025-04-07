.PHONY: build test lint format clean install

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=git-tidy

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/git-tidy

# Run tests
test:
	$(GOTEST) -v ./...

# Run linter
lint:
	golangci-lint run

# Format code
format:
	gofmt -w .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install the application
install:
	$(GOINSTALL) ./cmd/git-tidy
	ln -sf $(shell go env GOPATH)/bin/git-tidy $(shell go env GOPATH)/bin/git-tidy
	mkdir -p $(shell go env GOPATH)/bin
	ln -sf $(shell go env GOPATH)/bin/git-tidy $(shell go env GOPATH)/bin/git

# Install dependencies
deps:
	$(GOGET) -v -t -d ./...

# Run all checks
check: lint test format

# Default target
all: check build 
