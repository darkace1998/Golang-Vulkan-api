//go:build linux

package vulkan

/*
#cgo pkg-config: vulkan
#define VK_USE_PLATFORM_XLIB_KHR
#define VK_USE_PLATFORM_WAYLAND_KHR
#include <vulkan/vulkan.h>
#include <X11/Xlib.h>
#include <wayland-client.h>
*/
import "C"
import "unsafe"

// XlibSurfaceCreateInfoKHR contains parameters for creating an Xlib surface
type XlibSurfaceCreateInfoKHR struct {
	Dpy    unsafe.Pointer // *C.Display
	Window uintptr        // C.Window
}

// CreateXlibSurfaceKHR creates a Vulkan surface for an X11 window
func CreateXlibSurfaceKHR(instance Instance, createInfo *XlibSurfaceCreateInfoKHR) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.Dpy == nil {
		return nil, NewValidationError("createInfo.Dpy", "cannot be nil")
	}

	var cCreateInfo C.VkXlibSurfaceCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_XLIB_SURFACE_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.dpy = (*C.Display)(createInfo.Dpy)
	cCreateInfo.window = C.Window(createInfo.Window)

	var surface C.VkSurfaceKHR
	result := Result(C.vkCreateXlibSurfaceKHR(C.VkInstance(instance), &cCreateInfo, nil, &surface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateXlibSurfaceKHR", "failed to create Xlib surface")
	}
	return Surface(surface), nil
}

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
