package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
	fmt.Println("Initializing Vulkan multi-queue example...")

	// 2. Create Instance
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Multi-Queue Example",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "Go Vulkan Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
		EnabledLayerNames:     []string{},
		EnabledExtensionNames: []string{},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)

	// 3. Enumerate Physical Devices
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate physical devices: %v", err)
	}

	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}

	gpu := physicalDevices[0] // pick the first one

	// 4. Find Queue Families
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(gpu)

	graphicsFamilyIndex := ^uint32(0)
	transferFamilyIndex := ^uint32(0)

	// Find separate graphics and transfer queues if possible
	for i, qf := range queueFamilies {
		queueFlags := qf.QueueFlags

		// Find a graphics queue
		if (queueFlags&vulkan.QueueGraphicsBit) != 0 && graphicsFamilyIndex == ^uint32(0) {
			graphicsFamilyIndex = uint32(i)
		}

		// Find a dedicated transfer queue (has transfer bit, but NOT graphics or compute)
		if (queueFlags&vulkan.QueueTransferBit) != 0 &&
			(queueFlags&vulkan.QueueGraphicsBit) == 0 &&
			(queueFlags&vulkan.QueueComputeBit) == 0 &&
			transferFamilyIndex == ^uint32(0) {
			transferFamilyIndex = uint32(i)
		}
	}

	// Fallback if no dedicated transfer queue is found
	if transferFamilyIndex == ^uint32(0) {
		fmt.Println("No dedicated transfer queue found. Falling back to any queue with transfer capabilities.")
		for i, qf := range queueFamilies {
			if (qf.QueueFlags & vulkan.QueueTransferBit) != 0 {
				transferFamilyIndex = uint32(i)
				if transferFamilyIndex != graphicsFamilyIndex {
					break // Found a different queue family
				}
			}
		}
	}

	if graphicsFamilyIndex == ^uint32(0) || transferFamilyIndex == ^uint32(0) {
		log.Fatalf("Failed to find necessary queue families (Graphics: %v, Transfer: %v)", graphicsFamilyIndex, transferFamilyIndex)
	}

	fmt.Printf("Using Graphics Queue Family: %d\n", graphicsFamilyIndex)
	fmt.Printf("Using Transfer Queue Family: %d\n", transferFamilyIndex)

	// 5. Create Logical Device with Multiple Queues
	var queueCreateInfos []vulkan.DeviceQueueCreateInfo

	// Add graphics queue
	queueCreateInfos = append(queueCreateInfos, vulkan.DeviceQueueCreateInfo{
		QueueFamilyIndex: graphicsFamilyIndex,
		QueuePriorities:  []float32{1.0},
	})

	// Add transfer queue if it's from a different family
	if graphicsFamilyIndex != transferFamilyIndex {
		queueCreateInfos = append(queueCreateInfos, vulkan.DeviceQueueCreateInfo{
			QueueFamilyIndex: transferFamilyIndex,
			QueuePriorities:  []float32{1.0},
		})
	} else {
		fmt.Println("Graphics and Transfer queues are in the same family. Using a single queue.")
		// If they are the same family, we might need multiple queues from that family if supported
		if queueFamilies[graphicsFamilyIndex].QueueCount > 1 {
			queueCreateInfos[0].QueuePriorities = []float32{1.0, 0.5} // Example priorities
			fmt.Println("Allocating 2 queues from the same family.")
		}
	}

	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos: queueCreateInfos,
		EnabledFeatures:  &vulkan.PhysicalDeviceFeatures{},
	}

	device, err := vulkan.CreateDevice(gpu, deviceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create logical device: %v", err)
	}
	defer vulkan.DestroyDevice(device)

	// 6. Get Queue Handles
	var graphicsQueue vulkan.Queue
	var transferQueue vulkan.Queue

	graphicsQueue = vulkan.GetDeviceQueue(device, graphicsFamilyIndex, 0)

	if graphicsFamilyIndex != transferFamilyIndex {
		transferQueue = vulkan.GetDeviceQueue(device, transferFamilyIndex, 0)
	} else if queueFamilies[graphicsFamilyIndex].QueueCount > 1 {
		transferQueue = vulkan.GetDeviceQueue(device, graphicsFamilyIndex, 1)
	} else {
		// Fallback: use the same queue handle
		transferQueue = graphicsQueue
	}

	if graphicsQueue == nil || transferQueue == nil {
		log.Fatalf("Failed to get device queues")
	}

	fmt.Println("Successfully retrieved graphics and transfer queues.")

	// 7. Synchronization Primitives
	semaphoreInfo := &vulkan.SemaphoreCreateInfo{}

	transferCompleteSemaphore, err := vulkan.CreateSemaphore(device, semaphoreInfo)
	if err != nil {
		log.Fatalf("Failed to create semaphore: %v", err)
	}
	defer vulkan.DestroySemaphore(device, transferCompleteSemaphore)

	fenceInfo := &vulkan.FenceCreateInfo{}
	graphicsFence, err := vulkan.CreateFence(device, fenceInfo)
	if err != nil {
		log.Fatalf("Failed to create fence: %v", err)
	}
	defer vulkan.DestroyFence(device, graphicsFence)

	// --- Mocking Async Upload ---
	fmt.Println("Submitting mock transfer operation to Transfer Queue...")

	// Create a command pool for transfer operations
	transferPoolInfo := &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: transferFamilyIndex,
	}
	transferPool, err := vulkan.CreateCommandPool(device, transferPoolInfo)
	if err != nil {
		log.Fatalf("Failed to create transfer command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, transferPool)

	transferAllocInfo := &vulkan.CommandBufferAllocateInfo{
		CommandPool:        transferPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}
	transferCmds, err := vulkan.AllocateCommandBuffers(device, transferAllocInfo)
	if err != nil || len(transferCmds) == 0 {
		log.Fatalf("Failed to allocate transfer command buffer: %v", err)
	}
	transferCmd := transferCmds[0]

	err = vulkan.BeginCommandBuffer(transferCmd, &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	})
	if err != nil {
		log.Fatalf("Failed to begin transfer command buffer: %v", err)
	}
	// Normally you'd add vkCmdCopyBuffer / vkCmdCopyBufferToImage here
	err = vulkan.EndCommandBuffer(transferCmd)
	if err != nil {
		log.Fatalf("Failed to end transfer command buffer: %v", err)
	}

	transferSubmitInfo := vulkan.SubmitInfo{
		CommandBuffers:   []vulkan.CommandBuffer{transferCmd},
		SignalSemaphores: []vulkan.Semaphore{transferCompleteSemaphore},
	}

	err = vulkan.QueueSubmit(transferQueue, []vulkan.SubmitInfo{transferSubmitInfo}, nil)
	if err != nil {
		log.Fatalf("Failed to submit to transfer queue: %v", err)
	}

	// --- Mocking Graphics Operation ---
	fmt.Println("Submitting mock graphics operation to Graphics Queue (waiting for transfer)...")

	graphicsPoolInfo := &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsFamilyIndex,
	}
	graphicsPool, err := vulkan.CreateCommandPool(device, graphicsPoolInfo)
	if err != nil {
		log.Fatalf("Failed to create graphics command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, graphicsPool)

	graphicsAllocInfo := &vulkan.CommandBufferAllocateInfo{
		CommandPool:        graphicsPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}
	graphicsCmds, err := vulkan.AllocateCommandBuffers(device, graphicsAllocInfo)
	if err != nil || len(graphicsCmds) == 0 {
		log.Fatalf("Failed to allocate graphics command buffer: %v", err)
	}
	graphicsCmd := graphicsCmds[0]

	err = vulkan.BeginCommandBuffer(graphicsCmd, &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	})
	if err != nil {
		log.Fatalf("Failed to begin graphics command buffer: %v", err)
	}
	// Normally you'd add drawing commands here
	err = vulkan.EndCommandBuffer(graphicsCmd)
	if err != nil {
		log.Fatalf("Failed to end graphics command buffer: %v", err)
	}

	graphicsSubmitInfo := vulkan.SubmitInfo{
		WaitSemaphores:   []vulkan.Semaphore{transferCompleteSemaphore},
		WaitDstStageMask: []vulkan.PipelineStageFlags{vulkan.PipelineStageFragmentShaderBit},
		CommandBuffers:   []vulkan.CommandBuffer{graphicsCmd},
	}

	err = vulkan.QueueSubmit(graphicsQueue, []vulkan.SubmitInfo{graphicsSubmitInfo}, graphicsFence)
	if err != nil {
		log.Fatalf("Failed to submit to graphics queue: %v", err)
	}

	// 8. Wait for completion
	fmt.Println("Waiting for graphics operations to finish...")
	_, err = vulkan.WaitForFences(device, []vulkan.Fence{graphicsFence}, true, ^uint64(0))
	if err != nil {
		log.Fatalf("Failed to wait for fence: %v", err)
	}

	// Clean up
	err = vulkan.DeviceWaitIdle(device)
	if err != nil {
		log.Fatalf("Failed to wait for device idle: %v", err)
	}

	fmt.Println("Multi-queue async upload completed successfully!")
}
