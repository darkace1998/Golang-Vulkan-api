# Graphics Benchmark

GPU stress testing and benchmarking example for the Golang Vulkan API.

## Current state

- Entry point: `graphics_benchmark.go`
- Supports `stress` and `benchmark` modes
- Supports `-sim` for non-Vulkan environments
- Uses best-effort GPU sensor data when available

## Quick Start

```bash
go build -o gpu_stress_test ./examples/benchmark
./gpu_stress_test -help
./gpu_stress_test
```

## Features

- Stress testing with configurable quality levels
- Benchmark mode with performance scoring
- Resolution presets and custom `WIDTHxHEIGHT` input
- Best-effort monitoring of temperature, clocks, power, and fan data
- CSV export and output directory support
- Simulation mode for systems without Vulkan

## Usage

```bash
# Basic stress test
./gpu_stress_test

# 4K benchmark for 5 minutes with CSV export
./gpu_stress_test -mode=benchmark -resolution=4K -duration=5m -output=./results -csv

# Ultra quality stress test with artifact detection
./gpu_stress_test -quality=ultra -artifacts

# List preset resolutions
./gpu_stress_test -list-res

# Help
./gpu_stress_test -help
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-mode` | `stress` or `benchmark` | `stress` |
| `-quality` | `low`, `medium`, `high`, `ultra` | `high` |
| `-resolution` | `720p`, `1080p`, `1440p`, `4K`, or `WIDTHxHEIGHT` | `1080p` |
| `-duration` | Test duration (`0` for infinite stress test) | `0` |
| `-fps` | Target FPS for the test | `60` |
| `-output` | Output directory for logs, reports, and CSV exports | empty |
| `-csv` | Export performance data to CSV | `false` |
| `-artifacts` | Enable artifact detection mode | `false` |
| `-sim` | Force simulation mode (no Vulkan) | `false` |
| `-list-res` | List available resolutions | `false` |
| `-verbose` | Enable verbose logging | `false` |
| `-help` | Show detailed help information | `false` |

## Requirements

- Go 1.22+
- Vulkan SDK or development libraries
- Linux: `libvulkan-dev`
