.PHONY: help build build-verbose test clean setup lint format run-example run-video run-compute run-benchmark

help:
	@echo "Golang-Vulkan-api Build Targets"
	@echo "================================"
	@echo ""
	@echo "setup              - Install Vulkan headers and dependencies"
	@echo "build              - Build the entire project"
	@echo "build-verbose      - Build with verbose output"
	@echo "test               - Run all tests"
	@echo "test-verbose       - Run tests with verbose output"
	@echo "clean              - Clean build cache"
	@echo "lint               - Run linters (requires golangci-lint)"
	@echo "format             - Format code with gofmt"
	@echo "run-example        - Run basic example"
	@echo "run-simple         - Run simple example"
	@echo "run-compute        - Run compute example"
	@echo "run-video          - Run video example"
	@echo "run-type           - Run type example"
	@echo "run-benchmark      - Run benchmarks"
	@echo ""

setup:
	@echo "Installing Vulkan headers and dependencies..."
	@bash ./setup_build_environment.sh

build:
	@echo "Building project..."
	@go build -v ./...
	@echo "✓ Build complete"

build-verbose:
	@echo "Building project with verbose output..."
	@go build -v -x ./...
	@echo "✓ Build complete"

test:
	@echo "Running tests..."
	@go test ./...
	@echo "✓ Tests complete"

test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...
	@echo "✓ Tests complete"

clean:
	@echo "Cleaning build cache..."
	@go clean -cache
	@echo "✓ Cache cleaned"

lint:
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "✓ Lint check complete"

format:
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "✓ Code formatted"

run-example:
	@echo "Running basic example..."
	@go run ./examples/basic/main.go

run-simple:
	@echo "Running simple example..."
	@go run ./examples/simple/main.go

run-compute:
	@echo "Running compute example..."
	@go run ./examples/compute/main.go

run-video:
	@echo "Running video example..."
	@go run ./examples/video/main.go

run-type:
	@echo "Running type example..."
	@go run ./examples/type/main.go

run-benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./examples/benchmark/...

verify-vulkan:
	@echo "Verifying Vulkan installation..."
	@pkg-config --cflags vulkan
	@pkg-config --libs vulkan
	@echo "✓ Vulkan headers found"

all: setup build test
	@echo ""
	@echo "✓ Full build pipeline complete"
