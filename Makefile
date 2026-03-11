.PHONY: build test clean install fmt vet lint help

BINARY_NAME=jsonneat
GO=go

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GO) build -o $(BINARY_NAME) .

test: ## Run all tests
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

fmt: ## Format code with gofmt
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run

clean: ## Remove build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

install: ## Install the binary to GOPATH/bin
	$(GO) install .

run: build ## Build and run with example
	./$(BINARY_NAME) --help

all: fmt vet test build ## Run fmt, vet, test, and build
