# Vulkan Go API Reference

This document provides a comprehensive reference for all available functions in the Vulkan Go binding, including complete Vulkan 1.3 support.

## Table of Contents

- [Core Types](#core-types)
  - [Version Management](#version-management)
  - [Error Handling](#error-handling)
  - [Boolean Conversion](#boolean-conversion)
- [Instance Management](#instance-management)
- [Device Management](#device-management)
- [Memory Management](#memory-management)
- [Command Buffer Management](#command-buffer-management)
- [Synchronization](#synchronization)
- [Vulkan 1.3 Features ⭐ NEW](#vulkan-13-features--new)
- [Pipeline Management](#pipeline-management)
- [Descriptor Management](#descriptor-management)
- [Command Recording](#command-recording)
- [Compute Pipeline Management](#compute-pipeline-management)
- [Video Codec Support 🎬 NEW](#video-codec-support--new)
- [Utility Functions](#utility-functions)
- [Constants and Enums](#constants-and-enums)
- [Important Constants](#important-constants)
- [Notes](#notes)

## Core Types

### Version Management
- `MakeVersion(major, minor, patch uint32) Version` - Create version number
- `(v Version) Major() uint32` - Extract major version
- `(v Version) Minor() uint32` - Extract minor version  
- `(v Version) Patch() uint32` - Extract patch version

### Error Handling
- `(r Result) Error() string` - Get error message
- `(r Result) IsError() bool` - Check if result is error
- `(r Result) IsSuccess() bool` - Check if result is success

### Boolean Conversion
- `FromBool(b bool) Bool32` - Convert Go bool to Vulkan Bool32
- `(b Bool32) ToBool() bool` - Convert Vulkan Bool32 to Go bool

## Instance Management

### Instance Creation/Destruction
- `CreateInstance(createInfo *InstanceCreateInfo) (Instance, error)` - Create Vulkan instance
- `DestroyInstance(instance Instance)` - Destroy Vulkan instance

### Extension/Layer Enumeration
- `EnumerateInstanceExtensionProperties(layerName string) ([]ExtensionProperties, error)` - List instance extensions
- `EnumerateInstanceLayerProperties() ([]LayerProperties, error)` - List instance layers

### Physical Device Management
- `EnumeratePhysicalDevices(instance Instance) ([]PhysicalDevice, error)` - List physical devices
- `GetPhysicalDeviceProperties(physicalDevice PhysicalDevice) PhysicalDeviceProperties` - Get device properties
- `GetPhysicalDeviceFeatures(physicalDevice PhysicalDevice) PhysicalDeviceFeatures` - Get device features
- `GetPhysicalDeviceMemoryProperties(physicalDevice PhysicalDevice) PhysicalDeviceMemoryProperties` - Get memory properties
- `GetPhysicalDeviceQueueFamilyProperties(physicalDevice PhysicalDevice) []QueueFamilyProperties` - Get queue families
- `EnumerateDeviceExtensionProperties(physicalDevice PhysicalDevice, layerName string) ([]ExtensionProperties, error)` - List device extensions

## Device Management

### Device Creation/Destruction
- `CreateDevice(physicalDevice PhysicalDevice, createInfo *DeviceCreateInfo) (Device, error)` - Create logical device
- `DestroyDevice(device Device)` - Destroy logical device

### Queue Management
- `GetDeviceQueue(device Device, queueFamilyIndex, queueIndex uint32) Queue` - Get device queue
- `QueueWaitIdle(queue Queue) error` - Wait for queue to become idle
- `DeviceWaitIdle(device Device) error` - Wait for device to become idle

## Memory Management

### Buffer Operations
- `CreateBuffer(device Device, createInfo *BufferCreateInfo) (Buffer, error)` - Create buffer
- `DestroyBuffer(device Device, buffer Buffer)` - Destroy buffer
- `GetBufferMemoryRequirements(device Device, buffer Buffer) MemoryRequirements` - Get buffer memory requirements
- `BindBufferMemory(device Device, buffer Buffer, memory DeviceMemory, memoryOffset DeviceSize) error` - Bind buffer memory

### Image Operations
- `CreateImage(device Device, createInfo *ImageCreateInfo) (Image, error)` - Create image
- `DestroyImage(device Device, image Image)` - Destroy image
- `GetImageMemoryRequirements(device Device, image Image) MemoryRequirements` - Get image memory requirements
- `BindImageMemory(device Device, image Image, memory DeviceMemory, memoryOffset DeviceSize) error` - Bind image memory

### Memory Allocation
- `AllocateMemory(device Device, allocateInfo *MemoryAllocateInfo) (DeviceMemory, error)` - Allocate device memory
- `FreeMemory(device Device, memory DeviceMemory)` - Free device memory
- `MapMemory(device Device, memory DeviceMemory, offset, size DeviceSize, flags uint32) (unsafe.Pointer, error)` - Map memory
- `UnmapMemory(device Device, memory DeviceMemory)` - Unmap memory

### Utility Functions
- `FindMemoryType(memProperties PhysicalDeviceMemoryProperties, typeFilter uint32, properties MemoryPropertyFlags) (uint32, bool)` - Find suitable memory type

## Command Buffer Management

### Command Pool Operations
- `CreateCommandPool(device Device, createInfo *CommandPoolCreateInfo) (CommandPool, error)` - Create command pool
- `DestroyCommandPool(device Device, commandPool CommandPool)` - Destroy command pool

### Command Buffer Operations
- `AllocateCommandBuffers(device Device, allocateInfo *CommandBufferAllocateInfo) ([]CommandBuffer, error)` - Allocate command buffers
- `FreeCommandBuffers(device Device, commandPool CommandPool, commandBuffers []CommandBuffer)` - Free command buffers
- `BeginCommandBuffer(commandBuffer CommandBuffer, beginInfo *CommandBufferBeginInfo) error` - Begin recording
- `EndCommandBuffer(commandBuffer CommandBuffer) error` - End recording

### Queue Submission
- `QueueSubmit(queue Queue, submitInfos []SubmitInfo, fence Fence) error` - Submit command buffers to queue

## Synchronization

### Semaphore Operations
- `CreateSemaphore(device Device, createInfo *SemaphoreCreateInfo) (Semaphore, error)` - Create semaphore
- `DestroySemaphore(device Device, semaphore Semaphore)` - Destroy semaphore

### Fence Operations
- `CreateFence(device Device, createInfo *FenceCreateInfo) (Fence, error)` - Create fence
- `DestroyFence(device Device, fence Fence)` - Destroy fence
- `WaitForFences(device Device, fences []Fence, waitAll bool, timeout uint64) error` - Wait for fences
- `ResetFences(device Device, fences []Fence) error` - Reset fences
- `GetFenceStatus(device Device, fence Fence) Result` - Get fence status

## Vulkan 1.3 Features ⭐ NEW

### Dynamic Rendering
- `CmdBeginRendering(commandBuffer CommandBuffer, renderingInfo *RenderingInfo)` - Begin dynamic render pass
- `CmdEndRendering(commandBuffer CommandBuffer)` - End dynamic render pass

### Synchronization2 (Enhanced)
- `QueueSubmit2(queue Queue, submitInfos []SubmitInfo2, fence Fence) error` - Enhanced queue submission with timeline semantics

### Extended Dynamic State
- `CmdSetCullMode(commandBuffer CommandBuffer, cullMode CullModeFlags)` - Set cull mode dynamically
- `CmdSetFrontFace(commandBuffer CommandBuffer, frontFace FrontFace)` - Set front face orientation dynamically
- `CmdSetPrimitiveTopology(commandBuffer CommandBuffer, primitiveTopology PrimitiveTopology)` - Set primitive topology dynamically
- `CmdSetViewportWithCount(commandBuffer CommandBuffer, viewports []Viewport)` - Set viewports with count dynamically
- `CmdSetScissorWithCount(commandBuffer CommandBuffer, scissors []Rect2D)` - Set scissor rectangles with count dynamically
- `CmdBindVertexBuffers2(commandBuffer CommandBuffer, firstBinding uint32, buffers []Buffer, offsets []DeviceSize, sizes []DeviceSize, strides []DeviceSize)` - Bind vertex buffers with extended parameters
- `CmdSetDepthTestEnable(commandBuffer CommandBuffer, depthTestEnable bool)` - Set depth test enable state dynamically
- `CmdSetDepthWriteEnable(commandBuffer CommandBuffer, depthWriteEnable bool)` - Set depth write enable state dynamically
- `CmdSetDepthCompareOp(commandBuffer CommandBuffer, depthCompareOp CompareOp)` - Set depth compare operation dynamically
- `CmdSetDepthBoundsTestEnable(commandBuffer CommandBuffer, depthBoundsTestEnable bool)` - Set depth bounds test enable state dynamically
- `CmdSetStencilTestEnable(commandBuffer CommandBuffer, stencilTestEnable bool)` - Set stencil test enable state dynamically
- `CmdSetStencilOp(commandBuffer CommandBuffer, faceMask StencilFaceFlags, failOp, passOp, depthFailOp StencilOp, compareOp CompareOp)` - Set stencil operation dynamically

### Private Data
- `CreatePrivateDataSlot(device Device, createInfo *PrivateDataSlotCreateInfo) (PrivateDataSlot, error)` - Create private data slot
- `DestroyPrivateDataSlot(device Device, privateDataSlot PrivateDataSlot)` - Destroy private data slot
- `SetPrivateData(device Device, objectType ObjectType, objectHandle uint64, privateDataSlot PrivateDataSlot, data uint64) error` - Associate data with Vulkan object
- `GetPrivateData(device Device, objectType ObjectType, objectHandle uint64, privateDataSlot PrivateDataSlot) uint64` - Retrieve data associated with Vulkan object

### Maintenance4
- `GetDeviceBufferMemoryRequirements(device Device, bufferCreateInfo *BufferCreateInfo) MemoryRequirements` - Get buffer memory requirements without creating buffer
- `GetDeviceImageMemoryRequirements(device Device, imageCreateInfo *ImageCreateInfo) MemoryRequirements` - Get image memory requirements without creating image

## Pipeline Management

### Shader Modules
- `CreateShaderModule(device Device, createInfo *ShaderModuleCreateInfo) (ShaderModule, error)` - Create shader module
- `DestroyShaderModule(device Device, shaderModule ShaderModule)` - Destroy shader module

### Pipeline Layouts
- `CreatePipelineLayout(device Device, createInfo *PipelineLayoutCreateInfo) (PipelineLayout, error)` - Create pipeline layout
- `DestroyPipelineLayout(device Device, pipelineLayout PipelineLayout)` - Destroy pipeline layout

### Render Passes
- `CreateRenderPass(device Device, createInfo *RenderPassCreateInfo) (RenderPass, error)` - Create render pass
- `DestroyRenderPass(device Device, renderPass RenderPass)` - Destroy render pass

## Descriptor Management

### Image Views
- `CreateImageView(device Device, createInfo *ImageViewCreateInfo) (ImageView, error)` - Create image view
- `DestroyImageView(device Device, imageView ImageView)` - Destroy image view

### Samplers
- `CreateSampler(device Device, createInfo *SamplerCreateInfo) (Sampler, error)` - Create sampler
- `DestroySampler(device Device, sampler Sampler)` - Destroy sampler

### Descriptor Set Layouts
- `CreateDescriptorSetLayout(device Device, createInfo *DescriptorSetLayoutCreateInfo) (DescriptorSetLayout, error)` - Create descriptor set layout
- `DestroyDescriptorSetLayout(device Device, layout DescriptorSetLayout)` - Destroy descriptor set layout

### Descriptor Pools
- `CreateDescriptorPool(device Device, createInfo *DescriptorPoolCreateInfo) (DescriptorPool, error)` - Create descriptor pool
- `DestroyDescriptorPool(device Device, pool DescriptorPool)` - Destroy descriptor pool

## Command Recording

### Render Pass Commands
- `CmdBeginRenderPass(commandBuffer CommandBuffer, beginInfo *RenderPassBeginInfo, contents SubpassContents)` - Begin render pass
- `CmdEndRenderPass(commandBuffer CommandBuffer)` - End render pass

### Pipeline Commands
- `CmdBindPipeline(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, pipeline Pipeline)` - Bind pipeline

### Compute Commands
- `CmdDispatch(commandBuffer CommandBuffer, groupCountX, groupCountY, groupCountZ uint32)` - Dispatch compute work groups
- `CmdDispatchIndirect(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize)` - Dispatch compute work with parameters from buffer
- `CmdBindDescriptorSets(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, layout PipelineLayout, firstSet uint32, descriptorSets []DescriptorSet, dynamicOffsets []uint32)` - Bind descriptor sets

### State Commands
- `CmdSetViewport(commandBuffer CommandBuffer, firstViewport uint32, viewports []Viewport)` - Set viewport
- `CmdSetScissor(commandBuffer CommandBuffer, firstScissor uint32, scissors []Rect2D)` - Set scissor

### Buffer Binding Commands
- `CmdBindVertexBuffers(commandBuffer CommandBuffer, firstBinding uint32, buffers []Buffer, offsets []DeviceSize)` - Bind vertex buffers
- `CmdBindIndexBuffer(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, indexType IndexType)` - Bind index buffer

### Drawing Commands
- `CmdDraw(commandBuffer CommandBuffer, vertexCount, instanceCount, firstVertex, firstInstance uint32)` - Draw primitives
- `CmdDrawIndexed(commandBuffer CommandBuffer, indexCount, instanceCount, firstIndex uint32, vertexOffset int32, firstInstance uint32)` - Draw indexed

### Transfer Commands
- `CmdCopyBuffer(commandBuffer CommandBuffer, srcBuffer, dstBuffer Buffer, regions []BufferCopy)` - Copy buffer data

### Synchronization Commands
- `CmdPipelineBarrier(commandBuffer CommandBuffer, srcStageMask, dstStageMask PipelineStageFlags, dependencyFlags uint32)` - Insert pipeline barrier

## Compute Pipeline Management

### Compute Pipeline Creation
- `CreateComputePipelines(device Device, pipelineCache PipelineCache, createInfos []ComputePipelineCreateInfo) ([]Pipeline, error)` - Create compute pipelines
- `DestroyPipeline(device Device, pipeline Pipeline)` - Destroy pipeline (graphics or compute)

## Video Codec Support 🎬 NEW

### Video Codec Extensions

Supported video codec extensions:
- **H.264 (AVC)**: `VK_KHR_video_encode_h264` & `VK_KHR_video_decode_h264`
- **H.265 (HEVC)**: `VK_KHR_video_encode_h265` & `VK_KHR_video_decode_h265`
- **AV1**: `VK_KHR_video_encode_av1` & `VK_KHR_video_decode_av1`

### Video Codec Functions

#### Capability Queries
- `GetSupportedVideoCodecs(physicalDevice PhysicalDevice) ([]string, error)` - Get list of supported video codecs on the device
- `GetVideoCapabilities(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo) (*VideoCapabilities, error)` - Get video codec capabilities
- `IsVideoCodecSupported(extensionName string, availableExtensions []ExtensionProperties) bool` - Check if specific video codec is supported

#### Video Session Management
- `CreateVideoSession(device Device, createInfo *VideoSessionCreateInfo) (VideoSession, error)` - Create video session for encoding/decoding
- `DestroyVideoSession(device Device, videoSession VideoSession)` - Destroy video session
- `CreateVideoSessionParameters(device Device, createInfo *VideoSessionParametersCreateInfo) (VideoSessionParameters, error)` - Create video session parameters
- `DestroyVideoSessionParameters(device Device, videoSessionParameters VideoSessionParameters)` - Destroy video session parameters

#### Video Coding Commands
- `CmdBeginVideoCoding(commandBuffer CommandBuffer, videoSession VideoSession, videoSessionParameters VideoSessionParameters)` - Begin video coding operations
- `CmdEndVideoCoding(commandBuffer CommandBuffer)` - End video coding operations
- `CmdDecodeVideo(commandBuffer CommandBuffer, decodeInfo *VideoDecodeInfo)` - Perform video decode operation
- `CmdEncodeVideo(commandBuffer CommandBuffer, encodeInfo *VideoEncodeInfo)` - Perform video encode operation

### Video Types and Constants

#### Video Codec Operations
- `VideoCodecOperationDecodeH264Bit` - H.264 decode operation
- `VideoCodecOperationDecodeH265Bit` - H.265 decode operation
- `VideoCodecOperationDecodeAV1Bit` - AV1 decode operation
- `VideoCodecOperationEncodeH264Bit` - H.264 encode operation
- `VideoCodecOperationEncodeH265Bit` - H.265 encode operation
- `VideoCodecOperationEncodeAV1Bit` - AV1 encode operation

#### Chroma Subsampling
- `VideoChromaSubsamplingMonochrome` - Monochrome (no chroma)
- `VideoChromaSubsampling420` - 4:2:0 subsampling
- `VideoChromaSubsampling422` - 4:2:2 subsampling
- `VideoChromaSubsampling444` - 4:4:4 subsampling

#### Component Bit Depths
- `VideoComponentBitDepth8` - 8-bit component depth
- `VideoComponentBitDepth10` - 10-bit component depth
- `VideoComponentBitDepth12` - 12-bit component depth

### Example Usage

```go
// Check supported video codecs
supportedCodecs, err := vulkan.GetSupportedVideoCodecs(physicalDevice)
if err != nil {
    log.Fatal(err)
}

for _, codec := range supportedCodecs {
    fmt.Printf("Supported codec: %s\n", codec)
}

// Check if H.264 decode is available
extensions, _ := vulkan.EnumerateDeviceExtensionProperties(physicalDevice, "")
if vulkan.IsVideoCodecSupported(vulkan.ExtensionNameVideoDecodeH264, extensions) {
    fmt.Println("H.264 hardware decode is supported")
}
```

**Note**: Full video codec functionality requires the Vulkan Video extensions to be enabled on the device and supported by the GPU driver. Hardware support varies by GPU model and driver version.

## Utility Functions

### Version and Feature Queries
- `GetAPIVersion() Version` - Get supported API version
- `IsExtensionSupported(extensionName string, availableExtensions []ExtensionProperties) bool` - Check extension support
- `IsLayerSupported(layerName string, availableLayers []LayerProperties) bool` - Check layer support

## Constants and Enums

### API Versions
- `Version10`, `Version11`, `Version12`, `Version13`, `Version14` - Predefined API versions

### Result Codes
- `Success`, `NotReady`, `Timeout`, `EventSet`, `EventReset`, `Incomplete`
- Various error codes: `ErrorOutOfHostMemory`, `ErrorOutOfDeviceMemory`, etc.

### Boolean Values
- `True`, `False` - Vulkan boolean constants

### Queue Flags
- `QueueGraphicsBit`, `QueueComputeBit`, `QueueTransferBit`, `QueueSparseBindingBit`

### Pipeline Bind Points
- `PipelineBindPointGraphics`, `PipelineBindPointCompute`

### Shader Stages
- `ShaderStageVertexBit`, `ShaderStageFragmentBit`, `ShaderStageComputeBit`
- `ShaderStageTessellationControlBit`, `ShaderStageTessellationEvaluationBit`
- `ShaderStageGeometryBit`, `ShaderStageAllGraphics`, `ShaderStageAll`

### Buffer Usage Flags
- `BufferUsageTransferSrcBit`, `BufferUsageTransferDstBit`
- `BufferUsageUniformBufferBit`, `BufferUsageStorageBufferBit`
- `BufferUsageVertexBufferBit`, `BufferUsageIndexBufferBit`

### Memory Property Flags
- `MemoryPropertyDeviceLocalBit`, `MemoryPropertyHostVisibleBit`
- `MemoryPropertyHostCoherentBit`, `MemoryPropertyHostCachedBit`

### Image Usage Flags
- `ImageUsageTransferSrcBit`, `ImageUsageTransferDstBit`
- `ImageUsageSampledBit`, `ImageUsageStorageBit`
- `ImageUsageColorAttachmentBit`, `ImageUsageDepthStencilAttachmentBit`

### Formats
- `FormatUndefined`, `FormatR8G8B8A8Unorm`, `FormatB8G8R8A8Unorm`
- `FormatD16Unorm`, `FormatD32Sfloat`, `FormatD24UnormS8Uint`

### Sample Counts
- `SampleCount1Bit`, `SampleCount2Bit`, `SampleCount4Bit`, `SampleCount8Bit`

### Pipeline Stages
- `PipelineStageTopOfPipeBit`, `PipelineStageBottomOfPipeBit`
- `PipelineStageVertexShaderBit`, `PipelineStageFragmentShaderBit`
- `PipelineStageComputeShaderBit`, `PipelineStageTransferBit`
- `PipelineStageColorAttachmentOutputBit`

### Descriptor Types
- `DescriptorTypeSampler`, `DescriptorTypeCombinedImageSampler`
- `DescriptorTypeUniformBuffer`, `DescriptorTypeStorageBuffer`
- `DescriptorTypeUniformBufferDynamic`, `DescriptorTypeStorageBufferDynamic`
- `DescriptorTypeSampledImage`, `DescriptorTypeStorageImage`

### Access Flags
- `AccessShaderReadBit`, `AccessShaderWriteBit`
- `AccessColorAttachmentReadBit`, `AccessColorAttachmentWriteBit`
- `AccessTransferReadBit`, `AccessTransferWriteBit`

## Important Constants
- `MaxMemoryTypes` (32)
- `MaxMemoryHeaps` (16)
- `MaxPhysicalDeviceNameSize` (256)
- `UuidSize` (16)
- `WholeSize` (18446744073709551615)

## Notes

1. All functions follow Go error handling conventions where applicable
2. Memory management is manual - you must destroy what you create
3. The binding is designed to be as close to the C API as possible while remaining idiomatic Go
4. CGO is required and Vulkan development libraries must be installed
5. Some advanced features may require additional implementation
6. This binding supports Vulkan 1.0 through 1.4 (where available on the system)