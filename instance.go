package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>

// Helper function to convert Go string slice to C char**
char** makeCharArray(int size) {
    char** result = calloc(size, sizeof(char*));
    return result; // Returns NULL if allocation fails
}

// Helper function to set string in char array
void setArrayString(char **a, char *s, int n) {
    a[n] = s;
}

// Helper function to free char array
void freeCharArray(char **a, int size) {
    if (a == NULL) {
        return; // Safely handle NULL pointer
    }
    for (int i = 0; i < size; i++) {
        free(a[i]);
    }
    free(a);
}
*/
import "C"

import (
	"unsafe"
)

// ApplicationInfo contains application information
type ApplicationInfo struct {
	ApplicationName    string
	ApplicationVersion Version
	EngineName         string
	EngineVersion      Version
	APIVersion         Version
}

// InstanceCreateInfo contains instance creation information
type InstanceCreateInfo struct {
	ApplicationInfo       *ApplicationInfo
	EnabledLayerNames     []string
	EnabledExtensionNames []string
}

// ExtensionProperties contains extension information
type ExtensionProperties struct {
	ExtensionName string
	SpecVersion   uint32
}

// LayerProperties contains layer information
type LayerProperties struct {
	LayerName             string
	SpecVersion           Version
	ImplementationVersion Version
	Description           string
}

// PhysicalDeviceType represents the type of physical device
type PhysicalDeviceType int32

const (
	PhysicalDeviceTypeOther         PhysicalDeviceType = C.VK_PHYSICAL_DEVICE_TYPE_OTHER
	PhysicalDeviceTypeIntegratedGPU PhysicalDeviceType = C.VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	PhysicalDeviceTypeDiscreteGPU   PhysicalDeviceType = C.VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	PhysicalDeviceTypeVirtualGPU    PhysicalDeviceType = C.VK_PHYSICAL_DEVICE_TYPE_VIRTUAL_GPU
	PhysicalDeviceTypeCPU           PhysicalDeviceType = C.VK_PHYSICAL_DEVICE_TYPE_CPU
)

// PhysicalDeviceProperties contains physical device properties
type PhysicalDeviceProperties struct {
	APIVersion        Version
	DriverVersion     Version
	VendorID          uint32
	DeviceID          uint32
	DeviceType        PhysicalDeviceType
	DeviceName        string
	PipelineCacheUUID [UuidSize]uint8
	Limits            PhysicalDeviceLimits
	SparseProperties  PhysicalDeviceSparseProperties
}

// PhysicalDeviceLimits contains physical device limits
type PhysicalDeviceLimits struct {
	MaxImageDimension1D                             uint32
	MaxImageDimension2D                             uint32
	MaxImageDimension3D                             uint32
	MaxImageDimensionCube                           uint32
	MaxImageArrayLayers                             uint32
	MaxTexelBufferElements                          uint32
	MaxUniformBufferRange                           uint32
	MaxStorageBufferRange                           uint32
	MaxPushConstantsSize                            uint32
	MaxMemoryAllocationCount                        uint32
	MaxSamplerAllocationCount                       uint32
	BufferImageGranularity                          DeviceSize
	SparseAddressSpaceSize                          DeviceSize
	MaxBoundDescriptorSets                          uint32
	MaxPerStageDescriptorSamplers                   uint32
	MaxPerStageDescriptorUniformBuffers             uint32
	MaxPerStageDescriptorStorageBuffers             uint32
	MaxPerStageDescriptorSampledImages              uint32
	MaxPerStageDescriptorStorageImages              uint32
	MaxPerStageDescriptorInputAttachments           uint32
	MaxPerStageResources                            uint32
	MaxDescriptorSetSamplers                        uint32
	MaxDescriptorSetUniformBuffers                  uint32
	MaxDescriptorSetUniformBuffersDynamic           uint32
	MaxDescriptorSetStorageBuffers                  uint32
	MaxDescriptorSetStorageBuffersDynamic           uint32
	MaxDescriptorSetSampledImages                   uint32
	MaxDescriptorSetStorageImages                   uint32
	MaxDescriptorSetInputAttachments                uint32
	MaxVertexInputAttributes                        uint32
	MaxVertexInputBindings                          uint32
	MaxVertexInputAttributeOffset                   uint32
	MaxVertexInputBindingStride                     uint32
	MaxVertexOutputComponents                       uint32
	MaxTessellationGenerationLevel                  uint32
	MaxTessellationPatchSize                        uint32
	MaxTessellationControlPerVertexInputComponents  uint32
	MaxTessellationControlPerVertexOutputComponents uint32
	MaxTessellationControlPerPatchOutputComponents  uint32
	MaxTessellationControlTotalOutputComponents     uint32
	MaxTessellationEvaluationInputComponents        uint32
	MaxTessellationEvaluationOutputComponents       uint32
	MaxGeometryShaderInvocations                    uint32
	MaxGeometryInputComponents                      uint32
	MaxGeometryOutputComponents                     uint32
	MaxGeometryOutputVertices                       uint32
	MaxGeometryTotalOutputComponents                uint32
	MaxFragmentInputComponents                      uint32
	MaxFragmentOutputAttachments                    uint32
	MaxFragmentDualSrcAttachments                   uint32
	MaxFragmentCombinedOutputResources              uint32
	MaxComputeSharedMemorySize                      uint32
	MaxComputeWorkGroupCount                        [3]uint32
	MaxComputeWorkGroupInvocations                  uint32
	MaxComputeWorkGroupSize                         [3]uint32
	SubPixelPrecisionBits                           uint32
	SubTexelPrecisionBits                           uint32
	MipmapPrecisionBits                             uint32
	MaxDrawIndexedIndexValue                        uint32
	MaxDrawIndirectCount                            uint32
	MaxSamplerLodBias                               float32
	MaxSamplerAnisotropy                            float32
	MaxViewports                                    uint32
	MaxViewportDimensions                           [2]uint32
	ViewportBoundsRange                             [2]float32
	ViewportSubPixelBits                            uint32
	MinMemoryMapAlignment                           uintptr
	MinTexelBufferOffsetAlignment                   DeviceSize
	MinUniformBufferOffsetAlignment                 DeviceSize
	MinStorageBufferOffsetAlignment                 DeviceSize
	MinTexelOffset                                  int32
	MaxTexelOffset                                  uint32
	MinTexelGatherOffset                            int32
	MaxTexelGatherOffset                            uint32
	MinInterpolationOffset                          float32
	MaxInterpolationOffset                          float32
	SubPixelInterpolationOffsetBits                 uint32
	MaxFramebufferWidth                             uint32
	MaxFramebufferHeight                            uint32
	MaxFramebufferLayers                            uint32
	FramebufferColorSampleCounts                    SampleCountFlags
	FramebufferDepthSampleCounts                    SampleCountFlags
	FramebufferStencilSampleCounts                  SampleCountFlags
	FramebufferNoAttachmentsSampleCounts            SampleCountFlags
	MaxColorAttachments                             uint32
	SampledImageColorSampleCounts                   SampleCountFlags
	SampledImageIntegerSampleCounts                 SampleCountFlags
	SampledImageDepthSampleCounts                   SampleCountFlags
	SampledImageStencilSampleCounts                 SampleCountFlags
	StorageImageSampleCounts                        SampleCountFlags
	MaxSampleMaskWords                              uint32
	TimestampComputeAndGraphics                     Bool32
	TimestampPeriod                                 float32
	MaxClipDistances                                uint32
	MaxCullDistances                                uint32
	MaxCombinedClipAndCullDistances                 uint32
	DiscreteQueuePriorities                         uint32
	PointSizeRange                                  [2]float32
	LineWidthRange                                  [2]float32
	PointSizeGranularity                            float32
	LineWidthGranularity                            float32
	StrictLines                                     Bool32
	StandardSampleLocations                         Bool32
	OptimalBufferCopyOffsetAlignment                DeviceSize
	OptimalBufferCopyRowPitchAlignment              DeviceSize
	NonCoherentAtomSize                             DeviceSize
}

// PhysicalDeviceSparseProperties contains sparse resource properties
type PhysicalDeviceSparseProperties struct {
	ResidencyStandard2DBlockShape            Bool32
	ResidencyStandard2DMultisampleBlockShape Bool32
	ResidencyStandard3DBlockShape            Bool32
	ResidencyAlignedMipSize                  Bool32
	ResidencyNonResidentStrict               Bool32
}

// QueueFamilyProperties contains queue family properties
type QueueFamilyProperties struct {
	QueueFlags                  QueueFlags
	QueueCount                  uint32
	TimestampValidBits          uint32
	MinImageTransferGranularity Extent3D
}

// QueueFlags represents queue capability flags
type QueueFlags uint32

const (
	QueueGraphicsBit       QueueFlags = C.VK_QUEUE_GRAPHICS_BIT
	QueueComputeBit        QueueFlags = C.VK_QUEUE_COMPUTE_BIT
	QueueTransferBit       QueueFlags = C.VK_QUEUE_TRANSFER_BIT
	QueueSparseBindingBit  QueueFlags = C.VK_QUEUE_SPARSE_BINDING_BIT
	QueueProtectedBit      QueueFlags = C.VK_QUEUE_PROTECTED_BIT
	QueueVideoDecodeBitKHR QueueFlags = C.VK_QUEUE_VIDEO_DECODE_BIT_KHR
	QueueVideoEncodeBitKHR QueueFlags = C.VK_QUEUE_VIDEO_ENCODE_BIT_KHR
)

// Extent3D represents a 3D extent
type Extent3D struct {
	Width  uint32
	Height uint32
	Depth  uint32
}

// stringSliceToCharArray converts Go string slice to C char**
func stringSliceToCharArray(strs []string) **C.char {
	if len(strs) == 0 {
		return nil
	}

	// Add bounds checking for large allocations
	const maxAllowedSize = 10000 // Reasonable limit for string arrays
	if len(strs) > maxAllowedSize {
		return nil
	}

	cArray := C.makeCharArray(C.int(len(strs)))
	// Check if allocation failed
	if cArray == nil {
		return nil
	}

	// Convert strings and track successful allocations for cleanup on failure
	for i, str := range strs {
		cStr := C.CString(str)
		// Check if string allocation failed
		if cStr == nil {
			// Clean up the array - freeCharArray handles partial cleanup
			C.freeCharArray(cArray, C.int(i)) // Only free up to current index
			return nil
		}
		C.setArrayString(cArray, cStr, C.int(i))
	}
	return cArray
}

// freeCharArray frees a C char** array
func freeStringArray(cArray **C.char, size int) {
	if cArray != nil {
		C.freeCharArray(cArray, C.int(size))
	}
}

// CreateInstance creates a Vulkan instance
func CreateInstance(createInfo *InstanceCreateInfo) (Instance, error) {
	// Input validation
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	// Validate application info if provided
	if createInfo.ApplicationInfo != nil {
		if len(createInfo.ApplicationInfo.ApplicationName) > 256 {
			return nil, NewValidationError("ApplicationInfo.ApplicationName", "exceeds maximum length of 256 characters")
		}
		if len(createInfo.ApplicationInfo.EngineName) > 256 {
			return nil, NewValidationError("ApplicationInfo.EngineName", "exceeds maximum length of 256 characters")
		}
	}

	// Validate layer names
	const maxLayers = 64 // Reasonable limit
	if len(createInfo.EnabledLayerNames) > maxLayers {
		return nil, NewValidationError("EnabledLayerNames", "exceeds maximum of 64 layers")
	}
	for _, layer := range createInfo.EnabledLayerNames {
		if len(layer) > 256 {
			return nil, NewValidationError("EnabledLayerNames", "layer name at index exceeds maximum length of 256 characters")
		}
	}

	// Validate extension names
	const maxExtensions = 256 // Reasonable limit
	if len(createInfo.EnabledExtensionNames) > maxExtensions {
		return nil, NewValidationError("EnabledExtensionNames", "exceeds maximum of 256 extensions")
	}
	for _, ext := range createInfo.EnabledExtensionNames {
		if len(ext) > 256 {
			return nil, NewValidationError("EnabledExtensionNames", "extension name at index exceeds maximum length of 256 characters")
		}
	}

	var cCreateInfo C.VkInstanceCreateInfo
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0

	// Application info - allocate on heap to avoid Go pointer issues
	var cAppInfo *C.VkApplicationInfo
	var appNamePtr, engineNamePtr *C.char
	if createInfo.ApplicationInfo != nil {
		cAppInfo = (*C.VkApplicationInfo)(C.malloc(C.size_t(unsafe.Sizeof(C.VkApplicationInfo{}))))
		if cAppInfo == nil {
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateInstance", "failed to allocate memory for application info")
		}
		cAppInfo.sType = C.VK_STRUCTURE_TYPE_APPLICATION_INFO
		cAppInfo.pNext = nil
		cAppInfo.pApplicationName = nil
		cAppInfo.pEngineName = nil

		if createInfo.ApplicationInfo.ApplicationName != "" {
			appNamePtr = C.CString(createInfo.ApplicationInfo.ApplicationName)
			if appNamePtr == nil {
				C.free(unsafe.Pointer(cAppInfo))
				return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateInstance", "failed to allocate memory for application name")
			}
			cAppInfo.pApplicationName = appNamePtr
		}
		cAppInfo.applicationVersion = C.uint32_t(createInfo.ApplicationInfo.ApplicationVersion)

		if createInfo.ApplicationInfo.EngineName != "" {
			engineNamePtr = C.CString(createInfo.ApplicationInfo.EngineName)
			if engineNamePtr == nil {
				if appNamePtr != nil {
					C.free(unsafe.Pointer(appNamePtr))
				}
				C.free(unsafe.Pointer(cAppInfo))
				return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateInstance", "failed to allocate memory for engine name")
			}
			cAppInfo.pEngineName = engineNamePtr
		}
		cAppInfo.engineVersion = C.uint32_t(createInfo.ApplicationInfo.EngineVersion)
		cAppInfo.apiVersion = C.uint32_t(createInfo.ApplicationInfo.APIVersion)

		cCreateInfo.pApplicationInfo = cAppInfo
	}

	// Enabled layers
	var cLayers **C.char
	if len(createInfo.EnabledLayerNames) > 0 {
		cLayers = stringSliceToCharArray(createInfo.EnabledLayerNames)
		if cLayers == nil {
			// Clean up already allocated memory
			if appNamePtr != nil {
				C.free(unsafe.Pointer(appNamePtr))
			}
			if engineNamePtr != nil {
				C.free(unsafe.Pointer(engineNamePtr))
			}
			if cAppInfo != nil {
				C.free(unsafe.Pointer(cAppInfo))
			}
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateInstance", "failed to allocate memory for layer names")
		}
		cCreateInfo.enabledLayerCount = C.uint32_t(len(createInfo.EnabledLayerNames))
		cCreateInfo.ppEnabledLayerNames = cLayers
	}

	// Enabled extensions
	var cExtensions **C.char
	if len(createInfo.EnabledExtensionNames) > 0 {
		cExtensions = stringSliceToCharArray(createInfo.EnabledExtensionNames)
		if cExtensions == nil {
			// Clean up already allocated memory
			if appNamePtr != nil {
				C.free(unsafe.Pointer(appNamePtr))
			}
			if engineNamePtr != nil {
				C.free(unsafe.Pointer(engineNamePtr))
			}
			if cAppInfo != nil {
				C.free(unsafe.Pointer(cAppInfo))
			}
			if cLayers != nil {
				freeStringArray(cLayers, len(createInfo.EnabledLayerNames))
			}
			return nil, NewVulkanError(ErrorOutOfHostMemory, "CreateInstance", "failed to allocate memory for extension names")
		}
		cCreateInfo.enabledExtensionCount = C.uint32_t(len(createInfo.EnabledExtensionNames))
		cCreateInfo.ppEnabledExtensionNames = cExtensions
	}

	var instance C.VkInstance
	result := Result(C.vkCreateInstance(&cCreateInfo, nil, &instance))

	// Clean up memory
	if appNamePtr != nil {
		C.free(unsafe.Pointer(appNamePtr))
	}
	if engineNamePtr != nil {
		C.free(unsafe.Pointer(engineNamePtr))
	}
	if cAppInfo != nil {
		C.free(unsafe.Pointer(cAppInfo))
	}
	if cLayers != nil {
		freeStringArray(cLayers, len(createInfo.EnabledLayerNames))
	}
	if cExtensions != nil {
		freeStringArray(cExtensions, len(createInfo.EnabledExtensionNames))
	}

	if result != Success {
		return nil, NewVulkanError(result, "CreateInstance", "Vulkan instance creation failed")
	}

	return Instance(instance), nil
}

// DestroyInstance destroys a Vulkan instance
func DestroyInstance(instance Instance) {
	C.vkDestroyInstance(C.VkInstance(instance), nil)
}

// EnumerateInstanceExtensionProperties enumerates available instance extensions
func EnumerateInstanceExtensionProperties(layerName string) ([]ExtensionProperties, error) {
	var cLayerName *C.char
	if layerName != "" {
		cLayerName = C.CString(layerName)
		defer C.free(unsafe.Pointer(cLayerName))
	}

	var propertyCount C.uint32_t
	result := Result(C.vkEnumerateInstanceExtensionProperties(cLayerName, &propertyCount, nil))
	if result != Success {
		return nil, result
	}

	if propertyCount == 0 {
		return nil, nil
	}

	cProperties := make([]C.VkExtensionProperties, propertyCount)
	result = Result(C.vkEnumerateInstanceExtensionProperties(cLayerName, &propertyCount, &cProperties[0]))
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

// EnumerateInstanceLayerProperties enumerates available instance layers
func EnumerateInstanceLayerProperties() ([]LayerProperties, error) {
	var propertyCount C.uint32_t
	result := Result(C.vkEnumerateInstanceLayerProperties(&propertyCount, nil))
	if result != Success {
		return nil, result
	}

	if propertyCount == 0 {
		return nil, nil
	}

	cProperties := make([]C.VkLayerProperties, propertyCount)
	result = Result(C.vkEnumerateInstanceLayerProperties(&propertyCount, &cProperties[0]))
	if result != Success {
		return nil, result
	}

	properties := make([]LayerProperties, propertyCount)
	for i := range properties {
		properties[i].LayerName = C.GoString(&cProperties[i].layerName[0])
		properties[i].SpecVersion = Version(cProperties[i].specVersion)
		properties[i].ImplementationVersion = Version(cProperties[i].implementationVersion)
		properties[i].Description = C.GoString(&cProperties[i].description[0])
	}

	return properties, nil
}

// EnumeratePhysicalDevices enumerates physical devices
func EnumeratePhysicalDevices(instance Instance) ([]PhysicalDevice, error) {
	var deviceCount C.uint32_t
	result := Result(C.vkEnumeratePhysicalDevices(C.VkInstance(instance), &deviceCount, nil))
	if result != Success {
		return nil, result
	}

	if deviceCount == 0 {
		return nil, nil
	}

	cDevices := make([]C.VkPhysicalDevice, deviceCount)
	result = Result(C.vkEnumeratePhysicalDevices(C.VkInstance(instance), &deviceCount, &cDevices[0]))
	if result != Success {
		return nil, result
	}

	devices := make([]PhysicalDevice, deviceCount)
	for i := range devices {
		devices[i] = PhysicalDevice(cDevices[i])
	}

	return devices, nil
}

// GetPhysicalDeviceProperties gets physical device properties
func GetPhysicalDeviceProperties(physicalDevice PhysicalDevice) PhysicalDeviceProperties {
	var cProperties C.VkPhysicalDeviceProperties
	C.vkGetPhysicalDeviceProperties(C.VkPhysicalDevice(physicalDevice), &cProperties)

	properties := PhysicalDeviceProperties{
		APIVersion:    Version(cProperties.apiVersion),
		DriverVersion: Version(cProperties.driverVersion),
		VendorID:      uint32(cProperties.vendorID),
		DeviceID:      uint32(cProperties.deviceID),
		DeviceType:    PhysicalDeviceType(cProperties.deviceType),
		DeviceName:    C.GoString(&cProperties.deviceName[0]),
	}

	// Copy UUID
	for i := 0; i < UuidSize; i++ {
		properties.PipelineCacheUUID[i] = uint8(cProperties.pipelineCacheUUID[i])
	}

	// Convert limits
	properties.Limits = PhysicalDeviceLimits{
		MaxImageDimension1D:                             uint32(cProperties.limits.maxImageDimension1D),
		MaxImageDimension2D:                             uint32(cProperties.limits.maxImageDimension2D),
		MaxImageDimension3D:                             uint32(cProperties.limits.maxImageDimension3D),
		MaxImageDimensionCube:                           uint32(cProperties.limits.maxImageDimensionCube),
		MaxImageArrayLayers:                             uint32(cProperties.limits.maxImageArrayLayers),
		MaxTexelBufferElements:                          uint32(cProperties.limits.maxTexelBufferElements),
		MaxUniformBufferRange:                           uint32(cProperties.limits.maxUniformBufferRange),
		MaxStorageBufferRange:                           uint32(cProperties.limits.maxStorageBufferRange),
		MaxPushConstantsSize:                            uint32(cProperties.limits.maxPushConstantsSize),
		MaxMemoryAllocationCount:                        uint32(cProperties.limits.maxMemoryAllocationCount),
		MaxSamplerAllocationCount:                       uint32(cProperties.limits.maxSamplerAllocationCount),
		BufferImageGranularity:                          DeviceSize(cProperties.limits.bufferImageGranularity),
		SparseAddressSpaceSize:                          DeviceSize(cProperties.limits.sparseAddressSpaceSize),
		MaxBoundDescriptorSets:                          uint32(cProperties.limits.maxBoundDescriptorSets),
		MaxPerStageDescriptorSamplers:                   uint32(cProperties.limits.maxPerStageDescriptorSamplers),
		MaxPerStageDescriptorUniformBuffers:             uint32(cProperties.limits.maxPerStageDescriptorUniformBuffers),
		MaxPerStageDescriptorStorageBuffers:             uint32(cProperties.limits.maxPerStageDescriptorStorageBuffers),
		MaxPerStageDescriptorSampledImages:              uint32(cProperties.limits.maxPerStageDescriptorSampledImages),
		MaxPerStageDescriptorStorageImages:              uint32(cProperties.limits.maxPerStageDescriptorStorageImages),
		MaxPerStageDescriptorInputAttachments:           uint32(cProperties.limits.maxPerStageDescriptorInputAttachments),
		MaxPerStageResources:                            uint32(cProperties.limits.maxPerStageResources),
		MaxDescriptorSetSamplers:                        uint32(cProperties.limits.maxDescriptorSetSamplers),
		MaxDescriptorSetUniformBuffers:                  uint32(cProperties.limits.maxDescriptorSetUniformBuffers),
		MaxDescriptorSetUniformBuffersDynamic:           uint32(cProperties.limits.maxDescriptorSetUniformBuffersDynamic),
		MaxDescriptorSetStorageBuffers:                  uint32(cProperties.limits.maxDescriptorSetStorageBuffers),
		MaxDescriptorSetStorageBuffersDynamic:           uint32(cProperties.limits.maxDescriptorSetStorageBuffersDynamic),
		MaxDescriptorSetSampledImages:                   uint32(cProperties.limits.maxDescriptorSetSampledImages),
		MaxDescriptorSetStorageImages:                   uint32(cProperties.limits.maxDescriptorSetStorageImages),
		MaxDescriptorSetInputAttachments:                uint32(cProperties.limits.maxDescriptorSetInputAttachments),
		MaxVertexInputAttributes:                        uint32(cProperties.limits.maxVertexInputAttributes),
		MaxVertexInputBindings:                          uint32(cProperties.limits.maxVertexInputBindings),
		MaxVertexInputAttributeOffset:                   uint32(cProperties.limits.maxVertexInputAttributeOffset),
		MaxVertexInputBindingStride:                     uint32(cProperties.limits.maxVertexInputBindingStride),
		MaxVertexOutputComponents:                       uint32(cProperties.limits.maxVertexOutputComponents),
		MaxTessellationGenerationLevel:                  uint32(cProperties.limits.maxTessellationGenerationLevel),
		MaxTessellationPatchSize:                        uint32(cProperties.limits.maxTessellationPatchSize),
		MaxTessellationControlPerVertexInputComponents:  uint32(cProperties.limits.maxTessellationControlPerVertexInputComponents),
		MaxTessellationControlPerVertexOutputComponents: uint32(cProperties.limits.maxTessellationControlPerVertexOutputComponents),
		MaxTessellationControlPerPatchOutputComponents:  uint32(cProperties.limits.maxTessellationControlPerPatchOutputComponents),
		MaxTessellationControlTotalOutputComponents:     uint32(cProperties.limits.maxTessellationControlTotalOutputComponents),
		MaxTessellationEvaluationInputComponents:        uint32(cProperties.limits.maxTessellationEvaluationInputComponents),
		MaxTessellationEvaluationOutputComponents:       uint32(cProperties.limits.maxTessellationEvaluationOutputComponents),
		MaxGeometryShaderInvocations:                    uint32(cProperties.limits.maxGeometryShaderInvocations),
		MaxGeometryInputComponents:                      uint32(cProperties.limits.maxGeometryInputComponents),
		MaxGeometryOutputComponents:                     uint32(cProperties.limits.maxGeometryOutputComponents),
		MaxGeometryOutputVertices:                       uint32(cProperties.limits.maxGeometryOutputVertices),
		MaxGeometryTotalOutputComponents:                uint32(cProperties.limits.maxGeometryTotalOutputComponents),
		MaxFragmentInputComponents:                      uint32(cProperties.limits.maxFragmentInputComponents),
		MaxFragmentOutputAttachments:                    uint32(cProperties.limits.maxFragmentOutputAttachments),
		MaxFragmentDualSrcAttachments:                   uint32(cProperties.limits.maxFragmentDualSrcAttachments),
		MaxFragmentCombinedOutputResources:              uint32(cProperties.limits.maxFragmentCombinedOutputResources),
		MaxComputeSharedMemorySize:                      uint32(cProperties.limits.maxComputeSharedMemorySize),
		MaxComputeWorkGroupInvocations:                  uint32(cProperties.limits.maxComputeWorkGroupInvocations),
		SubPixelPrecisionBits:                           uint32(cProperties.limits.subPixelPrecisionBits),
		SubTexelPrecisionBits:                           uint32(cProperties.limits.subTexelPrecisionBits),
		MipmapPrecisionBits:                             uint32(cProperties.limits.mipmapPrecisionBits),
		MaxDrawIndexedIndexValue:                        uint32(cProperties.limits.maxDrawIndexedIndexValue),
		MaxDrawIndirectCount:                            uint32(cProperties.limits.maxDrawIndirectCount),
		MaxSamplerLodBias:                               float32(cProperties.limits.maxSamplerLodBias),
		MaxSamplerAnisotropy:                            float32(cProperties.limits.maxSamplerAnisotropy),
		MaxViewports:                                    uint32(cProperties.limits.maxViewports),
		ViewportSubPixelBits:                            uint32(cProperties.limits.viewportSubPixelBits),
		MinMemoryMapAlignment:                           uintptr(cProperties.limits.minMemoryMapAlignment),
		MinTexelBufferOffsetAlignment:                   DeviceSize(cProperties.limits.minTexelBufferOffsetAlignment),
		MinUniformBufferOffsetAlignment:                 DeviceSize(cProperties.limits.minUniformBufferOffsetAlignment),
		MinStorageBufferOffsetAlignment:                 DeviceSize(cProperties.limits.minStorageBufferOffsetAlignment),
		MinTexelOffset:                                  int32(cProperties.limits.minTexelOffset),
		MaxTexelOffset:                                  uint32(cProperties.limits.maxTexelOffset),
		MinTexelGatherOffset:                            int32(cProperties.limits.minTexelGatherOffset),
		MaxTexelGatherOffset:                            uint32(cProperties.limits.maxTexelGatherOffset),
		MinInterpolationOffset:                          float32(cProperties.limits.minInterpolationOffset),
		MaxInterpolationOffset:                          float32(cProperties.limits.maxInterpolationOffset),
		SubPixelInterpolationOffsetBits:                 uint32(cProperties.limits.subPixelInterpolationOffsetBits),
		MaxFramebufferWidth:                             uint32(cProperties.limits.maxFramebufferWidth),
		MaxFramebufferHeight:                            uint32(cProperties.limits.maxFramebufferHeight),
		MaxFramebufferLayers:                            uint32(cProperties.limits.maxFramebufferLayers),
		FramebufferColorSampleCounts:                    SampleCountFlags(cProperties.limits.framebufferColorSampleCounts),
		FramebufferDepthSampleCounts:                    SampleCountFlags(cProperties.limits.framebufferDepthSampleCounts),
		FramebufferStencilSampleCounts:                  SampleCountFlags(cProperties.limits.framebufferStencilSampleCounts),
		FramebufferNoAttachmentsSampleCounts:            SampleCountFlags(cProperties.limits.framebufferNoAttachmentsSampleCounts),
		MaxColorAttachments:                             uint32(cProperties.limits.maxColorAttachments),
		SampledImageColorSampleCounts:                   SampleCountFlags(cProperties.limits.sampledImageColorSampleCounts),
		SampledImageIntegerSampleCounts:                 SampleCountFlags(cProperties.limits.sampledImageIntegerSampleCounts),
		SampledImageDepthSampleCounts:                   SampleCountFlags(cProperties.limits.sampledImageDepthSampleCounts),
		SampledImageStencilSampleCounts:                 SampleCountFlags(cProperties.limits.sampledImageStencilSampleCounts),
		StorageImageSampleCounts:                        SampleCountFlags(cProperties.limits.storageImageSampleCounts),
		MaxSampleMaskWords:                              uint32(cProperties.limits.maxSampleMaskWords),
		TimestampComputeAndGraphics:                     Bool32(cProperties.limits.timestampComputeAndGraphics),
		TimestampPeriod:                                 float32(cProperties.limits.timestampPeriod),
		MaxClipDistances:                                uint32(cProperties.limits.maxClipDistances),
		MaxCullDistances:                                uint32(cProperties.limits.maxCullDistances),
		MaxCombinedClipAndCullDistances:                 uint32(cProperties.limits.maxCombinedClipAndCullDistances),
		DiscreteQueuePriorities:                         uint32(cProperties.limits.discreteQueuePriorities),
		PointSizeGranularity:                            float32(cProperties.limits.pointSizeGranularity),
		LineWidthGranularity:                            float32(cProperties.limits.lineWidthGranularity),
		StrictLines:                                     Bool32(cProperties.limits.strictLines),
		StandardSampleLocations:                         Bool32(cProperties.limits.standardSampleLocations),
		OptimalBufferCopyOffsetAlignment:                DeviceSize(cProperties.limits.optimalBufferCopyOffsetAlignment),
		OptimalBufferCopyRowPitchAlignment:              DeviceSize(cProperties.limits.optimalBufferCopyRowPitchAlignment),
		NonCoherentAtomSize:                             DeviceSize(cProperties.limits.nonCoherentAtomSize),
	}

	// Copy arrays
	for i := 0; i < 3; i++ {
		properties.Limits.MaxComputeWorkGroupCount[i] = uint32(cProperties.limits.maxComputeWorkGroupCount[i])
		properties.Limits.MaxComputeWorkGroupSize[i] = uint32(cProperties.limits.maxComputeWorkGroupSize[i])
	}
	for i := 0; i < 2; i++ {
		properties.Limits.MaxViewportDimensions[i] = uint32(cProperties.limits.maxViewportDimensions[i])
		properties.Limits.ViewportBoundsRange[i] = float32(cProperties.limits.viewportBoundsRange[i])
		properties.Limits.PointSizeRange[i] = float32(cProperties.limits.pointSizeRange[i])
		properties.Limits.LineWidthRange[i] = float32(cProperties.limits.lineWidthRange[i])
	}

	// Sparse properties
	properties.SparseProperties = PhysicalDeviceSparseProperties{
		ResidencyStandard2DBlockShape:            Bool32(cProperties.sparseProperties.residencyStandard2DBlockShape),
		ResidencyStandard2DMultisampleBlockShape: Bool32(cProperties.sparseProperties.residencyStandard2DMultisampleBlockShape),
		ResidencyStandard3DBlockShape:            Bool32(cProperties.sparseProperties.residencyStandard3DBlockShape),
		ResidencyAlignedMipSize:                  Bool32(cProperties.sparseProperties.residencyAlignedMipSize),
		ResidencyNonResidentStrict:               Bool32(cProperties.sparseProperties.residencyNonResidentStrict),
	}

	return properties
}

// GetPhysicalDeviceQueueFamilyProperties gets queue family properties
func GetPhysicalDeviceQueueFamilyProperties(physicalDevice PhysicalDevice) []QueueFamilyProperties {
	var queueFamilyCount C.uint32_t
	C.vkGetPhysicalDeviceQueueFamilyProperties(C.VkPhysicalDevice(physicalDevice), &queueFamilyCount, nil)

	if queueFamilyCount == 0 {
		return nil
	}

	cProperties := make([]C.VkQueueFamilyProperties, queueFamilyCount)
	C.vkGetPhysicalDeviceQueueFamilyProperties(C.VkPhysicalDevice(physicalDevice), &queueFamilyCount, &cProperties[0])

	properties := make([]QueueFamilyProperties, queueFamilyCount)
	for i := range properties {
		properties[i] = QueueFamilyProperties{
			QueueFlags:         QueueFlags(cProperties[i].queueFlags),
			QueueCount:         uint32(cProperties[i].queueCount),
			TimestampValidBits: uint32(cProperties[i].timestampValidBits),
			MinImageTransferGranularity: Extent3D{
				Width:  uint32(cProperties[i].minImageTransferGranularity.width),
				Height: uint32(cProperties[i].minImageTransferGranularity.height),
				Depth:  uint32(cProperties[i].minImageTransferGranularity.depth),
			},
		}
	}

	return properties
}

// ============================================================================
// Instance & Device Enhancements
// ============================================================================

// DebugUtilsMessageSeverityFlags represents debug message severity levels
type DebugUtilsMessageSeverityFlags uint32

const (
	DebugUtilsMessageSeverityVerbose DebugUtilsMessageSeverityFlags = C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_VERBOSE_BIT_EXT
	DebugUtilsMessageSeverityInfo    DebugUtilsMessageSeverityFlags = C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_INFO_BIT_EXT
	DebugUtilsMessageSeverityWarning DebugUtilsMessageSeverityFlags = C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_WARNING_BIT_EXT
	DebugUtilsMessageSeverityError   DebugUtilsMessageSeverityFlags = C.VK_DEBUG_UTILS_MESSAGE_SEVERITY_ERROR_BIT_EXT
)

// DebugUtilsMessageTypeFlags represents debug message types
type DebugUtilsMessageTypeFlags uint32

const (
	DebugUtilsMessageTypeGeneral     DebugUtilsMessageTypeFlags = C.VK_DEBUG_UTILS_MESSAGE_TYPE_GENERAL_BIT_EXT
	DebugUtilsMessageTypeValidation  DebugUtilsMessageTypeFlags = C.VK_DEBUG_UTILS_MESSAGE_TYPE_VALIDATION_BIT_EXT
	DebugUtilsMessageTypePerformance DebugUtilsMessageTypeFlags = C.VK_DEBUG_UTILS_MESSAGE_TYPE_PERFORMANCE_BIT_EXT
)

// DebugUtilsMessengerCallbackData contains data passed to debug callback
type DebugUtilsMessengerCallbackData struct {
	MessageIDName   string
	MessageIDNumber int32
	Message         string
}

// DebugUtilsMessengerCreateInfo contains debug messenger creation information
type DebugUtilsMessengerCreateInfo struct {
	MessageSeverity DebugUtilsMessageSeverityFlags
	MessageType     DebugUtilsMessageTypeFlags
}

// SurfaceCapabilities describes the capabilities of a surface
type SurfaceCapabilities struct {
	MinImageCount           uint32
	MaxImageCount           uint32
	CurrentExtent           Extent2D
	MinImageExtent          Extent2D
	MaxImageExtent          Extent2D
	MaxImageArrayLayers     uint32
	SupportedTransforms     uint32
	CurrentTransform        uint32
	SupportedCompositeAlpha uint32
	SupportedUsageFlags     ImageUsageFlags
}

// SurfaceFormat describes a surface format and color space
type SurfaceFormat struct {
	Format     Format
	ColorSpace uint32
}

// PresentMode represents presentation modes
type PresentMode uint32

const (
	PresentModeImmediate   PresentMode = C.VK_PRESENT_MODE_IMMEDIATE_KHR
	PresentModeMailbox     PresentMode = C.VK_PRESENT_MODE_MAILBOX_KHR
	PresentModeFIFO        PresentMode = C.VK_PRESENT_MODE_FIFO_KHR
	PresentModeFIFORelaxed PresentMode = C.VK_PRESENT_MODE_FIFO_RELAXED_KHR
)

// DestroySurface destroys a surface
func DestroySurface(instance Instance, surface Surface) {
	if instance == nil || surface == nil {
		return
	}
	C.vkDestroySurfaceKHR(C.VkInstance(instance), C.VkSurfaceKHR(surface), nil)
}

// GetPhysicalDeviceSurfaceSupport queries if a queue family supports presentation
func GetPhysicalDeviceSurfaceSupport(physicalDevice PhysicalDevice, queueFamilyIndex uint32, surface Surface) (bool, error) {
	if physicalDevice == nil {
		return false, NewValidationError("physicalDevice", "cannot be nil")
	}
	if surface == nil {
		return false, NewValidationError("surface", "cannot be nil")
	}

	var supported C.VkBool32
	result := Result(C.vkGetPhysicalDeviceSurfaceSupportKHR(
		C.VkPhysicalDevice(physicalDevice),
		C.uint32_t(queueFamilyIndex),
		C.VkSurfaceKHR(surface),
		&supported,
	))
	if result != Success {
		return false, NewVulkanError(result, "GetPhysicalDeviceSurfaceSupport", "Vulkan surface support query failed")
	}

	return supported == C.VK_TRUE, nil
}

// GetPhysicalDeviceSurfaceCapabilities gets surface capabilities
func GetPhysicalDeviceSurfaceCapabilities(physicalDevice PhysicalDevice, surface Surface) (SurfaceCapabilities, error) {
	if physicalDevice == nil {
		return SurfaceCapabilities{}, NewValidationError("physicalDevice", "cannot be nil")
	}
	if surface == nil {
		return SurfaceCapabilities{}, NewValidationError("surface", "cannot be nil")
	}

	var cCaps C.VkSurfaceCapabilitiesKHR
	result := Result(C.vkGetPhysicalDeviceSurfaceCapabilitiesKHR(
		C.VkPhysicalDevice(physicalDevice),
		C.VkSurfaceKHR(surface),
		&cCaps,
	))
	if result != Success {
		return SurfaceCapabilities{}, NewVulkanError(result, "GetPhysicalDeviceSurfaceCapabilities", "Vulkan surface capabilities query failed")
	}

	return SurfaceCapabilities{
		MinImageCount: uint32(cCaps.minImageCount),
		MaxImageCount: uint32(cCaps.maxImageCount),
		CurrentExtent: Extent2D{
			Width:  uint32(cCaps.currentExtent.width),
			Height: uint32(cCaps.currentExtent.height),
		},
		MinImageExtent: Extent2D{
			Width:  uint32(cCaps.minImageExtent.width),
			Height: uint32(cCaps.minImageExtent.height),
		},
		MaxImageExtent: Extent2D{
			Width:  uint32(cCaps.maxImageExtent.width),
			Height: uint32(cCaps.maxImageExtent.height),
		},
		MaxImageArrayLayers:     uint32(cCaps.maxImageArrayLayers),
		SupportedTransforms:     uint32(cCaps.supportedTransforms),
		CurrentTransform:        uint32(cCaps.currentTransform),
		SupportedCompositeAlpha: uint32(cCaps.supportedCompositeAlpha),
		SupportedUsageFlags:     ImageUsageFlags(cCaps.supportedUsageFlags),
	}, nil
}

// GetPhysicalDeviceSurfaceFormats gets surface formats
func GetPhysicalDeviceSurfaceFormats(physicalDevice PhysicalDevice, surface Surface) ([]SurfaceFormat, error) {
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if surface == nil {
		return nil, NewValidationError("surface", "cannot be nil")
	}

	var formatCount C.uint32_t
	result := Result(C.vkGetPhysicalDeviceSurfaceFormatsKHR(
		C.VkPhysicalDevice(physicalDevice),
		C.VkSurfaceKHR(surface),
		&formatCount,
		nil,
	))
	if result != Success {
		return nil, NewVulkanError(result, "GetPhysicalDeviceSurfaceFormats", "Vulkan surface formats query failed")
	}

	if formatCount == 0 {
		return []SurfaceFormat{}, nil
	}

	cFormats := make([]C.VkSurfaceFormatKHR, formatCount)
	result = Result(C.vkGetPhysicalDeviceSurfaceFormatsKHR(
		C.VkPhysicalDevice(physicalDevice),
		C.VkSurfaceKHR(surface),
		&formatCount,
		&cFormats[0],
	))
	if result != Success {
		return nil, NewVulkanError(result, "GetPhysicalDeviceSurfaceFormats", "Vulkan surface formats query failed")
	}

	formats := make([]SurfaceFormat, formatCount)
	for i := range formats {
		formats[i] = SurfaceFormat{
			Format:     Format(cFormats[i].format),
			ColorSpace: uint32(cFormats[i].colorSpace),
		}
	}

	return formats, nil
}

// GetPhysicalDeviceSurfacePresentModes gets surface present modes
func GetPhysicalDeviceSurfacePresentModes(physicalDevice PhysicalDevice, surface Surface) ([]PresentMode, error) {
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if surface == nil {
		return nil, NewValidationError("surface", "cannot be nil")
	}

	var modeCount C.uint32_t
	result := Result(C.vkGetPhysicalDeviceSurfacePresentModesKHR(
		C.VkPhysicalDevice(physicalDevice),
		C.VkSurfaceKHR(surface),
		&modeCount,
		nil,
	))
	if result != Success {
		return nil, NewVulkanError(result, "GetPhysicalDeviceSurfacePresentModes", "Vulkan present modes query failed")
	}

	if modeCount == 0 {
		return []PresentMode{}, nil
	}

	cModes := make([]C.VkPresentModeKHR, modeCount)
	result = Result(C.vkGetPhysicalDeviceSurfacePresentModesKHR(
		C.VkPhysicalDevice(physicalDevice),
		C.VkSurfaceKHR(surface),
		&modeCount,
		&cModes[0],
	))
	if result != Success {
		return nil, NewVulkanError(result, "GetPhysicalDeviceSurfacePresentModes", "Vulkan present modes query failed")
	}

	modes := make([]PresentMode, modeCount)
	for i := range modes {
		modes[i] = PresentMode(cModes[i])
	}

	return modes, nil
}

// PhysicalDeviceVulkan11Features contains Vulkan 1.1 features
type PhysicalDeviceVulkan11Features struct {
	StorageBuffer16BitAccess           bool
	UniformAndStorageBuffer16BitAccess bool
	StoragePushConstant16              bool
	StorageInputOutput16               bool
	Multiview                          bool
	MultiviewGeometryShader            bool
	MultiviewTessellationShader        bool
	VariablePointersStorageBuffer      bool
	VariablePointers                   bool
	ProtectedMemory                    bool
	SamplerYcbcrConversion             bool
	ShaderDrawParameters               bool
}

// PhysicalDeviceVulkan12Features contains Vulkan 1.2 features
type PhysicalDeviceVulkan12Features struct {
	SamplerMirrorClampToEdge                           bool
	DrawIndirectCount                                  bool
	StorageBuffer8BitAccess                            bool
	UniformAndStorageBuffer8BitAccess                  bool
	StoragePushConstant8                               bool
	ShaderBufferInt64Atomics                           bool
	ShaderSharedInt64Atomics                           bool
	ShaderFloat16                                      bool
	ShaderInt8                                         bool
	DescriptorIndexing                                 bool
	ShaderInputAttachmentArrayDynamicIndexing          bool
	ShaderUniformTexelBufferArrayDynamicIndexing       bool
	ShaderStorageTexelBufferArrayDynamicIndexing       bool
	ShaderUniformBufferArrayNonUniformIndexing         bool
	ShaderSampledImageArrayNonUniformIndexing          bool
	ShaderStorageBufferArrayNonUniformIndexing         bool
	ShaderStorageImageArrayNonUniformIndexing          bool
	ShaderInputAttachmentArrayNonUniformIndexing       bool
	ShaderUniformTexelBufferArrayNonUniformIndexing    bool
	ShaderStorageTexelBufferArrayNonUniformIndexing    bool
	DescriptorBindingUniformBufferUpdateAfterBind      bool
	DescriptorBindingSampledImageUpdateAfterBind       bool
	DescriptorBindingStorageImageUpdateAfterBind       bool
	DescriptorBindingStorageBufferUpdateAfterBind      bool
	DescriptorBindingUniformTexelBufferUpdateAfterBind bool
	DescriptorBindingStorageTexelBufferUpdateAfterBind bool
	DescriptorBindingUpdateUnusedWhilePending          bool
	DescriptorBindingPartiallyBound                    bool
	DescriptorBindingVariableDescriptorCount           bool
	RuntimeDescriptorArray                             bool
	SamplerFilterMinmax                                bool
	ScalarBlockLayout                                  bool
	ImagelessFramebuffer                               bool
	UniformBufferStandardLayout                        bool
	ShaderSubgroupExtendedTypes                        bool
	SeparateDepthStencilLayouts                        bool
	HostQueryReset                                     bool
	TimelineSemaphore                                  bool
	BufferDeviceAddress                                bool
	BufferDeviceAddressCaptureReplay                   bool
	BufferDeviceAddressMultiDevice                     bool
	VulkanMemoryModel                                  bool
	VulkanMemoryModelDeviceScope                       bool
	VulkanMemoryModelAvailabilityVisibilityChains      bool
	ShaderOutputViewportIndex                          bool
	ShaderOutputLayer                                  bool
	SubgroupBroadcastDynamicId                         bool
}

// PhysicalDeviceVulkan13Features contains Vulkan 1.3 features
type PhysicalDeviceVulkan13Features struct {
	RobustImageAccess                                  bool
	InlineUniformBlock                                 bool
	DescriptorBindingInlineUniformBlockUpdateAfterBind bool
	PipelineCreationCacheControl                       bool
	PrivateData                                        bool
	ShaderDemoteToHelperInvocation                     bool
	ShaderTerminateInvocation                          bool
	SubgroupSizeControl                                bool
	ComputeFullSubgroups                               bool
	Synchronization2                                   bool
	TextureCompressionASTC_HDR                         bool
	ShaderZeroInitializeWorkgroupMemory                bool
	DynamicRendering                                   bool
	ShaderIntegerDotProduct                            bool
	Maintenance4                                       bool
}

// GetPhysicalDeviceFeatures2 gets extended physical device features (Vulkan 1.1+)
func GetPhysicalDeviceFeatures2(physicalDevice PhysicalDevice) (PhysicalDeviceFeatures, error) {
	if physicalDevice == nil {
		return PhysicalDeviceFeatures{}, NewValidationError("physicalDevice", "cannot be nil")
	}

	var cFeatures2 C.VkPhysicalDeviceFeatures2
	cFeatures2.sType = C.VK_STRUCTURE_TYPE_PHYSICAL_DEVICE_FEATURES_2
	cFeatures2.pNext = nil

	C.vkGetPhysicalDeviceFeatures2(C.VkPhysicalDevice(physicalDevice), &cFeatures2)

	return physicalDeviceFeaturesFromC(&cFeatures2.features), nil
}

// PhysicalDeviceGroupProperties contains physical device group information
type PhysicalDeviceGroupProperties struct {
	PhysicalDeviceCount uint32
	PhysicalDevices     []PhysicalDevice
	SubsetAllocation    bool
}

// EnumeratePhysicalDeviceGroups enumerates physical device groups for multi-GPU
func EnumeratePhysicalDeviceGroups(instance Instance) ([]PhysicalDeviceGroupProperties, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}

	var groupCount C.uint32_t
	result := Result(C.vkEnumeratePhysicalDeviceGroups(C.VkInstance(instance), &groupCount, nil))
	if result != Success {
		return nil, NewVulkanError(result, "EnumeratePhysicalDeviceGroups", "Vulkan device group enumeration failed")
	}

	if groupCount == 0 {
		return []PhysicalDeviceGroupProperties{}, nil
	}

	cGroups := make([]C.VkPhysicalDeviceGroupProperties, groupCount)
	for i := range cGroups {
		cGroups[i].sType = C.VK_STRUCTURE_TYPE_PHYSICAL_DEVICE_GROUP_PROPERTIES
		cGroups[i].pNext = nil
	}

	result = Result(C.vkEnumeratePhysicalDeviceGroups(C.VkInstance(instance), &groupCount, &cGroups[0]))
	if result != Success {
		return nil, NewVulkanError(result, "EnumeratePhysicalDeviceGroups", "Vulkan device group enumeration failed")
	}

	groups := make([]PhysicalDeviceGroupProperties, groupCount)
	for i := range groups {
		deviceCount := uint32(cGroups[i].physicalDeviceCount)
		devices := make([]PhysicalDevice, deviceCount)
		for j := uint32(0); j < deviceCount; j++ {
			devices[j] = PhysicalDevice(cGroups[i].physicalDevices[j])
		}
		groups[i] = PhysicalDeviceGroupProperties{
			PhysicalDeviceCount: deviceCount,
			PhysicalDevices:     devices,
			SubsetAllocation:    cGroups[i].subsetAllocation == C.VK_TRUE,
		}
	}

	return groups, nil
}

// DeviceGroupDeviceCreateInfo contains device group creation information
type DeviceGroupDeviceCreateInfo struct {
	PhysicalDevices []PhysicalDevice
}
