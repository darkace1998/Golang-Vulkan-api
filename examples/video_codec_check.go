package main

import (
	"fmt"
	"log"

	vulkan "github.com/darkace1998/Golang-Vulkan-api"
)

func main() {
	fmt.Println("=== Vulkan Video Codec Support Checker ===")
	fmt.Println()

	// Create Vulkan instance
	instanceCreateInfo := &vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Video Codec Checker",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "No Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	}

	instance, err := vulkan.CreateInstance(instanceCreateInfo)
	if err != nil {
		log.Fatal("Failed to create Vulkan instance:", err)
	}
	defer vulkan.DestroyInstance(instance)

	// Enumerate physical devices
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatal("Failed to enumerate physical devices:", err)
	}

	if len(physicalDevices) == 0 {
		log.Fatal("No Vulkan-capable devices found")
	}

	fmt.Printf("Found %d Vulkan-capable device(s)\n\n", len(physicalDevices))

	// Check video codec support for each device
	for i, device := range physicalDevices {
		props := vulkan.GetPhysicalDeviceProperties(device)
		fmt.Printf("Device %d: %s\n", i, props.DeviceName)
		fmt.Printf("  API Version: %d.%d.%d\n",
			props.APIVersion.Major(),
			props.APIVersion.Minor(),
			props.APIVersion.Patch())
		fmt.Printf("  Driver Version: %d\n", props.DriverVersion)
		fmt.Printf("  Vendor ID: 0x%X\n", props.VendorID)
		fmt.Printf("  Device ID: 0x%X\n", props.DeviceID)
		fmt.Println()

		// Get supported video codecs
		supportedCodecs, err := vulkan.GetSupportedVideoCodecs(device)
		if err != nil {
			fmt.Printf("  Error checking video codec support: %v\n\n", err)
			continue
		}

		if len(supportedCodecs) == 0 {
			fmt.Println("  ❌ No video codec extensions detected")
			fmt.Println("     This GPU may not support hardware-accelerated video encoding/decoding")
			fmt.Println("     or the Vulkan drivers may need to be updated.")
		} else {
			fmt.Println("  ✅ Supported Video Codecs:")
			for _, codec := range supportedCodecs {
				fmt.Printf("     • %s\n", codec)
			}
		}

		// Check queue family support for video operations
		queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(device)
		hasVideoDecode := false
		hasVideoEncode := false

		for j, qf := range queueFamilies {
			if qf.QueueFlags&vulkan.QueueVideoDecodeBitKHR != 0 {
				hasVideoDecode = true
				fmt.Printf("  ✅ Queue Family %d: Video Decode support\n", j)
			}
			if qf.QueueFlags&vulkan.QueueVideoEncodeBitKHR != 0 {
				hasVideoEncode = true
				fmt.Printf("  ✅ Queue Family %d: Video Encode support\n", j)
			}
		}

		if !hasVideoDecode && !hasVideoEncode {
			fmt.Println("  ℹ️  No video queue families detected")
		}

		fmt.Println()
	}

	// Print summary
	fmt.Println("=== Video Codec Extension Information ===")
	fmt.Println()
	fmt.Println("Fully Supported Codecs (Encode & Decode):")
	fmt.Println("  • H.264 (AVC)")
	fmt.Println("    - Decode: VK_KHR_video_decode_h264")
	fmt.Println("    - Encode: VK_KHR_video_encode_h264")
	fmt.Println("  • H.265 (HEVC)")
	fmt.Println("    - Decode: VK_KHR_video_decode_h265")
	fmt.Println("    - Encode: VK_KHR_video_encode_h265")
	fmt.Println("  • AV1")
	fmt.Println("    - Decode: VK_KHR_video_decode_av1")
	fmt.Println("    - Encode: VK_KHR_video_encode_av1")
	fmt.Println()
	fmt.Println("Note: Hardware support depends on GPU capabilities and driver version.")
	fmt.Println("      Extension availability does not guarantee hardware acceleration.")
}
