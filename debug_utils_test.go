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

func TestSetDebugUtilsObjectNameValidation(t *testing.T) {
	err := SetDebugUtilsObjectNameEXT(nil, &DebugUtilsObjectNameInfo{})
	if err == nil {
		t.Error("Expected error when setting object name with nil device")
	}

	err = SetDebugUtilsObjectNameEXT(Device(fakeHandle()), nil)
	if err == nil {
		t.Error("Expected error when setting object name with nil info")
	}

	// Should not panic with valid fake device
	err = SetDebugUtilsObjectNameEXT(Device(fakeHandle()), &DebugUtilsObjectNameInfo{
		ObjectType:   ObjectTypeBuffer,
		ObjectHandle: 0,
		ObjectName:   "TestBuffer",
	})
	// This will likely return an error because the extension is not present/loaded,
	// but it should not crash.
	_ = err
}

func TestDebugUtilsLabelsValidation(t *testing.T) {
	// These should not crash or panic when called with nil handles or info
	CmdBeginDebugUtilsLabelEXT(nil, &DebugUtilsLabel{})
	CmdBeginDebugUtilsLabelEXT(CommandBuffer(fakeHandle()), nil)

	CmdEndDebugUtilsLabelEXT(nil)

	CmdInsertDebugUtilsLabelEXT(nil, &DebugUtilsLabel{})
	CmdInsertDebugUtilsLabelEXT(CommandBuffer(fakeHandle()), nil)

	QueueBeginDebugUtilsLabelEXT(nil, &DebugUtilsLabel{})
	QueueBeginDebugUtilsLabelEXT(Queue(fakeHandle()), nil)

	QueueEndDebugUtilsLabelEXT(nil)

	QueueInsertDebugUtilsLabelEXT(nil, &DebugUtilsLabel{})
	QueueInsertDebugUtilsLabelEXT(Queue(fakeHandle()), nil)

	// And valid inputs (with fake handles) shouldn't crash
	labelInfo := &DebugUtilsLabel{
		LabelName: "TestLabel",
		Color:     [4]float32{1.0, 0.0, 0.0, 1.0},
	}
	CmdBeginDebugUtilsLabelEXT(CommandBuffer(fakeHandle()), labelInfo)
	CmdEndDebugUtilsLabelEXT(CommandBuffer(fakeHandle()))
	CmdInsertDebugUtilsLabelEXT(CommandBuffer(fakeHandle()), labelInfo)

	QueueBeginDebugUtilsLabelEXT(Queue(fakeHandle()), labelInfo)
	QueueEndDebugUtilsLabelEXT(Queue(fakeHandle()))
	QueueInsertDebugUtilsLabelEXT(Queue(fakeHandle()), labelInfo)
}
