# Windows Build Fix for GPU Benchmark Example

## Problem
The original `graphics_benchmark.go` example failed to build on Windows due to:
- NVML dependency requiring `dlfcn.h` (Unix-only header)
- CGO compilation issues with Vulkan SDK requirements

## Solution
Created platform-specific GPU monitoring and a Windows-compatible example:

### Files Added:
1. **`gpu_monitoring_unix.go`** - Unix/Linux NVML-based GPU monitoring
2. **`gpu_monitoring_windows.go`** - Windows-compatible simulated GPU monitoring  
3. **`graphics_benchmark_windows.go`** - Windows-specific example without Vulkan dependency

### How to Build on Windows:

#### Option 1: Windows-specific example (Recommended)
```bash
go build -o bench.exe graphics_benchmark_windows.go gpu_monitoring_windows.go
```

#### Option 2: Original example (requires Vulkan SDK)
```bash
go build -o bench.exe graphics_benchmark.go gpu_monitoring_windows.go
```

### Features:
- ✅ Builds successfully on Windows without external dependencies
- ✅ Provides realistic simulated GPU statistics
- ✅ Full benchmark and stress testing functionality
- ✅ CSV export and performance analysis
- ✅ Cross-platform compatible

### Usage:
```bash
# Basic stress test
bench.exe

# 5-minute benchmark with CSV export
bench.exe -mode=benchmark -duration=5m -csv -output=./results

# Ultra quality stress test
bench.exe -quality=ultra -artifacts
```

The Windows version runs in simulation mode, providing CPU-based performance testing while maintaining the same interface and reporting capabilities as the hardware-accelerated version.