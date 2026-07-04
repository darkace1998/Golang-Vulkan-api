/*
Package vulkan provides a type-safe Go interface to the Vulkan 1.3+ graphics and compute APIs.

It is designed to be used as a library for other Go projects that need low-level graphics and compute functionality, bridging the gap between Go and the underlying C Vulkan API (libvulkan).

# Overview

This library features:
  - Full core Vulkan 1.3 API coverage needed for most 3D and compute applications.
  - Dynamic Rendering (VK_KHR_dynamic_rendering) out of the box.
  - Enhanced synchronization with Synchronization2 (VK_KHR_synchronization2).
  - Explicit memory management with LeakTracker integration for safe resource tracking.
  - Compute shader capabilities suitable for AI/ML and parallel tasks.
  - Hardware-accelerated video decoding/encoding through Vulkan Video extensions.

# Initializing Vulkan

The first step in any Vulkan application is initializing the library by creating an Instance:

	appInfo := &vulkan.ApplicationInfo{
		ApplicationName:    "My First Vulkan App",
		ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
		EngineName:         "No Engine",
		EngineVersion:      vulkan.MakeVersion(1, 0, 0),
		APIVersion:         vulkan.Version13, // Target Vulkan 1.3
	}

	createInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: appInfo,
	}

	instance, err := vulkan.CreateInstance(createInfo)
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)

# Selecting a Device

After initialization, you must select a physical device (GPU) and create a logical device interface:

	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil || len(physicalDevices) == 0 {
		log.Fatal("Failed to find GPUs with Vulkan support")
	}

	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		// configure queue create infos, features, and extensions
	}

	device, err := vulkan.CreateDevice(physicalDevices[0], deviceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create logical device: %v", err)
	}
	defer vulkan.DestroyDevice(device)

# Error Handling

This package uses two main error types:
  - ValidationError: Indicates that API detected invalid input (e.g., nil pointers) before calling the Vulkan C API.
  - VulkanError: Indicates that the underlying Vulkan C API call failed.

You can inspect the error with errors.As or functions like IsVulkanError(). Transient errors like VK_ERROR_DEVICE_LOST or VK_ERROR_OUT_OF_DATE_KHR can be handled by rebuilding the context or swapchain.

# Thread Safety

The package is largely thread-safe for reading. Functions that create or destroy Vulkan objects are thread-safe with respect to the parent Instance/Device. However, modifying the same Vulkan object concurrently from multiple goroutines (e.g., recording to the same CommandBuffer simultaneously) requires explicit external synchronization (e.g. sync.Mutex). Note that video extension loading functions (LoadVideoDeviceFunctions, LoadVideoInstanceFunctions) must be executed from a single thread.
*/
package vulkan
