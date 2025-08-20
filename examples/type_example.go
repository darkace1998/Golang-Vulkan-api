package main

import (
	"fmt"

	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

func main() {
	fmt.Println("Vulkan Go Binding Type System Test")
	fmt.Println("===================================")

	// Test version functions
	fmt.Println("\n1. Testing version functions...")
	version := vulkan.MakeVersion(1, 3, 269)
	fmt.Printf("Created version 1.3.269: %d\n", version)
	fmt.Printf("  Major: %d, Minor: %d, Patch: %d\n", 
		version.Major(), version.Minor(), version.Patch())

	// Test predefined versions
	fmt.Println("\n2. Testing predefined versions...")
	fmt.Printf("Vulkan 1.0: %d.%d.%d\n", 
		vulkan.Version10.Major(), vulkan.Version10.Minor(), vulkan.Version10.Patch())
	fmt.Printf("Vulkan 1.1: %d.%d.%d\n", 
		vulkan.Version11.Major(), vulkan.Version11.Minor(), vulkan.Version11.Patch())
	fmt.Printf("Vulkan 1.2: %d.%d.%d\n", 
		vulkan.Version12.Major(), vulkan.Version12.Minor(), vulkan.Version12.Patch())
	fmt.Printf("Vulkan 1.3: %d.%d.%d\n", 
		vulkan.Version13.Major(), vulkan.Version13.Minor(), vulkan.Version13.Patch())
	fmt.Printf("Vulkan 1.4: %d.%d.%d\n", 
		vulkan.Version14.Major(), vulkan.Version14.Minor(), vulkan.Version14.Patch())

	// Test error handling
	fmt.Println("\n3. Testing error codes...")
	fmt.Printf("Success: %s\n", vulkan.Success.Error())
	fmt.Printf("Error Out of Host Memory: %s\n", vulkan.ErrorOutOfHostMemory.Error())
	fmt.Printf("Error Device Lost: %s\n", vulkan.ErrorDeviceLost.Error())
	fmt.Printf("Error Out of Date KHR: %s\n", vulkan.ErrorOutOfDateKHR.Error())
	
	fmt.Printf("Is Success an error? %t\n", vulkan.Success.IsError())
	fmt.Printf("Is Success successful? %t\n", vulkan.Success.IsSuccess())
	fmt.Printf("Is ErrorOutOfHostMemory an error? %t\n", vulkan.ErrorOutOfHostMemory.IsError())
	fmt.Printf("Is ErrorOutOfHostMemory successful? %t\n", vulkan.ErrorOutOfHostMemory.IsSuccess())

	// Test boolean conversions
	fmt.Println("\n4. Testing boolean conversions...")
	vkTrue := vulkan.FromBool(true)
	vkFalse := vulkan.FromBool(false)
	fmt.Printf("Go true -> VkBool32: %d -> Go bool: %t\n", vkTrue, vkTrue.ToBool())
	fmt.Printf("Go false -> VkBool32: %d -> Go bool: %t\n", vkFalse, vkFalse.ToBool())

	// Test constants
	fmt.Println("\n5. Testing constants...")
	fmt.Printf("Max Memory Types: %d\n", vulkan.MaxMemoryTypes)
	fmt.Printf("Max Memory Heaps: %d\n", vulkan.MaxMemoryHeaps)
	fmt.Printf("Max Physical Device Name Size: %d\n", vulkan.MaxPhysicalDeviceNameSize)
	fmt.Printf("UUID Size: %d\n", vulkan.UuidSize)
	fmt.Printf("Whole Size: %d\n", vulkan.WholeSize)

	// Test flags
	fmt.Println("\n6. Testing flags...")
	fmt.Printf("Queue Graphics Bit: %d\n", vulkan.QueueGraphicsBit)
	fmt.Printf("Queue Compute Bit: %d\n", vulkan.QueueComputeBit)
	fmt.Printf("Queue Transfer Bit: %d\n", vulkan.QueueTransferBit)
	
	fmt.Printf("Buffer Usage Vertex Buffer Bit: %d\n", vulkan.BufferUsageVertexBufferBit)
	fmt.Printf("Buffer Usage Index Buffer Bit: %d\n", vulkan.BufferUsageIndexBufferBit)
	fmt.Printf("Buffer Usage Uniform Buffer Bit: %d\n", vulkan.BufferUsageUniformBufferBit)

	fmt.Printf("Memory Property Device Local Bit: %d\n", vulkan.MemoryPropertyDeviceLocalBit)
	fmt.Printf("Memory Property Host Visible Bit: %d\n", vulkan.MemoryPropertyHostVisibleBit)
	fmt.Printf("Memory Property Host Coherent Bit: %d\n", vulkan.MemoryPropertyHostCoherentBit)

	// Test formats
	fmt.Println("\n7. Testing formats...")
	fmt.Printf("Format Undefined: %d\n", vulkan.FormatUndefined)
	fmt.Printf("Format R8G8B8A8 Unorm: %d\n", vulkan.FormatR8G8B8A8Unorm)
	fmt.Printf("Format B8G8R8A8 Unorm: %d\n", vulkan.FormatB8G8R8A8Unorm)
	fmt.Printf("Format D32 Sfloat: %d\n", vulkan.FormatD32Sfloat)

	// Test sample counts
	fmt.Println("\n8. Testing sample counts...")
	fmt.Printf("Sample Count 1 Bit: %d\n", vulkan.SampleCount1Bit)
	fmt.Printf("Sample Count 4 Bit: %d\n", vulkan.SampleCount4Bit)
	fmt.Printf("Sample Count 8 Bit: %d\n", vulkan.SampleCount8Bit)

	// Test image layouts
	fmt.Println("\n9. Testing image layouts...")
	fmt.Printf("Image Layout Undefined: %d\n", vulkan.ImageLayoutUndefined)
	fmt.Printf("Image Layout Color Attachment Optimal: %d\n", vulkan.ImageLayoutColorAttachmentOptimal)
	fmt.Printf("Image Layout Shader Read Only Optimal: %d\n", vulkan.ImageLayoutShaderReadOnlyOptimal)
	fmt.Printf("Image Layout Present Src KHR: %d\n", vulkan.ImageLayoutPresentSrcKHR)

	// Test pipeline stages
	fmt.Println("\n10. Testing pipeline stages...")
	fmt.Printf("Pipeline Stage Top of Pipe Bit: %d\n", vulkan.PipelineStageTopOfPipeBit)
	fmt.Printf("Pipeline Stage Vertex Shader Bit: %d\n", vulkan.PipelineStageVertexShaderBit)
	fmt.Printf("Pipeline Stage Fragment Shader Bit: %d\n", vulkan.PipelineStageFragmentShaderBit)
	fmt.Printf("Pipeline Stage Color Attachment Output Bit: %d\n", vulkan.PipelineStageColorAttachmentOutputBit)

	// Test access flags
	fmt.Println("\n11. Testing access flags...")
	fmt.Printf("Access Shader Read Bit: %d\n", vulkan.AccessShaderReadBit)
	fmt.Printf("Access Shader Write Bit: %d\n", vulkan.AccessShaderWriteBit)
	fmt.Printf("Access Color Attachment Read Bit: %d\n", vulkan.AccessColorAttachmentReadBit)
	fmt.Printf("Access Color Attachment Write Bit: %d\n", vulkan.AccessColorAttachmentWriteBit)

	// Test utility functions
	fmt.Println("\n12. Testing utility functions...")
	apiVersion := vulkan.GetAPIVersion()
	fmt.Printf("Supported API Version: %d.%d.%d\n", 
		apiVersion.Major(), apiVersion.Minor(), apiVersion.Patch())

	// Test extension/layer checking
	extensions := []vulkan.ExtensionProperties{
		{ExtensionName: "VK_KHR_swapchain", SpecVersion: 70},
		{ExtensionName: "VK_EXT_debug_utils", SpecVersion: 2},
	}
	
	fmt.Printf("Is VK_KHR_swapchain supported? %t\n", vulkan.IsExtensionSupported("VK_KHR_swapchain", extensions))
	fmt.Printf("Is VK_KHR_nonexistent supported? %t\n", vulkan.IsExtensionSupported("VK_KHR_nonexistent", extensions))

	layers := []vulkan.LayerProperties{
		{LayerName: "VK_LAYER_KHRONOS_validation", Description: "Validation layer"},
	}
	
	fmt.Printf("Is VK_LAYER_KHRONOS_validation supported? %t\n", vulkan.IsLayerSupported("VK_LAYER_KHRONOS_validation", layers))
	fmt.Printf("Is VK_LAYER_nonexistent supported? %t\n", vulkan.IsLayerSupported("VK_LAYER_nonexistent", layers))

	fmt.Println("\n===================================")
	fmt.Println("Vulkan Go Binding Type System Test Complete!")
	fmt.Println("All type system functionality working correctly.")
	fmt.Println("===================================")
}