package vulkan

/*
#include <vulkan/vulkan.h>
*/
import "C"

import "unsafe"

// ============================================================================
// Image Layout Transition Helpers
// ============================================================================

// TransitionImageLayout transitions an image from one layout to another
// This is a helper function for common layout transitions
func TransitionImageLayout(
	commandBuffer CommandBuffer,
	image Image,
	format Format,
	oldLayout ImageLayout,
	newLayout ImageLayout,
	subresourceRange ImageSubresourceRange,
) {
	if commandBuffer == nil || image == nil {
		return
	}

	var barrier ImageMemoryBarrier
	barrier.OldLayout = oldLayout
	barrier.NewLayout = newLayout
	barrier.SrcQueueFamilyIndex = QueueFamilyIgnored
	barrier.DstQueueFamilyIndex = QueueFamilyIgnored
	barrier.Image = image
	barrier.SubresourceRange = subresourceRange

	var srcStage PipelineStageFlags
	var dstStage PipelineStageFlags

	// Determine access masks and pipeline stages based on the transition
	switch {
	case oldLayout == ImageLayoutUndefined && newLayout == ImageLayoutTransferDstOptimal:
		barrier.SrcAccessMask = 0
		barrier.DstAccessMask = AccessTransferWriteBit
		srcStage = PipelineStageTopOfPipeBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutTransferDstOptimal && newLayout == ImageLayoutShaderReadOnlyOptimal:
		barrier.SrcAccessMask = AccessTransferWriteBit
		barrier.DstAccessMask = AccessShaderReadBit
		srcStage = PipelineStageTransferBit
		dstStage = PipelineStageFragmentShaderBit

	case oldLayout == ImageLayoutUndefined && newLayout == ImageLayoutDepthStencilAttachmentOptimal:
		barrier.SrcAccessMask = 0
		barrier.DstAccessMask = AccessDepthStencilAttachmentReadBit | AccessDepthStencilAttachmentWriteBit
		srcStage = PipelineStageTopOfPipeBit
		dstStage = PipelineStageEarlyFragmentTestsBit

	case oldLayout == ImageLayoutUndefined && newLayout == ImageLayoutColorAttachmentOptimal:
		barrier.SrcAccessMask = 0
		barrier.DstAccessMask = AccessColorAttachmentReadBit | AccessColorAttachmentWriteBit
		srcStage = PipelineStageTopOfPipeBit
		dstStage = PipelineStageColorAttachmentOutputBit

	case oldLayout == ImageLayoutColorAttachmentOptimal && newLayout == ImageLayoutPresentSrcKHR:
		barrier.SrcAccessMask = AccessColorAttachmentWriteBit
		barrier.DstAccessMask = 0
		srcStage = PipelineStageColorAttachmentOutputBit
		dstStage = PipelineStageBottomOfPipeBit

	case oldLayout == ImageLayoutTransferSrcOptimal && newLayout == ImageLayoutShaderReadOnlyOptimal:
		barrier.SrcAccessMask = AccessTransferReadBit
		barrier.DstAccessMask = AccessShaderReadBit
		srcStage = PipelineStageTransferBit
		dstStage = PipelineStageFragmentShaderBit

	case oldLayout == ImageLayoutShaderReadOnlyOptimal && newLayout == ImageLayoutTransferSrcOptimal:
		barrier.SrcAccessMask = AccessShaderReadBit
		barrier.DstAccessMask = AccessTransferReadBit
		srcStage = PipelineStageFragmentShaderBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutShaderReadOnlyOptimal && newLayout == ImageLayoutTransferDstOptimal:
		barrier.SrcAccessMask = AccessShaderReadBit
		barrier.DstAccessMask = AccessTransferWriteBit
		srcStage = PipelineStageFragmentShaderBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutGeneral && newLayout == ImageLayoutTransferSrcOptimal:
		barrier.SrcAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		barrier.DstAccessMask = AccessTransferReadBit
		srcStage = PipelineStageAllCommandsBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutGeneral && newLayout == ImageLayoutTransferDstOptimal:
		barrier.SrcAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		barrier.DstAccessMask = AccessTransferWriteBit
		srcStage = PipelineStageAllCommandsBit
		dstStage = PipelineStageTransferBit

	default:
		// Generic fallback - use ALL_COMMANDS for safety
		barrier.SrcAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		barrier.DstAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		srcStage = PipelineStageAllCommandsBit
		dstStage = PipelineStageAllCommandsBit
	}

	CmdPipelineBarrierFull(commandBuffer, srcStage, dstStage, 0, nil, nil, []ImageMemoryBarrier{barrier})
}

// TransitionImageLayoutSimple transitions an image layout using default settings
// Convenience function that uses the color aspect and full mip/array range
func TransitionImageLayoutSimple(
	commandBuffer CommandBuffer,
	image Image,
	oldLayout ImageLayout,
	newLayout ImageLayout,
) {
	TransitionImageLayout(commandBuffer, image, FormatUndefined, oldLayout, newLayout, ImageSubresourceRange{
		AspectMask:     ImageAspectColorBit,
		BaseMipLevel:   0,
		LevelCount:     1,
		BaseArrayLayer: 0,
		LayerCount:     1,
	})
}

// ============================================================================
// Image Operations
// ============================================================================

// ImageBlit describes an image blit operation
type ImageBlit struct {
	SrcSubresource ImageSubresourceLayers
	SrcOffsets     [2]Offset3D
	DstSubresource ImageSubresourceLayers
	DstOffsets     [2]Offset3D
}

// ImageSubresourceLayers specifies image subresource layers
type ImageSubresourceLayers struct {
	AspectMask     ImageAspectFlags
	MipLevel       uint32
	BaseArrayLayer uint32
	LayerCount     uint32
}

// Offset3D represents a 3D offset
type Offset3D struct {
	X int32
	Y int32
	Z int32
}

// FilterCubic is a cubic filter mode (requires extension)
const FilterCubic Filter = C.VK_FILTER_CUBIC_IMG

// CmdBlitImage copies regions of an image with potential format conversion and scaling
func CmdBlitImage(
	commandBuffer CommandBuffer,
	srcImage Image,
	srcImageLayout ImageLayout,
	dstImage Image,
	dstImageLayout ImageLayout,
	regions []ImageBlit,
	filter Filter,
) {
	if commandBuffer == nil || srcImage == nil || dstImage == nil || len(regions) == 0 {
		return
	}

	cRegions := make([]C.VkImageBlit, len(regions))
	for i, region := range regions {
		cRegions[i].srcSubresource.aspectMask = C.VkImageAspectFlags(region.SrcSubresource.AspectMask)
		cRegions[i].srcSubresource.mipLevel = C.uint32_t(region.SrcSubresource.MipLevel)
		cRegions[i].srcSubresource.baseArrayLayer = C.uint32_t(region.SrcSubresource.BaseArrayLayer)
		cRegions[i].srcSubresource.layerCount = C.uint32_t(region.SrcSubresource.LayerCount)
		cRegions[i].srcOffsets[0].x = C.int32_t(region.SrcOffsets[0].X)
		cRegions[i].srcOffsets[0].y = C.int32_t(region.SrcOffsets[0].Y)
		cRegions[i].srcOffsets[0].z = C.int32_t(region.SrcOffsets[0].Z)
		cRegions[i].srcOffsets[1].x = C.int32_t(region.SrcOffsets[1].X)
		cRegions[i].srcOffsets[1].y = C.int32_t(region.SrcOffsets[1].Y)
		cRegions[i].srcOffsets[1].z = C.int32_t(region.SrcOffsets[1].Z)

		cRegions[i].dstSubresource.aspectMask = C.VkImageAspectFlags(region.DstSubresource.AspectMask)
		cRegions[i].dstSubresource.mipLevel = C.uint32_t(region.DstSubresource.MipLevel)
		cRegions[i].dstSubresource.baseArrayLayer = C.uint32_t(region.DstSubresource.BaseArrayLayer)
		cRegions[i].dstSubresource.layerCount = C.uint32_t(region.DstSubresource.LayerCount)
		cRegions[i].dstOffsets[0].x = C.int32_t(region.DstOffsets[0].X)
		cRegions[i].dstOffsets[0].y = C.int32_t(region.DstOffsets[0].Y)
		cRegions[i].dstOffsets[0].z = C.int32_t(region.DstOffsets[0].Z)
		cRegions[i].dstOffsets[1].x = C.int32_t(region.DstOffsets[1].X)
		cRegions[i].dstOffsets[1].y = C.int32_t(region.DstOffsets[1].Y)
		cRegions[i].dstOffsets[1].z = C.int32_t(region.DstOffsets[1].Z)
	}

	C.vkCmdBlitImage(
		C.VkCommandBuffer(commandBuffer),
		C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout),
		C.VkImage(dstImage),
		C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)),
		&cRegions[0],
		C.VkFilter(filter),
	)
}

// ImageResolve describes an image resolve operation
type ImageResolve struct {
	SrcSubresource ImageSubresourceLayers
	SrcOffset      Offset3D
	DstSubresource ImageSubresourceLayers
	DstOffset      Offset3D
	Extent         Extent3D
}

// CmdResolveImage resolves a multisample image to a non-multisample image
func CmdResolveImage(
	commandBuffer CommandBuffer,
	srcImage Image,
	srcImageLayout ImageLayout,
	dstImage Image,
	dstImageLayout ImageLayout,
	regions []ImageResolve,
) {
	if commandBuffer == nil || srcImage == nil || dstImage == nil || len(regions) == 0 {
		return
	}

	cRegions := make([]C.VkImageResolve, len(regions))
	for i, region := range regions {
		cRegions[i].srcSubresource.aspectMask = C.VkImageAspectFlags(region.SrcSubresource.AspectMask)
		cRegions[i].srcSubresource.mipLevel = C.uint32_t(region.SrcSubresource.MipLevel)
		cRegions[i].srcSubresource.baseArrayLayer = C.uint32_t(region.SrcSubresource.BaseArrayLayer)
		cRegions[i].srcSubresource.layerCount = C.uint32_t(region.SrcSubresource.LayerCount)
		cRegions[i].srcOffset.x = C.int32_t(region.SrcOffset.X)
		cRegions[i].srcOffset.y = C.int32_t(region.SrcOffset.Y)
		cRegions[i].srcOffset.z = C.int32_t(region.SrcOffset.Z)

		cRegions[i].dstSubresource.aspectMask = C.VkImageAspectFlags(region.DstSubresource.AspectMask)
		cRegions[i].dstSubresource.mipLevel = C.uint32_t(region.DstSubresource.MipLevel)
		cRegions[i].dstSubresource.baseArrayLayer = C.uint32_t(region.DstSubresource.BaseArrayLayer)
		cRegions[i].dstSubresource.layerCount = C.uint32_t(region.DstSubresource.LayerCount)
		cRegions[i].dstOffset.x = C.int32_t(region.DstOffset.X)
		cRegions[i].dstOffset.y = C.int32_t(region.DstOffset.Y)
		cRegions[i].dstOffset.z = C.int32_t(region.DstOffset.Z)

		cRegions[i].extent.width = C.uint32_t(region.Extent.Width)
		cRegions[i].extent.height = C.uint32_t(region.Extent.Height)
		cRegions[i].extent.depth = C.uint32_t(region.Extent.Depth)
	}

	C.vkCmdResolveImage(
		C.VkCommandBuffer(commandBuffer),
		C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout),
		C.VkImage(dstImage),
		C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)),
		&cRegions[0],
	)
}

// BufferImageCopy describes a buffer to image or image to buffer copy operation
type BufferImageCopy struct {
	BufferOffset      DeviceSize
	BufferRowLength   uint32
	BufferImageHeight uint32
	ImageSubresource  ImageSubresourceLayers
	ImageOffset       Offset3D
	ImageExtent       Extent3D
}

// CmdCopyBufferToImage copies data from a buffer to an image
func CmdCopyBufferToImage(
	commandBuffer CommandBuffer,
	srcBuffer Buffer,
	dstImage Image,
	dstImageLayout ImageLayout,
	regions []BufferImageCopy,
) {
	if commandBuffer == nil || srcBuffer == nil || dstImage == nil || len(regions) == 0 {
		return
	}

	cRegions := make([]C.VkBufferImageCopy, len(regions))
	for i, region := range regions {
		cRegions[i].bufferOffset = C.VkDeviceSize(region.BufferOffset)
		cRegions[i].bufferRowLength = C.uint32_t(region.BufferRowLength)
		cRegions[i].bufferImageHeight = C.uint32_t(region.BufferImageHeight)
		cRegions[i].imageSubresource.aspectMask = C.VkImageAspectFlags(region.ImageSubresource.AspectMask)
		cRegions[i].imageSubresource.mipLevel = C.uint32_t(region.ImageSubresource.MipLevel)
		cRegions[i].imageSubresource.baseArrayLayer = C.uint32_t(region.ImageSubresource.BaseArrayLayer)
		cRegions[i].imageSubresource.layerCount = C.uint32_t(region.ImageSubresource.LayerCount)
		cRegions[i].imageOffset.x = C.int32_t(region.ImageOffset.X)
		cRegions[i].imageOffset.y = C.int32_t(region.ImageOffset.Y)
		cRegions[i].imageOffset.z = C.int32_t(region.ImageOffset.Z)
		cRegions[i].imageExtent.width = C.uint32_t(region.ImageExtent.Width)
		cRegions[i].imageExtent.height = C.uint32_t(region.ImageExtent.Height)
		cRegions[i].imageExtent.depth = C.uint32_t(region.ImageExtent.Depth)
	}

	C.vkCmdCopyBufferToImage(
		C.VkCommandBuffer(commandBuffer),
		C.VkBuffer(srcBuffer),
		C.VkImage(dstImage),
		C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)),
		&cRegions[0],
	)
}

// CmdCopyImageToBuffer copies data from an image to a buffer
func CmdCopyImageToBuffer(
	commandBuffer CommandBuffer,
	srcImage Image,
	srcImageLayout ImageLayout,
	dstBuffer Buffer,
	regions []BufferImageCopy,
) {
	if commandBuffer == nil || srcImage == nil || dstBuffer == nil || len(regions) == 0 {
		return
	}

	cRegions := make([]C.VkBufferImageCopy, len(regions))
	for i, region := range regions {
		cRegions[i].bufferOffset = C.VkDeviceSize(region.BufferOffset)
		cRegions[i].bufferRowLength = C.uint32_t(region.BufferRowLength)
		cRegions[i].bufferImageHeight = C.uint32_t(region.BufferImageHeight)
		cRegions[i].imageSubresource.aspectMask = C.VkImageAspectFlags(region.ImageSubresource.AspectMask)
		cRegions[i].imageSubresource.mipLevel = C.uint32_t(region.ImageSubresource.MipLevel)
		cRegions[i].imageSubresource.baseArrayLayer = C.uint32_t(region.ImageSubresource.BaseArrayLayer)
		cRegions[i].imageSubresource.layerCount = C.uint32_t(region.ImageSubresource.LayerCount)
		cRegions[i].imageOffset.x = C.int32_t(region.ImageOffset.X)
		cRegions[i].imageOffset.y = C.int32_t(region.ImageOffset.Y)
		cRegions[i].imageOffset.z = C.int32_t(region.ImageOffset.Z)
		cRegions[i].imageExtent.width = C.uint32_t(region.ImageExtent.Width)
		cRegions[i].imageExtent.height = C.uint32_t(region.ImageExtent.Height)
		cRegions[i].imageExtent.depth = C.uint32_t(region.ImageExtent.Depth)
	}

	C.vkCmdCopyImageToBuffer(
		C.VkCommandBuffer(commandBuffer),
		C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout),
		C.VkBuffer(dstBuffer),
		C.uint32_t(len(cRegions)),
		&cRegions[0],
	)
}

// ImageCopy describes an image to image copy operation
type ImageCopy struct {
	SrcSubresource ImageSubresourceLayers
	SrcOffset      Offset3D
	DstSubresource ImageSubresourceLayers
	DstOffset      Offset3D
	Extent         Extent3D
}

// CmdCopyImage copies data between images
func CmdCopyImage(
	commandBuffer CommandBuffer,
	srcImage Image,
	srcImageLayout ImageLayout,
	dstImage Image,
	dstImageLayout ImageLayout,
	regions []ImageCopy,
) {
	if commandBuffer == nil || srcImage == nil || dstImage == nil || len(regions) == 0 {
		return
	}

	cRegions := make([]C.VkImageCopy, len(regions))
	for i, region := range regions {
		cRegions[i].srcSubresource.aspectMask = C.VkImageAspectFlags(region.SrcSubresource.AspectMask)
		cRegions[i].srcSubresource.mipLevel = C.uint32_t(region.SrcSubresource.MipLevel)
		cRegions[i].srcSubresource.baseArrayLayer = C.uint32_t(region.SrcSubresource.BaseArrayLayer)
		cRegions[i].srcSubresource.layerCount = C.uint32_t(region.SrcSubresource.LayerCount)
		cRegions[i].srcOffset.x = C.int32_t(region.SrcOffset.X)
		cRegions[i].srcOffset.y = C.int32_t(region.SrcOffset.Y)
		cRegions[i].srcOffset.z = C.int32_t(region.SrcOffset.Z)

		cRegions[i].dstSubresource.aspectMask = C.VkImageAspectFlags(region.DstSubresource.AspectMask)
		cRegions[i].dstSubresource.mipLevel = C.uint32_t(region.DstSubresource.MipLevel)
		cRegions[i].dstSubresource.baseArrayLayer = C.uint32_t(region.DstSubresource.BaseArrayLayer)
		cRegions[i].dstSubresource.layerCount = C.uint32_t(region.DstSubresource.LayerCount)
		cRegions[i].dstOffset.x = C.int32_t(region.DstOffset.X)
		cRegions[i].dstOffset.y = C.int32_t(region.DstOffset.Y)
		cRegions[i].dstOffset.z = C.int32_t(region.DstOffset.Z)

		cRegions[i].extent.width = C.uint32_t(region.Extent.Width)
		cRegions[i].extent.height = C.uint32_t(region.Extent.Height)
		cRegions[i].extent.depth = C.uint32_t(region.Extent.Depth)
	}

	C.vkCmdCopyImage(
		C.VkCommandBuffer(commandBuffer),
		C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout),
		C.VkImage(dstImage),
		C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)),
		&cRegions[0],
	)
}

// ============================================================================
// Buffer Operations
// ============================================================================

// CmdFillBuffer fills a buffer with a fixed 32-bit value
// size must be a multiple of 4, or WholeSize to fill to the end
func CmdFillBuffer(
	commandBuffer CommandBuffer,
	dstBuffer Buffer,
	dstOffset DeviceSize,
	size DeviceSize,
	data uint32,
) {
	if commandBuffer == nil || dstBuffer == nil {
		return
	}

	C.vkCmdFillBuffer(
		C.VkCommandBuffer(commandBuffer),
		C.VkBuffer(dstBuffer),
		C.VkDeviceSize(dstOffset),
		C.VkDeviceSize(size),
		C.uint32_t(data),
	)
}

// CmdUpdateBuffer updates buffer contents inline from host memory
// The data size must be less than or equal to 65536 bytes and a multiple of 4
func CmdUpdateBuffer(
	commandBuffer CommandBuffer,
	dstBuffer Buffer,
	dstOffset DeviceSize,
	data []byte,
) {
	if commandBuffer == nil || dstBuffer == nil || len(data) == 0 {
		return
	}

	// CmdUpdateBuffer is limited to 65536 bytes
	const maxUpdateSize = 65536
	if len(data) > maxUpdateSize {
		return
	}

	C.vkCmdUpdateBuffer(
		C.VkCommandBuffer(commandBuffer),
		C.VkBuffer(dstBuffer),
		C.VkDeviceSize(dstOffset),
		C.VkDeviceSize(len(data)),
		unsafe.Pointer(&data[0]),
	)
}

// ============================================================================
// Sparse Memory Binding Support
// ============================================================================

// SparseMemoryBind specifies a sparse memory bind operation
type SparseMemoryBind struct {
	ResourceOffset DeviceSize
	Size           DeviceSize
	Memory         DeviceMemory
	MemoryOffset   DeviceSize
	Flags          SparseMemoryBindFlags
}

// SparseMemoryBindFlags represents sparse memory bind flags
type SparseMemoryBindFlags uint32

const (
	SparseMemoryBindMetadataBit SparseMemoryBindFlags = C.VK_SPARSE_MEMORY_BIND_METADATA_BIT
)

// SparseBufferMemoryBindInfo specifies sparse buffer memory binding info
type SparseBufferMemoryBindInfo struct {
	Buffer Buffer
	Binds  []SparseMemoryBind
}

// SparseImageOpaqueMemoryBindInfo specifies sparse image opaque memory binding info
type SparseImageOpaqueMemoryBindInfo struct {
	Image Image
	Binds []SparseMemoryBind
}

// SparseImageMemoryBind specifies a sparse image memory bind
type SparseImageMemoryBind struct {
	Subresource  ImageSubresource
	Offset       Offset3D
	Extent       Extent3D
	Memory       DeviceMemory
	MemoryOffset DeviceSize
	Flags        SparseMemoryBindFlags
}

// ImageSubresource represents an image subresource
type ImageSubresource struct {
	AspectMask ImageAspectFlags
	MipLevel   uint32
	ArrayLayer uint32
}

// SparseImageMemoryBindInfo specifies sparse image memory binding info
type SparseImageMemoryBindInfo struct {
	Image Image
	Binds []SparseImageMemoryBind
}

// BindSparseInfo describes a sparse binding operation
type BindSparseInfo struct {
	WaitSemaphores       []Semaphore
	BufferBinds          []SparseBufferMemoryBindInfo
	ImageOpaqueBinds     []SparseImageOpaqueMemoryBindInfo
	ImageBinds           []SparseImageMemoryBindInfo
	SignalSemaphores     []Semaphore
}

// QueueBindSparse binds sparse resources on a queue
func QueueBindSparse(queue Queue, bindInfos []BindSparseInfo, fence Fence) error {
	if queue == nil {
		return NewValidationError("queue", "cannot be nil")
	}

	if len(bindInfos) == 0 {
		// Nothing to bind, still valid
		cFence := C.VkFence(nil)
		if fence != nil {
			cFence = C.VkFence(fence)
		}
		result := Result(C.vkQueueBindSparse(C.VkQueue(queue), 0, nil, cFence))
		if result != Success {
			return NewVulkanError(result, "QueueBindSparse", "Vulkan queue bind sparse failed")
		}
		return nil
	}

	// Convert bind infos
	cBindInfos := make([]C.VkBindSparseInfo, len(bindInfos))

	// We need to keep these slices alive for the duration of the call
	type bindInfoArrays struct {
		waitSemaphores              []C.VkSemaphore
		signalSemaphores            []C.VkSemaphore
		bufferBindInfos             []C.VkSparseBufferMemoryBindInfo
		imageOpaqueBindInfos        []C.VkSparseImageOpaqueMemoryBindInfo
		imageBindInfos              []C.VkSparseImageMemoryBindInfo
		sparseMemoryBinds           [][]C.VkSparseMemoryBind
		sparseImageOpaqueMemoryBinds [][]C.VkSparseMemoryBind
		sparseImageMemoryBinds      [][]C.VkSparseImageMemoryBind
	}
	arrays := make([]bindInfoArrays, len(bindInfos))

	for i, info := range bindInfos {
		cBindInfos[i].sType = C.VK_STRUCTURE_TYPE_BIND_SPARSE_INFO
		cBindInfos[i].pNext = nil

		// Wait semaphores
		if len(info.WaitSemaphores) > 0 {
			arrays[i].waitSemaphores = make([]C.VkSemaphore, len(info.WaitSemaphores))
			for j, sem := range info.WaitSemaphores {
				arrays[i].waitSemaphores[j] = C.VkSemaphore(sem)
			}
			cBindInfos[i].waitSemaphoreCount = C.uint32_t(len(arrays[i].waitSemaphores))
			cBindInfos[i].pWaitSemaphores = &arrays[i].waitSemaphores[0]
		}

		// Buffer binds
		if len(info.BufferBinds) > 0 {
			arrays[i].bufferBindInfos = make([]C.VkSparseBufferMemoryBindInfo, len(info.BufferBinds))
			arrays[i].sparseMemoryBinds = make([][]C.VkSparseMemoryBind, len(info.BufferBinds))
			for j, bufBind := range info.BufferBinds {
				arrays[i].bufferBindInfos[j].buffer = C.VkBuffer(bufBind.Buffer)
				if len(bufBind.Binds) > 0 {
					arrays[i].sparseMemoryBinds[j] = make([]C.VkSparseMemoryBind, len(bufBind.Binds))
					for k, bind := range bufBind.Binds {
						arrays[i].sparseMemoryBinds[j][k].resourceOffset = C.VkDeviceSize(bind.ResourceOffset)
						arrays[i].sparseMemoryBinds[j][k].size = C.VkDeviceSize(bind.Size)
						arrays[i].sparseMemoryBinds[j][k].memory = C.VkDeviceMemory(bind.Memory)
						arrays[i].sparseMemoryBinds[j][k].memoryOffset = C.VkDeviceSize(bind.MemoryOffset)
						arrays[i].sparseMemoryBinds[j][k].flags = C.VkSparseMemoryBindFlags(bind.Flags)
					}
					arrays[i].bufferBindInfos[j].bindCount = C.uint32_t(len(arrays[i].sparseMemoryBinds[j]))
					arrays[i].bufferBindInfos[j].pBinds = &arrays[i].sparseMemoryBinds[j][0]
				}
			}
			cBindInfos[i].bufferBindCount = C.uint32_t(len(arrays[i].bufferBindInfos))
			cBindInfos[i].pBufferBinds = &arrays[i].bufferBindInfos[0]
		}

		// Image opaque binds
		if len(info.ImageOpaqueBinds) > 0 {
			arrays[i].imageOpaqueBindInfos = make([]C.VkSparseImageOpaqueMemoryBindInfo, len(info.ImageOpaqueBinds))
			arrays[i].sparseImageOpaqueMemoryBinds = make([][]C.VkSparseMemoryBind, len(info.ImageOpaqueBinds))
			for j, imgBind := range info.ImageOpaqueBinds {
				arrays[i].imageOpaqueBindInfos[j].image = C.VkImage(imgBind.Image)
				if len(imgBind.Binds) > 0 {
					arrays[i].sparseImageOpaqueMemoryBinds[j] = make([]C.VkSparseMemoryBind, len(imgBind.Binds))
					for k, bind := range imgBind.Binds {
						arrays[i].sparseImageOpaqueMemoryBinds[j][k].resourceOffset = C.VkDeviceSize(bind.ResourceOffset)
						arrays[i].sparseImageOpaqueMemoryBinds[j][k].size = C.VkDeviceSize(bind.Size)
						arrays[i].sparseImageOpaqueMemoryBinds[j][k].memory = C.VkDeviceMemory(bind.Memory)
						arrays[i].sparseImageOpaqueMemoryBinds[j][k].memoryOffset = C.VkDeviceSize(bind.MemoryOffset)
						arrays[i].sparseImageOpaqueMemoryBinds[j][k].flags = C.VkSparseMemoryBindFlags(bind.Flags)
					}
					arrays[i].imageOpaqueBindInfos[j].bindCount = C.uint32_t(len(arrays[i].sparseImageOpaqueMemoryBinds[j]))
					arrays[i].imageOpaqueBindInfos[j].pBinds = &arrays[i].sparseImageOpaqueMemoryBinds[j][0]
				}
			}
			cBindInfos[i].imageOpaqueBindCount = C.uint32_t(len(arrays[i].imageOpaqueBindInfos))
			cBindInfos[i].pImageOpaqueBinds = &arrays[i].imageOpaqueBindInfos[0]
		}

		// Image binds
		if len(info.ImageBinds) > 0 {
			arrays[i].imageBindInfos = make([]C.VkSparseImageMemoryBindInfo, len(info.ImageBinds))
			arrays[i].sparseImageMemoryBinds = make([][]C.VkSparseImageMemoryBind, len(info.ImageBinds))
			for j, imgBind := range info.ImageBinds {
				arrays[i].imageBindInfos[j].image = C.VkImage(imgBind.Image)
				if len(imgBind.Binds) > 0 {
					arrays[i].sparseImageMemoryBinds[j] = make([]C.VkSparseImageMemoryBind, len(imgBind.Binds))
					for k, bind := range imgBind.Binds {
						arrays[i].sparseImageMemoryBinds[j][k].subresource.aspectMask = C.VkImageAspectFlags(bind.Subresource.AspectMask)
						arrays[i].sparseImageMemoryBinds[j][k].subresource.mipLevel = C.uint32_t(bind.Subresource.MipLevel)
						arrays[i].sparseImageMemoryBinds[j][k].subresource.arrayLayer = C.uint32_t(bind.Subresource.ArrayLayer)
						arrays[i].sparseImageMemoryBinds[j][k].offset.x = C.int32_t(bind.Offset.X)
						arrays[i].sparseImageMemoryBinds[j][k].offset.y = C.int32_t(bind.Offset.Y)
						arrays[i].sparseImageMemoryBinds[j][k].offset.z = C.int32_t(bind.Offset.Z)
						arrays[i].sparseImageMemoryBinds[j][k].extent.width = C.uint32_t(bind.Extent.Width)
						arrays[i].sparseImageMemoryBinds[j][k].extent.height = C.uint32_t(bind.Extent.Height)
						arrays[i].sparseImageMemoryBinds[j][k].extent.depth = C.uint32_t(bind.Extent.Depth)
						arrays[i].sparseImageMemoryBinds[j][k].memory = C.VkDeviceMemory(bind.Memory)
						arrays[i].sparseImageMemoryBinds[j][k].memoryOffset = C.VkDeviceSize(bind.MemoryOffset)
						arrays[i].sparseImageMemoryBinds[j][k].flags = C.VkSparseMemoryBindFlags(bind.Flags)
					}
					arrays[i].imageBindInfos[j].bindCount = C.uint32_t(len(arrays[i].sparseImageMemoryBinds[j]))
					arrays[i].imageBindInfos[j].pBinds = &arrays[i].sparseImageMemoryBinds[j][0]
				}
			}
			cBindInfos[i].imageBindCount = C.uint32_t(len(arrays[i].imageBindInfos))
			cBindInfos[i].pImageBinds = &arrays[i].imageBindInfos[0]
		}

		// Signal semaphores
		if len(info.SignalSemaphores) > 0 {
			arrays[i].signalSemaphores = make([]C.VkSemaphore, len(info.SignalSemaphores))
			for j, sem := range info.SignalSemaphores {
				arrays[i].signalSemaphores[j] = C.VkSemaphore(sem)
			}
			cBindInfos[i].signalSemaphoreCount = C.uint32_t(len(arrays[i].signalSemaphores))
			cBindInfos[i].pSignalSemaphores = &arrays[i].signalSemaphores[0]
		}
	}

	cFence := C.VkFence(nil)
	if fence != nil {
		cFence = C.VkFence(fence)
	}

	result := Result(C.vkQueueBindSparse(C.VkQueue(queue), C.uint32_t(len(cBindInfos)), &cBindInfos[0], cFence))
	if result != Success {
		return NewVulkanError(result, "QueueBindSparse", "Vulkan queue bind sparse failed")
	}

	return nil
}
