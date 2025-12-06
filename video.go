package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>

// Function pointers for video KHR extension functions
// These need to be loaded dynamically at runtime.
//
// IMPORTANT: These are global static pointers and NOT thread-safe during loading.
// LoadVideoInstanceFunctions/LoadVideoDeviceFunctions must be called from a single
// thread during initialization before any concurrent video API usage.
//
// NOTE: Only one Vulkan instance/device with video support is supported at a time.
// Calling the load functions multiple times will overwrite previous function pointers.
// Per-device function pointers are not currently supported.
static PFN_vkGetPhysicalDeviceVideoCapabilitiesKHR pfn_vkGetPhysicalDeviceVideoCapabilitiesKHR = NULL;
static PFN_vkCreateVideoSessionKHR pfn_vkCreateVideoSessionKHR = NULL;
static PFN_vkDestroyVideoSessionKHR pfn_vkDestroyVideoSessionKHR = NULL;
static PFN_vkGetVideoSessionMemoryRequirementsKHR pfn_vkGetVideoSessionMemoryRequirementsKHR = NULL;
static PFN_vkBindVideoSessionMemoryKHR pfn_vkBindVideoSessionMemoryKHR = NULL;
static PFN_vkCreateVideoSessionParametersKHR pfn_vkCreateVideoSessionParametersKHR = NULL;
static PFN_vkDestroyVideoSessionParametersKHR pfn_vkDestroyVideoSessionParametersKHR = NULL;
static PFN_vkCmdBeginVideoCodingKHR pfn_vkCmdBeginVideoCodingKHR = NULL;
static PFN_vkCmdEndVideoCodingKHR pfn_vkCmdEndVideoCodingKHR = NULL;
static PFN_vkCmdControlVideoCodingKHR pfn_vkCmdControlVideoCodingKHR = NULL;
static PFN_vkCmdDecodeVideoKHR pfn_vkCmdDecodeVideoKHR = NULL;
static PFN_vkCmdEncodeVideoKHR pfn_vkCmdEncodeVideoKHR = NULL;

// Helper functions to load extension functions
static int loadVideoInstanceFunctions(VkInstance instance) {
    if (instance == VK_NULL_HANDLE) {
        return 0;
    }
    pfn_vkGetPhysicalDeviceVideoCapabilitiesKHR = (PFN_vkGetPhysicalDeviceVideoCapabilitiesKHR)
        vkGetInstanceProcAddr(instance, "vkGetPhysicalDeviceVideoCapabilitiesKHR");
    return pfn_vkGetPhysicalDeviceVideoCapabilitiesKHR != NULL;
}

static int loadVideoDeviceFunctions(VkDevice device) {
    if (device == VK_NULL_HANDLE) {
        return 0;
    }
    pfn_vkCreateVideoSessionKHR = (PFN_vkCreateVideoSessionKHR)
        vkGetDeviceProcAddr(device, "vkCreateVideoSessionKHR");
    pfn_vkDestroyVideoSessionKHR = (PFN_vkDestroyVideoSessionKHR)
        vkGetDeviceProcAddr(device, "vkDestroyVideoSessionKHR");
    pfn_vkGetVideoSessionMemoryRequirementsKHR = (PFN_vkGetVideoSessionMemoryRequirementsKHR)
        vkGetDeviceProcAddr(device, "vkGetVideoSessionMemoryRequirementsKHR");
    pfn_vkBindVideoSessionMemoryKHR = (PFN_vkBindVideoSessionMemoryKHR)
        vkGetDeviceProcAddr(device, "vkBindVideoSessionMemoryKHR");
    pfn_vkCreateVideoSessionParametersKHR = (PFN_vkCreateVideoSessionParametersKHR)
        vkGetDeviceProcAddr(device, "vkCreateVideoSessionParametersKHR");
    pfn_vkDestroyVideoSessionParametersKHR = (PFN_vkDestroyVideoSessionParametersKHR)
        vkGetDeviceProcAddr(device, "vkDestroyVideoSessionParametersKHR");
    pfn_vkCmdBeginVideoCodingKHR = (PFN_vkCmdBeginVideoCodingKHR)
        vkGetDeviceProcAddr(device, "vkCmdBeginVideoCodingKHR");
    pfn_vkCmdEndVideoCodingKHR = (PFN_vkCmdEndVideoCodingKHR)
        vkGetDeviceProcAddr(device, "vkCmdEndVideoCodingKHR");
    pfn_vkCmdControlVideoCodingKHR = (PFN_vkCmdControlVideoCodingKHR)
        vkGetDeviceProcAddr(device, "vkCmdControlVideoCodingKHR");
    pfn_vkCmdDecodeVideoKHR = (PFN_vkCmdDecodeVideoKHR)
        vkGetDeviceProcAddr(device, "vkCmdDecodeVideoKHR");
    pfn_vkCmdEncodeVideoKHR = (PFN_vkCmdEncodeVideoKHR)
        vkGetDeviceProcAddr(device, "vkCmdEncodeVideoKHR");

    // Validate ALL loaded function pointers - returns false if any function failed to load.
    // All functions are considered critical for proper video support.
    return pfn_vkCreateVideoSessionKHR != NULL &&
           pfn_vkDestroyVideoSessionKHR != NULL &&
           pfn_vkGetVideoSessionMemoryRequirementsKHR != NULL &&
           pfn_vkBindVideoSessionMemoryKHR != NULL &&
           pfn_vkCreateVideoSessionParametersKHR != NULL &&
           pfn_vkDestroyVideoSessionParametersKHR != NULL &&
           pfn_vkCmdBeginVideoCodingKHR != NULL &&
           pfn_vkCmdEndVideoCodingKHR != NULL &&
           pfn_vkCmdControlVideoCodingKHR != NULL &&
           pfn_vkCmdDecodeVideoKHR != NULL &&
           pfn_vkCmdEncodeVideoKHR != NULL;
}

// Wrapper functions that use the dynamically loaded function pointers
static VkResult call_vkGetPhysicalDeviceVideoCapabilitiesKHR(
    VkPhysicalDevice physicalDevice,
    const VkVideoProfileInfoKHR* pVideoProfile,
    VkVideoCapabilitiesKHR* pCapabilities) {
    if (pfn_vkGetPhysicalDeviceVideoCapabilitiesKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkGetPhysicalDeviceVideoCapabilitiesKHR(physicalDevice, pVideoProfile, pCapabilities);
}

static VkResult call_vkCreateVideoSessionKHR(
    VkDevice device,
    const VkVideoSessionCreateInfoKHR* pCreateInfo,
    const VkAllocationCallbacks* pAllocator,
    VkVideoSessionKHR* pVideoSession) {
    if (pfn_vkCreateVideoSessionKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkCreateVideoSessionKHR(device, pCreateInfo, pAllocator, pVideoSession);
}

static void call_vkDestroyVideoSessionKHR(
    VkDevice device,
    VkVideoSessionKHR videoSession,
    const VkAllocationCallbacks* pAllocator) {
    if (pfn_vkDestroyVideoSessionKHR != NULL) {
        pfn_vkDestroyVideoSessionKHR(device, videoSession, pAllocator);
    }
}

static VkResult call_vkGetVideoSessionMemoryRequirementsKHR(
    VkDevice device,
    VkVideoSessionKHR videoSession,
    uint32_t* pMemoryRequirementsCount,
    VkVideoSessionMemoryRequirementsKHR* pMemoryRequirements) {
    if (pfn_vkGetVideoSessionMemoryRequirementsKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkGetVideoSessionMemoryRequirementsKHR(device, videoSession, pMemoryRequirementsCount, pMemoryRequirements);
}

static VkResult call_vkBindVideoSessionMemoryKHR(
    VkDevice device,
    VkVideoSessionKHR videoSession,
    uint32_t bindSessionMemoryInfoCount,
    const VkBindVideoSessionMemoryInfoKHR* pBindSessionMemoryInfos) {
    if (pfn_vkBindVideoSessionMemoryKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkBindVideoSessionMemoryKHR(device, videoSession, bindSessionMemoryInfoCount, pBindSessionMemoryInfos);
}

static VkResult call_vkCreateVideoSessionParametersKHR(
    VkDevice device,
    const VkVideoSessionParametersCreateInfoKHR* pCreateInfo,
    const VkAllocationCallbacks* pAllocator,
    VkVideoSessionParametersKHR* pVideoSessionParameters) {
    if (pfn_vkCreateVideoSessionParametersKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkCreateVideoSessionParametersKHR(device, pCreateInfo, pAllocator, pVideoSessionParameters);
}

static void call_vkDestroyVideoSessionParametersKHR(
    VkDevice device,
    VkVideoSessionParametersKHR videoSessionParameters,
    const VkAllocationCallbacks* pAllocator) {
    if (pfn_vkDestroyVideoSessionParametersKHR != NULL) {
        pfn_vkDestroyVideoSessionParametersKHR(device, videoSessionParameters, pAllocator);
    }
}

// Command buffer wrapper functions return 1 on success, 0 if function pointer is NULL.
// Callers should check return value to detect if LoadVideoDeviceFunctions was not called.
static int call_vkCmdBeginVideoCodingKHR(
    VkCommandBuffer commandBuffer,
    const VkVideoBeginCodingInfoKHR* pBeginInfo) {
    if (pfn_vkCmdBeginVideoCodingKHR == NULL) {
        return 0;
    }
    pfn_vkCmdBeginVideoCodingKHR(commandBuffer, pBeginInfo);
    return 1;
}

static int call_vkCmdEndVideoCodingKHR(
    VkCommandBuffer commandBuffer,
    const VkVideoEndCodingInfoKHR* pEndCodingInfo) {
    if (pfn_vkCmdEndVideoCodingKHR == NULL) {
        return 0;
    }
    pfn_vkCmdEndVideoCodingKHR(commandBuffer, pEndCodingInfo);
    return 1;
}

static int call_vkCmdControlVideoCodingKHR(
    VkCommandBuffer commandBuffer,
    const VkVideoCodingControlInfoKHR* pCodingControlInfo) {
    if (pfn_vkCmdControlVideoCodingKHR == NULL) {
        return 0;
    }
    pfn_vkCmdControlVideoCodingKHR(commandBuffer, pCodingControlInfo);
    return 1;
}

static int call_vkCmdDecodeVideoKHR(
    VkCommandBuffer commandBuffer,
    const VkVideoDecodeInfoKHR* pDecodeInfo) {
    if (pfn_vkCmdDecodeVideoKHR == NULL) {
        return 0;
    }
    pfn_vkCmdDecodeVideoKHR(commandBuffer, pDecodeInfo);
    return 1;
}

static int call_vkCmdEncodeVideoKHR(
    VkCommandBuffer commandBuffer,
    const VkVideoEncodeInfoKHR* pEncodeInfo) {
    if (pfn_vkCmdEncodeVideoKHR == NULL) {
        return 0;
    }
    pfn_vkCmdEncodeVideoKHR(commandBuffer, pEncodeInfo);
    return 1;
}
*/
import "C"

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

// VideoPictureResource contains video picture resource information
type VideoPictureResource struct {
	ImageView      ImageView
	ImageLayout    ImageLayout
	CodedOffset    Offset2D
	CodedExtent    Extent2D
	BaseArrayLayer uint32
}

// VideoDecodeInfo contains parameters for video decode operations
type VideoDecodeInfo struct {
	SrcBuffer          Buffer
	SrcBufferOffset    DeviceSize
	SrcBufferRange     DeviceSize
	DstPictureResource VideoPictureResource
	ReferenceSlots     []struct {
		SlotIndex   int32
		ImageView   ImageView
		ImageLayout ImageLayout
	}
}

// VideoEncodeInfo contains parameters for video encode operations
type VideoEncodeInfo struct {
	SrcPictureResource VideoPictureResource
	DstBuffer          Buffer
	DstBufferOffset    DeviceSize
	DstBufferRange     DeviceSize
	ReferenceSlots     []struct {
		SlotIndex   int32
		ImageView   ImageView
		ImageLayout ImageLayout
	}
}

// LoadVideoInstanceFunctions loads video extension functions that require a Vulkan instance.
//
// This function MUST be called after creating a Vulkan instance and before using any video-related
// functionality. If this function is not called, all video API calls will fail.
//
// IMPORTANT: This function is NOT thread-safe. It must be called from a single thread during
// initialization before any concurrent video API usage. Only one instance is supported at a time;
// calling this function again will overwrite previously loaded function pointers.
//
// Returns false if the video extension functions could not be loaded (e.g., if the Vulkan
// implementation does not support the VK_KHR_video_queue extension).
//
// Example usage:
//
//	instance, _ := vulkan.CreateInstance(...)
//	if !vulkan.LoadVideoInstanceFunctions(instance) {
//	    log.Fatal("Failed to load video instance functions - video extensions not supported")
//	}
func LoadVideoInstanceFunctions(instance Instance) bool {
	return C.loadVideoInstanceFunctions(C.VkInstance(instance)) != 0
}

// LoadVideoDeviceFunctions loads video extension functions that require a Vulkan device.
//
// This function MUST be called after creating a logical device and before using any video-related
// functionality. If this function is not called, all video API calls will fail.
//
// IMPORTANT: This function is NOT thread-safe. It must be called from a single thread during
// initialization before any concurrent video API usage. Only one device is supported at a time;
// calling this function again will overwrite previously loaded function pointers.
//
// Returns false if any video extension function could not be loaded. This indicates the device
// does not fully support the VK_KHR_video_queue extension.
//
// Example usage:
//
//	device, _ := vulkan.CreateDevice(...)
//	if !vulkan.LoadVideoDeviceFunctions(device) {
//	    log.Fatal("Failed to load video device functions - video extensions not supported")
//	}
func LoadVideoDeviceFunctions(device Device) bool {
	return C.loadVideoDeviceFunctions(C.VkDevice(device)) != 0
}

// GetVideoCapabilities retrieves video codec capabilities for a physical device
func GetVideoCapabilities(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo) (*VideoCapabilities, error) {
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if videoProfile == nil {
		return nil, NewValidationError("videoProfile", "cannot be nil")
	}

	// Create C structures for video profile
	var cVideoProfile C.VkVideoProfileInfoKHR
	cVideoProfile.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_INFO_KHR
	cVideoProfile.pNext = nil
	cVideoProfile.videoCodecOperation = C.VkVideoCodecOperationFlagBitsKHR(videoProfile.VideoCodecOperation)
	cVideoProfile.chromaSubsampling = C.VkVideoChromaSubsamplingFlagsKHR(videoProfile.ChromaSubsampling)
	cVideoProfile.lumaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(videoProfile.LumaBitDepth)
	cVideoProfile.chromaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(videoProfile.ChromaBitDepth)

	var cCaps C.VkVideoCapabilitiesKHR
	cCaps.sType = C.VK_STRUCTURE_TYPE_VIDEO_CAPABILITIES_KHR
	cCaps.pNext = nil

	result := Result(C.call_vkGetPhysicalDeviceVideoCapabilitiesKHR(
		C.VkPhysicalDevice(physicalDevice),
		&cVideoProfile,
		&cCaps,
	))

	if result != Success {
		return nil, NewVulkanError(result, "GetVideoCapabilities", "failed to get video capabilities")
	}

	caps := &VideoCapabilities{
		Flags:                         uint32(cCaps.flags),
		MinBitstreamBufferOffsetAlign: DeviceSize(cCaps.minBitstreamBufferOffsetAlignment),
		MinBitstreamBufferSizeAlign:   DeviceSize(cCaps.minBitstreamBufferSizeAlignment),
		PictureAccessGranularity: Extent2D{
			Width:  uint32(cCaps.pictureAccessGranularity.width),
			Height: uint32(cCaps.pictureAccessGranularity.height),
		},
		MinCodedExtent: Extent2D{
			Width:  uint32(cCaps.minCodedExtent.width),
			Height: uint32(cCaps.minCodedExtent.height),
		},
		MaxCodedExtent: Extent2D{
			Width:  uint32(cCaps.maxCodedExtent.width),
			Height: uint32(cCaps.maxCodedExtent.height),
		},
		MaxDpbSlots:                uint32(cCaps.maxDpbSlots),
		MaxActiveReferencePictures: uint32(cCaps.maxActiveReferencePictures),
	}

	return caps, nil
}

// CreateVideoSession creates a video session for encoding or decoding
func CreateVideoSession(device Device, createInfo *VideoSessionCreateInfo) (VideoSession, error) {
	if device == nil {
		return VideoSession(NullHandle), NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.VideoProfile == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo.VideoProfile", "cannot be nil")
	}

	// Create C video profile structure
	var cVideoProfile C.VkVideoProfileInfoKHR
	cVideoProfile.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_INFO_KHR
	cVideoProfile.pNext = nil
	cVideoProfile.videoCodecOperation = C.VkVideoCodecOperationFlagBitsKHR(createInfo.VideoProfile.VideoCodecOperation)
	cVideoProfile.chromaSubsampling = C.VkVideoChromaSubsamplingFlagsKHR(createInfo.VideoProfile.ChromaSubsampling)
	cVideoProfile.lumaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(createInfo.VideoProfile.LumaBitDepth)
	cVideoProfile.chromaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(createInfo.VideoProfile.ChromaBitDepth)

	// Create C video session create info
	var cCreateInfo C.VkVideoSessionCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_SESSION_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.queueFamilyIndex = C.uint32_t(createInfo.QueueFamilyIndex)
	cCreateInfo.pVideoProfile = &cVideoProfile
	cCreateInfo.pictureFormat = C.VkFormat(createInfo.PictureFormat)
	cCreateInfo.maxCodedExtent.width = C.uint32_t(createInfo.MaxCodedExtent.Width)
	cCreateInfo.maxCodedExtent.height = C.uint32_t(createInfo.MaxCodedExtent.Height)
	cCreateInfo.referencePictureFormat = C.VkFormat(createInfo.ReferencePictureFormat)
	cCreateInfo.maxDpbSlots = C.uint32_t(createInfo.MaxDpbSlots)
	cCreateInfo.maxActiveReferencePictures = C.uint32_t(createInfo.MaxActiveReferences)

	var videoSession C.VkVideoSessionKHR
	result := Result(C.call_vkCreateVideoSessionKHR(
		C.VkDevice(device),
		&cCreateInfo,
		nil,
		&videoSession,
	))

	if result != Success {
		return VideoSession(NullHandle), NewVulkanError(result, "CreateVideoSession", "failed to create video session")
	}

	return VideoSession(videoSession), nil
}

// DestroyVideoSession destroys a video session
func DestroyVideoSession(device Device, videoSession VideoSession) {
	if device == nil || videoSession == VideoSession(NullHandle) {
		return
	}
	C.call_vkDestroyVideoSessionKHR(C.VkDevice(device), C.VkVideoSessionKHR(videoSession), nil)
}

// GetVideoSessionMemoryRequirements gets memory requirements for a video session
func GetVideoSessionMemoryRequirements(device Device, videoSession VideoSession) ([]MemoryRequirements, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if videoSession == VideoSession(NullHandle) {
		return nil, NewValidationError("videoSession", "cannot be null")
	}

	var memReqCount C.uint32_t
	result := Result(C.call_vkGetVideoSessionMemoryRequirementsKHR(
		C.VkDevice(device),
		C.VkVideoSessionKHR(videoSession),
		&memReqCount,
		nil,
	))

	if result != Success {
		return nil, NewVulkanError(result, "GetVideoSessionMemoryRequirements", "failed to get memory requirements count")
	}

	if memReqCount == 0 {
		return []MemoryRequirements{}, nil
	}

	cMemReqs := make([]C.VkVideoSessionMemoryRequirementsKHR, memReqCount)
	for i := range cMemReqs {
		cMemReqs[i].sType = C.VK_STRUCTURE_TYPE_VIDEO_SESSION_MEMORY_REQUIREMENTS_KHR
		cMemReqs[i].pNext = nil
	}

	result = Result(C.call_vkGetVideoSessionMemoryRequirementsKHR(
		C.VkDevice(device),
		C.VkVideoSessionKHR(videoSession),
		&memReqCount,
		&cMemReqs[0],
	))

	if result != Success {
		return nil, NewVulkanError(result, "GetVideoSessionMemoryRequirements", "failed to get memory requirements")
	}

	memReqs := make([]MemoryRequirements, memReqCount)
	for i := range memReqs {
		memReqs[i] = MemoryRequirements{
			Size:           DeviceSize(cMemReqs[i].memoryRequirements.size),
			Alignment:      DeviceSize(cMemReqs[i].memoryRequirements.alignment),
			MemoryTypeBits: uint32(cMemReqs[i].memoryRequirements.memoryTypeBits),
		}
	}

	return memReqs, nil
}

// BindVideoSessionMemory binds memory to a video session
func BindVideoSessionMemory(device Device, videoSession VideoSession, bindInfos []VideoBindMemoryInfo) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if videoSession == VideoSession(NullHandle) {
		return NewValidationError("videoSession", "cannot be null")
	}
	if len(bindInfos) == 0 {
		return NewValidationError("bindInfos", "must have at least one bind info")
	}

	cBindInfos := make([]C.VkBindVideoSessionMemoryInfoKHR, len(bindInfos))
	for i, info := range bindInfos {
		cBindInfos[i].sType = C.VK_STRUCTURE_TYPE_BIND_VIDEO_SESSION_MEMORY_INFO_KHR
		cBindInfos[i].pNext = nil
		cBindInfos[i].memoryBindIndex = C.uint32_t(info.MemoryBindIndex)
		cBindInfos[i].memory = C.VkDeviceMemory(info.Memory)
		cBindInfos[i].memoryOffset = C.VkDeviceSize(info.MemoryOffset)
		cBindInfos[i].memorySize = C.VkDeviceSize(info.MemorySize)
	}

	result := Result(C.call_vkBindVideoSessionMemoryKHR(
		C.VkDevice(device),
		C.VkVideoSessionKHR(videoSession),
		C.uint32_t(len(bindInfos)),
		&cBindInfos[0],
	))

	if result != Success {
		return NewVulkanError(result, "BindVideoSessionMemory", "failed to bind video session memory")
	}

	return nil
}

// VideoBindMemoryInfo contains video session memory binding information
type VideoBindMemoryInfo struct {
	MemoryBindIndex uint32
	Memory          DeviceMemory
	MemoryOffset    DeviceSize
	MemorySize      DeviceSize
}

// CreateVideoSessionParameters creates video session parameters
func CreateVideoSessionParameters(device Device, createInfo *VideoSessionParametersCreateInfo) (VideoSessionParameters, error) {
	if device == nil {
		return VideoSessionParameters(NullHandle), NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return VideoSessionParameters(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	var cCreateInfo C.VkVideoSessionParametersCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_SESSION_PARAMETERS_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.videoSessionParametersTemplate = C.VkVideoSessionParametersKHR(createInfo.VideoSessionParameters)
	cCreateInfo.videoSession = C.VkVideoSessionKHR(createInfo.VideoSession)

	var videoSessionParams C.VkVideoSessionParametersKHR
	result := Result(C.call_vkCreateVideoSessionParametersKHR(
		C.VkDevice(device),
		&cCreateInfo,
		nil,
		&videoSessionParams,
	))

	if result != Success {
		return VideoSessionParameters(NullHandle), NewVulkanError(result, "CreateVideoSessionParameters", "failed to create video session parameters")
	}

	return VideoSessionParameters(videoSessionParams), nil
}

// DestroyVideoSessionParameters destroys video session parameters
func DestroyVideoSessionParameters(device Device, videoSessionParameters VideoSessionParameters) {
	if device == nil || videoSessionParameters == VideoSessionParameters(NullHandle) {
		return
	}
	C.call_vkDestroyVideoSessionParametersKHR(C.VkDevice(device), C.VkVideoSessionParametersKHR(videoSessionParameters), nil)
}

// VideoCodingControlInfo contains video coding control information
type VideoCodingControlInfo struct {
	Flags uint32
}

// CmdBeginVideoCoding begins video coding operations in a command buffer.
// Returns an error if LoadVideoDeviceFunctions was not called or video extensions are not supported.
func CmdBeginVideoCoding(commandBuffer CommandBuffer, beginInfo *VideoBeginCodingInfo) error {
	if commandBuffer == nil {
		return NewValidationError("commandBuffer", "cannot be nil")
	}
	if beginInfo == nil {
		return NewValidationError("beginInfo", "cannot be nil")
	}

	var cBeginInfo C.VkVideoBeginCodingInfoKHR
	cBeginInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_BEGIN_CODING_INFO_KHR
	cBeginInfo.pNext = nil
	cBeginInfo.flags = 0
	cBeginInfo.videoSession = C.VkVideoSessionKHR(beginInfo.VideoSession)
	cBeginInfo.videoSessionParameters = C.VkVideoSessionParametersKHR(beginInfo.VideoSessionParameters)
	cBeginInfo.referenceSlotCount = 0
	cBeginInfo.pReferenceSlots = nil

	if C.call_vkCmdBeginVideoCodingKHR(C.VkCommandBuffer(commandBuffer), &cBeginInfo) == 0 {
		return NewVulkanError(ErrorExtensionNotPresent, "CmdBeginVideoCoding", "video extension not loaded - call LoadVideoDeviceFunctions first")
	}
	return nil
}

// VideoBeginCodingInfo contains video begin coding information
type VideoBeginCodingInfo struct {
	VideoSession           VideoSession
	VideoSessionParameters VideoSessionParameters
}

// CmdEndVideoCoding ends video coding operations in a command buffer.
// Returns an error if LoadVideoDeviceFunctions was not called or video extensions are not supported.
func CmdEndVideoCoding(commandBuffer CommandBuffer) error {
	if commandBuffer == nil {
		return NewValidationError("commandBuffer", "cannot be nil")
	}

	var cEndInfo C.VkVideoEndCodingInfoKHR
	cEndInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_END_CODING_INFO_KHR
	cEndInfo.pNext = nil
	cEndInfo.flags = 0

	if C.call_vkCmdEndVideoCodingKHR(C.VkCommandBuffer(commandBuffer), &cEndInfo) == 0 {
		return NewVulkanError(ErrorExtensionNotPresent, "CmdEndVideoCoding", "video extension not loaded - call LoadVideoDeviceFunctions first")
	}
	return nil
}

// CmdControlVideoCoding controls video coding operations.
// Returns an error if LoadVideoDeviceFunctions was not called or video extensions are not supported.
func CmdControlVideoCoding(commandBuffer CommandBuffer, controlInfo *VideoCodingControlInfo) error {
	if commandBuffer == nil {
		return NewValidationError("commandBuffer", "cannot be nil")
	}
	if controlInfo == nil {
		return NewValidationError("controlInfo", "cannot be nil")
	}

	var cControlInfo C.VkVideoCodingControlInfoKHR
	cControlInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_CODING_CONTROL_INFO_KHR
	cControlInfo.pNext = nil
	cControlInfo.flags = C.VkVideoCodingControlFlagsKHR(controlInfo.Flags)

	if C.call_vkCmdControlVideoCodingKHR(C.VkCommandBuffer(commandBuffer), &cControlInfo) == 0 {
		return NewVulkanError(ErrorExtensionNotPresent, "CmdControlVideoCoding", "video extension not loaded - call LoadVideoDeviceFunctions first")
	}
	return nil
}

// CmdDecodeVideo performs video decode operation in a command buffer.
// Returns an error if LoadVideoDeviceFunctions was not called or video extensions are not supported.
func CmdDecodeVideo(commandBuffer CommandBuffer, decodeInfo *VideoDecodeInfo) error {
	if commandBuffer == nil {
		return NewValidationError("commandBuffer", "cannot be nil")
	}
	if decodeInfo == nil {
		return NewValidationError("decodeInfo", "cannot be nil")
	}

	var cDecodeInfo C.VkVideoDecodeInfoKHR
	cDecodeInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_INFO_KHR
	cDecodeInfo.pNext = nil
	cDecodeInfo.flags = 0
	cDecodeInfo.srcBuffer = C.VkBuffer(decodeInfo.SrcBuffer)
	cDecodeInfo.srcBufferOffset = C.VkDeviceSize(decodeInfo.SrcBufferOffset)
	cDecodeInfo.srcBufferRange = C.VkDeviceSize(decodeInfo.SrcBufferRange)

	// Setup destination picture resource
	var cDstPictureResource C.VkVideoPictureResourceInfoKHR
	cDstPictureResource.sType = C.VK_STRUCTURE_TYPE_VIDEO_PICTURE_RESOURCE_INFO_KHR
	cDstPictureResource.pNext = nil
	cDstPictureResource.codedOffset.x = C.int32_t(decodeInfo.DstPictureResource.CodedOffset.X)
	cDstPictureResource.codedOffset.y = C.int32_t(decodeInfo.DstPictureResource.CodedOffset.Y)
	cDstPictureResource.codedExtent.width = C.uint32_t(decodeInfo.DstPictureResource.CodedExtent.Width)
	cDstPictureResource.codedExtent.height = C.uint32_t(decodeInfo.DstPictureResource.CodedExtent.Height)
	cDstPictureResource.baseArrayLayer = C.uint32_t(decodeInfo.DstPictureResource.BaseArrayLayer)
	cDstPictureResource.imageViewBinding = C.VkImageView(decodeInfo.DstPictureResource.ImageView)

	cDecodeInfo.dstPictureResource = cDstPictureResource
	cDecodeInfo.pSetupReferenceSlot = nil
	// Note: Reference slots are not yet implemented. Any provided decodeInfo.ReferenceSlots are ignored.
	// Future implementation should iterate over ReferenceSlots and populate C structures.
	cDecodeInfo.referenceSlotCount = 0
	cDecodeInfo.pReferenceSlots = nil

	if C.call_vkCmdDecodeVideoKHR(C.VkCommandBuffer(commandBuffer), &cDecodeInfo) == 0 {
		return NewVulkanError(ErrorExtensionNotPresent, "CmdDecodeVideo", "video extension not loaded - call LoadVideoDeviceFunctions first")
	}
	return nil
}

// CmdEncodeVideo performs video encode operation in a command buffer.
// Returns an error if LoadVideoDeviceFunctions was not called or video extensions are not supported.
func CmdEncodeVideo(commandBuffer CommandBuffer, encodeInfo *VideoEncodeInfo) error {
	if commandBuffer == nil {
		return NewValidationError("commandBuffer", "cannot be nil")
	}
	if encodeInfo == nil {
		return NewValidationError("encodeInfo", "cannot be nil")
	}

	var cEncodeInfo C.VkVideoEncodeInfoKHR
	cEncodeInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_ENCODE_INFO_KHR
	cEncodeInfo.pNext = nil
	cEncodeInfo.flags = 0

	// Setup source picture resource
	var cSrcPictureResource C.VkVideoPictureResourceInfoKHR
	cSrcPictureResource.sType = C.VK_STRUCTURE_TYPE_VIDEO_PICTURE_RESOURCE_INFO_KHR
	cSrcPictureResource.pNext = nil
	cSrcPictureResource.codedOffset.x = C.int32_t(encodeInfo.SrcPictureResource.CodedOffset.X)
	cSrcPictureResource.codedOffset.y = C.int32_t(encodeInfo.SrcPictureResource.CodedOffset.Y)
	cSrcPictureResource.codedExtent.width = C.uint32_t(encodeInfo.SrcPictureResource.CodedExtent.Width)
	cSrcPictureResource.codedExtent.height = C.uint32_t(encodeInfo.SrcPictureResource.CodedExtent.Height)
	cSrcPictureResource.baseArrayLayer = C.uint32_t(encodeInfo.SrcPictureResource.BaseArrayLayer)
	cSrcPictureResource.imageViewBinding = C.VkImageView(encodeInfo.SrcPictureResource.ImageView)

	cEncodeInfo.srcPictureResource = cSrcPictureResource
	cEncodeInfo.pSetupReferenceSlot = nil
	// Note: Reference slots are not yet implemented. Any provided encodeInfo.ReferenceSlots are ignored.
	// Future implementation should iterate over ReferenceSlots and populate C structures.
	cEncodeInfo.referenceSlotCount = 0
	cEncodeInfo.pReferenceSlots = nil
	cEncodeInfo.dstBuffer = C.VkBuffer(encodeInfo.DstBuffer)
	cEncodeInfo.dstBufferOffset = C.VkDeviceSize(encodeInfo.DstBufferOffset)
	cEncodeInfo.dstBufferRange = C.VkDeviceSize(encodeInfo.DstBufferRange)

	if C.call_vkCmdEncodeVideoKHR(C.VkCommandBuffer(commandBuffer), &cEncodeInfo) == 0 {
		return NewVulkanError(ErrorExtensionNotPresent, "CmdEncodeVideo", "video extension not loaded - call LoadVideoDeviceFunctions first")
	}
	return nil
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
	if IsExtensionSupported(ExtensionNameVideoDecodeH264, extensions) {
		supportedCodecs = append(supportedCodecs, "H.264 (AVC) Decode")
	}
	if IsExtensionSupported(ExtensionNameVideoEncodeH264, extensions) {
		supportedCodecs = append(supportedCodecs, "H.264 (AVC) Encode")
	}

	// Check H.265 support
	if IsExtensionSupported(ExtensionNameVideoDecodeH265, extensions) {
		supportedCodecs = append(supportedCodecs, "H.265 (HEVC) Decode")
	}
	if IsExtensionSupported(ExtensionNameVideoEncodeH265, extensions) {
		supportedCodecs = append(supportedCodecs, "H.265 (HEVC) Encode")
	}

	// Check AV1 support
	if IsExtensionSupported(ExtensionNameVideoDecodeAV1, extensions) {
		supportedCodecs = append(supportedCodecs, "AV1 Decode")
	}
	if IsExtensionSupported(ExtensionNameVideoEncodeAV1, extensions) {
		supportedCodecs = append(supportedCodecs, "AV1 Encode")
	}

	return supportedCodecs, nil
}
