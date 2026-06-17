package vulkan

import (
	"errors"
	"fmt"
	"testing"
)

func TestVulkanError(t *testing.T) {
	// Test NewVulkanError
	result := ErrorOutOfHostMemory
	op := "vkCreateInstance"
	details := "failed to allocate memory"

	err := NewVulkanError(result, op, details)

	if err.Result != result {
		t.Errorf("Expected Result %v, got %v", result, err.Result)
	}
	if err.Operation != op {
		t.Errorf("Expected Operation %v, got %v", op, err.Operation)
	}
	if err.Details != details {
		t.Errorf("Expected Details %v, got %v", details, err.Details)
	}

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
	// True case
	vkErr := NewVulkanError(Success, "test", "")
	if !IsVulkanError(vkErr) {
		t.Errorf("Expected IsVulkanError to return true for VulkanError")
	}

	// False case - standard error
	stdErr := errors.New("standard error")
	if IsVulkanError(stdErr) {
		t.Errorf("Expected IsVulkanError to return false for standard error")
	}

	// False case - nil error
	if IsVulkanError(nil) {
		t.Errorf("Expected IsVulkanError to return false for nil error")
	}

	// False case - other custom error
	valErr := NewValidationError("field", "reason")
	if IsVulkanError(valErr) {
		t.Errorf("Expected IsVulkanError to return false for ValidationError")
	}
}

func TestValidationError(t *testing.T) {
	// Test NewValidationError
	field := "imageInfo.format"
	reason := "unsupported format"

	err := NewValidationError(field, reason)

	if err.Field != field {
		t.Errorf("Expected Field %v, got %v", field, err.Field)
	}
	if err.Reason != reason {
		t.Errorf("Expected Reason %v, got %v", reason, err.Reason)
	}

	// Test Error()
	expectedError := "vulkan validation error: imageInfo.format unsupported format"
	if err.Error() != expectedError {
		t.Errorf("Expected Error() string %q, got %q", expectedError, err.Error())
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
