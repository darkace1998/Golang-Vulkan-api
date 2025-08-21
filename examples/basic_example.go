package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

func main() {
	fmt.Println("Vulkan Go Binding Example")
	fmt.Println("=========================")

	// Check available instance extensions
	fmt.Println("\n1. Checking available instance extensions...")
	extensions, err := vulkan.EnumerateInstanceExtensionProperties("")
	if err != nil {
		log.Fatalf("Failed to enumerate instance extensions: %v", err)
	}
	fmt.Printf("Found %d instance extensions:\n", len(extensions))
	for _, ext := range extensions[:min(10, len(extensions))] { // Show first 10
		fmt.Printf("  - %s (spec version: %d)\n", ext.ExtensionName, ext.SpecVersion)
	}
	if len(extensions) > 10 {
		fmt.Printf("  ... and %d more\n", len(extensions)-10)
	}

	// Check available instance layers
	fmt.Println("\n2. Checking available instance layers...")
	layers, err := vulkan.EnumerateInstanceLayerProperties()
	if err != nil {
		log.Fatalf("Failed to enumerate instance layers: %v", err)
	}
	fmt.Printf("Found %d instance layers:\n", len(layers))
	for _, layer := range layers {
		fmt.Printf("  - %s: %s\n", layer.LayerName, layer.Description)
	}

	// Create Vulkan instance
	fmt.Println("\n3. Creating Vulkan instance...")
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Vulkan Go Test",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "Go Vulkan Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13, // Use 1.3 since 1.4 may not be available
		},
		// Leave layers and extensions empty for now to avoid CGO issues
		EnabledLayerNames:     []string{},
		EnabledExtensionNames: []string{},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)
	fmt.Println("✓ Vulkan instance created successfully")

	// Enumerate physical devices
	fmt.Println("\n4. Enumerating physical devices...")
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate physical devices: %v", err)
	}

	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}

	fmt.Printf("Found %d physical device(s):\n", len(physicalDevices))

	var selectedDevice vulkan.PhysicalDevice
	for i, device := range physicalDevices {
		props := vulkan.GetPhysicalDeviceProperties(device)
		fmt.Printf("  Device %d: %s\n", i, props.DeviceName)
		fmt.Printf("    Type: %d, Vendor ID: 0x%x, Device ID: 0x%x\n",
			props.DeviceType, props.VendorID, props.DeviceID)
		fmt.Printf("    API Version: %d.%d.%d, Driver Version: %d.%d.%d\n",
			props.APIVersion.Major(), props.APIVersion.Minor(), props.APIVersion.Patch(),
			props.DriverVersion.Major(), props.DriverVersion.Minor(), props.DriverVersion.Patch())

		if i == 0 {
			selectedDevice = device // Use the first device
		}
	}

	// Get device features
	fmt.Println("\n5. Checking device features...")
	features := vulkan.GetPhysicalDeviceFeatures(selectedDevice)
	fmt.Printf("Device features (selected highlights):\n")
	fmt.Printf("  - Geometry Shader: %t\n", features.GeometryShader)
	fmt.Printf("  - Tessellation Shader: %t\n", features.TessellationShader)
	fmt.Printf("  - Multi Viewport: %t\n", features.MultiViewport)
	fmt.Printf("  - Sampler Anisotropy: %t\n", features.SamplerAnisotropy)

	// Get memory properties
	fmt.Println("\n6. Checking memory properties...")
	memProps := vulkan.GetPhysicalDeviceMemoryProperties(selectedDevice)
	fmt.Printf("Memory properties:\n")
	fmt.Printf("  - Memory types: %d\n", memProps.MemoryTypeCount)
	fmt.Printf("  - Memory heaps: %d\n", memProps.MemoryHeapCount)

	for i := uint32(0); i < memProps.MemoryHeapCount; i++ {
		heap := memProps.MemoryHeaps[i]
		fmt.Printf("    Heap %d: %d MB, flags: %d\n", i, heap.Size/(1024*1024), heap.Flags)
	}

	// Get queue families
	fmt.Println("\n7. Checking queue families...")
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(selectedDevice)
	fmt.Printf("Found %d queue families:\n", len(queueFamilies))

	var graphicsQueueFamily uint32 = ^uint32(0) // Invalid index
	for i, qf := range queueFamilies {
		fmt.Printf("  Queue family %d: %d queues, flags: %d\n", i, qf.QueueCount, qf.QueueFlags)
		if qf.QueueFlags&vulkan.QueueGraphicsBit != 0 && graphicsQueueFamily == ^uint32(0) {
			graphicsQueueFamily = uint32(i)
			fmt.Printf("    ✓ Graphics queue family found at index %d\n", i)
		}
	}

	if graphicsQueueFamily == ^uint32(0) {
		log.Fatal("No graphics queue family found")
	}

	// Create logical device
	fmt.Println("\n8. Creating logical device...")
	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: graphicsQueueFamily,
				QueuePriorities:  []float32{1.0},
			},
		},
		EnabledLayerNames:     []string{},
		EnabledExtensionNames: []string{},
		EnabledFeatures:       &vulkan.PhysicalDeviceFeatures{}, // Use default features
	}

	device, err := vulkan.CreateDevice(selectedDevice, deviceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create logical device: %v", err)
	}
	defer vulkan.DestroyDevice(device)
	fmt.Println("✓ Logical device created successfully")

	// Get device queue
	fmt.Println("\n9. Getting device queue...")
	queue := vulkan.GetDeviceQueue(device, graphicsQueueFamily, 0)
	if queue == nil {
		log.Fatal("Failed to get device queue")
	}
	fmt.Println("✓ Device queue obtained successfully")

	// Create command pool
	fmt.Println("\n10. Creating command pool...")
	commandPoolCreateInfo := &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsQueueFamily,
	}

	commandPool, err := vulkan.CreateCommandPool(device, commandPoolCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)
	fmt.Println("✓ Command pool created successfully")

	// Allocate command buffer
	fmt.Println("\n11. Allocating command buffer...")
	commandBufferAllocInfo := &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	}

	commandBuffers, err := vulkan.AllocateCommandBuffers(device, commandBufferAllocInfo)
	if err != nil {
		log.Fatalf("Failed to allocate command buffer: %v", err)
	}
	fmt.Printf("✓ Allocated %d command buffer(s) successfully\n", len(commandBuffers))

	// Create a simple buffer
	fmt.Println("\n12. Creating buffer...")
	bufferCreateInfo := &vulkan.BufferCreateInfo{
		Flags:       0,    // Default flags
		Size:        1024, // 1KB buffer
		Usage:       vulkan.BufferUsageVertexBufferBit,
		SharingMode: vulkan.SharingModeExclusive,
	}

	buffer, err := vulkan.CreateBuffer(device, bufferCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create buffer: %v", err)
	}
	defer vulkan.DestroyBuffer(device, buffer)
	fmt.Println("✓ Buffer created successfully")

	// Get buffer memory requirements
	fmt.Println("\n13. Getting buffer memory requirements...")
	memRequirements := vulkan.GetBufferMemoryRequirements(device, buffer)
	fmt.Printf("Buffer memory requirements:\n")
	fmt.Printf("  - Size: %d bytes\n", memRequirements.Size)
	fmt.Printf("  - Alignment: %d bytes\n", memRequirements.Alignment)
	fmt.Printf("  - Memory type bits: 0x%x\n", memRequirements.MemoryTypeBits)

	// Find suitable memory type
	memTypeIndex, found := vulkan.FindMemoryType(memProps, memRequirements.MemoryTypeBits,
		vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit)
	if !found {
		log.Fatal("Failed to find suitable memory type")
	}
	fmt.Printf("✓ Found suitable memory type at index %d\n", memTypeIndex)

	// Allocate memory
	fmt.Println("\n14. Allocating memory...")
	memAllocInfo := &vulkan.MemoryAllocateInfo{
		AllocationSize:  memRequirements.Size,
		MemoryTypeIndex: memTypeIndex,
	}

	memory, err := vulkan.AllocateMemory(device, memAllocInfo)
	if err != nil {
		log.Fatalf("Failed to allocate memory: %v", err)
	}
	defer vulkan.FreeMemory(device, memory)
	fmt.Println("✓ Memory allocated successfully")

	// Bind buffer memory
	fmt.Println("\n15. Binding buffer memory...")
	err = vulkan.BindBufferMemory(device, buffer, memory, 0)
	if err != nil {
		log.Fatalf("Failed to bind buffer memory: %v", err)
	}
	fmt.Println("✓ Buffer memory bound successfully")

	// Test synchronization objects
	fmt.Println("\n16. Creating synchronization objects...")

	// Create semaphore
	semaphoreCreateInfo := &vulkan.SemaphoreCreateInfo{}
	semaphore, err := vulkan.CreateSemaphore(device, semaphoreCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create semaphore: %v", err)
	}
	defer vulkan.DestroySemaphore(device, semaphore)
	fmt.Println("✓ Semaphore created successfully")

	// Create fence
	fenceCreateInfo := &vulkan.FenceCreateInfo{
		Flags: 0, // Unsignaled
	}
	fence, err := vulkan.CreateFence(device, fenceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create fence: %v", err)
	}
	defer vulkan.DestroyFence(device, fence)
	fmt.Println("✓ Fence created successfully")

	// Test command buffer recording
	fmt.Println("\n17. Recording command buffer...")
	beginInfo := &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	}

	err = vulkan.BeginCommandBuffer(commandBuffers[0], beginInfo)
	if err != nil {
		log.Fatalf("Failed to begin command buffer: %v", err)
	}

	err = vulkan.EndCommandBuffer(commandBuffers[0])
	if err != nil {
		log.Fatalf("Failed to end command buffer: %v", err)
	}
	fmt.Println("✓ Command buffer recorded successfully")

	// Submit command buffer
	fmt.Println("\n18. Submitting command buffer...")
	submitInfo := vulkan.SubmitInfo{
		CommandBuffers: []vulkan.CommandBuffer{commandBuffers[0]},
	}

	err = vulkan.QueueSubmit(queue, []vulkan.SubmitInfo{submitInfo}, fence)
	if err != nil {
		log.Fatalf("Failed to submit command buffer: %v", err)
	}
	fmt.Println("✓ Command buffer submitted successfully")

	// Wait for completion
	fmt.Println("\n19. Waiting for command completion...")
	err = vulkan.WaitForFences(device, []vulkan.Fence{fence}, true, ^uint64(0)) // Wait indefinitely
	if err != nil {
		log.Fatalf("Failed to wait for fence: %v", err)
	}
	fmt.Println("✓ Command completed successfully")

	// Clean up
	fmt.Println("\n20. Waiting for device idle...")
	err = vulkan.DeviceWaitIdle(device)
	if err != nil {
		log.Fatalf("Failed to wait for device idle: %v", err)
	}
	fmt.Println("✓ Device is idle")

	fmt.Println("\n=========================")
	fmt.Println("Vulkan Go Binding Test Complete!")
	fmt.Println("All core functionality working correctly.")
	fmt.Println("=========================")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
