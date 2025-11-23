package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>
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

	var cVideoProfileList C.VkVideoProfileListInfoKHR
	cVideoProfileList.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_LIST_INFO_KHR
	cVideoProfileList.pNext = nil
	cVideoProfileList.profileCount = 1
	cVideoProfileList.pProfiles = &cVideoProfile

	var cCaps C.VkVideoCapabilitiesKHR
	cCaps.sType = C.VK_STRUCTURE_TYPE_VIDEO_CAPABILITIES_KHR
	cCaps.pNext = nil

	result := Result(C.vkGetPhysicalDeviceVideoCapabilitiesKHR(
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

	var cVideoProfileList C.VkVideoProfileListInfoKHR
	cVideoProfileList.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_LIST_INFO_KHR
	cVideoProfileList.pNext = nil
	cVideoProfileList.profileCount = 1
	cVideoProfileList.pProfiles = &cVideoProfile

	// Create C video session create info
	var cCreateInfo C.VkVideoSessionCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_SESSION_CREATE_INFO_KHR
	cCreateInfo.pNext = unsafe.Pointer(&cVideoProfileList)
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
	result := Result(C.vkCreateVideoSessionKHR(
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
	C.vkDestroyVideoSessionKHR(C.VkDevice(device), C.VkVideoSessionKHR(videoSession), nil)
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
	result := Result(C.vkGetVideoSessionMemoryRequirementsKHR(
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

	result = Result(C.vkGetVideoSessionMemoryRequirementsKHR(
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

	result := Result(C.vkBindVideoSessionMemoryKHR(
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
	result := Result(C.vkCreateVideoSessionParametersKHR(
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
	C.vkDestroyVideoSessionParametersKHR(C.VkDevice(device), C.VkVideoSessionParametersKHR(videoSessionParameters), nil)
}

// VideoCodingControlInfo contains video coding control information
type VideoCodingControlInfo struct {
	Flags uint32
}

// CmdBeginVideoCoding begins video coding operations in a command buffer
func CmdBeginVideoCoding(commandBuffer CommandBuffer, beginInfo *VideoBeginCodingInfo) {
	if commandBuffer == nil || beginInfo == nil {
		return
	}

	var cBeginInfo C.VkVideoBeginCodingInfoKHR
	cBeginInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_BEGIN_CODING_INFO_KHR
	cBeginInfo.pNext = nil
	cBeginInfo.flags = 0
	cBeginInfo.videoSession = C.VkVideoSessionKHR(beginInfo.VideoSession)
	cBeginInfo.videoSessionParameters = C.VkVideoSessionParametersKHR(beginInfo.VideoSessionParameters)
	cBeginInfo.referenceSlotCount = 0
	cBeginInfo.pReferenceSlots = nil

	C.vkCmdBeginVideoCodingKHR(C.VkCommandBuffer(commandBuffer), &cBeginInfo)
}

// VideoBeginCodingInfo contains video begin coding information
type VideoBeginCodingInfo struct {
	VideoSession           VideoSession
	VideoSessionParameters VideoSessionParameters
}

// CmdEndVideoCoding ends video coding operations in a command buffer
func CmdEndVideoCoding(commandBuffer CommandBuffer) {
	if commandBuffer == nil {
		return
	}

	var cEndInfo C.VkVideoEndCodingInfoKHR
	cEndInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_END_CODING_INFO_KHR
	cEndInfo.pNext = nil
	cEndInfo.flags = 0

	C.vkCmdEndVideoCodingKHR(C.VkCommandBuffer(commandBuffer), &cEndInfo)
}

// CmdControlVideoCoding controls video coding operations
func CmdControlVideoCoding(commandBuffer CommandBuffer, controlInfo *VideoCodingControlInfo) {
	if commandBuffer == nil || controlInfo == nil {
		return
	}

	var cControlInfo C.VkVideoCodingControlInfoKHR
	cControlInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_CODING_CONTROL_INFO_KHR
	cControlInfo.pNext = nil
	cControlInfo.flags = C.VkVideoCodingControlFlagsKHR(controlInfo.Flags)

	C.vkCmdControlVideoCodingKHR(C.VkCommandBuffer(commandBuffer), &cControlInfo)
}

// CmdDecodeVideo performs video decode operation in a command buffer
func CmdDecodeVideo(commandBuffer CommandBuffer, decodeInfo *VideoDecodeInfo) {
	if commandBuffer == nil || decodeInfo == nil {
		return
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
	cDecodeInfo.referenceSlotCount = 0
	cDecodeInfo.pReferenceSlots = nil

	C.vkCmdDecodeVideoKHR(C.VkCommandBuffer(commandBuffer), &cDecodeInfo)
}

// CmdEncodeVideo performs video encode operation in a command buffer
func CmdEncodeVideo(commandBuffer CommandBuffer, encodeInfo *VideoEncodeInfo) {
	if commandBuffer == nil || encodeInfo == nil {
		return
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
	cEncodeInfo.referenceSlotCount = 0
	cEncodeInfo.pReferenceSlots = nil
	cEncodeInfo.dstBuffer = C.VkBuffer(encodeInfo.DstBuffer)
	cEncodeInfo.dstBufferOffset = C.VkDeviceSize(encodeInfo.DstBufferOffset)
	cEncodeInfo.dstBufferRange = C.VkDeviceSize(encodeInfo.DstBufferRange)

	C.vkCmdEncodeVideoKHR(C.VkCommandBuffer(commandBuffer), &cEncodeInfo)
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
