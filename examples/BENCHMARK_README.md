# Vulkan Graphics Benchmark

A comprehensive graphics benchmark application built using the Golang Vulkan API. This benchmark renders a dynamic scene and provides real-time performance monitoring including FPS, GPU temperatures, and clock speeds.

## Features

- **Dynamic Scene Rendering**: Renders an animated scene with rotating geometry
- **Real-time FPS Monitoring**: Displays current and average frame rates
- **GPU Monitoring**: Shows GPU temperature, clock speeds, and memory usage (NVIDIA GPUs)
- **Cross-platform Support**: Works on Linux, Windows, and macOS
- **Graceful Fallbacks**: Functions even without GPU drivers for testing

## Requirements

### System Requirements
- Go 1.24.6 or later
- Vulkan SDK and drivers
- pkg-config

### For GPU Monitoring (Optional)
- NVIDIA GPU with drivers (for hardware monitoring)
- NVIDIA Management Library (NVML)

## Installation

1. Install Vulkan development libraries:

```bash
# Ubuntu/Debian
sudo apt install libvulkan-dev vulkan-tools pkg-config

# Fedora/RHEL
sudo dnf install vulkan-devel vulkan-tools pkgconfig

# macOS (with Homebrew)
brew install vulkan-headers vulkan-loader molten-vk
```

2. Build and run the benchmark:

```bash
cd examples
go run graphics_benchmark.go
```

## Usage

### Running the Benchmark

```bash
go run examples/graphics_benchmark.go
```

The benchmark will:
1. Initialize Vulkan (if drivers are available)
2. Set up GPU monitoring (if NVIDIA GPU detected)
3. Start rendering a dynamic scene at 60 FPS target
4. Display real-time statistics including:
   - Runtime duration
   - Total frames rendered
   - Current and average FPS
   - Scene rotation angle
   - GPU temperature, clock speeds, and memory usage
   - System memory usage

### Running Tests

```bash
# Run functionality tests
go test -v graphics_benchmark_test.go graphics_benchmark.go

# Run performance benchmarks
go test -bench=. graphics_benchmark_test.go graphics_benchmark.go
```

## Output Example

```
Vulkan Graphics Benchmark - Live Stats
=====================================
Runtime: 15s
Total Frames: 901
Average FPS: 60.1
Current FPS: 59.8
Rotation Angle: 3.14 radians

GPU Statistics:
Temperature: 65°C
Graphics Clock: 1800 MHz
Memory Clock: 7000 MHz
GPU Utilization: 85%
Memory Used: 2048.0 MB / 8192.0 MB (25.0%)

System Memory: 45.2 MB allocated
Goroutines: 1
```

## Architecture

The benchmark consists of several key components:

### BenchmarkApp Structure
- **Vulkan Context**: Instance, physical device, logical device, and command pool
- **Performance Monitoring**: Frame counting, FPS calculation, and timing
- **Scene Animation**: Dynamic rotation and rendering simulation
- **GPU Monitoring**: NVIDIA GPU statistics via NVML

### Key Functions
- `initVulkan()`: Sets up Vulkan context and enumerates devices
- `renderFrame()`: Simulates rendering work and updates scene state
- `getGPUStats()`: Retrieves real-time GPU performance data
- `displayStats()`: Shows formatted performance statistics

### Graceful Degradation
The benchmark handles various scenarios gracefully:
- **No Vulkan Drivers**: Runs in simulation mode
- **No NVIDIA GPU**: Skips GPU monitoring
- **Missing NVML**: Uses fallback monitoring
- **Performance Issues**: Maintains target frame rate

## Extending the Benchmark

### Adding New GPU Vendors

To add support for AMD or Intel GPUs:

1. Add vendor-specific monitoring libraries
2. Implement vendor detection in `initGPUMonitoring()`
3. Add vendor-specific stats collection in `getGPUStats()`

### Customizing the Scene

To modify the rendered scene:

1. Update `simulateRenderingWork()` for different workloads
2. Modify animation parameters in `renderFrame()`
3. Add new performance metrics as needed

### Performance Tuning

- Adjust `targetFPS` for different frame rate targets
- Modify `workAmount` calculation for varying GPU loads
- Add more complex animations for stress testing

## Testing

The benchmark includes comprehensive tests:

- **Unit Tests**: Verify core functionality
- **Integration Tests**: Test Vulkan initialization and GPU monitoring
- **Performance Tests**: Benchmark rendering performance
- **Simulation Tests**: Validate behavior without hardware

## Troubleshooting

### Common Issues

1. **VK_ERROR_INCOMPATIBLE_DRIVER**: No Vulkan drivers installed
2. **ERROR_LIBRARY_NOT_FOUND**: NVML library not available
3. **No GPU Stats**: Non-NVIDIA GPU or missing drivers

### Solutions

1. Install appropriate GPU drivers and Vulkan SDK
2. For NVIDIA: Install CUDA toolkit or NVIDIA drivers
3. Run in simulation mode for testing without hardware

## Performance Notes

- Target frame rate: 60 FPS
- Rendering simulation scales with scene complexity
- GPU monitoring adds minimal overhead (~1% CPU usage)
- Memory usage typically under 50MB for the benchmark itself

## License

This benchmark is part of the Golang Vulkan API project and follows the same licensing terms.