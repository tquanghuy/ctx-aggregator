.PHONY: all test lint fmt clean examples coverage bench check help

# Default target
all: check

# Run all checks (test, lint, fmt)
check: test lint fmt
	@echo "All checks passed!"

# Run tests with race detection
test:
	go test -v -race ./...

# Run tests with coverage
coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Run all examples
examples:
	@echo "Running basic example..."
	@cd examples/basic && go run main.go
	@echo ""
	@echo "Running concurrent example..."
	@cd examples/concurrent && go run main.go
	@echo ""
	@echo "Running multiple aggregators example..."
	@cd examples/multiple && go run main.go

# Build examples (verify they compile)
build-examples:
	@echo "Building examples..."
	@cd examples/basic && go build -o /dev/null .
	@cd examples/concurrent && go build -o /dev/null .
	@cd examples/multiple && go build -o /dev/null .
	@echo "All examples built successfully!"

# Clean build artifacts
clean:
	go clean
	rm -f coverage.out coverage.html

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Display help
help:
	@echo "Available targets:"
	@echo "  make all           - Run all checks (default)"
	@echo "  make check         - Run test, lint, and fmt"
	@echo "  make test          - Run tests with race detection"
	@echo "  make coverage      - Generate test coverage report"
	@echo "  make bench         - Run benchmarks"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make fmt           - Format code with go fmt and goimports"
	@echo "  make examples      - Run all example programs"
	@echo "  make build-examples - Build all examples (verify compilation)"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make install-tools - Install development tools"
	@echo "  make help          - Display this help message"
