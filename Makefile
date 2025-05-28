# Skeleton Framework Makefile
# Comprehensive build, test, and development workflow automation

# Project configuration
PROJECT_NAME := skeleton
MODULE_NAME := github.com/fintechain/skeleton
GO_VERSION := 1.24.2

# Directories
BIN_DIR := bin
CMD_DIR := cmd
PKG_DIR := pkg
INTERNAL_DIR := internal
TEST_DIR := test
DOCS_DIR := docs
COVERAGE_DIR := coverage

# Build configuration
LDFLAGS := -w -s
BUILD_FLAGS := -ldflags "$(LDFLAGS)"
TEST_FLAGS := -race -timeout=30s
COVERAGE_FLAGS := -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic

# Tools
GOLANGCI_LINT_VERSION := v1.55.2
MOCKGEN_VERSION := v1.6.0

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

.PHONY: help
help: ## Display this help message
	@echo "$(CYAN)Skeleton Framework - Development Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# Build Targets
# =============================================================================

.PHONY: build
build: ## Build all binaries
	@echo "$(BLUE)Building all binaries...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/fx-example ./$(CMD_DIR)/fx-example
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/server ./$(CMD_DIR)/server
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/client ./$(CMD_DIR)/client
	@echo "$(GREEN)✓ Build completed$(NC)"

.PHONY: build-fx-example
build-fx-example: ## Build the FX integration example
	@echo "$(BLUE)Building FX example...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/fx-example ./$(CMD_DIR)/fx-example
	@echo "$(GREEN)✓ FX example built$(NC)"

.PHONY: build-server
build-server: ## Build the server binary
	@echo "$(BLUE)Building server...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/server ./$(CMD_DIR)/server
	@echo "$(GREEN)✓ Server built$(NC)"

.PHONY: build-client
build-client: ## Build the client binary
	@echo "$(BLUE)Building client...$(NC)"
	@mkdir -p $(BIN_DIR)
	@go build $(BUILD_FLAGS) -o $(BIN_DIR)/client ./$(CMD_DIR)/client
	@echo "$(GREEN)✓ Client built$(NC)"

.PHONY: install
install: ## Install binaries to GOPATH/bin
	@echo "$(BLUE)Installing binaries...$(NC)"
	@go install $(BUILD_FLAGS) ./$(CMD_DIR)/fx-example
	@go install $(BUILD_FLAGS) ./$(CMD_DIR)/server
	@go install $(BUILD_FLAGS) ./$(CMD_DIR)/client
	@echo "$(GREEN)✓ Binaries installed$(NC)"

# =============================================================================
# Test Targets
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running all tests...$(NC)"
	@go test $(TEST_FLAGS) ./...
	@echo "$(GREEN)✓ All tests passed$(NC)"

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(PKG_DIR)/... ./$(INTERNAL_DIR)/...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

.PHONY: test-integration
test-integration: ## Run integration tests only
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(TEST_DIR)/integration/...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

.PHONY: test-system
test-system: ## Run system integration tests
	@echo "$(BLUE)Running system integration tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(TEST_DIR)/integration/system/...
	@echo "$(GREEN)✓ System tests passed$(NC)"

.PHONY: test-component
test-component: ## Run component tests
	@echo "$(BLUE)Running component tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(INTERNAL_DIR)/domain/component/...
	@echo "$(GREEN)✓ Component tests passed$(NC)"

.PHONY: test-storage
test-storage: ## Run storage tests
	@echo "$(BLUE)Running storage tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(INTERNAL_DIR)/domain/storage/... ./$(INTERNAL_DIR)/infrastructure/storage/...
	@echo "$(GREEN)✓ Storage tests passed$(NC)"

.PHONY: test-plugin
test-plugin: ## Run plugin tests
	@echo "$(BLUE)Running plugin tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(INTERNAL_DIR)/domain/plugin/...
	@echo "$(GREEN)✓ Plugin tests passed$(NC)"

.PHONY: test-service
test-service: ## Run service tests
	@echo "$(BLUE)Running service tests...$(NC)"
	@go test $(TEST_FLAGS) ./$(INTERNAL_DIR)/domain/service/...
	@echo "$(GREEN)✓ Service tests passed$(NC)"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	@go test -v $(TEST_FLAGS) ./...

.PHONY: test-short
test-short: ## Run tests with short flag (skip long-running tests)
	@echo "$(BLUE)Running short tests...$(NC)"
	@go test -short $(TEST_FLAGS) ./...

# =============================================================================
# Coverage Targets
# =============================================================================

.PHONY: coverage
coverage: ## Generate test coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: coverage-unit
coverage-unit: ## Generate coverage for unit tests only
	@echo "$(BLUE)Generating unit test coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) ./$(PKG_DIR)/... ./$(INTERNAL_DIR)/...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage-unit.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Unit test coverage: $(COVERAGE_DIR)/coverage-unit.html$(NC)"

.PHONY: coverage-integration
coverage-integration: ## Generate coverage for integration tests
	@echo "$(BLUE)Generating integration test coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test $(COVERAGE_FLAGS) ./$(TEST_DIR)/integration/...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage-integration.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1
	@echo "$(GREEN)✓ Integration test coverage: $(COVERAGE_DIR)/coverage-integration.html$(NC)"

.PHONY: coverage-show
coverage-show: ## Show coverage in browser
	@echo "$(BLUE)Opening coverage report in browser...$(NC)"
	@open $(COVERAGE_DIR)/coverage.html || xdg-open $(COVERAGE_DIR)/coverage.html

# =============================================================================
# Code Quality Targets
# =============================================================================

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

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "$(BLUE)Tidying go modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

.PHONY: mod-verify
mod-verify: ## Verify go modules
	@echo "$(BLUE)Verifying go modules...$(NC)"
	@go mod verify
	@echo "$(GREEN)✓ Modules verified$(NC)"

.PHONY: mod-download
mod-download: ## Download go modules
	@echo "$(BLUE)Downloading go modules...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Modules downloaded$(NC)"

# =============================================================================
# Development Targets
# =============================================================================

.PHONY: run-fx-example
run-fx-example: build-fx-example ## Run the FX integration example
	@echo "$(BLUE)Running FX example...$(NC)"
	@./$(BIN_DIR)/fx-example

.PHONY: run-server
run-server: build-server ## Run the server
	@echo "$(BLUE)Running server...$(NC)"
	@./$(BIN_DIR)/server

.PHONY: run-client
run-client: build-client ## Run the client
	@echo "$(BLUE)Running client...$(NC)"
	@./$(BIN_DIR)/client

.PHONY: dev
dev: clean fmt vet lint test build ## Full development cycle (clean, format, vet, lint, test, build)
	@echo "$(GREEN)✓ Development cycle completed$(NC)"

.PHONY: ci
ci: mod-verify fmt vet lint test-short coverage ## CI pipeline (verify, format, vet, lint, test, coverage)
	@echo "$(GREEN)✓ CI pipeline completed$(NC)"

# =============================================================================
# Documentation Targets
# =============================================================================

.PHONY: docs
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@go doc -all ./... > $(DOCS_DIR)/api.txt
	@echo "$(GREEN)✓ Documentation generated$(NC)"

.PHONY: docs-serve
docs-serve: ## Serve documentation locally
	@echo "$(BLUE)Serving documentation on http://localhost:6060$(NC)"
	@godoc -http=:6060

# =============================================================================
# Mock Generation Targets
# =============================================================================

.PHONY: mocks
mocks: ## Generate mocks for testing
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go generate ./...
	@echo "$(GREEN)✓ Mocks generated$(NC)"

.PHONY: mocks-component
mocks-component: ## Generate component mocks
	@echo "$(BLUE)Generating component mocks...$(NC)"
	@mockgen -source=$(INTERNAL_DIR)/domain/component/component.go -destination=$(INTERNAL_DIR)/domain/component/mocks/component_mock.go
	@mockgen -source=$(INTERNAL_DIR)/domain/component/registry.go -destination=$(INTERNAL_DIR)/domain/component/mocks/registry_mock.go
	@echo "$(GREEN)✓ Component mocks generated$(NC)"

.PHONY: mocks-storage
mocks-storage: ## Generate storage mocks
	@echo "$(BLUE)Generating storage mocks...$(NC)"
	@mockgen -source=$(INTERNAL_DIR)/domain/storage/store.go -destination=$(INTERNAL_DIR)/domain/storage/mocks/store_mock.go
	@mockgen -source=$(INTERNAL_DIR)/domain/storage/multistore.go -destination=$(INTERNAL_DIR)/domain/storage/mocks/multistore_mock.go
	@echo "$(GREEN)✓ Storage mocks generated$(NC)"

# =============================================================================
# Benchmark Targets
# =============================================================================

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completed$(NC)"

.PHONY: bench-component
bench-component: ## Run component benchmarks
	@echo "$(BLUE)Running component benchmarks...$(NC)"
	@go test -bench=. -benchmem ./$(INTERNAL_DIR)/domain/component/...

.PHONY: bench-storage
bench-storage: ## Run storage benchmarks
	@echo "$(BLUE)Running storage benchmarks...$(NC)"
	@go test -bench=. -benchmem ./$(INTERNAL_DIR)/infrastructure/storage/...

# =============================================================================
# Tool Installation Targets
# =============================================================================

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
	@echo "$(GREEN)✓ Development tools installed$(NC)"

.PHONY: check-tools
check-tools: ## Check if required tools are installed
	@echo "$(BLUE)Checking required tools...$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 || { echo "$(RED)golangci-lint not found. Run 'make install-tools'$(NC)"; exit 1; }
	@command -v mockgen >/dev/null 2>&1 || { echo "$(RED)mockgen not found. Run 'make install-tools'$(NC)"; exit 1; }
	@echo "$(GREEN)✓ All required tools are installed$(NC)"

# =============================================================================
# Cleanup Targets
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BIN_DIR)
	@rm -rf $(COVERAGE_DIR)
	@go clean -cache
	@go clean -testcache
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

.PHONY: clean-mocks
clean-mocks: ## Clean generated mocks
	@echo "$(BLUE)Cleaning generated mocks...$(NC)"
	@find . -name "*_mock.go" -type f -delete
	@echo "$(GREEN)✓ Mocks cleaned$(NC)"

# =============================================================================
# Release Targets
# =============================================================================

.PHONY: version
version: ## Show version information
	@echo "$(CYAN)Project: $(PROJECT_NAME)$(NC)"
	@echo "$(CYAN)Module: $(MODULE_NAME)$(NC)"
	@echo "$(CYAN)Go Version: $(GO_VERSION)$(NC)"
	@echo "$(CYAN)Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')$(NC)"
	@echo "$(CYAN)Git Branch: $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'unknown')$(NC)"

.PHONY: tag
tag: ## Create a git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "$(RED)VERSION is required. Usage: make tag VERSION=v1.0.0$(NC)"; exit 1; fi
	@echo "$(BLUE)Creating tag $(VERSION)...$(NC)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "$(GREEN)✓ Tag $(VERSION) created and pushed$(NC)"

# =============================================================================
# Docker Targets (if needed in future)
# =============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(PROJECT_NAME):latest .
	@echo "$(GREEN)✓ Docker image built$(NC)"

# =============================================================================
# Default Target
# =============================================================================

.DEFAULT_GOAL := help 