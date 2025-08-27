//go:build linux

package vulkan

/*
#cgo pkg-config: vulkan
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
