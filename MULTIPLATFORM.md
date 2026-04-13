# Multiplatform Build Support

This document explains the platform-specific build guidance for the Golang-Vulkan-api library. The current verified repo state passes `go build ./...`, `go test ./...`, and `go test -race ./...` on Linux with `libvulkan-dev` installed. For detailed installation instructions, see the [main README](../README.md#platform-specific-setup).

## Architecture

The library uses Go build tags to provide platform-specific CGO directives:

- `cgo_linux.go`: Linux-specific build configuration using pkg-config and the system Vulkan development package
- `cgo_darwin.go`: macOS-specific build configuration using pkg-config and the Vulkan SDK
- `cgo_windows.go`: Windows-specific build configuration using the Vulkan SDK and `vulkan-1.lib`
- `cgo_unix.go`: Fallback for other Unix-like systems (FreeBSD, OpenBSD, etc.) with pkg-config and Vulkan libraries available

## Platform-Specific Notes

### Linux
Verified on Linux with `libvulkan-dev` installed. `pkg-config` is still used when available:
```bash
# Install required packages
sudo apt-get install libvulkan-dev pkg-config
# Or for other distributions:
sudo yum install vulkan-devel pkgconf-pkg-config
sudo pacman -S vulkan-headers vulkan-validation-layers pkg-config

# Build and test
go build ./...
go test ./...
go test -race ./...
```

### Windows
Not reverified in the current repo state. Uses the Vulkan SDK and `vulkan-1.lib`:
```cmd
# Install Vulkan SDK from https://vulkan.lunarg.com/
# Make sure vulkan-1.lib is in your library path

# For custom SDK locations, you may need:
# set CGO_CFLAGS=-I"C:\VulkanSDK\1.3.290.0\Include"
# set CGO_LDFLAGS=-L"C:\VulkanSDK\1.3.290.0\Lib" -lvulkan-1

# Build
go build
```

### macOS
Not reverified in the current repo state. Uses pkg-config with MoltenVK support:
```bash
# Install Vulkan SDK with MoltenVK
# Install pkg-config if needed
brew install pkg-config

# Build
go build
```

### Other Unix Systems
Not reverified in the current repo state. Uses pkg-config as the fallback:
```bash
# Install Vulkan development libraries for your system
# Build
go build
```

## Testing Multiplatform Support

Run the included test script when validating additional targets:
```bash
./test_multiplatform.sh
```

This remains a documentation/example check; only the Linux build/test/race path has been verified in the current repo state.

## Documentation Gaps

Cross-platform verification is still incomplete. Windows, macOS, and other Unix-like platforms should be revalidated before this guide is treated as fully current.
