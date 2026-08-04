package vulkan

/*
#include <stdlib.h>
#include <vulkan/vulkan.h>

static PFN_vkCreateRayTracingPipelinesKHR load_vkCreateRayTracingPipelinesKHR(VkDevice device) {
    return (PFN_vkCreateRayTracingPipelinesKHR)vkGetDeviceProcAddr(device, "vkCreateRayTracingPipelinesKHR");
}

static PFN_vkCmdTraceRaysKHR load_vkCmdTraceRaysKHR(VkDevice device) {
    return (PFN_vkCmdTraceRaysKHR)vkGetDeviceProcAddr(device, "vkCmdTraceRaysKHR");
}

static VkResult call_vkCreateRayTracingPipelinesKHR(PFN_vkCreateRayTracingPipelinesKHR pfn, VkDevice device, VkDeferredOperationKHR deferredOperation, VkPipelineCache pipelineCache, uint32_t createInfoCount, const VkRayTracingPipelineCreateInfoKHR* pCreateInfos, const VkAllocationCallbacks* pAllocator, VkPipeline* pPipelines) {
    if (pfn != NULL) {
        return pfn(device, deferredOperation, pipelineCache, createInfoCount, pCreateInfos, pAllocator, pPipelines);
    }
    return VK_ERROR_EXTENSION_NOT_PRESENT;
}

static void call_vkCmdTraceRaysKHR(PFN_vkCmdTraceRaysKHR pfn, VkCommandBuffer commandBuffer, const VkStridedDeviceAddressRegionKHR* pRaygenShaderBindingTable, const VkStridedDeviceAddressRegionKHR* pMissShaderBindingTable, const VkStridedDeviceAddressRegionKHR* pHitShaderBindingTable, const VkStridedDeviceAddressRegionKHR* pCallableShaderBindingTable, uint32_t width, uint32_t height, uint32_t depth) {
    pfn(commandBuffer, pRaygenShaderBindingTable, pMissShaderBindingTable, pHitShaderBindingTable, pCallableShaderBindingTable, width, height, depth);
}
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"
)

// RayTracingFunctions holds the device-level VK_KHR_ray_tracing_pipeline
// function pointers for one specific VkDevice. Device-level function pointers
// are only valid for the device they were queried from, so applications using
// multiple devices must use one RayTracingFunctions per device.
type RayTracingFunctions struct {
	device          Device
	createPipelines C.PFN_vkCreateRayTracingPipelinesKHR
	traceRays       C.PFN_vkCmdTraceRaysKHR
}

var (
	rayTracingFuncsMu      sync.RWMutex
	rayTracingFuncsByDev   = make(map[Device]*RayTracingFunctions)
	rayTracingFuncsDefault *RayTracingFunctions
)

// LoadRayTracingPipelineFunctions resolves the device-level ray tracing
// pipeline functions for the given device and returns them. The result is
// cached per device; loading is idempotent and thread-safe.
//
// The first successfully loaded device also becomes the dispatch target for
// the package-level CmdTraceRaysKHR convenience function. Applications with
// more than one device must call methods on the returned RayTracingFunctions
// instead of the package-level functions.
//
// Returns an error if the device is nil or the extension is unavailable.
func LoadRayTracingPipelineFunctions(device Device) (*RayTracingFunctions, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}

	rayTracingFuncsMu.RLock()
	f := rayTracingFuncsByDev[device]
	rayTracingFuncsMu.RUnlock()
	if f != nil {
		return f, nil
	}

	f = &RayTracingFunctions{
		device:          device,
		createPipelines: C.load_vkCreateRayTracingPipelinesKHR(C.VkDevice(device)),
		traceRays:       C.load_vkCmdTraceRaysKHR(C.VkDevice(device)),
	}
	if f.createPipelines == nil {
		return nil, NewVulkanError(ErrorExtensionNotPresent, "LoadRayTracingPipelineFunctions", "vkCreateRayTracingPipelinesKHR not available (is VK_KHR_ray_tracing_pipeline enabled?)")
	}

	rayTracingFuncsMu.Lock()
	if existing := rayTracingFuncsByDev[device]; existing != nil {
		f = existing
	} else {
		rayTracingFuncsByDev[device] = f
		if rayTracingFuncsDefault == nil {
			rayTracingFuncsDefault = f
		}
	}
	rayTracingFuncsMu.Unlock()
	return f, nil
}

// unloadRayTracingFunctions drops cached function pointers for a destroyed
// device.
func unloadRayTracingFunctions(device Device) {
	rayTracingFuncsMu.Lock()
	defer rayTracingFuncsMu.Unlock()
	delete(rayTracingFuncsByDev, device)
	if rayTracingFuncsDefault != nil && rayTracingFuncsDefault.device == device {
		rayTracingFuncsDefault = nil
		for _, f := range rayTracingFuncsByDev {
			rayTracingFuncsDefault = f
			break
		}
	}
}

// rayTracingFunctionsForDevice returns the cached functions for a device,
// loading them on first use.
func rayTracingFunctionsForDevice(device Device) (*RayTracingFunctions, error) {
	rayTracingFuncsMu.RLock()
	f := rayTracingFuncsByDev[device]
	rayTracingFuncsMu.RUnlock()
	if f != nil {
		return f, nil
	}
	return LoadRayTracingPipelineFunctions(device)
}

// defaultRayTracingFunctions returns the functions for the first loaded
// device, used by the package-level convenience wrappers.
func defaultRayTracingFunctions() *RayTracingFunctions {
	rayTracingFuncsMu.RLock()
	defer rayTracingFuncsMu.RUnlock()
	return rayTracingFuncsDefault
}

// RayTracingShaderGroupTypeKHR represents the type of a ray tracing shader group.
type RayTracingShaderGroupTypeKHR int32

const (
	RayTracingShaderGroupTypeGeneralKHR            RayTracingShaderGroupTypeKHR = C.VK_RAY_TRACING_SHADER_GROUP_TYPE_GENERAL_KHR
	RayTracingShaderGroupTypeTrianglesHitGroupKHR  RayTracingShaderGroupTypeKHR = C.VK_RAY_TRACING_SHADER_GROUP_TYPE_TRIANGLES_HIT_GROUP_KHR
	RayTracingShaderGroupTypeProceduralHitGroupKHR RayTracingShaderGroupTypeKHR = C.VK_RAY_TRACING_SHADER_GROUP_TYPE_PROCEDURAL_HIT_GROUP_KHR
)

const ShaderUnusedKHR uint32 = C.VK_SHADER_UNUSED_KHR

// RayTracingShaderGroupCreateInfoKHR represents the VkRayTracingShaderGroupCreateInfoKHR structure.
type RayTracingShaderGroupCreateInfoKHR struct {
	Type                RayTracingShaderGroupTypeKHR
	GeneralShader       uint32
	ClosestHitShader    uint32
	AnyHitShader        uint32
	IntersectionShader  uint32
	AnyHitShaderDefault uint32
}

// PipelineCreateFlags represents pipeline creation flags.
type PipelineCreateFlags uint32

const (
	PipelineCreateDisableOptimizationBit PipelineCreateFlags = C.VK_PIPELINE_CREATE_DISABLE_OPTIMIZATION_BIT
	PipelineCreateAllowDerivativesBit    PipelineCreateFlags = C.VK_PIPELINE_CREATE_ALLOW_DERIVATIVES_BIT
	PipelineCreateDerivativeBit          PipelineCreateFlags = C.VK_PIPELINE_CREATE_DERIVATIVE_BIT
)

// RayTracingPipelineCreateInfoKHR represents the VkRayTracingPipelineCreateInfoKHR structure.
type RayTracingPipelineCreateInfoKHR struct {
	Flags                        PipelineCreateFlags
	Stages                       []PipelineShaderStageCreateInfo
	Groups                       []RayTracingShaderGroupCreateInfoKHR
	MaxPipelineRayRecursionDepth uint32
	LibraryInfo                  *PipelineLibraryCreateInfoKHR
	LibraryInterface             *RayTracingPipelineInterfaceCreateInfoKHR
	DynamicState                 *PipelineDynamicStateCreateInfo
	Layout                       PipelineLayout
	BasePipelineHandle           Pipeline
	BasePipelineIndex            int32
}

// PipelineLibraryCreateInfoKHR represents VkPipelineLibraryCreateInfoKHR (stub)
type PipelineLibraryCreateInfoKHR struct{}

// RayTracingPipelineInterfaceCreateInfoKHR represents VkRayTracingPipelineInterfaceCreateInfoKHR (stub)
type RayTracingPipelineInterfaceCreateInfoKHR struct{}

// StridedDeviceAddressRegionKHR represents the VkStridedDeviceAddressRegionKHR structure.
type StridedDeviceAddressRegionKHR struct {
	DeviceAddress DeviceAddress
	Stride        DeviceSize
	Size          DeviceSize
}

// CreateRayTracingPipelinesKHR creates ray tracing pipelines. The functions
// for the device are loaded on first use (per device).
func CreateRayTracingPipelinesKHR(device Device, pipelineCache PipelineCache, createInfos []RayTracingPipelineCreateInfoKHR) ([]Pipeline, error) {
	if device == nil {
		return nil, &ValidationError{Field: "Device", Reason: "cannot be nil"}
	}
	if len(createInfos) == 0 {
		return nil, &ValidationError{Field: "createInfos", Reason: "cannot be empty"}
	}

	funcs, err := rayTracingFunctionsForDevice(device)
	if err != nil {
		return nil, err
	}

	cCreateInfos := make([]C.VkRayTracingPipelineCreateInfoKHR, len(createInfos))

	// Nested arrays and name strings must be pinned because their addresses
	// are stored inside cCreateInfos, which is Go memory passed to C.
	var pinner runtime.Pinner
	defer pinner.Unpin()

	for i, ci := range createInfos {
		cCreateInfos[i].sType = C.VK_STRUCTURE_TYPE_RAY_TRACING_PIPELINE_CREATE_INFO_KHR
		cCreateInfos[i].flags = C.VkPipelineCreateFlags(ci.Flags)
		cCreateInfos[i].maxPipelineRayRecursionDepth = C.uint32_t(ci.MaxPipelineRayRecursionDepth)
		cCreateInfos[i].layout = C.VkPipelineLayout(ci.Layout)
		cCreateInfos[i].basePipelineHandle = C.VkPipeline(ci.BasePipelineHandle)
		cCreateInfos[i].basePipelineIndex = C.int32_t(ci.BasePipelineIndex)

		// Stages
		cCreateInfos[i].stageCount = C.uint32_t(len(ci.Stages))
		if len(ci.Stages) > 0 {
			cStages := make([]C.VkPipelineShaderStageCreateInfo, len(ci.Stages))
			for j, stage := range ci.Stages {
				cStages[j].sType = C.VK_STRUCTURE_TYPE_PIPELINE_SHADER_STAGE_CREATE_INFO
				cStages[j].stage = C.VkShaderStageFlagBits(stage.Stage)
				cStages[j].module = C.VkShaderModule(stage.Module)

				nameBytes := []byte(stage.Name + "\x00")
				pinner.Pin(&nameBytes[0])
				cStages[j].pName = (*C.char)(unsafe.Pointer(&nameBytes[0]))
			}
			pinner.Pin(&cStages[0])
			cCreateInfos[i].pStages = &cStages[0]
		}

		// Groups
		cCreateInfos[i].groupCount = C.uint32_t(len(ci.Groups))
		if len(ci.Groups) > 0 {
			cGroups := make([]C.VkRayTracingShaderGroupCreateInfoKHR, len(ci.Groups))
			for j, group := range ci.Groups {
				cGroups[j].sType = C.VK_STRUCTURE_TYPE_RAY_TRACING_SHADER_GROUP_CREATE_INFO_KHR
				cGroups[j]._type = C.VkRayTracingShaderGroupTypeKHR(group.Type)
				cGroups[j].generalShader = C.uint32_t(group.GeneralShader)
				cGroups[j].closestHitShader = C.uint32_t(group.ClosestHitShader)
				cGroups[j].anyHitShader = C.uint32_t(group.AnyHitShader)
				cGroups[j].intersectionShader = C.uint32_t(group.IntersectionShader)
			}
			pinner.Pin(&cGroups[0])
			cCreateInfos[i].pGroups = &cGroups[0]
		}
	}

	cPipelines := make([]C.VkPipeline, len(createInfos))

	result := C.call_vkCreateRayTracingPipelinesKHR(
		funcs.createPipelines,
		C.VkDevice(device),
		C.VkDeferredOperationKHR(nil), // deferred operation not supported yet
		C.VkPipelineCache(pipelineCache),
		C.uint32_t(len(createInfos)),
		&cCreateInfos[0],
		nil,
		&cPipelines[0],
	)

	if Result(result) != Success {
		// A failed batch create may still have created some pipelines
		// (failed entries are VK_NULL_HANDLE); destroy them to avoid leaks.
		for _, p := range cPipelines {
			if p != nil {
				C.vkDestroyPipeline(C.VkDevice(device), p, nil)
			}
		}
		return nil, &VulkanError{Result: Result(result)}
	}

	pipelines := make([]Pipeline, len(createInfos))
	for i, p := range cPipelines {
		pipelines[i] = Pipeline(p)
		trackResource("Pipeline", unsafe.Pointer(p))
	}

	return pipelines, nil
}

// CmdTraceRaysKHR records a trace-rays command using this device's function
// pointers.
func (f *RayTracingFunctions) CmdTraceRaysKHR(commandBuffer CommandBuffer, raygen, miss, hit, callable *StridedDeviceAddressRegionKHR, width, height, depth uint32) {
	if f == nil || f.traceRays == nil || commandBuffer == nil {
		return
	}

	var cRaygen, cMiss, cHit, cCallable C.VkStridedDeviceAddressRegionKHR

	if raygen != nil {
		cRaygen.deviceAddress = C.VkDeviceAddress(raygen.DeviceAddress)
		cRaygen.stride = C.VkDeviceSize(raygen.Stride)
		cRaygen.size = C.VkDeviceSize(raygen.Size)
	}
	if miss != nil {
		cMiss.deviceAddress = C.VkDeviceAddress(miss.DeviceAddress)
		cMiss.stride = C.VkDeviceSize(miss.Stride)
		cMiss.size = C.VkDeviceSize(miss.Size)
	}
	if hit != nil {
		cHit.deviceAddress = C.VkDeviceAddress(hit.DeviceAddress)
		cHit.stride = C.VkDeviceSize(hit.Stride)
		cHit.size = C.VkDeviceSize(hit.Size)
	}
	if callable != nil {
		cCallable.deviceAddress = C.VkDeviceAddress(callable.DeviceAddress)
		cCallable.stride = C.VkDeviceSize(callable.Stride)
		cCallable.size = C.VkDeviceSize(callable.Size)
	}

	C.call_vkCmdTraceRaysKHR(
		f.traceRays,
		C.VkCommandBuffer(commandBuffer),
		&cRaygen,
		&cMiss,
		&cHit,
		&cCallable,
		C.uint32_t(width),
		C.uint32_t(height),
		C.uint32_t(depth),
	)
}

// CmdTraceRaysKHR records a trace-rays command using the functions of the
// first device passed to LoadRayTracingPipelineFunctions. Single-device
// convenience; multi-device applications must use RayTracingFunctions
// methods.
func CmdTraceRaysKHR(commandBuffer CommandBuffer, raygen, miss, hit, callable *StridedDeviceAddressRegionKHR, width, height, depth uint32) {
	defaultRayTracingFunctions().CmdTraceRaysKHR(commandBuffer, raygen, miss, hit, callable, width, height, depth)
}
