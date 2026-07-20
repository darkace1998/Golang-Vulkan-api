package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// ============================================================================
// Query Pool Types
// ============================================================================

// QueryType represents the type of queries managed by a query pool
type QueryType uint32

// QueryType constants represent the types of queries supported
const (
	QueryTypeOcclusion          QueryType = C.VK_QUERY_TYPE_OCCLUSION
	QueryTypePipelineStatistics QueryType = C.VK_QUERY_TYPE_PIPELINE_STATISTICS
	QueryTypeTimestamp          QueryType = C.VK_QUERY_TYPE_TIMESTAMP
)

// QueryPoolCreateFlags represents query pool creation flags
type QueryPoolCreateFlags uint32

// QueryPoolCreateInfo contains query pool creation parameters
type QueryPoolCreateInfo struct {
	Flags              QueryPoolCreateFlags
	QueryType          QueryType
	QueryCount         uint32
	PipelineStatistics QueryPipelineStatisticFlags
}

// QueryResultFlags represents query result retrieval flags
type QueryResultFlags uint32

// QueryResultFlags constants represent flags for retrieving query results
const (
	QueryResult64Bit            QueryResultFlags = C.VK_QUERY_RESULT_64_BIT
	QueryResultWait             QueryResultFlags = C.VK_QUERY_RESULT_WAIT_BIT
	QueryResultWithAvailability QueryResultFlags = C.VK_QUERY_RESULT_WITH_AVAILABILITY_BIT
	QueryResultPartial          QueryResultFlags = C.VK_QUERY_RESULT_PARTIAL_BIT
	QueryResultWithStatusKHR    QueryResultFlags = 0x00000010 // VK_QUERY_RESULT_WITH_STATUS_BIT_KHR
)

// ============================================================================
// Query Pool Functions
// ============================================================================

// CreateQueryPool creates a query pool for managing a number of queries
func CreateQueryPool(device Device, createInfo *QueryPoolCreateInfo) (QueryPool, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.QueryCount == 0 {
		return nil, NewValidationError("createInfo.QueryCount", "must be greater than 0")
	}

	var cCreateInfo C.VkQueryPoolCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_QUERY_POOL_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkQueryPoolCreateFlags(createInfo.Flags)
	cCreateInfo.queryType = C.VkQueryType(createInfo.QueryType)
	cCreateInfo.queryCount = C.uint32_t(createInfo.QueryCount)
	cCreateInfo.pipelineStatistics = C.VkQueryPipelineStatisticFlags(createInfo.PipelineStatistics)

	var queryPool C.VkQueryPool
	result := Result(C.vkCreateQueryPool(C.VkDevice(device), &cCreateInfo, nil, &queryPool))
	if result != Success {
		return nil, NewVulkanError(result, "CreateQueryPool", "failed to create query pool")
	}

	return QueryPool(queryPool), nil
}

// DestroyQueryPool destroys a query pool
func DestroyQueryPool(device Device, queryPool QueryPool) {
	if device == nil || queryPool == nil {
		return
	}
	C.vkDestroyQueryPool(C.VkDevice(device), C.VkQueryPool(queryPool), nil)
}

// GetQueryPoolResults retrieves results from a query pool
// Returns the query results as a byte slice, or an error if the operation fails
// Use QueryResult64Bit flag for 64-bit results, otherwise 32-bit results are returned
func GetQueryPoolResults(device Device, queryPool QueryPool, firstQuery, queryCount uint32, dataSize uint64, flags QueryResultFlags) ([]byte, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if queryPool == nil {
		return nil, NewValidationError("queryPool", "cannot be nil")
	}
	if queryCount == 0 {
		return nil, NewValidationError("queryCount", "must be greater than 0")
	}
	if dataSize == 0 {
		return nil, NewValidationError("dataSize", "must be greater than 0")
	}

	data := make([]byte, dataSize)
	var stride C.VkDeviceSize
	if flags&QueryResult64Bit != 0 {
		stride = 8
	} else {
		stride = 4
	}
	// Add availability if requested
	if flags&QueryResultWithAvailability != 0 {
		stride *= 2
	}

	result := Result(C.vkGetQueryPoolResults(
		C.VkDevice(device),
		C.VkQueryPool(queryPool),
		C.uint32_t(firstQuery),
		C.uint32_t(queryCount),
		C.size_t(dataSize),
		unsafe.Pointer(&data[0]),
		stride,
		C.VkQueryResultFlags(flags),
	))

	// NotReady is not an error for non-wait queries
	if result != Success && result != NotReady {
		return nil, NewVulkanError(result, "GetQueryPoolResults", "failed to get query pool results")
	}

	return data, nil
}

// validateQueryPoolResultsParams validates common parameters for query pool results
func validateQueryPoolResultsParams(device Device, queryPool QueryPool, queryCount uint32) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if queryPool == nil {
		return NewValidationError("queryPool", "cannot be nil")
	}
	if queryCount == 0 {
		return NewValidationError("queryCount", "must be greater than 0")
	}
	return nil
}

// GetQueryPoolResultsUint32 retrieves 32-bit query results
func GetQueryPoolResultsUint32(device Device, queryPool QueryPool, firstQuery, queryCount uint32, flags QueryResultFlags) ([]uint32, error) {
	if err := validateQueryPoolResultsParams(device, queryPool, queryCount); err != nil {
		return nil, err
	}
	flags = flags &^ QueryResult64Bit // Mask off 64-bit flag
	resultCount := queryCount
	stride := C.VkDeviceSize(4)
	if flags&QueryResultWithAvailability != 0 {
		resultCount = queryCount * 2
		stride = 8
	}
	results := make([]uint32, resultCount)
	result := Result(C.vkGetQueryPoolResults(C.VkDevice(device), C.VkQueryPool(queryPool),
		C.uint32_t(firstQuery), C.uint32_t(queryCount), C.size_t(len(results)*4),
		unsafe.Pointer(&results[0]), stride, C.VkQueryResultFlags(flags)))
	if result != Success && result != NotReady {
		return nil, NewVulkanError(result, "GetQueryPoolResultsUint32", "failed to get query pool results")
	}
	return results, nil
}

// GetQueryPoolResultsUint64 retrieves 64-bit query results
func GetQueryPoolResultsUint64(device Device, queryPool QueryPool, firstQuery, queryCount uint32, flags QueryResultFlags) ([]uint64, error) {
	if err := validateQueryPoolResultsParams(device, queryPool, queryCount); err != nil {
		return nil, err
	}
	flags = flags | QueryResult64Bit // Ensure 64-bit flag is set
	resultCount := queryCount
	stride := C.VkDeviceSize(8)
	if flags&QueryResultWithAvailability != 0 {
		resultCount = queryCount * 2
		stride = 16
	}
	results := make([]uint64, resultCount)
	result := Result(C.vkGetQueryPoolResults(C.VkDevice(device), C.VkQueryPool(queryPool),
		C.uint32_t(firstQuery), C.uint32_t(queryCount), C.size_t(len(results)*8),
		unsafe.Pointer(&results[0]), stride, C.VkQueryResultFlags(flags)))
	if result != Success && result != NotReady {
		return nil, NewVulkanError(result, "GetQueryPoolResultsUint64", "failed to get query pool results")
	}
	return results, nil
}

// ============================================================================
// Query Commands
// ============================================================================

// CmdBeginQuery begins a query
func CmdBeginQuery(commandBuffer CommandBuffer, queryPool QueryPool, query uint32, flags QueryControlFlags) {
	if commandBuffer == nil || queryPool == nil {
		return
	}
	C.vkCmdBeginQuery(
		C.VkCommandBuffer(commandBuffer),
		C.VkQueryPool(queryPool),
		C.uint32_t(query),
		C.VkQueryControlFlags(flags),
	)
}

// CmdEndQuery ends a query
func CmdEndQuery(commandBuffer CommandBuffer, queryPool QueryPool, query uint32) {
	if commandBuffer == nil || queryPool == nil {
		return
	}
	C.vkCmdEndQuery(
		C.VkCommandBuffer(commandBuffer),
		C.VkQueryPool(queryPool),
		C.uint32_t(query),
	)
}

// CmdResetQueryPool resets a range of queries in a query pool on the GPU
func CmdResetQueryPool(commandBuffer CommandBuffer, queryPool QueryPool, firstQuery, queryCount uint32) {
	if commandBuffer == nil || queryPool == nil {
		return
	}
	C.vkCmdResetQueryPool(
		C.VkCommandBuffer(commandBuffer),
		C.VkQueryPool(queryPool),
		C.uint32_t(firstQuery),
		C.uint32_t(queryCount),
	)
}

// CmdWriteTimestamp writes a device timestamp into a query object
func CmdWriteTimestamp(commandBuffer CommandBuffer, pipelineStage PipelineStageFlags, queryPool QueryPool, query uint32) {
	if commandBuffer == nil || queryPool == nil {
		return
	}
	C.vkCmdWriteTimestamp(
		C.VkCommandBuffer(commandBuffer),
		C.VkPipelineStageFlagBits(pipelineStage),
		C.VkQueryPool(queryPool),
		C.uint32_t(query),
	)
}

// CmdCopyQueryPoolResults copies the results of queries in a query pool to a buffer object
func CmdCopyQueryPoolResults(commandBuffer CommandBuffer, queryPool QueryPool, firstQuery, queryCount uint32, dstBuffer Buffer, dstOffset DeviceSize, stride DeviceSize, flags QueryResultFlags) {
	if commandBuffer == nil || queryPool == nil || dstBuffer == nil {
		return
	}
	C.vkCmdCopyQueryPoolResults(
		C.VkCommandBuffer(commandBuffer),
		C.VkQueryPool(queryPool),
		C.uint32_t(firstQuery),
		C.uint32_t(queryCount),
		C.VkBuffer(dstBuffer),
		C.VkDeviceSize(dstOffset),
		C.VkDeviceSize(stride),
		C.VkQueryResultFlags(flags),
	)
}

// ============================================================================
// Host Query Reset (Vulkan 1.2+)
// ============================================================================

// ResetQueryPool resets a range of queries in a query pool on the host (Vulkan 1.2+)
// This requires the hostQueryReset feature to be enabled
func ResetQueryPool(device Device, queryPool QueryPool, firstQuery, queryCount uint32) {
	if device == nil || queryPool == nil {
		return
	}
	C.vkResetQueryPool(
		C.VkDevice(device),
		C.VkQueryPool(queryPool),
		C.uint32_t(firstQuery),
		C.uint32_t(queryCount),
	)
}
