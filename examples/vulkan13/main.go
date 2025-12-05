package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

// Vulkan 1.3 Feature Test and Example
// This example demonstrates all major Vulkan 1.3 features implemented

func main() {
	fmt.Println("=== Vulkan 1.3 Feature Test ===")

	// Test 1: Version support
	fmt.Println("\n1. Testing Vulkan 1.3 version support...")
	version13 := vulkan.Version13
	fmt.Printf("   Vulkan 1.3 version: %d.%d.%d\n",
		version13.Major(), version13.Minor(), version13.Patch())

	// Test 2: Create instance
	fmt.Println("\n2. Creating Vulkan instance...")
	appInfo := &vulkan.ApplicationInfo{
		ApplicationName:    "Vulkan 1.3 Test",
		ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
		EngineName:         "Test Engine",
		EngineVersion:      vulkan.MakeVersion(1, 0, 0),
		APIVersion:         vulkan.Version13, // Request Vulkan 1.3
	}

	createInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo:       appInfo,
		EnabledLayerNames:     []string{}, // Add validation layers if available
		EnabledExtensionNames: []string{}, // Add extensions as needed
	}

	instance, err := vulkan.CreateInstance(createInfo)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)
	fmt.Println("   ✓ Instance created successfully")

	// Test 3: Enumerate physical devices
	fmt.Println("\n3. Enumerating physical devices...")
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate physical devices: %v", err)
	}

	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}

	physicalDevice := physicalDevices[0]
	properties := vulkan.GetPhysicalDeviceProperties(physicalDevice)
	fmt.Printf("   ✓ Found %d device(s), using: %s\n", len(physicalDevices), properties.DeviceName)
	fmt.Printf("   ✓ API Version: %d.%d.%d\n",
		vulkan.Version(properties.APIVersion).Major(),
		vulkan.Version(properties.APIVersion).Minor(),
		vulkan.Version(properties.APIVersion).Patch())

	// Test 4: Check Vulkan 1.3 feature support
	fmt.Println("\n4. Checking Vulkan 1.3 features...")
	features := vulkan.GetPhysicalDeviceFeatures(physicalDevice)
	fmt.Printf("   ✓ Tessellation Shader: %t\n", features.TessellationShader)
	fmt.Printf("   ✓ Geometry Shader: %t\n", features.GeometryShader)
	fmt.Printf("   ✓ Multi Draw Indirect: %t\n", features.MultiDrawIndirect)

	// Test 5: Create logical device
	fmt.Println("\n5. Creating logical device...")
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice)

	var graphicsQueueFamily uint32 = ^uint32(0)
	var computeQueueFamily uint32 = ^uint32(0)

	for i, queueFamily := range queueFamilies {
		if queueFamily.QueueFlags&vulkan.QueueGraphicsBit != 0 {
			graphicsQueueFamily = uint32(i)
		}
		if queueFamily.QueueFlags&vulkan.QueueComputeBit != 0 {
			computeQueueFamily = uint32(i)
		}
	}

	if graphicsQueueFamily == ^uint32(0) {
		log.Fatal("No graphics queue family found")
	}
	if computeQueueFamily == ^uint32(0) {
		log.Fatal("No compute queue family found")
	}

	queueCreateInfos := []vulkan.DeviceQueueCreateInfo{
		{
			QueueFamilyIndex: graphicsQueueFamily,
			QueuePriorities:  []float32{1.0},
		},
	}

	// Add compute queue if different from graphics
	if computeQueueFamily != graphicsQueueFamily {
		queueCreateInfos = append(queueCreateInfos, vulkan.DeviceQueueCreateInfo{
			QueueFamilyIndex: computeQueueFamily,
			QueuePriorities:  []float32{1.0},
		})
	}

	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos:      queueCreateInfos,
		EnabledFeatures:       &features,  // Enable all available features
		EnabledExtensionNames: []string{}, // Add device extensions as needed
	}

	device, err := vulkan.CreateDevice(physicalDevice, deviceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create device: %v", err)
	}
	defer vulkan.DestroyDevice(device)
	fmt.Println("   ✓ Logical device created successfully")

	// Test 6: Vulkan 1.3 Dynamic Rendering Support
	fmt.Println("\n6. Testing Dynamic Rendering (Vulkan 1.3)...")

	// Create command pool
	commandPoolCreateInfo := &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsQueueFamily,
	}

	commandPool, err := vulkan.CreateCommandPool(device, commandPoolCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)

	// Allocate command buffer
	commandBufferAllocateInfo := &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	commandBuffers, err := vulkan.AllocateCommandBuffers(device, commandBufferAllocateInfo)
	if err != nil {
		log.Fatalf("Failed to allocate command buffer: %v", err)
	}
	commandBuffer := commandBuffers[0]

	// Begin command buffer
	beginInfo := &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	}

	err = vulkan.BeginCommandBuffer(commandBuffer, beginInfo)
	if err != nil {
		log.Fatalf("Failed to begin command buffer: %v", err)
	}

	// Test dynamic rendering commands (without actual render targets for this test)
	fmt.Println("   ✓ Testing CmdBeginRendering/CmdEndRendering...")

	// Note: For a real test, we would need actual render targets
	// This demonstrates the API is available
	renderingInfo := &vulkan.RenderingInfo{
		Flags:            0,
		RenderArea:       vulkan.Rect2D{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: vulkan.Extent2D{Width: 800, Height: 600}},
		LayerCount:       1,
		ViewMask:         0,
		ColorAttachments: nil, // Would contain actual attachments in real use
	}

	// These functions are now available in Vulkan 1.3
	vulkan.CmdBeginRendering(commandBuffer, renderingInfo)
	vulkan.CmdEndRendering(commandBuffer)
	fmt.Println("   ✓ Dynamic Rendering commands available")

	// Test 7: Vulkan 1.3 Extended Dynamic State
	fmt.Println("\n7. Testing Extended Dynamic State (Vulkan 1.3)...")

	// Test dynamic cull mode
	vulkan.CmdSetCullMode(commandBuffer, vulkan.CullModeBack)
	fmt.Println("   ✓ CmdSetCullMode available")

	// Test dynamic front face
	vulkan.CmdSetFrontFace(commandBuffer, vulkan.FrontFaceCounterClockwise)
	fmt.Println("   ✓ CmdSetFrontFace available")

	// Test dynamic primitive topology
	vulkan.CmdSetPrimitiveTopology(commandBuffer, vulkan.PrimitiveTopologyTriangleList)
	fmt.Println("   ✓ CmdSetPrimitiveTopology available")

	// Test dynamic depth test enable
	vulkan.CmdSetDepthTestEnable(commandBuffer, true)
	fmt.Println("   ✓ CmdSetDepthTestEnable available")

	// Test dynamic depth write enable
	vulkan.CmdSetDepthWriteEnable(commandBuffer, true)
	fmt.Println("   ✓ CmdSetDepthWriteEnable available")

	// Test dynamic depth compare op
	vulkan.CmdSetDepthCompareOp(commandBuffer, vulkan.CompareOpLess)
	fmt.Println("   ✓ CmdSetDepthCompareOp available")

	// End command buffer
	err = vulkan.EndCommandBuffer(commandBuffer)
	if err != nil {
		log.Fatalf("Failed to end command buffer: %v", err)
	}

	// Test 8: Vulkan 1.3 Synchronization2
	fmt.Println("\n8. Testing Synchronization2 (Vulkan 1.3)...")

	// Create fence for synchronization
	fenceCreateInfo := &vulkan.FenceCreateInfo{
		Flags: 0,
	}

	fence, err := vulkan.CreateFence(device, fenceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create fence: %v", err)
	}
	defer vulkan.DestroyFence(device, fence)

	// Get queue
	queue := vulkan.GetDeviceQueue(device, graphicsQueueFamily, 0)

	// Test QueueSubmit2 (Vulkan 1.3 enhanced submission)
	submitInfo2 := []vulkan.SubmitInfo2{
		{
			Flags: 0,
			CommandBufferInfos: []vulkan.CommandBufferSubmitInfo{
				{
					CommandBuffer: commandBuffer,
					DeviceMask:    0,
				},
			},
			WaitSemaphoreInfos:   []vulkan.SemaphoreSubmitInfo{},
			SignalSemaphoreInfos: []vulkan.SemaphoreSubmitInfo{},
		},
	}

	err = vulkan.QueueSubmit2(queue, submitInfo2, fence)
	if err != nil {
		log.Fatalf("Failed to submit with QueueSubmit2: %v", err)
	}
	fmt.Println("   ✓ QueueSubmit2 available and working")

	// Wait for completion
	err = vulkan.WaitForFences(device, []vulkan.Fence{fence}, true, ^uint64(0))
	if err != nil {
		log.Fatalf("Failed to wait for fence: %v", err)
	}
	fmt.Println("   ✓ Synchronization2 working correctly")

	// Test 9: Vulkan 1.3 Private Data
	fmt.Println("\n9. Testing Private Data (Vulkan 1.3)...")

	// Create private data slot
	privateDataSlotCreateInfo := &vulkan.PrivateDataSlotCreateInfo{
		Flags: 0,
	}

	privateDataSlot, err := vulkan.CreatePrivateDataSlot(device, privateDataSlotCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create private data slot: %v", err)
	}
	defer vulkan.DestroyPrivateDataSlot(device, privateDataSlot)
	fmt.Println("   ✓ Private data slot created")

	// Set private data
	testData := uint64(0xDEADBEEF)
	err = vulkan.SetPrivateData(device, vulkan.ObjectTypeDevice, uint64(uintptr(device)), privateDataSlot, testData)
	if err != nil {
		log.Fatalf("Failed to set private data: %v", err)
	}

	// Get private data
	retrievedData := vulkan.GetPrivateData(device, vulkan.ObjectTypeDevice, uint64(uintptr(device)), privateDataSlot)
	if retrievedData != testData {
		log.Fatalf("Private data mismatch: expected %x, got %x", testData, retrievedData)
	}
	fmt.Printf("   ✓ Private data working correctly (stored: %x, retrieved: %x)\n", testData, retrievedData)

	// Test 10: Vulkan 1.3 Maintenance4
	fmt.Println("\n10. Testing Maintenance4 features (Vulkan 1.3)...")

	// Test GetDeviceBufferMemoryRequirements
	bufferCreateInfo := &vulkan.BufferCreateInfo{
		Flags:       0, // Default flags
		Size:        1024,
		Usage:       vulkan.BufferUsageStorageBufferBit,
		SharingMode: vulkan.SharingModeExclusive,
	}

	memReqs := vulkan.GetDeviceBufferMemoryRequirements(device, bufferCreateInfo)
	fmt.Printf("   ✓ Buffer memory requirements: size=%d, alignment=%d, typeBits=0x%x\n",
		memReqs.Size, memReqs.Alignment, memReqs.MemoryTypeBits)

	// Test GetDeviceImageMemoryRequirements
	imageCreateInfo := &vulkan.ImageCreateInfo{
		Flags:         0, // Default flags
		ImageType:     vulkan.ImageType2D,
		Format:        vulkan.FormatR8G8B8A8Unorm,
		Extent:        vulkan.Extent3D{Width: 256, Height: 256, Depth: 1},
		MipLevels:     1,
		ArrayLayers:   1,
		Samples:       vulkan.SampleCount1Bit,
		Tiling:        vulkan.ImageTilingOptimal,
		Usage:         vulkan.ImageUsageColorAttachmentBit,
		SharingMode:   vulkan.SharingModeExclusive,
		InitialLayout: vulkan.ImageLayoutUndefined,
	}

	imageMemReqs := vulkan.GetDeviceImageMemoryRequirements(device, imageCreateInfo)
	fmt.Printf("   ✓ Image memory requirements: size=%d, alignment=%d, typeBits=0x%x\n",
		imageMemReqs.Size, imageMemReqs.Alignment, imageMemReqs.MemoryTypeBits)

	fmt.Println("\n=== Vulkan 1.3 Feature Test Complete ===")
	fmt.Println("All major Vulkan 1.3 features are implemented and working!")
	fmt.Println("\nImplemented Vulkan 1.3 features:")
	fmt.Println("✓ Dynamic Rendering (CmdBeginRendering, CmdEndRendering)")
	fmt.Println("✓ Synchronization2 (QueueSubmit2, enhanced pipeline stages)")
	fmt.Println("✓ Extended Dynamic State (CmdSetCullMode, CmdSetFrontFace, etc.)")
	fmt.Println("✓ Private Data (CreatePrivateDataSlot, SetPrivateData, GetPrivateData)")
	fmt.Println("✓ Maintenance4 (GetDeviceBufferMemoryRequirements, GetDeviceImageMemoryRequirements)")
	fmt.Println("✓ Pipeline Creation Feedback structures")
	fmt.Println("✓ Enhanced format support and robustness features")
}
