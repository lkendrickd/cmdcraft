# CmdCraft Makefile
# -----------------------------------------------------
# This Makefile contains targets for building, testing,
# and managing the CmdCraft CLI framework project.
#
# Usage:
#   make <target>
#
# Run 'make' or 'make help' to see available targets.
# -----------------------------------------------------

# Default target
.DEFAULT_GOAL := help

# Binary output configuration
BIN := build/cmdcraft
BINARY_NAME := cmdcraft

# OS and architecture for cross-compilation
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Linker flags for smaller binaries
LDFLAGS := -ldflags "-s -w"

# golangci-lint version (pinned for reproducibility)
GOLANGCI_LINT_VERSION := v1.64.8

# Go module name
MODULE := $(shell go list -m)

##@ General

.PHONY: help
help: ## Show this help message
	@echo "CmdCraft CLI Framework - Available targets:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
	@echo ""
	@echo "Cross-compilation (override with environment variables):"
	@echo "  GOOS=$(GOOS)  GOARCH=$(GOARCH)"

##@ Development

.PHONY: build
build: ## Build the example binary
	@mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BIN) cmd/main.go
	@echo "Built: $(BIN)"

.PHONY: build-all
build-all: ## Build for linux, darwin, and windows (amd64)
	@mkdir -p build
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-linux-amd64 cmd/main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-darwin-amd64 cmd/main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-darwin-arm64 cmd/main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o build/$(BINARY_NAME)-windows-amd64.exe cmd/main.go
	@echo "Built binaries in build/"
	@ls -lh build/

.PHONY: run
run: build ## Build and run the example application
	./$(BIN) --help

.PHONY: clean
clean: ## Remove build artifacts and generated files
	rm -rf build/
	rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

##@ Testing

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: test-verbose
test-verbose: ## Run unit tests with verbose output
	go test -v ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	go test -race ./...

.PHONY: coverage
coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@rm -f coverage.out

.PHONY: coverage-html
coverage-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: bench
bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

##@ Code Quality

.PHONY: fmt
fmt: ## Format Go source files
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (installs if not found)
	@which golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))
	@PATH="$(shell go env GOPATH)/bin:$$PATH" golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@which golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))
	@PATH="$(shell go env GOPATH)/bin:$$PATH" golangci-lint run --fix ./...

.PHONY: check
check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "All checks passed!"

##@ Dependencies

.PHONY: deps
deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

.PHONY: deps-update
deps-update: ## Update all dependencies to latest
	go get -u ./...
	go mod tidy

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	go mod verify

##@ Installation

.PHONY: install
install: ## Install the example binary to GOPATH/bin
	go install $(LDFLAGS) ./cmd/...
	@echo "Installed to $(shell go env GOPATH)/bin/"

.PHONY: uninstall
uninstall: ## Remove the installed binary
	rm -f $(shell go env GOPATH)/bin/main
	@echo "Uninstalled from $(shell go env GOPATH)/bin/"
