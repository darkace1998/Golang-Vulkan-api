package vulkan

import (
	"testing"
)

// ============================================================================
// Extension and Layer Support Tests
// ============================================================================

func TestIsExtensionSupported(t *testing.T) {
	availableExtensions := []ExtensionProperties{
		{ExtensionName: "VK_KHR_surface", SpecVersion: 25},
		{ExtensionName: "VK_KHR_swapchain", SpecVersion: 68},
	}

	tests := []struct {
		name          string
		extensionName string
		expected      bool
	}{
		{"SupportedExtension", "VK_KHR_surface", true},
		{"AnotherSupportedExtension", "VK_KHR_swapchain", true},
		{"UnsupportedExtension", "VK_KHR_video_queue", false},
		{"EmptyExtensionName", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExtensionSupported(tt.extensionName, availableExtensions)
			if result != tt.expected {
				t.Errorf("Expected IsExtensionSupported(%s) to be %v, got %v", tt.extensionName, tt.expected, result)
			}
		})
	}
}

func TestIsLayerSupported(t *testing.T) {
	availableLayers := []LayerProperties{
		{LayerName: "VK_LAYER_KHRONOS_validation", SpecVersion: 1, ImplementationVersion: 1, Description: "Validation Layer"},
		{LayerName: "VK_LAYER_LUNARG_api_dump", SpecVersion: 1, ImplementationVersion: 1, Description: "API Dump Layer"},
	}

	tests := []struct {
		name      string
		layerName string
		expected  bool
	}{
		{"SupportedLayer", "VK_LAYER_KHRONOS_validation", true},
		{"AnotherSupportedLayer", "VK_LAYER_LUNARG_api_dump", true},
		{"UnsupportedLayer", "VK_LAYER_RENDERDOC_Capture", false},
		{"EmptyLayerName", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLayerSupported(tt.layerName, availableLayers)
			if result != tt.expected {
				t.Errorf("Expected IsLayerSupported(%s) to be %v, got %v", tt.layerName, tt.expected, result)
			}
		})
	}
}

// ============================================================================
// Video Device Functions Tests
// ============================================================================

func TestCreateVideoDeviceFunctions(t *testing.T) {
	// Test nil device
	funcs, err := CreateVideoDeviceFunctions(nil)
	if err == nil {
		t.Error("expected error for nil device, got nil")
	}
	if funcs != nil {
		t.Error("expected nil functions for nil device, got non-nil")
	}

	// Setup fake device
	device := fakeDevice()

	// Clear the global map before testing
	videoDeviceFunctionsMapLock.Lock()
	videoDeviceFunctionsMap = make(map[Device]*VideoDeviceFunctions)
	videoDeviceFunctionsMapLock.Unlock()

	// Test successful creation
	funcs1, err := CreateVideoDeviceFunctions(device)
	if err != nil {
		t.Errorf("unexpected error creating video device functions: %v", err)
	}
	if funcs1 == nil {
		t.Fatal("expected non-nil functions, got nil")
	}

	// Verify the loaded state
	loaded := funcs1.IsLoaded()
	// Just calling it to ensure coverage, whether it's loaded depends on if the system actually has Vulkan loaded,
	// for a fake handle, C.loadAdditionalVideoDeviceFunctions likely returns 0, so IsLoaded should be false.
	if loaded {
		t.Log("Note: IsLoaded returned true")
	}

	// Test idempotency (creating again should return the same instance)
	funcs2, err := CreateVideoDeviceFunctions(device)
	if err != nil {
		t.Errorf("unexpected error calling CreateVideoDeviceFunctions twice: %v", err)
	}
	if funcs1 != funcs2 {
		t.Error("expected second call to return the exact same instance")
	}

	// Test GetVideoDeviceFunctions
	funcs3 := GetVideoDeviceFunctions(device)
	if funcs1 != funcs3 {
		t.Error("expected GetVideoDeviceFunctions to return the created instance")
	}

	// Test GetVideoDeviceFunctions for unknown device
	unknownDevice := Device(fakeHandle()) // creating a new unique handle
	funcs4 := GetVideoDeviceFunctions(unknownDevice)
	if funcs4 != nil {
		t.Error("expected GetVideoDeviceFunctions to return nil for unknown device")
	}
}
