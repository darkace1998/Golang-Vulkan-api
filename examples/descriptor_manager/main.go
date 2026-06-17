package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
	fmt.Println("=== Vulkan Descriptor Pool Manager Example ===")

	// 1. Initialize Vulkan
	fmt.Println("Initializing Vulkan instance and device...")

	appInfo := &vulkan.ApplicationInfo{
		ApplicationName:    "Descriptor Pool Manager Example",
		ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
		EngineName:         "No Engine",
		EngineVersion:      vulkan.MakeVersion(1, 0, 0),
		APIVersion:         vulkan.Version10,
	}

	instanceInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: appInfo,
	}

	instance, err := vulkan.CreateInstance(instanceInfo)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)

	devices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil || len(devices) == 0 {
		log.Fatalf("Failed to find physical devices: %v", err)
	}
	physicalDevice := devices[0]

	deviceInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: 0,
				QueuePriorities:  []float32{1.0},
			},
		},
	}

	device, err := vulkan.CreateDevice(physicalDevice, deviceInfo)
	if err != nil {
		log.Fatalf("Failed to create logical device: %v", err)
	}
	defer vulkan.DestroyDevice(device)

	// 2. Create a Descriptor Set Layout
	fmt.Println("Creating descriptor set layout...")
	layoutBinding := vulkan.DescriptorSetLayoutBinding{
		Binding:         0,
		DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      vulkan.ShaderStageVertexBit,
	}

	layoutInfo := &vulkan.DescriptorSetLayoutCreateInfo{
		Bindings: []vulkan.DescriptorSetLayoutBinding{layoutBinding},
	}

	layout, err := vulkan.CreateDescriptorSetLayout(device, layoutInfo)
	if err != nil {
		log.Fatalf("Failed to create descriptor set layout: %v", err)
	}
	defer vulkan.DestroyDescriptorSetLayout(device, layout)

	// 3. Create a Descriptor Pool Manager
	fmt.Println("Creating descriptor pool manager (max 2 sets per pool)...")
	manager, err := vulkan.NewDescriptorPoolManager(device, 2, []vulkan.DescriptorPoolSize{
		{
			Type:            vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 2,
		},
	}, 0)
	if err != nil {
		log.Fatalf("Failed to create descriptor pool manager: %v", err)
	}
	defer func() {
		fmt.Println("Destroying descriptor pool manager...")
		manager.Destroy()
	}()

	// 4. Allocate Descriptor Sets
	fmt.Println("Allocating 3 descriptor sets (this will trigger a new pool creation under the hood)...")

	fmt.Println("Allocating first 2 sets...")
	sets1, err := manager.AllocateDescriptorSets([]vulkan.DescriptorSetLayout{layout, layout})
	if err != nil {
		log.Fatalf("Failed to allocate first 2 sets: %v", err)
	}
	fmt.Printf("Successfully allocated %d sets.\n", len(sets1))

	fmt.Println("Allocating 1 more set (requires new pool)...")
	sets2, err := manager.AllocateDescriptorSets([]vulkan.DescriptorSetLayout{layout})
	if err != nil {
		log.Fatalf("Failed to allocate 1 more set: %v", err)
	}
	fmt.Printf("Successfully allocated %d sets from new pool.\n", len(sets2))

	// 5. Reset the Manager
	fmt.Println("Resetting descriptor pool manager (releasing all pools to free list)...")
	err = manager.Reset()
	if err != nil {
		log.Fatalf("Failed to reset descriptor pool manager: %v", err)
	}
	fmt.Println("Reset successful.")

	fmt.Println("=== Example Completed Successfully ===")
}
