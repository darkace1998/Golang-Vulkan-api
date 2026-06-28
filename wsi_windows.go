//go:build windows

package vulkan

/*
#cgo LDFLAGS: -lvulkan-1
#cgo CFLAGS: -DVK_USE_PLATFORM_WIN32_KHR
#include <vulkan/vulkan.h>
#include <vulkan/vulkan_win32.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// CreateWin32SurfaceKHR creates a Windows surface
func CreateWin32SurfaceKHR(instance Instance, hinstance unsafe.Pointer, hwnd unsafe.Pointer) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if hinstance == nil {
		return nil, NewValidationError("hinstance", "cannot be nil")
	}
	if hwnd == nil {
		return nil, NewValidationError("hwnd", "cannot be nil")
	}

	var createInfo C.VkWin32SurfaceCreateInfoKHR
	createInfo.sType = C.VK_STRUCTURE_TYPE_WIN32_SURFACE_CREATE_INFO_KHR
	createInfo.pNext = nil
	createInfo.flags = 0
	createInfo.hinstance = (C.HINSTANCE)(hinstance)
	createInfo.hwnd = (C.HWND)(hwnd)

	var cSurface C.VkSurfaceKHR
	result := Result(C.vkCreateWin32SurfaceKHR(C.VkInstance(instance), &createInfo, nil, &cSurface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateWin32SurfaceKHR", "Vulkan create Win32 surface failed")
	}

	trackResource("Surface", unsafe.Pointer(cSurface))
	return Surface(cSurface), nil
}
