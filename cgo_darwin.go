//go:build darwin

package vulkan

/*
#cgo pkg-config: vulkan
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>

// Helper function to convert Go string slice to C char**
char** makeCharArray(int size) {
    return calloc(sizeof(char*), size);
}

// Helper function to set string in char array
void setArrayString(char **a, char *s, int n) {
    a[n] = s;
}

// Helper function to free char array
void freeCharArray(char **a, int size) {
    for (int i = 0; i < size; i++) {
        free(a[i]);
    }
    free(a);
}
*/
import "C"
