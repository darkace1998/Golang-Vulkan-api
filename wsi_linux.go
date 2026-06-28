//go:build linux || freebsd || openbsd || netbsd

package vulkan

/*
#cgo pkg-config: vulkan
#cgo CFLAGS: -DVK_USE_PLATFORM_XLIB_KHR -DVK_USE_PLATFORM_XCB_KHR -DVK_USE_PLATFORM_WAYLAND_KHR
#include <vulkan/vulkan.h>
#include <vulkan/vulkan_xlib.h>
#include <vulkan/vulkan_xcb.h>
#include <vulkan/vulkan_wayland.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// CreateXlibSurfaceKHR creates an Xlib surface
func CreateXlibSurfaceKHR(instance Instance, dpy unsafe.Pointer, window uint64) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if dpy == nil {
		return nil, NewValidationError("dpy", "cannot be nil")
	}

	var createInfo C.VkXlibSurfaceCreateInfoKHR
	createInfo.sType = C.VK_STRUCTURE_TYPE_XLIB_SURFACE_CREATE_INFO_KHR
	createInfo.pNext = nil
	createInfo.flags = 0
	createInfo.dpy = (*C.Display)(dpy)
	createInfo.window = C.Window(window)

	var cSurface C.VkSurfaceKHR
	result := Result(C.vkCreateXlibSurfaceKHR(C.VkInstance(instance), &createInfo, nil, &cSurface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateXlibSurfaceKHR", "Vulkan create Xlib surface failed")
	}

	trackResource("Surface", unsafe.Pointer(cSurface))
	return Surface(cSurface), nil
}

// CreateXcbSurfaceKHR creates an XCB surface
func CreateXcbSurfaceKHR(instance Instance, connection unsafe.Pointer, window uint32) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if connection == nil {
		return nil, NewValidationError("connection", "cannot be nil")
	}

	var createInfo C.VkXcbSurfaceCreateInfoKHR
	createInfo.sType = C.VK_STRUCTURE_TYPE_XCB_SURFACE_CREATE_INFO_KHR
	createInfo.pNext = nil
	createInfo.flags = 0
	createInfo.connection = (*C.xcb_connection_t)(connection)
	createInfo.window = C.xcb_window_t(window)

	var cSurface C.VkSurfaceKHR
	result := Result(C.vkCreateXcbSurfaceKHR(C.VkInstance(instance), &createInfo, nil, &cSurface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateXcbSurfaceKHR", "Vulkan create XCB surface failed")
	}

	trackResource("Surface", unsafe.Pointer(cSurface))
	return Surface(cSurface), nil
}

// CreateWaylandSurfaceKHR creates a Wayland surface
func CreateWaylandSurfaceKHR(instance Instance, display unsafe.Pointer, surface unsafe.Pointer) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if display == nil {
		return nil, NewValidationError("display", "cannot be nil")
	}
	if surface == nil {
		return nil, NewValidationError("surface", "cannot be nil")
	}

	var createInfo C.VkWaylandSurfaceCreateInfoKHR
	createInfo.sType = C.VK_STRUCTURE_TYPE_WAYLAND_SURFACE_CREATE_INFO_KHR
	createInfo.pNext = nil
	createInfo.flags = 0
	createInfo.display = (*C.struct_wl_display)(display)
	createInfo.surface = (*C.struct_wl_surface)(surface)

	var cSurface C.VkSurfaceKHR
	result := Result(C.vkCreateWaylandSurfaceKHR(C.VkInstance(instance), &createInfo, nil, &cSurface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateWaylandSurfaceKHR", "Vulkan create Wayland surface failed")
	}

	trackResource("Surface", unsafe.Pointer(cSurface))
	return Surface(cSurface), nil
}
