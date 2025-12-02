package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

// Example compute shader (SPIR-V bytecode would be loaded here)
// This is a simple example that doubles input values
var computeShaderCode = []uint32{
	// This would be actual SPIR-V bytecode compiled from GLSL/HLSL
	// For demonstration, we use placeholder data
	0x07230203, 0x00010000, 0x000d000a, 0x0000002e,
	// ... actual shader bytecode would go here
}

func main() {
	fmt.Println("=== Vulkan Compute Layer Example ===")
	fmt.Println("Demonstrating compute shaders for AI workloads")

	// Create Vulkan instance
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Vulkan Compute AI Example",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "AI Compute Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatal("Failed to create Vulkan instance:", err)
	}
	defer vulkan.DestroyInstance(instance)

	// Get physical devices
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatal("Failed to enumerate physical devices:", err)
	}

	if len(physicalDevices) == 0 {
		log.Fatal("No Vulkan-capable devices found")
	}

	// Use the first device and check for compute queue support
	physicalDevice := physicalDevices[0]
	properties := vulkan.GetPhysicalDeviceProperties(physicalDevice)
	fmt.Printf("Using device: %s\n", properties.DeviceName)

	// Check compute capabilities
	fmt.Printf("Max compute shared memory: %d bytes\n", properties.Limits.MaxComputeSharedMemorySize)
	fmt.Printf("Max compute work group size: [%d, %d, %d]\n",
		properties.Limits.MaxComputeWorkGroupSize[0],
		properties.Limits.MaxComputeWorkGroupSize[1],
		properties.Limits.MaxComputeWorkGroupSize[2])
	fmt.Printf("Max compute work group invocations: %d\n", properties.Limits.MaxComputeWorkGroupInvocations)

	// Find compute queue family
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice)
	var computeQueueFamily uint32 = ^uint32(0) // Invalid index

	for i, family := range queueFamilies {
		if family.QueueFlags&vulkan.QueueComputeBit != 0 {
			computeQueueFamily = uint32(i)
			fmt.Printf("Found compute queue family at index %d\n", i)
			fmt.Printf("Queue count: %d\n", family.QueueCount)
			break
		}
	}

	if computeQueueFamily == ^uint32(0) {
		log.Fatal("No compute queue family found")
	}

	// Create logical device with compute queue
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
		log.Fatal("Failed to create device:", err)
	}
	defer vulkan.DestroyDevice(device)

	// Get compute queue
	computeQueue := vulkan.GetDeviceQueue(device, computeQueueFamily, 0)
	fmt.Println("Compute queue obtained successfully")

	// Create storage buffers for compute workload
	bufferSize := vulkan.DeviceSize(1024 * 1024) // 1MB buffer for AI data

	// Input buffer
	inputBuffer, err := vulkan.CreateBuffer(device, &vulkan.BufferCreateInfo{
		Flags:       0, // Default flags
		Size:        bufferSize,
		Usage:       vulkan.BufferUsageStorageBufferBit | vulkan.BufferUsageTransferDstBit,
		SharingMode: vulkan.SharingModeExclusive,
	})
	if err != nil {
		log.Fatal("Failed to create input buffer:", err)
	}
	defer vulkan.DestroyBuffer(device, inputBuffer)

	// Output buffer
	outputBuffer, err := vulkan.CreateBuffer(device, &vulkan.BufferCreateInfo{
		Flags:       0, // Default flags
		Size:        bufferSize,
		Usage:       vulkan.BufferUsageStorageBufferBit | vulkan.BufferUsageTransferSrcBit,
		SharingMode: vulkan.SharingModeExclusive,
	})
	if err != nil {
		log.Fatal("Failed to create output buffer:", err)
	}
	defer vulkan.DestroyBuffer(device, outputBuffer)

	// Get memory requirements and allocate memory
	memProps := vulkan.GetPhysicalDeviceMemoryProperties(physicalDevice)

	inputMemReqs := vulkan.GetBufferMemoryRequirements(device, inputBuffer)
	memoryType, found := vulkan.FindMemoryType(memProps, inputMemReqs.MemoryTypeBits,
		vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit)
	if !found {
		log.Fatal("Failed to find suitable memory type for input buffer")
	}

	inputMemory, err := vulkan.AllocateMemory(device, &vulkan.MemoryAllocateInfo{
		AllocationSize:  inputMemReqs.Size,
		MemoryTypeIndex: memoryType,
	})
	if err != nil {
		log.Fatal("Failed to allocate input memory:", err)
	}
	defer vulkan.FreeMemory(device, inputMemory)

	err = vulkan.BindBufferMemory(device, inputBuffer, inputMemory, 0)
	if err != nil {
		log.Fatal("Failed to bind input buffer memory:", err)
	}

	// Allocate output memory
	outputMemReqs := vulkan.GetBufferMemoryRequirements(device, outputBuffer)
	outputMemory, err := vulkan.AllocateMemory(device, &vulkan.MemoryAllocateInfo{
		AllocationSize:  outputMemReqs.Size,
		MemoryTypeIndex: memoryType,
	})
	if err != nil {
		log.Fatal("Failed to allocate output memory:", err)
	}
	defer vulkan.FreeMemory(device, outputMemory)

	err = vulkan.BindBufferMemory(device, outputBuffer, outputMemory, 0)
	if err != nil {
		log.Fatal("Failed to bind output buffer memory:", err)
	}

	fmt.Println("Storage buffers created and memory allocated successfully")

	// Create compute shader module (in a real application, you'd load compiled SPIR-V)
	shaderModule, err := vulkan.CreateShaderModule(device, &vulkan.ShaderModuleCreateInfo{
		CodeSize: uint32(len(computeShaderCode) * 4), // 4 bytes per uint32
		Code:     computeShaderCode,
	})
	if err != nil {
		log.Fatal("Failed to create shader module:", err)
	}
	defer vulkan.DestroyShaderModule(device, shaderModule)

	// Create descriptor set layout for storage buffers
	descriptorSetLayout, err := vulkan.CreateDescriptorSetLayout(device, &vulkan.DescriptorSetLayoutCreateInfo{
		Bindings: []vulkan.DescriptorSetLayoutBinding{
			{
				Binding:         0,
				DescriptorType:  vulkan.DescriptorTypeStorageBuffer,
				DescriptorCount: 1,
				StageFlags:      vulkan.ShaderStageComputeBit,
			},
			{
				Binding:         1,
				DescriptorType:  vulkan.DescriptorTypeStorageBuffer,
				DescriptorCount: 1,
				StageFlags:      vulkan.ShaderStageComputeBit,
			},
		},
	})
	if err != nil {
		log.Fatal("Failed to create descriptor set layout:", err)
	}
	defer vulkan.DestroyDescriptorSetLayout(device, descriptorSetLayout)

	// Create pipeline layout
	pipelineLayout, err := vulkan.CreatePipelineLayout(device, &vulkan.PipelineLayoutCreateInfo{
		SetLayouts: []vulkan.DescriptorSetLayout{descriptorSetLayout},
	})
	if err != nil {
		log.Fatal("Failed to create pipeline layout:", err)
	}
	defer vulkan.DestroyPipelineLayout(device, pipelineLayout)

	// Create compute pipeline
	computePipelines, err := vulkan.CreateComputePipelines(device, vulkan.PipelineCache(nil), []vulkan.ComputePipelineCreateInfo{
		{
			Stage: vulkan.PipelineShaderStageCreateInfo{
				Stage:  vulkan.ShaderStageComputeBit,
				Module: shaderModule,
				Name:   "main",
			},
			Layout: pipelineLayout,
		},
	})
	if err != nil {
		log.Fatal("Failed to create compute pipeline:", err)
	}
	defer vulkan.DestroyPipeline(device, computePipelines[0])

	fmt.Println("Compute pipeline created successfully")

	// Create command pool and command buffer
	commandPool, err := vulkan.CreateCommandPool(device, &vulkan.CommandPoolCreateInfo{
		QueueFamilyIndex: computeQueueFamily,
	})
	if err != nil {
		log.Fatal("Failed to create command pool:", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)

	commandBuffers, err := vulkan.AllocateCommandBuffers(device, &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		log.Fatal("Failed to allocate command buffer:", err)
	}

	commandBuffer := commandBuffers[0]

	// Record compute commands
	err = vulkan.BeginCommandBuffer(commandBuffer, &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	})
	if err != nil {
		log.Fatal("Failed to begin command buffer:", err)
	}

	// Bind compute pipeline
	vulkan.CmdBindPipeline(commandBuffer, vulkan.PipelineBindPointCompute, computePipelines[0])

	// Bind descriptor sets (would be created and updated in a complete implementation)
	// vulkan.CmdBindDescriptorSets(commandBuffer, vulkan.PipelineBindPointCompute, pipelineLayout, 0, []vulkan.DescriptorSet{descriptorSet}, nil)

	// Dispatch compute work
	// For AI workloads, you'd typically dispatch work groups that process chunks of data
	workGroupSize := uint32(64)       // Typical size for AI compute
	numElements := uint32(1024 * 256) // Number of elements to process
	numWorkGroups := (numElements + workGroupSize - 1) / workGroupSize

	vulkan.CmdDispatch(commandBuffer, numWorkGroups, 1, 1)

	// Add pipeline barrier to ensure compute work is complete
	vulkan.CmdPipelineBarrier(commandBuffer,
		vulkan.PipelineStageComputeShaderBit,
		vulkan.PipelineStageTransferBit,
		0)

	err = vulkan.EndCommandBuffer(commandBuffer)
	if err != nil {
		log.Fatal("Failed to end command buffer:", err)
	}

	fmt.Printf("Compute commands recorded successfully\n")
	fmt.Printf("Work groups dispatched: %d\n", numWorkGroups)
	fmt.Printf("Elements per work group: %d\n", workGroupSize)

	// Create fence for synchronization
	fence, err := vulkan.CreateFence(device, &vulkan.FenceCreateInfo{})
	if err != nil {
		log.Fatal("Failed to create fence:", err)
	}
	defer vulkan.DestroyFence(device, fence)

	// Submit compute work to queue
	err = vulkan.QueueSubmit(computeQueue, []vulkan.SubmitInfo{
		{
			CommandBuffers: []vulkan.CommandBuffer{commandBuffer},
		},
	}, fence)
	if err != nil {
		log.Fatal("Failed to submit to queue:", err)
	}

	// Wait for completion
	err = vulkan.WaitForFences(device, []vulkan.Fence{fence}, true, ^uint64(0))
	if err != nil {
		log.Fatal("Failed to wait for fence:", err)
	}

	fmt.Println("✅ Compute shader execution completed successfully!")
	fmt.Println()
	fmt.Println("=== Vulkan Compute Layer Features Demonstrated ===")
	fmt.Println("✅ Compute queue family detection and device creation")
	fmt.Println("✅ Storage buffer creation for large AI datasets")
	fmt.Println("✅ Compute shader module creation")
	fmt.Println("✅ Compute pipeline creation and binding")
	fmt.Println("✅ Compute dispatch commands (CmdDispatch)")
	fmt.Println("✅ Pipeline barriers for compute synchronization")
	fmt.Println("✅ GPU/CPU synchronization with fences")
	fmt.Println()
	fmt.Println("This implementation provides all the building blocks needed for:")
	fmt.Println("- Neural network inference")
	fmt.Println("- Matrix operations and linear algebra")
	fmt.Println("- Parallel data processing")
	fmt.Println("- Custom AI compute kernels")
	fmt.Println("- Memory-efficient large dataset processing")
}
