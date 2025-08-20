package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

// ShaderModuleCreateInfo contains shader module creation information
type ShaderModuleCreateInfo struct {
	CodeSize uint32
	Code     []uint32
}

// PipelineShaderStageCreateInfo contains pipeline shader stage creation information
type PipelineShaderStageCreateInfo struct {
	Stage  ShaderStageFlags
	Module ShaderModule
	Name   string
}

// ShaderStageFlags represents shader stage flags
type ShaderStageFlags uint32

const (
	ShaderStageVertexBit                 ShaderStageFlags = C.VK_SHADER_STAGE_VERTEX_BIT
	ShaderStageTessellationControlBit    ShaderStageFlags = C.VK_SHADER_STAGE_TESSELLATION_CONTROL_BIT
	ShaderStageTessellationEvaluationBit ShaderStageFlags = C.VK_SHADER_STAGE_TESSELLATION_EVALUATION_BIT
	ShaderStageGeometryBit               ShaderStageFlags = C.VK_SHADER_STAGE_GEOMETRY_BIT
	ShaderStageFragmentBit               ShaderStageFlags = C.VK_SHADER_STAGE_FRAGMENT_BIT
	ShaderStageComputeBit                ShaderStageFlags = C.VK_SHADER_STAGE_COMPUTE_BIT
	ShaderStageAllGraphics               ShaderStageFlags = C.VK_SHADER_STAGE_ALL_GRAPHICS
	ShaderStageAll                       ShaderStageFlags = C.VK_SHADER_STAGE_ALL
)

// PipelineLayoutCreateInfo contains pipeline layout creation information
type PipelineLayoutCreateInfo struct {
	SetLayouts    []DescriptorSetLayout
	PushConstants []PushConstantRange
}

// PushConstantRange represents a push constant range
type PushConstantRange struct {
	StageFlags ShaderStageFlags
	Offset     uint32
	Size       uint32
}

// RenderPassCreateInfo contains render pass creation information
type RenderPassCreateInfo struct {
	Attachments  []AttachmentDescription
	Subpasses    []SubpassDescription
	Dependencies []SubpassDependency
}

// AttachmentDescription describes a render pass attachment
type AttachmentDescription struct {
	Format         Format
	Samples        SampleCountFlags
	LoadOp         AttachmentLoadOp
	StoreOp        AttachmentStoreOp
	StencilLoadOp  AttachmentLoadOp
	StencilStoreOp AttachmentStoreOp
	InitialLayout  ImageLayout
	FinalLayout    ImageLayout
}

// AttachmentLoadOp represents attachment load operations
type AttachmentLoadOp int32

const (
	AttachmentLoadOpLoad     AttachmentLoadOp = C.VK_ATTACHMENT_LOAD_OP_LOAD
	AttachmentLoadOpClear    AttachmentLoadOp = C.VK_ATTACHMENT_LOAD_OP_CLEAR
	AttachmentLoadOpDontCare AttachmentLoadOp = C.VK_ATTACHMENT_LOAD_OP_DONT_CARE
)

// AttachmentStoreOp represents attachment store operations
type AttachmentStoreOp int32

const (
	AttachmentStoreOpStore    AttachmentStoreOp = C.VK_ATTACHMENT_STORE_OP_STORE
	AttachmentStoreOpDontCare AttachmentStoreOp = C.VK_ATTACHMENT_STORE_OP_DONT_CARE
)

// SubpassDescription describes a subpass
type SubpassDescription struct {
	PipelineBindPoint    PipelineBindPoint
	InputAttachments     []AttachmentReference
	ColorAttachments     []AttachmentReference
	ResolveAttachments   []AttachmentReference
	DepthStencilAttachment *AttachmentReference
	PreserveAttachments  []uint32
}

// PipelineBindPoint represents pipeline bind points
type PipelineBindPoint int32

const (
	PipelineBindPointGraphics PipelineBindPoint = C.VK_PIPELINE_BIND_POINT_GRAPHICS
	PipelineBindPointCompute  PipelineBindPoint = C.VK_PIPELINE_BIND_POINT_COMPUTE
)

// AttachmentReference references an attachment
type AttachmentReference struct {
	Attachment uint32
	Layout     ImageLayout
}

// SubpassDependency describes subpass dependencies
type SubpassDependency struct {
	SrcSubpass    uint32
	DstSubpass    uint32
	SrcStageMask  PipelineStageFlags
	DstStageMask  PipelineStageFlags
	SrcAccessMask AccessFlags
	DstAccessMask AccessFlags
}

// AccessFlags represents memory access flags
type AccessFlags uint32

const (
	AccessIndirectCommandReadBit         AccessFlags = C.VK_ACCESS_INDIRECT_COMMAND_READ_BIT
	AccessIndexReadBit                   AccessFlags = C.VK_ACCESS_INDEX_READ_BIT
	AccessVertexAttributeReadBit         AccessFlags = C.VK_ACCESS_VERTEX_ATTRIBUTE_READ_BIT
	AccessUniformReadBit                 AccessFlags = C.VK_ACCESS_UNIFORM_READ_BIT
	AccessInputAttachmentReadBit         AccessFlags = C.VK_ACCESS_INPUT_ATTACHMENT_READ_BIT
	AccessShaderReadBit                  AccessFlags = C.VK_ACCESS_SHADER_READ_BIT
	AccessShaderWriteBit                 AccessFlags = C.VK_ACCESS_SHADER_WRITE_BIT
	AccessColorAttachmentReadBit         AccessFlags = C.VK_ACCESS_COLOR_ATTACHMENT_READ_BIT
	AccessColorAttachmentWriteBit        AccessFlags = C.VK_ACCESS_COLOR_ATTACHMENT_WRITE_BIT
	AccessDepthStencilAttachmentReadBit  AccessFlags = C.VK_ACCESS_DEPTH_STENCIL_ATTACHMENT_READ_BIT
	AccessDepthStencilAttachmentWriteBit AccessFlags = C.VK_ACCESS_DEPTH_STENCIL_ATTACHMENT_WRITE_BIT
	AccessTransferReadBit                AccessFlags = C.VK_ACCESS_TRANSFER_READ_BIT
	AccessTransferWriteBit               AccessFlags = C.VK_ACCESS_TRANSFER_WRITE_BIT
	AccessHostReadBit                    AccessFlags = C.VK_ACCESS_HOST_READ_BIT
	AccessHostWriteBit                   AccessFlags = C.VK_ACCESS_HOST_WRITE_BIT
	AccessMemoryReadBit                  AccessFlags = C.VK_ACCESS_MEMORY_READ_BIT
	AccessMemoryWriteBit                 AccessFlags = C.VK_ACCESS_MEMORY_WRITE_BIT
)

// CreateShaderModule creates a shader module
func CreateShaderModule(device Device, createInfo *ShaderModuleCreateInfo) (ShaderModule, error) {
	var cCreateInfo C.VkShaderModuleCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.codeSize = C.size_t(createInfo.CodeSize)
	if len(createInfo.Code) > 0 {
		cCreateInfo.pCode = (*C.uint32_t)(&createInfo.Code[0])
	}

	var shaderModule C.VkShaderModule
	result := Result(C.vkCreateShaderModule(C.VkDevice(device), &cCreateInfo, nil, &shaderModule))
	if result != Success {
		return nil, result
	}

	return ShaderModule(shaderModule), nil
}

// DestroyShaderModule destroys a shader module
func DestroyShaderModule(device Device, shaderModule ShaderModule) {
	C.vkDestroyShaderModule(C.VkDevice(device), C.VkShaderModule(shaderModule), nil)
}

// CreatePipelineLayout creates a pipeline layout
func CreatePipelineLayout(device Device, createInfo *PipelineLayoutCreateInfo) (PipelineLayout, error) {
	var cCreateInfo C.VkPipelineLayoutCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_PIPELINE_LAYOUT_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0

	// Set layouts
	var cSetLayouts []C.VkDescriptorSetLayout
	if len(createInfo.SetLayouts) > 0 {
		cSetLayouts = make([]C.VkDescriptorSetLayout, len(createInfo.SetLayouts))
		for i, layout := range createInfo.SetLayouts {
			cSetLayouts[i] = C.VkDescriptorSetLayout(layout)
		}
		cCreateInfo.setLayoutCount = C.uint32_t(len(cSetLayouts))
		cCreateInfo.pSetLayouts = &cSetLayouts[0]
	}

	// Push constant ranges
	var cPushConstants []C.VkPushConstantRange
	if len(createInfo.PushConstants) > 0 {
		cPushConstants = make([]C.VkPushConstantRange, len(createInfo.PushConstants))
		for i, pc := range createInfo.PushConstants {
			cPushConstants[i].stageFlags = C.VkShaderStageFlags(pc.StageFlags)
			cPushConstants[i].offset = C.uint32_t(pc.Offset)
			cPushConstants[i].size = C.uint32_t(pc.Size)
		}
		cCreateInfo.pushConstantRangeCount = C.uint32_t(len(cPushConstants))
		cCreateInfo.pPushConstantRanges = &cPushConstants[0]
	}

	var pipelineLayout C.VkPipelineLayout
	result := Result(C.vkCreatePipelineLayout(C.VkDevice(device), &cCreateInfo, nil, &pipelineLayout))
	if result != Success {
		return nil, result
	}

	return PipelineLayout(pipelineLayout), nil
}

// DestroyPipelineLayout destroys a pipeline layout
func DestroyPipelineLayout(device Device, pipelineLayout PipelineLayout) {
	C.vkDestroyPipelineLayout(C.VkDevice(device), C.VkPipelineLayout(pipelineLayout), nil)
}

// CreateRenderPass creates a render pass
func CreateRenderPass(device Device, createInfo *RenderPassCreateInfo) (RenderPass, error) {
	var cCreateInfo C.VkRenderPassCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_RENDER_PASS_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0

	// Attachments
	var cAttachments []C.VkAttachmentDescription
	if len(createInfo.Attachments) > 0 {
		cAttachments = make([]C.VkAttachmentDescription, len(createInfo.Attachments))
		for i, att := range createInfo.Attachments {
			cAttachments[i].flags = 0
			cAttachments[i].format = C.VkFormat(att.Format)
			cAttachments[i].samples = C.VkSampleCountFlagBits(att.Samples)
			cAttachments[i].loadOp = C.VkAttachmentLoadOp(att.LoadOp)
			cAttachments[i].storeOp = C.VkAttachmentStoreOp(att.StoreOp)
			cAttachments[i].stencilLoadOp = C.VkAttachmentLoadOp(att.StencilLoadOp)
			cAttachments[i].stencilStoreOp = C.VkAttachmentStoreOp(att.StencilStoreOp)
			cAttachments[i].initialLayout = C.VkImageLayout(att.InitialLayout)
			cAttachments[i].finalLayout = C.VkImageLayout(att.FinalLayout)
		}
		cCreateInfo.attachmentCount = C.uint32_t(len(cAttachments))
		cCreateInfo.pAttachments = &cAttachments[0]
	}

	// Note: Subpass implementation simplified for this basic version
	// Full implementation would handle all attachment references properly

	var renderPass C.VkRenderPass
	result := Result(C.vkCreateRenderPass(C.VkDevice(device), &cCreateInfo, nil, &renderPass))
	if result != Success {
		return nil, result
	}

	return RenderPass(renderPass), nil
}

// DestroyRenderPass destroys a render pass
func DestroyRenderPass(device Device, renderPass RenderPass) {
	C.vkDestroyRenderPass(C.VkDevice(device), C.VkRenderPass(renderPass), nil)
}

// Additional utility functions for common operations

// GetAPIVersion returns the supported Vulkan API version
func GetAPIVersion() Version {
	return Version13 // This system supports up to Vulkan 1.3
}

// IsExtensionSupported checks if an extension is supported
func IsExtensionSupported(extensionName string, availableExtensions []ExtensionProperties) bool {
	for _, ext := range availableExtensions {
		if ext.ExtensionName == extensionName {
			return true
		}
	}
	return false
}

// IsLayerSupported checks if a layer is supported
func IsLayerSupported(layerName string, availableLayers []LayerProperties) bool {
	for _, layer := range availableLayers {
		if layer.LayerName == layerName {
			return true
		}
	}
	return false
}