//go:build linux && !vk_headless && !vk_no_wayland

package vulkan

/*
#cgo pkg-config: vulkan
#define VK_USE_PLATFORM_WAYLAND_KHR
#include <vulkan/vulkan.h>
#include <wayland-client.h>
*/
import "C"
import "unsafe"

// WaylandSurfaceCreateInfoKHR contains parameters for creating a Wayland surface
type WaylandSurfaceCreateInfoKHR struct {
	Display unsafe.Pointer // *C.struct_wl_display
	Surface unsafe.Pointer // *C.struct_wl_surface
}

// CreateWaylandSurfaceKHR creates a Vulkan surface for a Wayland window
func CreateWaylandSurfaceKHR(instance Instance, createInfo *WaylandSurfaceCreateInfoKHR) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.Display == nil {
		return nil, NewValidationError("createInfo.Display", "cannot be nil")
	}
	if createInfo.Surface == nil {
		return nil, NewValidationError("createInfo.Surface", "cannot be nil")
	}

	var cCreateInfo C.VkWaylandSurfaceCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_WAYLAND_SURFACE_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.display = (*C.struct_wl_display)(createInfo.Display)
	cCreateInfo.surface = (*C.struct_wl_surface)(createInfo.Surface)

	var surface C.VkSurfaceKHR
	result := Result(C.vkCreateWaylandSurfaceKHR(C.VkInstance(instance), &cCreateInfo, nil, &surface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateWaylandSurfaceKHR", "failed to create Wayland surface")
	}
	return Surface(surface), nil
}
