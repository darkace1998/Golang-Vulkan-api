# Windows Build Instructions for GPU Benchmark Example

## Problem Solved
The original `graphics_benchmark.go` example failed to build on Windows due to:
- NVML dependency requiring `dlfcn.h` (Unix-only header)
- Complex Vulkan SDK configuration requirements

## Solutions

### Option 1: Hardware-Accelerated Mode (Requires Vulkan SDK)
For users who want **real GPU hardware acceleration** on Windows:

1. **Install Vulkan SDK**: Download from [LunarG Vulkan SDK](https://vulkan.lunarg.com/)
2. **Set environment variables** (if not done automatically):
   ```cmd
   set VULKAN_SDK=C:\VulkanSDK\1.3.x.x
   set PATH=%PATH%;%VULKAN_SDK%\Bin
   ```
3. **Build with hardware acceleration**:
   ```bash
   cd examples
   go build -tags vulkan_hardware -o bench.exe graphics_benchmark_hardware.go gpu_monitoring_windows.go
   ```

   **Alternative** (if above fails):
   ```bash
   cd examples
   go build -o bench.exe graphics_benchmark.go gpu_monitoring_windows.go
   ```

This gives you:
- ✅ **Real Vulkan GPU acceleration**
- ✅ **Hardware GPU monitoring** (simulated on Windows)
- ✅ Both stress testing and benchmark modes
- ✅ Full performance scoring
- ✅ **Automatic fallback** to simulation if Vulkan setup fails

### Option 2: Simulation Mode (No Dependencies Required)
For users who want **cross-platform compatibility** without Vulkan SDK:

```bash
cd examples
go build -o bench.exe graphics_benchmark_windows.go gpu_monitoring_windows.go
```

This gives you:
- ✅ **CPU-based simulation** (no GPU hardware required)
- ✅ **Cross-platform compatibility**
- ✅ Both stress testing and benchmark modes
- ✅ Realistic simulated GPU statistics
- ✅ Clear guidance to upgrade to hardware mode

## Usage Examples:

### Hardware-Accelerated Mode (Option 1)
```bash
# Hardware-accelerated 60-second benchmark
bench.exe -mode=benchmark -duration=60s -quality=high

# If hardware acceleration fails, it automatically falls back to simulation
# Force simulation mode if desired
bench.exe -mode=benchmark -duration=60s -sim
```

### Simulation Mode (Option 2)
```bash
# Simulation benchmark (always runs in simulation mode)
bench.exe -mode=benchmark -duration=60s -quality=high

# Note: This version will show a message about hardware acceleration
# and recommend using Option 1 for real GPU acceleration
```

## Troubleshooting Hardware Mode

If hardware mode fails to build or run:

1. **Build Error - Vulkan not found**: Install Vulkan SDK and ensure environment variables are set
2. **Runtime Error - Falls back to simulation**: This is normal behavior when Vulkan setup isn't complete
3. **Want simulation anyway**: Use the `-sim` flag to force simulation mode

**The hardware version is designed to gracefully fallback to simulation mode if Vulkan setup isn't working properly.**

### Available Options:
- **Test Modes**: `stress` (infinite) or `benchmark` (timed)
- **Quality Levels**: `low`, `medium`, `high`, `ultra`
- **Resolutions**: `720p`, `1080p`, `1440p`, `4K`, or custom `WIDTHxHEIGHT`
- **Hardware Control**: `-sim` flag to force simulation mode (Option 1 only)
- **Export**: `-csv` flag to export performance data
- **Help**: `bench.exe -help` for detailed options

## Which Option Should I Choose?

### Choose Option 1 (Hardware Acceleration) if:
- You have or can install Vulkan SDK
- You want maximum performance testing
- You need real GPU hardware validation
- You're developing GPU applications

### Choose Option 2 (Simulation Mode) if:
- You don't want to install Vulkan SDK
- You need cross-platform compatibility
- You're running in CI/CD environments
- You want quick setup without dependencies

## Understanding the Modes:

### Execution Modes:
- **Hardware Mode**: Uses Vulkan API for real GPU acceleration
- **Simulation Mode**: CPU-based cross-platform performance testing

### Test Modes:
- **Stress Test**: Runs indefinitely until you stop it (Ctrl+C) - good for thermal testing
- **Benchmark**: Runs for a fixed duration and provides a final performance score

Both test modes work in both execution modes, providing meaningful performance testing and system stability validation.