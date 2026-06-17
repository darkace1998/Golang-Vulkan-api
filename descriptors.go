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

// DescriptorPoolCreateFlags represents descriptor pool creation flags
type DescriptorPoolCreateFlags uint32

const (
	DescriptorPoolCreateFreeDescriptorSetBit DescriptorPoolCreateFlags = C.VK_DESCRIPTOR_POOL_CREATE_FREE_DESCRIPTOR_SET_BIT
	DescriptorPoolCreateUpdateAfterBindBit   DescriptorPoolCreateFlags = C.VK_DESCRIPTOR_POOL_CREATE_UPDATE_AFTER_BIND_BIT
)

// DescriptorPoolCreateInfo contains descriptor pool creation information
type DescriptorPoolCreateInfo struct {
	Flags     DescriptorPoolCreateFlags
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
		return nil, NewVulkanError(result, "CreateImageView", "Vulkan image view creation failed")
	}

	return ImageView(imageView), nil
}

// DestroyImageView destroys an image view
func DestroyImageView(device Device, imageView ImageView) {
	if device == nil || imageView == nil {
		return
	}
	C.vkDestroyImageView(C.VkDevice(device), C.VkImageView(imageView), nil)
}

// CreateSampler creates a sampler
func CreateSampler(device Device, createInfo *SamplerCreateInfo) (Sampler, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

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
		return nil, NewVulkanError(result, "CreateSampler", "Vulkan sampler creation failed")
	}

	return Sampler(sampler), nil
}

// DestroySampler destroys a sampler
func DestroySampler(device Device, sampler Sampler) {
	if device == nil || sampler == nil {
		return
	}
	C.vkDestroySampler(C.VkDevice(device), C.VkSampler(sampler), nil)
}

// CreateDescriptorSetLayout creates a descriptor set layout
func CreateDescriptorSetLayout(device Device, createInfo *DescriptorSetLayoutCreateInfo) (DescriptorSetLayout, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

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
		return nil, NewVulkanError(result, "CreateDescriptorSetLayout", "Vulkan descriptor set layout creation failed")
	}

	return DescriptorSetLayout(layout), nil
}

// DestroyDescriptorSetLayout destroys a descriptor set layout
func DestroyDescriptorSetLayout(device Device, layout DescriptorSetLayout) {
	if device == nil || layout == nil {
		return
	}
	C.vkDestroyDescriptorSetLayout(C.VkDevice(device), C.VkDescriptorSetLayout(layout), nil)
}

// CreateDescriptorPool creates a descriptor pool
func CreateDescriptorPool(device Device, createInfo *DescriptorPoolCreateInfo) (DescriptorPool, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	var cCreateInfo C.VkDescriptorPoolCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_DESCRIPTOR_POOL_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkDescriptorPoolCreateFlags(createInfo.Flags)
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
		return nil, NewVulkanError(result, "CreateDescriptorPool", "Vulkan descriptor pool creation failed")
	}

	return DescriptorPool(pool), nil
}

// DestroyDescriptorPool destroys a descriptor pool
func DestroyDescriptorPool(device Device, pool DescriptorPool) {
	if device == nil || pool == nil {
		return
	}
	C.vkDestroyDescriptorPool(C.VkDevice(device), C.VkDescriptorPool(pool), nil)
}

// DescriptorSetAllocateInfo contains descriptor set allocation information
type DescriptorSetAllocateInfo struct {
	DescriptorPool DescriptorPool
	SetLayouts     []DescriptorSetLayout
}

// AllocateDescriptorSets allocates one or more descriptor sets
func AllocateDescriptorSets(device Device, allocateInfo *DescriptorSetAllocateInfo) ([]DescriptorSet, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if allocateInfo == nil {
		return nil, NewValidationError("allocateInfo", "cannot be nil")
	}
	if allocateInfo.DescriptorPool == nil {
		return nil, NewValidationError("DescriptorPool", "cannot be nil")
	}
	if len(allocateInfo.SetLayouts) == 0 {
		return nil, NewValidationError("SetLayouts", "cannot be empty")
	}

	var cAllocateInfo C.VkDescriptorSetAllocateInfo
	cAllocateInfo.sType = C.VK_STRUCTURE_TYPE_DESCRIPTOR_SET_ALLOCATE_INFO
	cAllocateInfo.pNext = nil
	cAllocateInfo.descriptorPool = C.VkDescriptorPool(allocateInfo.DescriptorPool)
	cAllocateInfo.descriptorSetCount = C.uint32_t(len(allocateInfo.SetLayouts))

	cSetLayouts := make([]C.VkDescriptorSetLayout, len(allocateInfo.SetLayouts))
	for i, layout := range allocateInfo.SetLayouts {
		cSetLayouts[i] = C.VkDescriptorSetLayout(layout)
	}
	cAllocateInfo.pSetLayouts = &cSetLayouts[0]

	cDescriptorSets := make([]C.VkDescriptorSet, len(allocateInfo.SetLayouts))
	result := Result(C.vkAllocateDescriptorSets(C.VkDevice(device), &cAllocateInfo, &cDescriptorSets[0]))
	if result != Success {
		return nil, NewVulkanError(result, "AllocateDescriptorSets", "Vulkan descriptor set allocation failed")
	}

	descriptorSets := make([]DescriptorSet, len(cDescriptorSets))
	for i, set := range cDescriptorSets {
		descriptorSets[i] = DescriptorSet(set)
	}

	return descriptorSets, nil
}

// FreeDescriptorSets frees one or more descriptor sets
func FreeDescriptorSets(device Device, descriptorPool DescriptorPool, descriptorSets []DescriptorSet) error {
	// Input validation
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if descriptorPool == nil {
		return NewValidationError("descriptorPool", "cannot be nil")
	}
	if len(descriptorSets) == 0 {
		return nil // Empty array is valid, just do nothing
	}

	cDescriptorSets := make([]C.VkDescriptorSet, len(descriptorSets))
	for i, set := range descriptorSets {
		cDescriptorSets[i] = C.VkDescriptorSet(set)
	}

	result := Result(C.vkFreeDescriptorSets(
		C.VkDevice(device),
		C.VkDescriptorPool(descriptorPool),
		C.uint32_t(len(cDescriptorSets)),
		&cDescriptorSets[0],
	))
	if result != Success {
		return NewVulkanError(result, "FreeDescriptorSets", "Vulkan descriptor set free failed")
	}

	return nil
}

// ResetDescriptorPool resets a descriptor pool
func ResetDescriptorPool(device Device, descriptorPool DescriptorPool) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if descriptorPool == nil {
		return NewValidationError("descriptorPool", "cannot be nil")
	}

	result := Result(C.vkResetDescriptorPool(C.VkDevice(device), C.VkDescriptorPool(descriptorPool), 0))
	if result != Success {
		return NewVulkanError(result, "ResetDescriptorPool", "Vulkan descriptor pool reset failed")
	}

	return nil
}

// DescriptorBufferInfo describes a buffer descriptor
type DescriptorBufferInfo struct {
	Buffer Buffer
	Offset DeviceSize
	Range  DeviceSize
}

// DescriptorImageInfo describes an image descriptor
type DescriptorImageInfo struct {
	Sampler     Sampler
	ImageView   ImageView
	ImageLayout ImageLayout
}

// WriteDescriptorSet describes a descriptor set write operation
type WriteDescriptorSet struct {
	DstSet          DescriptorSet
	DstBinding      uint32
	DstArrayElement uint32
	DescriptorCount uint32
	DescriptorType  DescriptorType
	ImageInfo       []DescriptorImageInfo
	BufferInfo      []DescriptorBufferInfo
	TexelBufferView []BufferView
}

// CopyDescriptorSet describes a descriptor set copy operation
type CopyDescriptorSet struct {
	SrcSet          DescriptorSet
	SrcBinding      uint32
	SrcArrayElement uint32
	DstSet          DescriptorSet
	DstBinding      uint32
	DstArrayElement uint32
	DescriptorCount uint32
}

// UpdateDescriptorSets updates descriptor sets with write and copy operations
// Note: This function follows the Vulkan API which is void and doesn't return errors.
// If device is nil, the function returns early without performing any operation.
func UpdateDescriptorSets(device Device, writes []WriteDescriptorSet, copies []CopyDescriptorSet) {
	if device == nil {
		return // Device is required but vkUpdateDescriptorSets is void, so we silently return
	}

	var cWrites []C.VkWriteDescriptorSet
	var allImageInfos [][]C.VkDescriptorImageInfo
	var allBufferInfos [][]C.VkDescriptorBufferInfo
	var allTexelBufferViews [][]C.VkBufferView

	if len(writes) > 0 {
		cWrites = make([]C.VkWriteDescriptorSet, len(writes))
		allImageInfos = make([][]C.VkDescriptorImageInfo, len(writes))
		allBufferInfos = make([][]C.VkDescriptorBufferInfo, len(writes))
		allTexelBufferViews = make([][]C.VkBufferView, len(writes))

		for i, write := range writes {
			cWrites[i].sType = C.VK_STRUCTURE_TYPE_WRITE_DESCRIPTOR_SET
			cWrites[i].pNext = nil
			cWrites[i].dstSet = C.VkDescriptorSet(write.DstSet)
			cWrites[i].dstBinding = C.uint32_t(write.DstBinding)
			cWrites[i].dstArrayElement = C.uint32_t(write.DstArrayElement)
			cWrites[i].descriptorCount = C.uint32_t(write.DescriptorCount)
			cWrites[i].descriptorType = C.VkDescriptorType(write.DescriptorType)

			// Handle image info
			if len(write.ImageInfo) > 0 {
				allImageInfos[i] = make([]C.VkDescriptorImageInfo, len(write.ImageInfo))
				for j, imgInfo := range write.ImageInfo {
					allImageInfos[i][j].sampler = C.VkSampler(imgInfo.Sampler)
					allImageInfos[i][j].imageView = C.VkImageView(imgInfo.ImageView)
					allImageInfos[i][j].imageLayout = C.VkImageLayout(imgInfo.ImageLayout)
				}
				cWrites[i].pImageInfo = &allImageInfos[i][0]
			}

			// Handle buffer info
			if len(write.BufferInfo) > 0 {
				allBufferInfos[i] = make([]C.VkDescriptorBufferInfo, len(write.BufferInfo))
				for j, bufInfo := range write.BufferInfo {
					allBufferInfos[i][j].buffer = C.VkBuffer(bufInfo.Buffer)
					allBufferInfos[i][j].offset = C.VkDeviceSize(bufInfo.Offset)
					allBufferInfos[i][j]._range = C.VkDeviceSize(bufInfo.Range)
				}
				cWrites[i].pBufferInfo = &allBufferInfos[i][0]
			}

			// Handle texel buffer views
			if len(write.TexelBufferView) > 0 {
				allTexelBufferViews[i] = make([]C.VkBufferView, len(write.TexelBufferView))
				for j, view := range write.TexelBufferView {
					allTexelBufferViews[i][j] = C.VkBufferView(view)
				}
				cWrites[i].pTexelBufferView = &allTexelBufferViews[i][0]
			}
		}
	}

	var cCopies []C.VkCopyDescriptorSet
	if len(copies) > 0 {
		cCopies = make([]C.VkCopyDescriptorSet, len(copies))
		for i, copy := range copies {
			cCopies[i].sType = C.VK_STRUCTURE_TYPE_COPY_DESCRIPTOR_SET
			cCopies[i].pNext = nil
			cCopies[i].srcSet = C.VkDescriptorSet(copy.SrcSet)
			cCopies[i].srcBinding = C.uint32_t(copy.SrcBinding)
			cCopies[i].srcArrayElement = C.uint32_t(copy.SrcArrayElement)
			cCopies[i].dstSet = C.VkDescriptorSet(copy.DstSet)
			cCopies[i].dstBinding = C.uint32_t(copy.DstBinding)
			cCopies[i].dstArrayElement = C.uint32_t(copy.DstArrayElement)
			cCopies[i].descriptorCount = C.uint32_t(copy.DescriptorCount)
		}
	}

	var pWrites *C.VkWriteDescriptorSet
	if len(cWrites) > 0 {
		pWrites = &cWrites[0]
	}

	var pCopies *C.VkCopyDescriptorSet
	if len(cCopies) > 0 {
		pCopies = &cCopies[0]
	}

	C.vkUpdateDescriptorSets(
		C.VkDevice(device),
		C.uint32_t(len(cWrites)),
		pWrites,
		C.uint32_t(len(cCopies)),
		pCopies,
	)
}
