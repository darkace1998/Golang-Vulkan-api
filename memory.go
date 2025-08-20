package vulkan

/*
#cgo pkg-config: vulkan
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
	Size        DeviceSize
	Usage       BufferUsageFlags
	SharingMode SharingMode
}

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

// Format represents pixel formats
type Format int32

const (
	FormatUndefined                Format = C.VK_FORMAT_UNDEFINED
	FormatR4G4UnormPack8           Format = C.VK_FORMAT_R4G4_UNORM_PACK8
	FormatR4G4B4A4UnormPack16      Format = C.VK_FORMAT_R4G4B4A4_UNORM_PACK16
	FormatB4G4R4A4UnormPack16      Format = C.VK_FORMAT_B4G4R4A4_UNORM_PACK16
	FormatR5G6B5UnormPack16        Format = C.VK_FORMAT_R5G6B5_UNORM_PACK16
	FormatB5G6R5UnormPack16        Format = C.VK_FORMAT_B5G6R5_UNORM_PACK16
	FormatR5G5B5A1UnormPack16      Format = C.VK_FORMAT_R5G5B5A1_UNORM_PACK16
	FormatB5G5R5A1UnormPack16      Format = C.VK_FORMAT_B5G5R5A1_UNORM_PACK16
	FormatA1R5G5B5UnormPack16      Format = C.VK_FORMAT_A1R5G5B5_UNORM_PACK16
	FormatR8Unorm                  Format = C.VK_FORMAT_R8_UNORM
	FormatR8Snorm                  Format = C.VK_FORMAT_R8_SNORM
	FormatR8Uscaled                Format = C.VK_FORMAT_R8_USCALED
	FormatR8Sscaled                Format = C.VK_FORMAT_R8_SSCALED
	FormatR8Uint                   Format = C.VK_FORMAT_R8_UINT
	FormatR8Sint                   Format = C.VK_FORMAT_R8_SINT
	FormatR8Srgb                   Format = C.VK_FORMAT_R8_SRGB
	FormatR8G8Unorm                Format = C.VK_FORMAT_R8G8_UNORM
	FormatR8G8Snorm                Format = C.VK_FORMAT_R8G8_SNORM
	FormatR8G8Uscaled              Format = C.VK_FORMAT_R8G8_USCALED
	FormatR8G8Sscaled              Format = C.VK_FORMAT_R8G8_SSCALED
	FormatR8G8Uint                 Format = C.VK_FORMAT_R8G8_UINT
	FormatR8G8Sint                 Format = C.VK_FORMAT_R8G8_SINT
	FormatR8G8Srgb                 Format = C.VK_FORMAT_R8G8_SRGB
	FormatR8G8B8Unorm              Format = C.VK_FORMAT_R8G8B8_UNORM
	FormatR8G8B8Snorm              Format = C.VK_FORMAT_R8G8B8_SNORM
	FormatR8G8B8Uscaled            Format = C.VK_FORMAT_R8G8B8_USCALED
	FormatR8G8B8Sscaled            Format = C.VK_FORMAT_R8G8B8_SSCALED
	FormatR8G8B8Uint               Format = C.VK_FORMAT_R8G8B8_UINT
	FormatR8G8B8Sint               Format = C.VK_FORMAT_R8G8B8_SINT
	FormatR8G8B8Srgb               Format = C.VK_FORMAT_R8G8B8_SRGB
	FormatB8G8R8Unorm              Format = C.VK_FORMAT_B8G8R8_UNORM
	FormatB8G8R8Snorm              Format = C.VK_FORMAT_B8G8R8_SNORM
	FormatB8G8R8Uscaled            Format = C.VK_FORMAT_B8G8R8_USCALED
	FormatB8G8R8Sscaled            Format = C.VK_FORMAT_B8G8R8_SSCALED
	FormatB8G8R8Uint               Format = C.VK_FORMAT_B8G8R8_UINT
	FormatB8G8R8Sint               Format = C.VK_FORMAT_B8G8R8_SINT
	FormatB8G8R8Srgb               Format = C.VK_FORMAT_B8G8R8_SRGB
	FormatR8G8B8A8Unorm            Format = C.VK_FORMAT_R8G8B8A8_UNORM
	FormatR8G8B8A8Snorm            Format = C.VK_FORMAT_R8G8B8A8_SNORM
	FormatR8G8B8A8Uscaled          Format = C.VK_FORMAT_R8G8B8A8_USCALED
	FormatR8G8B8A8Sscaled          Format = C.VK_FORMAT_R8G8B8A8_SSCALED
	FormatR8G8B8A8Uint             Format = C.VK_FORMAT_R8G8B8A8_UINT
	FormatR8G8B8A8Sint             Format = C.VK_FORMAT_R8G8B8A8_SINT
	FormatR8G8B8A8Srgb             Format = C.VK_FORMAT_R8G8B8A8_SRGB
	FormatB8G8R8A8Unorm            Format = C.VK_FORMAT_B8G8R8A8_UNORM
	FormatB8G8R8A8Snorm            Format = C.VK_FORMAT_B8G8R8A8_SNORM
	FormatB8G8R8A8Uscaled          Format = C.VK_FORMAT_B8G8R8A8_USCALED
	FormatB8G8R8A8Sscaled          Format = C.VK_FORMAT_B8G8R8A8_SSCALED
	FormatB8G8R8A8Uint             Format = C.VK_FORMAT_B8G8R8A8_UINT
	FormatB8G8R8A8Sint             Format = C.VK_FORMAT_B8G8R8A8_SINT
	FormatB8G8R8A8Srgb             Format = C.VK_FORMAT_B8G8R8A8_SRGB
	FormatD16Unorm                 Format = C.VK_FORMAT_D16_UNORM
	FormatX8D24UnormPack32         Format = C.VK_FORMAT_X8_D24_UNORM_PACK32
	FormatD32Sfloat                Format = C.VK_FORMAT_D32_SFLOAT
	FormatS8Uint                   Format = C.VK_FORMAT_S8_UINT
	FormatD16UnormS8Uint           Format = C.VK_FORMAT_D16_UNORM_S8_UINT
	FormatD24UnormS8Uint           Format = C.VK_FORMAT_D24_UNORM_S8_UINT
	FormatD32SfloatS8Uint          Format = C.VK_FORMAT_D32_SFLOAT_S8_UINT
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
	ImageLayoutUndefined                        ImageLayout = C.VK_IMAGE_LAYOUT_UNDEFINED
	ImageLayoutGeneral                          ImageLayout = C.VK_IMAGE_LAYOUT_GENERAL
	ImageLayoutColorAttachmentOptimal           ImageLayout = C.VK_IMAGE_LAYOUT_COLOR_ATTACHMENT_OPTIMAL
	ImageLayoutDepthStencilAttachmentOptimal    ImageLayout = C.VK_IMAGE_LAYOUT_DEPTH_STENCIL_ATTACHMENT_OPTIMAL
	ImageLayoutDepthStencilReadOnlyOptimal      ImageLayout = C.VK_IMAGE_LAYOUT_DEPTH_STENCIL_READ_ONLY_OPTIMAL
	ImageLayoutShaderReadOnlyOptimal            ImageLayout = C.VK_IMAGE_LAYOUT_SHADER_READ_ONLY_OPTIMAL
	ImageLayoutTransferSrcOptimal               ImageLayout = C.VK_IMAGE_LAYOUT_TRANSFER_SRC_OPTIMAL
	ImageLayoutTransferDstOptimal               ImageLayout = C.VK_IMAGE_LAYOUT_TRANSFER_DST_OPTIMAL
	ImageLayoutPreinitialized                   ImageLayout = C.VK_IMAGE_LAYOUT_PREINITIALIZED
	ImageLayoutPresentSrcKHR                    ImageLayout = C.VK_IMAGE_LAYOUT_PRESENT_SRC_KHR
)

// CreateBuffer creates a buffer
func CreateBuffer(device Device, createInfo *BufferCreateInfo) (Buffer, error) {
	var cCreateInfo C.VkBufferCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.size = C.VkDeviceSize(createInfo.Size)
	cCreateInfo.usage = C.VkBufferUsageFlags(createInfo.Usage)
	cCreateInfo.sharingMode = C.VkSharingMode(createInfo.SharingMode)
	cCreateInfo.queueFamilyIndexCount = 0
	cCreateInfo.pQueueFamilyIndices = nil

	var buffer C.VkBuffer
	result := Result(C.vkCreateBuffer(C.VkDevice(device), &cCreateInfo, nil, &buffer))
	if result != Success {
		return nil, result
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
	cCreateInfo.flags = 0
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