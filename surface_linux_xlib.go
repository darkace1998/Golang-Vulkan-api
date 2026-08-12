//go:build linux && !vk_headless && !vk_no_xlib

package vulkan

/*
#cgo pkg-config: vulkan
#define VK_USE_PLATFORM_XLIB_KHR
#include <vulkan/vulkan.h>
#include <X11/Xlib.h>
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
