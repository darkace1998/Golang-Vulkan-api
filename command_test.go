package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyCommandPoolNilArgs tests that DestroyCommandPool handles nil gracefully
func TestDestroyCommandPoolNilArgs(t *testing.T) {
	DestroyCommandPool(nil, nil)
	DestroyCommandPool(nil, fakeCommandPool())
	DestroyCommandPool(fakeDevice(), nil)
}

// TestDestroySemaphoreNilArgs tests that DestroySemaphore handles nil gracefully
func TestDestroySemaphoreNilArgs(t *testing.T) {
	DestroySemaphore(nil, nil)
	DestroySemaphore(nil, fakeSemaphore())
	DestroySemaphore(fakeDevice(), nil)
}

// TestDestroyFenceNilArgs tests that DestroyFence handles nil gracefully
func TestDestroyFenceNilArgs(t *testing.T) {
	DestroyFence(nil, nil)
	DestroyFence(nil, fakeFence())
	DestroyFence(fakeDevice(), nil)
}

// ============================================================================
// Validation Tests for Create Functions
// ============================================================================

// TestCreateCommandPoolValidation tests nil parameter validation
func TestCreateCommandPoolValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *CommandPoolCreateInfo
		expectParam string
	}{
		{testNilDeviceStr, nil, &CommandPoolCreateInfo{}, testDeviceParameter},
		{"nil createInfo", fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateCommandPool(tt.device, tt.createInfo)
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

// TestAllocateCommandBuffersValidation tests nil parameter validation
func TestAllocateCommandBuffersValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		allocInfo   *CommandBufferAllocateInfo
		expectParam string
	}{
		{testNilDeviceStr, nil, &CommandBufferAllocateInfo{CommandBufferCount: 1}, testDeviceParameter},
		{"nil allocateInfo", fakeDevice(), nil, "allocateInfo"},
		{"zero count", fakeDevice(), &CommandBufferAllocateInfo{CommandBufferCount: 0}, "allocateInfo.CommandBufferCount"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AllocateCommandBuffers(tt.device, tt.allocInfo)
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

// TestBeginCommandBufferValidation tests nil parameter validation
func TestBeginCommandBufferValidation(t *testing.T) {
	tests := []struct {
		name          string
		commandBuffer CommandBuffer
		beginInfo     *CommandBufferBeginInfo
		expectParam   string
	}{
		{"nil commandBuffer", nil, &CommandBufferBeginInfo{}, "commandBuffer"},
		{"nil beginInfo", fakeCommandBuffer(), nil, "beginInfo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BeginCommandBuffer(tt.commandBuffer, tt.beginInfo)
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

// TestEndCommandBufferValidation tests nil parameter validation
func TestEndCommandBufferValidation(t *testing.T) {
	err := EndCommandBuffer(nil)
	if err == nil {
		t.Fatal("Expected error for nil commandBuffer")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Parameter != "commandBuffer" {
		t.Errorf("Expected param 'commandBuffer', got '%s'", valErr.Parameter)
	}
}

// TestQueueSubmitValidation tests nil parameter validation
func TestQueueSubmitValidation(t *testing.T) {
	err := QueueSubmit(nil, nil, nil)
	if err == nil {
		t.Fatal("Expected error for nil queue")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Parameter != "queue" {
		t.Errorf("Expected param 'queue', got '%s'", valErr.Parameter)
	}
}

// TestCreateSemaphoreValidation tests nil parameter validation
func TestCreateSemaphoreValidation(t *testing.T) {
	_, err := CreateSemaphore(nil, nil)
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

// TestCreateFenceValidation tests nil parameter validation
func TestCreateFenceValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *FenceCreateInfo
		expectParam string
	}{
		{testNilDeviceStr, nil, &FenceCreateInfo{}, testDeviceParameter},
		{"nil createInfo", fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateFence(tt.device, tt.createInfo)
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

// TestWaitForFencesValidation tests nil parameter validation
func TestWaitForFencesValidation(t *testing.T) {
	err := WaitForFences(nil, []Fence{fakeFence()}, true, 1000)
	if err == nil {
		t.Fatal("Expected error for nil device")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
}

// TestWaitForFencesEmptySlice tests that empty fences is a no-op
func TestWaitForFencesEmptySlice(t *testing.T) {
	err := WaitForFences(fakeDevice(), []Fence{}, true, 1000)
	if err != nil {
		t.Fatalf("Expected nil for empty fences, got: %v", err)
	}
}

// TestResetFencesValidation tests nil parameter validation
func TestResetFencesValidation(t *testing.T) {
	err := ResetFences(nil, []Fence{fakeFence()})
	if err == nil {
		t.Fatal("Expected error for nil device")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
}

// TestResetFencesEmptySlice tests that empty fences is a no-op
func TestResetFencesEmptySlice(t *testing.T) {
	err := ResetFences(fakeDevice(), []Fence{})
	if err != nil {
		t.Fatalf("Expected nil for empty fences, got: %v", err)
	}
}

// TestFreeCommandBuffersNilArgs tests FreeCommandBuffers handles nil gracefully
func TestFreeCommandBuffersNilArgs(t *testing.T) {
	// All these should not panic
	FreeCommandBuffers(nil, nil, nil)
	FreeCommandBuffers(nil, fakeCommandPool(), []CommandBuffer{fakeCommandBuffer()})
	FreeCommandBuffers(fakeDevice(), nil, []CommandBuffer{fakeCommandBuffer()})
	FreeCommandBuffers(fakeDevice(), fakeCommandPool(), []CommandBuffer{})
	FreeCommandBuffers(fakeDevice(), fakeCommandPool(), nil)
}
