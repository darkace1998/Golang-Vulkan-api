# Golang-Vulkan-api

A comprehensive Go binding for the Vulkan 1.4 graphics and compute API.

## Overview

This library provides a complete, type-safe Go interface to the Vulkan API, supporting Vulkan versions 1.0 through 1.4. It's designed to be used as a library for other Go projects that need low-level graphics and compute functionality.

## Features

- ✅ **Complete Vulkan 1.4 Support**: All essential Vulkan functions and types
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
- `type_example.go`: Type system and constant validation
- `simple_example.go`: Minimal Vulkan instance creation

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

The library uses CGO to interface with the Vulkan C API. Make sure you have:

1. CGO enabled (`CGO_ENABLED=1`)
2. Vulkan development libraries installed
3. pkg-config available (Linux/macOS)

```bash
# Build the library
go build

# Run tests
go test

# Build examples
cd examples
go build basic_test.go
```

## Platform-Specific Notes

### Linux
```bash
# Install Vulkan development libraries
sudo apt-get install libvulkan-dev

# Or on other distributions
sudo yum install vulkan-devel
sudo pacman -S vulkan-headers vulkan-validation-layers
```

### Windows
- Install the Vulkan SDK from LunarG
- Make sure the SDK is in your PATH
- May need to set CGO_LDFLAGS manually if using custom install location

### macOS
- Install Vulkan SDK with MoltenVK support
- Vulkan runs on top of Metal via MoltenVK translation layer

## Contributing

Contributions are welcome! Please feel free to submit pull requests, report bugs, or suggest features.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Based on the official Vulkan specification
- Inspired by other Vulkan bindings in the Go ecosystem
- Thanks to the Vulkan community for excellent documentation
