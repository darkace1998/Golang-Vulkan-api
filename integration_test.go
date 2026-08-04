package vulkan

import (
	"os"
	"testing"
)

// TestIntegrationInstanceCreationWithLavapipe integration tests real instance and device creation using a software renderer like lavapipe
func TestIntegrationInstanceCreationWithLavapipe(t *testing.T) {
	// Only run this test if explicitly requested or if lavapipe is expected to be available
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests; set RUN_INTEGRATION_TESTS=1 to enable")
	}

	appInfo := ApplicationInfo{
		ApplicationName: "LavapipeTestApp",
		EngineName:      "NoEngine",
		APIVersion:      Version13,
	}

	createInfo := InstanceCreateInfo{
		ApplicationInfo: &appInfo,
	}

	instance, err := CreateInstance(&createInfo)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer DestroyInstance(instance)

	physicalDevices, err := EnumeratePhysicalDevices(instance)
	if err != nil {
		t.Fatalf("Failed to enumerate physical devices: %v", err)
	}

	if len(physicalDevices) == 0 {
		t.Fatal("No physical devices found (lavapipe might not be configured correctly)")
	}

	t.Logf("Found %d physical devices", len(physicalDevices))

	// Test getting properties of the first device to confirm it works
	props := GetPhysicalDeviceProperties(physicalDevices[0])

	// Note: DeviceName is usually a byte array

	if len(props.DeviceName) == 0 {
		t.Fatal("Device name is empty")
	}
}

// TestIntegrationDeviceCreationWithLavapipe tests creating a logical device and doing basic operations
func TestIntegrationDeviceCreationWithLavapipe(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests; set RUN_INTEGRATION_TESTS=1 to enable")
	}

	appInfo := ApplicationInfo{
		ApplicationName: "LavapipeTestApp",
		EngineName:      "NoEngine",
		APIVersion:      Version13,
	}

	createInfo := InstanceCreateInfo{
		ApplicationInfo: &appInfo,
	}

	instance, err := CreateInstance(&createInfo)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer DestroyInstance(instance)

	physicalDevices, err := EnumeratePhysicalDevices(instance)
	if err != nil || len(physicalDevices) == 0 {
		t.Fatalf("Failed to enumerate physical devices: %v", err)
	}
	physicalDevice := physicalDevices[0]

	// Find a queue family
	queueFamilies := GetPhysicalDeviceQueueFamilyProperties(physicalDevice)
	if len(queueFamilies) == 0 {
		t.Fatal("No queue families found")
	}

	// Create logical device
	queuePriority := float32(1.0)
	queueCreateInfo := DeviceQueueCreateInfo{
		QueueFamilyIndex: 0,
		QueuePriorities:  []float32{queuePriority},
	}

	deviceCreateInfo := DeviceCreateInfo{
		QueueCreateInfos: []DeviceQueueCreateInfo{queueCreateInfo},
	}

	device, err := CreateDevice(physicalDevice, &deviceCreateInfo)
	if err != nil {
		t.Fatalf("Failed to create logical device: %v", err)
	}
	defer DestroyDevice(device)

	// Get device queue
	queue := GetDeviceQueue(device, 0, 0)
	if queue == nil {
		t.Fatal("Failed to get device queue: handle is nil")
	}

	// Ensure we can wait on the device without errors
	err = DeviceWaitIdle(device)
	if err != nil {
		t.Fatalf("DeviceWaitIdle failed: %v", err)
	}
}

// TestIntegrationCgoPointerRules exercises call paths that build C structs
// containing pointers to nested Go arrays. Before these arrays were pinned
// with runtime.Pinner, every one of these calls panicked under the default
// GODEBUG=cgocheck=1 with "cgo argument has Go pointer to unpinned Go pointer".
func TestIntegrationCgoPointerRules(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests; set RUN_INTEGRATION_TESTS=1 to enable")
	}

	instance, err := CreateInstance(&InstanceCreateInfo{
		ApplicationInfo: &ApplicationInfo{
			ApplicationName: "CgoPointerRulesTest",
			EngineName:      "NoEngine",
			APIVersion:      Version13,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer DestroyInstance(instance)

	physicalDevices, err := EnumeratePhysicalDevices(instance)
	if err != nil || len(physicalDevices) == 0 {
		t.Fatalf("Failed to enumerate physical devices: %v", err)
	}

	device, err := CreateDevice(physicalDevices[0], &DeviceCreateInfo{
		QueueCreateInfos: []DeviceQueueCreateInfo{
			{QueueFamilyIndex: 0, QueuePriorities: []float32{1.0}},
		},
		EnableTimelineSemaphores: true,
	})
	if err != nil {
		t.Fatalf("Failed to create logical device: %v", err)
	}
	defer DestroyDevice(device)

	// Descriptor set layout with a non-empty bindings array (pBindings).
	setLayout, err := CreateDescriptorSetLayout(device, &DescriptorSetLayoutCreateInfo{
		Bindings: []DescriptorSetLayoutBinding{
			{Binding: 0, DescriptorType: DescriptorTypeUniformBuffer, DescriptorCount: 1, StageFlags: ShaderStageVertexBit},
		},
	})
	if err != nil {
		t.Fatalf("CreateDescriptorSetLayout failed: %v", err)
	}
	defer DestroyDescriptorSetLayout(device, setLayout)

	// Pipeline layout with non-empty pSetLayouts and pPushConstantRanges.
	pipelineLayout, err := CreatePipelineLayout(device, &PipelineLayoutCreateInfo{
		SetLayouts: []DescriptorSetLayout{setLayout},
		PushConstants: []PushConstantRange{
			{StageFlags: ShaderStageVertexBit, Offset: 0, Size: 16},
		},
	})
	if err != nil {
		t.Fatalf("CreatePipelineLayout failed: %v", err)
	}
	defer DestroyPipelineLayout(device, pipelineLayout)

	// Timeline semaphore create (pNext chain) + wait (pSemaphores/pValues).
	timeline, err := CreateTimelineSemaphore(device, 5)
	if err != nil {
		t.Fatalf("CreateTimelineSemaphore failed: %v", err)
	}
	defer DestroySemaphore(device, timeline)

	waitResult, err := WaitSemaphores(device, &SemaphoreWaitInfo{
		Semaphores: []Semaphore{timeline},
		Values:     []uint64{5},
	}, 1_000_000_000)
	if err != nil {
		t.Fatalf("WaitSemaphores failed: %v", err)
	}
	if waitResult != Success {
		t.Fatalf("WaitSemaphores returned %v, want Success", waitResult)
	}

	// QueueSubmit with a non-empty command buffer array (nested Go arrays in
	// VkSubmitInfo) plus a signal semaphore.
	queue := GetDeviceQueue(device, 0, 0)
	commandPool, err := CreateCommandPool(device, &CommandPoolCreateInfo{QueueFamilyIndex: 0})
	if err != nil {
		t.Fatalf("CreateCommandPool failed: %v", err)
	}
	defer DestroyCommandPool(device, commandPool)

	commandBuffers, err := AllocateCommandBuffers(device, &CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		t.Fatalf("AllocateCommandBuffers failed: %v", err)
	}

	if err := BeginCommandBuffer(commandBuffers[0], &CommandBufferBeginInfo{}); err != nil {
		t.Fatalf("BeginCommandBuffer failed: %v", err)
	}
	if err := EndCommandBuffer(commandBuffers[0]); err != nil {
		t.Fatalf("EndCommandBuffer failed: %v", err)
	}

	signalSem, err := CreateSemaphore(device, &SemaphoreCreateInfo{})
	if err != nil {
		t.Fatalf("CreateSemaphore failed: %v", err)
	}
	defer DestroySemaphore(device, signalSem)

	fence, err := CreateFence(device, &FenceCreateInfo{})
	if err != nil {
		t.Fatalf("CreateFence failed: %v", err)
	}
	defer DestroyFence(device, fence)

	err = QueueSubmit(queue, []SubmitInfo{
		{
			CommandBuffers:   []CommandBuffer{commandBuffers[0]},
			SignalSemaphores: []Semaphore{signalSem},
		},
	}, fence)
	if err != nil {
		t.Fatalf("QueueSubmit failed: %v", err)
	}

	if _, err := WaitForFences(device, []Fence{fence}, true, 1_000_000_000); err != nil {
		t.Fatalf("WaitForFences failed: %v", err)
	}
	if err := DeviceWaitIdle(device); err != nil {
		t.Fatalf("DeviceWaitIdle failed: %v", err)
	}
}
