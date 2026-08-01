# Golang-Vulkan-api

A Go binding for Vulkan 1.3+ graphics and compute APIs with broad coverage of the features implemented in this repository.

## Overview

This library provides a type-safe Go interface to the Vulkan APIs used by the project examples and tests. It's designed to be used as a library for other Go projects that need low-level graphics and compute functionality.

## Verified status

- `libvulkan-dev` is installed on Linux
- `go build ./...` passes
- `go test ./...` passes
- `go test -race ./...` passes
- Repository is synced with `origin/main`

## Features

- ✅ **Vulkan 1.3 API Coverage**: Core Vulkan 1.3 functions and types used by this repository
- ✅ **Dynamic Rendering**: Modern renderpass-free rendering (VK_KHR_dynamic_rendering)
- ✅ **Synchronization2**: Enhanced timeline semaphores and submission (VK_KHR_synchronization2)
- ✅ **Extended Dynamic State**: More pipeline state that can be set dynamically
- ✅ **Private Data**: Associate private data with Vulkan objects
- ✅ **Maintenance4**: Enhanced buffer/image memory requirements without object creation
- ✅ **Type Safety**: Go-idiomatic types with proper error handling
- ✅ **Memory Management**: Safe memory allocation and management functions
- ✅ **Command Buffers**: Full command buffer recording and submission
- ✅ **Synchronization**: Semaphores, fences, and other sync primitives
- ✅ **Device Management**: Physical and logical device enumeration and creation
- ✅ **Buffer/Image Operations**: Buffer and image management helpers
- ✅ **Queue Operations**: Graphics, compute, and transfer queue support
- ✅ **Compute Shaders**: Compute pipeline support and example workloads
- ✅ **Storage Buffers**: Large dataset handling for compute operations
- ✅ **Dispatch Commands**: Efficient compute work group dispatching
- ✅ **Ray Tracing**: Ray tracing pipelines and commands (VK_KHR_ray_tracing_pipeline)
- ✅ **Acceleration Structures**: Create and manage acceleration structures (VK_KHR_acceleration_structure)
- ✅ **Platform Setup**: Linux, Windows, and macOS setup notes are included

## Video Codec Support 🎬

### Supported on Compatible Hardware

These codecs have ratified extensions for both operations on compatible hardware:

- **H.264 (AVC)** - VK_KHR_video_encode_h264 & VK_KHR_video_decode_h264
- **H.265 (HEVC)** - VK_KHR_video_encode_h265 & VK_KHR_video_decode_h265
- **AV1** - VK_KHR_video_encode_av1 & VK_KHR_video_decode_av1

Hardware-accelerated video encoding and decoding is available through Vulkan Video extensions on compatible GPUs and drivers.

### Checking Video Codec Support

Use the provided API to check which codecs are supported on your hardware:

```go
// Get supported video codecs for a physical device
supportedCodecs, err := vulkan.GetSupportedVideoCodecs(physicalDevice)
if err != nil {
    log.Fatal(err)
}

for _, codec := range supportedCodecs {
    fmt.Printf("Supported: %s\n", codec)
}
```

See `examples/video/main.go` for a working example that detects and displays supported video codecs on your system.

**Note**: Actual hardware support depends on your GPU model and driver version. Extension availability does not guarantee hardware acceleration.

## Vulkan 1.3 Features

### Dynamic Rendering
Replace traditional render passes with flexible dynamic rendering:
```go
renderingInfo := &vulkan.RenderingInfo{
    RenderArea: vulkan.Rect2D{
        Offset: vulkan.Offset2D{X: 0, Y: 0}, 
        Extent: vulkan.Extent2D{Width: 800, Height: 600},
    },
    LayerCount: 1,
    ColorAttachments: []vulkan.RenderingAttachmentInfo{
        {
            ImageView:   colorImageView,
            ImageLayout: vulkan.ImageLayoutColorAttachmentOptimal,
            LoadOp:      vulkan.AttachmentLoadOpClear,
            StoreOp:     vulkan.AttachmentStoreOpStore,
        },
    },
}

vulkan.CmdBeginRendering(commandBuffer, renderingInfo)
// Draw commands here
vulkan.CmdEndRendering(commandBuffer)
```

### Synchronization2 (Enhanced Timeline Semaphores)
Modern submission with enhanced synchronization:
```go
submitInfo := []vulkan.SubmitInfo2{
    {
        CommandBufferInfos: []vulkan.CommandBufferSubmitInfo{
            {CommandBuffer: commandBuffer, DeviceMask: 0},
        },
        WaitSemaphoreInfos: []vulkan.SemaphoreSubmitInfo{
            {
                Semaphore: waitSemaphore,
                Value:     waitValue,
                StageMask: vulkan.PipelineStage2FragmentShader,
            },
        },
    },
}

err := vulkan.QueueSubmit2(queue, submitInfo, fence)
```

### Extended Dynamic State
Set more pipeline state dynamically:
```go
vulkan.CmdSetCullMode(commandBuffer, vulkan.CullModeBack)
vulkan.CmdSetFrontFace(commandBuffer, vulkan.FrontFaceCounterClockwise)
vulkan.CmdSetPrimitiveTopology(commandBuffer, vulkan.PrimitiveTopologyTriangleList)
vulkan.CmdSetDepthTestEnable(commandBuffer, true)
vulkan.CmdSetDepthCompareOp(commandBuffer, vulkan.CompareOpLess)
```

### Private Data
Associate application data with Vulkan objects:
```go
slot, err := vulkan.CreatePrivateDataSlot(device, &vulkan.PrivateDataSlotCreateInfo{})
err = vulkan.SetPrivateData(device, vulkan.ObjectTypeBuffer, uint64(buffer), slot, myData)
retrievedData := vulkan.GetPrivateData(device, vulkan.ObjectTypeBuffer, uint64(buffer), slot)
```

### Maintenance4
Get memory requirements without creating objects:
```go
memReqs := vulkan.GetDeviceBufferMemoryRequirements(device, &vulkan.BufferCreateInfo{
    Size:  1024 * 1024, // 1MB buffer
    Usage: vulkan.BufferUsageStorageBufferBit,
})

imageMemReqs := vulkan.GetDeviceImageMemoryRequirements(device, &vulkan.ImageCreateInfo{
    ImageType: vulkan.ImageType2D,
    Format:    vulkan.FormatR8G8B8A8Unorm,
    Extent:    vulkan.Extent3D{Width: 512, Height: 512, Depth: 1},
    Usage:     vulkan.ImageUsageColorAttachmentBit,
})
```

## Requirements

- Go 1.22 or later
- CGO enabled
- Vulkan SDK or development libraries installed (Linux: `libvulkan-dev`)
  - Linux: `libvulkan-dev` package. Testing requires `mesa-vulkan-drivers vulkan-tools libwayland-dev libx11-dev`.
  - Windows: Vulkan SDK from LunarG
  - macOS: Vulkan SDK with MoltenVK

## Installation

```bash
go get github.com/darkace1998/golang-vulkan-api
```


## Debugging and Resource Tracking

The library includes a built-in `LeakTracker` utility to help you monitor Vulkan object allocations and identify potential memory leaks (e.g., calling `CreateBuffer` without a corresponding `DestroyBuffer`).

```go
package main

import (
    "fmt"
    vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
    // Enable tracking before creating objects
    vulkan.EnableLeakTracker()
    defer func() {
        // Report any un-freed resources at exit
        fmt.Println(vulkan.ReportLeaks())
    }()

    // ... your Vulkan code ...
}
```

## Error Handling Patterns

See [ERROR_HANDLING.md](ERROR_HANDLING.md) for idiomatic Go patterns for handling `VulkanError` vs. `ValidationError`, including retry logic for transient failures like `VK_ERROR_DEVICE_LOST`.

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for solutions to common issues related to CGO, package dependencies, Vulkan drivers, and runtime segmentation faults.

## Performance Tuning

See [PERFORMANCE_TUNING.md](PERFORMANCE_TUNING.md) for detailed strategies and tips for optimizing Vulkan compute workloads, particularly for AI/ML and general parallel processing tasks.

## Thread Safety

See [THREAD_SAFETY.md](THREAD_SAFETY.md) for detailed information about thread safety guarantees, host synchronization requirements, and specific details regarding video codec function loading.


## Getting Started

If you are new to Vulkan or this library, check out the **[Getting Started Tutorial](GETTING_STARTED.md)** for a step-by-step guide to writing your first Vulkan application.
## Vulkan 1.4 Readiness

See [VULKAN_1_4_READINESS.md](VULKAN_1_4_READINESS.md) for detailed information about the current state of Vulkan 1.4 support and our roadmap for future implementation.

## Additional Documentation

The repository includes comprehensive documentation files to help you better understand and utilize the API:

- [API_REFERENCE.md](API_REFERENCE.md): Comprehensive API mappings and documentation for core features.
- [SECURITY.md](SECURITY.md): Security analysis and posture.
- [MULTIPLATFORM.md](MULTIPLATFORM.md): Cross-platform build support details.
- [todo.md](todo.md): Roadmap for planned features and API coverage tracking.

## Architecture

See [ARCHITECTURE_DIAGRAMS.md](ARCHITECTURE_DIAGRAMS.md) for visual representations of the extension loading mechanism, error handling paths, and thread safety models used by this library.

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
    // Create Vulkan instance
    instanceCreateInfo := &vulkan.InstanceCreateInfo{
        ApplicationInfo: &vulkan.ApplicationInfo{
            ApplicationName:    "My Vulkan App",
            ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
            EngineName:         "My Engine",
            EngineVersion:      vulkan.MakeVersion(1, 0, 0),
            APIVersion:         vulkan.Version13,
        },
    }

    instance, err := vulkan.CreateInstance(instanceCreateInfo)
    if err != nil {
        log.Fatal("Failed to create Vulkan instance:", err)
    }
    defer vulkan.DestroyInstance(instance)

    // Enumerate physical devices
    physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
    if err != nil {
        log.Fatal("Failed to enumerate physical devices:", err)
    }

    fmt.Printf("Found %d physical device(s)\n", len(physicalDevices))
    
    // Get device properties
    for i, device := range physicalDevices {
        props := vulkan.GetPhysicalDeviceProperties(device)
        fmt.Printf("Device %d: %s\n", i, props.DeviceName)
    }
}
```

## Core Components

### Instance Management
- Create and destroy Vulkan instances
- Enumerate extensions and layers
- Physical device enumeration

### Device Management
- Physical device properties and features
- Logical device creation
- Queue family management

### Memory Management
- Buffer and image creation
- Memory allocation and binding
- Memory type selection utilities

### Command Buffers
- Command pool management
- Command buffer allocation and recording
- Queue submission and synchronization

### Compute Pipeline Usage
- Compute pipeline helpers
- Storage buffer management
- Dispatch commands for parallel processing
- Pipeline barriers for compute synchronization

### Synchronization
- Semaphores for GPU-GPU synchronization
- Fences for CPU-GPU synchronization
- Pipeline barriers and memory barriers

## Examples

See the `examples/` directory for example programs:

- **`basic`**: Basic Vulkan setup and physical device enumeration.
- **`benchmark`**: A GPU stress testing and benchmarking tool. See `examples/benchmark/README.md` for more details.
- **`compute`**: Demonstrates how to run a compute shader, including buffer creation and memory binding.
- **`descriptor_manager`**: Example demonstrating the usage of the high-level `DescriptorPoolManager` to easily allocate descriptor sets.
- **`descriptor_update`**: Demonstrates how to bind uniform buffers and combined image samplers by updating descriptor sets.
- **`graphics_pipeline`**: A comprehensive graphics pipeline example showing offscreen rendering with vertex buffers, shader modules, render pass, framebuffer, graphics pipeline creation, and draw commands.
- **`multi_queue`**: Shows how to discover, create, and use multiple Vulkan queues, specifically focusing on parallel transfer and graphics operations.
- **`pipeline_cache`**: Example demonstrating pipeline cache creation, retrieval, merging, and loading data.
- **`push_constants`**: Demonstrates how to use push constants to pass small amounts of data to shaders efficiently.
- **`render_to_texture`**: Demonstrates framebuffer creation, render pass, and reading back pixels without a window/surface.
- **`secondary_command_buffer`**: Shows how to record and execute secondary command buffers.
- **`simple`**: Minimal example for Vulkan instance creation.
- **`subpass_dependencies`**: Demonstrates creating a multi-subpass render pass with a dependency chain and self-dependency.
- **`swapchain`**: Demonstrates all swapchain types, constants, input validation, synchronization objects, and the full present loop workflow.
- **`type`**: Type system and constant validation example.
- **`video`**: Demonstrates video codec support detection.
- **`vulkan13`**: Vulkan 1.3 feature demonstration, including dynamic state and rendering info.

See [examples/benchmark/README.md](examples/benchmark/README.md) for detailed information about the GPU benchmark tool.

## Testing

The implementation includes comprehensive tests and build checks:

```bash
# Build everything
go build ./...

# Run tests
go test ./...

# Run the race detector
go test -race ./...

# Run examples
go run ./examples/basic
go run ./examples/compute
go run ./examples/video
go run ./examples/benchmark -help
```

## API Reference

### Version Management

```go
// Create version numbers
version := vulkan.MakeVersion(1, 3, 0)
major := version.Major()    // 1
minor := version.Minor()    // 3
patch := version.Patch()    // 0

// Predefined versions
vulkan.Version10  // Vulkan 1.0
vulkan.Version11  // Vulkan 1.1
vulkan.Version12  // Vulkan 1.2
vulkan.Version13  // Vulkan 1.3
vulkan.Version14  // Vulkan 1.4 (when available)
```

### Error Handling

```go
result := vulkan.SomeFunction()
if result != vulkan.Success {
    fmt.Printf("Error: %s\n", result.Error())
}

// Or for functions that return (value, error)
value, err := vulkan.SomeOtherFunction()
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

### Instance Creation

```go
instance, err := vulkan.CreateInstance(&vulkan.InstanceCreateInfo{
    ApplicationInfo: &vulkan.ApplicationInfo{
        ApplicationName:    "My App",
        ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
        EngineName:         "My Engine", 
        EngineVersion:      vulkan.MakeVersion(1, 0, 0),
        APIVersion:         vulkan.Version13,
    },
    EnabledLayerNames:     []string{"VK_LAYER_KHRONOS_validation"},
    EnabledExtensionNames: []string{"VK_EXT_debug_utils"},
})
```

### Device Creation

```go
device, err := vulkan.CreateDevice(physicalDevice, &vulkan.DeviceCreateInfo{
    QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
        {
            QueueFamilyIndex: graphicsQueueFamily,
            QueuePriorities:  []float32{1.0},
        },
    },
    EnabledExtensionNames: []string{"VK_KHR_swapchain"},
    EnabledFeatures:       &features,
})
```

### Buffer Management

```go
// Create buffer
buffer, err := vulkan.CreateBuffer(device, &vulkan.BufferCreateInfo{
    Size:        1024,
    Usage:       vulkan.BufferUsageVertexBufferBit,
    SharingMode: vulkan.SharingModeExclusive,
})

// Get memory requirements
memReqs := vulkan.GetBufferMemoryRequirements(device, buffer)

// Allocate and bind memory
memory, err := vulkan.AllocateMemory(device, &vulkan.MemoryAllocateInfo{
    AllocationSize:  memReqs.Size,
    MemoryTypeIndex: suitableMemoryType,
})

err = vulkan.BindBufferMemory(device, buffer, memory, 0)
```

### Compute Pipeline Example

```go
// Create compute shader module (from compiled SPIR-V bytecode)
shaderModule, err := vulkan.CreateShaderModule(device, &vulkan.ShaderModuleCreateInfo{
    CodeSize: uint32(len(shaderCode) * 4),
    Code:     shaderCode, // SPIR-V bytecode
})

// Create descriptor set layout for storage buffers
descriptorSetLayout, err := vulkan.CreateDescriptorSetLayout(device, &vulkan.DescriptorSetLayoutCreateInfo{
    Bindings: []vulkan.DescriptorSetLayoutBinding{
        {
            Binding:         0,
            DescriptorType:  vulkan.DescriptorTypeStorageBuffer,
            DescriptorCount: 1,
            StageFlags:      vulkan.ShaderStageComputeBit,
        },
    },
})

// Create compute pipeline
computePipelines, err := vulkan.CreateComputePipelines(device, nil, []vulkan.ComputePipelineCreateInfo{
    {
        Stage: vulkan.PipelineShaderStageCreateInfo{
            Stage:  vulkan.ShaderStageComputeBit,
            Module: shaderModule,
            Name:   "main",
        },
        Layout: pipelineLayout,
    },
})

// Record and dispatch compute work
vulkan.CmdBindPipeline(commandBuffer, vulkan.PipelineBindPointCompute, computePipelines[0])
vulkan.CmdDispatch(commandBuffer, workGroupsX, workGroupsY, workGroupsZ)
```

## Building

The library uses CGO to interface with the Vulkan C API and is designed to work across multiple platforms. Make sure you have:

1. CGO enabled (`CGO_ENABLED=1`)
2. Vulkan development libraries installed
3. A supported Go compiler (Go 1.22+)

```bash
# Build the repository
go build ./...

# Run tests
go test ./...

# Run the race detector
go test -race ./...

# Run an example
go run ./examples/basic
```

The library automatically configures build settings for your platform using Go build tags.

## Platform-Specific Setup

### Linux
```bash
# Install Vulkan development libraries
sudo apt-get install libvulkan-dev pkg-config

# Or on other distributions
sudo yum install vulkan-devel pkgconf-pkg-config
sudo pacman -S vulkan-headers vulkan-validation-layers pkg-config
```

### Windows
1. Install the Vulkan SDK from [LunarG](https://vulkan.lunarg.com/)
2. Make sure the SDK is in your PATH
3. Ensure Vulkan libraries are available:
   ```cmd
   # The library will automatically link vulkan-1.lib
   # No additional configuration needed if SDK is installed properly
   ```

### macOS
1. Install Vulkan SDK with MoltenVK support from [LunarG](https://vulkan.lunarg.com/)
2. Install pkg-config if not available:
   ```bash
   brew install pkg-config
   ```
3. Vulkan runs on top of Metal via MoltenVK translation layer

### Other Unix Systems
Other Unix-like systems may work if pkg-config and Vulkan development libraries are available.

## Contributing

Contributions are welcome! Please feel free to submit pull requests, report bugs, or suggest features. See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on building, testing, and submitting code.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Based on the official Vulkan specification
- Inspired by other Vulkan bindings in the Go ecosystem
- Thanks to the Vulkan community for excellent documentation
