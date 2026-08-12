package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
	fmt.Println("=== Vulkan Secondary Command Buffer Example ===")

	// 1. Create Vulkan Instance
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "SecondaryCommandBufferTest",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "NoEngine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatal("Failed to create instance:", err)
	}
	defer vulkan.DestroyInstance(instance)

	// 2. Find suitable physical device and queue family
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatal("Failed to enumerate physical devices:", err)
	}
	if len(physicalDevices) == 0 {
		log.Fatal("No Vulkan physical devices found")
	}

	physicalDevice := physicalDevices[0]
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice)

	graphicsQueueFamily := ^uint32(0)
	for i, family := range queueFamilies {
		if family.QueueFlags&vulkan.QueueGraphicsBit != 0 {
			graphicsQueueFamily = uint32(i)
			break
		}
	}
	if graphicsQueueFamily == ^uint32(0) {
		log.Fatal("Failed to find graphics queue family")
	}

	// 3. Create Logical Device
	deviceCreateInfo := &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: graphicsQueueFamily,
				QueuePriorities:  []float32{1.0},
			},
		},
	}

	device, err := vulkan.CreateDevice(physicalDevice, deviceCreateInfo)
	if err != nil {
		log.Fatal("Failed to create logical device:", err)
	}
	defer vulkan.DestroyDevice(device)

	// 4. Create Command Pool
	commandPool, err := vulkan.CreateCommandPool(device, &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsQueueFamily,
	})
	if err != nil {
		log.Fatal("Failed to create command pool:", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)

	// 5. Allocate Primary Command Buffer
	primaryCmdBuffers, err := vulkan.AllocateCommandBuffers(device, &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		log.Fatal("Failed to allocate primary command buffer:", err)
	}
	primaryCmdBuffer := primaryCmdBuffers[0]

	// 6. Allocate Secondary Command Buffer
	secondaryCmdBuffers, err := vulkan.AllocateCommandBuffers(device, &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelSecondary,
		CommandBufferCount: 1,
	})
	if err != nil {
		log.Fatal("Failed to allocate secondary command buffer:", err)
	}
	secondaryCmdBuffer := secondaryCmdBuffers[0]

	// 7. Record Secondary Command Buffer
	// In a real application, you'd provide actual RenderPass and Framebuffer objects.
	// For this example, we provide a valid InheritanceInfo structure but with nil handles.
	inheritanceInfo := &vulkan.CommandBufferInheritanceInfo{
		RenderPass:           vulkan.RenderPass(nil),
		Subpass:              0,
		Framebuffer:          vulkan.Framebuffer(nil),
		OcclusionQueryEnable: false,
	}

	err = vulkan.BeginCommandBuffer(secondaryCmdBuffer, &vulkan.CommandBufferBeginInfo{
		Flags:           vulkan.CommandBufferUsageOneTimeSubmitBit | vulkan.CommandBufferUsageRenderPassContinueBit,
		InheritanceInfo: inheritanceInfo,
	})
	if err != nil {
		log.Fatal("Failed to begin secondary command buffer:", err)
	}

	// (Record rendering commands here: e.g., CmdBindPipeline, CmdDraw...)
	fmt.Println("   ✓ Secondary command buffer recorded")

	err = vulkan.EndCommandBuffer(secondaryCmdBuffer)
	if err != nil {
		log.Fatal("Failed to end secondary command buffer:", err)
	}

	// 8. Record Primary Command Buffer and Execute Secondary Command Buffer
	err = vulkan.BeginCommandBuffer(primaryCmdBuffer, &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	})
	if err != nil {
		log.Fatal("Failed to begin primary command buffer:", err)
	}

	fmt.Println("   ✓ Executing secondary command buffer from primary...")
	vulkan.CmdExecuteCommands(primaryCmdBuffer, secondaryCmdBuffers)

	err = vulkan.EndCommandBuffer(primaryCmdBuffer)
	if err != nil {
		log.Fatal("Failed to end primary command buffer:", err)
	}

	fmt.Println("\n=== Secondary Command Buffer Example Complete ===")
	fmt.Println("✓ Primary and Secondary Command Buffers allocated")
	fmt.Println("✓ CommandBufferInheritanceInfo utilized")
	fmt.Println("✓ CmdExecuteCommands executed successfully")
}
