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
