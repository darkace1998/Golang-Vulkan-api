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

// Helper function to convert Go string slice to C char**
char** makeCharArray(int size) {
    char** result = calloc(size, sizeof(char*));
    return result; // Returns NULL if allocation fails
}

// Helper function to set string in char array
void setArrayString(char **a, char *s, int n) {
    a[n] = s;
}

// Helper function to free char array
void freeCharArray(char **a, int size) {
    if (a == NULL) {
        return; // Safely handle NULL pointer
    }
    for (int i = 0; i < size; i++) {
        free(a[i]);
    }
    free(a);
}
*/
import "C"
