package vulkan

import (
	"testing"
)

func TestDebugUtilsMessengerValidation(t *testing.T) {
	var testCallback DebugCallbackFunc = func(
		messageSeverity DebugUtilsMessageSeverityFlags,
		messageType DebugUtilsMessageTypeFlags,
		callbackData *DebugUtilsMessengerCallbackData,
	) bool {
		return false
	}

	_, err := CreateDebugUtilsMessengerEXT(nil, &DebugUtilsMessengerCreateInfo{}, testCallback)
	if err == nil {
		t.Error("Expected error when creating debug messenger with nil instance")
	}

	_, err = CreateDebugUtilsMessengerEXT(Instance(fakeHandle()), nil, testCallback)
	if err == nil {
		t.Error("Expected error when creating debug messenger with nil createInfo")
	}

	_, err = CreateDebugUtilsMessengerEXT(Instance(fakeHandle()), &DebugUtilsMessengerCreateInfo{}, nil)
	if err == nil {
		t.Error("Expected error when creating debug messenger with nil callback")
	}
}
