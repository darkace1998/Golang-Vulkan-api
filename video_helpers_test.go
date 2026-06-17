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
