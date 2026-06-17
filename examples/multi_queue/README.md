# Multi-Queue Example

This example demonstrates how to discover, create, and use multiple Vulkan queues, specifically focusing on parallel transfer and graphics operations.

## Overview

Modern Vulkan applications often upload resources (like textures or buffers) asynchronously on a dedicated transfer queue while rendering continues on a graphics queue. This example shows the boilerplate required to set this up safely.

It covers:
1. Enumerating physical device queue families.
2. Identifying separate graphics and transfer queue families (with fallback logic if they share the same family).
3. Creating a logical device with multiple queues.
4. Setting up `CommandPool`s and allocating `CommandBuffer`s for both transfer and graphics operations.
5. Synchronizing the queues using a `VkSemaphore` to ensure the graphics operations only start once the transfer operations have completed.

## Running the Example

```bash
go run main.go
```

Note: In a headless environment or one without compatible hardware Vulkan drivers, this will return `VK_ERROR_INCOMPATIBLE_DRIVER`. This is expected unless a software renderer (like lavapipe) is available.
