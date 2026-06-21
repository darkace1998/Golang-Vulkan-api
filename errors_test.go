package vulkan

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewVulkanError(t *testing.T) {
	tests := []struct {
		name      string
		result    Result
		operation string
		details   string
	}{
		{
			name:      "all fields provided",
			result:    ErrorOutOfHostMemory,
			operation: "vkCreateInstance",
			details:   "failed to allocate memory",
		},
		{
			name:      "empty details",
			result:    ErrorDeviceLost,
			operation: "vkQueueSubmit",
			details:   "",
		},
		{
			name:      "empty operation",
			result:    ErrorInitializationFailed,
			operation: "",
			details:   "driver not found",
		},
		{
			name:      "success result",
			result:    Success,
			operation: "vkAllocateMemory",
			details:   "allocated successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewVulkanError(tt.result, tt.operation, tt.details)

			if err == nil {
				t.Fatal("Expected NewVulkanError to return a non-nil error")
			}

			if err.Result != tt.result {
				t.Errorf("Expected Result %v, got %v", tt.result, err.Result)
			}
			if err.Operation != tt.operation {
				t.Errorf("Expected Operation %q, got %q", tt.operation, err.Operation)
			}
			if err.Details != tt.details {
				t.Errorf("Expected Details %q, got %q", tt.details, err.Details)
			}
		})
	}
}

func TestVulkanError(t *testing.T) {
	result := ErrorOutOfHostMemory
	op := "vkCreateInstance"
	details := "failed to allocate memory"

	err := NewVulkanError(result, op, details)

	// Test Error() with details
	expectedErrorWithDetails := "vkCreateInstance failed: VK_ERROR_OUT_OF_HOST_MEMORY (failed to allocate memory)"
	if err.Error() != expectedErrorWithDetails {
		t.Errorf("Expected Error() string %q, got %q", expectedErrorWithDetails, err.Error())
	}

	// Test Error() without details
	errNoDetails := NewVulkanError(result, op, "")
	expectedErrorNoDetails := "vkCreateInstance failed: VK_ERROR_OUT_OF_HOST_MEMORY"
	if errNoDetails.Error() != expectedErrorNoDetails {
		t.Errorf("Expected Error() string %q, got %q", expectedErrorNoDetails, errNoDetails.Error())
	}

	// Test Unwrap()
	if unwrapped := err.Unwrap(); unwrapped != result {
		t.Errorf("Expected Unwrap() to return %v, got %v", result, unwrapped)
	}
}

func TestIsVulkanError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "vulkan error with success",
			err:  NewVulkanError(Success, "test", ""),
			want: true,
		},
		{
			name: "standard error",
			err:  errors.New("standard error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "validation error",
			err:  NewValidationError("field", "reason"),
			want: false,
		},
		{
			name: "vulkan error with failure",
			err:  NewVulkanError(ErrorOutOfHostMemory, "vkCreateInstance", ""),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVulkanError(tt.err); got != tt.want {
				t.Errorf("IsVulkanError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		reason   string
		expected string
	}{
		{
			name:     "normal validation error",
			field:    "imageInfo.format",
			reason:   "unsupported format",
			expected: "vulkan validation error: imageInfo.format unsupported format",
		},
		{
			name:     "empty field",
			field:    "",
			reason:   "missing field",
			expected: "vulkan validation error:  missing field",
		},
		{
			name:     "empty reason",
			field:    "flags",
			reason:   "",
			expected: "vulkan validation error: flags ",
		},
		{
			name:     "empty field and reason",
			field:    "",
			reason:   "",
			expected: "vulkan validation error:  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.field, tt.reason)

			if err.Field != tt.field {
				t.Errorf("Expected Field %v, got %v", tt.field, err.Field)
			}
			if err.Reason != tt.reason {
				t.Errorf("Expected Reason %v, got %v", tt.reason, err.Reason)
			}

			if err.Error() != tt.expected {
				t.Errorf("Expected Error() string %q, got %q", tt.expected, err.Error())
			}
		})
	}
}

func TestVulkanErrorUnwrapIsAs(t *testing.T) {
	// Setup
	result := ErrorOutOfDeviceMemory
	op := "vkAllocateMemory"
	details := "failed to allocate 1024 bytes"

	vkErr := NewVulkanError(result, op, details)

	// Test errors.Is
	if !errors.Is(vkErr, result) {
		t.Errorf("Expected errors.Is to return true for %v", result)
	}

	if errors.Is(vkErr, ErrorOutOfHostMemory) {
		t.Errorf("Expected errors.Is to return false for %v", ErrorOutOfHostMemory)
	}

	// Test errors.As
	var extractedVkErr *VulkanError
	if !errors.As(vkErr, &extractedVkErr) {
		t.Errorf("Expected errors.As to return true and extract VulkanError")
	}

	if extractedVkErr.Result != result || extractedVkErr.Operation != op || extractedVkErr.Details != details {
		t.Errorf("Expected extracted error to match original. Got %+v", extractedVkErr)
	}

	// Test wrapping in fmt.Errorf
	wrappedErr := fmt.Errorf("allocation failed: %w", vkErr)

	if !errors.Is(wrappedErr, result) {
		t.Errorf("Expected wrapped error to match underlying result %v via errors.Is", result)
	}

	var extractedWrappedVkErr *VulkanError
	if !errors.As(wrappedErr, &extractedWrappedVkErr) {
		t.Errorf("Expected errors.As to extract VulkanError from wrapped error")
	}
}

func TestValidationErrorAs(t *testing.T) {
	field := "bufferInfo.size"
	reason := "must be greater than 0"
	valErr := NewValidationError(field, reason)

	// Test errors.As
	var extractedValErr *ValidationError
	if !errors.As(valErr, &extractedValErr) {
		t.Errorf("Expected errors.As to extract ValidationError")
	}

	if extractedValErr.Field != field || extractedValErr.Reason != reason {
		t.Errorf("Expected extracted error to match original. Got %+v", extractedValErr)
	}

	// Test wrapping in fmt.Errorf
	wrappedErr := fmt.Errorf("parameter validation failed: %w", valErr)

	var extractedWrappedValErr *ValidationError
	if !errors.As(wrappedErr, &extractedWrappedValErr) {
		t.Errorf("Expected errors.As to extract ValidationError from wrapped error")
	}
}

func TestVulkanErrorStringGeneration(t *testing.T) {
	// 1. Without details
	result := ErrorOutOfHostMemory
	op := "vkAllocateMemory"
	errNoDetails := NewVulkanError(result, op, "")
	expectedNoDetails := "vkAllocateMemory failed: VK_ERROR_OUT_OF_HOST_MEMORY"

	if errNoDetails.Error() != expectedNoDetails {
		t.Errorf("Expected string %q, got %q", expectedNoDetails, errNoDetails.Error())
	}

	// 2. With details
	details := "system out of memory"
	errWithDetails := NewVulkanError(result, op, details)
	expectedWithDetails := "vkAllocateMemory failed: VK_ERROR_OUT_OF_HOST_MEMORY (system out of memory)"

	if errWithDetails.Error() != expectedWithDetails {
		t.Errorf("Expected string %q, got %q", expectedWithDetails, errWithDetails.Error())
	}
}

func TestValidationErrorStringGeneration(t *testing.T) {
	err := NewValidationError("image.format", "unsupported format")
	expected := "vulkan validation error: image.format unsupported format"
	if err.Error() != expected {
		t.Errorf("Expected string %q, got %q", expected, err.Error())
	}
}

func TestVulkanError_errorsUnwrap(t *testing.T) {
	// Setup
	result := ErrorOutOfDeviceMemory
	op := "vkAllocateMemory"
	details := "failed to allocate 1024 bytes"

	vkErr := NewVulkanError(result, op, details)

	// Test standard library unwrapping natively
	unwrapped := errors.Unwrap(vkErr)
	if unwrapped != result {
		t.Errorf("Expected errors.Unwrap to yield Result %v, but got %v", result, unwrapped)
	}

	// Test standard library unwrapping when wrapped with fmt.Errorf
	wrappedErr := fmt.Errorf("allocation failed: %w", vkErr)
	unwrappedFromFmt := errors.Unwrap(wrappedErr)
	if unwrappedFromFmt != vkErr {
		t.Errorf("Expected errors.Unwrap of fmt.Errorf to yield original VulkanError, but got %v", unwrappedFromFmt)
	}

	// Unwrapping twice should get the result
	unwrappedTwice := errors.Unwrap(unwrappedFromFmt)
	if unwrappedTwice != result {
		t.Errorf("Expected errors.Unwrap twice to yield Result %v, but got %v", result, unwrappedTwice)
	}
}
