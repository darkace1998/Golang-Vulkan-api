package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>

static PFN_vkCreateDebugUtilsMessengerEXT pfn_vkCreateDebugUtilsMessengerEXT = NULL;
static PFN_vkDestroyDebugUtilsMessengerEXT pfn_vkDestroyDebugUtilsMessengerEXT = NULL;
static PFN_vkSetDebugUtilsObjectNameEXT pfn_vkSetDebugUtilsObjectNameEXT = NULL;
static PFN_vkCmdBeginDebugUtilsLabelEXT pfn_vkCmdBeginDebugUtilsLabelEXT = NULL;
static PFN_vkCmdEndDebugUtilsLabelEXT pfn_vkCmdEndDebugUtilsLabelEXT = NULL;
static PFN_vkCmdInsertDebugUtilsLabelEXT pfn_vkCmdInsertDebugUtilsLabelEXT = NULL;
static PFN_vkQueueBeginDebugUtilsLabelEXT pfn_vkQueueBeginDebugUtilsLabelEXT = NULL;
static PFN_vkQueueEndDebugUtilsLabelEXT pfn_vkQueueEndDebugUtilsLabelEXT = NULL;
static PFN_vkQueueInsertDebugUtilsLabelEXT pfn_vkQueueInsertDebugUtilsLabelEXT = NULL;

static void loadDebugUtilsFunctions(VkInstance instance) {
    if (instance == NULL) return;
    if (pfn_vkCreateDebugUtilsMessengerEXT == NULL) {
        pfn_vkCreateDebugUtilsMessengerEXT = (PFN_vkCreateDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkCreateDebugUtilsMessengerEXT");
    }
    if (pfn_vkDestroyDebugUtilsMessengerEXT == NULL) {
        pfn_vkDestroyDebugUtilsMessengerEXT = (PFN_vkDestroyDebugUtilsMessengerEXT)vkGetInstanceProcAddr(instance, "vkDestroyDebugUtilsMessengerEXT");
    }
    if (pfn_vkSetDebugUtilsObjectNameEXT == NULL) {
        pfn_vkSetDebugUtilsObjectNameEXT = (PFN_vkSetDebugUtilsObjectNameEXT)vkGetInstanceProcAddr(instance, "vkSetDebugUtilsObjectNameEXT");
    }
    if (pfn_vkCmdBeginDebugUtilsLabelEXT == NULL) {
        pfn_vkCmdBeginDebugUtilsLabelEXT = (PFN_vkCmdBeginDebugUtilsLabelEXT)vkGetInstanceProcAddr(instance, "vkCmdBeginDebugUtilsLabelEXT");
    }
    if (pfn_vkCmdEndDebugUtilsLabelEXT == NULL) {
        pfn_vkCmdEndDebugUtilsLabelEXT = (PFN_vkCmdEndDebugUtilsLabelEXT)vkGetInstanceProcAddr(instance, "vkCmdEndDebugUtilsLabelEXT");
    }
    if (pfn_vkCmdInsertDebugUtilsLabelEXT == NULL) {
        pfn_vkCmdInsertDebugUtilsLabelEXT = (PFN_vkCmdInsertDebugUtilsLabelEXT)vkGetInstanceProcAddr(instance, "vkCmdInsertDebugUtilsLabelEXT");
    }
    if (pfn_vkQueueBeginDebugUtilsLabelEXT == NULL) {
        pfn_vkQueueBeginDebugUtilsLabelEXT = (PFN_vkQueueBeginDebugUtilsLabelEXT)vkGetInstanceProcAddr(instance, "vkQueueBeginDebugUtilsLabelEXT");
    }
    if (pfn_vkQueueEndDebugUtilsLabelEXT == NULL) {
        pfn_vkQueueEndDebugUtilsLabelEXT = (PFN_vkQueueEndDebugUtilsLabelEXT)vkGetInstanceProcAddr(instance, "vkQueueEndDebugUtilsLabelEXT");
    }
    if (pfn_vkQueueInsertDebugUtilsLabelEXT == NULL) {
        pfn_vkQueueInsertDebugUtilsLabelEXT = (PFN_vkQueueInsertDebugUtilsLabelEXT)vkGetInstanceProcAddr(instance, "vkQueueInsertDebugUtilsLabelEXT");
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

static VkResult call_vkSetDebugUtilsObjectNameEXT(VkDevice device, const VkDebugUtilsObjectNameInfoEXT* pNameInfo) {
    if (pfn_vkSetDebugUtilsObjectNameEXT == NULL) return VK_ERROR_EXTENSION_NOT_PRESENT;
    return pfn_vkSetDebugUtilsObjectNameEXT(device, pNameInfo);
}

static void call_vkCmdBeginDebugUtilsLabelEXT(VkCommandBuffer commandBuffer, const VkDebugUtilsLabelEXT* pLabelInfo) {
    if (pfn_vkCmdBeginDebugUtilsLabelEXT != NULL) {
        pfn_vkCmdBeginDebugUtilsLabelEXT(commandBuffer, pLabelInfo);
    }
}

static void call_vkCmdEndDebugUtilsLabelEXT(VkCommandBuffer commandBuffer) {
    if (pfn_vkCmdEndDebugUtilsLabelEXT != NULL) {
        pfn_vkCmdEndDebugUtilsLabelEXT(commandBuffer);
    }
}

static void call_vkCmdInsertDebugUtilsLabelEXT(VkCommandBuffer commandBuffer, const VkDebugUtilsLabelEXT* pLabelInfo) {
    if (pfn_vkCmdInsertDebugUtilsLabelEXT != NULL) {
        pfn_vkCmdInsertDebugUtilsLabelEXT(commandBuffer, pLabelInfo);
    }
}

static void call_vkQueueBeginDebugUtilsLabelEXT(VkQueue queue, const VkDebugUtilsLabelEXT* pLabelInfo) {
    if (pfn_vkQueueBeginDebugUtilsLabelEXT != NULL) {
        pfn_vkQueueBeginDebugUtilsLabelEXT(queue, pLabelInfo);
    }
}

static void call_vkQueueEndDebugUtilsLabelEXT(VkQueue queue) {
    if (pfn_vkQueueEndDebugUtilsLabelEXT != NULL) {
        pfn_vkQueueEndDebugUtilsLabelEXT(queue);
    }
}

static void call_vkQueueInsertDebugUtilsLabelEXT(VkQueue queue, const VkDebugUtilsLabelEXT* pLabelInfo) {
    if (pfn_vkQueueInsertDebugUtilsLabelEXT != NULL) {
        pfn_vkQueueInsertDebugUtilsLabelEXT(queue, pLabelInfo);
    }
}
*/
import "C"
import (
	"runtime/cgo"
	"sync"
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
	val := handle.Value()
	callback, ok := val.(DebugCallbackFunc)
	if !ok {
		return C.VK_FALSE
	}

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

	messengerUserDataMu.Lock()
	messengerUserData[DebugUtilsMessengerEXT(messenger)] = messengerCallbackData{handle: handle, pUserData: pUserData}
	messengerUserDataMu.Unlock()

	return DebugUtilsMessengerEXT(messenger), nil
}

// messengerCallbackData tracks the cgo.Handle and C allocation backing a
// messenger's pUserData so both can be released on destroy.
type messengerCallbackData struct {
	handle    cgo.Handle
	pUserData unsafe.Pointer
}

var (
	messengerUserDataMu sync.Mutex
	messengerUserData   = make(map[DebugUtilsMessengerEXT]messengerCallbackData)
)

// DestroyDebugUtilsMessengerEXT destroys a debug messenger
func DestroyDebugUtilsMessengerEXT(instance Instance, messenger DebugUtilsMessengerEXT) {
	if instance == nil || messenger == nil {
		return
	}
	C.call_vkDestroyDebugUtilsMessengerEXT(C.VkInstance(instance), C.VkDebugUtilsMessengerEXT(messenger), nil)

	// Release the cgo.Handle and the C memory that held it.
	messengerUserDataMu.Lock()
	data, ok := messengerUserData[messenger]
	delete(messengerUserData, messenger)
	messengerUserDataMu.Unlock()
	if ok {
		data.handle.Delete()
		C.free(data.pUserData)
	}
}

// DebugUtilsObjectNameInfo defines parameters for naming an object
type DebugUtilsObjectNameInfo struct {
	ObjectType   ObjectType
	ObjectHandle uint64
	ObjectName   string
}

// DebugUtilsLabel specifies parameters for a debug label
type DebugUtilsLabel struct {
	LabelName string
	Color     [4]float32
}

// SetDebugUtilsObjectNameEXT gives a user-friendly name to an object
func SetDebugUtilsObjectNameEXT(device Device, nameInfo *DebugUtilsObjectNameInfo) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if nameInfo == nil {
		return NewValidationError("nameInfo", "cannot be nil")
	}

	cName := C.CString(nameInfo.ObjectName)
	defer C.free(unsafe.Pointer(cName))

	var cInfo C.VkDebugUtilsObjectNameInfoEXT
	cInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_OBJECT_NAME_INFO_EXT
	cInfo.pNext = nil
	cInfo.objectType = C.VkObjectType(nameInfo.ObjectType)
	cInfo.objectHandle = C.uint64_t(nameInfo.ObjectHandle)
	cInfo.pObjectName = cName

	res := C.call_vkSetDebugUtilsObjectNameEXT(C.VkDevice(device), &cInfo)
	if Result(res) != Success {
		return NewVulkanError(Result(res), "vkSetDebugUtilsObjectNameEXT", "failed to set debug utils object name")
	}
	return nil
}

// CmdBeginDebugUtilsLabelEXT opens a command buffer debug label region
func CmdBeginDebugUtilsLabelEXT(commandBuffer CommandBuffer, labelInfo *DebugUtilsLabel) {
	if commandBuffer == nil || labelInfo == nil {
		return
	}

	cLabelName := C.CString(labelInfo.LabelName)
	defer C.free(unsafe.Pointer(cLabelName))

	var cInfo C.VkDebugUtilsLabelEXT
	cInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_LABEL_EXT
	cInfo.pNext = nil
	cInfo.pLabelName = cLabelName
	cInfo.color[0] = C.float(labelInfo.Color[0])
	cInfo.color[1] = C.float(labelInfo.Color[1])
	cInfo.color[2] = C.float(labelInfo.Color[2])
	cInfo.color[3] = C.float(labelInfo.Color[3])

	C.call_vkCmdBeginDebugUtilsLabelEXT(C.VkCommandBuffer(commandBuffer), &cInfo)
}

// CmdEndDebugUtilsLabelEXT closes a command buffer debug label region
func CmdEndDebugUtilsLabelEXT(commandBuffer CommandBuffer) {
	if commandBuffer == nil {
		return
	}
	C.call_vkCmdEndDebugUtilsLabelEXT(C.VkCommandBuffer(commandBuffer))
}

// CmdInsertDebugUtilsLabelEXT inserts a single debug label into a command buffer
func CmdInsertDebugUtilsLabelEXT(commandBuffer CommandBuffer, labelInfo *DebugUtilsLabel) {
	if commandBuffer == nil || labelInfo == nil {
		return
	}

	cLabelName := C.CString(labelInfo.LabelName)
	defer C.free(unsafe.Pointer(cLabelName))

	var cInfo C.VkDebugUtilsLabelEXT
	cInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_LABEL_EXT
	cInfo.pNext = nil
	cInfo.pLabelName = cLabelName
	cInfo.color[0] = C.float(labelInfo.Color[0])
	cInfo.color[1] = C.float(labelInfo.Color[1])
	cInfo.color[2] = C.float(labelInfo.Color[2])
	cInfo.color[3] = C.float(labelInfo.Color[3])

	C.call_vkCmdInsertDebugUtilsLabelEXT(C.VkCommandBuffer(commandBuffer), &cInfo)
}

// QueueBeginDebugUtilsLabelEXT opens a queue debug label region
func QueueBeginDebugUtilsLabelEXT(queue Queue, labelInfo *DebugUtilsLabel) {
	if queue == nil || labelInfo == nil {
		return
	}

	cLabelName := C.CString(labelInfo.LabelName)
	defer C.free(unsafe.Pointer(cLabelName))

	var cInfo C.VkDebugUtilsLabelEXT
	cInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_LABEL_EXT
	cInfo.pNext = nil
	cInfo.pLabelName = cLabelName
	cInfo.color[0] = C.float(labelInfo.Color[0])
	cInfo.color[1] = C.float(labelInfo.Color[1])
	cInfo.color[2] = C.float(labelInfo.Color[2])
	cInfo.color[3] = C.float(labelInfo.Color[3])

	C.call_vkQueueBeginDebugUtilsLabelEXT(C.VkQueue(queue), &cInfo)
}

// QueueEndDebugUtilsLabelEXT closes a queue debug label region
func QueueEndDebugUtilsLabelEXT(queue Queue) {
	if queue == nil {
		return
	}
	C.call_vkQueueEndDebugUtilsLabelEXT(C.VkQueue(queue))
}

// QueueInsertDebugUtilsLabelEXT inserts a single debug label into a queue
func QueueInsertDebugUtilsLabelEXT(queue Queue, labelInfo *DebugUtilsLabel) {
	if queue == nil || labelInfo == nil {
		return
	}

	cLabelName := C.CString(labelInfo.LabelName)
	defer C.free(unsafe.Pointer(cLabelName))

	var cInfo C.VkDebugUtilsLabelEXT
	cInfo.sType = C.VK_STRUCTURE_TYPE_DEBUG_UTILS_LABEL_EXT
	cInfo.pNext = nil
	cInfo.pLabelName = cLabelName
	cInfo.color[0] = C.float(labelInfo.Color[0])
	cInfo.color[1] = C.float(labelInfo.Color[1])
	cInfo.color[2] = C.float(labelInfo.Color[2])
	cInfo.color[3] = C.float(labelInfo.Color[3])

	C.call_vkQueueInsertDebugUtilsLabelEXT(C.VkQueue(queue), &cInfo)
}
