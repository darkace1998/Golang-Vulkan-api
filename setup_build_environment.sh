#!/bin/bash
# setup_build_environment.sh
# Sets up the build environment for Golang-Vulkan-api with Vulkan headers

set -e  # Exit on error

echo "========================================="
echo "Golang-Vulkan-api Build Environment Setup"
echo "========================================="

# Detect OS
OS_TYPE=$(uname -s)
echo "Detected OS: $OS_TYPE"

if [[ "$OS_TYPE" == "Linux" ]]; then
    echo ""
    echo "Installing dependencies for Linux..."
    sudo apt-get update
    sudo apt-get install -y \
        vulkan-headers \
        libvulkan-dev \
        vulkan-tools \
        pkg-config \
        build-essential
    
    echo ""
    echo "Verifying Vulkan installation..."
    pkg-config --cflags vulkan
    pkg-config --libs vulkan
    echo "✓ Vulkan headers installed successfully"
    
elif [[ "$OS_TYPE" == "Darwin" ]]; then
    echo ""
    echo "Installing dependencies for macOS..."
    brew install vulkan-headers
    
    echo ""
    echo "Verifying Vulkan installation..."
    pkg-config --cflags vulkan
    pkg-config --libs vulkan
    echo "✓ Vulkan headers installed successfully"
    
else
    echo "Unsupported OS: $OS_TYPE"
    echo "Please manually install Vulkan SDK from https://www.lunarg.com/vulkan-sdk/"
    exit 1
fi

echo ""
echo "Checking Go installation..."
go version

echo ""
echo "========================================="
echo "✓ Build environment ready!"
echo "========================================="
echo ""
echo "To build the project:"
echo "  go build ./..."
echo ""
echo "To run tests:"
echo "  go test ./..."
echo ""
echo "To run examples:"
echo "  go run ./examples/basic/main.go"
echo "  go run ./examples/simple/main.go"
echo "  go run ./examples/compute/main.go"
echo ""
