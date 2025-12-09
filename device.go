package vulkan

/*
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
	// Input validation
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	// Validate queue create infos
	const maxQueues = 16 // Reasonable limit for queue families
	if len(createInfo.QueueCreateInfos) > maxQueues {
		return nil, NewValidationError("QueueCreateInfos", "exceeds maximum of 16 queue families")
	}
	for i, qci := range createInfo.QueueCreateInfos {
		if len(qci.QueuePriorities) == 0 {
			return nil, NewValidationError("QueueCreateInfos", "queue family must have at least one queue")
		}
		const maxQueuesPerFamily = 16
		if len(qci.QueuePriorities) > maxQueuesPerFamily {
			return nil, NewValidationError("QueueCreateInfos", "queue family exceeds maximum of 16 queues")
		}
		// Validate queue priorities are in range [0.0, 1.0]
		for j, priority := range qci.QueuePriorities {
			if priority < 0.0 || priority > 1.0 {
				return nil, NewValidationError("QueueCreateInfos", "queue priority must be between 0.0 and 1.0")
			}
			_ = j // avoid unused variable
		}
		_ = i // avoid unused variable
	}

	// Validate layers (reuse same validation as CreateInstance)
	const maxLayers = 64
	if len(createInfo.EnabledLayerNames) > maxLayers {
		return nil, NewValidationError("EnabledLayerNames", "exceeds maximum of 64 layers")
	}
	for _, layer := range createInfo.EnabledLayerNames {
		if len(layer) > 256 {
			return nil, NewValidationError("EnabledLayerNames", "layer name exceeds maximum length of 256 characters")
		}
	}

	// Validate extensions
	const maxExtensions = 256
	if len(createInfo.EnabledExtensionNames) > maxExtensions {
		return nil, NewValidationError("EnabledExtensionNames", "exceeds maximum of 256 extensions")
	}
	for _, ext := range createInfo.EnabledExtensionNames {
		if len(ext) > 256 {
			return nil, NewValidationError("EnabledExtensionNames", "extension name exceeds maximum length of 256 characters")
		}
	}

	// Allocate create info in C memory to avoid Go pointer issues
	cCreateInfoPtr := (*C.VkDeviceCreateInfo)(C.malloc(C.sizeof_VkDeviceCreateInfo))
	if cCreateInfoPtr == nil {
		return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateDevice", "failed to allocate memory for device create info")
	}
	defer C.free(unsafe.Pointer(cCreateInfoPtr))

	// Zero-initialize the entire structure
	C.memset(unsafe.Pointer(cCreateInfoPtr), 0, C.sizeof_VkDeviceCreateInfo)

	cCreateInfoPtr.sType = C.VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO
	cCreateInfoPtr.pNext = nil
	cCreateInfoPtr.flags = 0

	// Queue create infos - allocate in C memory
	var cQueueCreateInfosPtr *C.VkDeviceQueueCreateInfo
	var cPrioritiesArray []*C.float
	var cPrioritiesToFree []*C.float // Track allocated priorities for cleanup
	
	if len(createInfo.QueueCreateInfos) > 0 {
		cQueueCreateInfosPtr = (*C.VkDeviceQueueCreateInfo)(C.malloc(C.size_t(len(createInfo.QueueCreateInfos)) * C.sizeof_VkDeviceQueueCreateInfo))
		if cQueueCreateInfosPtr == nil {
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateDevice", "failed to allocate memory for queue create infos")
		}
		defer C.free(unsafe.Pointer(cQueueCreateInfosPtr))

		// Zero-initialize the queue create info structures
		C.memset(unsafe.Pointer(cQueueCreateInfosPtr), 0, C.size_t(len(createInfo.QueueCreateInfos))*C.sizeof_VkDeviceQueueCreateInfo)

		cPrioritiesArray = make([]*C.float, len(createInfo.QueueCreateInfos))

		for i, qci := range createInfo.QueueCreateInfos {
			// Use pointer arithmetic to access array elements (in bytes)
			offset := uintptr(i) * uintptr(C.sizeof_VkDeviceQueueCreateInfo)
			cQueueInfo := (*C.VkDeviceQueueCreateInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cQueueCreateInfosPtr)) + offset))
			cQueueInfo.sType = C.VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO
			cQueueInfo.pNext = nil
			cQueueInfo.flags = 0
			cQueueInfo.queueFamilyIndex = C.uint32_t(qci.QueueFamilyIndex)
			cQueueInfo.queueCount = C.uint32_t(len(qci.QueuePriorities))

			if len(qci.QueuePriorities) > 0 {
				cPrioritiesPtr := (*C.float)(C.malloc(C.size_t(len(qci.QueuePriorities)) * C.sizeof_float))
				if cPrioritiesPtr == nil {
					// Clean up allocated priorities before returning
					for _, ptr := range cPrioritiesToFree {
						C.free(unsafe.Pointer(ptr))
					}
					return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateDevice", "failed to allocate memory for queue priorities")
				}
				// Zero-initialize the priorities array
				C.memset(unsafe.Pointer(cPrioritiesPtr), 0, C.size_t(len(qci.QueuePriorities))*C.sizeof_float)
				cPrioritiesToFree = append(cPrioritiesToFree, cPrioritiesPtr)
				cPrioritiesArray[i] = cPrioritiesPtr

				for j, priority := range qci.QueuePriorities {
					cPriority := (*C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(cPrioritiesPtr)) + uintptr(j)*uintptr(C.sizeof_float)))
					*cPriority = C.float(priority)
				}
				cQueueInfo.pQueuePriorities = cPrioritiesPtr
			}
		}
		cCreateInfoPtr.queueCreateInfoCount = C.uint32_t(len(createInfo.QueueCreateInfos))
		cCreateInfoPtr.pQueueCreateInfos = cQueueCreateInfosPtr
	}
	
	// Defer cleanup of priority arrays
	defer func() {
		for _, ptr := range cPrioritiesToFree {
			C.free(unsafe.Pointer(ptr))
		}
	}()

	// Enabled layers
	var cLayers **C.char
	if len(createInfo.EnabledLayerNames) > 0 {
		cLayers = stringSliceToCharArray(createInfo.EnabledLayerNames)
		if cLayers == nil {
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateDevice", "failed to allocate memory for layer names")
		}
		defer freeStringArray(cLayers, len(createInfo.EnabledLayerNames))
		cCreateInfoPtr.enabledLayerCount = C.uint32_t(len(createInfo.EnabledLayerNames))
		cCreateInfoPtr.ppEnabledLayerNames = cLayers
	}

	// Enabled extensions
	var cExtensions **C.char
	if len(createInfo.EnabledExtensionNames) > 0 {
		cExtensions = stringSliceToCharArray(createInfo.EnabledExtensionNames)
		if cExtensions == nil {
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateDevice", "failed to allocate memory for extension names")
		}
		defer freeStringArray(cExtensions, len(createInfo.EnabledExtensionNames))
		cCreateInfoPtr.enabledExtensionCount = C.uint32_t(len(createInfo.EnabledExtensionNames))
		cCreateInfoPtr.ppEnabledExtensionNames = cExtensions
	}

	// Enabled features - allocate in C memory
	var cFeaturesPtr *C.VkPhysicalDeviceFeatures
	if createInfo.EnabledFeatures != nil {
		cFeaturesPtr = (*C.VkPhysicalDeviceFeatures)(C.malloc(C.sizeof_VkPhysicalDeviceFeatures))
		if cFeaturesPtr == nil {
			// Clean up priorities before returning
			for _, ptr := range cPrioritiesToFree {
				C.free(unsafe.Pointer(ptr))
			}
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateDevice", "failed to allocate memory for physical device features")
		}
		*cFeaturesPtr = physicalDeviceFeaturesToC(createInfo.EnabledFeatures)
		cCreateInfoPtr.pEnabledFeatures = cFeaturesPtr
		
		// Defer cleanup of features
		defer C.free(unsafe.Pointer(cFeaturesPtr))
	}

	var device C.VkDevice
	result := Result(C.vkCreateDevice(C.VkPhysicalDevice(physicalDevice), cCreateInfoPtr, nil, &device))
	
	if result != Success {
		return nil, NewVulkanError(result, "CreateDevice", "Vulkan device creation failed")
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
