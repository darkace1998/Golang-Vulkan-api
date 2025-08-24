//go:build unix && !linux && !darwin

package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
