# Graphics Benchmark

A GPU stress testing and benchmarking application using the Golang Vulkan API.

## Quick Start

```bash
cd examples
go build -o gpu_stress_test graphics_benchmark.go
./gpu_stress_test
```

## Features

- **Stress Testing**: Intensive GPU workload with thermal monitoring
- **Benchmarking**: Performance scoring and detailed reports
- **Real-time Monitoring**: FPS, GPU temperature, memory usage
- **Quality Levels**: Low, Medium, High, Ultra GPU stress levels
- **Cross-platform**: Works on Linux, Windows, macOS

## Usage

```bash
# Basic stress test
./gpu_stress_test

# 4K benchmark for 5 minutes
./gpu_stress_test -mode=benchmark -resolution=4K -duration=5m

# Ultra quality stress test
./gpu_stress_test -quality=ultra

# Help
./gpu_stress_test -help
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-mode` | 'stress' or 'benchmark' | stress |
| `-quality` | 'low', 'medium', 'high', 'ultra' | high |
| `-resolution` | Resolution (e.g., '4K', '1920x1080') | 1080p |
| `-duration` | Test duration (0 for infinite) | 0 |
| `-csv` | Export performance data to CSV | false |
| `-verbose` | Enable verbose logging | false |

## Requirements

- Go 1.19+
- Vulkan SDK and drivers
- Optional: NVIDIA GPU for enhanced monitoring