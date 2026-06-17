# Getting Started with Golang Vulkan API

Welcome to the `golang-vulkan-api` bindings! This tutorial will walk you through the fundamentals of writing a Vulkan application in Go using this library.

Vulkan is a low-level, high-performance API. Setting up a Vulkan application involves several explicit steps. In this guide, we'll cover the core sequence:

1. **Instance Creation**: Initializing the Vulkan library.
2. **Device Selection**: Choosing a physical GPU and creating a logical device.
3. **Memory Allocation**: Allocating GPU memory and binding it to a buffer.
4. **Command Buffers**: Recording commands for the GPU.
5. **First Compute Dispatch**: Submitting a simple compute job to the GPU.

---

## 1. Instance Creation

The very first step is to initialize the Vulkan library by creating a `vulkan.Instance`. This object stores per-application state.

```go
package main

import (
	"log"
	vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
	// Describe your application
	appInfo := &vulkan.ApplicationInfo{
		ApplicationName:    "My First Vulkan App",
		ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
		EngineName:         "No Engine",
		EngineVersion:      vulkan.MakeVersion(1, 0, 0),
		APIVersion:         vulkan.Version13, // Target Vulkan 1.3
	}

	// Define instance creation details
	createInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: appInfo,
		// Enable validation layers for debugging
		// EnabledLayerNames: []string{"VK_LAYER_KHRONOS_validation"},
	}

	// Create the instance
	instance, err := vulkan.CreateInstance(createInfo)
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance) // Always clean up!

    // ... continue to Step 2
}
```

## 2. Device Selection

Next, we need to find the graphics hardware (Physical Device) and initialize a software interface to it (Logical Device).

```go
	// Enumerate available physical devices
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil || len(physicalDevices) == 0 {
		log.Fatal("Failed to find GPUs with Vulkan support")
	}

	// Simply pick the first one (in a real app, you'd score them based on features)
	physicalDevice := physicalDevices[0]

	// Find a queue family that supports compute operations
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice)
	var computeQueueFamily uint32 = ^uint32(0)

	for i, qf := range queueFamilies {
		if qf.QueueFlags & vulkan.QueueComputeBit != 0 {
			computeQueueFamily = uint32(i)
			break
		}
	}

	if computeQueueFamily == ^uint32(0) {
		log.Fatal("Could not find a compute queue family")
	}

	// Create the logical device
	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: computeQueueFamily,
				QueuePriorities:  []float32{1.0},
			},
		},
	}

	device, err := vulkan.CreateDevice(physicalDevice, deviceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create logical device: %v", err)
	}
	defer vulkan.DestroyDevice(device)

	// Get the queue handle
	queue := vulkan.GetDeviceQueue(device, computeQueueFamily, 0)
```

## 3. Memory Allocation

Unlike higher-level APIs, Vulkan requires you to manually manage memory. Creating a buffer is a two-step process: you create the buffer object, and then you allocate and bind actual memory to it.

```go
	// 1. Create the buffer object
	bufferCreateInfo := &vulkan.BufferCreateInfo{
		Size:        1024, // 1KB
		Usage:       vulkan.BufferUsageStorageBufferBit,
		SharingMode: vulkan.SharingModeExclusive,
	}

	buffer, err := vulkan.CreateBuffer(device, bufferCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create buffer: %v", err)
	}
	defer vulkan.DestroyBuffer(device, buffer)

	// 2. Find memory requirements
	memReqs := vulkan.GetBufferMemoryRequirements(device, buffer)
	memProps := vulkan.GetPhysicalDeviceMemoryProperties(physicalDevice)

	// Find memory type that is host-visible (so CPU can write to it)
	memTypeIndex, found := vulkan.FindMemoryType(
		memProps,
		memReqs.MemoryTypeBits,
		vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit,
	)
	if !found {
		log.Fatal("Suitable memory type not found")
	}

	// 3. Allocate the memory
	allocInfo := &vulkan.MemoryAllocateInfo{
		AllocationSize:  memReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	}

	memory, err := vulkan.AllocateMemory(device, allocInfo)
	if err != nil {
		log.Fatalf("Failed to allocate memory: %v", err)
	}
	defer vulkan.FreeMemory(device, memory)

	// 4. Bind the memory to the buffer
	err = vulkan.BindBufferMemory(device, buffer, memory, 0)
	if err != nil {
		log.Fatalf("Failed to bind buffer memory: %v", err)
	}
```

## 4. Command Buffers

Vulkan doesn't let you issue commands directly to the GPU. Instead, you record them into a `CommandBuffer` and submit the buffer to a queue.

```go
	// Command buffers are allocated from a Command Pool
	poolInfo := &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: computeQueueFamily,
	}

	commandPool, err := vulkan.CreateCommandPool(device, poolInfo)
	if err != nil {
		log.Fatalf("Failed to create command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)

	// Allocate a single command buffer
	allocCmdInfo := &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	cmdBuffers, err := vulkan.AllocateCommandBuffers(device, allocCmdInfo)
	if err != nil {
		log.Fatalf("Failed to allocate command buffer: %v", err)
	}
	cmdBuffer := cmdBuffers[0]
```

## 5. First Compute Dispatch

Now we record a command into our command buffer and submit it! Here we're setting up a dispatch command to run a compute shader.

*(Note: In a full program, you would also need to compile a SPIR-V compute shader, create a compute pipeline, and bind a descriptor set. See the `examples/compute/` folder for a complete runnable shader example.)*

```go
	// Begin recording commands
	beginInfo := &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	}

	if err := vulkan.BeginCommandBuffer(cmdBuffer, beginInfo); err != nil {
		log.Fatalf("Failed to begin command buffer: %v", err)
	}

	// Normally you'd bind your pipeline and descriptor sets here:
	// vulkan.CmdBindPipeline(cmdBuffer, vulkan.PipelineBindPointCompute, pipeline)
	// vulkan.CmdBindDescriptorSets(...)

	// Dispatch compute work! (e.g., 1 workgroup in X, Y, and Z)
	vulkan.CmdDispatch(cmdBuffer, 1, 1, 1)

	// Finish recording
	if err := vulkan.EndCommandBuffer(cmdBuffer); err != nil {
		log.Fatalf("Failed to end command buffer: %v", err)
	}

	// Submit the command buffer to the GPU
	submitInfo := vulkan.SubmitInfo{
		CommandBuffers: []vulkan.CommandBuffer{cmdBuffer},
	}

	// A Fence lets the CPU know when the GPU has finished
	fence, err := vulkan.CreateFence(device, &vulkan.FenceCreateInfo{})
	if err != nil {
		log.Fatalf("Failed to create fence: %v", err)
	}
	defer vulkan.DestroyFence(device, fence)

	if err := vulkan.QueueSubmit(queue, []vulkan.SubmitInfo{submitInfo}, fence); err != nil {
		log.Fatalf("Failed to submit to queue: %v", err)
	}

	// Wait for the GPU to finish execution
	vulkan.WaitForFences(device, []vulkan.Fence{fence}, true, ^uint64(0))

	log.Println("Compute dispatch completed successfully!")
```

## Next Steps

Now that you understand the core mechanics of initializing Vulkan, managing memory, and submitting command buffers:

* Take a look at the **[Examples Directory](examples/)** to see runnable code, including setting up an actual shader in `examples/compute`.
* Learn how to render graphics to a screen in `examples/swapchain` and `examples/graphics_pipeline`.
* Check out the [API Reference](API_REFERENCE.md) to explore the wrapper in detail.
