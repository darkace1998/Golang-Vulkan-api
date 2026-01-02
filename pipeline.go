package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

import (
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

// ComputePipelineCreateInfo contains compute pipeline creation information
type ComputePipelineCreateInfo struct {
	Stage  PipelineShaderStageCreateInfo
	Layout PipelineLayout
}

// CreateComputePipelines creates compute pipelines
func CreateComputePipelines(device Device, pipelineCache PipelineCache, createInfos []ComputePipelineCreateInfo) ([]Pipeline, error) {
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
		return nil, result
	}

	pipelines := make([]Pipeline, len(cPipelines))
	for i, pipeline := range cPipelines {
		pipelines[i] = Pipeline(pipeline)
	}

	return pipelines, nil
}

// DestroyPipeline destroys a pipeline
func DestroyPipeline(device Device, pipeline Pipeline) {
	C.vkDestroyPipeline(C.VkDevice(device), C.VkPipeline(pipeline), nil)
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

// CreateGraphicsPipelines creates graphics pipelines
func CreateGraphicsPipelines(device Device, pipelineCache PipelineCache, createInfos []GraphicsPipelineCreateInfo) ([]Pipeline, error) {
	if len(createInfos) == 0 {
		return nil, nil
	}

	cCreateInfos := make([]C.VkGraphicsPipelineCreateInfo, len(createInfos))
	cPipelines := make([]C.VkPipeline, len(createInfos))

	// Storage for C strings and arrays that need to stay alive during the API call
	var cNames []*C.char
	var cStageArrays [][]C.VkPipelineShaderStageCreateInfo
	var cBindingArrays [][]C.VkVertexInputBindingDescription
	var cAttributeArrays [][]C.VkVertexInputAttributeDescription
	var cViewportArrays [][]C.VkViewport
	var cScissorArrays [][]C.VkRect2D
	var cBlendAttachmentArrays [][]C.VkPipelineColorBlendAttachmentState
	var cDynamicStateArrays [][]C.VkDynamicState
	var cSampleMaskArrays [][]C.VkSampleMask

	// Storage for state structs
	cVertexInputStates := make([]C.VkPipelineVertexInputStateCreateInfo, len(createInfos))
	cInputAssemblyStates := make([]C.VkPipelineInputAssemblyStateCreateInfo, len(createInfos))
	cTessellationStates := make([]C.VkPipelineTessellationStateCreateInfo, len(createInfos))
	cViewportStates := make([]C.VkPipelineViewportStateCreateInfo, len(createInfos))
	cRasterizationStates := make([]C.VkPipelineRasterizationStateCreateInfo, len(createInfos))
	cMultisampleStates := make([]C.VkPipelineMultisampleStateCreateInfo, len(createInfos))
	cDepthStencilStates := make([]C.VkPipelineDepthStencilStateCreateInfo, len(createInfos))
	cColorBlendStates := make([]C.VkPipelineColorBlendStateCreateInfo, len(createInfos))
	cDynamicStates := make([]C.VkPipelineDynamicStateCreateInfo, len(createInfos))

	for i, info := range createInfos {
		cCreateInfos[i].sType = C.VK_STRUCTURE_TYPE_GRAPHICS_PIPELINE_CREATE_INFO
		cCreateInfos[i].pNext = nil
		cCreateInfos[i].flags = 0

		// Shader stages
		if len(info.Stages) > 0 {
			cStages := make([]C.VkPipelineShaderStageCreateInfo, len(info.Stages))
			for j, stage := range info.Stages {
				cStages[j].sType = C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO
				cStages[j].pNext = nil
				cStages[j].flags = 0
				cStages[j].stage = C.VkShaderStageFlagBits(stage.Stage)
				cStages[j].module = C.VkShaderModule(stage.Module)
				cName := C.CString(stage.Name)
				cNames = append(cNames, cName)
				cStages[j].pName = cName
				cStages[j].pSpecializationInfo = nil
			}
			cStageArrays = append(cStageArrays, cStages)
			cCreateInfos[i].stageCount = C.uint32_t(len(cStages))
			cCreateInfos[i].pStages = &cStageArrays[len(cStageArrays)-1][0]
		}

		// Vertex input state
		cVertexInputStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_VERTEX_INPUT_STATE_CREATE_INFO
		cVertexInputStates[i].pNext = nil
		cVertexInputStates[i].flags = 0
		if info.VertexInputState != nil {
			if len(info.VertexInputState.VertexBindingDescriptions) > 0 {
				cBindings := make([]C.VkVertexInputBindingDescription, len(info.VertexInputState.VertexBindingDescriptions))
				for j, binding := range info.VertexInputState.VertexBindingDescriptions {
					cBindings[j].binding = C.uint32_t(binding.Binding)
					cBindings[j].stride = C.uint32_t(binding.Stride)
					cBindings[j].inputRate = C.VkVertexInputRate(binding.InputRate)
				}
				cBindingArrays = append(cBindingArrays, cBindings)
				cVertexInputStates[i].vertexBindingDescriptionCount = C.uint32_t(len(cBindings))
				cVertexInputStates[i].pVertexBindingDescriptions = &cBindingArrays[len(cBindingArrays)-1][0]
			}
			if len(info.VertexInputState.VertexAttributeDescriptions) > 0 {
				cAttributes := make([]C.VkVertexInputAttributeDescription, len(info.VertexInputState.VertexAttributeDescriptions))
				for j, attr := range info.VertexInputState.VertexAttributeDescriptions {
					cAttributes[j].location = C.uint32_t(attr.Location)
					cAttributes[j].binding = C.uint32_t(attr.Binding)
					cAttributes[j].format = C.VkFormat(attr.Format)
					cAttributes[j].offset = C.uint32_t(attr.Offset)
				}
				cAttributeArrays = append(cAttributeArrays, cAttributes)
				cVertexInputStates[i].vertexAttributeDescriptionCount = C.uint32_t(len(cAttributes))
				cVertexInputStates[i].pVertexAttributeDescriptions = &cAttributeArrays[len(cAttributeArrays)-1][0]
			}
		}
		cCreateInfos[i].pVertexInputState = &cVertexInputStates[i]

		// Input assembly state
		cInputAssemblyStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_INPUT_ASSEMBLY_STATE_CREATE_INFO
		cInputAssemblyStates[i].pNext = nil
		cInputAssemblyStates[i].flags = 0
		if info.InputAssemblyState != nil {
			cInputAssemblyStates[i].topology = C.VkPrimitiveTopology(info.InputAssemblyState.Topology)
			cInputAssemblyStates[i].primitiveRestartEnable = boolToVkBool32(info.InputAssemblyState.PrimitiveRestartEnable)
		} else {
			cInputAssemblyStates[i].topology = C.VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST
			cInputAssemblyStates[i].primitiveRestartEnable = C.VK_FALSE
		}
		cCreateInfos[i].pInputAssemblyState = &cInputAssemblyStates[i]

		// Tessellation state (optional)
		if info.TessellationState != nil {
			cTessellationStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_TESSELLATION_STATE_CREATE_INFO
			cTessellationStates[i].pNext = nil
			cTessellationStates[i].flags = 0
			cTessellationStates[i].patchControlPoints = C.uint32_t(info.TessellationState.PatchControlPoints)
			cCreateInfos[i].pTessellationState = &cTessellationStates[i]
		}

		// Viewport state
		cViewportStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_VIEWPORT_STATE_CREATE_INFO
		cViewportStates[i].pNext = nil
		cViewportStates[i].flags = 0
		if info.ViewportState != nil {
			if len(info.ViewportState.Viewports) > 0 {
				cViewports := make([]C.VkViewport, len(info.ViewportState.Viewports))
				for j, vp := range info.ViewportState.Viewports {
					cViewports[j].x = C.float(vp.X)
					cViewports[j].y = C.float(vp.Y)
					cViewports[j].width = C.float(vp.Width)
					cViewports[j].height = C.float(vp.Height)
					cViewports[j].minDepth = C.float(vp.MinDepth)
					cViewports[j].maxDepth = C.float(vp.MaxDepth)
				}
				cViewportArrays = append(cViewportArrays, cViewports)
				cViewportStates[i].viewportCount = C.uint32_t(len(cViewports))
				cViewportStates[i].pViewports = &cViewportArrays[len(cViewportArrays)-1][0]
			} else {
				// Dynamic viewport - set count to 1 but no pointer
				cViewportStates[i].viewportCount = 1
				cViewportStates[i].pViewports = nil
			}
			if len(info.ViewportState.Scissors) > 0 {
				cScissors := make([]C.VkRect2D, len(info.ViewportState.Scissors))
				for j, sc := range info.ViewportState.Scissors {
					cScissors[j].offset.x = C.int32_t(sc.Offset.X)
					cScissors[j].offset.y = C.int32_t(sc.Offset.Y)
					cScissors[j].extent.width = C.uint32_t(sc.Extent.Width)
					cScissors[j].extent.height = C.uint32_t(sc.Extent.Height)
				}
				cScissorArrays = append(cScissorArrays, cScissors)
				cViewportStates[i].scissorCount = C.uint32_t(len(cScissors))
				cViewportStates[i].pScissors = &cScissorArrays[len(cScissorArrays)-1][0]
			} else {
				// Dynamic scissor - set count to 1 but no pointer
				cViewportStates[i].scissorCount = 1
				cViewportStates[i].pScissors = nil
			}
		} else {
			// Default viewport state for dynamic viewport/scissor
			cViewportStates[i].viewportCount = 1
			cViewportStates[i].pViewports = nil
			cViewportStates[i].scissorCount = 1
			cViewportStates[i].pScissors = nil
		}
		cCreateInfos[i].pViewportState = &cViewportStates[i]

		// Rasterization state
		cRasterizationStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_RASTERIZATION_STATE_CREATE_INFO
		cRasterizationStates[i].pNext = nil
		cRasterizationStates[i].flags = 0
		if info.RasterizationState != nil {
			cRasterizationStates[i].depthClampEnable = boolToVkBool32(info.RasterizationState.DepthClampEnable)
			cRasterizationStates[i].rasterizerDiscardEnable = boolToVkBool32(info.RasterizationState.RasterizerDiscardEnable)
			cRasterizationStates[i].polygonMode = C.VkPolygonMode(info.RasterizationState.PolygonMode)
			cRasterizationStates[i].cullMode = C.VkCullModeFlags(info.RasterizationState.CullMode)
			cRasterizationStates[i].frontFace = C.VkFrontFace(info.RasterizationState.FrontFace)
			cRasterizationStates[i].depthBiasEnable = boolToVkBool32(info.RasterizationState.DepthBiasEnable)
			cRasterizationStates[i].depthBiasConstantFactor = C.float(info.RasterizationState.DepthBiasConstantFactor)
			cRasterizationStates[i].depthBiasClamp = C.float(info.RasterizationState.DepthBiasClamp)
			cRasterizationStates[i].depthBiasSlopeFactor = C.float(info.RasterizationState.DepthBiasSlopeFactor)
			cRasterizationStates[i].lineWidth = C.float(info.RasterizationState.LineWidth)
		} else {
			// Default rasterization state
			cRasterizationStates[i].depthClampEnable = C.VK_FALSE
			cRasterizationStates[i].rasterizerDiscardEnable = C.VK_FALSE
			cRasterizationStates[i].polygonMode = C.VK_POLYGON_MODE_FILL
			cRasterizationStates[i].cullMode = C.VkCullModeFlags(C.VK_CULL_MODE_BACK_BIT)
			cRasterizationStates[i].frontFace = C.VK_FRONT_FACE_COUNTER_CLOCKWISE
			cRasterizationStates[i].depthBiasEnable = C.VK_FALSE
			cRasterizationStates[i].lineWidth = 1.0
		}
		cCreateInfos[i].pRasterizationState = &cRasterizationStates[i]

		// Multisample state
		cMultisampleStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_MULTISAMPLE_STATE_CREATE_INFO
		cMultisampleStates[i].pNext = nil
		cMultisampleStates[i].flags = 0
		if info.MultisampleState != nil {
			cMultisampleStates[i].rasterizationSamples = C.VkSampleCountFlagBits(info.MultisampleState.RasterizationSamples)
			cMultisampleStates[i].sampleShadingEnable = boolToVkBool32(info.MultisampleState.SampleShadingEnable)
			cMultisampleStates[i].minSampleShading = C.float(info.MultisampleState.MinSampleShading)
			if len(info.MultisampleState.SampleMask) > 0 {
				cSampleMask := make([]C.VkSampleMask, len(info.MultisampleState.SampleMask))
				for j, mask := range info.MultisampleState.SampleMask {
					cSampleMask[j] = C.VkSampleMask(mask)
				}
				cSampleMaskArrays = append(cSampleMaskArrays, cSampleMask)
				cMultisampleStates[i].pSampleMask = &cSampleMaskArrays[len(cSampleMaskArrays)-1][0]
			}
			cMultisampleStates[i].alphaToCoverageEnable = boolToVkBool32(info.MultisampleState.AlphaToCoverageEnable)
			cMultisampleStates[i].alphaToOneEnable = boolToVkBool32(info.MultisampleState.AlphaToOneEnable)
		} else {
			// Default multisample state
			cMultisampleStates[i].rasterizationSamples = C.VK_SAMPLE_COUNT_1_BIT
			cMultisampleStates[i].sampleShadingEnable = C.VK_FALSE
			cMultisampleStates[i].minSampleShading = 1.0
			cMultisampleStates[i].pSampleMask = nil
			cMultisampleStates[i].alphaToCoverageEnable = C.VK_FALSE
			cMultisampleStates[i].alphaToOneEnable = C.VK_FALSE
		}
		cCreateInfos[i].pMultisampleState = &cMultisampleStates[i]

		// Depth stencil state (optional)
		if info.DepthStencilState != nil {
			cDepthStencilStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_DEPTH_STENCIL_STATE_CREATE_INFO
			cDepthStencilStates[i].pNext = nil
			cDepthStencilStates[i].flags = 0
			cDepthStencilStates[i].depthTestEnable = boolToVkBool32(info.DepthStencilState.DepthTestEnable)
			cDepthStencilStates[i].depthWriteEnable = boolToVkBool32(info.DepthStencilState.DepthWriteEnable)
			cDepthStencilStates[i].depthCompareOp = C.VkCompareOp(info.DepthStencilState.DepthCompareOp)
			cDepthStencilStates[i].depthBoundsTestEnable = boolToVkBool32(info.DepthStencilState.DepthBoundsTestEnable)
			cDepthStencilStates[i].stencilTestEnable = boolToVkBool32(info.DepthStencilState.StencilTestEnable)
			cDepthStencilStates[i].front.failOp = C.VkStencilOp(info.DepthStencilState.Front.FailOp)
			cDepthStencilStates[i].front.passOp = C.VkStencilOp(info.DepthStencilState.Front.PassOp)
			cDepthStencilStates[i].front.depthFailOp = C.VkStencilOp(info.DepthStencilState.Front.DepthFailOp)
			cDepthStencilStates[i].front.compareOp = C.VkCompareOp(info.DepthStencilState.Front.CompareOp)
			cDepthStencilStates[i].front.compareMask = C.uint32_t(info.DepthStencilState.Front.CompareMask)
			cDepthStencilStates[i].front.writeMask = C.uint32_t(info.DepthStencilState.Front.WriteMask)
			cDepthStencilStates[i].front.reference = C.uint32_t(info.DepthStencilState.Front.Reference)
			cDepthStencilStates[i].back.failOp = C.VkStencilOp(info.DepthStencilState.Back.FailOp)
			cDepthStencilStates[i].back.passOp = C.VkStencilOp(info.DepthStencilState.Back.PassOp)
			cDepthStencilStates[i].back.depthFailOp = C.VkStencilOp(info.DepthStencilState.Back.DepthFailOp)
			cDepthStencilStates[i].back.compareOp = C.VkCompareOp(info.DepthStencilState.Back.CompareOp)
			cDepthStencilStates[i].back.compareMask = C.uint32_t(info.DepthStencilState.Back.CompareMask)
			cDepthStencilStates[i].back.writeMask = C.uint32_t(info.DepthStencilState.Back.WriteMask)
			cDepthStencilStates[i].back.reference = C.uint32_t(info.DepthStencilState.Back.Reference)
			cDepthStencilStates[i].minDepthBounds = C.float(info.DepthStencilState.MinDepthBounds)
			cDepthStencilStates[i].maxDepthBounds = C.float(info.DepthStencilState.MaxDepthBounds)
			cCreateInfos[i].pDepthStencilState = &cDepthStencilStates[i]
		}

		// Color blend state
		cColorBlendStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_COLOR_BLEND_STATE_CREATE_INFO
		cColorBlendStates[i].pNext = nil
		cColorBlendStates[i].flags = 0
		if info.ColorBlendState != nil {
			cColorBlendStates[i].logicOpEnable = boolToVkBool32(info.ColorBlendState.LogicOpEnable)
			cColorBlendStates[i].logicOp = C.VkLogicOp(info.ColorBlendState.LogicOp)
			if len(info.ColorBlendState.Attachments) > 0 {
				cBlendAttachments := make([]C.VkPipelineColorBlendAttachmentState, len(info.ColorBlendState.Attachments))
				for j, att := range info.ColorBlendState.Attachments {
					cBlendAttachments[j].blendEnable = boolToVkBool32(att.BlendEnable)
					cBlendAttachments[j].srcColorBlendFactor = C.VkBlendFactor(att.SrcColorBlendFactor)
					cBlendAttachments[j].dstColorBlendFactor = C.VkBlendFactor(att.DstColorBlendFactor)
					cBlendAttachments[j].colorBlendOp = C.VkBlendOp(att.ColorBlendOp)
					cBlendAttachments[j].srcAlphaBlendFactor = C.VkBlendFactor(att.SrcAlphaBlendFactor)
					cBlendAttachments[j].dstAlphaBlendFactor = C.VkBlendFactor(att.DstAlphaBlendFactor)
					cBlendAttachments[j].alphaBlendOp = C.VkBlendOp(att.AlphaBlendOp)
					cBlendAttachments[j].colorWriteMask = C.VkColorComponentFlags(att.ColorWriteMask)
				}
				cBlendAttachmentArrays = append(cBlendAttachmentArrays, cBlendAttachments)
				cColorBlendStates[i].attachmentCount = C.uint32_t(len(cBlendAttachments))
				cColorBlendStates[i].pAttachments = &cBlendAttachmentArrays[len(cBlendAttachmentArrays)-1][0]
			}
			for j := 0; j < 4; j++ {
				cColorBlendStates[i].blendConstants[j] = C.float(info.ColorBlendState.BlendConstants[j])
			}
		} else {
			// Default color blend state - no blending, write all color components
			cColorBlendStates[i].logicOpEnable = C.VK_FALSE
			cColorBlendStates[i].logicOp = C.VK_LOGIC_OP_COPY
			// Create a default attachment for one color attachment
			defaultAttachment := make([]C.VkPipelineColorBlendAttachmentState, 1)
			defaultAttachment[0].blendEnable = C.VK_FALSE
			defaultAttachment[0].colorWriteMask = C.VkColorComponentFlags(ColorComponentAll)
			cBlendAttachmentArrays = append(cBlendAttachmentArrays, defaultAttachment)
			cColorBlendStates[i].attachmentCount = 1
			cColorBlendStates[i].pAttachments = &cBlendAttachmentArrays[len(cBlendAttachmentArrays)-1][0]
		}
		cCreateInfos[i].pColorBlendState = &cColorBlendStates[i]

		// Dynamic state (optional)
		if info.DynamicState != nil && len(info.DynamicState.DynamicStates) > 0 {
			cDynamicStates[i].sType = C.VK_STRUCTURE_TYPE_PIPELINE_DYNAMIC_STATE_CREATE_INFO
			cDynamicStates[i].pNext = nil
			cDynamicStates[i].flags = 0
			cDynStates := make([]C.VkDynamicState, len(info.DynamicState.DynamicStates))
			for j, ds := range info.DynamicState.DynamicStates {
				cDynStates[j] = C.VkDynamicState(ds)
			}
			cDynamicStateArrays = append(cDynamicStateArrays, cDynStates)
			cDynamicStates[i].dynamicStateCount = C.uint32_t(len(cDynStates))
			cDynamicStates[i].pDynamicStates = &cDynamicStateArrays[len(cDynamicStateArrays)-1][0]
			cCreateInfos[i].pDynamicState = &cDynamicStates[i]
		}

		// Layout and render pass
		cCreateInfos[i].layout = C.VkPipelineLayout(info.Layout)
		cCreateInfos[i].renderPass = C.VkRenderPass(info.RenderPass)
		cCreateInfos[i].subpass = C.uint32_t(info.Subpass)
		cCreateInfos[i].basePipelineHandle = C.VkPipeline(info.BasePipelineHandle)
		cCreateInfos[i].basePipelineIndex = C.int32_t(info.BasePipelineIndex)
	}

	// Free all C strings after API call regardless of success/failure
	defer func() {
		for _, cName := range cNames {
			if cName != nil {
				C.free(unsafe.Pointer(cName))
			}
		}
	}()

	result := Result(C.vkCreateGraphicsPipelines(
		C.VkDevice(device),
		C.VkPipelineCache(pipelineCache),
		C.uint32_t(len(cCreateInfos)),
		&cCreateInfos[0],
		nil,
		&cPipelines[0],
	))

	if result != Success {
		return nil, result
	}

	pipelines := make([]Pipeline, len(cPipelines))
	for i, pipeline := range cPipelines {
		pipelines[i] = Pipeline(pipeline)
	}

	return pipelines, nil
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

	// Handle attachments
	var cAttachments []C.VkImageView
	if len(createInfo.Attachments) > 0 {
		cAttachments = make([]C.VkImageView, len(createInfo.Attachments))
		for i, attachment := range createInfo.Attachments {
			cAttachments[i] = C.VkImageView(attachment)
		}
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

	return Framebuffer(framebuffer), nil
}

// DestroyFramebuffer destroys a framebuffer
func DestroyFramebuffer(device Device, framebuffer Framebuffer) {
	if device != nil && framebuffer != nil {
		C.vkDestroyFramebuffer(C.VkDevice(device), C.VkFramebuffer(framebuffer), nil)
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
