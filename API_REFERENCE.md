# Vulkan Go API Reference

This document provides a reference for the exported functions in the Vulkan Go binding. The current verified repo state builds and tests on Linux with `libvulkan-dev` installed.

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
- [Vulkan 1.3 Features](#vulkan-13-features)
- [Query Management](#query-management)
- [Swapchain Management](#swapchain-management)
- [Pipeline Management](#pipeline-management)
- [Descriptor Management](#descriptor-management)
- [Command Recording](#command-recording)
- [Compute Pipeline Management](#compute-pipeline-management)
- [Video Codec Support](#video-codec-support)
- [Utility Functions](#utility-functions)
  - [Debugging and Resource Tracking](#debugging-and-resource-tracking)
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

- `IsErrorDeviceLost(err error) bool` - IsErrorDeviceLost checks if an error indicates that the Vulkan device has been lost
- `IsErrorOutOfDate(err error) bool` - IsErrorOutOfDate checks if an error indicates that the Vulkan swapchain is out of date
- `IsVulkanError(err error) bool` - IsVulkanError checks if an error is a VulkanError
- `NewValidationError(field, reason string) *ValidationError` - NewValidationError creates a new ValidationError
- `NewVulkanError(result Result, operation string, details string) *VulkanError` - NewVulkanError creates a new VulkanError
- `(e *VulkanError) Unwrap() error` - Returns the underlying Result as an error for error unwrapping
### Boolean Conversion
- `FromBool(b bool) Bool32` - Convert Go bool to Vulkan Bool32
- `(b Bool32) ToBool() bool` - Convert Vulkan Bool32 to Go bool

## Instance Management

### Instance Creation/Destruction
- `CreateInstance(createInfo *InstanceCreateInfo) (Instance, error)` - Create Vulkan instance
- `DestroyInstance(instance Instance)` - Destroy Vulkan instance

### Debug Utilities
- `LoadDebugUtilsFunctions(instance Instance)` - Load VK_EXT_debug_utils functions
- `CreateDebugUtilsMessengerEXT(instance Instance, createInfo *DebugUtilsMessengerCreateInfo, callback DebugCallbackFunc) (DebugUtilsMessengerEXT, error)` - Create debug messenger
- `DestroyDebugUtilsMessengerEXT(instance Instance, messenger DebugUtilsMessengerEXT)` - Destroy debug messenger
- `SetDebugUtilsObjectNameEXT(device Device, nameInfo *DebugUtilsObjectNameInfo) error` - Gives a user-friendly name to an object
- `CmdBeginDebugUtilsLabelEXT(commandBuffer CommandBuffer, labelInfo *DebugUtilsLabel)` - Opens a command buffer debug label region
- `CmdEndDebugUtilsLabelEXT(commandBuffer CommandBuffer)` - Closes a command buffer debug label region
- `CmdInsertDebugUtilsLabelEXT(commandBuffer CommandBuffer, labelInfo *DebugUtilsLabel)` - Inserts a single debug label into a command buffer
- `QueueBeginDebugUtilsLabelEXT(queue Queue, labelInfo *DebugUtilsLabel)` - Opens a queue debug label region
- `QueueEndDebugUtilsLabelEXT(queue Queue)` - Closes a queue debug label region
- `QueueInsertDebugUtilsLabelEXT(queue Queue, labelInfo *DebugUtilsLabel)` - Inserts a single debug label into a queue

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

- `EnumeratePhysicalDeviceGroups(instance Instance) ([]PhysicalDeviceGroupProperties, error)` - EnumeratePhysicalDeviceGroups enumerates physical device groups for multi-GPU
- `GetPhysicalDeviceFeatures2(physicalDevice PhysicalDevice) (PhysicalDeviceFeatures, error)` - GetPhysicalDeviceFeatures2 gets extended physical device features (Vulkan 1.1+)
## Device Management

- `CreateGraphicsPipelines(device Device, pipelineCache PipelineCache, createInfos []GraphicsPipelineCreateInfo) ([]Pipeline, error)` - CreateGraphicsPipelines creates graphics pipelines
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

- `CopyDataToStagingBuffer(stagingBuffer *StagingBuffer, data []byte) error` - CopyDataToStagingBuffer copies data to a staging buffer
- `CreateBufferView(device Device, createInfo *BufferViewCreateInfo) (BufferView, error)` - CreateBufferView creates a buffer view
- `CreateStagingBuffer(device Device, physicalDevice PhysicalDevice, size DeviceSize) (*StagingBuffer, error)` - CreateStagingBuffer creates a staging buffer for host-to-device transfers
- `DestroyBufferView(device Device, bufferView BufferView)` - DestroyBufferView destroys a buffer view
- `DestroyStagingBuffer(device Device, stagingBuffer *StagingBuffer)` - DestroyStagingBuffer destroys a staging buffer and frees its memory
### Image Operations
- `CreateImage(device Device, createInfo *ImageCreateInfo) (Image, error)` - Create image
- `DestroyImage(device Device, image Image)` - Destroy image
- `GetImageMemoryRequirements(device Device, image Image) MemoryRequirements` - Get image memory requirements
- `BindImageMemory(device Device, image Image, memory DeviceMemory, memoryOffset DeviceSize) error` - Bind image memory

- `CreateFramebuffer(device Device, createInfo *FramebufferCreateInfo) (Framebuffer, error)` - CreateFramebuffer creates a framebuffer
- `DestroyFramebuffer(device Device, framebuffer Framebuffer)` - DestroyFramebuffer destroys a framebuffer
- `GetPhysicalDeviceFormatProperties(physicalDevice PhysicalDevice, format Format) FormatProperties` - GetPhysicalDeviceFormatProperties returns format properties for a physical device
- `GetPhysicalDeviceImageFormatProperties(physicalDevice PhysicalDevice, format Format, imageType ImageType, tiling ImageTiling, usage ImageUsageFlags, flags ImageCreateFlags) (ImageFormatProperties, error)` - GetPhysicalDeviceImageFormatProperties returns image format properties for a physical device
- `GetPhysicalDeviceSparseImageFormatProperties(physicalDevice PhysicalDevice, format Format, imageType ImageType, samples SampleCountFlags, usage ImageUsageFlags, tiling ImageTiling) []SparseImageFormatProperties` - GetPhysicalDeviceSparseImageFormatProperties returns sparse image format properties
### Memory Allocation
- `AllocateMemory(device Device, allocateInfo *MemoryAllocateInfo) (DeviceMemory, error)` - Allocate device memory
- `FreeMemory(device Device, memory DeviceMemory)` - Free device memory
- `MapMemory(device Device, memory DeviceMemory, offset, size DeviceSize, flags uint32) (unsafe.Pointer, error)` - Map memory
- `UnmapMemory(device Device, memory DeviceMemory)` - Unmap memory
- `GetDeviceMemoryCommitment(device Device, memory DeviceMemory) DeviceSize` - Query current memory commitment

- `CreateMemoryPool(device Device, size DeviceSize, memoryTypeIndex uint32, alignment DeviceSize) (*MemoryPool, error)` - CreateMemoryPool creates a memory pool for efficient sub-allocations
- `FindMemoryTypeForUsage(memProperties PhysicalDeviceMemoryProperties, typeFilter uint32, usage MemoryUsage) (uint32, bool)` - FindMemoryTypeForUsage finds a suitable memory type based on common usage patterns
- `FlushMappedMemoryRanges(device Device, memoryRanges []MappedMemoryRange) error` - FlushMappedMemoryRanges flushes mapped memory ranges to make host writes visible to device
- `GetImageSparseMemoryRequirements(device Device, image Image) []SparseImageMemoryRequirements` - GetImageSparseMemoryRequirements returns sparse memory requirements for an image
- `InvalidateMappedMemoryRanges(device Device, memoryRanges []MappedMemoryRange) error` - InvalidateMappedMemoryRanges invalidates mapped memory ranges to make device writes visible to host
### Utility Functions
- `FindMemoryType(memProperties PhysicalDeviceMemoryProperties, typeFilter uint32, properties MemoryPropertyFlags) (uint32, bool)` - Find suitable memory type

## Command Buffer Management

### Command Pool Operations
- `CreateCommandPool(device Device, createInfo *CommandPoolCreateInfo) (CommandPool, error)` - Create command pool
- `DestroyCommandPool(device Device, commandPool CommandPool)` - Destroy command pool

- `CreateThreadLocalCommandPool(device Device, queueFamilyIndex uint32) (*ThreadLocalCommandPool, error)` - CreateThreadLocalCommandPool creates a thread local command pool
- `(pool *ThreadLocalCommandPool) AllocatePrimaryCommandBuffer() (CommandBuffer, error)` - Allocates a primary command buffer from the pool
- `(pool *ThreadLocalCommandPool) AllocateSecondaryCommandBuffer() (CommandBuffer, error)` - Allocates a secondary command buffer from the pool
- `ResetCommandPool(device Device, commandPool CommandPool, flags CommandPoolResetFlags) error` - ResetCommandPool resets a command pool
- `TrimCommandPool(device Device, commandPool CommandPool)` - TrimCommandPool trims a command pool (Vulkan 1.1+)
### Command Buffer Operations
- `AllocateCommandBuffers(device Device, allocateInfo *CommandBufferAllocateInfo) ([]CommandBuffer, error)` - Allocate command buffers
- `FreeCommandBuffers(device Device, commandPool CommandPool, commandBuffers []CommandBuffer)` - Free command buffers
- `BeginCommandBuffer(commandBuffer CommandBuffer, beginInfo *CommandBufferBeginInfo) error` - Begin recording
- `EndCommandBuffer(commandBuffer CommandBuffer) error` - End recording

- `CmdBeginQuery(commandBuffer CommandBuffer, queryPool QueryPool, query uint32, flags QueryControlFlags)` - CmdBeginQuery begins a query
- `CmdCopyQueryPoolResults(commandBuffer CommandBuffer, queryPool QueryPool, firstQuery, queryCount uint32, dstBuffer Buffer, dstOffset DeviceSize, stride DeviceSize, flags QueryResultFlags)` - CmdCopyQueryPoolResults copies the results of queries in a query pool to a buffer object
- `CmdEndQuery(commandBuffer CommandBuffer, queryPool QueryPool, query uint32)` - CmdEndQuery ends a query
- `CmdResetQueryPool(commandBuffer CommandBuffer, queryPool QueryPool, firstQuery, queryCount uint32)` - CmdResetQueryPool resets a range of queries in a query pool on the GPU
- `CmdWriteTimestamp(commandBuffer CommandBuffer, pipelineStage PipelineStageFlags, queryPool QueryPool, query uint32)` - CmdWriteTimestamp writes a device timestamp into a query object
### Queue Submission
- `QueueSubmit(queue Queue, submitInfos []SubmitInfo, fence Fence) error` - Submit command buffers to queue

## Synchronization

### Semaphore Operations
- `CreateSemaphore(device Device, createInfo *SemaphoreCreateInfo) (Semaphore, error)` - Create semaphore
- `DestroySemaphore(device Device, semaphore Semaphore)` - Destroy semaphore

- `CreateTimelineSemaphore(device Device, initialValue uint64) (Semaphore, error)` - CreateTimelineSemaphore creates a timeline semaphore (Vulkan 1.2+)
- `GetSemaphoreCounterValue(device Device, semaphore Semaphore) (uint64, error)` - GetSemaphoreCounterValue gets the current counter value of a timeline semaphore (Vulkan 1.2+)
- `SignalSemaphore(device Device, signalInfo *SemaphoreSignalInfo) error` - SignalSemaphore signals a timeline semaphore (Vulkan 1.2+)
- `WaitSemaphores(device Device, waitInfo *SemaphoreWaitInfo, timeout uint64) error` - WaitSemaphores waits for timeline semaphores (Vulkan 1.2+)

### Event Operations
- `CmdResetEvent(commandBuffer CommandBuffer, event Event, stageMask PipelineStageFlags)` - CmdResetEvent resets an event object to unsignaled state from the device
- `CmdSetEvent(commandBuffer CommandBuffer, event Event, stageMask PipelineStageFlags)` - CmdSetEvent sets an event object to signaled state from the device
- `CreateEvent(device Device, createInfo *EventCreateInfo) (Event, error)` - CreateEvent creates an event object
- `DestroyEvent(device Device, event Event)` - DestroyEvent destroys an event object
- `GetEventStatus(device Device, event Event) (Result, error)` - GetEventStatus gets the status of an event
- `ResetEvent(device Device, event Event) error` - ResetEvent resets an event to unsignaled state from the host
- `SetEvent(device Device, event Event) error` - SetEvent sets an event to signaled state from the host
### Fence Operations
- `CreateFence(device Device, createInfo *FenceCreateInfo) (Fence, error)` - Create fence
- `DestroyFence(device Device, fence Fence)` - Destroy fence
- `WaitForFences(device Device, fences []Fence, waitAll bool, timeout uint64) error` - Wait for fences
- `ResetFences(device Device, fences []Fence) error` - Reset fences
- `GetFenceStatus(device Device, fence Fence) (Result, error)` - Get fence status

- `QueueBindSparse(queue Queue, bindInfos []BindSparseInfo, fence Fence) error` - QueueBindSparse binds sparse resources on a queue
## Vulkan 1.3 Features

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

## Query Management
- `CreateQueryPool(device Device, createInfo *QueryPoolCreateInfo) (QueryPool, error)` - Create a query pool
- `DestroyQueryPool(device Device, queryPool QueryPool)` - Destroy a query pool
- `GetQueryPoolResults(device Device, queryPool QueryPool, firstQuery, queryCount uint32, dataSize uint64, flags QueryResultFlags) ([]byte, error)` - Retrieve query results
- `GetQueryPoolResultsUint32(device Device, queryPool QueryPool, firstQuery, queryCount uint32, flags QueryResultFlags) ([]uint32, error)` - Retrieve 32-bit query results
- `GetQueryPoolResultsUint64(device Device, queryPool QueryPool, firstQuery, queryCount uint32, flags QueryResultFlags) ([]uint64, error)` - Retrieve 64-bit query results

- `ResetQueryPool(device Device, queryPool QueryPool, firstQuery, queryCount uint32)` - ResetQueryPool resets a range of queries in a query pool on the host (Vulkan 1.2+)
## Swapchain Management
- `CreateSwapchain(device Device, createInfo *SwapchainCreateInfo) (Swapchain, error)` - Create a swapchain
- `DestroySwapchain(device Device, swapchain Swapchain)` - Destroy a swapchain
- `GetSwapchainImages(device Device, swapchain Swapchain) ([]Image, error)` - Get swapchain images
- `AcquireNextImage(device Device, swapchain Swapchain, timeout uint64, semaphore Semaphore, fence Fence) (uint32, bool, error)` - Acquire next presentable image
- `QueuePresent(queue Queue, presentInfo *PresentInfo) (bool, error)` - Queue an image for presentation

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
- `GetRenderAreaGranularity(device Device, renderPass RenderPass) Extent2D` - Get render area granularity

### Pipeline Cache
- `CreatePipelineCache(device Device, createInfo *PipelineCacheCreateInfo) (PipelineCache, error)` - Create pipeline cache
- `DestroyPipelineCache(device Device, pipelineCache PipelineCache)` - Destroy pipeline cache
- `MergePipelineCaches(device Device, dstCache PipelineCache, srcCaches []PipelineCache) error` - Merge multiple pipeline caches
- `GetPipelineCacheData(device Device, pipelineCache PipelineCache) ([]byte, error)` - Retrieve data from pipeline cache

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

- `FreeDescriptorSets(device Device, descriptorPool DescriptorPool, descriptorSets []DescriptorSet) error` - FreeDescriptorSets frees one or more descriptor sets
- `ResetDescriptorPool(device Device, descriptorPool DescriptorPool) error` - ResetDescriptorPool resets a descriptor pool
### Descriptor Sets
- `AllocateDescriptorSets(device Device, allocateInfo *DescriptorSetAllocateInfo) ([]DescriptorSet, error)` - Allocate descriptor sets
- `UpdateDescriptorSets(device Device, writes []WriteDescriptorSet, copies []CopyDescriptorSet)` - Update descriptor sets with write and copy operations

### Descriptor Pool Manager (High-Level)
- `NewDescriptorPoolManager(device Device, maxSetsPerPool uint32, poolSizes []DescriptorPoolSize, flags DescriptorPoolCreateFlags) (*DescriptorPoolManager, error)` - Create a high-level manager to dynamically allocate from multiple pools
- `(m *DescriptorPoolManager) AllocateDescriptorSets(layouts []DescriptorSetLayout) ([]DescriptorSet, error)` - Safely allocate descriptor sets, creating new pools as needed
- `(m *DescriptorPoolManager) Reset() error` - Reset all managed descriptor pools
- `(m *DescriptorPoolManager) Destroy()` - Destroy all managed descriptor pools

## Command Recording
- `LoadMeshShaderFunctions(device Device)` - LoadMeshShaderFunctions loads the device-level mesh shader functions.
- `CmdDrawMeshTasksEXT(commandBuffer CommandBuffer, groupCountX, groupCountY, groupCountZ uint32)` - CmdDrawMeshTasksEXT draws mesh tasks.
- `CmdDrawMeshTasksIndirectEXT(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, drawCount, stride uint32)` - CmdDrawMeshTasksIndirectEXT draws mesh tasks with indirect parameters.
- `CmdDrawMeshTasksIndirectCountEXT(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, countBuffer Buffer, countBufferOffset DeviceSize, maxDrawCount, stride uint32)` - CmdDrawMeshTasksIndirectCountEXT draws mesh tasks with indirect parameters and indirect count.

### Render Pass Commands
- `CmdBeginRenderPass(commandBuffer CommandBuffer, beginInfo *RenderPassBeginInfo, contents SubpassContents)` - Begin render pass
- `CmdNextSubpass(commandBuffer CommandBuffer, contents SubpassContents)` - Advances to the next subpass in a render pass
- `CmdEndRenderPass(commandBuffer CommandBuffer)` - End render pass


### Clear Commands
- `CmdClearAttachments(commandBuffer CommandBuffer, attachments []ClearAttachment, rects []ClearRect)` - Clear attachment regions within a render pass
- `CmdClearColorImage(commandBuffer CommandBuffer, image Image, imageLayout ImageLayout, color *ClearColorValue, ranges []ImageSubresourceRange)` - Clear a color image outside of a render pass
- `CmdClearDepthStencilImage(commandBuffer CommandBuffer, image Image, imageLayout ImageLayout, depthStencil *ClearDepthStencilValue, ranges []ImageSubresourceRange)` - Clear a depth/stencil image outside of a render pass

### Execution Commands
- `CmdExecuteCommands(commandBuffer CommandBuffer, commandBuffers []CommandBuffer)` - Executes secondary command buffers from a primary command buffer

### Pipeline Commands
- `CmdBindPipeline(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, pipeline Pipeline)` - Bind pipeline
- `CmdPushConstants(commandBuffer CommandBuffer, layout PipelineLayout, stageFlags ShaderStageFlags, offset uint32, data []byte)` - Update the values of push constants
- `CmdPushConstantsTyped[T any](commandBuffer CommandBuffer, layout PipelineLayout, stageFlags ShaderStageFlags, offset uint32, value *T)` - Convenience wrapper around CmdPushConstants for common use cases

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
- `CmdDrawIndirect(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, drawCount, stride uint32)` - Draw primitives with indirect parameters
- `CmdDrawIndexedIndirect(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, drawCount, stride uint32)` - Draw indexed primitives with indirect parameters
- `CmdDrawIndirectCount(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, countBuffer Buffer, countBufferOffset DeviceSize, maxDrawCount, stride uint32)` - Draw primitives with indirect parameters and draw count
- `CmdDrawIndexedIndirectCount(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, countBuffer Buffer, countBufferOffset DeviceSize, maxDrawCount, stride uint32)` - Draw indexed primitives with indirect parameters and draw count

### Transfer Commands
- `CmdCopyBuffer(commandBuffer CommandBuffer, srcBuffer, dstBuffer Buffer, regions []BufferCopy)` - Copy buffer data
- `CmdBlitImage(commandBuffer CommandBuffer, srcImage Image, srcImageLayout ImageLayout, dstImage Image, dstImageLayout ImageLayout, regions []ImageBlit, filter Filter)` - Copies regions of an image with potential format conversion and scaling
- `CmdCopyBufferToImage(commandBuffer CommandBuffer, srcBuffer Buffer, dstImage Image, dstImageLayout ImageLayout, regions []BufferImageCopy)` - Copies data from a buffer to an image
- `CmdCopyImage(commandBuffer CommandBuffer, srcImage Image, srcImageLayout ImageLayout, dstImage Image, dstImageLayout ImageLayout, regions []ImageCopy)` - Copies data between images
- `CmdCopyImageToBuffer(commandBuffer CommandBuffer, srcImage Image, srcImageLayout ImageLayout, dstBuffer Buffer, regions []BufferImageCopy)` - Copies data from an image to a buffer
- `CmdFillBuffer(commandBuffer CommandBuffer, dstBuffer Buffer, dstOffset DeviceSize, size DeviceSize, data uint32)` - Fills a buffer with a fixed 32-bit value
- `CmdResolveImage(commandBuffer CommandBuffer, srcImage Image, srcImageLayout ImageLayout, dstImage Image, dstImageLayout ImageLayout, regions []ImageResolve)` - Resolves a multisample image to a non-multisample image
- `CmdUpdateBuffer(commandBuffer CommandBuffer, dstBuffer Buffer, dstOffset DeviceSize, data []byte)` - Updates buffer contents inline from host memory

### Synchronization Commands
- `CmdPipelineBarrier(commandBuffer CommandBuffer, srcStageMask, dstStageMask PipelineStageFlags, dependencyFlags uint32)` - Insert pipeline barrier
- `CmdPipelineBarrierFull(commandBuffer CommandBuffer, srcStageMask, dstStageMask PipelineStageFlags, dependencyFlags uint32, memoryBarriers []MemoryBarrier, bufferMemoryBarriers []BufferMemoryBarrier, imageMemoryBarriers []ImageMemoryBarrier)` - Insert a full pipeline barrier with memory barriers
- `CmdWaitEvents(commandBuffer CommandBuffer, events []Event, srcStageMask, dstStageMask PipelineStageFlags, memoryBarriers []MemoryBarrier, bufferMemoryBarriers []BufferMemoryBarrier, imageMemoryBarriers []ImageMemoryBarrier)` - Waits for one or more events and inserts a set of memory barriers

## Compute Pipeline Management

### Compute Pipeline Creation
- `CreateComputePipelines(device Device, pipelineCache PipelineCache, createInfos []ComputePipelineCreateInfo) ([]Pipeline, error)` - Create compute pipelines
- `DestroyPipeline(device Device, pipeline Pipeline)` - Destroy pipeline (graphics or compute)

## Video Codec Support

### Video Codec Extensions

Supported video codec extensions:
- **H.264 (AVC)**: `VK_KHR_video_encode_h264` & `VK_KHR_video_decode_h264`
- **H.265 (HEVC)**: `VK_KHR_video_encode_h265` & `VK_KHR_video_decode_h265`
- **AV1**: `VK_KHR_video_encode_av1` & `VK_KHR_video_decode_av1`

### Video Codec Functions

- `CmdControlVideoCodingReset(commandBuffer CommandBuffer) error` - CmdControlVideoCodingReset issues a reset control command for video coding
- `CreateVideoDeviceFunctions(device Device) (*VideoDeviceFunctions, error)` - CreateVideoDeviceFunctions creates and loads video functions for a device
- `CreateVideoPictureResource(imageView ImageView, imageLayout ImageLayout, codedExtent Extent2D) VideoPictureResource` - CreateVideoPictureResource creates a VideoPictureResource from an image view
- `CreateVideoPictureResourceWithOffset(imageView ImageView, imageLayout ImageLayout, codedOffset Offset2D, codedExtent Extent2D, baseArrayLayer uint32) VideoPictureResource` - CreateVideoPictureResourceWithOffset creates a VideoPictureResource with a specific offset
- `FindVideoDecodeQueueFamily(physicalDevice PhysicalDevice) (uint32, bool)` - FindVideoDecodeQueueFamily finds a queue family that supports video decode
- `FindVideoEncodeQueueFamily(physicalDevice PhysicalDevice) (uint32, bool)` - FindVideoEncodeQueueFamily finds a queue family that supports video encode
- `GetBitDepthForYUVFormat(yuvFormat YUVFormat) VideoComponentBitDepth` - GetBitDepthForYUVFormat returns the luma bit depth for a YUV format
- `GetChromaSubsamplingForYUVFormat(yuvFormat YUVFormat) VideoChromaSubsampling` - GetChromaSubsamplingForYUVFormat returns the chroma subsampling for a YUV format
- `GetVideoDeviceFunctions(device Device) *VideoDeviceFunctions` - GetVideoDeviceFunctions returns the video functions for a device
- `(vdf *VideoDeviceFunctions) IsLoaded() bool` - Returns whether the video functions are loaded
- `GetVideoFormatProperties(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo, imageUsage ImageUsageFlags) ([]VideoFormatProperties, error)` - GetVideoFormatProperties queries the video format properties for a physical device
- `LoadVideoDeviceFunctions(device Device) bool` - LoadVideoDeviceFunctions loads video device functions
- `LoadVideoFormatFunctions(instance Instance) bool` - LoadVideoFormatFunctions loads video format query functions
- `LoadVideoInstanceFunctions(instance Instance) bool` - LoadVideoInstanceFunctions loads video instance functions
- `ResetVideoDeviceFunctions()` - ResetVideoDeviceFunctions resets the device function loader
- `ResetVideoInstanceFunctions()` - ResetVideoInstanceFunctions resets the instance function loader
- `YUVFormatToVulkanFormat(yuvFormat YUVFormat) Format` - YUVFormatToVulkanFormat converts a YUV format to the corresponding Vulkan format
#### Capability Queries
- `GetSupportedVideoCodecs(physicalDevice PhysicalDevice) ([]string, error)` - Get list of supported video codecs on the device
- `GetVideoCapabilities(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo) (*VideoCapabilities, error)` - Get video codec capabilities

**Note**: To check if a specific video codec extension is supported, use `IsExtensionSupported(extensionName, availableExtensions)` with the appropriate extension name constant (e.g., `ExtensionNameVideoDecodeH264`).

#### Video Session Management
- `CreateVideoSession(device Device, createInfo *VideoSessionCreateInfo) (VideoSession, error)` - Create video session for encoding/decoding
- `DestroyVideoSession(device Device, videoSession VideoSession)` - Destroy video session
- `GetVideoSessionMemoryRequirements(device Device, videoSession VideoSession) ([]MemoryRequirements, error)` - Get memory requirements for video session
- `BindVideoSessionMemory(device Device, videoSession VideoSession, bindInfos []VideoBindMemoryInfo) error` - Bind memory to video session
- `CreateVideoSessionParameters(device Device, createInfo *VideoSessionParametersCreateInfo) (VideoSessionParameters, error)` - Create video session parameters
- `DestroyVideoSessionParameters(device Device, videoSessionParameters VideoSessionParameters)` - Destroy video session parameters

- `CreateAV1DecodeSession(device Device, createInfo *AV1DecodeSessionCreateInfo) (VideoSession, error)` - CreateAV1DecodeSession creates an AV1 decode session with the given configuration
- `CreateAV1EncodeSession(device Device, createInfo *AV1EncodeSessionCreateInfo) (VideoSession, error)` - CreateAV1EncodeSession creates an AV1 encode session with the given configuration
- `CreateDPBManager(maxSlots uint32) *DPBManager` - CreateDPBManager creates a new DPB manager with the specified number of slots
- `(dpb *DPBManager) AddSlot(imageView ImageView, imageLayout ImageLayout, poc int32) (*DPBSlot, error)` - Adds a picture to the DPB
- `(dpb *DPBManager) CalculatePOC() int32` - Calculates the Picture Order Count for the next frame
- `(dpb *DPBManager) GetReferenceSlots() []DPBSlot` - Returns all current reference slots
- `(dpb *DPBManager) MarkAsLongTerm(slotIndex int32)` - Marks a slot as a long-term reference
- `(dpb *DPBManager) RemoveOldestReference()` - Removes the oldest short-term reference from the DPB
- `CreateH264DecodeSession(device Device, createInfo *H264DecodeSessionCreateInfo) (VideoSession, error)` - CreateH264DecodeSession creates an H.264 decode session with the given configuration
- `CreateH264EncodeSession(device Device, createInfo *H264EncodeSessionCreateInfo) (VideoSession, error)` - CreateH264EncodeSession creates an H.264 encode session with the given configuration
- `CreateH265DecodeSession(device Device, createInfo *H265DecodeSessionCreateInfo) (VideoSession, error)` - CreateH265DecodeSession creates an H.265 decode session with the given configuration
- `CreateH265EncodeSession(device Device, createInfo *H265EncodeSessionCreateInfo) (VideoSession, error)` - CreateH265EncodeSession creates an H.265 encode session with the given configuration
- `DefaultAV1DecodeSessionCreateInfo(width, height uint32) *AV1DecodeSessionCreateInfo` - DefaultAV1DecodeSessionCreateInfo returns a default AV1 decode session configuration
- `DefaultAV1EncodeSessionCreateInfo(width, height uint32) *AV1EncodeSessionCreateInfo` - DefaultAV1EncodeSessionCreateInfo returns a default AV1 encode session configuration
- `DefaultH264DecodeSessionCreateInfo(width, height uint32) *H264DecodeSessionCreateInfo` - DefaultH264DecodeSessionCreateInfo returns a default H.264 decode session configuration
- `DefaultH264EncodeSessionCreateInfo(width, height uint32) *H264EncodeSessionCreateInfo` - DefaultH264EncodeSessionCreateInfo returns a default H.264 encode session configuration
- `DefaultH265DecodeSessionCreateInfo(width, height uint32) *H265DecodeSessionCreateInfo` - DefaultH265DecodeSessionCreateInfo returns a default H.265 decode session configuration
- `DefaultH265EncodeSessionCreateInfo(width, height uint32) *H265EncodeSessionCreateInfo` - DefaultH265EncodeSessionCreateInfo returns a default H.265 encode session configuration
- `UpdateVideoSessionParameters(device Device, videoSessionParameters VideoSessionParameters, updateInfo *VideoSessionParametersUpdateInfo) error` - UpdateVideoSessionParameters updates video session parameters
#### Video Coding Commands
- `CmdBeginVideoCoding(commandBuffer CommandBuffer, beginInfo *VideoBeginCodingInfo)` - Begin video coding operations
- `CmdEndVideoCoding(commandBuffer CommandBuffer)` - End video coding operations
- `CmdControlVideoCoding(commandBuffer CommandBuffer, controlInfo *VideoCodingControlInfo)` - Control video coding operations
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
if vulkan.IsExtensionSupported(vulkan.ExtensionNameVideoDecodeH264, extensions) {
    fmt.Println("H.264 hardware decode is supported")
    
    // Get video capabilities
    videoProfile := &vulkan.VideoProfileInfo{
        VideoCodecOperation: vulkan.VideoCodecOperationDecodeH264Bit,
        ChromaSubsampling:   vulkan.VideoChromaSubsampling420,
        LumaBitDepth:        vulkan.VideoComponentBitDepth8,
        ChromaBitDepth:      vulkan.VideoComponentBitDepth8,
    }
    
    caps, err := vulkan.GetVideoCapabilities(physicalDevice, videoProfile)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Max DPB slots: %d\n", caps.MaxDpbSlots)
    fmt.Printf("Max active references: %d\n", caps.MaxActiveReferencePictures)
    
    // Create video session (requires device with video queue extension enabled)
    // Note: Use an appropriate format for your video codec (e.g., NV12 for YUV 4:2:0)
    // Prerequisites:
    // - device: must be created with video queue extension enabled
    // - queueFamilyIndex: obtained from GetPhysicalDeviceQueueFamilyProperties,
    //   selecting a queue family that supports video decode operations (with QueueVideoDecodeBitKHR)
    createInfo := &vulkan.VideoSessionCreateInfo{
        QueueFamilyIndex:       queueFamilyIndex,
        VideoProfile:           videoProfile,
        PictureFormat:          vulkan.FormatR8G8B8A8Unorm,
        MaxCodedExtent:         vulkan.Extent2D{Width: 1920, Height: 1080},
        ReferencePictureFormat: vulkan.FormatR8G8B8A8Unorm,
        MaxDpbSlots:            caps.MaxDpbSlots,
        MaxActiveReferences:    caps.MaxActiveReferencePictures,
    }
    
    videoSession, err := vulkan.CreateVideoSession(device, createInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer vulkan.DestroyVideoSession(device, videoSession)
    
    // Get and bind memory for video session
    memReqs, err := vulkan.GetVideoSessionMemoryRequirements(device, videoSession)
    if err != nil {
        log.Fatal(err)
    }
    
    // Allocate and bind memory (example for first requirement)
    // Note: You can use FindMemoryType from memory.go or implement your own selector.
    // Example implementation:
    //
    // func findMemoryType(physicalDevice vulkan.PhysicalDevice, typeBits uint32, properties vulkan.MemoryPropertyFlags) (uint32, bool) {
    //     memProps := vulkan.GetPhysicalDeviceMemoryProperties(physicalDevice)
    //     for i := uint32(0); i < memProps.MemoryTypeCount; i++ {
    //         if (typeBits & (1 << i)) != 0 && (memProps.MemoryTypes[i].PropertyFlags & properties) == properties {
    //             return i, true
    //         }
    //     }
    //     return 0, false
    // }
    //
    if len(memReqs) > 0 {
        // Get memory properties from the physical device
        memProps := vulkan.GetPhysicalDeviceMemoryProperties(physicalDevice)
        
        // Use FindMemoryType from memory.go
        memTypeIndex, found := vulkan.FindMemoryType(memProps, memReqs[0].MemoryTypeBits, 0)
        if !found {
            log.Fatal("Failed to find suitable memory type")
        }
        memory, err := vulkan.AllocateMemory(device, &vulkan.MemoryAllocateInfo{
            AllocationSize:  memReqs[0].Size,
            MemoryTypeIndex: memTypeIndex,
        })
        if err != nil {
            log.Fatal(err)
        }
        
        bindInfo := []vulkan.VideoBindMemoryInfo{{
            MemoryBindIndex: 0,
            Memory:          memory,
            MemoryOffset:    0,
            MemorySize:      memReqs[0].Size,
        }}
        
        err = vulkan.BindVideoSessionMemory(device, videoSession, bindInfo)
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

**Note**: Full video codec functionality requires the Vulkan Video extensions to be enabled on the device and supported by the GPU driver. Hardware support varies by GPU model and driver version.

## Utility Functions

### Debugging and Resource Tracking
- `EnableLeakTracker()` - Turn on tracking of Vulkan object allocations
- `DisableLeakTracker()` - Turn off tracking of Vulkan object allocations
- `ClearLeaks()` - Reset the current list of tracked allocations
- `ReportLeaks() string` - Return a formatted string containing information about any un-freed resources

### Version and Feature Queries
- `GetAPIVersion() Version` - Get supported API version
- `IsExtensionSupported(extensionName string, availableExtensions []ExtensionProperties) bool` - Check extension support
- `IsLayerSupported(layerName string, availableLayers []LayerProperties) bool` - Check layer support

### Surface and WSI Integration
- `CreateXlibSurfaceKHR(instance Instance, createInfo *XlibSurfaceCreateInfoKHR) (Surface, error)` - Create Xlib surface (Linux)
- `CreateWaylandSurfaceKHR(instance Instance, createInfo *WaylandSurfaceCreateInfoKHR) (Surface, error)` - Create Wayland surface (Linux)
- `CreateWin32SurfaceKHR(instance Instance, createInfo *Win32SurfaceCreateInfoKHR) (Surface, error)` - Create Win32 surface (Windows)
- `CreateMetalSurfaceEXT(instance Instance, createInfo *MetalSurfaceCreateInfoEXT) (Surface, error)` - Create Metal surface (macOS/iOS)
- `CreateXcbSurfaceKHR(instance Instance, connection unsafe.Pointer, window uint32) (Surface, error)` - Create XCB surface (Linux)
- `GetRenderAreaGranularity(device Device, renderPass RenderPass) Extent2D` - Get render area granularity
- `GetDeviceMemoryCommitment(device Device, memory DeviceMemory) DeviceSize` - Query device memory commitment
- `IsErrorSurfaceLost(err error) bool` - Check if an error indicates that the Vulkan surface has been lost

### Surface Queries
- `GetPhysicalDeviceSurfaceSupport(physicalDevice PhysicalDevice, queueFamilyIndex uint32, surface Surface) (bool, error)` - Query if a queue family supports a surface
- `GetPhysicalDeviceSurfaceCapabilities(physicalDevice PhysicalDevice, surface Surface) (SurfaceCapabilities, error)` - Get surface capabilities
- `GetPhysicalDeviceSurfaceFormats(physicalDevice PhysicalDevice, surface Surface) ([]SurfaceFormat, error)` - Get supported surface formats
- `GetPhysicalDeviceSurfacePresentModes(physicalDevice PhysicalDevice, surface Surface) ([]PresentMode, error)` - Get supported present modes
- `DestroySurface(instance Instance, surface Surface)` - Destroy surface

### Image Layout and Subresources
- `GetImageSubresourceLayout(device Device, image Image, subresource *ImageSubresource) SubresourceLayout` - Query image subresource layout
- `TransitionImageLayout(commandBuffer CommandBuffer, image Image, format Format, oldLayout, newLayout ImageLayout, subresourceRange ImageSubresourceRange)` - Transitions an image from one layout to another

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
4. CGO is required and the appropriate Vulkan development libraries must be installed for your platform
5. The current verified Linux setup uses `libvulkan-dev`; other platforms still need their own SDK/loader validation
6. Some advanced features may require additional implementation or device extensions, and runtime support depends on the installed loader, driver, and enabled extensions





## Miscellaneous

*   **`vulkan.CmdClearAttachments`**: Clears attachment regions within a render pass.
*   **`vulkan.CmdClearColorImage`**: Clears a color image outside of a render pass.
*   **`vulkan.CmdClearDepthStencilImage`**: Clears a depth/stencil image outside of a render pass.
*   **`vulkan.CreatePipelineCache`**: Creates a pipeline cache.
*   **`vulkan.DestroyPipelineCache`**: Destroys a pipeline cache.
*   **`vulkan.GetPipelineCacheData`**: Retrieves the data from a pipeline cache.
*   **`vulkan.MergePipelineCaches`**: Merges multiple pipeline caches into a destination cache.
*   **`vulkan.CreateBufferView`**: Creates a buffer view.
*   **`vulkan.DestroyBufferView`**: Destroys a buffer view.
*   **`vulkan.GetPhysicalDeviceFormatProperties`**: Returns format properties for a physical device.
*   **`vulkan.GetPhysicalDeviceImageFormatProperties`**: Returns image format properties for a physical device.
*   **`vulkan.GetPhysicalDeviceSparseImageFormatProperties`**: Returns sparse image format properties.
*   **`vulkan.GetImageSparseMemoryRequirements`**: Returns sparse memory requirements for an image.

## Swapchain & Presentation

*   **`vulkan.CreateSwapchain`**: Creates a swapchain.
*   **`vulkan.DestroySwapchain`**: Destroys a swapchain.
*   **`vulkan.GetSwapchainImages`**: Gets the swapchain images.
*   **`vulkan.AcquireNextImage`**: Acquires the next presentable image from a swapchain. Returns the index of the next image to use, and whether the swapchain is suboptimal.
*   **`vulkan.QueuePresent`**: Queues an image for presentation. Returns true if the swapchain is suboptimal.

## Queries

*   **`vulkan.CreateQueryPool`**: Creates a query pool for managing a number of queries.
*   **`vulkan.DestroyQueryPool`**: Destroys a query pool.
*   **`vulkan.GetQueryPoolResults`**: Retrieves results from a query pool as a byte slice, or an error if the operation fails. Use QueryResult64Bit flag for 64-bit results, otherwise 32-bit results are returned.
*   **`vulkan.GetQueryPoolResultsUint32`**: Retrieves 32-bit query results.
*   **`vulkan.GetQueryPoolResultsUint64`**: Retrieves 64-bit query results.
*   **`vulkan.CmdBeginQuery`**: Begins a query.
*   **`vulkan.CmdEndQuery`**: Ends a query.
*   **`vulkan.CmdResetQueryPool`**: Resets a range of queries in a query pool on the GPU.
*   **`vulkan.CmdWriteTimestamp`**: Writes a device timestamp into a query object.
*   **`vulkan.CmdCopyQueryPoolResults`**: Copies the results of queries in a query pool to a buffer object.
*   **`vulkan.ResetQueryPool`**: Resets a range of queries in a query pool on the host (Vulkan 1.2+).
