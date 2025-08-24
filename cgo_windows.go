//go:build windows

package vulkan

/*
// For Vulkan SDK installed in standard locations:
#cgo LDFLAGS: -lvulkan-1
// Alternative if Vulkan SDK is in a custom location:
// #cgo CFLAGS: -I"C:/VulkanSDK/1.3.290.0/Include"
// #cgo LDFLAGS: -L"C:/VulkanSDK/1.3.290.0/Lib" -lvulkan-1
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
