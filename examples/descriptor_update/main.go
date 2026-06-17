package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
	fmt.Println("=== Vulkan Descriptor Update Example ===")

	// 1. Initialize Vulkan
	appInfo := &vulkan.ApplicationInfo{
		ApplicationName:    "Descriptor Update Example",
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

	// 2. Create Descriptor Set Layout
	// We'll create a layout with:
	// - Binding 0: Uniform Buffer
	// - Binding 1: Combined Image Sampler
	layoutBindings := []vulkan.DescriptorSetLayoutBinding{
		{
			Binding:         0,
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
			StageFlags:      vulkan.ShaderStageVertexBit,
		},
		{
			Binding:         1,
			DescriptorType:  vulkan.DescriptorTypeCombinedImageSampler,
			DescriptorCount: 1,
			StageFlags:      vulkan.ShaderStageFragmentBit,
		},
	}

	layoutInfo := &vulkan.DescriptorSetLayoutCreateInfo{
		Bindings: layoutBindings,
	}

	layout, err := vulkan.CreateDescriptorSetLayout(device, layoutInfo)
	if err != nil {
		log.Fatalf("Failed to create descriptor set layout: %v", err)
	}
	defer vulkan.DestroyDescriptorSetLayout(device, layout)

	// 3. Create Descriptor Pool
	poolSizes := []vulkan.DescriptorPoolSize{
		{
			Type:            vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
		},
		{
			Type:            vulkan.DescriptorTypeCombinedImageSampler,
			DescriptorCount: 1,
		},
	}

	poolInfo := &vulkan.DescriptorPoolCreateInfo{
		MaxSets:   1,
		PoolSizes: poolSizes,
	}

	descriptorPool, err := vulkan.CreateDescriptorPool(device, poolInfo)
	if err != nil {
		log.Fatalf("Failed to create descriptor pool: %v", err)
	}
	defer vulkan.DestroyDescriptorPool(device, descriptorPool)

	// 4. Allocate Descriptor Set
	allocInfo := &vulkan.DescriptorSetAllocateInfo{
		DescriptorPool: descriptorPool,
		SetLayouts:     []vulkan.DescriptorSetLayout{layout},
	}

	descriptorSets, err := vulkan.AllocateDescriptorSets(device, allocInfo)
	if err != nil || len(descriptorSets) == 0 {
		log.Fatalf("Failed to allocate descriptor sets: %v", err)
	}
	descriptorSet := descriptorSets[0]

	// 5. Create Mock Resources for Update
	// Dummy Buffer
	bufferInfo := &vulkan.BufferCreateInfo{
		Size:        256,
		Usage:       vulkan.BufferUsageUniformBufferBit,
		SharingMode: vulkan.SharingModeExclusive,
	}
	uniformBuffer, err := vulkan.CreateBuffer(device, bufferInfo)
	if err != nil {
		log.Fatalf("Failed to create buffer: %v", err)
	}
	defer vulkan.DestroyBuffer(device, uniformBuffer)

	// Dummy Sampler
	samplerInfo := &vulkan.SamplerCreateInfo{
		MagFilter:    vulkan.FilterLinear,
		MinFilter:    vulkan.FilterLinear,
		AddressModeU: vulkan.SamplerAddressModeRepeat,
		AddressModeV: vulkan.SamplerAddressModeRepeat,
		AddressModeW: vulkan.SamplerAddressModeRepeat,
	}
	sampler, err := vulkan.CreateSampler(device, samplerInfo)
	if err != nil {
		log.Fatalf("Failed to create sampler: %v", err)
	}
	defer vulkan.DestroySampler(device, sampler)

	fmt.Println("Preparing descriptor set update structures...")

	// 6. Update Descriptor Set
	bufferWriteInfo := []vulkan.DescriptorBufferInfo{
		{
			Buffer: uniformBuffer,
			Offset: 0,
			Range:  256, // or vulkan.WholeSize
		},
	}

	imageWriteInfo := []vulkan.DescriptorImageInfo{
		{
			Sampler: sampler,
			// Normally you'd set a valid ImageView here. We use a mock 0 handle for demonstration.
			// ImageView: imageView,
			ImageLayout: vulkan.ImageLayoutShaderReadOnlyOptimal,
		},
	}

	descriptorWrites := []vulkan.WriteDescriptorSet{
		{
			DstSet:          descriptorSet,
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			BufferInfo:      bufferWriteInfo,
		},
		{
			DstSet:          descriptorSet,
			DstBinding:      1,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vulkan.DescriptorTypeCombinedImageSampler,
			ImageInfo:       imageWriteInfo,
		},
	}

	// We prepare the structures and we can update them
	// vulkan.UpdateDescriptorSets(device, descriptorWrites, nil)
	_ = descriptorWrites

	fmt.Println("Descriptor set update structures created successfully.")
	fmt.Println("Use `vulkan.UpdateDescriptorSets(device, descriptorWrites, nil)` to apply updates.")

	fmt.Println("=== Example Completed Successfully ===")
}
