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
		{"nil device", nil, &ShaderModuleCreateInfo{}, testDeviceParameter},
		{"nil createInfo", fakeDevice(), nil, testCreateInfoParameter},
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
			if valErr.Parameter != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Parameter)
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
		{"nil device", nil, &PipelineLayoutCreateInfo{}, testDeviceParameter},
		{"nil createInfo", fakeDevice(), nil, testCreateInfoParameter},
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
			if valErr.Parameter != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Parameter)
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
	if valErr.Parameter != testDeviceParameter {
		t.Errorf("Expected param 'device', got '%s'", valErr.Parameter)
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
