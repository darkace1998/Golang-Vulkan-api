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
		{testNilDevice, nil, &CommandPoolCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
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
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
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
		{testNilDevice, nil, &CommandBufferAllocateInfo{CommandBufferCount: 1}, testDeviceParameter},
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
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
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
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
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
	if valErr.Field != "commandBuffer" {
		t.Errorf("Expected param 'commandBuffer', got '%s'", valErr.Field)
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
	if valErr.Field != "queue" {
		t.Errorf("Expected param 'queue', got '%s'", valErr.Field)
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
	if valErr.Field != testDeviceParameter {
		t.Errorf("Expected param 'device', got '%s'", valErr.Field)
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
		{testNilDevice, nil, &FenceCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
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
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
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

// TestCmdExecuteCommandsValidation validates CmdExecuteCommands function
func TestCmdExecuteCommandsValidation(t *testing.T) {
	// These shouldn't panic
	CmdExecuteCommands(nil, nil)
	CmdExecuteCommands(fakeCommandBuffer(), nil)
}

func TestGetFenceStatusValidation(t *testing.T) {
	_, err := GetFenceStatus(nil, fakeFence())
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	_, err = GetFenceStatus(fakeDevice(), nil)
	if err == nil {
		t.Errorf("Expected error when fence is nil")
	} else if err.Error() != testNilFenceError {
		t.Errorf("Unexpected error message: %v", err)
	}
}
// ============================================================================
// Benchmarks
// ============================================================================

// TestCreateThreadLocalCommandPool tests the creation of a thread-local command pool
func TestCreateThreadLocalCommandPool(t *testing.T) {
	originalCreateCommandPool := createCommandPoolFunc
	defer func() { createCommandPoolFunc = originalCreateCommandPool }()

	t.Run("Validation error", testCreateThreadLocalCommandPoolValidation)
	t.Run("Success path", testCreateThreadLocalCommandPoolSuccess)
	t.Run("Error path", testCreateThreadLocalCommandPoolError)
}

func testCreateThreadLocalCommandPoolValidation(t *testing.T) {
	_, err := CreateThreadLocalCommandPool(nil, 0)
	if err == nil {
		t.Fatal("Expected error for nil device")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Field != testDeviceParameter {
		t.Errorf("Expected param '%s', got '%s'", testDeviceParameter, valErr.Field)
	}
}

func testCreateThreadLocalCommandPoolSuccess(t *testing.T) {
	expectedPool := fakeCommandPool()
	mockDevice := fakeDevice()

	createCommandPoolFunc = func(device Device, createInfo *CommandPoolCreateInfo) (CommandPool, error) {
		if device != mockDevice {
			t.Errorf("Expected device %v, got %v", mockDevice, device)
		}
		if createInfo.QueueFamilyIndex != 1 {
			t.Errorf("Expected queue family index 1, got %d", createInfo.QueueFamilyIndex)
		}
		if createInfo.Flags != (CommandPoolCreateTransientBit | CommandPoolCreateResetCommandBufferBit) {
			t.Errorf("Expected correct flags, got %x", createInfo.Flags)
		}
		return expectedPool, nil
	}

	pool, err := CreateThreadLocalCommandPool(mockDevice, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pool == nil {
		t.Fatal("Expected pool, got nil")
	}
	if pool.Device != mockDevice {
		t.Errorf("Expected device %v, got %v", mockDevice, pool.Device)
	}
	if pool.CommandPool != expectedPool {
		t.Errorf("Expected command pool %v, got %v", expectedPool, pool.CommandPool)
	}
	if pool.CommandBuffers == nil {
		t.Error("CommandBuffers should be initialized to an empty slice")
	}
}

func testCreateThreadLocalCommandPoolError(t *testing.T) {
	mockDevice := fakeDevice()
	expectedErr := errors.New("mock error")

	createCommandPoolFunc = func(device Device, createInfo *CommandPoolCreateInfo) (CommandPool, error) {
		return nil, expectedErr
	}

	pool, err := CreateThreadLocalCommandPool(mockDevice, 0)
	if err != expectedErr {
		t.Fatalf("Expected error %v, got %v", expectedErr, err)
	}
	if pool != nil {
		t.Fatal("Expected nil pool on error")
	}
}
