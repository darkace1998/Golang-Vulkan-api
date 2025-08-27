//go:build linux

package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
