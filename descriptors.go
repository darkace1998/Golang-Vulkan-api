package vulkan

/*
#include <vulkan/vulkan.h>
*/
import "C"

// ImageViewCreateInfo contains image view creation information
type ImageViewCreateInfo struct {
	Image            Image
	ViewType         ImageViewType
	Format           Format
	SubresourceRange ImageSubresourceRange
}

// ImageViewType represents image view types
type ImageViewType int32

const (
	ImageViewType1D        ImageViewType = C.VK_IMAGE_VIEW_TYPE_1D
	ImageViewType2D        ImageViewType = C.VK_IMAGE_VIEW_TYPE_2D
	ImageViewType3D        ImageViewType = C.VK_IMAGE_VIEW_TYPE_3D
	ImageViewTypeCube      ImageViewType = C.VK_IMAGE_VIEW_TYPE_CUBE
	ImageViewType1DArray   ImageViewType = C.VK_IMAGE_VIEW_TYPE_1D_ARRAY
	ImageViewType2DArray   ImageViewType = C.VK_IMAGE_VIEW_TYPE_2D_ARRAY
	ImageViewTypeCubeArray ImageViewType = C.VK_IMAGE_VIEW_TYPE_CUBE_ARRAY
)

// ImageSubresourceRange describes an image subresource range
type ImageSubresourceRange struct {
	AspectMask     ImageAspectFlags
	BaseMipLevel   uint32
	LevelCount     uint32
	BaseArrayLayer uint32
	LayerCount     uint32
}

// ImageAspectFlags represents image aspect flags
type ImageAspectFlags uint32

const (
	ImageAspectColorBit   ImageAspectFlags = C.VK_IMAGE_ASPECT_COLOR_BIT
	ImageAspectDepthBit   ImageAspectFlags = C.VK_IMAGE_ASPECT_DEPTH_BIT
	ImageAspectStencilBit ImageAspectFlags = C.VK_IMAGE_ASPECT_STENCIL_BIT
)

// SamplerCreateInfo contains sampler creation information
type SamplerCreateInfo struct {
	MagFilter    Filter
	MinFilter    Filter
	AddressModeU SamplerAddressMode
	AddressModeV SamplerAddressMode
	AddressModeW SamplerAddressMode
}

// Filter represents texture filtering modes
type Filter int32

const (
	FilterNearest Filter = C.VK_FILTER_NEAREST
	FilterLinear  Filter = C.VK_FILTER_LINEAR
)

// SamplerAddressMode represents sampler address modes
type SamplerAddressMode int32

const (
	SamplerAddressModeRepeat            SamplerAddressMode = C.VK_SAMPLER_ADDRESS_MODE_REPEAT
	SamplerAddressModeMirroredRepeat    SamplerAddressMode = C.VK_SAMPLER_ADDRESS_MODE_MIRRORED_REPEAT
	SamplerAddressModeClampToEdge       SamplerAddressMode = C.VK_SAMPLER_ADDRESS_MODE_CLAMP_TO_EDGE
	SamplerAddressModeClampToBorder     SamplerAddressMode = C.VK_SAMPLER_ADDRESS_MODE_CLAMP_TO_BORDER
	SamplerAddressModeMirrorClampToEdge SamplerAddressMode = C.VK_SAMPLER_ADDRESS_MODE_MIRROR_CLAMP_TO_EDGE
)

// DescriptorSetLayoutCreateInfo contains descriptor set layout creation information
type DescriptorSetLayoutCreateInfo struct {
	Bindings []DescriptorSetLayoutBinding
}

// DescriptorSetLayoutBinding describes a descriptor set layout binding
type DescriptorSetLayoutBinding struct {
	Binding         uint32
	DescriptorType  DescriptorType
	DescriptorCount uint32
	StageFlags      ShaderStageFlags
}

// DescriptorType represents descriptor types
type DescriptorType int32

const (
	DescriptorTypeSampler              DescriptorType = C.VK_DESCRIPTOR_TYPE_SAMPLER
	DescriptorTypeCombinedImageSampler DescriptorType = C.VK_DESCRIPTOR_TYPE_COMBINED_IMAGE_SAMPLER
	DescriptorTypeSampledImage         DescriptorType = C.VK_DESCRIPTOR_TYPE_SAMPLED_IMAGE
	DescriptorTypeStorageImage         DescriptorType = C.VK_DESCRIPTOR_TYPE_STORAGE_IMAGE
	DescriptorTypeUniformTexelBuffer   DescriptorType = C.VK_DESCRIPTOR_TYPE_UNIFORM_TEXEL_BUFFER
	DescriptorTypeStorageTexelBuffer   DescriptorType = C.VK_DESCRIPTOR_TYPE_STORAGE_TEXEL_BUFFER
	DescriptorTypeUniformBuffer        DescriptorType = C.VK_DESCRIPTOR_TYPE_UNIFORM_BUFFER
	DescriptorTypeStorageBuffer        DescriptorType = C.VK_DESCRIPTOR_TYPE_STORAGE_BUFFER
	DescriptorTypeUniformBufferDynamic DescriptorType = C.VK_DESCRIPTOR_TYPE_UNIFORM_BUFFER_DYNAMIC
	DescriptorTypeStorageBufferDynamic DescriptorType = C.VK_DESCRIPTOR_TYPE_STORAGE_BUFFER_DYNAMIC
	DescriptorTypeInputAttachment      DescriptorType = C.VK_DESCRIPTOR_TYPE_INPUT_ATTACHMENT
)

// DescriptorPoolCreateInfo contains descriptor pool creation information
type DescriptorPoolCreateInfo struct {
	MaxSets   uint32
	PoolSizes []DescriptorPoolSize
}

// DescriptorPoolSize describes a descriptor pool size
type DescriptorPoolSize struct {
	Type            DescriptorType
	DescriptorCount uint32
}

// CreateImageView creates an image view
func CreateImageView(device Device, createInfo *ImageViewCreateInfo) (ImageView, error) {
	var cCreateInfo C.VkImageViewCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_IMAGE_VIEW_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.image = C.VkImage(createInfo.Image)
	cCreateInfo.viewType = C.VkImageViewType(createInfo.ViewType)
	cCreateInfo.format = C.VkFormat(createInfo.Format)

	// Component mapping (identity)
	cCreateInfo.components.r = C.VK_COMPONENT_SWIZZLE_IDENTITY
	cCreateInfo.components.g = C.VK_COMPONENT_SWIZZLE_IDENTITY
	cCreateInfo.components.b = C.VK_COMPONENT_SWIZZLE_IDENTITY
	cCreateInfo.components.a = C.VK_COMPONENT_SWIZZLE_IDENTITY

	// Subresource range
	cCreateInfo.subresourceRange.aspectMask = C.VkImageAspectFlags(createInfo.SubresourceRange.AspectMask)
	cCreateInfo.subresourceRange.baseMipLevel = C.uint32_t(createInfo.SubresourceRange.BaseMipLevel)
	cCreateInfo.subresourceRange.levelCount = C.uint32_t(createInfo.SubresourceRange.LevelCount)
	cCreateInfo.subresourceRange.baseArrayLayer = C.uint32_t(createInfo.SubresourceRange.BaseArrayLayer)
	cCreateInfo.subresourceRange.layerCount = C.uint32_t(createInfo.SubresourceRange.LayerCount)

	var imageView C.VkImageView
	result := Result(C.vkCreateImageView(C.VkDevice(device), &cCreateInfo, nil, &imageView))
	if result != Success {
		return nil, result
	}

	return ImageView(imageView), nil
}

// DestroyImageView destroys an image view
func DestroyImageView(device Device, imageView ImageView) {
	C.vkDestroyImageView(C.VkDevice(device), C.VkImageView(imageView), nil)
}

// CreateSampler creates a sampler
func CreateSampler(device Device, createInfo *SamplerCreateInfo) (Sampler, error) {
	var cCreateInfo C.VkSamplerCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_SAMPLER_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.magFilter = C.VkFilter(createInfo.MagFilter)
	cCreateInfo.minFilter = C.VkFilter(createInfo.MinFilter)
	cCreateInfo.mipmapMode = C.VK_SAMPLER_MIPMAP_MODE_LINEAR
	cCreateInfo.addressModeU = C.VkSamplerAddressMode(createInfo.AddressModeU)
	cCreateInfo.addressModeV = C.VkSamplerAddressMode(createInfo.AddressModeV)
	cCreateInfo.addressModeW = C.VkSamplerAddressMode(createInfo.AddressModeW)
	cCreateInfo.mipLodBias = 0.0
	cCreateInfo.anisotropyEnable = C.VK_FALSE
	cCreateInfo.maxAnisotropy = 1.0
	cCreateInfo.compareEnable = C.VK_FALSE
	cCreateInfo.compareOp = C.VK_COMPARE_OP_ALWAYS
	cCreateInfo.minLod = 0.0
	cCreateInfo.maxLod = 0.0
	cCreateInfo.borderColor = C.VK_BORDER_COLOR_INT_OPAQUE_BLACK
	cCreateInfo.unnormalizedCoordinates = C.VK_FALSE

	var sampler C.VkSampler
	result := Result(C.vkCreateSampler(C.VkDevice(device), &cCreateInfo, nil, &sampler))
	if result != Success {
		return nil, result
	}

	return Sampler(sampler), nil
}

// DestroySampler destroys a sampler
func DestroySampler(device Device, sampler Sampler) {
	C.vkDestroySampler(C.VkDevice(device), C.VkSampler(sampler), nil)
}

// CreateDescriptorSetLayout creates a descriptor set layout
func CreateDescriptorSetLayout(device Device, createInfo *DescriptorSetLayoutCreateInfo) (DescriptorSetLayout, error) {
	var cCreateInfo C.VkDescriptorSetLayoutCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_DESCRIPTOR_SET_LAYOUT_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0

	var cBindings []C.VkDescriptorSetLayoutBinding
	if len(createInfo.Bindings) > 0 {
		cBindings = make([]C.VkDescriptorSetLayoutBinding, len(createInfo.Bindings))
		for i, binding := range createInfo.Bindings {
			cBindings[i].binding = C.uint32_t(binding.Binding)
			cBindings[i].descriptorType = C.VkDescriptorType(binding.DescriptorType)
			cBindings[i].descriptorCount = C.uint32_t(binding.DescriptorCount)
			cBindings[i].stageFlags = C.VkShaderStageFlags(binding.StageFlags)
			cBindings[i].pImmutableSamplers = nil
		}
		cCreateInfo.bindingCount = C.uint32_t(len(cBindings))
		cCreateInfo.pBindings = &cBindings[0]
	}

	var layout C.VkDescriptorSetLayout
	result := Result(C.vkCreateDescriptorSetLayout(C.VkDevice(device), &cCreateInfo, nil, &layout))
	if result != Success {
		return nil, result
	}

	return DescriptorSetLayout(layout), nil
}

// DestroyDescriptorSetLayout destroys a descriptor set layout
func DestroyDescriptorSetLayout(device Device, layout DescriptorSetLayout) {
	C.vkDestroyDescriptorSetLayout(C.VkDevice(device), C.VkDescriptorSetLayout(layout), nil)
}

// CreateDescriptorPool creates a descriptor pool
func CreateDescriptorPool(device Device, createInfo *DescriptorPoolCreateInfo) (DescriptorPool, error) {
	var cCreateInfo C.VkDescriptorPoolCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_DESCRIPTOR_POOL_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.maxSets = C.uint32_t(createInfo.MaxSets)

	var cPoolSizes []C.VkDescriptorPoolSize
	if len(createInfo.PoolSizes) > 0 {
		cPoolSizes = make([]C.VkDescriptorPoolSize, len(createInfo.PoolSizes))
		for i, poolSize := range createInfo.PoolSizes {
			cPoolSizes[i]._type = C.VkDescriptorType(poolSize.Type)
			cPoolSizes[i].descriptorCount = C.uint32_t(poolSize.DescriptorCount)
		}
		cCreateInfo.poolSizeCount = C.uint32_t(len(cPoolSizes))
		cCreateInfo.pPoolSizes = &cPoolSizes[0]
	}

	var pool C.VkDescriptorPool
	result := Result(C.vkCreateDescriptorPool(C.VkDevice(device), &cCreateInfo, nil, &pool))
	if result != Success {
		return nil, result
	}

	return DescriptorPool(pool), nil
}

// DestroyDescriptorPool destroys a descriptor pool
func DestroyDescriptorPool(device Device, pool DescriptorPool) {
	C.vkDestroyDescriptorPool(C.VkDevice(device), C.VkDescriptorPool(pool), nil)
}
