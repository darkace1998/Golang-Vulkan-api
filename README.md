# Golang-Vulkan-api

A comprehensive Go binding for the Vulkan 1.3+ graphics and compute API with **complete Vulkan 1.3 feature support**.

## Overview

This library provides a complete, type-safe Go interface to the Vulkan API, supporting Vulkan versions 1.0 through 1.3. It's designed to be used as a library for other Go projects that need low-level graphics and compute functionality.

## Features

- ✅ **Complete Vulkan 1.3 Support**: All essential Vulkan 1.3 functions and types
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
- ✅ **Buffer/Image Operations**: Complete buffer and image management
- ✅ **Queue Operations**: Graphics, compute, and transfer queue support
- ✅ **Compute Shaders**: Complete compute pipeline support for AI/ML workloads
- ✅ **Storage Buffers**: Large dataset handling for compute operations
- ✅ **Dispatch Commands**: Efficient compute work group dispatching
- ✅ **Cross-Platform**: Works on Linux, Windows, and macOS (where Vulkan is supported)

## Video Codec Support 🎬

### Fully Supported (Encode & Decode)

These codecs have ratified extensions for both operations:

- **H.264 (AVC)** - VK_KHR_video_encode_h264 & VK_KHR_video_decode_h264
- **H.265 (HEVC)** - VK_KHR_video_encode_h265 & VK_KHR_video_decode_h265
- **AV1** - VK_KHR_video_encode_av1 & VK_KHR_video_decode_av1

Hardware-accelerated video encoding and decoding is supported through Vulkan Video extensions, enabling efficient media processing on compatible GPUs.

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

See `examples/video_codec_check.go` for a complete example that detects and displays supported video codecs on your system.

**Note**: Actual hardware support depends on your GPU model and driver version. Extension availability does not guarantee hardware acceleration.

## Vulkan 1.3 Features ⭐ NEW

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

- Go 1.19 or later
- CGO enabled
- Vulkan SDK or development libraries installed
  - Linux: `libvulkan-dev` package
  - Windows: Vulkan SDK from LunarG
  - macOS: Vulkan SDK with MoltenVK

## Installation

```bash
go get github.com/darkace1998/Golang-Vulkan-api
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    vulkan "github.com/darkace1998/Golang-Vulkan-api"
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

### Compute Pipeline
- Compute shader support for AI/ML workloads
- Storage buffer management for large datasets
- Dispatch commands for parallel processing
- Pipeline barriers for compute synchronization

### Synchronization
- Semaphores for GPU-GPU synchronization
- Fences for CPU-GPU synchronization
- Pipeline barriers and memory barriers

## Examples

See the `examples/` directory for complete working examples:

- `basic_example.go`: Basic Vulkan setup and device enumeration
- `compute_example.go`: **Compute shader example for AI/ML workloads**
- `vulkan13_test.go`: **Complete Vulkan 1.3 feature test and demonstration**
- `video_codec_check.go`: **Video codec support detection tool** 🎬
- `type_example.go`: Type system and constant validation
- `simple_example.go`: Minimal Vulkan instance creation
- `graphics_benchmark.go`: **GPU stress testing and benchmarking tool**

See [examples/BENCHMARK_README.md](examples/BENCHMARK_README.md) for detailed information about the GPU benchmark tool.

## Testing

The implementation includes comprehensive tests for all Vulkan 1.3 features:

```bash
# Test all Vulkan 1.3 features
go run examples/vulkan13_test.go

# Test basic functionality
go run examples/basic_example.go

# Test compute capabilities for AI workloads
go run examples/compute_example.go

# Check video codec support on your GPU
go run examples/video_codec_check.go
```

The Vulkan 1.3 test validates:
- ✅ Dynamic Rendering APIs
- ✅ Synchronization2 enhanced submission
- ✅ Extended Dynamic State commands
- ✅ Private Data management
- ✅ Maintenance4 memory requirements
- ✅ All new pipeline stage flags and constants

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
vulkan.Version14  // Vulkan 1.4
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

### Compute Pipeline for AI/ML Workloads

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
3. A supported Go compiler (Go 1.19+)

```bash
# Build the library
go build

# Run tests
go test

# Build examples
cd examples
go build basic_example.go
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
The library supports other Unix-like systems (FreeBSD, OpenBSD, etc.) that have pkg-config and Vulkan development libraries available.

## Contributing

Contributions are welcome! Please feel free to submit pull requests, report bugs, or suggest features.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Based on the official Vulkan specification
- Inspired by other Vulkan bindings in the Go ecosystem
- Thanks to the Vulkan community for excellent documentation
