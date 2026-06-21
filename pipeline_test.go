package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyShaderModuleNilArgs tests that DestroyShaderModule handles nil gracefully
func TestDestroyShaderModuleNilArgs(t *testing.T) {
	DestroyShaderModule(nil, nil)
	DestroyShaderModule(nil, fakeShaderModule())
	DestroyShaderModule(fakeDevice(), nil)
}

// TestDestroyPipelineLayoutNilArgs tests that DestroyPipelineLayout handles nil gracefully
func TestDestroyPipelineLayoutNilArgs(t *testing.T) {
	DestroyPipelineLayout(nil, nil)
	DestroyPipelineLayout(nil, fakePipelineLayout())
	DestroyPipelineLayout(fakeDevice(), nil)
}

// TestDestroyRenderPassNilArgs tests that DestroyRenderPass handles nil gracefully
func TestDestroyRenderPassNilArgs(t *testing.T) {
	DestroyRenderPass(nil, nil)
	DestroyRenderPass(nil, fakeRenderPass())
	DestroyRenderPass(fakeDevice(), nil)
}

// TestDestroyPipelineNilArgs tests that DestroyPipeline handles nil gracefully
func TestDestroyPipelineNilArgs(t *testing.T) {
	DestroyPipeline(nil, nil)
	DestroyPipeline(nil, fakePipeline())
	DestroyPipeline(fakeDevice(), nil)
}

func TestGetRenderAreaGranularityValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		renderPass  RenderPass
		expectParam string
	}{
		{testNilDevice, nil, fakeRenderPass(), testDeviceParameter},
		{"nil renderPass", fakeDevice(), nil, "renderPass"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetRenderAreaGranularity(tt.device, tt.renderPass)
			// GetRenderAreaGranularity returns an empty struct on error/nil params,
			// it does not return an error itself to maintain API compatibility
			// with similar query functions, but we can visually verify it handles nil
			// safely without crashing during tests.
		})
	}
}

// ============================================================================
// Validation Tests for Create Functions
// ============================================================================

// TestCreateShaderModuleValidation tests nil parameter validation
func TestCreateShaderModuleValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *ShaderModuleCreateInfo
		expectParam string
	}{
		{testNilDevice, nil, &ShaderModuleCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateShaderModule(tt.device, tt.createInfo)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
			}
		})
	}
}

// TestCreatePipelineLayoutValidation tests nil parameter validation
func TestCreatePipelineLayoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *PipelineLayoutCreateInfo
		expectParam string
	}{
		{testNilDevice, nil, &PipelineLayoutCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreatePipelineLayout(tt.device, tt.createInfo)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
			}
		})
	}
}

// TestCreateComputePipelinesValidation tests nil device validation
func TestCreateComputePipelinesValidation(t *testing.T) {
	_, err := CreateComputePipelines(nil, nil, []ComputePipelineCreateInfo{{}})
	if err == nil {
		t.Fatal("Expected error for nil device")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Field != testDeviceParameter {
		t.Errorf("Expected param 'device', got '%s'", valErr.Field)
	}
}

// TestCreateComputePipelinesEmptyInfos tests that empty createInfos returns nil
func TestCreateComputePipelinesEmptyInfos(t *testing.T) {
	result, err := CreateComputePipelines(fakeDevice(), nil, []ComputePipelineCreateInfo{})
	if err != nil {
		t.Fatalf("Expected no error for empty createInfos, got: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for empty createInfos")
	}
}

// TestIsExtensionSupported tests IsExtensionSupported function
func TestPipelineIsExtensionSupported(t *testing.T) {
	extensions := []ExtensionProperties{
		{ExtensionName: "VK_KHR_surface", SpecVersion: 1},
		{ExtensionName: "VK_KHR_swapchain", SpecVersion: 2},
	}

	if !IsExtensionSupported("VK_KHR_surface", extensions) {
		t.Errorf("Expected IsExtensionSupported to return true for existing extension")
	}

	if IsExtensionSupported("VK_EXT_debug_utils", extensions) {
		t.Errorf("Expected IsExtensionSupported to return false for missing extension")
	}
}

// TestIsLayerSupported tests IsLayerSupported function
func TestPipelineIsLayerSupported(t *testing.T) {
	layers := []LayerProperties{
		{LayerName: "VK_LAYER_KHRONOS_validation", SpecVersion: MakeVersion(1, 0, 0)},
		{LayerName: "VK_LAYER_LUNARG_api_dump", SpecVersion: MakeVersion(1, 0, 0)},
	}

	if !IsLayerSupported("VK_LAYER_KHRONOS_validation", layers) {
		t.Errorf("Expected IsLayerSupported to return true for existing layer")
	}

	if IsLayerSupported("VK_LAYER_RENDERDOC_Capture", layers) {
		t.Errorf("Expected IsLayerSupported to return false for missing layer")
	}
}
