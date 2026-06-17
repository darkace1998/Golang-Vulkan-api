package vulkan

import (
	"errors"
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
