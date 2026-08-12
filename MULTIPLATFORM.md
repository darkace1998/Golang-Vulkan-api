# Multiplatform Build Support

This document explains the platform-specific build guidance for the Golang-Vulkan-api library. The current verified repo state passes `go build ./...`, `go test ./...`, and `go test -race ./...` on Linux with `libvulkan-dev` installed. For detailed installation instructions, see the [main README](../README.md#platform-specific-setup).

## Architecture

### Supported CPU architectures: 64-bit only

The library models Vulkan non-dispatchable handles (buffers, images, semaphores, ...)
as pointers, which matches `VK_DEFINE_NON_DISPATCHABLE_HANDLE` only on 64-bit
targets. On 32-bit architectures Vulkan defines those handles as `uint64_t`, so
the package does not compile there. Supported `GOARCH` values are `amd64`,
`arm64`, `riscv64`, `loong64`, `ppc64le`, and `s390x`; building for a 32-bit
architecture (e.g. `386`, `arm`) fails with an explicit error pointing at this
document (see `arch_unsupported.go`).

### Platform build tags

The library uses Go build tags to provide platform-specific CGO directives:

- `cgo_linux.go`: Linux-specific build configuration using pkg-config and the system Vulkan development package
- `cgo_darwin.go`: macOS-specific build configuration using pkg-config and the Vulkan SDK
- `cgo_windows.go`: Windows-specific build configuration using the Vulkan SDK and `vulkan-1.lib`
- `cgo_unix.go`: Fallback for other Unix-like systems (FreeBSD, OpenBSD, etc.) with pkg-config and Vulkan libraries available

## Platform-Specific Notes

### Linux
Verified on Linux with `libvulkan-dev`, `libx11-dev`, and `libwayland-dev` installed. `pkg-config` is still used when available:
```bash
# Install required packages (X11/Wayland headers are needed by the default build)
sudo apt-get install libvulkan-dev pkg-config libx11-dev libwayland-dev
# Or for other distributions:
sudo yum install vulkan-devel pkgconf-pkg-config libX11-devel wayland-devel
sudo pacman -S vulkan-headers vulkan-validation-layers pkg-config libx11 wayland

# Build and test
go build ./...
go test ./...
go test -race ./...
```

Headless machines (servers, CI, slim containers) can build without any
display-server headers — only `libvulkan-dev` is required:

```bash
go build -tags vk_headless ./...    # no X11/Wayland surface support
go build -tags vk_no_xlib ./...     # skip only the X11 backend
go build -tags vk_no_wayland ./...  # skip only the Wayland backend
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

Run the following commands to validate additional targets:
```bash
# Linux (amd64)
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build ./...

# Windows (amd64)
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build ./...

# macOS (amd64)
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build ./...

# macOS (arm64)
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build ./...
```

This remains a documentation/example check; only the Linux build/test/race path has been verified in the current repo state.

## Documentation Gaps

Cross-platform verification is still incomplete. Windows, macOS, and other Unix-like platforms should be revalidated before this guide is treated as fully current.
