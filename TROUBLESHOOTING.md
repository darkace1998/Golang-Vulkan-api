# Troubleshooting Guide

This document lists common issues you might encounter when working with the `golang-vulkan-api` bindings, along with their solutions.

## Build and Compilation Errors

### `cgo: C compiler "gcc" not found` or similar CGO errors
**Symptom:** Build fails with errors indicating the C compiler is missing or CGO is disabled.
**Cause:** The library uses CGO to interface with the Vulkan C API. CGO requires a C compiler (like `gcc` or `clang`) to be installed and available in your `PATH`, and `CGO_ENABLED=1` environment variable must be set.
**Solution:**
1. Install a C compiler for your platform.
   - Ubuntu/Debian: `sudo apt-get install build-essential`
   - macOS: Install Xcode Command Line Tools (`xcode-select --install`)
   - Windows: Install MinGW-w64 or TDM-GCC.
2. Ensure `CGO_ENABLED=1` is set in your environment (it usually is by default if a compiler is found).

### `fatal error: vulkan/vulkan.h: No such file or directory`
**Symptom:** Build fails complaining about missing Vulkan headers.
**Cause:** The Vulkan development headers are not installed on your system.
**Solution:** Install the Vulkan development package for your platform.
- Ubuntu/Debian: `sudo apt-get install libvulkan-dev`
- Fedora/CentOS/RHEL: `sudo dnf install vulkan-headers vulkan-loader-devel`
- Arch Linux: `sudo pacman -S vulkan-headers`
- Windows/macOS: Install the Vulkan SDK from LunarG.

### `pkg-config: command not found` or package `vulkan` not found
**Symptom:** Build errors referencing `pkg-config`.
**Cause:** The Linux/macOS build relies on `pkg-config` to locate the Vulkan libraries.
**Solution:** Install `pkg-config`.
- Ubuntu/Debian: `sudo apt-get install pkg-config`
- macOS: `brew install pkg-config`

## Runtime and Test Errors

### `VK_ERROR_INCOMPATIBLE_DRIVER` or Tests Failing with "setup failed"
**Symptom:** Vulkan instance creation fails with `VK_ERROR_INCOMPATIBLE_DRIVER`, or tests crash/fail when trying to create a device, particularly in headless or CI environments.
**Cause:** The system lacks a Vulkan driver, or the requested API version (e.g., Vulkan 1.4) is not supported by the installed driver.
**Solution:**
- **In CI/Headless Environments:** Install a software Vulkan driver like `lavapipe` (Mesa) to allow Vulkan instance creation without a physical GPU.
  - Ubuntu: `sudo apt-get install mesa-vulkan-drivers vulkan-tools`
- **On Desktop:** Ensure your GPU drivers are up-to-date and support the requested Vulkan API version. You can verify your driver's capabilities by running `vulkaninfo`.

### `SIGABRT` or Segmentation Faults During Execution
**Symptom:** The program crashes abruptly with a segmentation fault (`SIGSEGV`) or abort signal (`SIGABRT`).
**Cause:** This is often caused by passing invalid pointers to the Vulkan C API via CGO. A very common mistake in Go bindings is passing the address of the first element of an empty slice (e.g., `&mySlice[0]` when `len(mySlice) == 0`). Go's slice memory representation can cause issues if accessed incorrectly when empty.
**Solution:**
- **Always check slice length before passing to CGO:**
  ```go
  var pData *C.float
  if len(goSlice) > 0 {
      pData = (*C.float)(&goSlice[0])
  } else {
      pData = nil // Or handle the zero-length case as required by the specific Vulkan API
  }
  ```
- Review the `ERROR_HANDLING.md` and ensure you are handling validation errors properly. Our library attempts to catch many of these before calling C, returning a `ValidationError`.

### Tests Output NVML CGO Deprecation Warnings
**Symptom:** When running `go test ./...`, you see warnings like:
```
# github.com/NVIDIA/go-nvml/pkg/nvml
... warning: passing ... to CGO ...
```
**Cause:** These warnings originate from the `go-nvml` dependency used in the benchmark examples for GPU monitoring.
**Solution:** **This is expected behavior and can be safely ignored.** It does not indicate a test failure or an issue with the core `golang-vulkan-api` library.

## Advanced Issues

- For issues related to recovering from GPU crashes (`VK_ERROR_DEVICE_LOST`), see **[ERROR_HANDLING.md](ERROR_HANDLING.md)**.
- For issues loading video extension functions (e.g., `vkCmdDecodeVideoKHR` is nil), refer to the thread safety and loading model in **[ARCHITECTURE_DIAGRAMS.md](ARCHITECTURE_DIAGRAMS.md)** and **[THREAD_SAFETY.md](THREAD_SAFETY.md)**.
