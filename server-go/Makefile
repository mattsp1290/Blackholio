# Blackholio Server Go - Makefile
# Provides convenient build targets for development and deployment

.PHONY: all build wasm test clean demo help install-deps generate publish logs

# Default target
all: test build

# Build regular Go binary for testing/demo
build:
	@echo "Building Go binary..."
	go build -o blackholio-demo .
	@echo "Built: blackholio-demo"

# Build WASM module for SpacetimeDB
wasm:
	@echo "Building WASM module..."
	GOOS=js GOARCH=wasm go build -o blackholio.wasm .
	@echo "Built: blackholio.wasm ($$(stat -f%z blackholio.wasm 2>/dev/null || stat -c%s blackholio.wasm) bytes)"

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test ./... -bench=. -benchmem

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f blackholio-demo blackholio.wasm coverage.out coverage.html
	@echo "Clean completed"

# Run demo application
demo: build
	@echo "Running demo application..."
	./blackholio-demo

# Generate WASM and client bindings
generate: wasm
	@echo "Generating client bindings..."
	./generate.sh

# Publish to SpacetimeDB
publish: wasm
	@echo "Publishing to SpacetimeDB..."
	./publish.sh

# View SpacetimeDB logs
logs:
	@echo "Viewing SpacetimeDB logs..."
	./logs.sh

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Check code quality
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Full development cycle: format, lint, test, build
dev: fmt lint test build

# Full deployment cycle: test, build WASM, publish
deploy: test wasm publish

# Display help
help:
	@echo "Blackholio Server Go - Available Make Targets:"
	@echo ""
	@echo "Development:"
	@echo "  build          - Build regular Go binary for testing"
	@echo "  wasm          - Build WASM module for SpacetimeDB"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  bench         - Run benchmark tests"
	@echo "  demo          - Run demo application"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  dev           - Full development cycle (fmt + lint + test + build)"
	@echo ""
	@echo "Deployment:"
	@echo "  generate      - Generate WASM and client bindings"
	@echo "  publish       - Publish to SpacetimeDB"
	@echo "  logs          - View SpacetimeDB logs"
	@echo "  deploy        - Full deployment cycle (test + wasm + publish)"
	@echo ""
	@echo "Setup:"
	@echo "  install-deps  - Install Go dependencies"
	@echo "  help          - Show this help message" 