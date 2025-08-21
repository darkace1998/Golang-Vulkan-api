package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"unsafe"
)

// DeviceQueueCreateInfo contains device queue creation information
type DeviceQueueCreateInfo struct {
	QueueFamilyIndex uint32
	QueuePriorities  []float32
}

// DeviceCreateInfo contains device creation information
type DeviceCreateInfo struct {
	QueueCreateInfos      []DeviceQueueCreateInfo
	EnabledLayerNames     []string
	EnabledExtensionNames []string
	EnabledFeatures       *PhysicalDeviceFeatures
}

// PhysicalDeviceFeatures contains physical device features
type PhysicalDeviceFeatures struct {
	RobustBufferAccess                      bool
	FullDrawIndexUint32                     bool
	ImageCubeArray                          bool
	IndependentBlend                        bool
	GeometryShader                          bool
	TessellationShader                      bool
	SampleRateShading                       bool
	DualSrcBlend                            bool
	LogicOp                                 bool
	MultiDrawIndirect                       bool
	DrawIndirectFirstInstance               bool
	DepthClamp                              bool
	DepthBiasClamp                          bool
	FillModeNonSolid                        bool
	DepthBounds                             bool
	WideLines                               bool
	LargePoints                             bool
	AlphaToOne                              bool
	MultiViewport                           bool
	SamplerAnisotropy                       bool
	TextureCompressionETC2                  bool
	TextureCompressionASTC_LDR              bool
	TextureCompressionBC                    bool
	OcclusionQueryPrecise                   bool
	PipelineStatisticsQuery                 bool
	VertexPipelineStoresAndAtomics          bool
	FragmentStoresAndAtomics                bool
	ShaderTessellationAndGeometryPointSize  bool
	ShaderImageGatherExtended               bool
	ShaderStorageImageExtendedFormats       bool
	ShaderStorageImageMultisample           bool
	ShaderStorageImageReadWithoutFormat     bool
	ShaderStorageImageWriteWithoutFormat    bool
	ShaderUniformBufferArrayDynamicIndexing bool
	ShaderSampledImageArrayDynamicIndexing  bool
	ShaderStorageBufferArrayDynamicIndexing bool
	ShaderStorageImageArrayDynamicIndexing  bool
	ShaderClipDistance                      bool
	ShaderCullDistance                      bool
	ShaderFloat64                           bool
	ShaderInt64                             bool
	ShaderInt16                             bool
	ShaderResourceResidency                 bool
	ShaderResourceMinLod                    bool
	SparseBinding                           bool
	SparseResidencyBuffer                   bool
	SparseResidencyImage2D                  bool
	SparseResidencyImage3D                  bool
	SparseResidency2Samples                 bool
	SparseResidency4Samples                 bool
	SparseResidency8Samples                 bool
	SparseResidency16Samples                bool
	SparseResidencyAliased                  bool
	VariableMultisampleRate                 bool
	InheritedQueries                        bool
}

// PhysicalDeviceMemoryProperties contains memory properties
type PhysicalDeviceMemoryProperties struct {
	MemoryTypeCount uint32
	MemoryTypes     [MaxMemoryTypes]MemoryType
	MemoryHeapCount uint32
	MemoryHeaps     [MaxMemoryHeaps]MemoryHeap
}

// MemoryType contains memory type information
type MemoryType struct {
	PropertyFlags MemoryPropertyFlags
	HeapIndex     uint32
}

// MemoryHeap contains memory heap information
type MemoryHeap struct {
	Size  DeviceSize
	Flags MemoryHeapFlags
}

// MemoryPropertyFlags represents memory property flags
type MemoryPropertyFlags uint32

const (
	MemoryPropertyDeviceLocalBit     MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_DEVICE_LOCAL_BIT
	MemoryPropertyHostVisibleBit     MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT
	MemoryPropertyHostCoherentBit    MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_HOST_COHERENT_BIT
	MemoryPropertyHostCachedBit      MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_HOST_CACHED_BIT
	MemoryPropertyLazilyAllocatedBit MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_LAZILY_ALLOCATED_BIT
	MemoryPropertyProtectedBit       MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_PROTECTED_BIT
	MemoryPropertyDeviceCoherentBit  MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_DEVICE_COHERENT_BIT_AMD
	MemoryPropertyDeviceUncachedBit  MemoryPropertyFlags = C.VK_MEMORY_PROPERTY_DEVICE_UNCACHED_BIT_AMD
)

// MemoryHeapFlags represents memory heap flags
type MemoryHeapFlags uint32

const (
	MemoryHeapDeviceLocalBit   MemoryHeapFlags = C.VK_MEMORY_HEAP_DEVICE_LOCAL_BIT
	MemoryHeapMultiInstanceBit MemoryHeapFlags = C.VK_MEMORY_HEAP_MULTI_INSTANCE_BIT
)

// CreateDevice creates a logical device
func CreateDevice(physicalDevice PhysicalDevice, createInfo *DeviceCreateInfo) (Device, error) {
	var cCreateInfo C.VkDeviceCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0

	// Queue create infos
	var cQueueCreateInfos []C.VkDeviceQueueCreateInfo
	var cPriorities [][]C.float
	if len(createInfo.QueueCreateInfos) > 0 {
		cQueueCreateInfos = make([]C.VkDeviceQueueCreateInfo, len(createInfo.QueueCreateInfos))
		cPriorities = make([][]C.float, len(createInfo.QueueCreateInfos))

		for i, qci := range createInfo.QueueCreateInfos {
			cQueueCreateInfos[i].sType = C.VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO
			cQueueCreateInfos[i].pNext = nil
			cQueueCreateInfos[i].flags = 0
			cQueueCreateInfos[i].queueFamilyIndex = C.uint32_t(qci.QueueFamilyIndex)
			cQueueCreateInfos[i].queueCount = C.uint32_t(len(qci.QueuePriorities))

			if len(qci.QueuePriorities) > 0 {
				cPriorities[i] = make([]C.float, len(qci.QueuePriorities))
				for j, priority := range qci.QueuePriorities {
					cPriorities[i][j] = C.float(priority)
				}
				cQueueCreateInfos[i].pQueuePriorities = &cPriorities[i][0]
			}
		}
		cCreateInfo.queueCreateInfoCount = C.uint32_t(len(cQueueCreateInfos))
		cCreateInfo.pQueueCreateInfos = &cQueueCreateInfos[0]
	}

	// Enabled layers
	var cLayers **C.char
	if len(createInfo.EnabledLayerNames) > 0 {
		cLayers = stringSliceToCharArray(createInfo.EnabledLayerNames)
		defer freeStringArray(cLayers, len(createInfo.EnabledLayerNames))
		cCreateInfo.enabledLayerCount = C.uint32_t(len(createInfo.EnabledLayerNames))
		cCreateInfo.ppEnabledLayerNames = cLayers
	}

	// Enabled extensions
	var cExtensions **C.char
	if len(createInfo.EnabledExtensionNames) > 0 {
		cExtensions = stringSliceToCharArray(createInfo.EnabledExtensionNames)
		defer freeStringArray(cExtensions, len(createInfo.EnabledExtensionNames))
		cCreateInfo.enabledExtensionCount = C.uint32_t(len(createInfo.EnabledExtensionNames))
		cCreateInfo.ppEnabledExtensionNames = cExtensions
	}

	// Enabled features
	var cFeatures C.VkPhysicalDeviceFeatures
	if createInfo.EnabledFeatures != nil {
		cFeatures = physicalDeviceFeaturesToC(createInfo.EnabledFeatures)
		cCreateInfo.pEnabledFeatures = &cFeatures
	}

	var device C.VkDevice
	result := Result(C.vkCreateDevice(C.VkPhysicalDevice(physicalDevice), &cCreateInfo, nil, &device))
	if result != Success {
		return nil, result
	}

	return Device(device), nil
}

// DestroyDevice destroys a logical device
func DestroyDevice(device Device) {
	C.vkDestroyDevice(C.VkDevice(device), nil)
}

// GetDeviceQueue gets a device queue
func GetDeviceQueue(device Device, queueFamilyIndex, queueIndex uint32) Queue {
	var queue C.VkQueue
	C.vkGetDeviceQueue(C.VkDevice(device), C.uint32_t(queueFamilyIndex), C.uint32_t(queueIndex), &queue)
	return Queue(queue)
}

// QueueWaitIdle waits for a queue to become idle
func QueueWaitIdle(queue Queue) error {
	result := Result(C.vkQueueWaitIdle(C.VkQueue(queue)))
	if result != Success {
		return result
	}
	return nil
}

// DeviceWaitIdle waits for a device to become idle
func DeviceWaitIdle(device Device) error {
	result := Result(C.vkDeviceWaitIdle(C.VkDevice(device)))
	if result != Success {
		return result
	}
	return nil
}

// GetPhysicalDeviceFeatures gets physical device features
func GetPhysicalDeviceFeatures(physicalDevice PhysicalDevice) PhysicalDeviceFeatures {
	var cFeatures C.VkPhysicalDeviceFeatures
	C.vkGetPhysicalDeviceFeatures(C.VkPhysicalDevice(physicalDevice), &cFeatures)
	return physicalDeviceFeaturesFromC(&cFeatures)
}

// GetPhysicalDeviceMemoryProperties gets physical device memory properties
func GetPhysicalDeviceMemoryProperties(physicalDevice PhysicalDevice) PhysicalDeviceMemoryProperties {
	var cProps C.VkPhysicalDeviceMemoryProperties
	C.vkGetPhysicalDeviceMemoryProperties(C.VkPhysicalDevice(physicalDevice), &cProps)

	props := PhysicalDeviceMemoryProperties{
		MemoryTypeCount: uint32(cProps.memoryTypeCount),
		MemoryHeapCount: uint32(cProps.memoryHeapCount),
	}

	for i := uint32(0); i < props.MemoryTypeCount; i++ {
		props.MemoryTypes[i] = MemoryType{
			PropertyFlags: MemoryPropertyFlags(cProps.memoryTypes[i].propertyFlags),
			HeapIndex:     uint32(cProps.memoryTypes[i].heapIndex),
		}
	}

	for i := uint32(0); i < props.MemoryHeapCount; i++ {
		props.MemoryHeaps[i] = MemoryHeap{
			Size:  DeviceSize(cProps.memoryHeaps[i].size),
			Flags: MemoryHeapFlags(cProps.memoryHeaps[i].flags),
		}
	}

	return props
}

// EnumerateDeviceExtensionProperties enumerates device extension properties
func EnumerateDeviceExtensionProperties(physicalDevice PhysicalDevice, layerName string) ([]ExtensionProperties, error) {
	var cLayerName *C.char
	if layerName != "" {
		cLayerName = C.CString(layerName)
		defer C.free(unsafe.Pointer(cLayerName))
	}

	var propertyCount C.uint32_t
	result := Result(C.vkEnumerateDeviceExtensionProperties(C.VkPhysicalDevice(physicalDevice), cLayerName, &propertyCount, nil))
	if result != Success {
		return nil, result
	}

	if propertyCount == 0 {
		return nil, nil
	}

	cProperties := make([]C.VkExtensionProperties, propertyCount)
	result = Result(C.vkEnumerateDeviceExtensionProperties(C.VkPhysicalDevice(physicalDevice), cLayerName, &propertyCount, &cProperties[0]))
	if result != Success {
		return nil, result
	}

	properties := make([]ExtensionProperties, propertyCount)
	for i := range properties {
		properties[i].ExtensionName = C.GoString(&cProperties[i].extensionName[0])
		properties[i].SpecVersion = uint32(cProperties[i].specVersion)
	}

	return properties, nil
}

// Helper function to convert Go PhysicalDeviceFeatures to C VkPhysicalDeviceFeatures
func physicalDeviceFeaturesToC(features *PhysicalDeviceFeatures) C.VkPhysicalDeviceFeatures {
	var cFeatures C.VkPhysicalDeviceFeatures
	cFeatures.robustBufferAccess = boolToVkBool32(features.RobustBufferAccess)
	cFeatures.fullDrawIndexUint32 = boolToVkBool32(features.FullDrawIndexUint32)
	cFeatures.imageCubeArray = boolToVkBool32(features.ImageCubeArray)
	cFeatures.independentBlend = boolToVkBool32(features.IndependentBlend)
	cFeatures.geometryShader = boolToVkBool32(features.GeometryShader)
	cFeatures.tessellationShader = boolToVkBool32(features.TessellationShader)
	cFeatures.sampleRateShading = boolToVkBool32(features.SampleRateShading)
	cFeatures.dualSrcBlend = boolToVkBool32(features.DualSrcBlend)
	cFeatures.logicOp = boolToVkBool32(features.LogicOp)
	cFeatures.multiDrawIndirect = boolToVkBool32(features.MultiDrawIndirect)
	cFeatures.drawIndirectFirstInstance = boolToVkBool32(features.DrawIndirectFirstInstance)
	cFeatures.depthClamp = boolToVkBool32(features.DepthClamp)
	cFeatures.depthBiasClamp = boolToVkBool32(features.DepthBiasClamp)
	cFeatures.fillModeNonSolid = boolToVkBool32(features.FillModeNonSolid)
	cFeatures.depthBounds = boolToVkBool32(features.DepthBounds)
	cFeatures.wideLines = boolToVkBool32(features.WideLines)
	cFeatures.largePoints = boolToVkBool32(features.LargePoints)
	cFeatures.alphaToOne = boolToVkBool32(features.AlphaToOne)
	cFeatures.multiViewport = boolToVkBool32(features.MultiViewport)
	cFeatures.samplerAnisotropy = boolToVkBool32(features.SamplerAnisotropy)
	cFeatures.textureCompressionETC2 = boolToVkBool32(features.TextureCompressionETC2)
	cFeatures.textureCompressionASTC_LDR = boolToVkBool32(features.TextureCompressionASTC_LDR)
	cFeatures.textureCompressionBC = boolToVkBool32(features.TextureCompressionBC)
	cFeatures.occlusionQueryPrecise = boolToVkBool32(features.OcclusionQueryPrecise)
	cFeatures.pipelineStatisticsQuery = boolToVkBool32(features.PipelineStatisticsQuery)
	cFeatures.vertexPipelineStoresAndAtomics = boolToVkBool32(features.VertexPipelineStoresAndAtomics)
	cFeatures.fragmentStoresAndAtomics = boolToVkBool32(features.FragmentStoresAndAtomics)
	cFeatures.shaderTessellationAndGeometryPointSize = boolToVkBool32(features.ShaderTessellationAndGeometryPointSize)
	cFeatures.shaderImageGatherExtended = boolToVkBool32(features.ShaderImageGatherExtended)
	cFeatures.shaderStorageImageExtendedFormats = boolToVkBool32(features.ShaderStorageImageExtendedFormats)
	cFeatures.shaderStorageImageMultisample = boolToVkBool32(features.ShaderStorageImageMultisample)
	cFeatures.shaderStorageImageReadWithoutFormat = boolToVkBool32(features.ShaderStorageImageReadWithoutFormat)
	cFeatures.shaderStorageImageWriteWithoutFormat = boolToVkBool32(features.ShaderStorageImageWriteWithoutFormat)
	cFeatures.shaderUniformBufferArrayDynamicIndexing = boolToVkBool32(features.ShaderUniformBufferArrayDynamicIndexing)
	cFeatures.shaderSampledImageArrayDynamicIndexing = boolToVkBool32(features.ShaderSampledImageArrayDynamicIndexing)
	cFeatures.shaderStorageBufferArrayDynamicIndexing = boolToVkBool32(features.ShaderStorageBufferArrayDynamicIndexing)
	cFeatures.shaderStorageImageArrayDynamicIndexing = boolToVkBool32(features.ShaderStorageImageArrayDynamicIndexing)
	cFeatures.shaderClipDistance = boolToVkBool32(features.ShaderClipDistance)
	cFeatures.shaderCullDistance = boolToVkBool32(features.ShaderCullDistance)
	cFeatures.shaderFloat64 = boolToVkBool32(features.ShaderFloat64)
	cFeatures.shaderInt64 = boolToVkBool32(features.ShaderInt64)
	cFeatures.shaderInt16 = boolToVkBool32(features.ShaderInt16)
	cFeatures.shaderResourceResidency = boolToVkBool32(features.ShaderResourceResidency)
	cFeatures.shaderResourceMinLod = boolToVkBool32(features.ShaderResourceMinLod)
	cFeatures.sparseBinding = boolToVkBool32(features.SparseBinding)
	cFeatures.sparseResidencyBuffer = boolToVkBool32(features.SparseResidencyBuffer)
	cFeatures.sparseResidencyImage2D = boolToVkBool32(features.SparseResidencyImage2D)
	cFeatures.sparseResidencyImage3D = boolToVkBool32(features.SparseResidencyImage3D)
	cFeatures.sparseResidency2Samples = boolToVkBool32(features.SparseResidency2Samples)
	cFeatures.sparseResidency4Samples = boolToVkBool32(features.SparseResidency4Samples)
	cFeatures.sparseResidency8Samples = boolToVkBool32(features.SparseResidency8Samples)
	cFeatures.sparseResidency16Samples = boolToVkBool32(features.SparseResidency16Samples)
	cFeatures.sparseResidencyAliased = boolToVkBool32(features.SparseResidencyAliased)
	cFeatures.variableMultisampleRate = boolToVkBool32(features.VariableMultisampleRate)
	cFeatures.inheritedQueries = boolToVkBool32(features.InheritedQueries)
	return cFeatures
}

// Helper function to convert C VkPhysicalDeviceFeatures to Go PhysicalDeviceFeatures
func physicalDeviceFeaturesFromC(cFeatures *C.VkPhysicalDeviceFeatures) PhysicalDeviceFeatures {
	return PhysicalDeviceFeatures{
		RobustBufferAccess:                      vkBool32ToBool(cFeatures.robustBufferAccess),
		FullDrawIndexUint32:                     vkBool32ToBool(cFeatures.fullDrawIndexUint32),
		ImageCubeArray:                          vkBool32ToBool(cFeatures.imageCubeArray),
		IndependentBlend:                        vkBool32ToBool(cFeatures.independentBlend),
		GeometryShader:                          vkBool32ToBool(cFeatures.geometryShader),
		TessellationShader:                      vkBool32ToBool(cFeatures.tessellationShader),
		SampleRateShading:                       vkBool32ToBool(cFeatures.sampleRateShading),
		DualSrcBlend:                            vkBool32ToBool(cFeatures.dualSrcBlend),
		LogicOp:                                 vkBool32ToBool(cFeatures.logicOp),
		MultiDrawIndirect:                       vkBool32ToBool(cFeatures.multiDrawIndirect),
		DrawIndirectFirstInstance:               vkBool32ToBool(cFeatures.drawIndirectFirstInstance),
		DepthClamp:                              vkBool32ToBool(cFeatures.depthClamp),
		DepthBiasClamp:                          vkBool32ToBool(cFeatures.depthBiasClamp),
		FillModeNonSolid:                        vkBool32ToBool(cFeatures.fillModeNonSolid),
		DepthBounds:                             vkBool32ToBool(cFeatures.depthBounds),
		WideLines:                               vkBool32ToBool(cFeatures.wideLines),
		LargePoints:                             vkBool32ToBool(cFeatures.largePoints),
		AlphaToOne:                              vkBool32ToBool(cFeatures.alphaToOne),
		MultiViewport:                           vkBool32ToBool(cFeatures.multiViewport),
		SamplerAnisotropy:                       vkBool32ToBool(cFeatures.samplerAnisotropy),
		TextureCompressionETC2:                  vkBool32ToBool(cFeatures.textureCompressionETC2),
		TextureCompressionASTC_LDR:              vkBool32ToBool(cFeatures.textureCompressionASTC_LDR),
		TextureCompressionBC:                    vkBool32ToBool(cFeatures.textureCompressionBC),
		OcclusionQueryPrecise:                   vkBool32ToBool(cFeatures.occlusionQueryPrecise),
		PipelineStatisticsQuery:                 vkBool32ToBool(cFeatures.pipelineStatisticsQuery),
		VertexPipelineStoresAndAtomics:          vkBool32ToBool(cFeatures.vertexPipelineStoresAndAtomics),
		FragmentStoresAndAtomics:                vkBool32ToBool(cFeatures.fragmentStoresAndAtomics),
		ShaderTessellationAndGeometryPointSize:  vkBool32ToBool(cFeatures.shaderTessellationAndGeometryPointSize),
		ShaderImageGatherExtended:               vkBool32ToBool(cFeatures.shaderImageGatherExtended),
		ShaderStorageImageExtendedFormats:       vkBool32ToBool(cFeatures.shaderStorageImageExtendedFormats),
		ShaderStorageImageMultisample:           vkBool32ToBool(cFeatures.shaderStorageImageMultisample),
		ShaderStorageImageReadWithoutFormat:     vkBool32ToBool(cFeatures.shaderStorageImageReadWithoutFormat),
		ShaderStorageImageWriteWithoutFormat:    vkBool32ToBool(cFeatures.shaderStorageImageWriteWithoutFormat),
		ShaderUniformBufferArrayDynamicIndexing: vkBool32ToBool(cFeatures.shaderUniformBufferArrayDynamicIndexing),
		ShaderSampledImageArrayDynamicIndexing:  vkBool32ToBool(cFeatures.shaderSampledImageArrayDynamicIndexing),
		ShaderStorageBufferArrayDynamicIndexing: vkBool32ToBool(cFeatures.shaderStorageBufferArrayDynamicIndexing),
		ShaderStorageImageArrayDynamicIndexing:  vkBool32ToBool(cFeatures.shaderStorageImageArrayDynamicIndexing),
		ShaderClipDistance:                      vkBool32ToBool(cFeatures.shaderClipDistance),
		ShaderCullDistance:                      vkBool32ToBool(cFeatures.shaderCullDistance),
		ShaderFloat64:                           vkBool32ToBool(cFeatures.shaderFloat64),
		ShaderInt64:                             vkBool32ToBool(cFeatures.shaderInt64),
		ShaderInt16:                             vkBool32ToBool(cFeatures.shaderInt16),
		ShaderResourceResidency:                 vkBool32ToBool(cFeatures.shaderResourceResidency),
		ShaderResourceMinLod:                    vkBool32ToBool(cFeatures.shaderResourceMinLod),
		SparseBinding:                           vkBool32ToBool(cFeatures.sparseBinding),
		SparseResidencyBuffer:                   vkBool32ToBool(cFeatures.sparseResidencyBuffer),
		SparseResidencyImage2D:                  vkBool32ToBool(cFeatures.sparseResidencyImage2D),
		SparseResidencyImage3D:                  vkBool32ToBool(cFeatures.sparseResidencyImage3D),
		SparseResidency2Samples:                 vkBool32ToBool(cFeatures.sparseResidency2Samples),
		SparseResidency4Samples:                 vkBool32ToBool(cFeatures.sparseResidency4Samples),
		SparseResidency8Samples:                 vkBool32ToBool(cFeatures.sparseResidency8Samples),
		SparseResidency16Samples:                vkBool32ToBool(cFeatures.sparseResidency16Samples),
		SparseResidencyAliased:                  vkBool32ToBool(cFeatures.sparseResidencyAliased),
		VariableMultisampleRate:                 vkBool32ToBool(cFeatures.variableMultisampleRate),
		InheritedQueries:                        vkBool32ToBool(cFeatures.inheritedQueries),
	}
}

// Helper function to convert Go bool to VkBool32
func boolToVkBool32(b bool) C.VkBool32 {
	if b {
		return C.VK_TRUE
	}
	return C.VK_FALSE
}

// Helper function to convert VkBool32 to Go bool
func vkBool32ToBool(b C.VkBool32) bool {
	return b == C.VK_TRUE
}
