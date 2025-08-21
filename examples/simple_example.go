package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

func main() {
	fmt.Println("Simple Vulkan Go Binding Test")
	fmt.Println("==============================")

	// Create Vulkan instance with minimal configuration
	fmt.Println("\n1. Creating Vulkan instance...")
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Simple Vulkan Test",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)
	fmt.Println("✓ Vulkan instance created successfully")

	// Enumerate physical devices
	fmt.Println("\n2. Enumerating physical devices...")
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate physical devices: %v", err)
	}

	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}

	fmt.Printf("Found %d physical device(s):\n", len(physicalDevices))

	for i, device := range physicalDevices {
		props := vulkan.GetPhysicalDeviceProperties(device)
		fmt.Printf("  Device %d: %s\n", i, props.DeviceName)
		fmt.Printf("    Type: %d, API Version: %d.%d.%d\n",
			props.DeviceType,
			props.APIVersion.Major(),
			props.APIVersion.Minor(),
			props.APIVersion.Patch())
	}

	// Test version functions
	fmt.Println("\n3. Testing version functions...")
	version := vulkan.MakeVersion(1, 3, 269)
	fmt.Printf("Created version 1.3.269: %d\n", version)
	fmt.Printf("  Major: %d, Minor: %d, Patch: %d\n",
		version.Major(), version.Minor(), version.Patch())

	// Test predefined versions
	fmt.Println("\n4. Testing predefined versions...")
	fmt.Printf("Vulkan 1.0: %d.%d.%d\n",
		vulkan.Version10.Major(), vulkan.Version10.Minor(), vulkan.Version10.Patch())
	fmt.Printf("Vulkan 1.3: %d.%d.%d\n",
		vulkan.Version13.Major(), vulkan.Version13.Minor(), vulkan.Version13.Patch())

	// Test error handling
	fmt.Println("\n5. Testing error codes...")
	fmt.Printf("Success: %s\n", vulkan.Success.Error())
	fmt.Printf("Error Out of Host Memory: %s\n", vulkan.ErrorOutOfHostMemory.Error())
	fmt.Printf("Is Success an error? %t\n", vulkan.Success.IsError())
	fmt.Printf("Is ErrorOutOfHostMemory an error? %t\n", vulkan.ErrorOutOfHostMemory.IsError())

	fmt.Println("\n==============================")
	fmt.Println("Simple Vulkan Go Binding Test Complete!")
	fmt.Println("Basic functionality working correctly.")
	fmt.Println("==============================")
}
