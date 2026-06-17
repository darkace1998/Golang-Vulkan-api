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
