//go:build darwin

package vulkan

/*
#cgo pkg-config: vulkan
#define VK_USE_PLATFORM_METAL_EXT
#include <vulkan/vulkan.h>
*/
import "C"
import "unsafe"

// MetalSurfaceCreateInfoEXT contains parameters for creating a Metal surface
type MetalSurfaceCreateInfoEXT struct {
	Layer unsafe.Pointer // *CAMetalLayer
}

// CreateMetalSurfaceEXT creates a Vulkan surface for a Metal layer
func CreateMetalSurfaceEXT(instance Instance, createInfo *MetalSurfaceCreateInfoEXT) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	var cCreateInfo C.VkMetalSurfaceCreateInfoEXT
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_METAL_SURFACE_CREATE_INFO_EXT
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.pLayer = (*C.CAMetalLayer)(createInfo.Layer)

	var surface C.VkSurfaceKHR
	result := Result(C.vkCreateMetalSurfaceEXT(C.VkInstance(instance), &cCreateInfo, nil, &surface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateMetalSurfaceEXT", "failed to create Metal surface")
	}
	return Surface(surface), nil
}
