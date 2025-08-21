package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
*/
import "C"

// ClearColorValue represents a clear color value
type ClearColorValue struct {
	Float32 [4]float32
	Int32   [4]int32
	Uint32  [4]uint32
}

// ClearDepthStencilValue represents a clear depth/stencil value
type ClearDepthStencilValue struct {
	Depth   float32
	Stencil uint32
}

// ClearValue represents a clear value union
type ClearValue struct {
	Color        ClearColorValue
	DepthStencil ClearDepthStencilValue
}

// RenderPassBeginInfo contains render pass begin information
type RenderPassBeginInfo struct {
	RenderPass  RenderPass
	Framebuffer Framebuffer
	RenderArea  Rect2D
	ClearValues []ClearValue
}

// Rect2D represents a 2D rectangle
type Rect2D struct {
	Offset Offset2D
	Extent Extent2D
}

// Offset2D represents a 2D offset
type Offset2D struct {
	X int32
	Y int32
}

// Extent2D represents a 2D extent
type Extent2D struct {
	Width  uint32
	Height uint32
}

// Viewport represents a viewport
type Viewport struct {
	X        float32
	Y        float32
	Width    float32
	Height   float32
	MinDepth float32
	MaxDepth float32
}

// SubpassContents represents subpass contents
type SubpassContents int32

const (
	SubpassContentsInline                  SubpassContents = C.VK_SUBPASS_CONTENTS_INLINE
	SubpassContentsSecondaryCommandBuffers SubpassContents = C.VK_SUBPASS_CONTENTS_SECONDARY_COMMAND_BUFFERS
)

// CmdBeginRenderPass begins a render pass
func CmdBeginRenderPass(commandBuffer CommandBuffer, beginInfo *RenderPassBeginInfo, contents SubpassContents) {
	var cBeginInfo C.VkRenderPassBeginInfo
	cBeginInfo.sType = C.VK_STRUCTURE_TYPE_RENDER_PASS_BEGIN_INFO
	cBeginInfo.pNext = nil
	cBeginInfo.renderPass = C.VkRenderPass(beginInfo.RenderPass)
	cBeginInfo.framebuffer = C.VkFramebuffer(beginInfo.Framebuffer)
	cBeginInfo.renderArea.offset.x = C.int32_t(beginInfo.RenderArea.Offset.X)
	cBeginInfo.renderArea.offset.y = C.int32_t(beginInfo.RenderArea.Offset.Y)
	cBeginInfo.renderArea.extent.width = C.uint32_t(beginInfo.RenderArea.Extent.Width)
	cBeginInfo.renderArea.extent.height = C.uint32_t(beginInfo.RenderArea.Extent.Height)

	// For simplicity, skip clear values for now - can be added later
	cBeginInfo.clearValueCount = 0
	cBeginInfo.pClearValues = nil

	C.vkCmdBeginRenderPass(C.VkCommandBuffer(commandBuffer), &cBeginInfo, C.VkSubpassContents(contents))
}

// CmdEndRenderPass ends a render pass
func CmdEndRenderPass(commandBuffer CommandBuffer) {
	C.vkCmdEndRenderPass(C.VkCommandBuffer(commandBuffer))
}

// CmdBindPipeline binds a pipeline
func CmdBindPipeline(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, pipeline Pipeline) {
	C.vkCmdBindPipeline(C.VkCommandBuffer(commandBuffer), C.VkPipelineBindPoint(pipelineBindPoint), C.VkPipeline(pipeline))
}

// CmdSetViewport sets the viewport
func CmdSetViewport(commandBuffer CommandBuffer, firstViewport uint32, viewports []Viewport) {
	if len(viewports) == 0 {
		return
	}

	cViewports := make([]C.VkViewport, len(viewports))
	for i, vp := range viewports {
		cViewports[i].x = C.float(vp.X)
		cViewports[i].y = C.float(vp.Y)
		cViewports[i].width = C.float(vp.Width)
		cViewports[i].height = C.float(vp.Height)
		cViewports[i].minDepth = C.float(vp.MinDepth)
		cViewports[i].maxDepth = C.float(vp.MaxDepth)
	}

	C.vkCmdSetViewport(C.VkCommandBuffer(commandBuffer), C.uint32_t(firstViewport), C.uint32_t(len(cViewports)), &cViewports[0])
}

// CmdSetScissor sets the scissor rectangles
func CmdSetScissor(commandBuffer CommandBuffer, firstScissor uint32, scissors []Rect2D) {
	if len(scissors) == 0 {
		return
	}

	cScissors := make([]C.VkRect2D, len(scissors))
	for i, scissor := range scissors {
		cScissors[i].offset.x = C.int32_t(scissor.Offset.X)
		cScissors[i].offset.y = C.int32_t(scissor.Offset.Y)
		cScissors[i].extent.width = C.uint32_t(scissor.Extent.Width)
		cScissors[i].extent.height = C.uint32_t(scissor.Extent.Height)
	}

	C.vkCmdSetScissor(C.VkCommandBuffer(commandBuffer), C.uint32_t(firstScissor), C.uint32_t(len(cScissors)), &cScissors[0])
}

// CmdBindVertexBuffers binds vertex buffers
func CmdBindVertexBuffers(commandBuffer CommandBuffer, firstBinding uint32, buffers []Buffer, offsets []DeviceSize) {
	if len(buffers) == 0 || len(buffers) != len(offsets) {
		return
	}

	cBuffers := make([]C.VkBuffer, len(buffers))
	cOffsets := make([]C.VkDeviceSize, len(offsets))

	for i, buffer := range buffers {
		cBuffers[i] = C.VkBuffer(buffer)
		cOffsets[i] = C.VkDeviceSize(offsets[i])
	}

	C.vkCmdBindVertexBuffers(C.VkCommandBuffer(commandBuffer), C.uint32_t(firstBinding), C.uint32_t(len(cBuffers)), &cBuffers[0], &cOffsets[0])
}

// CmdBindIndexBuffer binds an index buffer
func CmdBindIndexBuffer(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize, indexType IndexType) {
	C.vkCmdBindIndexBuffer(C.VkCommandBuffer(commandBuffer), C.VkBuffer(buffer), C.VkDeviceSize(offset), C.VkIndexType(indexType))
}

// IndexType represents index buffer types
type IndexType int32

const (
	IndexTypeUint16 IndexType = C.VK_INDEX_TYPE_UINT16
	IndexTypeUint32 IndexType = C.VK_INDEX_TYPE_UINT32
)

// CmdDraw records a draw command
func CmdDraw(commandBuffer CommandBuffer, vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	C.vkCmdDraw(C.VkCommandBuffer(commandBuffer), C.uint32_t(vertexCount), C.uint32_t(instanceCount), C.uint32_t(firstVertex), C.uint32_t(firstInstance))
}

// CmdDrawIndexed records an indexed draw command
func CmdDrawIndexed(commandBuffer CommandBuffer, indexCount, instanceCount, firstIndex uint32, vertexOffset int32, firstInstance uint32) {
	C.vkCmdDrawIndexed(C.VkCommandBuffer(commandBuffer), C.uint32_t(indexCount), C.uint32_t(instanceCount), C.uint32_t(firstIndex), C.int32_t(vertexOffset), C.uint32_t(firstInstance))
}

// CmdCopyBuffer copies data between buffers
func CmdCopyBuffer(commandBuffer CommandBuffer, srcBuffer, dstBuffer Buffer, regions []BufferCopy) {
	if len(regions) == 0 {
		return
	}

	cRegions := make([]C.VkBufferCopy, len(regions))
	for i, region := range regions {
		cRegions[i].srcOffset = C.VkDeviceSize(region.SrcOffset)
		cRegions[i].dstOffset = C.VkDeviceSize(region.DstOffset)
		cRegions[i].size = C.VkDeviceSize(region.Size)
	}

	C.vkCmdCopyBuffer(C.VkCommandBuffer(commandBuffer), C.VkBuffer(srcBuffer), C.VkBuffer(dstBuffer), C.uint32_t(len(cRegions)), &cRegions[0])
}

// BufferCopy describes a buffer copy region
type BufferCopy struct {
	SrcOffset DeviceSize
	DstOffset DeviceSize
	Size      DeviceSize
}

// CmdPipelineBarrier inserts a pipeline barrier
func CmdPipelineBarrier(commandBuffer CommandBuffer, srcStageMask, dstStageMask PipelineStageFlags, dependencyFlags uint32) {
	C.vkCmdPipelineBarrier(C.VkCommandBuffer(commandBuffer), C.VkPipelineStageFlags(srcStageMask), C.VkPipelineStageFlags(dstStageMask), C.VkDependencyFlags(dependencyFlags), 0, nil, 0, nil, 0, nil)
}

// Compute dispatch commands

// CmdDispatch dispatches compute work
func CmdDispatch(commandBuffer CommandBuffer, groupCountX, groupCountY, groupCountZ uint32) {
	C.vkCmdDispatch(C.VkCommandBuffer(commandBuffer), C.uint32_t(groupCountX), C.uint32_t(groupCountY), C.uint32_t(groupCountZ))
}

// CmdDispatchIndirect dispatches compute work with parameters from a buffer
func CmdDispatchIndirect(commandBuffer CommandBuffer, buffer Buffer, offset DeviceSize) {
	C.vkCmdDispatchIndirect(C.VkCommandBuffer(commandBuffer), C.VkBuffer(buffer), C.VkDeviceSize(offset))
}

// CmdBindDescriptorSets binds descriptor sets to a command buffer
func CmdBindDescriptorSets(commandBuffer CommandBuffer, pipelineBindPoint PipelineBindPoint, layout PipelineLayout, firstSet uint32, descriptorSets []DescriptorSet, dynamicOffsets []uint32) {
	if len(descriptorSets) == 0 {
		return
	}

	cDescriptorSets := make([]C.VkDescriptorSet, len(descriptorSets))
	for i, set := range descriptorSets {
		cDescriptorSets[i] = C.VkDescriptorSet(set)
	}

	var cDynamicOffsets []C.uint32_t
	if len(dynamicOffsets) > 0 {
		cDynamicOffsets = make([]C.uint32_t, len(dynamicOffsets))
		for i, offset := range dynamicOffsets {
			cDynamicOffsets[i] = C.uint32_t(offset)
		}
	}

	var pDynamicOffsets *C.uint32_t = nil
	if len(cDynamicOffsets) > 0 {
		pDynamicOffsets = &cDynamicOffsets[0]
	}

	C.vkCmdBindDescriptorSets(
		C.VkCommandBuffer(commandBuffer),
		C.VkPipelineBindPoint(pipelineBindPoint),
		C.VkPipelineLayout(layout),
		C.uint32_t(firstSet),
		C.uint32_t(len(cDescriptorSets)),
		&cDescriptorSets[0],
		C.uint32_t(len(cDynamicOffsets)),
		pDynamicOffsets,
	)
}
