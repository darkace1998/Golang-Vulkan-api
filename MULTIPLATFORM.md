# Multiplatform Build Support

This document explains how the Golang-Vulkan-api library supports multiple platforms through conditional compilation. For detailed installation instructions, see the [main README](../README.md#platform-specific-setup).

## Architecture

The library uses Go build tags to provide platform-specific CGO directives:

- `cgo_linux.go`: Linux-specific build configuration using pkg-config
- `cgo_darwin.go`: macOS-specific build configuration using pkg-config  
- `cgo_windows.go`: Windows-specific build configuration using direct linking
- `cgo_unix.go`: Fallback for other Unix-like systems (FreeBSD, OpenBSD, etc.)

## Platform-Specific Notes

### Linux
Uses pkg-config to find Vulkan libraries:
```bash
# Install required packages
sudo apt-get install libvulkan-dev pkg-config
# Or for other distributions:
sudo yum install vulkan-devel pkgconf-pkg-config
sudo pacman -S vulkan-headers vulkan-validation-layers pkg-config

# Build
go build
```

### Windows
Uses direct linking to vulkan-1.lib:
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
Uses pkg-config with MoltenVK support:
```bash
# Install Vulkan SDK with MoltenVK
# Install pkg-config if needed
brew install pkg-config

# Build
go build
```

### Other Unix Systems
Uses pkg-config as fallback:
```bash
# Install Vulkan development libraries for your system
# Build
go build
```

## Testing Multiplatform Support

Run the included test script:
```bash
./test_multiplatform.sh
```

This will verify that the build tags are correctly configured for each platform.