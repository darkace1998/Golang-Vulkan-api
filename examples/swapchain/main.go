package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

// Swapchain Example
// Demonstrates swapchain creation, image acquisition, and presentation flow.
// NOTE: A real swapchain requires a window surface (e.g. via GLFW or SDL).
// This example shows the full API usage and validates all types and constants,
// then performs the steps that can run without a window surface.

func main() {
	fmt.Println("=== Vulkan Swapchain Example ===")
	fmt.Println("Demonstrating swapchain types, constants, and API usage")
	fmt.Println()

	// -----------------------------------------------------------------------
	// 1. Create Vulkan instance
	// -----------------------------------------------------------------------
	fmt.Println("1. Creating Vulkan instance...")
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Swapchain Example",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "Example Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)
	fmt.Println("   ✓ Vulkan instance created")

	// -----------------------------------------------------------------------
	// 2. Select physical device and find a graphics queue family
	// -----------------------------------------------------------------------
	fmt.Println("\n2. Selecting physical device...")
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate physical devices: %v", err)
	}
	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}

	physicalDevice := physicalDevices[0]
	props := vulkan.GetPhysicalDeviceProperties(physicalDevice)
	fmt.Printf("   ✓ Using device: %s\n", props.DeviceName)

	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice)
	var graphicsQueueFamily uint32 = ^uint32(0)
	for i, qf := range queueFamilies {
		if qf.QueueFlags&vulkan.QueueGraphicsBit != 0 {
			graphicsQueueFamily = uint32(i)
			break
		}
	}
	if graphicsQueueFamily == ^uint32(0) {
		log.Fatal("No graphics queue family found")
	}
	fmt.Printf("   ✓ Graphics queue family: %d\n", graphicsQueueFamily)

	// -----------------------------------------------------------------------
	// 3. Create logical device
	// -----------------------------------------------------------------------
	fmt.Println("\n3. Creating logical device...")
	device, err := vulkan.CreateDevice(physicalDevice, &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: graphicsQueueFamily,
				QueuePriorities:  []float32{1.0},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create device: %v", err)
	}
	defer vulkan.DestroyDevice(device)
	fmt.Println("   ✓ Logical device created")

	// -----------------------------------------------------------------------
	// 4. Demonstrate swapchain types and constants
	// -----------------------------------------------------------------------
	fmt.Println("\n4. Validating swapchain types and constants...")

	// Color spaces
	colorSpaces := map[string]vulkan.ColorSpace{
		"SRGB Nonlinear":          vulkan.ColorSpaceSRGBNonlinear,
		"Display P3 Nonlinear":    vulkan.ColorSpaceDisplayP3Nonlinear,
		"Extended SRGB Linear":    vulkan.ColorSpaceExtendedSRGBLinear,
		"BT709 Linear":            vulkan.ColorSpaceBT709Linear,
		"BT2020 Linear":           vulkan.ColorSpaceBT2020Linear,
		"HDR10 ST2084":            vulkan.ColorSpaceHDR10ST2084,
		"Dolby Vision":            vulkan.ColorSpaceDolbyVision,
		"Adobe RGB Linear":        vulkan.ColorSpaceAdobeRGBLinear,
		"Adobe RGB Nonlinear":     vulkan.ColorSpaceAdobeRGBNonlinear,
		"Extended SRGB Nonlinear": vulkan.ColorSpaceExtendedSRGBNonlinear,
	}
	fmt.Printf("   ✓ Color spaces defined: %d\n", len(colorSpaces))

	// Surface transforms
	transforms := map[string]vulkan.SurfaceTransformFlags{
		"Identity":          vulkan.SurfaceTransformIdentity,
		"Rotate 90":         vulkan.SurfaceTransformRotate90,
		"Rotate 180":        vulkan.SurfaceTransformRotate180,
		"Rotate 270":        vulkan.SurfaceTransformRotate270,
		"Horizontal Mirror": vulkan.SurfaceTransformHorizontalMirror,
		"Inherit":           vulkan.SurfaceTransformInherit,
	}
	fmt.Printf("   ✓ Surface transforms defined: %d\n", len(transforms))

	// Composite alpha modes
	compositeAlphas := map[string]vulkan.CompositeAlphaFlags{
		"Opaque":          vulkan.CompositeAlphaOpaque,
		"Pre-multiplied":  vulkan.CompositeAlphaPreMultiplied,
		"Post-multiplied": vulkan.CompositeAlphaPostMultiplied,
		"Inherit":         vulkan.CompositeAlphaInherit,
	}
	fmt.Printf("   ✓ Composite alpha modes defined: %d\n", len(compositeAlphas))

	// Swapchain create flags
	_ = vulkan.SwapchainCreateSplitInstanceBindRegions
	_ = vulkan.SwapchainCreateProtected
	_ = vulkan.SwapchainCreateMutableFormat
	fmt.Println("   ✓ Swapchain create flags validated")

	// -----------------------------------------------------------------------
	// 5. Demonstrate SwapchainCreateInfo structure
	// -----------------------------------------------------------------------
	fmt.Println("\n5. Demonstrating SwapchainCreateInfo structure...")

	// In a real application you would obtain a Surface from a windowing system
	// (e.g. via GLFW: glfwCreateWindowSurface, or SDL: SDL_Vulkan_CreateSurface).
	// Here we show the full struct layout without calling CreateSwapchain because
	// we have no surface.
	fmt.Println("   SwapchainCreateInfo fields:")
	fmt.Println("     Flags              - SwapchainCreateFlags (e.g. MutableFormat)")
	fmt.Println("     Surface            - VkSurfaceKHR handle (from windowing system)")
	fmt.Println("     MinImageCount      - Minimum number of presentable images")
	fmt.Println("     ImageFormat        - e.g. FormatB8G8R8A8Srgb")
	fmt.Println("     ImageColorSpace    - e.g. ColorSpaceSRGBNonlinear")
	fmt.Println("     ImageExtent        - Width × Height of the swapchain images")
	fmt.Println("     ImageArrayLayers   - Number of image array layers (1 for non-stereo)")
	fmt.Println("     ImageUsage         - e.g. ImageUsageColorAttachmentBit")
	fmt.Println("     ImageSharingMode   - Exclusive or Concurrent")
	fmt.Println("     QueueFamilyIndices - Queue families (for concurrent mode)")
	fmt.Println("     PreTransform       - Surface transform (e.g. Identity)")
	fmt.Println("     CompositeAlpha     - e.g. CompositeAlphaOpaque")
	fmt.Println("     PresentMode        - e.g. PresentModeFIFO")
	fmt.Println("     Clipped            - Whether obscured pixels can be discarded")
	fmt.Println("     OldSwapchain       - Previous swapchain for recycling")

	// -----------------------------------------------------------------------
	// 6. Validate input validation on swapchain functions
	// -----------------------------------------------------------------------
	fmt.Println("\n6. Validating input validation...")

	// CreateSwapchain – nil device
	_, err = vulkan.CreateSwapchain(nil, &vulkan.SwapchainCreateInfo{})
	if err != nil {
		fmt.Printf("   ✓ CreateSwapchain(nil device) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ CreateSwapchain(nil device) should have been rejected")
	}

	// CreateSwapchain – nil createInfo
	_, err = vulkan.CreateSwapchain(device, nil)
	if err != nil {
		fmt.Printf("   ✓ CreateSwapchain(nil createInfo) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ CreateSwapchain(nil createInfo) should have been rejected")
	}

	// CreateSwapchain – nil surface
	_, err = vulkan.CreateSwapchain(device, &vulkan.SwapchainCreateInfo{
		MinImageCount:    2,
		ImageExtent:      vulkan.Extent2D{Width: 800, Height: 600},
		ImageArrayLayers: 1,
	})
	if err != nil {
		fmt.Printf("   ✓ CreateSwapchain(nil Surface) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ CreateSwapchain(nil Surface) should have been rejected")
	}

	// CreateSwapchain – zero extent
	_, err = vulkan.CreateSwapchain(device, &vulkan.SwapchainCreateInfo{
		MinImageCount:    2,
		ImageArrayLayers: 1,
	})
	if err != nil {
		fmt.Printf("   ✓ CreateSwapchain(zero extent) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ CreateSwapchain(zero extent) should have been rejected")
	}

	// GetSwapchainImages – nil device
	_, err = vulkan.GetSwapchainImages(nil, nil)
	if err != nil {
		fmt.Printf("   ✓ GetSwapchainImages(nil device) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ GetSwapchainImages(nil device) should have been rejected")
	}

	// AcquireNextImage – nil device
	_, _, err = vulkan.AcquireNextImage(nil, nil, 0, nil, nil)
	if err != nil {
		fmt.Printf("   ✓ AcquireNextImage(nil device) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ AcquireNextImage(nil device) should have been rejected")
	}

	// QueuePresent – nil queue
	_, err = vulkan.QueuePresent(nil, &vulkan.PresentInfo{
		Swapchains:   []vulkan.Swapchain{nil},
		ImageIndices: []uint32{0},
	})
	if err != nil {
		fmt.Printf("   ✓ QueuePresent(nil queue) correctly rejected: %v\n", err)
	} else {
		log.Fatal("   ✗ QueuePresent(nil queue) should have been rejected")
	}

	// -----------------------------------------------------------------------
	// 7. Show the typical swapchain render loop (pseudo-code)
	// -----------------------------------------------------------------------
	fmt.Println("\n7. Typical swapchain render loop:")
	fmt.Println("   ┌─────────────────────────────────────────────────────┐")
	fmt.Println("   │ 1. AcquireNextImage(device, swapchain, timeout,    │")
	fmt.Println("   │       imageAvailableSemaphore, nil)                 │")
	fmt.Println("   │ 2. Record command buffer for image[index]          │")
	fmt.Println("   │ 3. QueueSubmit(queue, submitInfo, renderFence)     │")
	fmt.Println("   │       waitSemaphores: [imageAvailableSemaphore]    │")
	fmt.Println("   │       signalSemaphores: [renderFinishedSemaphore]  │")
	fmt.Println("   │ 4. QueuePresent(queue, presentInfo)                │")
	fmt.Println("   │       waitSemaphores: [renderFinishedSemaphore]    │")
	fmt.Println("   │       swapchains: [swapchain]                      │")
	fmt.Println("   │       imageIndices: [index]                        │")
	fmt.Println("   │ 5. Handle suboptimal / out-of-date → recreate     │")
	fmt.Println("   └─────────────────────────────────────────────────────┘")

	// -----------------------------------------------------------------------
	// 8. Create synchronization objects that a real swapchain loop would use
	// -----------------------------------------------------------------------
	fmt.Println("\n8. Creating synchronization objects for swapchain loop...")

	imageAvailableSemaphore, err := vulkan.CreateSemaphore(device, &vulkan.SemaphoreCreateInfo{})
	if err != nil {
		log.Fatalf("Failed to create semaphore: %v", err)
	}
	defer vulkan.DestroySemaphore(device, imageAvailableSemaphore)
	fmt.Println("   ✓ Image-available semaphore created")

	renderFinishedSemaphore, err := vulkan.CreateSemaphore(device, &vulkan.SemaphoreCreateInfo{})
	if err != nil {
		log.Fatalf("Failed to create semaphore: %v", err)
	}
	defer vulkan.DestroySemaphore(device, renderFinishedSemaphore)
	fmt.Println("   ✓ Render-finished semaphore created")

	inFlightFence, err := vulkan.CreateFence(device, &vulkan.FenceCreateInfo{
		Flags: vulkan.FenceCreateSignaledBit,
	})
	if err != nil {
		log.Fatalf("Failed to create fence: %v", err)
	}
	defer vulkan.DestroyFence(device, inFlightFence)
	fmt.Println("   ✓ In-flight fence created (signaled)")

	// -----------------------------------------------------------------------
	// 9. Create command pool and buffer
	// -----------------------------------------------------------------------
	fmt.Println("\n9. Creating command pool and buffer...")
	commandPool, err := vulkan.CreateCommandPool(device, &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsQueueFamily,
	})
	if err != nil {
		log.Fatalf("Failed to create command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)

	commandBuffers, err := vulkan.AllocateCommandBuffers(device, &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		log.Fatalf("Failed to allocate command buffers: %v", err)
	}
	fmt.Printf("   ✓ Allocated %d command buffer(s)\n", len(commandBuffers))

	// -----------------------------------------------------------------------
	// Summary
	// -----------------------------------------------------------------------
	fmt.Println("\n=== Swapchain Example Complete ===")
	fmt.Println("✓ All swapchain types and constants validated")
	fmt.Println("✓ Input validation verified for all swapchain functions")
	fmt.Println("✓ Synchronization objects created for present loop")
	fmt.Println("✓ Command pool and buffer ready for rendering")
	fmt.Println()
	fmt.Println("To use a real swapchain, add a windowing library (GLFW, SDL)")
	fmt.Println("and create a VkSurfaceKHR before calling CreateSwapchain.")
	fmt.Println()
	fmt.Println("Swapchain API functions demonstrated:")
	fmt.Println("  vulkan.CreateSwapchain(device, createInfo)")
	fmt.Println("  vulkan.DestroySwapchain(device, swapchain)")
	fmt.Println("  vulkan.GetSwapchainImages(device, swapchain)")
	fmt.Println("  vulkan.AcquireNextImage(device, swapchain, timeout, sem, fence)")
	fmt.Println("  vulkan.QueuePresent(queue, presentInfo)")
}
