package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// ============================================================================
// Clear Commands
// ============================================================================

// ClearRect represents a clear rectangle
type ClearRect struct {
	Rect           Rect2D
	BaseArrayLayer uint32
	LayerCount     uint32
}

// ClearAttachment describes a clear attachment operation
type ClearAttachment struct {
	AspectMask      ImageAspectFlags
	ColorAttachment uint32
	ClearValue      ClearValue
}

// CmdClearAttachments clears attachment regions within a render pass
func CmdClearAttachments(commandBuffer CommandBuffer, attachments []ClearAttachment, rects []ClearRect) {
	if commandBuffer == nil {
		return
	}
	if len(attachments) == 0 || len(rects) == 0 {
		return
	}

	cAttachments := make([]C.VkClearAttachment, len(attachments))
	for i, att := range attachments {
		cAttachments[i].aspectMask = C.VkImageAspectFlags(att.AspectMask)
		cAttachments[i].colorAttachment = C.uint32_t(att.ColorAttachment)
		// Set clear value based on aspect
		if att.ClearValue.IsDepthStencil {
			// Use C struct field sizes for correct offset calculation to match VkClearDepthStencilValue layout
			cDepthStencil := (*C.VkClearDepthStencilValue)(unsafe.Pointer(&cAttachments[i].clearValue))
			cDepthStencil.depth = C.float(att.ClearValue.DepthStencil.Depth)
			cDepthStencil.stencil = C.uint32_t(att.ClearValue.DepthStencil.Stencil)
		} else {
			*(*[4]float32)(unsafe.Pointer(&cAttachments[i].clearValue)) = att.ClearValue.Color.Float32
		}
	}

	cRects := make([]C.VkClearRect, len(rects))
	for i, rect := range rects {
		cRects[i].rect.offset.x = C.int32_t(rect.Rect.Offset.X)
		cRects[i].rect.offset.y = C.int32_t(rect.Rect.Offset.Y)
		cRects[i].rect.extent.width = C.uint32_t(rect.Rect.Extent.Width)
		cRects[i].rect.extent.height = C.uint32_t(rect.Rect.Extent.Height)
		cRects[i].baseArrayLayer = C.uint32_t(rect.BaseArrayLayer)
		cRects[i].layerCount = C.uint32_t(rect.LayerCount)
	}

	C.vkCmdClearAttachments(
		C.VkCommandBuffer(commandBuffer),
		C.uint32_t(len(cAttachments)),
		&cAttachments[0],
		C.uint32_t(len(cRects)),
		&cRects[0],
	)
}

// CmdClearColorImage clears a color image outside of a render pass
func CmdClearColorImage(commandBuffer CommandBuffer, image Image, imageLayout ImageLayout, color *ClearColorValue, ranges []ImageSubresourceRange) {
	if commandBuffer == nil || image == nil || color == nil {
		return
	}
	if len(ranges) == 0 {
		return
	}

	var cColor C.VkClearColorValue
	*(*[4]float32)(unsafe.Pointer(&cColor)) = color.Float32

	cRanges := make([]C.VkImageSubresourceRange, len(ranges))
	for i, r := range ranges {
		cRanges[i].aspectMask = C.VkImageAspectFlags(r.AspectMask)
		cRanges[i].baseMipLevel = C.uint32_t(r.BaseMipLevel)
		cRanges[i].levelCount = C.uint32_t(r.LevelCount)
		cRanges[i].baseArrayLayer = C.uint32_t(r.BaseArrayLayer)
		cRanges[i].layerCount = C.uint32_t(r.LayerCount)
	}

	C.vkCmdClearColorImage(
		C.VkCommandBuffer(commandBuffer),
		C.VkImage(image),
		C.VkImageLayout(imageLayout),
		&cColor,
		C.uint32_t(len(cRanges)),
		&cRanges[0],
	)
}

// CmdClearDepthStencilImage clears a depth/stencil image outside of a render pass
func CmdClearDepthStencilImage(commandBuffer CommandBuffer, image Image, imageLayout ImageLayout, depthStencil *ClearDepthStencilValue, ranges []ImageSubresourceRange) {
	if commandBuffer == nil || image == nil || depthStencil == nil {
		return
	}
	if len(ranges) == 0 {
		return
	}

	var cDepthStencil C.VkClearDepthStencilValue
	cDepthStencil.depth = C.float(depthStencil.Depth)
	cDepthStencil.stencil = C.uint32_t(depthStencil.Stencil)

	cRanges := make([]C.VkImageSubresourceRange, len(ranges))
	for i, r := range ranges {
		cRanges[i].aspectMask = C.VkImageAspectFlags(r.AspectMask)
		cRanges[i].baseMipLevel = C.uint32_t(r.BaseMipLevel)
		cRanges[i].levelCount = C.uint32_t(r.LevelCount)
		cRanges[i].baseArrayLayer = C.uint32_t(r.BaseArrayLayer)
		cRanges[i].layerCount = C.uint32_t(r.LayerCount)
	}

	C.vkCmdClearDepthStencilImage(
		C.VkCommandBuffer(commandBuffer),
		C.VkImage(image),
		C.VkImageLayout(imageLayout),
		&cDepthStencil,
		C.uint32_t(len(cRanges)),
		&cRanges[0],
	)
}

// ============================================================================
// Pipeline Cache
// ============================================================================

// PipelineCacheCreateFlags represents pipeline cache creation flags
type PipelineCacheCreateFlags uint32

const (
	PipelineCacheCreateExternallySynchronized PipelineCacheCreateFlags = 0x00000001
)

// PipelineCacheCreateInfo contains pipeline cache creation information
type PipelineCacheCreateInfo struct {
	Flags       PipelineCacheCreateFlags
	InitialData []byte
}

// CreatePipelineCache creates a pipeline cache
func CreatePipelineCache(device Device, createInfo *PipelineCacheCreateInfo) (PipelineCache, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	var cCreateInfo C.VkPipelineCacheCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_PIPELINE_CACHE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkPipelineCacheCreateFlags(createInfo.Flags)

	if len(createInfo.InitialData) > 0 {
		cCreateInfo.initialDataSize = C.size_t(len(createInfo.InitialData))
		cCreateInfo.pInitialData = unsafe.Pointer(&createInfo.InitialData[0])
	} else {
		cCreateInfo.initialDataSize = 0
		cCreateInfo.pInitialData = nil
	}

	var pipelineCache C.VkPipelineCache
	result := Result(C.vkCreatePipelineCache(C.VkDevice(device), &cCreateInfo, nil, &pipelineCache))
	if result != Success {
		return nil, NewVulkanError(result, "CreatePipelineCache", "Vulkan pipeline cache creation failed")
	}

	return PipelineCache(pipelineCache), nil
}

// DestroyPipelineCache destroys a pipeline cache
func DestroyPipelineCache(device Device, pipelineCache PipelineCache) {
	if device != nil && pipelineCache != nil {
		C.vkDestroyPipelineCache(C.VkDevice(device), C.VkPipelineCache(pipelineCache), nil)
	}
}

// GetPipelineCacheData retrieves the data from a pipeline cache
func GetPipelineCacheData(device Device, pipelineCache PipelineCache) ([]byte, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if pipelineCache == nil {
		return nil, NewValidationError("pipelineCache", "cannot be nil")
	}

	// First, get the size of the cache data
	var dataSize C.size_t
	result := Result(C.vkGetPipelineCacheData(C.VkDevice(device), C.VkPipelineCache(pipelineCache), &dataSize, nil))
	if result != Success {
		return nil, NewVulkanError(result, "GetPipelineCacheData", "failed to get pipeline cache data size")
	}

	if dataSize == 0 {
		return []byte{}, nil
	}

	// Allocate buffer and retrieve data
	data := make([]byte, dataSize)
	result = Result(C.vkGetPipelineCacheData(C.VkDevice(device), C.VkPipelineCache(pipelineCache), &dataSize, unsafe.Pointer(&data[0])))
	if result != Success {
		return nil, NewVulkanError(result, "GetPipelineCacheData", "failed to get pipeline cache data")
	}

	return data[:dataSize], nil
}

// MergePipelineCaches merges multiple pipeline caches into a destination cache
func MergePipelineCaches(device Device, dstCache PipelineCache, srcCaches []PipelineCache) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if dstCache == nil {
		return NewValidationError("dstCache", "cannot be nil")
	}
	if len(srcCaches) == 0 {
		return nil // Nothing to merge
	}

	cSrcCaches := make([]C.VkPipelineCache, len(srcCaches))
	for i, cache := range srcCaches {
		if cache == nil {
			return NewValidationError("srcCaches", "contains nil cache")
		}
		cSrcCaches[i] = C.VkPipelineCache(cache)
	}

	result := Result(C.vkMergePipelineCaches(
		C.VkDevice(device),
		C.VkPipelineCache(dstCache),
		C.uint32_t(len(cSrcCaches)),
		&cSrcCaches[0],
	))
	if result != Success {
		return NewVulkanError(result, "MergePipelineCaches", "failed to merge pipeline caches")
	}

	return nil
}

// ============================================================================
// Buffer View
// ============================================================================

// BufferViewCreateInfo contains buffer view creation information
type BufferViewCreateInfo struct {
	Buffer Buffer
	Format Format
	Offset DeviceSize
	Range  DeviceSize
}

// CreateBufferView creates a buffer view
func CreateBufferView(device Device, createInfo *BufferViewCreateInfo) (BufferView, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.Buffer == nil {
		return nil, NewValidationError("Buffer", "cannot be nil")
	}

	var cCreateInfo C.VkBufferViewCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_BUFFER_VIEW_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.buffer = C.VkBuffer(createInfo.Buffer)
	cCreateInfo.format = C.VkFormat(createInfo.Format)
	cCreateInfo.offset = C.VkDeviceSize(createInfo.Offset)
	cCreateInfo._range = C.VkDeviceSize(createInfo.Range)

	var bufferView C.VkBufferView
	result := Result(C.vkCreateBufferView(C.VkDevice(device), &cCreateInfo, nil, &bufferView))
	if result != Success {
		return nil, NewVulkanError(result, "CreateBufferView", "Vulkan buffer view creation failed")
	}

	return BufferView(bufferView), nil
}

// DestroyBufferView destroys a buffer view
func DestroyBufferView(device Device, bufferView BufferView) {
	if device != nil && bufferView != nil {
		C.vkDestroyBufferView(C.VkDevice(device), C.VkBufferView(bufferView), nil)
	}
}

// ============================================================================
// Format Queries
// ============================================================================

// FormatFeatureFlags represents format feature flags
type FormatFeatureFlags uint32

const (
	FormatFeatureSampledImageBit                            FormatFeatureFlags = C.VK_FORMAT_FEATURE_SAMPLED_IMAGE_BIT
	FormatFeatureStorageImageBit                            FormatFeatureFlags = C.VK_FORMAT_FEATURE_STORAGE_IMAGE_BIT
	FormatFeatureStorageImageAtomicBit                      FormatFeatureFlags = C.VK_FORMAT_FEATURE_STORAGE_IMAGE_ATOMIC_BIT
	FormatFeatureUniformTexelBufferBit                      FormatFeatureFlags = C.VK_FORMAT_FEATURE_UNIFORM_TEXEL_BUFFER_BIT
	FormatFeatureStorageTexelBufferBit                      FormatFeatureFlags = C.VK_FORMAT_FEATURE_STORAGE_TEXEL_BUFFER_BIT
	FormatFeatureStorageTexelBufferAtomicBit                FormatFeatureFlags = C.VK_FORMAT_FEATURE_STORAGE_TEXEL_BUFFER_ATOMIC_BIT
	FormatFeatureVertexBufferBit                            FormatFeatureFlags = C.VK_FORMAT_FEATURE_VERTEX_BUFFER_BIT
	FormatFeatureColorAttachmentBit                         FormatFeatureFlags = C.VK_FORMAT_FEATURE_COLOR_ATTACHMENT_BIT
	FormatFeatureColorAttachmentBlendBit                    FormatFeatureFlags = C.VK_FORMAT_FEATURE_COLOR_ATTACHMENT_BLEND_BIT
	FormatFeatureDepthStencilAttachmentBit                  FormatFeatureFlags = C.VK_FORMAT_FEATURE_DEPTH_STENCIL_ATTACHMENT_BIT
	FormatFeatureBlitSrcBit                                 FormatFeatureFlags = C.VK_FORMAT_FEATURE_BLIT_SRC_BIT
	FormatFeatureBlitDstBit                                 FormatFeatureFlags = C.VK_FORMAT_FEATURE_BLIT_DST_BIT
	FormatFeatureSampledImageFilterLinearBit                FormatFeatureFlags = C.VK_FORMAT_FEATURE_SAMPLED_IMAGE_FILTER_LINEAR_BIT
	FormatFeatureTransferSrcBit                             FormatFeatureFlags = C.VK_FORMAT_FEATURE_TRANSFER_SRC_BIT
	FormatFeatureTransferDstBit                             FormatFeatureFlags = C.VK_FORMAT_FEATURE_TRANSFER_DST_BIT
	FormatFeatureMidpointChromaSamplesBit                   FormatFeatureFlags = C.VK_FORMAT_FEATURE_MIDPOINT_CHROMA_SAMPLES_BIT
	FormatFeatureSampledImageYcbcrConversionLinearFilterBit FormatFeatureFlags = C.VK_FORMAT_FEATURE_SAMPLED_IMAGE_YCBCR_CONVERSION_LINEAR_FILTER_BIT
	FormatFeatureSampledImageFilterMinmaxBit                FormatFeatureFlags = C.VK_FORMAT_FEATURE_SAMPLED_IMAGE_FILTER_MINMAX_BIT
)

// FormatProperties contains format properties
type FormatProperties struct {
	LinearTilingFeatures  FormatFeatureFlags
	OptimalTilingFeatures FormatFeatureFlags
	BufferFeatures        FormatFeatureFlags
}

// GetPhysicalDeviceFormatProperties returns format properties for a physical device
func GetPhysicalDeviceFormatProperties(physicalDevice PhysicalDevice, format Format) FormatProperties {
	if physicalDevice == nil {
		return FormatProperties{}
	}

	var cProps C.VkFormatProperties
	C.vkGetPhysicalDeviceFormatProperties(C.VkPhysicalDevice(physicalDevice), C.VkFormat(format), &cProps)

	return FormatProperties{
		LinearTilingFeatures:  FormatFeatureFlags(cProps.linearTilingFeatures),
		OptimalTilingFeatures: FormatFeatureFlags(cProps.optimalTilingFeatures),
		BufferFeatures:        FormatFeatureFlags(cProps.bufferFeatures),
	}
}

// ImageFormatProperties contains image format properties
type ImageFormatProperties struct {
	MaxExtent       Extent3D
	MaxMipLevels    uint32
	MaxArrayLayers  uint32
	SampleCounts    SampleCountFlags
	MaxResourceSize DeviceSize
}

// GetPhysicalDeviceImageFormatProperties returns image format properties for a physical device
func GetPhysicalDeviceImageFormatProperties(physicalDevice PhysicalDevice, format Format, imageType ImageType, tiling ImageTiling, usage ImageUsageFlags, flags ImageCreateFlags) (ImageFormatProperties, error) {
	if physicalDevice == nil {
		return ImageFormatProperties{}, NewValidationError("physicalDevice", "cannot be nil")
	}

	var cProps C.VkImageFormatProperties
	result := Result(C.vkGetPhysicalDeviceImageFormatProperties(
		C.VkPhysicalDevice(physicalDevice),
		C.VkFormat(format),
		C.VkImageType(imageType),
		C.VkImageTiling(tiling),
		C.VkImageUsageFlags(usage),
		C.VkImageCreateFlags(flags),
		&cProps,
	))

	if result != Success {
		return ImageFormatProperties{}, NewVulkanError(result, "GetPhysicalDeviceImageFormatProperties", "failed to get image format properties")
	}

	return ImageFormatProperties{
		MaxExtent: Extent3D{
			Width:  uint32(cProps.maxExtent.width),
			Height: uint32(cProps.maxExtent.height),
			Depth:  uint32(cProps.maxExtent.depth),
		},
		MaxMipLevels:    uint32(cProps.maxMipLevels),
		MaxArrayLayers:  uint32(cProps.maxArrayLayers),
		SampleCounts:    SampleCountFlags(cProps.sampleCounts),
		MaxResourceSize: DeviceSize(cProps.maxResourceSize),
	}, nil
}

// ============================================================================
// Sparse Resources
// ============================================================================

// SparseImageFormatFlags represents sparse image format flags
type SparseImageFormatFlags uint32

const (
	SparseImageFormatSingleMiptailBit        SparseImageFormatFlags = C.VK_SPARSE_IMAGE_FORMAT_SINGLE_MIPTAIL_BIT
	SparseImageFormatAlignedMipSizeBit       SparseImageFormatFlags = C.VK_SPARSE_IMAGE_FORMAT_ALIGNED_MIP_SIZE_BIT
	SparseImageFormatNonstandardBlockSizeBit SparseImageFormatFlags = C.VK_SPARSE_IMAGE_FORMAT_NONSTANDARD_BLOCK_SIZE_BIT
)

// SparseImageFormatProperties contains sparse image format properties
type SparseImageFormatProperties struct {
	AspectMask       ImageAspectFlags
	ImageGranularity Extent3D
	Flags            SparseImageFormatFlags
}

// GetPhysicalDeviceSparseImageFormatProperties returns sparse image format properties
func GetPhysicalDeviceSparseImageFormatProperties(physicalDevice PhysicalDevice, format Format, imageType ImageType, samples SampleCountFlags, usage ImageUsageFlags, tiling ImageTiling) []SparseImageFormatProperties {
	if physicalDevice == nil {
		return nil
	}

	var count C.uint32_t
	C.vkGetPhysicalDeviceSparseImageFormatProperties(
		C.VkPhysicalDevice(physicalDevice),
		C.VkFormat(format),
		C.VkImageType(imageType),
		C.VkSampleCountFlagBits(samples),
		C.VkImageUsageFlags(usage),
		C.VkImageTiling(tiling),
		&count,
		nil,
	)

	if count == 0 {
		return nil
	}

	cProps := make([]C.VkSparseImageFormatProperties, count)
	C.vkGetPhysicalDeviceSparseImageFormatProperties(
		C.VkPhysicalDevice(physicalDevice),
		C.VkFormat(format),
		C.VkImageType(imageType),
		C.VkSampleCountFlagBits(samples),
		C.VkImageUsageFlags(usage),
		C.VkImageTiling(tiling),
		&count,
		&cProps[0],
	)

	props := make([]SparseImageFormatProperties, count)
	for i := range props {
		props[i] = SparseImageFormatProperties{
			AspectMask: ImageAspectFlags(cProps[i].aspectMask),
			ImageGranularity: Extent3D{
				Width:  uint32(cProps[i].imageGranularity.width),
				Height: uint32(cProps[i].imageGranularity.height),
				Depth:  uint32(cProps[i].imageGranularity.depth),
			},
			Flags: SparseImageFormatFlags(cProps[i].flags),
		}
	}

	return props
}

// SparseImageMemoryRequirements contains sparse image memory requirements
type SparseImageMemoryRequirements struct {
	FormatProperties     SparseImageFormatProperties
	ImageMipTailFirstLod uint32
	ImageMipTailSize     DeviceSize
	ImageMipTailOffset   DeviceSize
	ImageMipTailStride   DeviceSize
}

// GetImageSparseMemoryRequirements returns sparse memory requirements for an image
func GetImageSparseMemoryRequirements(device Device, image Image) []SparseImageMemoryRequirements {
	if device == nil || image == nil {
		return nil
	}

	var count C.uint32_t
	C.vkGetImageSparseMemoryRequirements(C.VkDevice(device), C.VkImage(image), &count, nil)

	if count == 0 {
		return nil
	}

	cReqs := make([]C.VkSparseImageMemoryRequirements, count)
	C.vkGetImageSparseMemoryRequirements(C.VkDevice(device), C.VkImage(image), &count, &cReqs[0])

	reqs := make([]SparseImageMemoryRequirements, count)
	for i := range reqs {
		reqs[i] = SparseImageMemoryRequirements{
			FormatProperties: SparseImageFormatProperties{
				AspectMask: ImageAspectFlags(cReqs[i].formatProperties.aspectMask),
				ImageGranularity: Extent3D{
					Width:  uint32(cReqs[i].formatProperties.imageGranularity.width),
					Height: uint32(cReqs[i].formatProperties.imageGranularity.height),
					Depth:  uint32(cReqs[i].formatProperties.imageGranularity.depth),
				},
				Flags: SparseImageFormatFlags(cReqs[i].formatProperties.flags),
			},
			ImageMipTailFirstLod: uint32(cReqs[i].imageMipTailFirstLod),
			ImageMipTailSize:     DeviceSize(cReqs[i].imageMipTailSize),
			ImageMipTailOffset:   DeviceSize(cReqs[i].imageMipTailOffset),
			ImageMipTailStride:   DeviceSize(cReqs[i].imageMipTailStride),
		}
	}

	return reqs
}
