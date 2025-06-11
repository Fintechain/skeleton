# Fintechain Skeleton Framework Makefile
# Streamlined build, test, and development workflow

# Project configuration
PROJECT_NAME := skeleton
MODULE_NAME := github.com/fintechain/skeleton
GO_VERSION := 1.21

# Directories
BIN_DIR := bin
COVERAGE_DIR := coverage
MOCKS_DIR := test/unit/mocks

# Build configuration
LDFLAGS := -w -s
BUILD_FLAGS := -ldflags "$(LDFLAGS)"
TEST_FLAGS := -race -timeout=30s

# Tools versions
GOLANGCI_LINT_VERSION := v1.62.2
MOCKERY_VERSION := v2.50.0

# Colors for output
GREEN := \033[0;32m
BLUE := \033[0;34m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: help
help: ## Display available commands
	@echo "$(BLUE)Fintechain Skeleton Framework$(NC)"
	@echo ""
	@echo "$(YELLOW)Development Commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# Core Development Workflow
# =============================================================================

.PHONY: dev
dev: clean fmt vet lint test ## Complete development cycle
	@echo "$(GREEN)✓ Development cycle completed$(NC)"

.PHONY: ci
ci: fmt vet lint test coverage ## CI pipeline
	@echo "$(GREEN)✓ CI pipeline completed$(NC)"

# =============================================================================
# Build Targets
# =============================================================================

.PHONY: build
build: ## Build example binaries
	@echo "$(BLUE)Building binaries...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/server ./cmd/server
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/client ./cmd/client
	@echo "$(GREEN)✓ Build completed$(NC)"

.PHONY: examples
examples: ## Run framework examples
	@echo "$(BLUE)Running FX daemon example...$(NC)"
	@go run examples/fx_usage.go daemon &
	@sleep 2 && pkill -f "fx_usage.go daemon" || true
	@echo "$(BLUE)Running FX command example...$(NC)"
	@go run examples/fx_usage.go command
	@echo "$(GREEN)✓ Examples completed$(NC)"

# =============================================================================
# Testing
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test $(TEST_FLAGS) ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(NC)"
	@go test $(TEST_FLAGS) ./test/unit/...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test $(TEST_FLAGS) ./test/integration/...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

.PHONY: test-short
test-short: ## Run tests (skip long-running)
	@echo "$(BLUE)Running short tests...$(NC)"
	@go test -short $(TEST_FLAGS) ./...
	@echo "$(GREEN)✓ Short tests passed$(NC)"

.PHONY: coverage
coverage: ## Generate test coverage report
	@echo "$(BLUE)Generating coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Coverage: $(COVERAGE_DIR)/coverage.html$(NC)"

# =============================================================================
# Code Quality
# =============================================================================

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet completed$(NC)"

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Linting completed$(NC)"

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "$(BLUE)Running linter with auto-fix...$(NC)"
	@golangci-lint run --fix ./...
	@echo "$(GREEN)✓ Linting with auto-fix completed$(NC)"

# =============================================================================
# Dependencies & Tools
# =============================================================================

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "$(BLUE)Tidying modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(BLUE)Installing tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	@echo "$(GREEN)✓ Tools installed$(NC)"

.PHONY: check-tools
check-tools: ## Check required tools
	@echo "$(BLUE)Checking tools...$(NC)"
	@command -v golangci-lint >/dev/null || { echo "$(RED)golangci-lint missing. Run 'make install-tools'$(NC)"; exit 1; }
	@command -v mockery >/dev/null || { echo "$(RED)mockery missing. Run 'make install-tools'$(NC)"; exit 1; }
	@echo "$(GREEN)✓ All tools available$(NC)"

# =============================================================================
# Mock Generation
# =============================================================================

.PHONY: mocks
mocks: ## Generate mocks using mockery
	@echo "$(BLUE)Generating mocks...$(NC)"
	@mockery --config .mockery.yaml
	@echo "$(GREEN)✓ Mocks generated$(NC)"

.PHONY: mocks-clean
mocks-clean: ## Clean generated mocks
	@echo "$(BLUE)Cleaning mocks...$(NC)"
	@find $(MOCKS_DIR) -name "*_mock.go" -delete 2>/dev/null || true
	@echo "$(GREEN)✓ Mocks cleaned$(NC)"

.PHONY: mocks-regen
mocks-regen: mocks-clean mocks ## Regenerate all mocks
	@echo "$(GREEN)✓ Mocks regenerated$(NC)"

# =============================================================================
# Benchmarks
# =============================================================================

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completed$(NC)"

# =============================================================================
# Cleanup
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BIN_DIR) $(COVERAGE_DIR)
	@go clean -cache -testcache
	@echo "$(GREEN)✓ Cleaned$(NC)"

# =============================================================================
# Information
# =============================================================================

.PHONY: version
version: ## Show version info
	@echo "$(BLUE)Project:$(NC) $(PROJECT_NAME)"
	@echo "$(BLUE)Module:$(NC) $(MODULE_NAME)"
	@echo "$(BLUE)Go Version:$(NC) $(GO_VERSION)"
	@echo "$(BLUE)Git Commit:$(NC) $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# =============================================================================
# Default
# =============================================================================

.DEFAULT_GOAL := help 