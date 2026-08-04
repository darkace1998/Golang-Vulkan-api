package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>

// Std header version descriptors required by VkVideoSessionCreateInfoKHR.
static VkExtensionProperties make_std_header(const char* name, uint32_t version) {
    VkExtensionProperties p;
    memset(&p, 0, sizeof(p));
    strncpy(p.extensionName, name, VK_MAX_EXTENSION_NAME_SIZE - 1);
    p.specVersion = version;
    return p;
}
static VkExtensionProperties stdHeaderH264Decode(void) {
    return make_std_header(VK_STD_VULKAN_VIDEO_CODEC_H264_DECODE_EXTENSION_NAME, VK_STD_VULKAN_VIDEO_CODEC_H264_DECODE_SPEC_VERSION);
}
static VkExtensionProperties stdHeaderH265Decode(void) {
    return make_std_header(VK_STD_VULKAN_VIDEO_CODEC_H265_DECODE_EXTENSION_NAME, VK_STD_VULKAN_VIDEO_CODEC_H265_DECODE_SPEC_VERSION);
}
static VkExtensionProperties stdHeaderH264Encode(void) {
    return make_std_header(VK_STD_VULKAN_VIDEO_CODEC_H264_ENCODE_EXTENSION_NAME, VK_STD_VULKAN_VIDEO_CODEC_H264_ENCODE_SPEC_VERSION);
}
static VkExtensionProperties stdHeaderH265Encode(void) {
    return make_std_header(VK_STD_VULKAN_VIDEO_CODEC_H265_ENCODE_EXTENSION_NAME, VK_STD_VULKAN_VIDEO_CODEC_H265_ENCODE_SPEC_VERSION);
}

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

import (
	"runtime"
	"sync"
	"unsafe"
)

// videoInstanceOnce and videoDeviceOnce ensure that video function pointers
// are loaded exactly once, preventing data races if multiple goroutines
// attempt to load them concurrently.
var (
	videoInstanceOnce   sync.Once
	videoInstanceLoaded bool
	videoDeviceOnce     sync.Once
	videoDeviceLoaded   bool
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

// VideoDecodeH264PictureLayoutFlags represents H.264 decode picture layouts
type VideoDecodeH264PictureLayoutFlags uint32

const (
	VideoDecodeH264PictureLayoutProgressive                VideoDecodeH264PictureLayoutFlags = 0
	VideoDecodeH264PictureLayoutInterlacedInterleavedLines VideoDecodeH264PictureLayoutFlags = 0x00000001
	VideoDecodeH264PictureLayoutInterlacedSeparatePlanes   VideoDecodeH264PictureLayoutFlags = 0x00000002
)

// VideoDecodeH264ProfileInfo is the codec-specific profile for H.264 decode
// (VkVideoDecodeH264ProfileInfoKHR).
type VideoDecodeH264ProfileInfo struct {
	StdProfileIdc H264Profile
	PictureLayout VideoDecodeH264PictureLayoutFlags
}

// VideoDecodeH265ProfileInfo is the codec-specific profile for H.265 decode
// (VkVideoDecodeH265ProfileInfoKHR).
type VideoDecodeH265ProfileInfo struct {
	StdProfileIdc H265Profile
}

// VideoEncodeH264ProfileInfo is the codec-specific profile for H.264 encode
// (VkVideoEncodeH264ProfileInfoKHR).
type VideoEncodeH264ProfileInfo struct {
	StdProfileIdc H264Profile
}

// VideoEncodeH265ProfileInfo is the codec-specific profile for H.265 encode
// (VkVideoEncodeH265ProfileInfoKHR).
type VideoEncodeH265ProfileInfo struct {
	StdProfileIdc H265Profile
}

// VideoProfileInfo describes a video profile.
//
// The Vulkan spec requires every VkVideoProfileInfoKHR to chain a
// codec-specific profile struct matching VideoCodecOperation. Set the
// matching codec field (e.g. DecodeH264 for VideoCodecOperationDecodeH264Bit)
// to control it; when left nil a documented default is chained instead
// (H.264: High profile, progressive layout; H.265: Main profile).
type VideoProfileInfo struct {
	VideoCodecOperation VideoCodecOperationFlags
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth

	// Codec-specific profile information; only the field matching
	// VideoCodecOperation is used.
	DecodeH264 *VideoDecodeH264ProfileInfo
	DecodeH265 *VideoDecodeH265ProfileInfo
	EncodeH264 *VideoEncodeH264ProfileInfo
	EncodeH265 *VideoEncodeH265ProfileInfo
}

// buildCVideoProfile fills cProfile from profile and chains the mandatory
// codec-specific profile struct (allocated on the Go heap and pinned into
// pinner, which must outlive the C call). Returns an error for codec
// operations this build cannot express.
func buildCVideoProfile(profile *VideoProfileInfo, pinner *runtime.Pinner, cProfile *C.VkVideoProfileInfoKHR) error {
	cProfile.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_INFO_KHR
	cProfile.pNext = nil
	cProfile.videoCodecOperation = C.VkVideoCodecOperationFlagBitsKHR(profile.VideoCodecOperation)
	cProfile.chromaSubsampling = C.VkVideoChromaSubsamplingFlagsKHR(profile.ChromaSubsampling)
	cProfile.lumaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(profile.LumaBitDepth)
	cProfile.chromaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(profile.ChromaBitDepth)

	switch profile.VideoCodecOperation {
	case VideoCodecOperationDecodeH264Bit:
		info := profile.DecodeH264
		if info == nil {
			info = &VideoDecodeH264ProfileInfo{StdProfileIdc: H264ProfileHigh, PictureLayout: VideoDecodeH264PictureLayoutProgressive}
		}
		cCodec := new(C.VkVideoDecodeH264ProfileInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_H264_PROFILE_INFO_KHR
		cCodec.pNext = nil
		cCodec.stdProfileIdc = C.StdVideoH264ProfileIdc(info.StdProfileIdc)
		cCodec.pictureLayout = C.VkVideoDecodeH264PictureLayoutFlagBitsKHR(info.PictureLayout)
		pinner.Pin(cCodec)
		cProfile.pNext = unsafe.Pointer(cCodec)

	case VideoCodecOperationDecodeH265Bit:
		info := profile.DecodeH265
		if info == nil {
			info = &VideoDecodeH265ProfileInfo{StdProfileIdc: H265ProfileMain}
		}
		cCodec := new(C.VkVideoDecodeH265ProfileInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_H265_PROFILE_INFO_KHR
		cCodec.pNext = nil
		cCodec.stdProfileIdc = C.StdVideoH265ProfileIdc(info.StdProfileIdc)
		pinner.Pin(cCodec)
		cProfile.pNext = unsafe.Pointer(cCodec)

	case VideoCodecOperationEncodeH264Bit:
		info := profile.EncodeH264
		if info == nil {
			info = &VideoEncodeH264ProfileInfo{StdProfileIdc: H264ProfileHigh}
		}
		cCodec := new(C.VkVideoEncodeH264ProfileInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_ENCODE_H264_PROFILE_INFO_KHR
		cCodec.pNext = nil
		cCodec.stdProfileIdc = C.StdVideoH264ProfileIdc(info.StdProfileIdc)
		pinner.Pin(cCodec)
		cProfile.pNext = unsafe.Pointer(cCodec)

	case VideoCodecOperationEncodeH265Bit:
		info := profile.EncodeH265
		if info == nil {
			info = &VideoEncodeH265ProfileInfo{StdProfileIdc: H265ProfileMain}
		}
		cCodec := new(C.VkVideoEncodeH265ProfileInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_ENCODE_H265_PROFILE_INFO_KHR
		cCodec.pNext = nil
		cCodec.stdProfileIdc = C.StdVideoH265ProfileIdc(info.StdProfileIdc)
		pinner.Pin(cCodec)
		cProfile.pNext = unsafe.Pointer(cCodec)

	case VideoCodecOperationDecodeAV1Bit, VideoCodecOperationEncodeAV1Bit:
		return NewValidationError("VideoCodecOperation", "AV1 requires Vulkan headers >= 1.3.277, which this build does not have")
	}

	return nil
}

// stdHeaderVersionFor returns the codec Std header version descriptor that
// VkVideoSessionCreateInfoKHR.pStdHeaderVersion requires.
func stdHeaderVersionFor(op VideoCodecOperationFlags) (*C.VkExtensionProperties, error) {
	props := new(C.VkExtensionProperties)
	switch op {
	case VideoCodecOperationDecodeH264Bit:
		*props = C.stdHeaderH264Decode()
	case VideoCodecOperationDecodeH265Bit:
		*props = C.stdHeaderH265Decode()
	case VideoCodecOperationEncodeH264Bit:
		*props = C.stdHeaderH264Encode()
	case VideoCodecOperationEncodeH265Bit:
		*props = C.stdHeaderH265Encode()
	default:
		return nil, NewValidationError("VideoCodecOperation", "unsupported codec operation for video session creation")
	}
	return props, nil
}

// VideoDecodeCapabilityFlags represents video decode capability flags
type VideoDecodeCapabilityFlags uint32

const (
	VideoDecodeCapabilityDpbAndOutputCoincideBit VideoDecodeCapabilityFlags = 0x00000001
	VideoDecodeCapabilityDpbAndOutputDistinctBit VideoDecodeCapabilityFlags = 0x00000002
)

// VideoDecodeCapabilities holds the decode-specific capabilities
// (VkVideoDecodeCapabilitiesKHR).
type VideoDecodeCapabilities struct {
	Flags VideoDecodeCapabilityFlags
}

// VideoDecodeH264Capabilities holds H.264 decode capabilities
// (VkVideoDecodeH264CapabilitiesKHR).
type VideoDecodeH264Capabilities struct {
	MaxLevelIdc            int32
	FieldOffsetGranularity Offset2D
}

// VideoDecodeH265Capabilities holds H.265 decode capabilities
// (VkVideoDecodeH265CapabilitiesKHR).
type VideoDecodeH265Capabilities struct {
	MaxLevelIdc int32
}

// VideoCapabilities represents video codec capabilities. The codec-specific
// sub-capabilities are populated according to the profile's codec operation.
type VideoCapabilities struct {
	Flags                         uint32
	MinBitstreamBufferOffsetAlign DeviceSize
	MinBitstreamBufferSizeAlign   DeviceSize
	PictureAccessGranularity      Extent2D
	MinCodedExtent                Extent2D
	MaxCodedExtent                Extent2D
	MaxDpbSlots                   uint32
	MaxActiveReferencePictures    uint32

	// Populated for decode profiles.
	Decode     *VideoDecodeCapabilities
	DecodeH264 *VideoDecodeH264Capabilities
	DecodeH265 *VideoDecodeH265Capabilities
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

// VideoDecodeH264SessionParametersCreateInfo sizes the H.264 decode parameter
// object (VkVideoDecodeH264SessionParametersCreateInfoKHR). Supplying actual
// SPS/PPS entries is not yet exposed; entries can be reserved here and the
// object updated later.
type VideoDecodeH264SessionParametersCreateInfo struct {
	MaxStdSPSCount uint32
	MaxStdPPSCount uint32
}

// VideoDecodeH265SessionParametersCreateInfo sizes the H.265 decode parameter
// object (VkVideoDecodeH265SessionParametersCreateInfoKHR).
type VideoDecodeH265SessionParametersCreateInfo struct {
	MaxStdVPSCount uint32
	MaxStdSPSCount uint32
	MaxStdPPSCount uint32
}

// VideoEncodeH264SessionParametersCreateInfo sizes the H.264 encode parameter
// object (VkVideoEncodeH264SessionParametersCreateInfoKHR).
type VideoEncodeH264SessionParametersCreateInfo struct {
	MaxStdSPSCount uint32
	MaxStdPPSCount uint32
}

// VideoEncodeH265SessionParametersCreateInfo sizes the H.265 encode parameter
// object (VkVideoEncodeH265SessionParametersCreateInfoKHR).
type VideoEncodeH265SessionParametersCreateInfo struct {
	MaxStdVPSCount uint32
	MaxStdSPSCount uint32
	MaxStdPPSCount uint32
}

// VideoSessionParametersCreateInfo contains parameters for video session
// parameters. The Vulkan spec requires the codec-specific create struct
// matching the session's codec operation to be chained; set exactly one of
// the codec fields.
type VideoSessionParametersCreateInfo struct {
	VideoSession           VideoSession
	VideoSessionParameters VideoSessionParameters

	// Codec-specific parameter capacities; set the field matching the video
	// session's codec operation.
	DecodeH264 *VideoDecodeH264SessionParametersCreateInfo
	DecodeH265 *VideoDecodeH265SessionParametersCreateInfo
	EncodeH264 *VideoEncodeH264SessionParametersCreateInfo
	EncodeH265 *VideoEncodeH265SessionParametersCreateInfo
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
// This function is thread-safe. The underlying C function pointers are loaded exactly once;
// subsequent calls return the cached result. Note that only one instance is supported at a time.
// If you need to reload for a different instance, use ResetVideoInstanceFunctions first.
//
// Returns false if the video extension functions could not be loaded (e.g., if the Vulkan
// implementation does not support the VK_KHR_video_queue extension).
func LoadVideoInstanceFunctions(instance Instance) bool {
	videoInstanceOnce.Do(func() {
		videoInstanceLoaded = C.loadVideoInstanceFunctions(C.VkInstance(instance)) != 0
	})
	return videoInstanceLoaded
}

// ResetVideoInstanceFunctions resets the instance function loader so that
// LoadVideoInstanceFunctions can be called again with a different instance.
// This is NOT thread-safe and must not be called concurrently with
// LoadVideoInstanceFunctions or any video API calls.
func ResetVideoInstanceFunctions() {
	videoInstanceOnce = sync.Once{}
	videoInstanceLoaded = false
}

// LoadVideoDeviceFunctions loads video extension functions that require a Vulkan device.
//
// This function MUST be called after creating a logical device and before using any video-related
// functionality. If this function is not called, all video API calls will fail.
//
// This function is thread-safe. The underlying C function pointers are loaded exactly once;
// subsequent calls return the cached result. Note that only one device is supported at a time.
// If you need to reload for a different device, use ResetVideoDeviceFunctions first.
// Returns false if any video extension function could not be loaded. This indicates the device
// does not fully support the VK_KHR_video_queue extension.
func LoadVideoDeviceFunctions(device Device) bool {
	videoDeviceOnce.Do(func() {
		videoDeviceLoaded = C.loadVideoDeviceFunctions(C.VkDevice(device)) != 0
	})
	return videoDeviceLoaded
}

// ResetVideoDeviceFunctions resets the device function loader so that
// LoadVideoDeviceFunctions can be called again with a different device.
// This is NOT thread-safe and must not be called concurrently with
// LoadVideoDeviceFunctions or any video API calls.
func ResetVideoDeviceFunctions() {
	videoDeviceOnce = sync.Once{}
	videoDeviceLoaded = false
}

// GetVideoCapabilities retrieves video codec capabilities for a physical device
func GetVideoCapabilities(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo) (*VideoCapabilities, error) {
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if videoProfile == nil {
		return nil, NewValidationError("videoProfile", "cannot be nil")
	}

	// Create C structures for video profile with the mandatory codec chain.
	var pinner runtime.Pinner
	defer pinner.Unpin()

	var cVideoProfile C.VkVideoProfileInfoKHR
	if err := buildCVideoProfile(videoProfile, &pinner, &cVideoProfile); err != nil {
		return nil, err
	}

	var cCaps C.VkVideoCapabilitiesKHR
	cCaps.sType = C.VK_STRUCTURE_TYPE_VIDEO_CAPABILITIES_KHR
	cCaps.pNext = nil

	// Decode profiles must chain VkVideoDecodeCapabilitiesKHR plus the codec
	// capabilities struct; both live on the Go heap and are pinned.
	var cDecodeCaps *C.VkVideoDecodeCapabilitiesKHR
	var cDecodeH264Caps *C.VkVideoDecodeH264CapabilitiesKHR
	var cDecodeH265Caps *C.VkVideoDecodeH265CapabilitiesKHR
	switch videoProfile.VideoCodecOperation {
	case VideoCodecOperationDecodeH264Bit:
		cDecodeH264Caps = new(C.VkVideoDecodeH264CapabilitiesKHR)
		cDecodeH264Caps.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_H264_CAPABILITIES_KHR
		pinner.Pin(cDecodeH264Caps)

		cDecodeCaps = new(C.VkVideoDecodeCapabilitiesKHR)
		cDecodeCaps.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_CAPABILITIES_KHR
		cDecodeCaps.pNext = unsafe.Pointer(cDecodeH264Caps)
		pinner.Pin(cDecodeCaps)
		cCaps.pNext = unsafe.Pointer(cDecodeCaps)
	case VideoCodecOperationDecodeH265Bit:
		cDecodeH265Caps = new(C.VkVideoDecodeH265CapabilitiesKHR)
		cDecodeH265Caps.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_H265_CAPABILITIES_KHR
		pinner.Pin(cDecodeH265Caps)

		cDecodeCaps = new(C.VkVideoDecodeCapabilitiesKHR)
		cDecodeCaps.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_CAPABILITIES_KHR
		cDecodeCaps.pNext = unsafe.Pointer(cDecodeH265Caps)
		pinner.Pin(cDecodeCaps)
		cCaps.pNext = unsafe.Pointer(cDecodeCaps)
	}

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

	if cDecodeCaps != nil {
		caps.Decode = &VideoDecodeCapabilities{Flags: VideoDecodeCapabilityFlags(cDecodeCaps.flags)}
	}
	if cDecodeH264Caps != nil {
		caps.DecodeH264 = &VideoDecodeH264Capabilities{
			MaxLevelIdc: int32(cDecodeH264Caps.maxLevelIdc),
			FieldOffsetGranularity: Offset2D{
				X: int32(cDecodeH264Caps.fieldOffsetGranularity.x),
				Y: int32(cDecodeH264Caps.fieldOffsetGranularity.y),
			},
		}
	}
	if cDecodeH265Caps != nil {
		caps.DecodeH265 = &VideoDecodeH265Capabilities{MaxLevelIdc: int32(cDecodeH265Caps.maxLevelIdc)}
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

	// Create C video profile structure with the mandatory codec chain. It must
	// be pinned because its address is stored inside cCreateInfo, which is Go
	// memory passed to C.
	var pinner runtime.Pinner
	defer pinner.Unpin()

	var cVideoProfile C.VkVideoProfileInfoKHR
	if err := buildCVideoProfile(createInfo.VideoProfile, &pinner, &cVideoProfile); err != nil {
		return VideoSession(NullHandle), err
	}
	pinner.Pin(&cVideoProfile)

	// The spec requires pStdHeaderVersion to name the codec Std header the
	// application was built against.
	cStdHeader, err := stdHeaderVersionFor(createInfo.VideoProfile.VideoCodecOperation)
	if err != nil {
		return VideoSession(NullHandle), err
	}
	pinner.Pin(cStdHeader)

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
	cCreateInfo.pStdHeaderVersion = cStdHeader

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

	// Chain the codec-specific parameter create struct; it lives on the Go
	// heap and must be pinned for the duration of the call.
	var pinner runtime.Pinner
	defer pinner.Unpin()

	var cCreateInfo C.VkVideoSessionParametersCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_SESSION_PARAMETERS_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.videoSessionParametersTemplate = C.VkVideoSessionParametersKHR(createInfo.VideoSessionParameters)
	cCreateInfo.videoSession = C.VkVideoSessionKHR(createInfo.VideoSession)

	switch {
	case createInfo.DecodeH264 != nil:
		cCodec := new(C.VkVideoDecodeH264SessionParametersCreateInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_H264_SESSION_PARAMETERS_CREATE_INFO_KHR
		cCodec.maxStdSPSCount = C.uint32_t(createInfo.DecodeH264.MaxStdSPSCount)
		cCodec.maxStdPPSCount = C.uint32_t(createInfo.DecodeH264.MaxStdPPSCount)
		pinner.Pin(cCodec)
		cCreateInfo.pNext = unsafe.Pointer(cCodec)
	case createInfo.DecodeH265 != nil:
		cCodec := new(C.VkVideoDecodeH265SessionParametersCreateInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_DECODE_H265_SESSION_PARAMETERS_CREATE_INFO_KHR
		cCodec.maxStdVPSCount = C.uint32_t(createInfo.DecodeH265.MaxStdVPSCount)
		cCodec.maxStdSPSCount = C.uint32_t(createInfo.DecodeH265.MaxStdSPSCount)
		cCodec.maxStdPPSCount = C.uint32_t(createInfo.DecodeH265.MaxStdPPSCount)
		pinner.Pin(cCodec)
		cCreateInfo.pNext = unsafe.Pointer(cCodec)
	case createInfo.EncodeH264 != nil:
		cCodec := new(C.VkVideoEncodeH264SessionParametersCreateInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_ENCODE_H264_SESSION_PARAMETERS_CREATE_INFO_KHR
		cCodec.maxStdSPSCount = C.uint32_t(createInfo.EncodeH264.MaxStdSPSCount)
		cCodec.maxStdPPSCount = C.uint32_t(createInfo.EncodeH264.MaxStdPPSCount)
		pinner.Pin(cCodec)
		cCreateInfo.pNext = unsafe.Pointer(cCodec)
	case createInfo.EncodeH265 != nil:
		cCodec := new(C.VkVideoEncodeH265SessionParametersCreateInfoKHR)
		cCodec.sType = C.VK_STRUCTURE_TYPE_VIDEO_ENCODE_H265_SESSION_PARAMETERS_CREATE_INFO_KHR
		cCodec.maxStdVPSCount = C.uint32_t(createInfo.EncodeH265.MaxStdVPSCount)
		cCodec.maxStdSPSCount = C.uint32_t(createInfo.EncodeH265.MaxStdSPSCount)
		cCodec.maxStdPPSCount = C.uint32_t(createInfo.EncodeH265.MaxStdPPSCount)
		pinner.Pin(cCodec)
		cCreateInfo.pNext = unsafe.Pointer(cCodec)
	}

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

// CmdBeginVideoCoding executes the operation
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

// CmdEndVideoCoding executes the operation
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

// CmdControlVideoCoding executes the operation
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

// CmdDecodeVideo executes the operation
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

// CmdEncodeVideo executes the operation
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
