package vulkan

/*
#include <vulkan/vulkan.h>
*/
import "C"
import "unsafe"

// ============================================================================
// Timeline Semaphore Support (Vulkan 1.2+)
// ============================================================================

// SemaphoreType represents semaphore types
type SemaphoreType uint32

const (
	SemaphoreTypeBinary   SemaphoreType = C.VK_SEMAPHORE_TYPE_BINARY
	SemaphoreTypeTimeline SemaphoreType = C.VK_SEMAPHORE_TYPE_TIMELINE
)

// SemaphoreTypeCreateInfo specifies the type of a semaphore
type SemaphoreTypeCreateInfo struct {
	SemaphoreType SemaphoreType
	InitialValue  uint64
}

// CreateTimelineSemaphore creates a timeline semaphore (Vulkan 1.2+)
func CreateTimelineSemaphore(device Device, initialValue uint64) (Semaphore, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}

	var cTypeCreateInfo C.VkSemaphoreTypeCreateInfo
	cTypeCreateInfo.sType = C.VK_STRUCTURE_TYPE_SEMAPHORE_TYPE_CREATE_INFO
	cTypeCreateInfo.pNext = nil
	cTypeCreateInfo.semaphoreType = C.VK_SEMAPHORE_TYPE_TIMELINE
	cTypeCreateInfo.initialValue = C.uint64_t(initialValue)

	var cCreateInfo C.VkSemaphoreCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO
	cCreateInfo.pNext = unsafe.Pointer(&cTypeCreateInfo)
	cCreateInfo.flags = 0

	var semaphore C.VkSemaphore
	result := Result(C.vkCreateSemaphore(C.VkDevice(device), &cCreateInfo, nil, &semaphore))
	if result != Success {
		return nil, NewVulkanError(result, "CreateTimelineSemaphore", "Vulkan timeline semaphore creation failed")
	}

	return Semaphore(semaphore), nil
}

// SemaphoreWaitFlags represents semaphore wait flags
type SemaphoreWaitFlags uint32

const (
	SemaphoreWaitAnyBit SemaphoreWaitFlags = C.VK_SEMAPHORE_WAIT_ANY_BIT
)

// SemaphoreWaitInfo contains information for waiting on semaphores
type SemaphoreWaitInfo struct {
	Flags      SemaphoreWaitFlags
	Semaphores []Semaphore
	Values     []uint64
}

// WaitSemaphores waits for timeline semaphores (Vulkan 1.2+)
func WaitSemaphores(device Device, waitInfo *SemaphoreWaitInfo, timeout uint64) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if waitInfo == nil {
		return NewValidationError("waitInfo", "cannot be nil")
	}
	if len(waitInfo.Semaphores) == 0 {
		return NewValidationError("Semaphores", "cannot be empty")
	}
	if len(waitInfo.Values) != len(waitInfo.Semaphores) {
		return NewValidationError("Values", "must have same length as Semaphores")
	}

	cSemaphores := make([]C.VkSemaphore, len(waitInfo.Semaphores))
	for i, sem := range waitInfo.Semaphores {
		cSemaphores[i] = C.VkSemaphore(sem)
	}

	cValues := make([]C.uint64_t, len(waitInfo.Values))
	for i, val := range waitInfo.Values {
		cValues[i] = C.uint64_t(val)
	}

	var cWaitInfo C.VkSemaphoreWaitInfo
	cWaitInfo.sType = C.VK_STRUCTURE_TYPE_SEMAPHORE_WAIT_INFO
	cWaitInfo.pNext = nil
	cWaitInfo.flags = C.VkSemaphoreWaitFlags(waitInfo.Flags)
	cWaitInfo.semaphoreCount = C.uint32_t(len(cSemaphores))
	cWaitInfo.pSemaphores = &cSemaphores[0]
	cWaitInfo.pValues = &cValues[0]

	result := Result(C.vkWaitSemaphores(C.VkDevice(device), &cWaitInfo, C.uint64_t(timeout)))
	if result != Success {
		return NewVulkanError(result, "WaitSemaphores", "Vulkan semaphore wait failed")
	}

	return nil
}

// SemaphoreSignalInfo contains information for signaling a semaphore
type SemaphoreSignalInfo struct {
	Semaphore Semaphore
	Value     uint64
}

// SignalSemaphore signals a timeline semaphore (Vulkan 1.2+)
func SignalSemaphore(device Device, signalInfo *SemaphoreSignalInfo) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if signalInfo == nil {
		return NewValidationError("signalInfo", "cannot be nil")
	}
	if signalInfo.Semaphore == nil {
		return NewValidationError("Semaphore", "cannot be nil")
	}

	var cSignalInfo C.VkSemaphoreSignalInfo
	cSignalInfo.sType = C.VK_STRUCTURE_TYPE_SEMAPHORE_SIGNAL_INFO
	cSignalInfo.pNext = nil
	cSignalInfo.semaphore = C.VkSemaphore(signalInfo.Semaphore)
	cSignalInfo.value = C.uint64_t(signalInfo.Value)

	result := Result(C.vkSignalSemaphore(C.VkDevice(device), &cSignalInfo))
	if result != Success {
		return NewVulkanError(result, "SignalSemaphore", "Vulkan semaphore signal failed")
	}

	return nil
}

// GetSemaphoreCounterValue gets the current counter value of a timeline semaphore (Vulkan 1.2+)
func GetSemaphoreCounterValue(device Device, semaphore Semaphore) (uint64, error) {
	if device == nil {
		return 0, NewValidationError("device", "cannot be nil")
	}
	if semaphore == nil {
		return 0, NewValidationError("semaphore", "cannot be nil")
	}

	var value C.uint64_t
	result := Result(C.vkGetSemaphoreCounterValue(C.VkDevice(device), C.VkSemaphore(semaphore), &value))
	if result != Success {
		return 0, NewVulkanError(result, "GetSemaphoreCounterValue", "Vulkan semaphore counter value query failed")
	}

	return uint64(value), nil
}

// ============================================================================
// Event Objects
// ============================================================================

// EventCreateInfo contains event creation information
type EventCreateInfo struct {
	Flags EventCreateFlags
}

// EventCreateFlags represents event creation flags
type EventCreateFlags uint32

const (
	EventCreateDeviceOnlyBit EventCreateFlags = C.VK_EVENT_CREATE_DEVICE_ONLY_BIT
)

// CreateEvent creates an event object
func CreateEvent(device Device, createInfo *EventCreateInfo) (Event, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		createInfo = &EventCreateInfo{}
	}

	var cCreateInfo C.VkEventCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_EVENT_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkEventCreateFlags(createInfo.Flags)

	var event C.VkEvent
	result := Result(C.vkCreateEvent(C.VkDevice(device), &cCreateInfo, nil, &event))
	if result != Success {
		return nil, NewVulkanError(result, "CreateEvent", "Vulkan event creation failed")
	}

	return Event(event), nil
}

// DestroyEvent destroys an event object
func DestroyEvent(device Device, event Event) {
	if device == nil || event == nil {
		return
	}
	C.vkDestroyEvent(C.VkDevice(device), C.VkEvent(event), nil)
}

// SetEvent sets an event to signaled state from the host
func SetEvent(device Device, event Event) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if event == nil {
		return NewValidationError("event", "cannot be nil")
	}

	result := Result(C.vkSetEvent(C.VkDevice(device), C.VkEvent(event)))
	if result != Success {
		return NewVulkanError(result, "SetEvent", "Vulkan set event failed")
	}

	return nil
}

// ResetEvent resets an event to unsignaled state from the host
func ResetEvent(device Device, event Event) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if event == nil {
		return NewValidationError("event", "cannot be nil")
	}

	result := Result(C.vkResetEvent(C.VkDevice(device), C.VkEvent(event)))
	if result != Success {
		return NewVulkanError(result, "ResetEvent", "Vulkan reset event failed")
	}

	return nil
}

// GetEventStatus gets the status of an event
// Returns Success if the event is signaled, EventReset if unsignaled
func GetEventStatus(device Device, event Event) (Result, error) {
	if device == nil {
		return 0, NewValidationError("device", "cannot be nil")
	}
	if event == nil {
		return 0, NewValidationError("event", "cannot be nil")
	}

	result := Result(C.vkGetEventStatus(C.VkDevice(device), C.VkEvent(event)))
	if result != Success && result != EventReset {
		return result, NewVulkanError(result, "GetEventStatus", "Vulkan get event status failed")
	}

	return result, nil
}

// CmdSetEvent sets an event object to signaled state from the device
func CmdSetEvent(commandBuffer CommandBuffer, event Event, stageMask PipelineStageFlags) {
	if commandBuffer == nil || event == nil {
		return
	}
	C.vkCmdSetEvent(C.VkCommandBuffer(commandBuffer), C.VkEvent(event), C.VkPipelineStageFlags(stageMask))
}

// CmdResetEvent resets an event object to unsignaled state from the device
func CmdResetEvent(commandBuffer CommandBuffer, event Event, stageMask PipelineStageFlags) {
	if commandBuffer == nil || event == nil {
		return
	}
	C.vkCmdResetEvent(C.VkCommandBuffer(commandBuffer), C.VkEvent(event), C.VkPipelineStageFlags(stageMask))
}

// ============================================================================
// Memory Barriers with Queue Family Transfer
// ============================================================================

// MemoryBarrier represents a global memory barrier
type MemoryBarrier struct {
	SrcAccessMask AccessFlags
	DstAccessMask AccessFlags
}

// BufferMemoryBarrier represents a buffer memory barrier with queue family transfer support
type BufferMemoryBarrier struct {
	SrcAccessMask       AccessFlags
	DstAccessMask       AccessFlags
	SrcQueueFamilyIndex uint32
	DstQueueFamilyIndex uint32
	Buffer              Buffer
	Offset              uint64
	Size                uint64
}

// ImageMemoryBarrier represents an image memory barrier with queue family transfer support
type ImageMemoryBarrier struct {
	SrcAccessMask       AccessFlags
	DstAccessMask       AccessFlags
	OldLayout           ImageLayout
	NewLayout           ImageLayout
	SrcQueueFamilyIndex uint32
	DstQueueFamilyIndex uint32
	Image               Image
	SubresourceRange    ImageSubresourceRange
}

// buildImageMemoryBarriers converts Go ImageMemoryBarrier to C VkImageMemoryBarrier
func buildImageMemoryBarriers(imageMemoryBarriers []ImageMemoryBarrier) []C.VkImageMemoryBarrier {
	if len(imageMemoryBarriers) == 0 {
		return nil
	}
	cImageMemoryBarriers := make([]C.VkImageMemoryBarrier, len(imageMemoryBarriers))
	for i, imb := range imageMemoryBarriers {
		cImageMemoryBarriers[i].sType = C.VK_STRUCTURE_TYPE_IMAGE_MEMORY_BARRIER
		cImageMemoryBarriers[i].pNext = nil
		cImageMemoryBarriers[i].srcAccessMask = C.VkAccessFlags(imb.SrcAccessMask)
		cImageMemoryBarriers[i].dstAccessMask = C.VkAccessFlags(imb.DstAccessMask)
		cImageMemoryBarriers[i].oldLayout = C.VkImageLayout(imb.OldLayout)
		cImageMemoryBarriers[i].newLayout = C.VkImageLayout(imb.NewLayout)
		cImageMemoryBarriers[i].srcQueueFamilyIndex = C.uint32_t(imb.SrcQueueFamilyIndex)
		cImageMemoryBarriers[i].dstQueueFamilyIndex = C.uint32_t(imb.DstQueueFamilyIndex)
		cImageMemoryBarriers[i].image = C.VkImage(imb.Image)
		cImageMemoryBarriers[i].subresourceRange.aspectMask = C.VkImageAspectFlags(imb.SubresourceRange.AspectMask)
		cImageMemoryBarriers[i].subresourceRange.baseMipLevel = C.uint32_t(imb.SubresourceRange.BaseMipLevel)
		cImageMemoryBarriers[i].subresourceRange.levelCount = C.uint32_t(imb.SubresourceRange.LevelCount)
		cImageMemoryBarriers[i].subresourceRange.baseArrayLayer = C.uint32_t(imb.SubresourceRange.BaseArrayLayer)
		cImageMemoryBarriers[i].subresourceRange.layerCount = C.uint32_t(imb.SubresourceRange.LayerCount)
	}
	return cImageMemoryBarriers
}

// DependencyFlags represents dependency flags
type DependencyFlags uint32

const (
	DependencyByRegionBit    DependencyFlags = C.VK_DEPENDENCY_BY_REGION_BIT
	DependencyDeviceGroupBit DependencyFlags = C.VK_DEPENDENCY_DEVICE_GROUP_BIT
	DependencyViewLocalBit   DependencyFlags = C.VK_DEPENDENCY_VIEW_LOCAL_BIT
)

// CmdPipelineBarrierFull inserts a pipeline barrier with full memory barrier support
func CmdPipelineBarrierFull(
	commandBuffer CommandBuffer,
	srcStageMask PipelineStageFlags,
	dstStageMask PipelineStageFlags,
	dependencyFlags DependencyFlags,
	memoryBarriers []MemoryBarrier,
	bufferMemoryBarriers []BufferMemoryBarrier,
	imageMemoryBarriers []ImageMemoryBarrier,
) {
	if commandBuffer == nil {
		return
	}

	// Convert memory barriers
	var cMemoryBarriers []C.VkMemoryBarrier
	if len(memoryBarriers) > 0 {
		cMemoryBarriers = make([]C.VkMemoryBarrier, len(memoryBarriers))
		for i, mb := range memoryBarriers {
			cMemoryBarriers[i].sType = C.VK_STRUCTURE_TYPE_MEMORY_BARRIER
			cMemoryBarriers[i].pNext = nil
			cMemoryBarriers[i].srcAccessMask = C.VkAccessFlags(mb.SrcAccessMask)
			cMemoryBarriers[i].dstAccessMask = C.VkAccessFlags(mb.DstAccessMask)
		}
	}

	// Convert buffer memory barriers
	var cBufferMemoryBarriers []C.VkBufferMemoryBarrier
	if len(bufferMemoryBarriers) > 0 {
		cBufferMemoryBarriers = make([]C.VkBufferMemoryBarrier, len(bufferMemoryBarriers))
		for i, bmb := range bufferMemoryBarriers {
			cBufferMemoryBarriers[i].sType = C.VK_STRUCTURE_TYPE_BUFFER_MEMORY_BARRIER
			cBufferMemoryBarriers[i].pNext = nil
			cBufferMemoryBarriers[i].srcAccessMask = C.VkAccessFlags(bmb.SrcAccessMask)
			cBufferMemoryBarriers[i].dstAccessMask = C.VkAccessFlags(bmb.DstAccessMask)
			cBufferMemoryBarriers[i].srcQueueFamilyIndex = C.uint32_t(bmb.SrcQueueFamilyIndex)
			cBufferMemoryBarriers[i].dstQueueFamilyIndex = C.uint32_t(bmb.DstQueueFamilyIndex)
			cBufferMemoryBarriers[i].buffer = C.VkBuffer(bmb.Buffer)
			cBufferMemoryBarriers[i].offset = C.VkDeviceSize(bmb.Offset)
			cBufferMemoryBarriers[i].size = C.VkDeviceSize(bmb.Size)
		}
	}

	// Convert image memory barriers
	cImageMemoryBarriers := buildImageMemoryBarriers(imageMemoryBarriers)

	// Prepare pointers
	var pMemoryBarriers *C.VkMemoryBarrier
	var pBufferMemoryBarriers *C.VkBufferMemoryBarrier
	var pImageMemoryBarriers *C.VkImageMemoryBarrier

	if len(cMemoryBarriers) > 0 {
		pMemoryBarriers = &cMemoryBarriers[0]
	}
	if len(cBufferMemoryBarriers) > 0 {
		pBufferMemoryBarriers = &cBufferMemoryBarriers[0]
	}
	if len(cImageMemoryBarriers) > 0 {
		pImageMemoryBarriers = &cImageMemoryBarriers[0]
	}

	C.vkCmdPipelineBarrier(
		C.VkCommandBuffer(commandBuffer),
		C.VkPipelineStageFlags(srcStageMask),
		C.VkPipelineStageFlags(dstStageMask),
		C.VkDependencyFlags(dependencyFlags),
		C.uint32_t(len(cMemoryBarriers)),
		pMemoryBarriers,
		C.uint32_t(len(cBufferMemoryBarriers)),
		pBufferMemoryBarriers,
		C.uint32_t(len(cImageMemoryBarriers)),
		pImageMemoryBarriers,
	)
}

// CmdWaitEvents waits for one or more events and inserts a set of memory barriers
func CmdWaitEvents(
	commandBuffer CommandBuffer,
	events []Event,
	srcStageMask PipelineStageFlags,
	dstStageMask PipelineStageFlags,
	memoryBarriers []MemoryBarrier,
	bufferMemoryBarriers []BufferMemoryBarrier,
	imageMemoryBarriers []ImageMemoryBarrier,
) {
	if commandBuffer == nil || len(events) == 0 {
		return
	}

	// Convert events
	cEvents := make([]C.VkEvent, len(events))
	for i, e := range events {
		cEvents[i] = C.VkEvent(e)
	}

	// Convert memory barriers
	var cMemoryBarriers []C.VkMemoryBarrier
	if len(memoryBarriers) > 0 {
		cMemoryBarriers = make([]C.VkMemoryBarrier, len(memoryBarriers))
		for i, mb := range memoryBarriers {
			cMemoryBarriers[i].sType = C.VK_STRUCTURE_TYPE_MEMORY_BARRIER
			cMemoryBarriers[i].pNext = nil
			cMemoryBarriers[i].srcAccessMask = C.VkAccessFlags(mb.SrcAccessMask)
			cMemoryBarriers[i].dstAccessMask = C.VkAccessFlags(mb.DstAccessMask)
		}
	}

	// Convert buffer memory barriers
	var cBufferMemoryBarriers []C.VkBufferMemoryBarrier
	if len(bufferMemoryBarriers) > 0 {
		cBufferMemoryBarriers = make([]C.VkBufferMemoryBarrier, len(bufferMemoryBarriers))
		for i, bmb := range bufferMemoryBarriers {
			cBufferMemoryBarriers[i].sType = C.VK_STRUCTURE_TYPE_BUFFER_MEMORY_BARRIER
			cBufferMemoryBarriers[i].pNext = nil
			cBufferMemoryBarriers[i].srcAccessMask = C.VkAccessFlags(bmb.SrcAccessMask)
			cBufferMemoryBarriers[i].dstAccessMask = C.VkAccessFlags(bmb.DstAccessMask)
			cBufferMemoryBarriers[i].srcQueueFamilyIndex = C.uint32_t(bmb.SrcQueueFamilyIndex)
			cBufferMemoryBarriers[i].dstQueueFamilyIndex = C.uint32_t(bmb.DstQueueFamilyIndex)
			cBufferMemoryBarriers[i].buffer = C.VkBuffer(bmb.Buffer)
			cBufferMemoryBarriers[i].offset = C.VkDeviceSize(bmb.Offset)
			cBufferMemoryBarriers[i].size = C.VkDeviceSize(bmb.Size)
		}
	}

	// Convert image memory barriers
	cImageMemoryBarriers := buildImageMemoryBarriers(imageMemoryBarriers)

	// Prepare pointers
	var pMemoryBarriers *C.VkMemoryBarrier
	var pBufferMemoryBarriers *C.VkBufferMemoryBarrier
	var pImageMemoryBarriers *C.VkImageMemoryBarrier

	if len(cMemoryBarriers) > 0 {
		pMemoryBarriers = &cMemoryBarriers[0]
	}
	if len(cBufferMemoryBarriers) > 0 {
		pBufferMemoryBarriers = &cBufferMemoryBarriers[0]
	}
	if len(cImageMemoryBarriers) > 0 {
		pImageMemoryBarriers = &cImageMemoryBarriers[0]
	}

	C.vkCmdWaitEvents(
		C.VkCommandBuffer(commandBuffer),
		C.uint32_t(len(cEvents)),
		&cEvents[0],
		C.VkPipelineStageFlags(srcStageMask),
		C.VkPipelineStageFlags(dstStageMask),
		C.uint32_t(len(cMemoryBarriers)),
		pMemoryBarriers,
		C.uint32_t(len(cBufferMemoryBarriers)),
		pBufferMemoryBarriers,
		C.uint32_t(len(cImageMemoryBarriers)),
		pImageMemoryBarriers,
	)
}
