package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

// CommandPoolCreateInfo contains command pool creation information
type CommandPoolCreateInfo struct {
	Flags            CommandPoolCreateFlags
	QueueFamilyIndex uint32
}

// CommandPoolCreateFlags represents command pool creation flags
type CommandPoolCreateFlags uint32

const (
	CommandPoolCreateTransientBit            CommandPoolCreateFlags = C.VK_COMMAND_POOL_CREATE_TRANSIENT_BIT
	CommandPoolCreateResetCommandBufferBit   CommandPoolCreateFlags = C.VK_COMMAND_POOL_CREATE_RESET_COMMAND_BUFFER_BIT
	CommandPoolCreateProtectedBit            CommandPoolCreateFlags = C.VK_COMMAND_POOL_CREATE_PROTECTED_BIT
)

// CommandBufferAllocateInfo contains command buffer allocation information
type CommandBufferAllocateInfo struct {
	CommandPool        CommandPool
	Level              CommandBufferLevel
	CommandBufferCount uint32
}

// CommandBufferLevel represents command buffer levels
type CommandBufferLevel int32

const (
	CommandBufferLevelPrimary   CommandBufferLevel = C.VK_COMMAND_BUFFER_LEVEL_PRIMARY
	CommandBufferLevelSecondary CommandBufferLevel = C.VK_COMMAND_BUFFER_LEVEL_SECONDARY
)

// CommandBufferBeginInfo contains command buffer begin information
type CommandBufferBeginInfo struct {
	Flags CommandBufferUsageFlags
}

// CommandBufferUsageFlags represents command buffer usage flags
type CommandBufferUsageFlags uint32

const (
	CommandBufferUsageOneTimeSubmitBit      CommandBufferUsageFlags = C.VK_COMMAND_BUFFER_USAGE_ONE_TIME_SUBMIT_BIT
	CommandBufferUsageRenderPassContinueBit CommandBufferUsageFlags = C.VK_COMMAND_BUFFER_USAGE_RENDER_PASS_CONTINUE_BIT
	CommandBufferUsageSimultaneousUseBit    CommandBufferUsageFlags = C.VK_COMMAND_BUFFER_USAGE_SIMULTANEOUS_USE_BIT
)

// SubmitInfo contains queue submit information
type SubmitInfo struct {
	WaitSemaphores    []Semaphore
	WaitDstStageMask  []PipelineStageFlags
	CommandBuffers    []CommandBuffer
	SignalSemaphores  []Semaphore
}

// PipelineStageFlags represents pipeline stage flags
type PipelineStageFlags uint32

const (
	PipelineStageTopOfPipeBit                    PipelineStageFlags = C.VK_PIPELINE_STAGE_TOP_OF_PIPE_BIT
	PipelineStageDrawIndirectBit                 PipelineStageFlags = C.VK_PIPELINE_STAGE_DRAW_INDIRECT_BIT
	PipelineStageVertexInputBit                  PipelineStageFlags = C.VK_PIPELINE_STAGE_VERTEX_INPUT_BIT
	PipelineStageVertexShaderBit                 PipelineStageFlags = C.VK_PIPELINE_STAGE_VERTEX_SHADER_BIT
	PipelineStageTessellationControlShaderBit    PipelineStageFlags = C.VK_PIPELINE_STAGE_TESSELLATION_CONTROL_SHADER_BIT
	PipelineStageTessellationEvaluationShaderBit PipelineStageFlags = C.VK_PIPELINE_STAGE_TESSELLATION_EVALUATION_SHADER_BIT
	PipelineStageGeometryShaderBit               PipelineStageFlags = C.VK_PIPELINE_STAGE_GEOMETRY_SHADER_BIT
	PipelineStageFragmentShaderBit               PipelineStageFlags = C.VK_PIPELINE_STAGE_FRAGMENT_SHADER_BIT
	PipelineStageEarlyFragmentTestsBit           PipelineStageFlags = C.VK_PIPELINE_STAGE_EARLY_FRAGMENT_TESTS_BIT
	PipelineStageLateFragmentTestsBit            PipelineStageFlags = C.VK_PIPELINE_STAGE_LATE_FRAGMENT_TESTS_BIT
	PipelineStageColorAttachmentOutputBit        PipelineStageFlags = C.VK_PIPELINE_STAGE_COLOR_ATTACHMENT_OUTPUT_BIT
	PipelineStageComputeShaderBit                PipelineStageFlags = C.VK_PIPELINE_STAGE_COMPUTE_SHADER_BIT
	PipelineStageTransferBit                     PipelineStageFlags = C.VK_PIPELINE_STAGE_TRANSFER_BIT
	PipelineStageBottomOfPipeBit                 PipelineStageFlags = C.VK_PIPELINE_STAGE_BOTTOM_OF_PIPE_BIT
	PipelineStageHostBit                         PipelineStageFlags = C.VK_PIPELINE_STAGE_HOST_BIT
	PipelineStageAllGraphicsBit                  PipelineStageFlags = C.VK_PIPELINE_STAGE_ALL_GRAPHICS_BIT
	PipelineStageAllCommandsBit                  PipelineStageFlags = C.VK_PIPELINE_STAGE_ALL_COMMANDS_BIT
)

// SemaphoreCreateInfo contains semaphore creation information
type SemaphoreCreateInfo struct {
	// No fields needed for basic semaphore creation
}

// FenceCreateInfo contains fence creation information
type FenceCreateInfo struct {
	Flags FenceCreateFlags
}

// FenceCreateFlags represents fence creation flags
type FenceCreateFlags uint32

const (
	FenceCreateSignaledBit FenceCreateFlags = C.VK_FENCE_CREATE_SIGNALED_BIT
)

// CreateCommandPool creates a command pool
func CreateCommandPool(device Device, createInfo *CommandPoolCreateInfo) (CommandPool, error) {
	var cCreateInfo C.VkCommandPoolCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_COMMAND_POOL_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkCommandPoolCreateFlags(createInfo.Flags)
	cCreateInfo.queueFamilyIndex = C.uint32_t(createInfo.QueueFamilyIndex)

	var commandPool C.VkCommandPool
	result := Result(C.vkCreateCommandPool(C.VkDevice(device), &cCreateInfo, nil, &commandPool))
	if result != Success {
		return nil, result
	}

	return CommandPool(commandPool), nil
}

// DestroyCommandPool destroys a command pool
func DestroyCommandPool(device Device, commandPool CommandPool) {
	C.vkDestroyCommandPool(C.VkDevice(device), C.VkCommandPool(commandPool), nil)
}

// AllocateCommandBuffers allocates command buffers
func AllocateCommandBuffers(device Device, allocateInfo *CommandBufferAllocateInfo) ([]CommandBuffer, error) {
	var cAllocateInfo C.VkCommandBufferAllocateInfo
	cAllocateInfo.sType = C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_ALLOCATE_INFO
	cAllocateInfo.pNext = nil
	cAllocateInfo.commandPool = C.VkCommandPool(allocateInfo.CommandPool)
	cAllocateInfo.level = C.VkCommandBufferLevel(allocateInfo.Level)
	cAllocateInfo.commandBufferCount = C.uint32_t(allocateInfo.CommandBufferCount)

	cCommandBuffers := make([]C.VkCommandBuffer, allocateInfo.CommandBufferCount)
	result := Result(C.vkAllocateCommandBuffers(C.VkDevice(device), &cAllocateInfo, &cCommandBuffers[0]))
	if result != Success {
		return nil, result
	}

	commandBuffers := make([]CommandBuffer, allocateInfo.CommandBufferCount)
	for i := range commandBuffers {
		commandBuffers[i] = CommandBuffer(cCommandBuffers[i])
	}

	return commandBuffers, nil
}

// FreeCommandBuffers frees command buffers
func FreeCommandBuffers(device Device, commandPool CommandPool, commandBuffers []CommandBuffer) {
	if len(commandBuffers) == 0 {
		return
	}

	cCommandBuffers := make([]C.VkCommandBuffer, len(commandBuffers))
	for i, cb := range commandBuffers {
		cCommandBuffers[i] = C.VkCommandBuffer(cb)
	}

	C.vkFreeCommandBuffers(C.VkDevice(device), C.VkCommandPool(commandPool), C.uint32_t(len(cCommandBuffers)), &cCommandBuffers[0])
}

// BeginCommandBuffer begins recording a command buffer
func BeginCommandBuffer(commandBuffer CommandBuffer, beginInfo *CommandBufferBeginInfo) error {
	var cBeginInfo C.VkCommandBufferBeginInfo
	cBeginInfo.sType = C.VK_STRUCTURE_TYPE_COMMAND_BUFFER_BEGIN_INFO
	cBeginInfo.pNext = nil
	cBeginInfo.flags = C.VkCommandBufferUsageFlags(beginInfo.Flags)
	cBeginInfo.pInheritanceInfo = nil

	result := Result(C.vkBeginCommandBuffer(C.VkCommandBuffer(commandBuffer), &cBeginInfo))
	if result != Success {
		return result
	}
	return nil
}

// EndCommandBuffer ends recording a command buffer
func EndCommandBuffer(commandBuffer CommandBuffer) error {
	result := Result(C.vkEndCommandBuffer(C.VkCommandBuffer(commandBuffer)))
	if result != Success {
		return result
	}
	return nil
}

// QueueSubmit submits command buffers to a queue
func QueueSubmit(queue Queue, submitInfos []SubmitInfo, fence Fence) error {
	if len(submitInfos) == 0 {
		result := Result(C.vkQueueSubmit(C.VkQueue(queue), 0, nil, C.VkFence(fence)))
		if result != Success {
			return result
		}
		return nil
	}

	cSubmitInfos := make([]C.VkSubmitInfo, len(submitInfos))
	
	// We need to keep slices alive during the call
	var allWaitSemaphores [][]C.VkSemaphore
	var allWaitStages [][]C.VkPipelineStageFlags
	var allCommandBuffers [][]C.VkCommandBuffer
	var allSignalSemaphores [][]C.VkSemaphore

	for i, si := range submitInfos {
		cSubmitInfos[i].sType = C.VK_STRUCTURE_TYPE_SUBMIT_INFO
		cSubmitInfos[i].pNext = nil

		// Wait semaphores
		if len(si.WaitSemaphores) > 0 {
			waitSems := make([]C.VkSemaphore, len(si.WaitSemaphores))
			for j, sem := range si.WaitSemaphores {
				waitSems[j] = C.VkSemaphore(sem)
			}
			allWaitSemaphores = append(allWaitSemaphores, waitSems)
			cSubmitInfos[i].waitSemaphoreCount = C.uint32_t(len(waitSems))
			cSubmitInfos[i].pWaitSemaphores = &waitSems[0]
		}

		// Wait stages
		if len(si.WaitDstStageMask) > 0 {
			waitStages := make([]C.VkPipelineStageFlags, len(si.WaitDstStageMask))
			for j, stage := range si.WaitDstStageMask {
				waitStages[j] = C.VkPipelineStageFlags(stage)
			}
			allWaitStages = append(allWaitStages, waitStages)
			cSubmitInfos[i].pWaitDstStageMask = &waitStages[0]
		}

		// Command buffers
		if len(si.CommandBuffers) > 0 {
			cmdBufs := make([]C.VkCommandBuffer, len(si.CommandBuffers))
			for j, cb := range si.CommandBuffers {
				cmdBufs[j] = C.VkCommandBuffer(cb)
			}
			allCommandBuffers = append(allCommandBuffers, cmdBufs)
			cSubmitInfos[i].commandBufferCount = C.uint32_t(len(cmdBufs))
			cSubmitInfos[i].pCommandBuffers = &cmdBufs[0]
		}

		// Signal semaphores
		if len(si.SignalSemaphores) > 0 {
			signalSems := make([]C.VkSemaphore, len(si.SignalSemaphores))
			for j, sem := range si.SignalSemaphores {
				signalSems[j] = C.VkSemaphore(sem)
			}
			allSignalSemaphores = append(allSignalSemaphores, signalSems)
			cSubmitInfos[i].signalSemaphoreCount = C.uint32_t(len(signalSems))
			cSubmitInfos[i].pSignalSemaphores = &signalSems[0]
		}
	}

	result := Result(C.vkQueueSubmit(C.VkQueue(queue), C.uint32_t(len(cSubmitInfos)), &cSubmitInfos[0], C.VkFence(fence)))
	if result != Success {
		return result
	}
	return nil
}

// CreateSemaphore creates a semaphore
func CreateSemaphore(device Device, createInfo *SemaphoreCreateInfo) (Semaphore, error) {
	var cCreateInfo C.VkSemaphoreCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_SEMAPHORE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0

	var semaphore C.VkSemaphore
	result := Result(C.vkCreateSemaphore(C.VkDevice(device), &cCreateInfo, nil, &semaphore))
	if result != Success {
		return nil, result
	}

	return Semaphore(semaphore), nil
}

// DestroySemaphore destroys a semaphore
func DestroySemaphore(device Device, semaphore Semaphore) {
	C.vkDestroySemaphore(C.VkDevice(device), C.VkSemaphore(semaphore), nil)
}

// CreateFence creates a fence
func CreateFence(device Device, createInfo *FenceCreateInfo) (Fence, error) {
	var cCreateInfo C.VkFenceCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_FENCE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkFenceCreateFlags(createInfo.Flags)

	var fence C.VkFence
	result := Result(C.vkCreateFence(C.VkDevice(device), &cCreateInfo, nil, &fence))
	if result != Success {
		return nil, result
	}

	return Fence(fence), nil
}

// DestroyFence destroys a fence
func DestroyFence(device Device, fence Fence) {
	C.vkDestroyFence(C.VkDevice(device), C.VkFence(fence), nil)
}

// WaitForFences waits for fences to be signaled
func WaitForFences(device Device, fences []Fence, waitAll bool, timeout uint64) error {
	if len(fences) == 0 {
		return nil
	}

	cFences := make([]C.VkFence, len(fences))
	for i, fence := range fences {
		cFences[i] = C.VkFence(fence)
	}

	var cWaitAll C.VkBool32
	if waitAll {
		cWaitAll = C.VK_TRUE
	} else {
		cWaitAll = C.VK_FALSE
	}

	result := Result(C.vkWaitForFences(C.VkDevice(device), C.uint32_t(len(cFences)), &cFences[0], cWaitAll, C.uint64_t(timeout)))
	if result != Success {
		return result
	}
	return nil
}

// ResetFences resets fences
func ResetFences(device Device, fences []Fence) error {
	if len(fences) == 0 {
		return nil
	}

	cFences := make([]C.VkFence, len(fences))
	for i, fence := range fences {
		cFences[i] = C.VkFence(fence)
	}

	result := Result(C.vkResetFences(C.VkDevice(device), C.uint32_t(len(cFences)), &cFences[0]))
	if result != Success {
		return result
	}
	return nil
}

// GetFenceStatus gets fence status
func GetFenceStatus(device Device, fence Fence) Result {
	return Result(C.vkGetFenceStatus(C.VkDevice(device), C.VkFence(fence)))
}