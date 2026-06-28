//go:build darwin

package vulkan

/*
#cgo pkg-config: vulkan
#cgo CFLAGS: -DVK_USE_PLATFORM_METAL_EXT
#include <vulkan/vulkan.h>
#include <vulkan/vulkan_metal.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// CreateMetalSurfaceEXT creates a Metal surface on macOS/iOS
func CreateMetalSurfaceEXT(instance Instance, layer unsafe.Pointer) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if layer == nil {
		return nil, NewValidationError("layer", "cannot be nil")
	}

	var createInfo C.VkMetalSurfaceCreateInfoEXT
	createInfo.sType = C.VK_STRUCTURE_TYPE_METAL_SURFACE_CREATE_INFO_EXT
	createInfo.pNext = nil
	createInfo.flags = 0
	// The layer is expected to be a pointer to a CAMetalLayer
	createInfo.pLayer = (*C.CAMetalLayer)(layer)

	var cSurface C.VkSurfaceKHR
	result := Result(C.vkCreateMetalSurfaceEXT(C.VkInstance(instance), &createInfo, nil, &cSurface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateMetalSurfaceEXT", "Vulkan create Metal surface failed")
	}

	trackResource("Surface", unsafe.Pointer(cSurface))
	return Surface(cSurface), nil
}
