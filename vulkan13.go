package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// Vulkan 1.3 Features Implementation

// ============================================================================
// Dynamic Rendering (VK_KHR_dynamic_rendering promoted to core)
// ============================================================================

// RenderingFlags represents flags for dynamic rendering
type RenderingFlags uint32

const (
	RenderingContentsSecondaryCommandBuffers RenderingFlags = C.VK_RENDERING_CONTENTS_SECONDARY_COMMAND_BUFFERS_BIT
	RenderingSuspending                      RenderingFlags = C.VK_RENDERING_SUSPENDING_BIT
	RenderingResuming                        RenderingFlags = C.VK_RENDERING_RESUMING_BIT
)

// RenderingAttachmentInfo describes a single attachment for dynamic rendering
type RenderingAttachmentInfo struct {
	ImageView          ImageView
	ImageLayout        ImageLayout
	ResolveMode        ResolveModeFlagBits
	ResolveImageView   ImageView
	ResolveImageLayout ImageLayout
	LoadOp             AttachmentLoadOp
	StoreOp            AttachmentStoreOp
	ClearValue         ClearValue
}

// RenderingInfo contains information to begin a render pass instance
type RenderingInfo struct {
	Flags             RenderingFlags
	RenderArea        Rect2D
	LayerCount        uint32
	ViewMask          uint32
	ColorAttachments  []RenderingAttachmentInfo
	DepthAttachment   *RenderingAttachmentInfo
	StencilAttachment *RenderingAttachmentInfo
}

// CmdBeginRendering begins a render pass instance with dynamic rendering
func CmdBeginRendering(commandBuffer CommandBuffer, renderingInfo *RenderingInfo) {
	cRenderingInfo := C.VkRenderingInfo{
		sType:      C.VK_STRUCTURE_TYPE_RENDERING_INFO,
		pNext:      nil,
		flags:      C.VkRenderingFlags(renderingInfo.Flags),
		renderArea: *(*C.VkRect2D)(unsafe.Pointer(&renderingInfo.RenderArea)),
		layerCount: C.uint32_t(renderingInfo.LayerCount),
		viewMask:   C.uint32_t(renderingInfo.ViewMask),
	}

	// Handle color attachments
	if len(renderingInfo.ColorAttachments) > 0 {
		cColorAttachments := make([]C.VkRenderingAttachmentInfo, len(renderingInfo.ColorAttachments))
		for i, attachment := range renderingInfo.ColorAttachments {
			cColorAttachments[i] = C.VkRenderingAttachmentInfo{
				sType:              C.VK_STRUCTURE_TYPE_RENDERING_ATTACHMENT_INFO,
				pNext:              nil,
				imageView:          C.VkImageView(attachment.ImageView),
				imageLayout:        C.VkImageLayout(attachment.ImageLayout),
				resolveMode:        C.VkResolveModeFlagBits(attachment.ResolveMode),
				resolveImageView:   C.VkImageView(attachment.ResolveImageView),
				resolveImageLayout: C.VkImageLayout(attachment.ResolveImageLayout),
				loadOp:             C.VkAttachmentLoadOp(attachment.LoadOp),
				storeOp:            C.VkAttachmentStoreOp(attachment.StoreOp),
				clearValue:         *(*C.VkClearValue)(unsafe.Pointer(&attachment.ClearValue)),
			}
		}
		cRenderingInfo.colorAttachmentCount = C.uint32_t(len(cColorAttachments))
		if len(cColorAttachments) > 0 {
			cRenderingInfo.pColorAttachments = &cColorAttachments[0]
		}
	}

	// Handle depth attachment
	if renderingInfo.DepthAttachment != nil {
		cDepthAttachment := C.VkRenderingAttachmentInfo{
			sType:              C.VK_STRUCTURE_TYPE_RENDERING_ATTACHMENT_INFO,
			pNext:              nil,
			imageView:          C.VkImageView(renderingInfo.DepthAttachment.ImageView),
			imageLayout:        C.VkImageLayout(renderingInfo.DepthAttachment.ImageLayout),
			resolveMode:        C.VkResolveModeFlagBits(renderingInfo.DepthAttachment.ResolveMode),
			resolveImageView:   C.VkImageView(renderingInfo.DepthAttachment.ResolveImageView),
			resolveImageLayout: C.VkImageLayout(renderingInfo.DepthAttachment.ResolveImageLayout),
			loadOp:             C.VkAttachmentLoadOp(renderingInfo.DepthAttachment.LoadOp),
			storeOp:            C.VkAttachmentStoreOp(renderingInfo.DepthAttachment.StoreOp),
			clearValue:         *(*C.VkClearValue)(unsafe.Pointer(&renderingInfo.DepthAttachment.ClearValue)),
		}
		cRenderingInfo.pDepthAttachment = &cDepthAttachment
	}

	// Handle stencil attachment
	if renderingInfo.StencilAttachment != nil {
		cStencilAttachment := C.VkRenderingAttachmentInfo{
			sType:              C.VK_STRUCTURE_TYPE_RENDERING_ATTACHMENT_INFO,
			pNext:              nil,
			imageView:          C.VkImageView(renderingInfo.StencilAttachment.ImageView),
			imageLayout:        C.VkImageLayout(renderingInfo.StencilAttachment.ImageLayout),
			resolveMode:        C.VkResolveModeFlagBits(renderingInfo.StencilAttachment.ResolveMode),
			resolveImageView:   C.VkImageView(renderingInfo.StencilAttachment.ResolveImageView),
			resolveImageLayout: C.VkImageLayout(renderingInfo.StencilAttachment.ResolveImageLayout),
			loadOp:             C.VkAttachmentLoadOp(renderingInfo.StencilAttachment.LoadOp),
			storeOp:            C.VkAttachmentStoreOp(renderingInfo.StencilAttachment.StoreOp),
			clearValue:         *(*C.VkClearValue)(unsafe.Pointer(&renderingInfo.StencilAttachment.ClearValue)),
		}
		cRenderingInfo.pStencilAttachment = &cStencilAttachment
	}

	C.vkCmdBeginRendering(C.VkCommandBuffer(commandBuffer), &cRenderingInfo)
}

// CmdEndRendering ends a render pass instance with dynamic rendering
func CmdEndRendering(commandBuffer CommandBuffer) {
	C.vkCmdEndRendering(C.VkCommandBuffer(commandBuffer))
}

// ============================================================================
// Synchronization2 (VK_KHR_synchronization2 promoted to core)
// ============================================================================

// SubmitFlags represents flags for queue submission
type SubmitFlags uint32

const (
	SubmitProtected SubmitFlags = C.VK_SUBMIT_PROTECTED_BIT
)

// SemaphoreSubmitInfo describes a semaphore signal or wait operation
type SemaphoreSubmitInfo struct {
	Semaphore   Semaphore
	Value       uint64
	StageMask   PipelineStageFlags2
	DeviceIndex uint32
}

// CommandBufferSubmitInfo describes a command buffer submit operation
type CommandBufferSubmitInfo struct {
	CommandBuffer CommandBuffer
	DeviceMask    uint32
}

// SubmitInfo2 describes a queue submission operation with enhanced synchronization
type SubmitInfo2 struct {
	Flags                SubmitFlags
	WaitSemaphoreInfos   []SemaphoreSubmitInfo
	CommandBufferInfos   []CommandBufferSubmitInfo
	SignalSemaphoreInfos []SemaphoreSubmitInfo
}

// PipelineStageFlags2 represents enhanced pipeline stage flags
type PipelineStageFlags2 uint64

const (
	PipelineStage2None                         PipelineStageFlags2 = 0
	PipelineStage2TopOfPipe                    PipelineStageFlags2 = 0x00000001
	PipelineStage2DrawIndirect                 PipelineStageFlags2 = 0x00000002
	PipelineStage2VertexInput                  PipelineStageFlags2 = 0x00000004
	PipelineStage2VertexShader                 PipelineStageFlags2 = 0x00000008
	PipelineStage2TessellationControlShader    PipelineStageFlags2 = 0x00000010
	PipelineStage2TessellationEvaluationShader PipelineStageFlags2 = 0x00000020
	PipelineStage2GeometryShader               PipelineStageFlags2 = 0x00000040
	PipelineStage2FragmentShader               PipelineStageFlags2 = 0x00000080
	PipelineStage2EarlyFragmentTests           PipelineStageFlags2 = 0x00000100
	PipelineStage2LateFragmentTests            PipelineStageFlags2 = 0x00000200
	PipelineStage2ColorAttachmentOutput        PipelineStageFlags2 = 0x00000400
	PipelineStage2ComputeShader                PipelineStageFlags2 = 0x00000800
	PipelineStage2AllTransfer                  PipelineStageFlags2 = 0x00001000
	PipelineStage2BottomOfPipe                 PipelineStageFlags2 = 0x00002000
	PipelineStage2Host                         PipelineStageFlags2 = 0x00004000
	PipelineStage2AllGraphics                  PipelineStageFlags2 = 0x00008000
	PipelineStage2AllCommands                  PipelineStageFlags2 = 0x00010000
	PipelineStage2Copy                         PipelineStageFlags2 = 0x100000000
	PipelineStage2Resolve                      PipelineStageFlags2 = 0x200000000
	PipelineStage2Blit                         PipelineStageFlags2 = 0x400000000
	PipelineStage2Clear                        PipelineStageFlags2 = 0x800000000
	PipelineStage2IndexInput                   PipelineStageFlags2 = 0x1000000000
	PipelineStage2VertexAttributeInput         PipelineStageFlags2 = 0x2000000000
	PipelineStage2PreRasterizationShaders      PipelineStageFlags2 = 0x4000000000
)

// QueueSubmit2 submits command buffers to a queue with enhanced synchronization
func QueueSubmit2(queue Queue, submitInfos []SubmitInfo2, fence Fence) error {
	var cSubmitInfos []C.VkSubmitInfo2
	if len(submitInfos) > 0 {
		cSubmitInfos = make([]C.VkSubmitInfo2, len(submitInfos))

		for i, submitInfo := range submitInfos {
			cSubmitInfos[i].sType = C.VK_STRUCTURE_TYPE_SUBMIT_INFO_2
			cSubmitInfos[i].pNext = nil
			cSubmitInfos[i].flags = C.VkSubmitFlags(submitInfo.Flags)

			// Handle wait semaphores
			if len(submitInfo.WaitSemaphoreInfos) > 0 {
				cWaitSemaphoreInfos := make([]C.VkSemaphoreSubmitInfo, len(submitInfo.WaitSemaphoreInfos))
				for j, waitInfo := range submitInfo.WaitSemaphoreInfos {
					cWaitSemaphoreInfos[j] = C.VkSemaphoreSubmitInfo{
						sType:       C.VK_STRUCTURE_TYPE_SEMAPHORE_SUBMIT_INFO,
						pNext:       nil,
						semaphore:   C.VkSemaphore(waitInfo.Semaphore),
						value:       C.uint64_t(waitInfo.Value),
						stageMask:   C.VkPipelineStageFlags2(waitInfo.StageMask),
						deviceIndex: C.uint32_t(waitInfo.DeviceIndex),
					}
				}
				cSubmitInfos[i].waitSemaphoreInfoCount = C.uint32_t(len(cWaitSemaphoreInfos))
				if len(cWaitSemaphoreInfos) > 0 {
					cSubmitInfos[i].pWaitSemaphoreInfos = &cWaitSemaphoreInfos[0]
				}
			}

			// Handle command buffers
			if len(submitInfo.CommandBufferInfos) > 0 {
				cCommandBufferInfos := make([]C.VkCommandBufferSubmitInfo, len(submitInfo.CommandBufferInfos))
				for j, cmdInfo := range submitInfo.CommandBufferInfos {
					cCommandBufferInfos[j] = C.VkCommandBufferSubmitInfo{
						sType:         C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_SUBMIT_INFO,
						pNext:         nil,
						commandBuffer: C.VkCommandBuffer(cmdInfo.CommandBuffer),
						deviceMask:    C.uint32_t(cmdInfo.DeviceMask),
					}
				}
				cSubmitInfos[i].commandBufferInfoCount = C.uint32_t(len(cCommandBufferInfos))
				if len(cCommandBufferInfos) > 0 {
					cSubmitInfos[i].pCommandBufferInfos = &cCommandBufferInfos[0]
				}
			}

			// Handle signal semaphores
			if len(submitInfo.SignalSemaphoreInfos) > 0 {
				cSignalSemaphoreInfos := make([]C.VkSemaphoreSubmitInfo, len(submitInfo.SignalSemaphoreInfos))
				for j, signalInfo := range submitInfo.SignalSemaphoreInfos {
					cSignalSemaphoreInfos[j] = C.VkSemaphoreSubmitInfo{
						sType:       C.VK_STRUCTURE_TYPE_SEMAPHORE_SUBMIT_INFO,
						pNext:       nil,
						semaphore:   C.VkSemaphore(signalInfo.Semaphore),
						value:       C.uint64_t(signalInfo.Value),
						stageMask:   C.VkPipelineStageFlags2(signalInfo.StageMask),
						deviceIndex: C.uint32_t(signalInfo.DeviceIndex),
					}
				}
				cSubmitInfos[i].signalSemaphoreInfoCount = C.uint32_t(len(cSignalSemaphoreInfos))
				if len(cSignalSemaphoreInfos) > 0 {
					cSubmitInfos[i].pSignalSemaphoreInfos = &cSignalSemaphoreInfos[0]
				}
			}
		}
	}

	var pSubmitInfos *C.VkSubmitInfo2
	if len(cSubmitInfos) > 0 {
		pSubmitInfos = &cSubmitInfos[0]
	}

	result := C.vkQueueSubmit2(
		C.VkQueue(queue),
		C.uint32_t(len(cSubmitInfos)),
		pSubmitInfos,
		C.VkFence(fence),
	)

	if result != C.VK_SUCCESS {
		return Result(result)
	}
	return nil
}

// ============================================================================
// Extended Dynamic State (VK_EXT_extended_dynamic_state promoted to core)
// ============================================================================

// Additional dynamic state commands that were promoted in Vulkan 1.3

// CmdSetCullMode sets the cull mode dynamically
func CmdSetCullMode(commandBuffer CommandBuffer, cullMode CullModeFlags) {
	C.vkCmdSetCullMode(C.VkCommandBuffer(commandBuffer), C.VkCullModeFlags(cullMode))
}

// CmdSetFrontFace sets the front face orientation dynamically
func CmdSetFrontFace(commandBuffer CommandBuffer, frontFace FrontFace) {
	C.vkCmdSetFrontFace(C.VkCommandBuffer(commandBuffer), C.VkFrontFace(frontFace))
}

// CmdSetPrimitiveTopology sets the primitive topology dynamically
func CmdSetPrimitiveTopology(commandBuffer CommandBuffer, primitiveTopology PrimitiveTopology) {
	C.vkCmdSetPrimitiveTopology(C.VkCommandBuffer(commandBuffer), C.VkPrimitiveTopology(primitiveTopology))
}

// CmdSetViewportWithCount sets viewports with count dynamically
func CmdSetViewportWithCount(commandBuffer CommandBuffer, viewports []Viewport) {
	if len(viewports) == 0 {
		return
	}

	cViewports := make([]C.VkViewport, len(viewports))
	for i, viewport := range viewports {
		cViewports[i] = *(*C.VkViewport)(unsafe.Pointer(&viewport))
	}

	C.vkCmdSetViewportWithCount(
		C.VkCommandBuffer(commandBuffer),
		C.uint32_t(len(cViewports)),
		&cViewports[0],
	)
}

// CmdSetScissorWithCount sets scissor rectangles with count dynamically
func CmdSetScissorWithCount(commandBuffer CommandBuffer, scissors []Rect2D) {
	if len(scissors) == 0 {
		return
	}

	cScissors := make([]C.VkRect2D, len(scissors))
	for i, scissor := range scissors {
		cScissors[i] = *(*C.VkRect2D)(unsafe.Pointer(&scissor))
	}

	C.vkCmdSetScissorWithCount(
		C.VkCommandBuffer(commandBuffer),
		C.uint32_t(len(cScissors)),
		&cScissors[0],
	)
}

// CmdBindVertexBuffers2 binds vertex buffers with extended parameters
func CmdBindVertexBuffers2(commandBuffer CommandBuffer, firstBinding uint32, buffers []Buffer, offsets []DeviceSize, sizes []DeviceSize, strides []DeviceSize) {
	if len(buffers) == 0 {
		return
	}

	cBuffers := make([]C.VkBuffer, len(buffers))
	for i, buffer := range buffers {
		cBuffers[i] = C.VkBuffer(buffer)
	}

	var cOffsets []C.VkDeviceSize
	if len(offsets) > 0 {
		cOffsets = make([]C.VkDeviceSize, len(offsets))
		for i, offset := range offsets {
			cOffsets[i] = C.VkDeviceSize(offset)
		}
	}

	var cSizes []C.VkDeviceSize
	if len(sizes) > 0 {
		cSizes = make([]C.VkDeviceSize, len(sizes))
		for i, size := range sizes {
			cSizes[i] = C.VkDeviceSize(size)
		}
	}

	var cStrides []C.VkDeviceSize
	if len(strides) > 0 {
		cStrides = make([]C.VkDeviceSize, len(strides))
		for i, stride := range strides {
			cStrides[i] = C.VkDeviceSize(stride)
		}
	}

	var pOffsets *C.VkDeviceSize
	if len(cOffsets) > 0 {
		pOffsets = &cOffsets[0]
	}

	var pSizes *C.VkDeviceSize
	if len(cSizes) > 0 {
		pSizes = &cSizes[0]
	}

	var pStrides *C.VkDeviceSize
	if len(cStrides) > 0 {
		pStrides = &cStrides[0]
	}

	C.vkCmdBindVertexBuffers2(
		C.VkCommandBuffer(commandBuffer),
		C.uint32_t(firstBinding),
		C.uint32_t(len(cBuffers)),
		&cBuffers[0],
		pOffsets,
		pSizes,
		pStrides,
	)
}

// CmdSetDepthTestEnable sets depth test enable state dynamically
func CmdSetDepthTestEnable(commandBuffer CommandBuffer, depthTestEnable bool) {
	C.vkCmdSetDepthTestEnable(C.VkCommandBuffer(commandBuffer), boolToVkBool32(depthTestEnable))
}

// CmdSetDepthWriteEnable sets depth write enable state dynamically
func CmdSetDepthWriteEnable(commandBuffer CommandBuffer, depthWriteEnable bool) {
	C.vkCmdSetDepthWriteEnable(C.VkCommandBuffer(commandBuffer), boolToVkBool32(depthWriteEnable))
}

// CmdSetDepthCompareOp sets depth compare operation dynamically
func CmdSetDepthCompareOp(commandBuffer CommandBuffer, depthCompareOp CompareOp) {
	C.vkCmdSetDepthCompareOp(C.VkCommandBuffer(commandBuffer), C.VkCompareOp(depthCompareOp))
}

// CmdSetDepthBoundsTestEnable sets depth bounds test enable state dynamically
func CmdSetDepthBoundsTestEnable(commandBuffer CommandBuffer, depthBoundsTestEnable bool) {
	C.vkCmdSetDepthBoundsTestEnable(C.VkCommandBuffer(commandBuffer), boolToVkBool32(depthBoundsTestEnable))
}

// CmdSetStencilTestEnable sets stencil test enable state dynamically
func CmdSetStencilTestEnable(commandBuffer CommandBuffer, stencilTestEnable bool) {
	C.vkCmdSetStencilTestEnable(C.VkCommandBuffer(commandBuffer), boolToVkBool32(stencilTestEnable))
}

// CmdSetStencilOp sets stencil operation dynamically
func CmdSetStencilOp(commandBuffer CommandBuffer, faceMask StencilFaceFlags, failOp, passOp, depthFailOp StencilOp, compareOp CompareOp) {
	C.vkCmdSetStencilOp(
		C.VkCommandBuffer(commandBuffer),
		C.VkStencilFaceFlags(faceMask),
		C.VkStencilOp(failOp),
		C.VkStencilOp(passOp),
		C.VkStencilOp(depthFailOp),
		C.VkCompareOp(compareOp),
	)
}

// ============================================================================
// Private Data (VK_EXT_private_data promoted to core)
// ============================================================================

// PrivateDataSlotCreateFlags represents flags for private data slot creation
type PrivateDataSlotCreateFlags uint32

// PrivateDataSlotCreateInfo contains information for creating a private data slot
type PrivateDataSlotCreateInfo struct {
	Flags PrivateDataSlotCreateFlags
}

// CreatePrivateDataSlot creates a private data slot
func CreatePrivateDataSlot(device Device, createInfo *PrivateDataSlotCreateInfo) (PrivateDataSlot, error) {
	cCreateInfo := C.VkPrivateDataSlotCreateInfo{
		sType: C.VK_STRUCTURE_TYPE_PRIVATE_DATA_SLOT_CREATE_INFO,
		pNext: nil,
		flags: C.VkPrivateDataSlotCreateFlags(createInfo.Flags),
	}

	var cPrivateDataSlot C.VkPrivateDataSlot
	result := C.vkCreatePrivateDataSlot(
		C.VkDevice(device),
		&cCreateInfo,
		nil,
		&cPrivateDataSlot,
	)

	if result != C.VK_SUCCESS {
		return PrivateDataSlot(uintptr(0)), Result(result)
	}

	return PrivateDataSlot(cPrivateDataSlot), nil
}

// DestroyPrivateDataSlot destroys a private data slot
func DestroyPrivateDataSlot(device Device, privateDataSlot PrivateDataSlot) {
	C.vkDestroyPrivateDataSlot(
		C.VkDevice(device),
		C.VkPrivateDataSlot(privateDataSlot),
		nil,
	)
}

// SetPrivateData associates data with a Vulkan object
func SetPrivateData(device Device, objectType ObjectType, objectHandle uint64, privateDataSlot PrivateDataSlot, data uint64) error {
	result := C.vkSetPrivateData(
		C.VkDevice(device),
		C.VkObjectType(objectType),
		C.uint64_t(objectHandle),
		C.VkPrivateDataSlot(privateDataSlot),
		C.uint64_t(data),
	)

	if result != C.VK_SUCCESS {
		return Result(result)
	}

	return nil
}

// GetPrivateData retrieves data associated with a Vulkan object
func GetPrivateData(device Device, objectType ObjectType, objectHandle uint64, privateDataSlot PrivateDataSlot) uint64 {
	var data C.uint64_t
	C.vkGetPrivateData(
		C.VkDevice(device),
		C.VkObjectType(objectType),
		C.uint64_t(objectHandle),
		C.VkPrivateDataSlot(privateDataSlot),
		&data,
	)

	return uint64(data)
}

// ============================================================================
// Pipeline Creation Feedback (VK_EXT_pipeline_creation_feedback promoted to core)
// ============================================================================

// PipelineCreationFeedbackFlags represents pipeline creation feedback flags
type PipelineCreationFeedbackFlags uint32

const (
	PipelineCreationFeedbackValid                       PipelineCreationFeedbackFlags = C.VK_PIPELINE_CREATION_FEEDBACK_VALID_BIT
	PipelineCreationFeedbackApplicationPipelineCacheHit PipelineCreationFeedbackFlags = C.VK_PIPELINE_CREATION_FEEDBACK_APPLICATION_PIPELINE_CACHE_HIT_BIT
	PipelineCreationFeedbackBasePipelineAcceleration    PipelineCreationFeedbackFlags = C.VK_PIPELINE_CREATION_FEEDBACK_BASE_PIPELINE_ACCELERATION_BIT
)

// PipelineCreationFeedback provides feedback about pipeline creation
type PipelineCreationFeedback struct {
	Flags    PipelineCreationFeedbackFlags
	Duration uint64
}

// PipelineCreationFeedbackCreateInfo contains pipeline creation feedback information
type PipelineCreationFeedbackCreateInfo struct {
	PipelineCreationFeedback       *PipelineCreationFeedback
	PipelineStageCreationFeedbacks []PipelineCreationFeedback
}

// ============================================================================
// Maintenance4 (VK_KHR_maintenance4 promoted to core)
// ============================================================================

// GetDeviceBufferMemoryRequirements gets buffer memory requirements without creating a buffer (Vulkan 1.3)
func GetDeviceBufferMemoryRequirements(device Device, bufferCreateInfo *BufferCreateInfo) MemoryRequirements {
	cBufferCreateInfo := C.VkBufferCreateInfo{
		sType:                 C.VK_STRUCTURE_TYPE_BUFFER_CREATE_INFO,
		pNext:                 nil,
		flags:                 C.VkBufferCreateFlags(bufferCreateInfo.Flags),
		size:                  C.VkDeviceSize(bufferCreateInfo.Size),
		usage:                 C.VkBufferUsageFlags(bufferCreateInfo.Usage),
		sharingMode:           C.VkSharingMode(bufferCreateInfo.SharingMode),
		queueFamilyIndexCount: 0,
		pQueueFamilyIndices:   nil,
	}

	cDeviceBufferMemoryRequirements := C.VkDeviceBufferMemoryRequirements{
		sType:       C.VK_STRUCTURE_TYPE_DEVICE_BUFFER_MEMORY_REQUIREMENTS,
		pNext:       nil,
		pCreateInfo: &cBufferCreateInfo,
	}

	var cMemoryRequirements C.VkMemoryRequirements2
	cMemoryRequirements.sType = C.VK_STRUCTURE_TYPE_MEMORY_REQUIREMENTS_2
	cMemoryRequirements.pNext = nil

	C.vkGetDeviceBufferMemoryRequirements(
		C.VkDevice(device),
		&cDeviceBufferMemoryRequirements,
		&cMemoryRequirements,
	)

	return MemoryRequirements{
		Size:           DeviceSize(cMemoryRequirements.memoryRequirements.size),
		Alignment:      DeviceSize(cMemoryRequirements.memoryRequirements.alignment),
		MemoryTypeBits: uint32(cMemoryRequirements.memoryRequirements.memoryTypeBits),
	}
}

// GetDeviceImageMemoryRequirements gets image memory requirements without creating an image (Vulkan 1.3)
func GetDeviceImageMemoryRequirements(device Device, imageCreateInfo *ImageCreateInfo) MemoryRequirements {
	cImageCreateInfo := C.VkImageCreateInfo{
		sType:                 C.VK_STRUCTURE_TYPE_IMAGE_CREATE_INFO,
		pNext:                 nil,
		flags:                 C.VkImageCreateFlags(imageCreateInfo.Flags),
		imageType:             C.VkImageType(imageCreateInfo.ImageType),
		format:                C.VkFormat(imageCreateInfo.Format),
		mipLevels:             C.uint32_t(imageCreateInfo.MipLevels),
		arrayLayers:           C.uint32_t(imageCreateInfo.ArrayLayers),
		samples:               C.VkSampleCountFlagBits(imageCreateInfo.Samples),
		tiling:                C.VkImageTiling(imageCreateInfo.Tiling),
		usage:                 C.VkImageUsageFlags(imageCreateInfo.Usage),
		sharingMode:           C.VkSharingMode(imageCreateInfo.SharingMode),
		queueFamilyIndexCount: 0,
		pQueueFamilyIndices:   nil,
		initialLayout:         C.VkImageLayout(imageCreateInfo.InitialLayout),
	}

	// Set extent
	cImageCreateInfo.extent.width = C.uint32_t(imageCreateInfo.Extent.Width)
	cImageCreateInfo.extent.height = C.uint32_t(imageCreateInfo.Extent.Height)
	cImageCreateInfo.extent.depth = C.uint32_t(imageCreateInfo.Extent.Depth)

	cDeviceImageMemoryRequirements := C.VkDeviceImageMemoryRequirements{
		sType:       C.VK_STRUCTURE_TYPE_DEVICE_IMAGE_MEMORY_REQUIREMENTS,
		pNext:       nil,
		pCreateInfo: &cImageCreateInfo,
		planeAspect: C.VK_IMAGE_ASPECT_COLOR_BIT,
	}

	var cMemoryRequirements2 C.VkMemoryRequirements2
	cMemoryRequirements2.sType = C.VK_STRUCTURE_TYPE_MEMORY_REQUIREMENTS_2
	cMemoryRequirements2.pNext = nil

	C.vkGetDeviceImageMemoryRequirements(
		C.VkDevice(device),
		&cDeviceImageMemoryRequirements,
		&cMemoryRequirements2,
	)

	return MemoryRequirements{
		Size:           DeviceSize(cMemoryRequirements2.memoryRequirements.size),
		Alignment:      DeviceSize(cMemoryRequirements2.memoryRequirements.alignment),
		MemoryTypeBits: uint32(cMemoryRequirements2.memoryRequirements.memoryTypeBits),
	}
}
