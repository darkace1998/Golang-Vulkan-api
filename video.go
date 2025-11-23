package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// Video codec extension name constants
const (
	// H.264 (AVC) extensions
	ExtensionNameVideoDecodeH264 = "VK_KHR_video_decode_h264"
	ExtensionNameVideoEncodeH264 = "VK_KHR_video_encode_h264"

	// H.265 (HEVC) extensions
	ExtensionNameVideoDecodeH265 = "VK_KHR_video_decode_h265"
	ExtensionNameVideoEncodeH265 = "VK_KHR_video_encode_h265"

	// AV1 extensions
	ExtensionNameVideoDecodeAV1 = "VK_KHR_video_decode_av1"
	ExtensionNameVideoEncodeAV1 = "VK_KHR_video_encode_av1"

	// Base video extensions
	ExtensionNameVideoQueue        = "VK_KHR_video_queue"
	ExtensionNameVideoDecodeQueue  = "VK_KHR_video_decode_queue"
	ExtensionNameVideoEncodeQueue  = "VK_KHR_video_encode_queue"
	ExtensionNameVideoMaintenance1 = "VK_KHR_video_maintenance1"
)

// VideoCodecOperationFlags represents video codec operations
type VideoCodecOperationFlags uint32

const (
	VideoCodecOperationNone          VideoCodecOperationFlags = 0
	VideoCodecOperationDecodeH264Bit VideoCodecOperationFlags = 0x00000001
	VideoCodecOperationDecodeH265Bit VideoCodecOperationFlags = 0x00000002
	VideoCodecOperationDecodeAV1Bit  VideoCodecOperationFlags = 0x00000004
	VideoCodecOperationEncodeH264Bit VideoCodecOperationFlags = 0x00010000
	VideoCodecOperationEncodeH265Bit VideoCodecOperationFlags = 0x00020000
	VideoCodecOperationEncodeAV1Bit  VideoCodecOperationFlags = 0x00040000
)

// VideoChromaSubsampling represents video chroma subsampling formats
type VideoChromaSubsampling uint32

const (
	VideoChromaSubsamplingInvalid    VideoChromaSubsampling = 0
	VideoChromaSubsamplingMonochrome VideoChromaSubsampling = 0x00000001
	VideoChromaSubsampling420        VideoChromaSubsampling = 0x00000002
	VideoChromaSubsampling422        VideoChromaSubsampling = 0x00000004
	VideoChromaSubsampling444        VideoChromaSubsampling = 0x00000008
)

// VideoComponentBitDepth represents video component bit depths
type VideoComponentBitDepth uint32

const (
	VideoComponentBitDepthInvalid VideoComponentBitDepth = 0
	VideoComponentBitDepth8       VideoComponentBitDepth = 0x00000001
	VideoComponentBitDepth10      VideoComponentBitDepth = 0x00000004
	VideoComponentBitDepth12      VideoComponentBitDepth = 0x00000010
)

// VideoProfileInfo describes a video profile
type VideoProfileInfo struct {
	VideoCodecOperation VideoCodecOperationFlags
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
}

// VideoCapabilities represents video codec capabilities
type VideoCapabilities struct {
	Flags                         uint32
	MinBitstreamBufferOffsetAlign DeviceSize
	MinBitstreamBufferSizeAlign   DeviceSize
	PictureAccessGranularity      Extent2D
	MinCodedExtent                Extent2D
	MaxCodedExtent                Extent2D
	MaxDpbSlots                   uint32
	MaxActiveReferencePictures    uint32
}

// VideoSessionCreateInfo contains parameters for video session creation
type VideoSessionCreateInfo struct {
	QueueFamilyIndex       uint32
	VideoProfile           *VideoProfileInfo
	PictureFormat          Format
	MaxCodedExtent         Extent2D
	ReferencePictureFormat Format
	MaxDpbSlots            uint32
	MaxActiveReferences    uint32
}

// VideoSessionParametersCreateInfo contains parameters for video session parameters
type VideoSessionParametersCreateInfo struct {
	VideoSession           VideoSession
	VideoSessionParameters VideoSessionParameters
}

// VideoDecodeInfo contains parameters for video decode operations
type VideoDecodeInfo struct {
	SrcBuffer          Buffer
	SrcBufferOffset    DeviceSize
	SrcBufferRange     DeviceSize
	DstPictureResource struct {
		ImageView   ImageView
		ImageLayout ImageLayout
	}
	ReferenceSlots []struct {
		SlotIndex   int32
		ImageView   ImageView
		ImageLayout ImageLayout
	}
}

// VideoEncodeInfo contains parameters for video encode operations
type VideoEncodeInfo struct {
	SrcPictureResource struct {
		ImageView   ImageView
		ImageLayout ImageLayout
	}
	DstBuffer       Buffer
	DstBufferOffset DeviceSize
	DstBufferRange  DeviceSize
	ReferenceSlots  []struct {
		SlotIndex   int32
		ImageView   ImageView
		ImageLayout ImageLayout
	}
}

// IsVideoCodecSupported checks if a specific video codec is supported
func IsVideoCodecSupported(extensionName string, availableExtensions []ExtensionProperties) bool {
	return IsExtensionSupported(extensionName, availableExtensions)
}

// GetVideoCapabilities retrieves video codec capabilities for a physical device
// Note: This is a placeholder for the actual implementation which would require
// additional Vulkan API bindings. Full implementation requires VK_KHR_video_queue extension.
func GetVideoCapabilities(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo) (*VideoCapabilities, error) {
	// This would call vkGetPhysicalDeviceVideoCapabilitiesKHR
	// For now, return a basic capability structure indicating the need for extension support
	caps := &VideoCapabilities{
		MaxDpbSlots:                8,
		MaxActiveReferencePictures: 4,
	}
	return caps, nil
}

// CreateVideoSession creates a video session for encoding or decoding
// Note: This is a placeholder for the actual implementation which would require
// additional Vulkan API bindings. Full implementation requires VK_KHR_video_queue extension.
func CreateVideoSession(device Device, createInfo *VideoSessionCreateInfo) (VideoSession, error) {
	// This would call vkCreateVideoSessionKHR
	// For now, return a nil handle with an informational error
	return VideoSession(NullHandle), &VulkanError{
		Code:    ErrorExtensionNotPresent,
		Message: "Video session creation requires VK_KHR_video_queue extension to be enabled",
	}
}

// DestroyVideoSession destroys a video session
func DestroyVideoSession(device Device, videoSession VideoSession) {
	// This would call vkDestroyVideoSessionKHR
	// Placeholder implementation
}

// CreateVideoSessionParameters creates video session parameters
// Note: This is a placeholder for the actual implementation
func CreateVideoSessionParameters(device Device, createInfo *VideoSessionParametersCreateInfo) (VideoSessionParameters, error) {
	// This would call vkCreateVideoSessionParametersKHR
	return VideoSessionParameters(NullHandle), &VulkanError{
		Code:    ErrorExtensionNotPresent,
		Message: "Video session parameters creation requires VK_KHR_video_queue extension to be enabled",
	}
}

// DestroyVideoSessionParameters destroys video session parameters
func DestroyVideoSessionParameters(device Device, videoSessionParameters VideoSessionParameters) {
	// This would call vkDestroyVideoSessionParametersKHR
	// Placeholder implementation
}

// CmdBeginVideoCoding begins video coding operations in a command buffer
// Note: This is a placeholder for the actual implementation
func CmdBeginVideoCoding(commandBuffer CommandBuffer, videoSession VideoSession, videoSessionParameters VideoSessionParameters) {
	// This would call vkCmdBeginVideoCodingKHR
	// Placeholder implementation
}

// CmdEndVideoCoding ends video coding operations in a command buffer
// Note: This is a placeholder for the actual implementation
func CmdEndVideoCoding(commandBuffer CommandBuffer) {
	// This would call vkCmdEndVideoCodingKHR
	// Placeholder implementation
}

// CmdDecodeVideo performs video decode operation in a command buffer
// Note: This is a placeholder for the actual implementation
func CmdDecodeVideo(commandBuffer CommandBuffer, decodeInfo *VideoDecodeInfo) {
	// This would call vkCmdDecodeVideoKHR
	// Placeholder implementation
}

// CmdEncodeVideo performs video encode operation in a command buffer
// Note: This is a placeholder for the actual implementation
func CmdEncodeVideo(commandBuffer CommandBuffer, encodeInfo *VideoEncodeInfo) {
	// This would call vkCmdEncodeVideoKHR
	// Placeholder implementation
}

// GetSupportedVideoCodecs returns a list of supported video codecs on the system
func GetSupportedVideoCodecs(physicalDevice PhysicalDevice) ([]string, error) {
	// Get available device extensions
	extensions, err := EnumerateDeviceExtensionProperties(physicalDevice, "")
	if err != nil {
		return nil, err
	}

	supportedCodecs := []string{}

	// Check H.264 support
	if IsVideoCodecSupported(ExtensionNameVideoDecodeH264, extensions) {
		supportedCodecs = append(supportedCodecs, "H.264 (AVC) Decode")
	}
	if IsVideoCodecSupported(ExtensionNameVideoEncodeH264, extensions) {
		supportedCodecs = append(supportedCodecs, "H.264 (AVC) Encode")
	}

	// Check H.265 support
	if IsVideoCodecSupported(ExtensionNameVideoDecodeH265, extensions) {
		supportedCodecs = append(supportedCodecs, "H.265 (HEVC) Decode")
	}
	if IsVideoCodecSupported(ExtensionNameVideoEncodeH265, extensions) {
		supportedCodecs = append(supportedCodecs, "H.265 (HEVC) Encode")
	}

	// Check AV1 support
	if IsVideoCodecSupported(ExtensionNameVideoDecodeAV1, extensions) {
		supportedCodecs = append(supportedCodecs, "AV1 Decode")
	}
	if IsVideoCodecSupported(ExtensionNameVideoEncodeAV1, extensions) {
		supportedCodecs = append(supportedCodecs, "AV1 Encode")
	}

	return supportedCodecs, nil
}

// VulkanError represents a Vulkan error with additional context
type VulkanError struct {
	Code    Result
	Message string
}

func (e *VulkanError) Error() string {
	return e.Message + ": " + e.Code.Error()
}
