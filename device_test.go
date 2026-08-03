package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyDeviceNilArgs tests that DestroyDevice handles nil gracefully
func TestDestroyDeviceNilArgs(t *testing.T) {
	DestroyDevice(nil) // Should not panic
}

// TestQueueWaitIdleValidation tests nil parameter validation
func TestQueueWaitIdleValidation(t *testing.T) {
	err := QueueWaitIdle(nil)
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

// TestDeviceWaitIdleValidation tests nil parameter validation
func TestDeviceWaitIdleValidation(t *testing.T) {
	err := DeviceWaitIdle(nil)
	if err == nil {
		t.Fatalf("Expected error for nil %s", testDeviceParameter)
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Field != testDeviceParameter {
		t.Errorf("Expected param '%s', got '%s'", testDeviceParameter, valErr.Field)
	}
}

// TestValidateDeviceCreateInfo tests the validateDeviceCreateInfo function.
func TestValidateDeviceCreateInfo(t *testing.T) {
	physicalDevice := fakePhysicalDevice()

	tests := []struct {
		name           string
		physicalDevice PhysicalDevice
		createInfo     *DeviceCreateInfo
		expectErr      bool
		errField       string
	}{
		{
			name:           "Valid CreateInfo",
			physicalDevice: physicalDevice,
			createInfo: &DeviceCreateInfo{
				QueueCreateInfos: []DeviceQueueCreateInfo{
					{QueuePriorities: []float32{1.0}},
				},
			},
			expectErr: false,
		},
		{
			name:           "Empty QueueCreateInfos",
			physicalDevice: physicalDevice,
			createInfo:     &DeviceCreateInfo{},
			expectErr:      true,
			errField:       "QueueCreateInfos",
		},
		{
			name:           "Nil PhysicalDevice",
			physicalDevice: nil,
			createInfo:     &DeviceCreateInfo{},
			expectErr:      true,
			errField:       "physicalDevice",
		},
		{
			name:           "Nil CreateInfo",
			physicalDevice: physicalDevice,
			createInfo:     nil,
			expectErr:      true,
			errField:       testCreateInfoParameter,
		},
		{
			name:           "Invalid QueueCreateInfos",
			physicalDevice: physicalDevice,
			createInfo: &DeviceCreateInfo{
				QueueCreateInfos: []DeviceQueueCreateInfo{
					{
						QueuePriorities: []float32{1.5}, // Invalid priority (> 1.0)
					},
				},
			},
			expectErr: true,
			errField:  "QueueCreateInfos",
		},
		{
			name:           "Invalid Layers",
			physicalDevice: physicalDevice,
			createInfo: &DeviceCreateInfo{
				QueueCreateInfos: []DeviceQueueCreateInfo{
					{QueuePriorities: []float32{1.0}},
				},
				EnabledLayerNames: []string{
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // 257 chars long string
				},
			},
			expectErr: true,
			errField:  "EnabledLayerNames",
		},
		{
			name:           "Invalid Extensions",
			physicalDevice: physicalDevice,
			createInfo: &DeviceCreateInfo{
				QueueCreateInfos: []DeviceQueueCreateInfo{
					{QueuePriorities: []float32{1.0}},
				},
				EnabledExtensionNames: []string{
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // 257 chars long string
				},
			},
			expectErr: true,
			errField:  "EnabledExtensionNames",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeviceCreateInfo(tt.physicalDevice, tt.createInfo)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else {
					var valErr *ValidationError
					if !errors.As(err, &valErr) {
						t.Errorf("Expected ValidationError, got %T: %v", err, err)
					} else if valErr.Field != tt.errField {
						t.Errorf("Expected error on field '%s', got '%s'", tt.errField, valErr.Field)
					}
				}
			} else if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
