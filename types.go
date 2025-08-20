package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// Version represents Vulkan API version
type Version uint32

// Vulkan API versions
const (
	Version10 Version = C.VK_API_VERSION_1_0
	Version11 Version = C.VK_API_VERSION_1_1
	Version12 Version = C.VK_API_VERSION_1_2
	Version13 Version = C.VK_API_VERSION_1_3
	// Version14 will be available when system supports Vulkan 1.4
	Version14 Version = (1 << 22) | (4 << 12) // VK_MAKE_API_VERSION(0, 1, 4, 0)
)

// MakeVersion creates a version number from major, minor, and patch components
func MakeVersion(major, minor, patch uint32) Version {
	return Version((major << 22) | (minor << 12) | patch)
}

// VersionMajor extracts the major version number
func (v Version) Major() uint32 {
	return uint32((v >> 22) & 0x7F)
}

// VersionMinor extracts the minor version number
func (v Version) Minor() uint32 {
	return uint32((v >> 12) & 0x3FF)
}

// VersionPatch extracts the patch version number
func (v Version) Patch() uint32 {
	return uint32(v & 0xFFF)
}

// Result represents Vulkan result codes
type Result int32

// Vulkan result codes
const (
	Success                                    Result = C.VK_SUCCESS
	NotReady                                   Result = C.VK_NOT_READY
	Timeout                                    Result = C.VK_TIMEOUT
	EventSet                                   Result = C.VK_EVENT_SET
	EventReset                                 Result = C.VK_EVENT_RESET
	Incomplete                                 Result = C.VK_INCOMPLETE
	ErrorOutOfHostMemory                       Result = C.VK_ERROR_OUT_OF_HOST_MEMORY
	ErrorOutOfDeviceMemory                     Result = C.VK_ERROR_OUT_OF_DEVICE_MEMORY
	ErrorInitializationFailed                  Result = C.VK_ERROR_INITIALIZATION_FAILED
	ErrorDeviceLost                            Result = C.VK_ERROR_DEVICE_LOST
	ErrorMemoryMapFailed                       Result = C.VK_ERROR_MEMORY_MAP_FAILED
	ErrorLayerNotPresent                       Result = C.VK_ERROR_LAYER_NOT_PRESENT
	ErrorExtensionNotPresent                   Result = C.VK_ERROR_EXTENSION_NOT_PRESENT
	ErrorFeatureNotPresent                     Result = C.VK_ERROR_FEATURE_NOT_PRESENT
	ErrorIncompatibleDriver                    Result = C.VK_ERROR_INCOMPATIBLE_DRIVER
	ErrorTooManyObjects                        Result = C.VK_ERROR_TOO_MANY_OBJECTS
	ErrorFormatNotSupported                    Result = C.VK_ERROR_FORMAT_NOT_SUPPORTED
	ErrorFragmentedPool                        Result = C.VK_ERROR_FRAGMENTED_POOL
	ErrorUnknown                               Result = C.VK_ERROR_UNKNOWN
	ErrorOutOfPoolMemory                       Result = C.VK_ERROR_OUT_OF_POOL_MEMORY
	ErrorInvalidExternalHandle                 Result = C.VK_ERROR_INVALID_EXTERNAL_HANDLE
	ErrorFragmentation                         Result = C.VK_ERROR_FRAGMENTATION
	ErrorInvalidOpaqueCaptureAddress           Result = C.VK_ERROR_INVALID_OPAQUE_CAPTURE_ADDRESS
	ErrorSurfaceLostKHR                        Result = C.VK_ERROR_SURFACE_LOST_KHR
	ErrorNativeWindowInUseKHR                  Result = C.VK_ERROR_NATIVE_WINDOW_IN_USE_KHR
	SuboptimalKHR                              Result = C.VK_SUBOPTIMAL_KHR
	ErrorOutOfDateKHR                          Result = C.VK_ERROR_OUT_OF_DATE_KHR
	ErrorIncompatibleDisplayKHR                Result = C.VK_ERROR_INCOMPATIBLE_DISPLAY_KHR
	ErrorValidationFailedEXT                   Result = C.VK_ERROR_VALIDATION_FAILED_EXT
	ErrorInvalidShaderNV                       Result = C.VK_ERROR_INVALID_SHADER_NV
	ErrorInvalidDrmFormatModifierPlaneLayoutEXT Result = C.VK_ERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXT
	ErrorNotPermittedEXT                       Result = C.VK_ERROR_NOT_PERMITTED_EXT
	ErrorFullScreenExclusiveModeLostEXT        Result = C.VK_ERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXT
	ThreadIdleKHR                              Result = C.VK_THREAD_IDLE_KHR
	ThreadDoneKHR                              Result = C.VK_THREAD_DONE_KHR
	OperationDeferredKHR                       Result = C.VK_OPERATION_DEFERRED_KHR
	OperationNotDeferredKHR                    Result = C.VK_OPERATION_NOT_DEFERRED_KHR
	PipelineCompileRequiredEXT                 Result = C.VK_PIPELINE_COMPILE_REQUIRED_EXT
)

// Error returns the error message for the result
func (r Result) Error() string {
	switch r {
	case Success:
		return "VK_SUCCESS"
	case NotReady:
		return "VK_NOT_READY"
	case Timeout:
		return "VK_TIMEOUT"
	case EventSet:
		return "VK_EVENT_SET"
	case EventReset:
		return "VK_EVENT_RESET"
	case Incomplete:
		return "VK_INCOMPLETE"
	case ErrorOutOfHostMemory:
		return "VK_ERROR_OUT_OF_HOST_MEMORY"
	case ErrorOutOfDeviceMemory:
		return "VK_ERROR_OUT_OF_DEVICE_MEMORY"
	case ErrorInitializationFailed:
		return "VK_ERROR_INITIALIZATION_FAILED"
	case ErrorDeviceLost:
		return "VK_ERROR_DEVICE_LOST"
	case ErrorMemoryMapFailed:
		return "VK_ERROR_MEMORY_MAP_FAILED"
	case ErrorLayerNotPresent:
		return "VK_ERROR_LAYER_NOT_PRESENT"
	case ErrorExtensionNotPresent:
		return "VK_ERROR_EXTENSION_NOT_PRESENT"
	case ErrorFeatureNotPresent:
		return "VK_ERROR_FEATURE_NOT_PRESENT"
	case ErrorIncompatibleDriver:
		return "VK_ERROR_INCOMPATIBLE_DRIVER"
	case ErrorTooManyObjects:
		return "VK_ERROR_TOO_MANY_OBJECTS"
	case ErrorFormatNotSupported:
		return "VK_ERROR_FORMAT_NOT_SUPPORTED"
	case ErrorFragmentedPool:
		return "VK_ERROR_FRAGMENTED_POOL"
	case ErrorUnknown:
		return "VK_ERROR_UNKNOWN"
	case ErrorOutOfPoolMemory:
		return "VK_ERROR_OUT_OF_POOL_MEMORY"
	case ErrorInvalidExternalHandle:
		return "VK_ERROR_INVALID_EXTERNAL_HANDLE"
	case ErrorFragmentation:
		return "VK_ERROR_FRAGMENTATION"
	case ErrorInvalidOpaqueCaptureAddress:
		return "VK_ERROR_INVALID_OPAQUE_CAPTURE_ADDRESS"
	case ErrorSurfaceLostKHR:
		return "VK_ERROR_SURFACE_LOST_KHR"
	case ErrorNativeWindowInUseKHR:
		return "VK_ERROR_NATIVE_WINDOW_IN_USE_KHR"
	case SuboptimalKHR:
		return "VK_SUBOPTIMAL_KHR"
	case ErrorOutOfDateKHR:
		return "VK_ERROR_OUT_OF_DATE_KHR"
	case ErrorIncompatibleDisplayKHR:
		return "VK_ERROR_INCOMPATIBLE_DISPLAY_KHR"
	case ErrorValidationFailedEXT:
		return "VK_ERROR_VALIDATION_FAILED_EXT"
	case ErrorInvalidShaderNV:
		return "VK_ERROR_INVALID_SHADER_NV"
	case ErrorInvalidDrmFormatModifierPlaneLayoutEXT:
		return "VK_ERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXT"
	case ErrorNotPermittedEXT:
		return "VK_ERROR_NOT_PERMITTED_EXT"
	case ErrorFullScreenExclusiveModeLostEXT:
		return "VK_ERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXT"
	case ThreadIdleKHR:
		return "VK_THREAD_IDLE_KHR"
	case ThreadDoneKHR:
		return "VK_THREAD_DONE_KHR"
	case OperationDeferredKHR:
		return "VK_OPERATION_DEFERRED_KHR"
	case OperationNotDeferredKHR:
		return "VK_OPERATION_NOT_DEFERRED_KHR"
	case PipelineCompileRequiredEXT:
		return "VK_PIPELINE_COMPILE_REQUIRED_EXT"
	default:
		return "Unknown Vulkan error"
	}
}

// IsError returns true if the result represents an error condition
func (r Result) IsError() bool {
	return r < 0
}

// IsSuccess returns true if the result represents success
func (r Result) IsSuccess() bool {
	return r >= 0
}

// Bool type for Vulkan boolean values
type Bool32 uint32

const (
	False Bool32 = C.VK_FALSE
	True  Bool32 = C.VK_TRUE
)

// ToBool converts a Bool32 to a Go bool
func (b Bool32) ToBool() bool {
	return b == True
}

// FromBool converts a Go bool to Bool32
func FromBool(b bool) Bool32 {
	if b {
		return True
	}
	return False
}

// DeviceSize represents device memory size
type DeviceSize uint64

// DeviceAddress represents device memory address
type DeviceAddress uint64

// Flags represents generic flags
type Flags uint32

// SampleCount represents sample count flags
type SampleCountFlags uint32

const (
	SampleCount1Bit  SampleCountFlags = C.VK_SAMPLE_COUNT_1_BIT
	SampleCount2Bit  SampleCountFlags = C.VK_SAMPLE_COUNT_2_BIT
	SampleCount4Bit  SampleCountFlags = C.VK_SAMPLE_COUNT_4_BIT
	SampleCount8Bit  SampleCountFlags = C.VK_SAMPLE_COUNT_8_BIT
	SampleCount16Bit SampleCountFlags = C.VK_SAMPLE_COUNT_16_BIT
	SampleCount32Bit SampleCountFlags = C.VK_SAMPLE_COUNT_32_BIT
	SampleCount64Bit SampleCountFlags = C.VK_SAMPLE_COUNT_64_BIT
)

// Handle types
type (
	Instance               unsafe.Pointer
	PhysicalDevice         unsafe.Pointer
	Device                 unsafe.Pointer
	Queue                  unsafe.Pointer
	Semaphore              unsafe.Pointer
	CommandBuffer          unsafe.Pointer
	Fence                  unsafe.Pointer
	DeviceMemory           unsafe.Pointer
	Buffer                 unsafe.Pointer
	Image                  unsafe.Pointer
	Event                  unsafe.Pointer
	QueryPool              unsafe.Pointer
	BufferView             unsafe.Pointer
	ImageView              unsafe.Pointer
	ShaderModule           unsafe.Pointer
	PipelineCache          unsafe.Pointer
	PipelineLayout         unsafe.Pointer
	RenderPass             unsafe.Pointer
	Pipeline               unsafe.Pointer
	DescriptorSetLayout    unsafe.Pointer
	Sampler                unsafe.Pointer
	DescriptorPool         unsafe.Pointer
	DescriptorSet          unsafe.Pointer
	Framebuffer            unsafe.Pointer
	CommandPool            unsafe.Pointer
	Surface                unsafe.Pointer
	Swapchain              unsafe.Pointer
	Display                unsafe.Pointer
	DisplayMode            unsafe.Pointer
	DescriptorUpdateTemplate unsafe.Pointer
	SamplerYcbcrConversion unsafe.Pointer
	ValidationCache        unsafe.Pointer
	AccelerationStructure  unsafe.Pointer
	PerformanceConfiguration unsafe.Pointer
	DeferredOperation      unsafe.Pointer
	PrivateDataSlot        unsafe.Pointer
	VideoSession           unsafe.Pointer
	VideoSessionParameters unsafe.Pointer
	CuModule               unsafe.Pointer
	CuFunction             unsafe.Pointer
	OpticalFlowSession     unsafe.Pointer
	MicromapEXT            unsafe.Pointer
	ShaderEXT              unsafe.Pointer
)

// Null handle constants
var (
	NullHandle = unsafe.Pointer(nil)
)

// Constants
const (
	MaxMemoryTypes          = C.VK_MAX_MEMORY_TYPES
	MaxMemoryHeaps          = C.VK_MAX_MEMORY_HEAPS
	MaxPhysicalDeviceNameSize = C.VK_MAX_PHYSICAL_DEVICE_NAME_SIZE
	MaxExtensionNameSize    = C.VK_MAX_EXTENSION_NAME_SIZE
	MaxDescriptionSize      = C.VK_MAX_DESCRIPTION_SIZE
	UuidSize                = C.VK_UUID_SIZE
	LuidSize                = C.VK_LUID_SIZE
	MaxDriverNameSize       = C.VK_MAX_DRIVER_NAME_SIZE
	MaxDriverInfoSize       = C.VK_MAX_DRIVER_INFO_SIZE
	AttachmentUnused        = C.VK_ATTACHMENT_UNUSED
	SubpassExternal         = C.VK_SUBPASS_EXTERNAL
	QueueFamilyIgnored      = C.VK_QUEUE_FAMILY_IGNORED
	QueueFamilyExternal     = C.VK_QUEUE_FAMILY_EXTERNAL
	QueueFamilyForeignEXT   = C.VK_QUEUE_FAMILY_FOREIGN_EXT
	RemainingMipLevels      = C.VK_REMAINING_MIP_LEVELS
	RemainingArrayLayers    = C.VK_REMAINING_ARRAY_LAYERS
	WholeSize               = uint64(C.VK_WHOLE_SIZE)
)