# TypeScript-like Go Implementation
# Makefile for building and managing the project

.PHONY: help build run test clean fmt lint vet deps install dev watch demo

# Default target
help: ## Show help message
	@echo "TypeScript-like Go Implementation"
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the project
build: ## Build the project
	@echo "Building TypeScript-like Go project..."
	go build -o bin/typescript-golang .
	@echo "Build complete: bin/typescript-golang"

# Run the project
run: ## Run the project
	@echo "Running TypeScript-like Go project..."
	go run .

# Run with race detection
run-race: ## Run with race detection
	@echo "Running with race detection..."
	go run -race .

# Install dependencies
deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run benchmarks
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Format code
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Lint code
lint: ## Lint code
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Vet code
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Check code quality
check: fmt vet lint ## Run all code quality checks

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Install the project
install: ## Install the project
	@echo "Installing TypeScript-like Go project..."
	go install .

# Development mode (format, vet, test, build)
dev: fmt vet test build ## Development workflow: format, vet, test, build
	@echo "Development workflow complete!"

# Watch for changes and run development workflow
watch: ## Watch for changes and run dev workflow
	@which fswatch > /dev/null || (echo "fswatch not found. Install with: brew install fswatch (macOS) or apt-get install inotify-tools (Linux)" && exit 1)
	@echo "Watching for changes... (Ctrl+C to stop)"
	fswatch -o . --exclude='bin/' --exclude='coverage.*' --exclude='.git/' | xargs -n1 -I{} make dev

# Run demo examples
demo: ## Run demo examples
	@echo "Running TypeScript-like Go demos..."
	@go run . 2>/dev/null || true
	@echo "\n=== Array Utilities Demo ==="
	@go run examples/arrays_demo.go 2>/dev/null || echo "Create examples/arrays_demo.go to see array utilities in action"
	@echo "\n=== Async/Promise Demo ==="
	@go run examples/async_demo.go 2>/dev/null || echo "Create examples/async_demo.go to see async/promise utilities in action"
	@echo "\n=== Classes Demo ==="
	@go run examples/classes_demo.go 2>/dev/null || echo "Create examples/classes_demo.go to see class-like structures in action"

# Initialize project for development
init: ## Initialize project for development
	@echo "Initializing TypeScript-like Go project for development..."
	go mod tidy
	mkdir -p bin examples tests docs
	@echo "Project initialized!"

# Generate documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@which godoc > /dev/null || (echo "Installing godoc..." && go install golang.org/x/tools/cmd/godoc@latest)
	@echo "Starting documentation server at http://localhost:6060"
	@echo "Visit http://localhost:6060/pkg/typescript-golang/ to view documentation"
	godoc -http=:6060

# Update dependencies
update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Security check
security: ## Run security checks
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	@echo "Running security checks..."
	gosec ./...

# Performance profiling
profile: ## Run with CPU profiling
	@echo "Running with CPU profiling..."
	go run -cpuprofile=cpu.prof .
	@echo "CPU profile saved to cpu.prof"
	@echo "Use 'go tool pprof cpu.prof' to analyze"

# Memory profiling
profile-mem: ## Run with memory profiling
	@echo "Running with memory profiling..."
	go run -memprofile=mem.prof .
	@echo "Memory profile saved to mem.prof"
	@echo "Use 'go tool pprof mem.prof' to analyze"

# Docker build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t typescript-golang .

# Docker run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm typescript-golang

# Release build (optimized)
release: ## Build optimized release version
	@echo "Building release version..."
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/typescript-golang .
	@echo "Release build complete: bin/typescript-golang"

# Cross-compile for different platforms
cross-compile: ## Cross-compile for multiple platforms
	@echo "Cross-compiling for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/typescript-golang-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/typescript-golang-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/typescript-golang-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o bin/typescript-golang-darwin-arm64 .
	@echo "Cross-compilation complete! Binaries available in bin/"

# Show project structure
tree: ## Show project structure
	@echo "Project structure:"
	@tree -I 'bin|.git|coverage.*' . 2>/dev/null || find . -type f -not -path './bin/*' -not -path './.git/*' -not -name 'coverage.*' | sort

# Show project statistics
stats: ## Show project statistics
	@echo "Project Statistics:"
	@echo "=================="
	@echo "Go files: $(shell find . -name '*.go' -not -path './bin/*' | wc -l)"
	@echo "Total lines: $(shell find . -name '*.go' -not -path './bin/*' -exec wc -l {} + | tail -1 | awk '{print $$1}')"
	@echo "Packages: $(shell go list ./... | wc -l)"
	@echo "Dependencies: $(shell go list -m all | wc -l)"

# Package size analysis
analyze: ## Analyze binary size
	@echo "Analyzing binary size..."
	@if [ -f bin/typescript-golang ]; then \
		ls -lah bin/typescript-golang; \
		file bin/typescript-golang; \
	else \
		echo "Binary not found. Run 'make build' first."; \
	fi

# All-in-one quality check
quality: clean fmt vet lint test build ## Complete quality check workflow
	@echo "All quality checks passed!"

# CI/CD pipeline simulation
ci: deps fmt vet lint test build ## Simulate CI/CD pipeline
	@echo "CI pipeline completed successfully!"