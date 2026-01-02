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

// BufferCreateInfo contains buffer creation information
type BufferCreateInfo struct {
	Flags       BufferCreateFlags
	Size        DeviceSize
	Usage       BufferUsageFlags
	SharingMode SharingMode
}

// BufferCreateFlags represents buffer creation flags
type BufferCreateFlags uint32

const (
	BufferCreateSparseBindingBit              BufferCreateFlags = C.VK_BUFFER_CREATE_SPARSE_BINDING_BIT
	BufferCreateSparseResidencyBit            BufferCreateFlags = C.VK_BUFFER_CREATE_SPARSE_RESIDENCY_BIT
	BufferCreateSparseAliasedBit              BufferCreateFlags = C.VK_BUFFER_CREATE_SPARSE_ALIASED_BIT
	BufferCreateProtectedBit                  BufferCreateFlags = C.VK_BUFFER_CREATE_PROTECTED_BIT
	BufferCreateDeviceAddressCaptureReplayBit BufferCreateFlags = C.VK_BUFFER_CREATE_DEVICE_ADDRESS_CAPTURE_REPLAY_BIT
)

// BufferUsageFlags represents buffer usage flags
type BufferUsageFlags uint32

const (
	BufferUsageTransferSrcBit         BufferUsageFlags = C.VK_BUFFER_USAGE_TRANSFER_SRC_BIT
	BufferUsageTransferDstBit         BufferUsageFlags = C.VK_BUFFER_USAGE_TRANSFER_DST_BIT
	BufferUsageUniformTexelBufferBit  BufferUsageFlags = C.VK_BUFFER_USAGE_UNIFORM_TEXEL_BUFFER_BIT
	BufferUsageStorageTexelBufferBit  BufferUsageFlags = C.VK_BUFFER_USAGE_STORAGE_TEXEL_BUFFER_BIT
	BufferUsageUniformBufferBit       BufferUsageFlags = C.VK_BUFFER_USAGE_UNIFORM_BUFFER_BIT
	BufferUsageStorageBufferBit       BufferUsageFlags = C.VK_BUFFER_USAGE_STORAGE_BUFFER_BIT
	BufferUsageIndexBufferBit         BufferUsageFlags = C.VK_BUFFER_USAGE_INDEX_BUFFER_BIT
	BufferUsageVertexBufferBit        BufferUsageFlags = C.VK_BUFFER_USAGE_VERTEX_BUFFER_BIT
	BufferUsageIndirectBufferBit      BufferUsageFlags = C.VK_BUFFER_USAGE_INDIRECT_BUFFER_BIT
	BufferUsageShaderDeviceAddressBit BufferUsageFlags = C.VK_BUFFER_USAGE_SHADER_DEVICE_ADDRESS_BIT
)

// SharingMode represents resource sharing mode
type SharingMode int32

const (
	SharingModeExclusive  SharingMode = C.VK_SHARING_MODE_EXCLUSIVE
	SharingModeConcurrent SharingMode = C.VK_SHARING_MODE_CONCURRENT
)

// MemoryAllocateInfo contains memory allocation information
type MemoryAllocateInfo struct {
	AllocationSize  DeviceSize
	MemoryTypeIndex uint32
}

// MemoryRequirements contains memory requirements
type MemoryRequirements struct {
	Size           DeviceSize
	Alignment      DeviceSize
	MemoryTypeBits uint32
}

// ImageCreateInfo contains image creation information
type ImageCreateInfo struct {
	Flags         ImageCreateFlags
	ImageType     ImageType
	Format        Format
	Extent        Extent3D
	MipLevels     uint32
	ArrayLayers   uint32
	Samples       SampleCountFlags
	Tiling        ImageTiling
	Usage         ImageUsageFlags
	SharingMode   SharingMode
	InitialLayout ImageLayout
}

// ImageType represents image types
type ImageType int32

const (
	ImageType1D ImageType = C.VK_IMAGE_TYPE_1D
	ImageType2D ImageType = C.VK_IMAGE_TYPE_2D
	ImageType3D ImageType = C.VK_IMAGE_TYPE_3D
)

// ImageCreateFlags represents image creation flags
type ImageCreateFlags uint32

const (
	ImageCreateSparseBindingBit                     ImageCreateFlags = C.VK_IMAGE_CREATE_SPARSE_BINDING_BIT
	ImageCreateSparseResidencyBit                   ImageCreateFlags = C.VK_IMAGE_CREATE_SPARSE_RESIDENCY_BIT
	ImageCreateSparseAliasedBit                     ImageCreateFlags = C.VK_IMAGE_CREATE_SPARSE_ALIASED_BIT
	ImageCreateMutableFormatBit                     ImageCreateFlags = C.VK_IMAGE_CREATE_MUTABLE_FORMAT_BIT
	ImageCreateCubeCompatibleBit                    ImageCreateFlags = C.VK_IMAGE_CREATE_CUBE_COMPATIBLE_BIT
	ImageCreateAliasBit                             ImageCreateFlags = C.VK_IMAGE_CREATE_ALIAS_BIT
	ImageCreateSplitInstanceBindRegionsBit          ImageCreateFlags = C.VK_IMAGE_CREATE_SPLIT_INSTANCE_BIND_REGIONS_BIT
	ImageCreate2DArrayCompatibleBit                 ImageCreateFlags = C.VK_IMAGE_CREATE_2D_ARRAY_COMPATIBLE_BIT
	ImageCreateBlockTexelViewCompatibleBit          ImageCreateFlags = C.VK_IMAGE_CREATE_BLOCK_TEXEL_VIEW_COMPATIBLE_BIT
	ImageCreateExtendedUsageBit                     ImageCreateFlags = C.VK_IMAGE_CREATE_EXTENDED_USAGE_BIT
	ImageCreateProtectedBit                         ImageCreateFlags = C.VK_IMAGE_CREATE_PROTECTED_BIT
	ImageCreateDisjointBit                          ImageCreateFlags = C.VK_IMAGE_CREATE_DISJOINT_BIT
	ImageCreateCornerSampledBitNV                   ImageCreateFlags = C.VK_IMAGE_CREATE_CORNER_SAMPLED_BIT_NV
	ImageCreateSampleLocationsCompatibleDepthBitEXT ImageCreateFlags = C.VK_IMAGE_CREATE_SAMPLE_LOCATIONS_COMPATIBLE_DEPTH_BIT_EXT
	ImageCreateSubsampledBitEXT                     ImageCreateFlags = C.VK_IMAGE_CREATE_SUBSAMPLED_BIT_EXT
)

// Format represents pixel formats
type Format int32

const (
	FormatUndefined           Format = C.VK_FORMAT_UNDEFINED
	FormatR4G4UnormPack8      Format = C.VK_FORMAT_R4G4_UNORM_PACK8
	FormatR4G4B4A4UnormPack16 Format = C.VK_FORMAT_R4G4B4A4_UNORM_PACK16
	FormatB4G4R4A4UnormPack16 Format = C.VK_FORMAT_B4G4R4A4_UNORM_PACK16
	FormatR5G6B5UnormPack16   Format = C.VK_FORMAT_R5G6B5_UNORM_PACK16
	FormatB5G6R5UnormPack16   Format = C.VK_FORMAT_B5G6R5_UNORM_PACK16
	FormatR5G5B5A1UnormPack16 Format = C.VK_FORMAT_R5G5B5A1_UNORM_PACK16
	FormatB5G5R5A1UnormPack16 Format = C.VK_FORMAT_B5G5R5A1_UNORM_PACK16
	FormatA1R5G5B5UnormPack16 Format = C.VK_FORMAT_A1R5G5B5_UNORM_PACK16
	FormatR8Unorm             Format = C.VK_FORMAT_R8_UNORM
	FormatR8Snorm             Format = C.VK_FORMAT_R8_SNORM
	FormatR8Uscaled           Format = C.VK_FORMAT_R8_USCALED
	FormatR8Sscaled           Format = C.VK_FORMAT_R8_SSCALED
	FormatR8Uint              Format = C.VK_FORMAT_R8_UINT
	FormatR8Sint              Format = C.VK_FORMAT_R8_SINT
	FormatR8Srgb              Format = C.VK_FORMAT_R8_SRGB
	FormatR8G8Unorm           Format = C.VK_FORMAT_R8G8_UNORM
	FormatR8G8Snorm           Format = C.VK_FORMAT_R8G8_SNORM
	FormatR8G8Uscaled         Format = C.VK_FORMAT_R8G8_USCALED
	FormatR8G8Sscaled         Format = C.VK_FORMAT_R8G8_SSCALED
	FormatR8G8Uint            Format = C.VK_FORMAT_R8G8_UINT
	FormatR8G8Sint            Format = C.VK_FORMAT_R8G8_SINT
	FormatR8G8Srgb            Format = C.VK_FORMAT_R8G8_SRGB
	FormatR8G8B8Unorm         Format = C.VK_FORMAT_R8G8B8_UNORM
	FormatR8G8B8Snorm         Format = C.VK_FORMAT_R8G8B8_SNORM
	FormatR8G8B8Uscaled       Format = C.VK_FORMAT_R8G8B8_USCALED
	FormatR8G8B8Sscaled       Format = C.VK_FORMAT_R8G8B8_SSCALED
	FormatR8G8B8Uint          Format = C.VK_FORMAT_R8G8B8_UINT
	FormatR8G8B8Sint          Format = C.VK_FORMAT_R8G8B8_SINT
	FormatR8G8B8Srgb          Format = C.VK_FORMAT_R8G8B8_SRGB
	FormatB8G8R8Unorm         Format = C.VK_FORMAT_B8G8R8_UNORM
	FormatB8G8R8Snorm         Format = C.VK_FORMAT_B8G8R8_SNORM
	FormatB8G8R8Uscaled       Format = C.VK_FORMAT_B8G8R8_USCALED
	FormatB8G8R8Sscaled       Format = C.VK_FORMAT_B8G8R8_SSCALED
	FormatB8G8R8Uint          Format = C.VK_FORMAT_B8G8R8_UINT
	FormatB8G8R8Sint          Format = C.VK_FORMAT_B8G8R8_SINT
	FormatB8G8R8Srgb          Format = C.VK_FORMAT_B8G8R8_SRGB
	FormatR8G8B8A8Unorm       Format = C.VK_FORMAT_R8G8B8A8_UNORM
	FormatR8G8B8A8Snorm       Format = C.VK_FORMAT_R8G8B8A8_SNORM
	FormatR8G8B8A8Uscaled     Format = C.VK_FORMAT_R8G8B8A8_USCALED
	FormatR8G8B8A8Sscaled     Format = C.VK_FORMAT_R8G8B8A8_SSCALED
	FormatR8G8B8A8Uint        Format = C.VK_FORMAT_R8G8B8A8_UINT
	FormatR8G8B8A8Sint        Format = C.VK_FORMAT_R8G8B8A8_SINT
	FormatR8G8B8A8Srgb        Format = C.VK_FORMAT_R8G8B8A8_SRGB
	FormatB8G8R8A8Unorm       Format = C.VK_FORMAT_B8G8R8A8_UNORM
	FormatB8G8R8A8Snorm       Format = C.VK_FORMAT_B8G8R8A8_SNORM
	FormatB8G8R8A8Uscaled     Format = C.VK_FORMAT_B8G8R8A8_USCALED
	FormatB8G8R8A8Sscaled     Format = C.VK_FORMAT_B8G8R8A8_SSCALED
	FormatB8G8R8A8Uint        Format = C.VK_FORMAT_B8G8R8A8_UINT
	FormatB8G8R8A8Sint        Format = C.VK_FORMAT_B8G8R8A8_SINT
	FormatB8G8R8A8Srgb        Format = C.VK_FORMAT_B8G8R8A8_SRGB
	FormatD16Unorm            Format = C.VK_FORMAT_D16_UNORM
	FormatX8D24UnormPack32    Format = C.VK_FORMAT_X8_D24_UNORM_PACK32
	FormatD32Sfloat           Format = C.VK_FORMAT_D32_SFLOAT
	FormatS8Uint              Format = C.VK_FORMAT_S8_UINT
	FormatD16UnormS8Uint      Format = C.VK_FORMAT_D16_UNORM_S8_UINT
	FormatD24UnormS8Uint      Format = C.VK_FORMAT_D24_UNORM_S8_UINT
	FormatD32SfloatS8Uint     Format = C.VK_FORMAT_D32_SFLOAT_S8_UINT
)

// ImageTiling represents image tiling modes
type ImageTiling int32

const (
	ImageTilingOptimal ImageTiling = C.VK_IMAGE_TILING_OPTIMAL
	ImageTilingLinear  ImageTiling = C.VK_IMAGE_TILING_LINEAR
)

// ImageUsageFlags represents image usage flags
type ImageUsageFlags uint32

const (
	ImageUsageTransferSrcBit            ImageUsageFlags = C.VK_IMAGE_USAGE_TRANSFER_SRC_BIT
	ImageUsageTransferDstBit            ImageUsageFlags = C.VK_IMAGE_USAGE_TRANSFER_DST_BIT
	ImageUsageSampledBit                ImageUsageFlags = C.VK_IMAGE_USAGE_SAMPLED_BIT
	ImageUsageStorageBit                ImageUsageFlags = C.VK_IMAGE_USAGE_STORAGE_BIT
	ImageUsageColorAttachmentBit        ImageUsageFlags = C.VK_IMAGE_USAGE_COLOR_ATTACHMENT_BIT
	ImageUsageDepthStencilAttachmentBit ImageUsageFlags = C.VK_IMAGE_USAGE_DEPTH_STENCIL_ATTACHMENT_BIT
	ImageUsageTransientAttachmentBit    ImageUsageFlags = C.VK_IMAGE_USAGE_TRANSIENT_ATTACHMENT_BIT
	ImageUsageInputAttachmentBit        ImageUsageFlags = C.VK_IMAGE_USAGE_INPUT_ATTACHMENT_BIT
)

// ImageLayout represents image layouts
type ImageLayout int32

const (
	ImageLayoutUndefined                     ImageLayout = C.VK_IMAGE_LAYOUT_UNDEFINED
	ImageLayoutGeneral                       ImageLayout = C.VK_IMAGE_LAYOUT_GENERAL
	ImageLayoutColorAttachmentOptimal        ImageLayout = C.VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL
	ImageLayoutDepthStencilAttachmentOptimal ImageLayout = C.VK_IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL
	ImageLayoutDepthStencilReadOnlyOptimal   ImageLayout = C.VK_IMAGE_LAYOUT_DEPTH_STENCIL_READ_ONLY_OPTIMAL
	ImageLayoutShaderReadOnlyOptimal         ImageLayout = C.VK_IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL
	ImageLayoutTransferSrcOptimal            ImageLayout = C.VK_IMAGE_LAYOUT_TRANSFER_SRC_OPTIMAL
	ImageLayoutTransferDstOptimal            ImageLayout = C.VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL
	ImageLayoutPreinitialized                ImageLayout = C.VK_IMAGE_LAYOUT_PREINITIALIZED
	ImageLayoutPresentSrcKHR                 ImageLayout = C.VK_IMAGE_LAYOUT_PRESENT_SRC_KHR
)

// CreateBuffer creates a buffer
func CreateBuffer(device Device, createInfo *BufferCreateInfo) (Buffer, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	// Validate buffer size
	if createInfo.Size == 0 {
		return nil, NewValidationError("Size", "buffer size cannot be zero")
	}

	// Check for reasonable size limits (1GB limit for safety)
	const maxBufferSize = DeviceSize(1024 * 1024 * 1024)
	if createInfo.Size > maxBufferSize {
		return nil, NewValidationError("Size", "buffer size exceeds reasonable limit of 1GB")
	}

	// Validate usage flags (must have at least one usage bit set)
	if createInfo.Usage == 0 {
		return nil, NewValidationError("Usage", "buffer usage flags cannot be zero")
	}

	var cCreateInfo C.VkBufferCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkBufferCreateFlags(createInfo.Flags)
	cCreateInfo.size = C.VkDeviceSize(createInfo.Size)
	cCreateInfo.usage = C.VkBufferUsageFlags(createInfo.Usage)
	cCreateInfo.sharingMode = C.VkSharingMode(createInfo.SharingMode)
	cCreateInfo.queueFamilyIndexCount = 0
	cCreateInfo.pQueueFamilyIndices = nil

	var buffer C.VkBuffer
	result := Result(C.vkCreateBuffer(C.VkDevice(device), &cCreateInfo, nil, &buffer))
	if result != Success {
		return nil, NewVulkanError(result, "CreateBuffer", "Vulkan buffer creation failed")
	}

	return Buffer(buffer), nil
}

// DestroyBuffer destroys a buffer
func DestroyBuffer(device Device, buffer Buffer) {
	C.vkDestroyBuffer(C.VkDevice(device), C.VkBuffer(buffer), nil)
}

// GetBufferMemoryRequirements gets buffer memory requirements
func GetBufferMemoryRequirements(device Device, buffer Buffer) MemoryRequirements {
	var cReqs C.VkMemoryRequirements
	C.vkGetBufferMemoryRequirements(C.VkDevice(device), C.VkBuffer(buffer), &cReqs)

	return MemoryRequirements{
		Size:           DeviceSize(cReqs.size),
		Alignment:      DeviceSize(cReqs.alignment),
		MemoryTypeBits: uint32(cReqs.memoryTypeBits),
	}
}

// AllocateMemory allocates device memory
func AllocateMemory(device Device, allocateInfo *MemoryAllocateInfo) (DeviceMemory, error) {
	var cAllocateInfo C.VkMemoryAllocateInfo
	cAllocateInfo.sType = C.VK_STRUCTURE_TYPE_MEMORY_ALLOCATE_INFO
	cAllocateInfo.pNext = nil
	cAllocateInfo.allocationSize = C.VkDeviceSize(allocateInfo.AllocationSize)
	cAllocateInfo.memoryTypeIndex = C.uint32_t(allocateInfo.MemoryTypeIndex)

	var memory C.VkDeviceMemory
	result := Result(C.vkAllocateMemory(C.VkDevice(device), &cAllocateInfo, nil, &memory))
	if result != Success {
		return nil, result
	}

	return DeviceMemory(memory), nil
}

// FreeMemory frees device memory
func FreeMemory(device Device, memory DeviceMemory) {
	C.vkFreeMemory(C.VkDevice(device), C.VkDeviceMemory(memory), nil)
}

// BindBufferMemory binds buffer memory
func BindBufferMemory(device Device, buffer Buffer, memory DeviceMemory, memoryOffset DeviceSize) error {
	result := Result(C.vkBindBufferMemory(C.VkDevice(device), C.VkBuffer(buffer), C.VkDeviceMemory(memory), C.VkDeviceSize(memoryOffset)))
	if result != Success {
		return result
	}
	return nil
}

// MapMemory maps device memory
func MapMemory(device Device, memory DeviceMemory, offset, size DeviceSize, flags uint32) (unsafe.Pointer, error) {
	var data unsafe.Pointer
	result := Result(C.vkMapMemory(C.VkDevice(device), C.VkDeviceMemory(memory), C.VkDeviceSize(offset), C.VkDeviceSize(size), C.VkMemoryMapFlags(flags), &data))
	if result != Success {
		return nil, result
	}
	return data, nil
}

// UnmapMemory unmaps device memory
func UnmapMemory(device Device, memory DeviceMemory) {
	C.vkUnmapMemory(C.VkDevice(device), C.VkDeviceMemory(memory))
}

// CreateImage creates an image
func CreateImage(device Device, createInfo *ImageCreateInfo) (Image, error) {
	var cCreateInfo C.VkImageCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_IMAGE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkImageCreateFlags(createInfo.Flags)
	cCreateInfo.imageType = C.VkImageType(createInfo.ImageType)
	cCreateInfo.format = C.VkFormat(createInfo.Format)
	cCreateInfo.extent.width = C.uint32_t(createInfo.Extent.Width)
	cCreateInfo.extent.height = C.uint32_t(createInfo.Extent.Height)
	cCreateInfo.extent.depth = C.uint32_t(createInfo.Extent.Depth)
	cCreateInfo.mipLevels = C.uint32_t(createInfo.MipLevels)
	cCreateInfo.arrayLayers = C.uint32_t(createInfo.ArrayLayers)
	cCreateInfo.samples = C.VkSampleCountFlagBits(createInfo.Samples)
	cCreateInfo.tiling = C.VkImageTiling(createInfo.Tiling)
	cCreateInfo.usage = C.VkImageUsageFlags(createInfo.Usage)
	cCreateInfo.sharingMode = C.VkSharingMode(createInfo.SharingMode)
	cCreateInfo.queueFamilyIndexCount = 0
	cCreateInfo.pQueueFamilyIndices = nil
	cCreateInfo.initialLayout = C.VkImageLayout(createInfo.InitialLayout)

	var image C.VkImage
	result := Result(C.vkCreateImage(C.VkDevice(device), &cCreateInfo, nil, &image))
	if result != Success {
		return nil, result
	}

	return Image(image), nil
}

// DestroyImage destroys an image
func DestroyImage(device Device, image Image) {
	C.vkDestroyImage(C.VkDevice(device), C.VkImage(image), nil)
}

// GetImageMemoryRequirements gets image memory requirements
func GetImageMemoryRequirements(device Device, image Image) MemoryRequirements {
	var cReqs C.VkMemoryRequirements
	C.vkGetImageMemoryRequirements(C.VkDevice(device), C.VkImage(image), &cReqs)

	return MemoryRequirements{
		Size:           DeviceSize(cReqs.size),
		Alignment:      DeviceSize(cReqs.alignment),
		MemoryTypeBits: uint32(cReqs.memoryTypeBits),
	}
}

// BindImageMemory binds image memory
func BindImageMemory(device Device, image Image, memory DeviceMemory, memoryOffset DeviceSize) error {
	result := Result(C.vkBindImageMemory(C.VkDevice(device), C.VkImage(image), C.VkDeviceMemory(memory), C.VkDeviceSize(memoryOffset)))
	if result != Success {
		return result
	}
	return nil
}

// FindMemoryType finds a suitable memory type
func FindMemoryType(memProperties PhysicalDeviceMemoryProperties, typeFilter uint32, properties MemoryPropertyFlags) (uint32, bool) {
	for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
		if (typeFilter&(1<<i)) != 0 && (memProperties.MemoryTypes[i].PropertyFlags&properties) == properties {
			return i, true
		}
	}
	return 0, false
}

// ============================================================================
// Memory Management Enhancements
// ============================================================================

// MappedMemoryRange describes a mapped memory range for flush/invalidate operations
type MappedMemoryRange struct {
	Memory DeviceMemory
	Offset DeviceSize
	Size   DeviceSize
}

// FlushMappedMemoryRanges flushes mapped memory ranges to make host writes visible to device
// This is required for non-coherent memory after the host writes to mapped memory
func FlushMappedMemoryRanges(device Device, memoryRanges []MappedMemoryRange) error {
	// Input validation
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if len(memoryRanges) == 0 {
		return nil // Nothing to flush
	}

	cRanges := make([]C.VkMappedMemoryRange, len(memoryRanges))
	for i, r := range memoryRanges {
		if r.Memory == nil {
			return NewValidationError("Memory", "memory handle cannot be nil")
		}
		cRanges[i].sType = C.VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE
		cRanges[i].pNext = nil
		cRanges[i].memory = C.VkDeviceMemory(r.Memory)
		cRanges[i].offset = C.VkDeviceSize(r.Offset)
		cRanges[i].size = C.VkDeviceSize(r.Size)
	}

	result := Result(C.vkFlushMappedMemoryRanges(C.VkDevice(device), C.uint32_t(len(cRanges)), &cRanges[0]))
	if result != Success {
		return NewVulkanError(result, "FlushMappedMemoryRanges", "Vulkan flush mapped memory ranges failed")
	}

	return nil
}

// InvalidateMappedMemoryRanges invalidates mapped memory ranges to make device writes visible to host
// This is required for non-coherent memory before the host reads from mapped memory
func InvalidateMappedMemoryRanges(device Device, memoryRanges []MappedMemoryRange) error {
	// Input validation
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if len(memoryRanges) == 0 {
		return nil // Nothing to invalidate
	}

	cRanges := make([]C.VkMappedMemoryRange, len(memoryRanges))
	for i, r := range memoryRanges {
		if r.Memory == nil {
			return NewValidationError("Memory", "memory handle cannot be nil")
		}
		cRanges[i].sType = C.VK_STRUCTURE_TYPE_MAPPED_MEMORY_RANGE
		cRanges[i].pNext = nil
		cRanges[i].memory = C.VkDeviceMemory(r.Memory)
		cRanges[i].offset = C.VkDeviceSize(r.Offset)
		cRanges[i].size = C.VkDeviceSize(r.Size)
	}

	result := Result(C.vkInvalidateMappedMemoryRanges(C.VkDevice(device), C.uint32_t(len(cRanges)), &cRanges[0]))
	if result != Success {
		return NewVulkanError(result, "InvalidateMappedMemoryRanges", "Vulkan invalidate mapped memory ranges failed")
	}

	return nil
}

// MemoryUsage represents common memory usage patterns for automatic memory type selection
type MemoryUsage int

const (
	// MemoryUsageGPUOnly - Memory that is only accessible by the GPU (fastest for GPU operations)
	MemoryUsageGPUOnly MemoryUsage = iota
	// MemoryUsageCPUOnly - Memory that is only accessible by the CPU (for staging)
	MemoryUsageCPUOnly
	// MemoryUsageCPUToGPU - Memory for CPU-to-GPU data transfer (upload)
	MemoryUsageCPUToGPU
	// MemoryUsageGPUToCPU - Memory for GPU-to-CPU data transfer (readback)
	MemoryUsageGPUToCPU
)

// FindMemoryTypeForUsage finds a suitable memory type based on common usage patterns
// This provides automatic memory type selection for common use cases
func FindMemoryTypeForUsage(memProperties PhysicalDeviceMemoryProperties, typeFilter uint32, usage MemoryUsage) (uint32, bool) {
	var requiredFlags MemoryPropertyFlags
	var preferredFlags MemoryPropertyFlags

	switch usage {
	case MemoryUsageGPUOnly:
		// Device local memory for GPU-only access (textures, GPU-only buffers)
		requiredFlags = MemoryPropertyDeviceLocalBit
		preferredFlags = 0
	case MemoryUsageCPUOnly:
		// Host visible + coherent for CPU-only staging buffers
		requiredFlags = MemoryPropertyHostVisibleBit | MemoryPropertyHostCoherentBit
		preferredFlags = 0
	case MemoryUsageCPUToGPU:
		// Host visible + coherent, prefer device local for unified memory architectures
		requiredFlags = MemoryPropertyHostVisibleBit | MemoryPropertyHostCoherentBit
		preferredFlags = MemoryPropertyDeviceLocalBit
	case MemoryUsageGPUToCPU:
		// Host visible + cached for efficient readback
		requiredFlags = MemoryPropertyHostVisibleBit
		preferredFlags = MemoryPropertyHostCachedBit
	default:
		return 0, false
	}

	// First try to find memory with both required and preferred flags
	if preferredFlags != 0 {
		for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
			flags := memProperties.MemoryTypes[i].PropertyFlags
			if (typeFilter&(1<<i)) != 0 && (flags&(requiredFlags|preferredFlags)) == (requiredFlags|preferredFlags) {
				return i, true
			}
		}
	}

	// Fall back to just required flags
	for i := uint32(0); i < memProperties.MemoryTypeCount; i++ {
		if (typeFilter&(1<<i)) != 0 && (memProperties.MemoryTypes[i].PropertyFlags&requiredFlags) == requiredFlags {
			return i, true
		}
	}

	return 0, false
}

// StagingBuffer represents a staging buffer for host-to-device transfers
type StagingBuffer struct {
	Buffer Buffer
	Memory DeviceMemory
	Size   DeviceSize
	Data   unsafe.Pointer // Mapped pointer (nil if not mapped)
}

// CreateStagingBuffer creates a staging buffer for host-to-device transfers
// The buffer is created with TRANSFER_SRC usage and host-visible, coherent memory
func CreateStagingBuffer(device Device, physicalDevice PhysicalDevice, size DeviceSize) (*StagingBuffer, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if size == 0 {
		return nil, NewValidationError("size", "cannot be zero")
	}

	// Create the staging buffer
	buffer, err := CreateBuffer(device, &BufferCreateInfo{
		Size:        size,
		Usage:       BufferUsageTransferSrcBit,
		SharingMode: SharingModeExclusive,
	})
	if err != nil {
		return nil, err
	}

	// Get memory requirements
	memReqs := GetBufferMemoryRequirements(device, buffer)

	// Find suitable memory type
	memProps := GetPhysicalDeviceMemoryProperties(physicalDevice)
	memTypeIndex, found := FindMemoryTypeForUsage(memProps, memReqs.MemoryTypeBits, MemoryUsageCPUToGPU)
	if !found {
		DestroyBuffer(device, buffer)
		return nil, NewValidationError("memory", "no suitable memory type found for staging buffer")
	}

	// Allocate memory
	memory, err := AllocateMemory(device, &MemoryAllocateInfo{
		AllocationSize:  memReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	})
	if err != nil {
		DestroyBuffer(device, buffer)
		return nil, err
	}

	// Bind memory to buffer
	err = BindBufferMemory(device, buffer, memory, 0)
	if err != nil {
		FreeMemory(device, memory)
		DestroyBuffer(device, buffer)
		return nil, err
	}

	// Map the memory for easy access
	data, err := MapMemory(device, memory, 0, size, 0)
	if err != nil {
		FreeMemory(device, memory)
		DestroyBuffer(device, buffer)
		return nil, err
	}

	return &StagingBuffer{
		Buffer: buffer,
		Memory: memory,
		Size:   size,
		Data:   data,
	}, nil
}

// DestroyStagingBuffer destroys a staging buffer and frees its memory
func DestroyStagingBuffer(device Device, stagingBuffer *StagingBuffer) {
	if device == nil || stagingBuffer == nil {
		return
	}

	if stagingBuffer.Data != nil {
		UnmapMemory(device, stagingBuffer.Memory)
	}
	if stagingBuffer.Memory != nil {
		FreeMemory(device, stagingBuffer.Memory)
	}
	if stagingBuffer.Buffer != nil {
		DestroyBuffer(device, stagingBuffer.Buffer)
	}
}

// CopyDataToStagingBuffer copies data to a staging buffer
func CopyDataToStagingBuffer(stagingBuffer *StagingBuffer, data []byte) error {
	if stagingBuffer == nil {
		return NewValidationError("stagingBuffer", "cannot be nil")
	}
	if stagingBuffer.Data == nil {
		return NewValidationError("stagingBuffer.Data", "buffer is not mapped")
	}
	if len(data) == 0 {
		return nil // Nothing to copy
	}
	if DeviceSize(len(data)) > stagingBuffer.Size {
		return NewValidationError("data", "data size exceeds staging buffer size")
	}

	// Copy data to the mapped memory
	C.memcpy(stagingBuffer.Data, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	return nil
}

// DefaultMemoryAlignment is the default alignment for memory pool allocations
const DefaultMemoryAlignment DeviceSize = 256

// MemoryPool represents a simple memory pool for efficient allocations
type MemoryPool struct {
	Device          Device
	Memory          DeviceMemory
	Size            DeviceSize
	MemoryTypeIndex uint32
	Offset          DeviceSize // Current allocation offset
	Alignment       DeviceSize // Minimum allocation alignment
}

// isPowerOfTwo checks if a value is a power of two
func isPowerOfTwo(n DeviceSize) bool {
	return n > 0 && (n&(n-1)) == 0
}

// CreateMemoryPool creates a memory pool for efficient sub-allocations
func CreateMemoryPool(device Device, size DeviceSize, memoryTypeIndex uint32, alignment DeviceSize) (*MemoryPool, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if size == 0 {
		return nil, NewValidationError("size", "cannot be zero")
	}
	if alignment == 0 {
		alignment = DefaultMemoryAlignment
	}
	if !isPowerOfTwo(alignment) {
		return nil, NewValidationError("alignment", "must be a power of two")
	}

	// Allocate the pool memory
	memory, err := AllocateMemory(device, &MemoryAllocateInfo{
		AllocationSize:  size,
		MemoryTypeIndex: memoryTypeIndex,
	})
	if err != nil {
		return nil, err
	}

	return &MemoryPool{
		Device:          device,
		Memory:          memory,
		Size:            size,
		MemoryTypeIndex: memoryTypeIndex,
		Offset:          0,
		Alignment:       alignment,
	}, nil
}

// Allocate allocates memory from the pool
// Returns the offset within the pool memory, or an error if there's not enough space
func (pool *MemoryPool) Allocate(size DeviceSize, alignment DeviceSize) (DeviceSize, error) {
	if pool == nil {
		return 0, NewValidationError("pool", "cannot be nil")
	}
	if size == 0 {
		return 0, NewValidationError("size", "cannot be zero")
	}

	// Use pool's default alignment if not specified
	if alignment == 0 {
		alignment = pool.Alignment
	}
	if !isPowerOfTwo(alignment) {
		return 0, NewValidationError("alignment", "must be a power of two")
	}

	// Align the current offset
	alignedOffset := (pool.Offset + alignment - 1) & ^(alignment - 1)

	// Check if there's enough space
	if alignedOffset+size > pool.Size {
		return 0, NewValidationError("size", "not enough space in memory pool")
	}

	// Update the pool offset
	pool.Offset = alignedOffset + size

	return alignedOffset, nil
}

// Reset resets the pool for reuse (does not free memory)
func (pool *MemoryPool) Reset() {
	if pool != nil {
		pool.Offset = 0
	}
}

// Destroy destroys the memory pool and frees its memory
func (pool *MemoryPool) Destroy() {
	if pool != nil && pool.Device != nil && pool.Memory != nil {
		FreeMemory(pool.Device, pool.Memory)
		pool.Memory = nil
	}
}
