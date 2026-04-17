.PHONY: help build build-all build-linux build-darwin build-windows clean run install fmt lint test

# Variables
APP_NAME := fileserver
MAIN_PATH := ./cmd/fileserver/main.go
BUILD_DIR := ./build
BIN_DIR := ./bin
TMP_DIR := ./tmp
GO := go
GOFLAGS := -v
CGO_ENABLED := 0
GOARCH := amd64

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)$(APP_NAME) - Build & Development Targets$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make build              # Build for current OS"
	@echo "  make build-all          # Build for all platforms"
	@echo "  make run                # Run locally"
	@echo "  make clean              # Remove binaries and temp files"

## Build targets
build: ## Build binary for current OS
	@mkdir -p $(BIN_DIR) $(TMP_DIR)
	@echo "$(BLUE)Building $(APP_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)...$(NC)"
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) \
		-ldflags="-s -w" -trimpath \
		-o $(TMP_DIR)/$(APP_NAME)_local $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: $(TMP_DIR)/$(APP_NAME)_local$(NC)"

build-all: ## Build for all platforms (Linux, macOS, Windows)
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@$(MAKE) build-linux
	@$(MAKE) build-darwin
	@$(MAKE) build-windows
	@echo "$(GREEN)✓ All builds complete!$(NC)"
	@echo "  Binaries: $(BIN_DIR)/$(APP_NAME)_{win,mac,lnx}"

build-linux: ## Build for Linux
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) $(BUILD_DIR)/build.sh linux $(APP_NAME) $(MAIN_PATH)

build-darwin: ## Build for macOS
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) $(BUILD_DIR)/build.sh darwin $(APP_NAME) $(MAIN_PATH) --force-macos

build-windows: ## Build for Windows
	CGO_ENABLED=$(CGO_ENABLED) GOARCH=$(GOARCH) $(BUILD_DIR)/build.sh windows $(APP_NAME) $(MAIN_PATH)

## Development targets
run: build ## Build and run locally
	@echo "$(BLUE)Running $(APP_NAME)...$(NC)"
	@$(TMP_DIR)/$(APP_NAME)_local

install: build ## Build and install to GOBIN
	@echo "$(BLUE)Installing $(APP_NAME) to $(shell go env GOBIN)...$(NC)"
	@cp $(TMP_DIR)/$(APP_NAME)_local $(shell go env GOBIN)/$(APP_NAME)
	@echo "$(GREEN)✓ Installed to: $(shell go env GOBIN)/$(APP_NAME)$(NC)"

## Code quality
fmt: ## Format code with gofmt
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

lint: ## Run golangci-lint (if installed)
	@if command -v golangci-lint > /dev/null; then \
		echo "$(BLUE)Running linter...$(NC)"; \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install with: brew install golangci-lint$(NC)"; \
	fi

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	TESTING=1 $(GO) test -v ./...

## Maintenance
clean: ## Remove binaries and temporary files
	@echo "$(BLUE)Cleaning up...$(NC)"
	@rm -rf $(BIN_DIR) $(TMP_DIR)
	@echo "$(GREEN)✓ Cleaned$(NC)"

tidy: ## Tidy go.mod and go.sum
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies tidied$(NC)"

## Info
info: ## Show build info
	@echo "$(BLUE)Project Information:$(NC)"
	@echo "  App Name:      $(APP_NAME)"
	@echo "  Go Version:    $(shell $(GO) version)"
	@echo "  OS:            $(shell go env GOOS)"
	@echo "  Arch:          $(shell go env GOARCH)"
	@echo "  Bin Dir:       $(BIN_DIR)"
	@echo "  Tmp Dir:       $(TMP_DIR)"
