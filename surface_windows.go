//go:build windows

package vulkan

/*
#cgo LDFLAGS: -lvulkan-1
#define VK_USE_PLATFORM_WIN32_KHR
#include <vulkan/vulkan.h>
*/
import "C"
import "unsafe"

// Win32SurfaceCreateInfoKHR contains parameters for creating a Win32 surface
type Win32SurfaceCreateInfoKHR struct {
	Hinstance unsafe.Pointer // HINSTANCE
	Hwnd      unsafe.Pointer // HWND
}

// CreateWin32SurfaceKHR creates a Vulkan surface for a Win32 window
func CreateWin32SurfaceKHR(instance Instance, createInfo *Win32SurfaceCreateInfoKHR) (Surface, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}

	var cCreateInfo C.VkWin32SurfaceCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_WIN32_SURFACE_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.hinstance = (C.HINSTANCE)(createInfo.Hinstance)
	cCreateInfo.hwnd = (C.HWND)(createInfo.Hwnd)

	var surface C.VkSurfaceKHR
	result := Result(C.vkCreateWin32SurfaceKHR(C.VkInstance(instance), &cCreateInfo, nil, &surface))
	if result != Success {
		return nil, NewVulkanError(result, "CreateWin32SurfaceKHR", "failed to create Win32 surface")
	}
	return Surface(surface), nil
}
