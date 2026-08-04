package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

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
	ShaderStageMeshBitEXT                ShaderStageFlags = C.VK_SHADER_STAGE_MESH_BIT_EXT
	ShaderStageTaskBitEXT                ShaderStageFlags = C.VK_SHADER_STAGE_TASK_BIT_EXT
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
	PipelineBindPoint      PipelineBindPoint
	InputAttachments       []AttachmentReference
	ColorAttachments       []AttachmentReference
	ResolveAttachments     []AttachmentReference
	DepthStencilAttachment *AttachmentReference
	PreserveAttachments    []uint32
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
	SrcSubpass      uint32
	DstSubpass      uint32
	SrcStageMask    PipelineStageFlags
	DstStageMask    PipelineStageFlags
	SrcAccessMask   AccessFlags
	DstAccessMask   AccessFlags
	DependencyFlags DependencyFlags
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
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	// The code array must be pinned because its address is stored inside
	// cCreateInfo, which is Go memory passed to C (cgo pointer rules).
	var pinner runtime.Pinner
	defer pinner.Unpin()

	var cCreateInfo C.VkShaderModuleCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_SHADER_MODULE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.codeSize = C.size_t(createInfo.CodeSize)
	if len(createInfo.Code) > 0 {
		pinner.Pin(&createInfo.Code[0])
		cCreateInfo.pCode = (*C.uint32_t)(unsafe.Pointer(&createInfo.Code[0]))
	}

	var shaderModule C.VkShaderModule
	result := Result(C.vkCreateShaderModule(C.VkDevice(device), &cCreateInfo, nil, &shaderModule))
	if result != Success {
		return nil, NewVulkanError(result, "CreateShaderModule", "Vulkan shader module creation failed")
	}

	trackResource("ShaderModule", unsafe.Pointer(shaderModule))
	return ShaderModule(shaderModule), nil
}

// DestroyShaderModule destroys a shader module
func DestroyShaderModule(device Device, shaderModule ShaderModule) {
	if device == nil || shaderModule == nil {
		return
	}
	C.vkDestroyShaderModule(C.VkDevice(device), C.VkShaderModule(shaderModule), nil)
	untrackResource("ShaderModule", unsafe.Pointer(shaderModule))
}

// CreatePipelineLayout creates a pipeline layout
func CreatePipelineLayout(device Device, createInfo *PipelineLayoutCreateInfo) (PipelineLayout, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	// Nested arrays must be pinned because their addresses are stored inside
	// cCreateInfo, which is Go memory passed to C (cgo pointer rules).
	var pinner runtime.Pinner
	defer pinner.Unpin()

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
		pinner.Pin(&cSetLayouts[0])
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
		pinner.Pin(&cPushConstants[0])
		cCreateInfo.pushConstantRangeCount = C.uint32_t(len(cPushConstants))
		cCreateInfo.pPushConstantRanges = &cPushConstants[0]
	}

	var pipelineLayout C.VkPipelineLayout
	result := Result(C.vkCreatePipelineLayout(C.VkDevice(device), &cCreateInfo, nil, &pipelineLayout))
	if result != Success {
		return nil, NewVulkanError(result, "CreatePipelineLayout", "Vulkan pipeline layout creation failed")
	}

	trackResource("PipelineLayout", unsafe.Pointer(pipelineLayout))
	return PipelineLayout(pipelineLayout), nil
}

// DestroyPipelineLayout destroys a pipeline layout
func DestroyPipelineLayout(device Device, pipelineLayout PipelineLayout) {
	if device == nil || pipelineLayout == nil {
		return
	}
	C.vkDestroyPipelineLayout(C.VkDevice(device), C.VkPipelineLayout(pipelineLayout), nil)
	untrackResource("PipelineLayout", unsafe.Pointer(pipelineLayout))
}

// CreateRenderPass creates a render pass
func CreateRenderPass(device Device, createInfo *RenderPassCreateInfo) (RenderPass, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	// All nested arrays must be pinned because their addresses are stored in
	// cCreateInfo/cSubpasses, which are Go memory passed to C (cgo rules).
	var pinner runtime.Pinner
	defer pinner.Unpin()

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
		pinner.Pin(&cAttachments[0])
		cCreateInfo.attachmentCount = C.uint32_t(len(cAttachments))
		cCreateInfo.pAttachments = &cAttachments[0]
	}

	// Subpasses - store all attachment reference arrays for the lifetime of the call
	var cSubpasses []C.VkSubpassDescription
	var allInputAttachments [][]C.VkAttachmentReference
	var allColorAttachments [][]C.VkAttachmentReference
	var allResolveAttachments [][]C.VkAttachmentReference
	var allDepthStencilAttachments []C.VkAttachmentReference
	var allPreserveAttachments [][]C.uint32_t

	if len(createInfo.Subpasses) > 0 {
		cSubpasses = make([]C.VkSubpassDescription, len(createInfo.Subpasses))
		allInputAttachments = make([][]C.VkAttachmentReference, len(createInfo.Subpasses))
		allColorAttachments = make([][]C.VkAttachmentReference, len(createInfo.Subpasses))
		allResolveAttachments = make([][]C.VkAttachmentReference, len(createInfo.Subpasses))
		allDepthStencilAttachments = make([]C.VkAttachmentReference, len(createInfo.Subpasses))
		allPreserveAttachments = make([][]C.uint32_t, len(createInfo.Subpasses))

		for i, subpass := range createInfo.Subpasses {
			cSubpasses[i].flags = 0
			cSubpasses[i].pipelineBindPoint = C.VkPipelineBindPoint(subpass.PipelineBindPoint)

			// Input attachments
			if len(subpass.InputAttachments) > 0 {
				allInputAttachments[i] = make([]C.VkAttachmentReference, len(subpass.InputAttachments))
				for j, ref := range subpass.InputAttachments {
					allInputAttachments[i][j].attachment = C.uint32_t(ref.Attachment)
					allInputAttachments[i][j].layout = C.VkImageLayout(ref.Layout)
				}
				pinner.Pin(&allInputAttachments[i][0])
				cSubpasses[i].inputAttachmentCount = C.uint32_t(len(allInputAttachments[i]))
				cSubpasses[i].pInputAttachments = &allInputAttachments[i][0]
			}

			// Color attachments
			if len(subpass.ColorAttachments) > 0 {
				allColorAttachments[i] = make([]C.VkAttachmentReference, len(subpass.ColorAttachments))
				for j, ref := range subpass.ColorAttachments {
					allColorAttachments[i][j].attachment = C.uint32_t(ref.Attachment)
					allColorAttachments[i][j].layout = C.VkImageLayout(ref.Layout)
				}
				pinner.Pin(&allColorAttachments[i][0])
				cSubpasses[i].colorAttachmentCount = C.uint32_t(len(allColorAttachments[i]))
				cSubpasses[i].pColorAttachments = &allColorAttachments[i][0]
			}

			// Resolve attachments (must match color attachment count if provided)
			if len(subpass.ResolveAttachments) > 0 {
				// Validate that resolve attachments count matches color attachments count
				if len(subpass.ResolveAttachments) != len(subpass.ColorAttachments) {
					return nil, NewValidationError("ResolveAttachments", "count must match ColorAttachments count")
				}
				allResolveAttachments[i] = make([]C.VkAttachmentReference, len(subpass.ResolveAttachments))
				for j, ref := range subpass.ResolveAttachments {
					allResolveAttachments[i][j].attachment = C.uint32_t(ref.Attachment)
					allResolveAttachments[i][j].layout = C.VkImageLayout(ref.Layout)
				}
				pinner.Pin(&allResolveAttachments[i][0])
				cSubpasses[i].pResolveAttachments = &allResolveAttachments[i][0]
			}

			// Depth/stencil attachment
			if subpass.DepthStencilAttachment != nil {
				allDepthStencilAttachments[i].attachment = C.uint32_t(subpass.DepthStencilAttachment.Attachment)
				allDepthStencilAttachments[i].layout = C.VkImageLayout(subpass.DepthStencilAttachment.Layout)
				pinner.Pin(&allDepthStencilAttachments[i])
				cSubpasses[i].pDepthStencilAttachment = &allDepthStencilAttachments[i]
			}

			// Preserve attachments
			if len(subpass.PreserveAttachments) > 0 {
				allPreserveAttachments[i] = make([]C.uint32_t, len(subpass.PreserveAttachments))
				for j, att := range subpass.PreserveAttachments {
					allPreserveAttachments[i][j] = C.uint32_t(att)
				}
				pinner.Pin(&allPreserveAttachments[i][0])
				cSubpasses[i].preserveAttachmentCount = C.uint32_t(len(allPreserveAttachments[i]))
				cSubpasses[i].pPreserveAttachments = &allPreserveAttachments[i][0]
			}
		}
		pinner.Pin(&cSubpasses[0])
		cCreateInfo.subpassCount = C.uint32_t(len(cSubpasses))
		cCreateInfo.pSubpasses = &cSubpasses[0]
	}

	// Subpass dependencies
	var cDependencies []C.VkSubpassDependency
	if len(createInfo.Dependencies) > 0 {
		cDependencies = make([]C.VkSubpassDependency, len(createInfo.Dependencies))
		for i, dep := range createInfo.Dependencies {
			cDependencies[i].srcSubpass = C.uint32_t(dep.SrcSubpass)
			cDependencies[i].dstSubpass = C.uint32_t(dep.DstSubpass)
			cDependencies[i].srcStageMask = C.VkPipelineStageFlags(dep.SrcStageMask)
			cDependencies[i].dstStageMask = C.VkPipelineStageFlags(dep.DstStageMask)
			cDependencies[i].srcAccessMask = C.VkAccessFlags(dep.SrcAccessMask)
			cDependencies[i].dstAccessMask = C.VkAccessFlags(dep.DstAccessMask)
			cDependencies[i].dependencyFlags = C.VkDependencyFlags(dep.DependencyFlags)
		}
		pinner.Pin(&cDependencies[0])
		cCreateInfo.dependencyCount = C.uint32_t(len(cDependencies))
		cCreateInfo.pDependencies = &cDependencies[0]
	}

	var renderPass C.VkRenderPass
	result := Result(C.vkCreateRenderPass(C.VkDevice(device), &cCreateInfo, nil, &renderPass))
	if result != Success {
		return nil, NewVulkanError(result, "CreateRenderPass", "Vulkan render pass creation failed")
	}

	trackResource("RenderPass", unsafe.Pointer(renderPass))
	return RenderPass(renderPass), nil
}

// DestroyRenderPass destroys a render pass
func DestroyRenderPass(device Device, renderPass RenderPass) {
	if device == nil || renderPass == nil {
		return
	}
	C.vkDestroyRenderPass(C.VkDevice(device), C.VkRenderPass(renderPass), nil)
	untrackResource("RenderPass", unsafe.Pointer(renderPass))
}

// GetRenderAreaGranularity returns the render area granularity for a render pass
func GetRenderAreaGranularity(device Device, renderPass RenderPass) Extent2D {
	if device == nil || renderPass == nil {
		return Extent2D{}
	}
	var cGranularity C.VkExtent2D
	C.vkGetRenderAreaGranularity(C.VkDevice(device), C.VkRenderPass(renderPass), &cGranularity)

	return Extent2D{
		Width:  uint32(cGranularity.width),
		Height: uint32(cGranularity.height),
	}
}

// ComputePipelineCreateInfo contains compute pipeline creation information
type ComputePipelineCreateInfo struct {
	Stage  PipelineShaderStageCreateInfo
	Layout PipelineLayout
}

// CreateComputePipelines creates compute pipelines
func CreateComputePipelines(device Device, pipelineCache PipelineCache, createInfos []ComputePipelineCreateInfo) ([]Pipeline, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if len(createInfos) == 0 {
		return nil, nil
	}

	cCreateInfos := make([]C.VkComputePipelineCreateInfo, len(createInfos))
	cPipelines := make([]C.VkPipeline, len(createInfos))

	// Collect C strings for proper memory management
	cNames := make([]*C.char, len(createInfos))

	for i, info := range createInfos {
		cCreateInfos[i].sType = C.VK_STRUCTURE_TYPE_COMPUTE_PIPELINE_CREATE_INFO
		cCreateInfos[i].pNext = nil
		cCreateInfos[i].flags = 0

		// Set up shader stage
		cCreateInfos[i].stage.sType = C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO
		cCreateInfos[i].stage.pNext = nil
		cCreateInfos[i].stage.flags = 0
		cCreateInfos[i].stage.stage = C.VkShaderStageFlagBits(info.Stage.Stage)
		cCreateInfos[i].stage.module = C.VkShaderModule(info.Stage.Module)

		// Convert name to C string and store for later cleanup
		cNames[i] = C.CString(info.Stage.Name)
		cCreateInfos[i].stage.pName = cNames[i]
		cCreateInfos[i].stage.pSpecializationInfo = nil

		cCreateInfos[i].layout = C.VkPipelineLayout(info.Layout)
		cCreateInfos[i].basePipelineHandle = C.VkPipeline(nil)
		cCreateInfos[i].basePipelineIndex = -1
	}

	// Free all C strings after API call regardless of success/failure
	defer func() {
		for _, cName := range cNames {
			if cName != nil {
				C.free(unsafe.Pointer(cName))
			}
		}
	}()

	result := Result(C.vkCreateComputePipelines(
		C.VkDevice(device),
		C.VkPipelineCache(pipelineCache),
		C.uint32_t(len(cCreateInfos)),
		&cCreateInfos[0],
		nil,
		&cPipelines[0],
	))

	if result != Success {
		// A failed batch create may still have created some pipelines
		// (failed entries are VK_NULL_HANDLE); destroy them to avoid leaks.
		for _, pipeline := range cPipelines {
			if pipeline != nil {
				C.vkDestroyPipeline(C.VkDevice(device), pipeline, nil)
			}
		}
		return nil, NewVulkanError(result, "CreateComputePipelines", "Vulkan compute pipeline creation failed")
	}

	pipelines := make([]Pipeline, len(cPipelines))
	for i, pipeline := range cPipelines {
		pipelines[i] = Pipeline(pipeline)
		trackResource("Pipeline", unsafe.Pointer(pipeline))
	}

	return pipelines, nil
}

// DestroyPipeline destroys a pipeline
func DestroyPipeline(device Device, pipeline Pipeline) {
	if device == nil || pipeline == nil {
		return
	}
	C.vkDestroyPipeline(C.VkDevice(device), C.VkPipeline(pipeline), nil)
	untrackResource("Pipeline", unsafe.Pointer(pipeline))
}

// ============================================================================
// Graphics Pipeline Types and Functions
// ============================================================================

// VertexInputRate represents the rate at which vertex attributes are pulled from buffers
type VertexInputRate uint32

const (
	VertexInputRateVertex   VertexInputRate = C.VK_VERTEX_INPUT_RATE_VERTEX
	VertexInputRateInstance VertexInputRate = C.VK_VERTEX_INPUT_RATE_INSTANCE
)

// VertexInputBindingDescription describes a vertex input binding
type VertexInputBindingDescription struct {
	Binding   uint32
	Stride    uint32
	InputRate VertexInputRate
}

// VertexInputAttributeDescription describes a vertex input attribute
type VertexInputAttributeDescription struct {
	Location uint32
	Binding  uint32
	Format   Format
	Offset   uint32
}

// PipelineVertexInputStateCreateInfo contains vertex input state creation information
type PipelineVertexInputStateCreateInfo struct {
	VertexBindingDescriptions   []VertexInputBindingDescription
	VertexAttributeDescriptions []VertexInputAttributeDescription
}

// PipelineInputAssemblyStateCreateInfo contains input assembly state creation information
type PipelineInputAssemblyStateCreateInfo struct {
	Topology               PrimitiveTopology
	PrimitiveRestartEnable bool
}

// PipelineTessellationStateCreateInfo contains tessellation state creation information
type PipelineTessellationStateCreateInfo struct {
	PatchControlPoints uint32
}

// PipelineViewportStateCreateInfo contains viewport state creation information
type PipelineViewportStateCreateInfo struct {
	Viewports []Viewport
	Scissors  []Rect2D
}

// PolygonMode represents polygon rasterization mode
type PolygonMode uint32

const (
	PolygonModeFill  PolygonMode = C.VK_POLYGON_MODE_FILL
	PolygonModeLine  PolygonMode = C.VK_POLYGON_MODE_LINE
	PolygonModePoint PolygonMode = C.VK_POLYGON_MODE_POINT
)

// PipelineRasterizationStateCreateInfo contains rasterization state creation information
type PipelineRasterizationStateCreateInfo struct {
	DepthClampEnable        bool
	RasterizerDiscardEnable bool
	PolygonMode             PolygonMode
	CullMode                CullModeFlags
	FrontFace               FrontFace
	DepthBiasEnable         bool
	DepthBiasConstantFactor float32
	DepthBiasClamp          float32
	DepthBiasSlopeFactor    float32
	LineWidth               float32
}

// PipelineMultisampleStateCreateInfo contains multisample state creation information
type PipelineMultisampleStateCreateInfo struct {
	RasterizationSamples  SampleCountFlags
	SampleShadingEnable   bool
	MinSampleShading      float32
	SampleMask            []uint32
	AlphaToCoverageEnable bool
	AlphaToOneEnable      bool
}

// StencilOpState contains stencil operation state
type StencilOpState struct {
	FailOp      StencilOp
	PassOp      StencilOp
	DepthFailOp StencilOp
	CompareOp   CompareOp
	CompareMask uint32
	WriteMask   uint32
	Reference   uint32
}

// PipelineDepthStencilStateCreateInfo contains depth/stencil state creation information
type PipelineDepthStencilStateCreateInfo struct {
	DepthTestEnable       bool
	DepthWriteEnable      bool
	DepthCompareOp        CompareOp
	DepthBoundsTestEnable bool
	StencilTestEnable     bool
	Front                 StencilOpState
	Back                  StencilOpState
	MinDepthBounds        float32
	MaxDepthBounds        float32
}

// BlendFactor represents blend factors
type BlendFactor uint32

const (
	BlendFactorZero                  BlendFactor = C.VK_BLEND_FACTOR_ZERO
	BlendFactorOne                   BlendFactor = C.VK_BLEND_FACTOR_ONE
	BlendFactorSrcColor              BlendFactor = C.VK_BLEND_FACTOR_SRC_COLOR
	BlendFactorOneMinusSrcColor      BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_SRC_COLOR
	BlendFactorDstColor              BlendFactor = C.VK_BLEND_FACTOR_DST_COLOR
	BlendFactorOneMinusDstColor      BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_DST_COLOR
	BlendFactorSrcAlpha              BlendFactor = C.VK_BLEND_FACTOR_SRC_ALPHA
	BlendFactorOneMinusSrcAlpha      BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_SRC_ALPHA
	BlendFactorDstAlpha              BlendFactor = C.VK_BLEND_FACTOR_DST_ALPHA
	BlendFactorOneMinusDstAlpha      BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_DST_ALPHA
	BlendFactorConstantColor         BlendFactor = C.VK_BLEND_FACTOR_CONSTANT_COLOR
	BlendFactorOneMinusConstantColor BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_CONSTANT_COLOR
	BlendFactorConstantAlpha         BlendFactor = C.VK_BLEND_FACTOR_CONSTANT_ALPHA
	BlendFactorOneMinusConstantAlpha BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_CONSTANT_ALPHA
	BlendFactorSrcAlphaSaturate      BlendFactor = C.VK_BLEND_FACTOR_SRC_ALPHA_SATURATE
	BlendFactorSrc1Color             BlendFactor = C.VK_BLEND_FACTOR_SRC1_COLOR
	BlendFactorOneMinusSrc1Color     BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_SRC1_COLOR
	BlendFactorSrc1Alpha             BlendFactor = C.VK_BLEND_FACTOR_SRC1_ALPHA
	BlendFactorOneMinusSrc1Alpha     BlendFactor = C.VK_BLEND_FACTOR_ONE_MINUS_SRC1_ALPHA
)

// BlendOp represents blend operations
type BlendOp uint32

const (
	BlendOpAdd             BlendOp = C.VK_BLEND_OP_ADD
	BlendOpSubtract        BlendOp = C.VK_BLEND_OP_SUBTRACT
	BlendOpReverseSubtract BlendOp = C.VK_BLEND_OP_REVERSE_SUBTRACT
	BlendOpMin             BlendOp = C.VK_BLEND_OP_MIN
	BlendOpMax             BlendOp = C.VK_BLEND_OP_MAX
)

// ColorComponentFlags represents color component write mask
type ColorComponentFlags uint32

const (
	ColorComponentRBit ColorComponentFlags = C.VK_COLOR_COMPONENT_R_BIT
	ColorComponentGBit ColorComponentFlags = C.VK_COLOR_COMPONENT_G_BIT
	ColorComponentBBit ColorComponentFlags = C.VK_COLOR_COMPONENT_B_BIT
	ColorComponentABit ColorComponentFlags = C.VK_COLOR_COMPONENT_A_BIT
	ColorComponentAll  ColorComponentFlags = ColorComponentRBit | ColorComponentGBit | ColorComponentBBit | ColorComponentABit
)

// LogicOp represents logical operations
type LogicOp uint32

const (
	LogicOpClear        LogicOp = C.VK_LOGIC_OP_CLEAR
	LogicOpAnd          LogicOp = C.VK_LOGIC_OP_AND
	LogicOpAndReverse   LogicOp = C.VK_LOGIC_OP_AND_REVERSE
	LogicOpCopy         LogicOp = C.VK_LOGIC_OP_COPY
	LogicOpAndInverted  LogicOp = C.VK_LOGIC_OP_AND_INVERTED
	LogicOpNoOp         LogicOp = C.VK_LOGIC_OP_NO_OP
	LogicOpXor          LogicOp = C.VK_LOGIC_OP_XOR
	LogicOpOr           LogicOp = C.VK_LOGIC_OP_OR
	LogicOpNor          LogicOp = C.VK_LOGIC_OP_NOR
	LogicOpEquivalent   LogicOp = C.VK_LOGIC_OP_EQUIVALENT
	LogicOpInvert       LogicOp = C.VK_LOGIC_OP_INVERT
	LogicOpOrReverse    LogicOp = C.VK_LOGIC_OP_OR_REVERSE
	LogicOpCopyInverted LogicOp = C.VK_LOGIC_OP_COPY_INVERTED
	LogicOpOrInverted   LogicOp = C.VK_LOGIC_OP_OR_INVERTED
	LogicOpNand         LogicOp = C.VK_LOGIC_OP_NAND
	LogicOpSet          LogicOp = C.VK_LOGIC_OP_SET
)

// PipelineColorBlendAttachmentState contains color blend attachment state
type PipelineColorBlendAttachmentState struct {
	BlendEnable         bool
	SrcColorBlendFactor BlendFactor
	DstColorBlendFactor BlendFactor
	ColorBlendOp        BlendOp
	SrcAlphaBlendFactor BlendFactor
	DstAlphaBlendFactor BlendFactor
	AlphaBlendOp        BlendOp
	ColorWriteMask      ColorComponentFlags
}

// PipelineColorBlendStateCreateInfo contains color blend state creation information
type PipelineColorBlendStateCreateInfo struct {
	LogicOpEnable  bool
	LogicOp        LogicOp
	Attachments    []PipelineColorBlendAttachmentState
	BlendConstants [4]float32
}

// DynamicState represents dynamic pipeline states
type DynamicState uint32

const (
	DynamicStateViewport                 DynamicState = C.VK_DYNAMIC_STATE_VIEWPORT
	DynamicStateScissor                  DynamicState = C.VK_DYNAMIC_STATE_SCISSOR
	DynamicStateLineWidth                DynamicState = C.VK_DYNAMIC_STATE_LINE_WIDTH
	DynamicStateDepthBias                DynamicState = C.VK_DYNAMIC_STATE_DEPTH_BIAS
	DynamicStateBlendConstants           DynamicState = C.VK_DYNAMIC_STATE_BLEND_CONSTANTS
	DynamicStateDepthBounds              DynamicState = C.VK_DYNAMIC_STATE_DEPTH_BOUNDS
	DynamicStateStencilCompareMask       DynamicState = C.VK_DYNAMIC_STATE_STENCIL_COMPARE_MASK
	DynamicStateStencilWriteMask         DynamicState = C.VK_DYNAMIC_STATE_STENCIL_WRITE_MASK
	DynamicStateStencilReference         DynamicState = C.VK_DYNAMIC_STATE_STENCIL_REFERENCE
	DynamicStateCullMode                 DynamicState = C.VK_DYNAMIC_STATE_CULL_MODE
	DynamicStateFrontFace                DynamicState = C.VK_DYNAMIC_STATE_FRONT_FACE
	DynamicStatePrimitiveTopology        DynamicState = C.VK_DYNAMIC_STATE_PRIMITIVE_TOPOLOGY
	DynamicStateViewportWithCount        DynamicState = C.VK_DYNAMIC_STATE_VIEWPORT_WITH_COUNT
	DynamicStateScissorWithCount         DynamicState = C.VK_DYNAMIC_STATE_SCISSOR_WITH_COUNT
	DynamicStateVertexInputBindingStride DynamicState = C.VK_DYNAMIC_STATE_VERTEX_INPUT_BINDING_STRIDE
	DynamicStateDepthTestEnable          DynamicState = C.VK_DYNAMIC_STATE_DEPTH_TEST_ENABLE
	DynamicStateDepthWriteEnable         DynamicState = C.VK_DYNAMIC_STATE_DEPTH_WRITE_ENABLE
	DynamicStateDepthCompareOp           DynamicState = C.VK_DYNAMIC_STATE_DEPTH_COMPARE_OP
	DynamicStateDepthBoundsTestEnable    DynamicState = C.VK_DYNAMIC_STATE_DEPTH_BOUNDS_TEST_ENABLE
	DynamicStateStencilTestEnable        DynamicState = C.VK_DYNAMIC_STATE_STENCIL_TEST_ENABLE
	DynamicStateStencilOp                DynamicState = C.VK_DYNAMIC_STATE_STENCIL_OP
	DynamicStateRasterizerDiscardEnable  DynamicState = C.VK_DYNAMIC_STATE_RASTERIZER_DISCARD_ENABLE
	DynamicStateDepthBiasEnable          DynamicState = C.VK_DYNAMIC_STATE_DEPTH_BIAS_ENABLE
	DynamicStatePrimitiveRestartEnable   DynamicState = C.VK_DYNAMIC_STATE_PRIMITIVE_RESTART_ENABLE
)

// PipelineDynamicStateCreateInfo contains dynamic state creation information
type PipelineDynamicStateCreateInfo struct {
	DynamicStates []DynamicState
}

// GraphicsPipelineCreateInfo contains graphics pipeline creation information
type GraphicsPipelineCreateInfo struct {
	Stages             []PipelineShaderStageCreateInfo
	VertexInputState   *PipelineVertexInputStateCreateInfo
	InputAssemblyState *PipelineInputAssemblyStateCreateInfo
	TessellationState  *PipelineTessellationStateCreateInfo
	ViewportState      *PipelineViewportStateCreateInfo
	RasterizationState *PipelineRasterizationStateCreateInfo
	MultisampleState   *PipelineMultisampleStateCreateInfo
	DepthStencilState  *PipelineDepthStencilStateCreateInfo
	ColorBlendState    *PipelineColorBlendStateCreateInfo
	DynamicState       *PipelineDynamicStateCreateInfo
	Layout             PipelineLayout
	RenderPass         RenderPass
	Subpass            uint32
	BasePipelineHandle Pipeline
	BasePipelineIndex  int32
}

// graphicsPipelineBuilder holds state needed for building graphics pipelines.
// All nested arrays whose addresses end up inside structs passed to C are
// pinned via the pinner (cgo pointer rules).
type graphicsPipelineBuilder struct {
	pinner                 runtime.Pinner
	cNames                 []*C.char
	cStageArrays           [][]C.VkPipelineShaderStageCreateInfo
	cBindingArrays         [][]C.VkVertexInputBindingDescription
	cAttributeArrays       [][]C.VkVertexInputAttributeDescription
	cViewportArrays        [][]C.VkViewport
	cScissorArrays         [][]C.VkRect2D
	cBlendAttachmentArrays [][]C.VkPipelineColorBlendAttachmentState
	cDynamicStateArrays    [][]C.VkDynamicState
	cSampleMaskArrays      [][]C.VkSampleMask
}

// setupShaderStages configures shader stages for a pipeline
func (b *graphicsPipelineBuilder) setupShaderStages(cInfo *C.VkGraphicsPipelineCreateInfo, stages []PipelineShaderStageCreateInfo) {
	if len(stages) == 0 {
		return
	}
	cStages := make([]C.VkPipelineShaderStageCreateInfo, len(stages))
	for j, stage := range stages {
		cStages[j].sType = C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO
		cStages[j].pNext = nil
		cStages[j].flags = 0
		cStages[j].stage = C.VkShaderStageFlagBits(stage.Stage)
		cStages[j].module = C.VkShaderModule(stage.Module)

		cName := C.CString(stage.Name)
		b.cNames = append(b.cNames, cName)
		cStages[j].pName = cName
		cStages[j].pSpecializationInfo = nil
	}
	b.cStageArrays = append(b.cStageArrays, cStages)
	b.pinner.Pin(&cStages[0])
	cInfo.stageCount = C.uint32_t(len(cStages))
	cInfo.pStages = &b.cStageArrays[len(b.cStageArrays)-1][0]
}

// setupVertexInputState configures vertex input state for a pipeline
func (b *graphicsPipelineBuilder) setupVertexInputState(cState *C.VkPipelineVertexInputStateCreateInfo, info *PipelineVertexInputStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	if info == nil {
		return
	}
	if len(info.VertexBindingDescriptions) > 0 {
		cBindings := make([]C.VkVertexInputBindingDescription, len(info.VertexBindingDescriptions))
		for j, binding := range info.VertexBindingDescriptions {
			cBindings[j].binding = C.uint32_t(binding.Binding)
			cBindings[j].stride = C.uint32_t(binding.Stride)
			cBindings[j].inputRate = C.VkVertexInputRate(binding.InputRate)
		}
		b.cBindingArrays = append(b.cBindingArrays, cBindings)
		b.pinner.Pin(&cBindings[0])
		cState.vertexBindingDescriptionCount = C.uint32_t(len(cBindings))
		cState.pVertexBindingDescriptions = &b.cBindingArrays[len(b.cBindingArrays)-1][0]
	}
	if len(info.VertexAttributeDescriptions) > 0 {
		cAttributes := make([]C.VkVertexInputAttributeDescription, len(info.VertexAttributeDescriptions))
		for j, attr := range info.VertexAttributeDescriptions {
			cAttributes[j].location = C.uint32_t(attr.Location)
			cAttributes[j].binding = C.uint32_t(attr.Binding)
			cAttributes[j].format = C.VkFormat(attr.Format)
			cAttributes[j].offset = C.uint32_t(attr.Offset)
		}
		b.cAttributeArrays = append(b.cAttributeArrays, cAttributes)
		b.pinner.Pin(&cAttributes[0])
		cState.vertexAttributeDescriptionCount = C.uint32_t(len(cAttributes))
		cState.pVertexAttributeDescriptions = &b.cAttributeArrays[len(b.cAttributeArrays)-1][0]
	}
}

// setupViewportState configures viewport state for a pipeline
func (b *graphicsPipelineBuilder) setupViewportState(cState *C.VkPipelineViewportStateCreateInfo, info *PipelineViewportStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	if info == nil {
		cState.viewportCount = 1
		cState.pViewports = nil
		cState.scissorCount = 1
		cState.pScissors = nil
		return
	}
	if len(info.Viewports) > 0 {
		cViewports := make([]C.VkViewport, len(info.Viewports))
		for j, vp := range info.Viewports {
			cViewports[j].x = C.float(vp.X)
			cViewports[j].y = C.float(vp.Y)
			cViewports[j].width = C.float(vp.Width)
			cViewports[j].height = C.float(vp.Height)
			cViewports[j].minDepth = C.float(vp.MinDepth)
			cViewports[j].maxDepth = C.float(vp.MaxDepth)
		}
		b.cViewportArrays = append(b.cViewportArrays, cViewports)
		b.pinner.Pin(&cViewports[0])
		cState.viewportCount = C.uint32_t(len(cViewports))
		cState.pViewports = &b.cViewportArrays[len(b.cViewportArrays)-1][0]
	} else {
		cState.viewportCount = 1
		cState.pViewports = nil
	}
	if len(info.Scissors) > 0 {
		cScissors := make([]C.VkRect2D, len(info.Scissors))
		for j, sc := range info.Scissors {
			cScissors[j].offset.x = C.int32_t(sc.Offset.X)
			cScissors[j].offset.y = C.int32_t(sc.Offset.Y)
			cScissors[j].extent.width = C.uint32_t(sc.Extent.Width)
			cScissors[j].extent.height = C.uint32_t(sc.Extent.Height)
		}
		b.cScissorArrays = append(b.cScissorArrays, cScissors)
		b.pinner.Pin(&cScissors[0])
		cState.scissorCount = C.uint32_t(len(cScissors))
		cState.pScissors = &b.cScissorArrays[len(b.cScissorArrays)-1][0]
	} else {
		cState.scissorCount = 1
		cState.pScissors = nil
	}
}

// setupRasterizationState configures rasterization state for a pipeline
func setupRasterizationState(cState *C.VkPipelineRasterizationStateCreateInfo, info *PipelineRasterizationStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	if info != nil {
		cState.depthClampEnable = boolToVkBool32(info.DepthClampEnable)
		cState.rasterizerDiscardEnable = boolToVkBool32(info.RasterizerDiscardEnable)
		cState.polygonMode = C.VkPolygonMode(info.PolygonMode)
		cState.cullMode = C.VkCullModeFlags(info.CullMode)
		cState.frontFace = C.VkFrontFace(info.FrontFace)
		cState.depthBiasEnable = boolToVkBool32(info.DepthBiasEnable)
		cState.depthBiasConstantFactor = C.float(info.DepthBiasConstantFactor)
		cState.depthBiasClamp = C.float(info.DepthBiasClamp)
		cState.depthBiasSlopeFactor = C.float(info.DepthBiasSlopeFactor)
		cState.lineWidth = C.float(info.LineWidth)
	} else {
		cState.depthClampEnable = C.VK_FALSE
		cState.rasterizerDiscardEnable = C.VK_FALSE
		cState.polygonMode = C.VK_POLYGON_MODE_FILL
		cState.cullMode = C.VkCullModeFlags(C.VK_CULL_MODE_BACK_BIT)
		cState.frontFace = C.VK_FRONT_FACE_COUNTER_CLOCKWISE
		cState.depthBiasEnable = C.VK_FALSE
		cState.lineWidth = 1.0
	}
}

// setupMultisampleState configures multisample state for a pipeline
func (b *graphicsPipelineBuilder) setupMultisampleState(cState *C.VkPipelineMultisampleStateCreateInfo, info *PipelineMultisampleStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	if info != nil {
		cState.rasterizationSamples = C.VkSampleCountFlagBits(info.RasterizationSamples)
		cState.sampleShadingEnable = boolToVkBool32(info.SampleShadingEnable)
		cState.minSampleShading = C.float(info.MinSampleShading)
		if len(info.SampleMask) > 0 {
			cSampleMask := make([]C.VkSampleMask, len(info.SampleMask))
			for j, mask := range info.SampleMask {
				cSampleMask[j] = C.VkSampleMask(mask)
			}
			b.cSampleMaskArrays = append(b.cSampleMaskArrays, cSampleMask)
			b.pinner.Pin(&cSampleMask[0])
			cState.pSampleMask = &b.cSampleMaskArrays[len(b.cSampleMaskArrays)-1][0]
		}
		cState.alphaToCoverageEnable = boolToVkBool32(info.AlphaToCoverageEnable)
		cState.alphaToOneEnable = boolToVkBool32(info.AlphaToOneEnable)
	} else {
		cState.rasterizationSamples = C.VK_SAMPLE_COUNT_1_BIT
		cState.sampleShadingEnable = C.VK_FALSE
		cState.minSampleShading = 1.0
		cState.pSampleMask = nil
		cState.alphaToCoverageEnable = C.VK_FALSE
		cState.alphaToOneEnable = C.VK_FALSE
	}
}

// setupDepthStencilState configures depth stencil state for a pipeline
func setupDepthStencilState(cState *C.VkPipelineDepthStencilStateCreateInfo, info *PipelineDepthStencilStateCreateInfo) bool {
	if info == nil {
		return false
	}
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_DEPTH_STENCIL_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	cState.depthTestEnable = boolToVkBool32(info.DepthTestEnable)
	cState.depthWriteEnable = boolToVkBool32(info.DepthWriteEnable)
	cState.depthCompareOp = C.VkCompareOp(info.DepthCompareOp)
	cState.depthBoundsTestEnable = boolToVkBool32(info.DepthBoundsTestEnable)
	cState.stencilTestEnable = boolToVkBool32(info.StencilTestEnable)
	cState.front.failOp = C.VkStencilOp(info.Front.FailOp)
	cState.front.passOp = C.VkStencilOp(info.Front.PassOp)
	cState.front.depthFailOp = C.VkStencilOp(info.Front.DepthFailOp)
	cState.front.compareOp = C.VkCompareOp(info.Front.CompareOp)
	cState.front.compareMask = C.uint32_t(info.Front.CompareMask)
	cState.front.writeMask = C.uint32_t(info.Front.WriteMask)
	cState.front.reference = C.uint32_t(info.Front.Reference)
	cState.back.failOp = C.VkStencilOp(info.Back.FailOp)
	cState.back.passOp = C.VkStencilOp(info.Back.PassOp)
	cState.back.depthFailOp = C.VkStencilOp(info.Back.DepthFailOp)
	cState.back.compareOp = C.VkCompareOp(info.Back.CompareOp)
	cState.back.compareMask = C.uint32_t(info.Back.CompareMask)
	cState.back.writeMask = C.uint32_t(info.Back.WriteMask)
	cState.back.reference = C.uint32_t(info.Back.Reference)
	cState.minDepthBounds = C.float(info.MinDepthBounds)
	cState.maxDepthBounds = C.float(info.MaxDepthBounds)
	return true
}

// setupColorBlendState configures color blend state for a pipeline
func (b *graphicsPipelineBuilder) setupColorBlendState(cState *C.VkPipelineColorBlendStateCreateInfo, info *PipelineColorBlendStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	if info != nil {
		cState.logicOpEnable = boolToVkBool32(info.LogicOpEnable)
		cState.logicOp = C.VkLogicOp(info.LogicOp)
		if len(info.Attachments) > 0 {
			cBlendAttachments := make([]C.VkPipelineColorBlendAttachmentState, len(info.Attachments))
			for j, att := range info.Attachments {
				cBlendAttachments[j].blendEnable = boolToVkBool32(att.BlendEnable)
				cBlendAttachments[j].srcColorBlendFactor = C.VkBlendFactor(att.SrcColorBlendFactor)
				cBlendAttachments[j].dstColorBlendFactor = C.VkBlendFactor(att.DstColorBlendFactor)
				cBlendAttachments[j].colorBlendOp = C.VkBlendOp(att.ColorBlendOp)
				cBlendAttachments[j].srcAlphaBlendFactor = C.VkBlendFactor(att.SrcAlphaBlendFactor)
				cBlendAttachments[j].dstAlphaBlendFactor = C.VkBlendFactor(att.DstAlphaBlendFactor)
				cBlendAttachments[j].alphaBlendOp = C.VkBlendOp(att.AlphaBlendOp)
				cBlendAttachments[j].colorWriteMask = C.VkColorComponentFlags(att.ColorWriteMask)
			}
			b.cBlendAttachmentArrays = append(b.cBlendAttachmentArrays, cBlendAttachments)
			b.pinner.Pin(&cBlendAttachments[0])
			cState.attachmentCount = C.uint32_t(len(cBlendAttachments))
			cState.pAttachments = &b.cBlendAttachmentArrays[len(b.cBlendAttachmentArrays)-1][0]
		}
		for j := 0; j < 4; j++ {
			cState.blendConstants[j] = C.float(info.BlendConstants[j])
		}
	} else {
		cState.logicOpEnable = C.VK_FALSE
		cState.logicOp = C.VK_LOGIC_OP_COPY
		defaultAttachment := make([]C.VkPipelineColorBlendAttachmentState, 1)
		defaultAttachment[0].blendEnable = C.VK_FALSE
		defaultAttachment[0].colorWriteMask = C.VkColorComponentFlags(ColorComponentAll)
		b.cBlendAttachmentArrays = append(b.cBlendAttachmentArrays, defaultAttachment)
		b.pinner.Pin(&defaultAttachment[0])
		cState.attachmentCount = 1
		cState.pAttachments = &b.cBlendAttachmentArrays[len(b.cBlendAttachmentArrays)-1][0]
	}
}

// setupDynamicState configures dynamic state for a pipeline
func (b *graphicsPipelineBuilder) setupDynamicState(cState *C.VkPipelineDynamicStateCreateInfo, info *PipelineDynamicStateCreateInfo) bool {
	if info == nil || len(info.DynamicStates) == 0 {
		return false
	}
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	cDynStates := make([]C.VkDynamicState, len(info.DynamicStates))
	for j, ds := range info.DynamicStates {
		cDynStates[j] = C.VkDynamicState(ds)
	}
	b.cDynamicStateArrays = append(b.cDynamicStateArrays, cDynStates)
	b.pinner.Pin(&cDynStates[0])
	cState.dynamicStateCount = C.uint32_t(len(cDynStates))
	cState.pDynamicStates = &b.cDynamicStateArrays[len(b.cDynamicStateArrays)-1][0]
	return true
}

// freeNames frees all C strings allocated by the builder
func (b *graphicsPipelineBuilder) freeNames() {
	for _, cName := range b.cNames {
		if cName != nil {
			C.free(unsafe.Pointer(cName))
		}
	}
}

// CreateGraphicsPipelines creates graphics pipelines
func CreateGraphicsPipelines(device Device, pipelineCache PipelineCache, createInfos []GraphicsPipelineCreateInfo) ([]Pipeline, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if len(createInfos) == 0 {
		return nil, nil
	}

	builder := &graphicsPipelineBuilder{}
	defer builder.freeNames()
	defer builder.pinner.Unpin()

	cCreateInfos := make([]C.VkGraphicsPipelineCreateInfo, len(createInfos))
	cPipelines := make([]C.VkPipeline, len(createInfos))

	// Storage for state structs. These arrays must be pinned because their
	// element addresses are stored inside cCreateInfos, which is Go memory
	// passed to C (cgo pointer rules).
	cVertexInputStates := make([]C.VkPipelineVertexInputStateCreateInfo, len(createInfos))
	cInputAssemblyStates := make([]C.VkPipelineInputAssemblyStateCreateInfo, len(createInfos))
	cTessellationStates := make([]C.VkPipelineTessellationStateCreateInfo, len(createInfos))
	cViewportStates := make([]C.VkPipelineViewportStateCreateInfo, len(createInfos))
	cRasterizationStates := make([]C.VkPipelineRasterizationStateCreateInfo, len(createInfos))
	cMultisampleStates := make([]C.VkPipelineMultisampleStateCreateInfo, len(createInfos))
	cDepthStencilStates := make([]C.VkPipelineDepthStencilStateCreateInfo, len(createInfos))
	cColorBlendStates := make([]C.VkPipelineColorBlendStateCreateInfo, len(createInfos))
	cDynamicStates := make([]C.VkPipelineDynamicStateCreateInfo, len(createInfos))

	builder.pinner.Pin(&cVertexInputStates[0])
	builder.pinner.Pin(&cInputAssemblyStates[0])
	builder.pinner.Pin(&cTessellationStates[0])
	builder.pinner.Pin(&cViewportStates[0])
	builder.pinner.Pin(&cRasterizationStates[0])
	builder.pinner.Pin(&cMultisampleStates[0])
	builder.pinner.Pin(&cDepthStencilStates[0])
	builder.pinner.Pin(&cColorBlendStates[0])
	builder.pinner.Pin(&cDynamicStates[0])

	for i, info := range createInfos {
		cCreateInfos[i].sType = C.VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO
		cCreateInfos[i].pNext = nil
		cCreateInfos[i].flags = 0

		// Shader stages
		builder.setupShaderStages(&cCreateInfos[i], info.Stages)

		// Vertex input state
		builder.setupVertexInputState(&cVertexInputStates[i], info.VertexInputState)
		cCreateInfos[i].pVertexInputState = &cVertexInputStates[i]

		// Input assembly state
		setupInputAssemblyState(&cInputAssemblyStates[i], info.InputAssemblyState)
		cCreateInfos[i].pInputAssemblyState = &cInputAssemblyStates[i]

		// Tessellation state (optional)
		if info.TessellationState != nil {
			setupTessellationState(&cTessellationStates[i], info.TessellationState)
			cCreateInfos[i].pTessellationState = &cTessellationStates[i]
		}

		// Viewport state
		builder.setupViewportState(&cViewportStates[i], info.ViewportState)
		cCreateInfos[i].pViewportState = &cViewportStates[i]

		// Rasterization state
		setupRasterizationState(&cRasterizationStates[i], info.RasterizationState)
		cCreateInfos[i].pRasterizationState = &cRasterizationStates[i]

		// Multisample state
		builder.setupMultisampleState(&cMultisampleStates[i], info.MultisampleState)
		cCreateInfos[i].pMultisampleState = &cMultisampleStates[i]

		// Depth stencil state (optional)
		if setupDepthStencilState(&cDepthStencilStates[i], info.DepthStencilState) {
			cCreateInfos[i].pDepthStencilState = &cDepthStencilStates[i]
		}

		// Color blend state
		builder.setupColorBlendState(&cColorBlendStates[i], info.ColorBlendState)
		cCreateInfos[i].pColorBlendState = &cColorBlendStates[i]

		// Dynamic state (optional)
		if builder.setupDynamicState(&cDynamicStates[i], info.DynamicState) {
			cCreateInfos[i].pDynamicState = &cDynamicStates[i]
		}

		// Layout and render pass
		cCreateInfos[i].layout = C.VkPipelineLayout(info.Layout)
		cCreateInfos[i].renderPass = C.VkRenderPass(info.RenderPass)
		cCreateInfos[i].subpass = C.uint32_t(info.Subpass)
		cCreateInfos[i].basePipelineHandle = C.VkPipeline(info.BasePipelineHandle)
		cCreateInfos[i].basePipelineIndex = C.int32_t(info.BasePipelineIndex)
	}

	result := Result(C.vkCreateGraphicsPipelines(
		C.VkDevice(device),
		C.VkPipelineCache(pipelineCache),
		C.uint32_t(len(cCreateInfos)),
		&cCreateInfos[0],
		nil,
		&cPipelines[0],
	))

	if result != Success {
		// A failed batch create may still have created some pipelines
		// (failed entries are VK_NULL_HANDLE); destroy them to avoid leaks.
		for _, pipeline := range cPipelines {
			if pipeline != nil {
				C.vkDestroyPipeline(C.VkDevice(device), pipeline, nil)
			}
		}
		return nil, NewVulkanError(result, "CreateGraphicsPipelines", "Vulkan graphics pipeline creation failed")
	}

	pipelines := make([]Pipeline, len(cPipelines))
	for i, pipeline := range cPipelines {
		pipelines[i] = Pipeline(pipeline)
		trackResource("Pipeline", unsafe.Pointer(pipeline))
	}

	return pipelines, nil
}

// setupInputAssemblyState configures input assembly state
func setupInputAssemblyState(cState *C.VkPipelineInputAssemblyStateCreateInfo, info *PipelineInputAssemblyStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	if info != nil {
		cState.topology = C.VkPrimitiveTopology(info.Topology)
		cState.primitiveRestartEnable = boolToVkBool32(info.PrimitiveRestartEnable)
	} else {
		cState.topology = C.VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST
		cState.primitiveRestartEnable = C.VK_FALSE
	}
}

// setupTessellationState configures tessellation state
func setupTessellationState(cState *C.VkPipelineTessellationStateCreateInfo, info *PipelineTessellationStateCreateInfo) {
	cState.sType = C.VK_STRUCTURE_TYPE_PIPELINE_TESSELLATION_STATE_CREATE_INFO
	cState.pNext = nil
	cState.flags = 0
	cState.patchControlPoints = C.uint32_t(info.PatchControlPoints)
}

// ============================================================================
// Framebuffer Management
// ============================================================================

// FramebufferCreateInfo contains framebuffer creation information
type FramebufferCreateInfo struct {
	RenderPass  RenderPass
	Attachments []ImageView
	Width       uint32
	Height      uint32
	Layers      uint32
}

// CreateFramebuffer creates a framebuffer
func CreateFramebuffer(device Device, createInfo *FramebufferCreateInfo) (Framebuffer, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.RenderPass == nil {
		return nil, NewValidationError("RenderPass", "cannot be nil")
	}
	if createInfo.Width == 0 {
		return nil, NewValidationError("Width", "cannot be zero")
	}
	if createInfo.Height == 0 {
		return nil, NewValidationError("Height", "cannot be zero")
	}
	if createInfo.Layers == 0 {
		return nil, NewValidationError("Layers", "cannot be zero")
	}
	// Validate attachments are not nil
	for _, attachment := range createInfo.Attachments {
		if attachment == nil {
			return nil, NewValidationError("Attachments", "contains nil attachment")
		}
	}

	var cCreateInfo C.VkFramebufferCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_FRAMEBUFFER_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.renderPass = C.VkRenderPass(createInfo.RenderPass)
	cCreateInfo.width = C.uint32_t(createInfo.Width)
	cCreateInfo.height = C.uint32_t(createInfo.Height)
	cCreateInfo.layers = C.uint32_t(createInfo.Layers)

	// Handle attachments. The array must be pinned because its address is
	// stored inside cCreateInfo, which is Go memory passed to C.
	var pinner runtime.Pinner
	defer pinner.Unpin()
	var cAttachments []C.VkImageView
	if len(createInfo.Attachments) > 0 {
		cAttachments = make([]C.VkImageView, len(createInfo.Attachments))
		for i, attachment := range createInfo.Attachments {
			cAttachments[i] = C.VkImageView(attachment)
		}
		pinner.Pin(&cAttachments[0])
		cCreateInfo.attachmentCount = C.uint32_t(len(cAttachments))
		cCreateInfo.pAttachments = &cAttachments[0]
	} else {
		cCreateInfo.attachmentCount = 0
		cCreateInfo.pAttachments = nil
	}

	var framebuffer C.VkFramebuffer
	result := Result(C.vkCreateFramebuffer(C.VkDevice(device), &cCreateInfo, nil, &framebuffer))
	if result != Success {
		return nil, NewVulkanError(result, "CreateFramebuffer", "Vulkan framebuffer creation failed")
	}

	trackResource("Framebuffer", unsafe.Pointer(framebuffer))
	return Framebuffer(framebuffer), nil
}

// DestroyFramebuffer destroys a framebuffer
func DestroyFramebuffer(device Device, framebuffer Framebuffer) {
	if device != nil && framebuffer != nil {
		C.vkDestroyFramebuffer(C.VkDevice(device), C.VkFramebuffer(framebuffer), nil)
		untrackResource("Framebuffer", unsafe.Pointer(framebuffer))
	}
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
