package vulkan

/*
#include <vulkan/vulkan.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

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

	barrier.SrcAccessMask, barrier.DstAccessMask, srcStage, dstStage = getLayoutTransitionAccessAndStages(oldLayout, newLayout)

	CmdPipelineBarrierFull(commandBuffer, srcStage, dstStage, 0, nil, nil, []ImageMemoryBarrier{barrier})
}

// getLayoutTransitionAccessAndStages returns the source access mask, destination access mask,
// source pipeline stage, and destination pipeline stage for a given image layout transition.
func getLayoutTransitionAccessAndStages(oldLayout, newLayout ImageLayout) (AccessFlags, AccessFlags, PipelineStageFlags, PipelineStageFlags) {
	var srcAccessMask AccessFlags
	var dstAccessMask AccessFlags
	var srcStage PipelineStageFlags
	var dstStage PipelineStageFlags

	// Determine access masks and pipeline stages based on the transition
	switch {
	case oldLayout == ImageLayoutUndefined && newLayout == ImageLayoutTransferDstOptimal:
		srcAccessMask = 0
		dstAccessMask = AccessTransferWriteBit
		srcStage = PipelineStageTopOfPipeBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutTransferDstOptimal && newLayout == ImageLayoutShaderReadOnlyOptimal:
		srcAccessMask = AccessTransferWriteBit
		dstAccessMask = AccessShaderReadBit
		srcStage = PipelineStageTransferBit
		dstStage = PipelineStageFragmentShaderBit

	case oldLayout == ImageLayoutUndefined && newLayout == ImageLayoutDepthStencilAttachmentOptimal:
		srcAccessMask = 0
		dstAccessMask = AccessDepthStencilAttachmentReadBit | AccessDepthStencilAttachmentWriteBit
		srcStage = PipelineStageTopOfPipeBit
		dstStage = PipelineStageEarlyFragmentTestsBit

	case oldLayout == ImageLayoutUndefined && newLayout == ImageLayoutColorAttachmentOptimal:
		srcAccessMask = 0
		dstAccessMask = AccessColorAttachmentReadBit | AccessColorAttachmentWriteBit
		srcStage = PipelineStageTopOfPipeBit
		dstStage = PipelineStageColorAttachmentOutputBit

	case oldLayout == ImageLayoutColorAttachmentOptimal && newLayout == ImageLayoutPresentSrcKHR:
		srcAccessMask = AccessColorAttachmentWriteBit
		dstAccessMask = 0
		srcStage = PipelineStageColorAttachmentOutputBit
		dstStage = PipelineStageBottomOfPipeBit

	case oldLayout == ImageLayoutTransferSrcOptimal && newLayout == ImageLayoutShaderReadOnlyOptimal:
		srcAccessMask = AccessTransferReadBit
		dstAccessMask = AccessShaderReadBit
		srcStage = PipelineStageTransferBit
		dstStage = PipelineStageFragmentShaderBit

	case oldLayout == ImageLayoutShaderReadOnlyOptimal && newLayout == ImageLayoutTransferSrcOptimal:
		srcAccessMask = AccessShaderReadBit
		dstAccessMask = AccessTransferReadBit
		srcStage = PipelineStageFragmentShaderBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutShaderReadOnlyOptimal && newLayout == ImageLayoutTransferDstOptimal:
		srcAccessMask = AccessShaderReadBit
		dstAccessMask = AccessTransferWriteBit
		srcStage = PipelineStageFragmentShaderBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutGeneral && newLayout == ImageLayoutTransferSrcOptimal:
		srcAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		dstAccessMask = AccessTransferReadBit
		srcStage = PipelineStageAllCommandsBit
		dstStage = PipelineStageTransferBit

	case oldLayout == ImageLayoutGeneral && newLayout == ImageLayoutTransferDstOptimal:
		srcAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		dstAccessMask = AccessTransferWriteBit
		srcStage = PipelineStageAllCommandsBit
		dstStage = PipelineStageTransferBit

	default:
		// Generic fallback - use ALL_COMMANDS for safety
		srcAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		dstAccessMask = AccessMemoryReadBit | AccessMemoryWriteBit
		srcStage = PipelineStageAllCommandsBit
		dstStage = PipelineStageAllCommandsBit
	}

	return srcAccessMask, dstAccessMask, srcStage, dstStage
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

// setImageSubresourceLayers sets the image subresource layers on a C struct
func setImageSubresourceLayers(dst *C.VkImageSubresourceLayers, src ImageSubresourceLayers) {
	dst.aspectMask = C.VkImageAspectFlags(src.AspectMask)
	dst.mipLevel = C.uint32_t(src.MipLevel)
	dst.baseArrayLayer = C.uint32_t(src.BaseArrayLayer)
	dst.layerCount = C.uint32_t(src.LayerCount)
}

// setOffset3D sets the offset on a C struct
func setOffset3D(dst *C.VkOffset3D, src Offset3D) {
	dst.x = C.int32_t(src.X)
	dst.y = C.int32_t(src.Y)
	dst.z = C.int32_t(src.Z)
}

// setExtent3D sets the extent on a C struct
func setExtent3D(dst *C.VkExtent3D, src Extent3D) {
	dst.width = C.uint32_t(src.Width)
	dst.height = C.uint32_t(src.Height)
	dst.depth = C.uint32_t(src.Depth)
}

// buildImageResolveRegions converts Go ImageResolve to C VkImageResolve
func buildImageResolveRegions(regions []ImageResolve) []C.VkImageResolve {
	cRegions := make([]C.VkImageResolve, len(regions))
	for i, region := range regions {
		setImageSubresourceLayers(&cRegions[i].srcSubresource, region.SrcSubresource)
		setOffset3D(&cRegions[i].srcOffset, region.SrcOffset)
		setImageSubresourceLayers(&cRegions[i].dstSubresource, region.DstSubresource)
		setOffset3D(&cRegions[i].dstOffset, region.DstOffset)
		setExtent3D(&cRegions[i].extent, region.Extent)
	}
	return cRegions
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
	cRegions := buildImageResolveRegions(regions)
	C.vkCmdResolveImage(C.VkCommandBuffer(commandBuffer), C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout), C.VkImage(dstImage), C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)), &cRegions[0])
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

// buildBufferImageCopyRegions converts Go BufferImageCopy to C VkBufferImageCopy
func buildBufferImageCopyRegions(regions []BufferImageCopy) []C.VkBufferImageCopy {
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
	return cRegions
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
	cRegions := buildBufferImageCopyRegions(regions)
	C.vkCmdCopyBufferToImage(C.VkCommandBuffer(commandBuffer), C.VkBuffer(srcBuffer),
		C.VkImage(dstImage), C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)), &cRegions[0])
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
	cRegions := buildBufferImageCopyRegions(regions)
	C.vkCmdCopyImageToBuffer(C.VkCommandBuffer(commandBuffer), C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout), C.VkBuffer(dstBuffer),
		C.uint32_t(len(cRegions)), &cRegions[0])
}

// ImageCopy describes an image to image copy operation (same structure as ImageResolve)
type ImageCopy = ImageResolve

// buildImageCopyRegions converts Go ImageCopy to C VkImageCopy
func buildImageCopyRegions(regions []ImageCopy) []C.VkImageCopy {
	cRegions := make([]C.VkImageCopy, len(regions))
	for i, region := range regions {
		setImageSubresourceLayers(&cRegions[i].srcSubresource, region.SrcSubresource)
		setOffset3D(&cRegions[i].srcOffset, region.SrcOffset)
		setImageSubresourceLayers(&cRegions[i].dstSubresource, region.DstSubresource)
		setOffset3D(&cRegions[i].dstOffset, region.DstOffset)
		setExtent3D(&cRegions[i].extent, region.Extent)
	}
	return cRegions
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
	cRegions := buildImageCopyRegions(regions)
	C.vkCmdCopyImage(C.VkCommandBuffer(commandBuffer), C.VkImage(srcImage),
		C.VkImageLayout(srcImageLayout), C.VkImage(dstImage), C.VkImageLayout(dstImageLayout),
		C.uint32_t(len(cRegions)), &cRegions[0])
}

// ============================================================================
// Buffer Operations
// ============================================================================

// CmdFillBuffer executes the operation
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

// CmdUpdateBuffer executes the operation
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

// buildSparseMemoryBinds converts Go SparseMemoryBind slice to C VkSparseMemoryBind slice
func buildSparseMemoryBinds(binds []SparseMemoryBind) []C.VkSparseMemoryBind {
	cBinds := make([]C.VkSparseMemoryBind, len(binds))
	for k, bind := range binds {
		cBinds[k].resourceOffset = C.VkDeviceSize(bind.ResourceOffset)
		cBinds[k].size = C.VkDeviceSize(bind.Size)
		cBinds[k].memory = C.VkDeviceMemory(bind.Memory)
		cBinds[k].memoryOffset = C.VkDeviceSize(bind.MemoryOffset)
		cBinds[k].flags = C.VkSparseMemoryBindFlags(bind.Flags)
	}
	return cBinds
}

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

// SubresourceLayout represents an image subresource layout
type SubresourceLayout struct {
	Offset     DeviceSize
	Size       DeviceSize
	RowPitch   DeviceSize
	ArrayPitch DeviceSize
	DepthPitch DeviceSize
}

// GetImageSubresourceLayout queries the layout of an image subresource
func GetImageSubresourceLayout(device Device, image Image, subresource *ImageSubresource) SubresourceLayout {
	if device == nil || image == nil || subresource == nil {
		return SubresourceLayout{}
	}

	var cSubresource C.VkImageSubresource
	cSubresource.aspectMask = C.VkImageAspectFlags(subresource.AspectMask)
	cSubresource.mipLevel = C.uint32_t(subresource.MipLevel)
	cSubresource.arrayLayer = C.uint32_t(subresource.ArrayLayer)

	var cLayout C.VkSubresourceLayout
	C.vkGetImageSubresourceLayout(C.VkDevice(device), C.VkImage(image), &cSubresource, &cLayout)

	return SubresourceLayout{
		Offset:     DeviceSize(cLayout.offset),
		Size:       DeviceSize(cLayout.size),
		RowPitch:   DeviceSize(cLayout.rowPitch),
		ArrayPitch: DeviceSize(cLayout.arrayPitch),
		DepthPitch: DeviceSize(cLayout.depthPitch),
	}
}

// SparseImageMemoryBindInfo specifies sparse image memory binding info
type SparseImageMemoryBindInfo struct {
	Image Image
	Binds []SparseImageMemoryBind
}

// BindSparseInfo describes a sparse binding operation
type BindSparseInfo struct {
	WaitSemaphores   []Semaphore
	BufferBinds      []SparseBufferMemoryBindInfo
	ImageOpaqueBinds []SparseImageOpaqueMemoryBindInfo
	ImageBinds       []SparseImageMemoryBindInfo
	SignalSemaphores []Semaphore
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

	// All nested arrays referenced from cBindInfos must be pinned: cBindInfos
	// is Go memory passed to C and may not contain unpinned Go pointers.
	var pinner runtime.Pinner
	defer pinner.Unpin()

	// We need to keep these slices alive for the duration of the call
	type bindInfoArrays struct {
		waitSemaphores               []C.VkSemaphore
		signalSemaphores             []C.VkSemaphore
		bufferBindInfos              []C.VkSparseBufferMemoryBindInfo
		imageOpaqueBindInfos         []C.VkSparseImageOpaqueMemoryBindInfo
		imageBindInfos               []C.VkSparseImageMemoryBindInfo
		sparseMemoryBinds            [][]C.VkSparseMemoryBind
		sparseImageOpaqueMemoryBinds [][]C.VkSparseMemoryBind
		sparseImageMemoryBinds       [][]C.VkSparseImageMemoryBind
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
			pinner.Pin(&arrays[i].waitSemaphores[0])
			cBindInfos[i].waitSemaphoreCount = C.uint32_t(len(arrays[i].waitSemaphores))
			cBindInfos[i].pWaitSemaphores = &arrays[i].waitSemaphores[0]
		}

		// Buffer binds - populate sparse buffer memory bind info
		if bufBindCount := len(info.BufferBinds); bufBindCount > 0 {
			arrays[i].bufferBindInfos = make([]C.VkSparseBufferMemoryBindInfo, bufBindCount)
			arrays[i].sparseMemoryBinds = make([][]C.VkSparseMemoryBind, bufBindCount)
			for j, bufBind := range info.BufferBinds {
				arrays[i].bufferBindInfos[j].buffer = C.VkBuffer(bufBind.Buffer)
				if len(bufBind.Binds) > 0 {
					arrays[i].sparseMemoryBinds[j] = buildSparseMemoryBinds(bufBind.Binds)
					pinner.Pin(&arrays[i].sparseMemoryBinds[j][0])
					arrays[i].bufferBindInfos[j].bindCount = C.uint32_t(len(arrays[i].sparseMemoryBinds[j]))
					arrays[i].bufferBindInfos[j].pBinds = &arrays[i].sparseMemoryBinds[j][0]
				}
			}
			pinner.Pin(&arrays[i].bufferBindInfos[0])
			cBindInfos[i].bufferBindCount = C.uint32_t(len(arrays[i].bufferBindInfos))
			cBindInfos[i].pBufferBinds = &arrays[i].bufferBindInfos[0]
		}

		// Image opaque binds - populate sparse image opaque memory bind info
		if imgOpaqueBindCount := len(info.ImageOpaqueBinds); imgOpaqueBindCount > 0 {
			arrays[i].imageOpaqueBindInfos = make([]C.VkSparseImageOpaqueMemoryBindInfo, imgOpaqueBindCount)
			arrays[i].sparseImageOpaqueMemoryBinds = make([][]C.VkSparseMemoryBind, imgOpaqueBindCount)
			for j := 0; j < imgOpaqueBindCount; j++ {
				imgBind := info.ImageOpaqueBinds[j]
				arrays[i].imageOpaqueBindInfos[j].image = C.VkImage(imgBind.Image)
				if bindLen := len(imgBind.Binds); bindLen > 0 {
					arrays[i].sparseImageOpaqueMemoryBinds[j] = buildSparseMemoryBinds(imgBind.Binds)
					pinner.Pin(&arrays[i].sparseImageOpaqueMemoryBinds[j][0])
					arrays[i].imageOpaqueBindInfos[j].bindCount = C.uint32_t(bindLen)
					arrays[i].imageOpaqueBindInfos[j].pBinds = &arrays[i].sparseImageOpaqueMemoryBinds[j][0]
				}
			}
			pinner.Pin(&arrays[i].imageOpaqueBindInfos[0])
			cBindInfos[i].imageOpaqueBindCount = C.uint32_t(imgOpaqueBindCount)
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
					pinner.Pin(&arrays[i].sparseImageMemoryBinds[j][0])
					arrays[i].imageBindInfos[j].bindCount = C.uint32_t(len(arrays[i].sparseImageMemoryBinds[j]))
					arrays[i].imageBindInfos[j].pBinds = &arrays[i].sparseImageMemoryBinds[j][0]
				}
			}
			pinner.Pin(&arrays[i].imageBindInfos[0])
			cBindInfos[i].imageBindCount = C.uint32_t(len(arrays[i].imageBindInfos))
			cBindInfos[i].pImageBinds = &arrays[i].imageBindInfos[0]
		}

		// Signal semaphores
		if len(info.SignalSemaphores) > 0 {
			arrays[i].signalSemaphores = make([]C.VkSemaphore, len(info.SignalSemaphores))
			for j, sem := range info.SignalSemaphores {
				arrays[i].signalSemaphores[j] = C.VkSemaphore(sem)
			}
			pinner.Pin(&arrays[i].signalSemaphores[0])
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
