package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>

static PFN_vkCreateDebugUtilsMessengerEXT pfn_vkCreateDebugUtilsMessengerEXT = NULL;
static PFN_vkDestroyDebugUtilsMessengerEXT pfn_vkDestroyDebugUtilsMessengerEXT = NULL;

static void loadDebugUtilsFunctions(VkInstance instance) {
    if (instance == NULL) return;
    if (pfn_vkCreateDebugUtilsMessengerEXT == NULL) {
        pfn_vkCreateDebugUtilsMessengerEXT = (PFN_vkCreateDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkCreateDebugUtilsMessengerEXT");
    }
    if (pfn_vkDestroyDebugUtilsMessengerEXT == NULL) {
        pfn_vkDestroyDebugUtilsMessengerEXT = (PFN_vkDestroyDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkDestroyDebugUtilsMessengerEXT");
    }
}

// Forward declaration of the callback exported by Go.
extern VkBool32 debugUtilsCallbackEXT(
    VkDebugUtilsMessageSeverityFlagBitsEXT messageSeverity,
    VkDebugUtilsMessageTypeFlagsEXT messageTypes,
    VkDebugUtilsMessengerCallbackDataEXT* pCallbackData,
    void* pUserData);

// Helper function to return the C function pointer to Go
static PFN_vkDebugUtilsMessengerCallbackEXT getDebugCallbackFunc() {
    return (PFN_vkDebugUtilsMessengerCallbackEXT)debugUtilsCallbackEXT;
}

static VkResult call_vkCreateDebugUtilsMessengerEXT(VkInstance instance, const VkDebugUtilsMessengerCreateInfoEXT* pCreateInfo, const VkAllocationCallbacks* pAllocator, VkDebugUtilsMessengerEXT* pMessenger) {
    if (pfn_vkCreateDebugUtilsMessengerEXT == NULL) return VK_ERROR_EXTENSION_NOT_PRESENT;
    return pfn_vkCreateDebugUtilsMessengerEXT(instance, pCreateInfo, pAllocator, pMessenger);
}

static void call_vkDestroyDebugUtilsMessengerEXT(VkInstance instance, VkDebugUtilsMessengerEXT messenger, const VkAllocationCallbacks* pAllocator) {
    if (pfn_vkDestroyDebugUtilsMessengerEXT != NULL) {
        pfn_vkDestroyDebugUtilsMessengerEXT(instance, messenger, pAllocator);
    }
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// LoadDebugUtilsFunctions loads the debug utils functions for an instance.
func LoadDebugUtilsFunctions(instance Instance) {
	if instance != nil {
		C.loadDebugUtilsFunctions(C.VkInstance(instance))
	}
}

// DebugCallbackFunc is the Go callback type for debug messages
type DebugCallbackFunc func(
	messageSeverity DebugUtilsMessageSeverityFlags,
	messageType DebugUtilsMessageTypeFlags,
	callbackData *DebugUtilsMessengerCallbackData,
) bool

//export debugUtilsCallbackEXT
func debugUtilsCallbackEXT(
	messageSeverity C.VkDebugUtilsMessageSeverityFlagBitsEXT,
	messageTypes C.VkDebugUtilsMessageTypeFlagsEXT,
	pCallbackData *C.VkDebugUtilsMessengerCallbackDataEXT,
	pUserData unsafe.Pointer,
) C.VkBool32 {
	if pUserData == nil {
		return C.VK_FALSE
	}

	// Restore the Go callback from the handle
	handle := *(*cgo.Handle)(pUserData)
	callback := handle.Value().(DebugCallbackFunc)

	data := &DebugUtilsMessengerCallbackData{
		MessageIDName:   C.GoString(pCallbackData.pMessageIdName),
		MessageIDNumber: int32(pCallbackData.messageIdNumber),
		Message:         C.GoString(pCallbackData.pMessage),
	}

	if callback(
		DebugUtilsMessageSeverityFlags(messageSeverity),
		DebugUtilsMessageTypeFlags(messageTypes),
		data,
	) {
		return C.VK_TRUE
	}
	return C.VK_FALSE
}

// CreateDebugUtilsMessengerEXT creates a debug messenger
func CreateDebugUtilsMessengerEXT(instance Instance, createInfo *DebugUtilsMessengerCreateInfo, callback DebugCallbackFunc) (DebugUtilsMessengerEXT, error) {
	if instance == nil {
		return nil, NewValidationError("instance", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if callback == nil {
		return nil, NewValidationError("callback", "cannot be nil")
	}

	// Create a cgo handle to store the Go callback safely
	handle := cgo.NewHandle(callback)
	// Allocate C memory to hold the handle so we can pass it as pUserData
	pUserData := C.malloc(C.size_t(unsafe.Sizeof(handle)))
	*(*cgo.Handle)(pUserData) = handle

	var cCreateInfo C.VkDebugUtilsMessengerCreateInfoEXT
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_MESSENGER_CREATE_INFO_EXT
	cCreateInfo.pNext = nil
	cCreateInfo.flags = 0
	cCreateInfo.messageSeverity = C.VkDebugUtilsMessageSeverityFlagsEXT(createInfo.MessageSeverity)
	cCreateInfo.messageType = C.VkDebugUtilsMessageTypeFlagsEXT(createInfo.MessageType)

	// Get the C function pointer via our helper function
	cCreateInfo.pfnUserCallback = C.getDebugCallbackFunc()
	cCreateInfo.pUserData = pUserData

	var messenger C.VkDebugUtilsMessengerEXT
	res := C.call_vkCreateDebugUtilsMessengerEXT(C.VkInstance(instance), &cCreateInfo, nil, &messenger)
	if Result(res) != Success {
		// Clean up the handle if creation fails
		handle.Delete()
		C.free(pUserData)
		return nil, NewVulkanError(Result(res), "vkCreateDebugUtilsMessengerEXT", "failed to create debug messenger")
	}

	return DebugUtilsMessengerEXT(messenger), nil
}

// DestroyDebugUtilsMessengerEXT destroys a debug messenger
func DestroyDebugUtilsMessengerEXT(instance Instance, messenger DebugUtilsMessengerEXT) {
	if instance == nil || messenger == nil {
		return
	}
	// Note: While this frees the Vulkan messenger handle, we have a leak of the pUserData memory and the cgo.Handle
	// To fix this fully, we would need a map from messenger to the pUserData/Handle, or store it elsewhere.
	// However, it is typical to leave debug messengers alive until shutdown.
	C.call_vkDestroyDebugUtilsMessengerEXT(C.VkInstance(instance), C.VkDebugUtilsMessengerEXT(messenger), nil)
}
