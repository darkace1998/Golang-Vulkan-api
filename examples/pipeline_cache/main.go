package main

import (
	"fmt"
	"log"

	"github.com/darkace1998/golang-vulkan-api"
)

func main() {
	// Create Instance
	instance, err := vulkan.CreateInstance(&vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Pipeline Cache Example",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "No Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)
	fmt.Println("Vulkan instance created successfully.")

	// Pick a physical device
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate physical devices: %v", err)
	}
	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}
	physicalDevice := physicalDevices[0]

	// Create logical device
	queuePriority := float32(1.0)
	device, err := vulkan.CreateDevice(physicalDevice, &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: 0, // Assuming 0 is a valid queue family
				QueuePriorities:  []float32{queuePriority},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create logical device: %v", err)
	}
	defer vulkan.DestroyDevice(device)
	fmt.Println("Logical device created successfully.")

	// =========================================================================
	// Pipeline Cache Demonstration
	// =========================================================================

	// 1. Create an empty pipeline cache
	fmt.Println("Creating empty pipeline cache...")
	emptyCacheInfo := &vulkan.PipelineCacheCreateInfo{
		Flags: 0,
	}
	emptyCache, err := vulkan.CreatePipelineCache(device, emptyCacheInfo)
	if err != nil {
		log.Fatalf("Failed to create empty pipeline cache: %v", err)
	}
	defer vulkan.DestroyPipelineCache(device, emptyCache)

	// 2. Retrieve data from the pipeline cache (simulating saving to disk)
	// In a real application, you would create pipelines using this cache,
	// which populates it, and then you would retrieve the data to save to disk.
	fmt.Println("Retrieving pipeline cache data...")
	cacheData, err := vulkan.GetPipelineCacheData(device, emptyCache)
	if err != nil {
		log.Fatalf("Failed to get pipeline cache data: %v", err)
	}
	fmt.Printf("Retrieved %d bytes of pipeline cache data.\n", len(cacheData))

	// 3. Create a new pipeline cache using the retrieved data (simulating loading from disk)
	fmt.Println("Creating new pipeline cache from retrieved data...")
	loadedCacheInfo := &vulkan.PipelineCacheCreateInfo{
		Flags:       0,
		InitialData: cacheData,
	}
	loadedCache, err := vulkan.CreatePipelineCache(device, loadedCacheInfo)
	if err != nil {
		log.Fatalf("Failed to create loaded pipeline cache: %v", err)
	}
	defer vulkan.DestroyPipelineCache(device, loadedCache)
	fmt.Println("Loaded pipeline cache created successfully.")

	// 4. Merge pipeline caches
	fmt.Println("Merging pipeline caches...")

	// Create another empty cache to use as the destination
	dstCacheInfo := &vulkan.PipelineCacheCreateInfo{
		Flags: 0,
	}
	dstCache, err := vulkan.CreatePipelineCache(device, dstCacheInfo)
	if err != nil {
		log.Fatalf("Failed to create destination pipeline cache: %v", err)
	}
	defer vulkan.DestroyPipelineCache(device, dstCache)

	// Merge loadedCache into dstCache
	srcCaches := []vulkan.PipelineCache{loadedCache}
	err = vulkan.MergePipelineCaches(device, dstCache, srcCaches)
	if err != nil {
		log.Fatalf("Failed to merge pipeline caches: %v", err)
	}
	fmt.Println("Pipeline caches merged successfully.")

	fmt.Println("Pipeline cache example completed successfully.")
}
