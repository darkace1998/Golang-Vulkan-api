package vulkan

import (
	"errors"
	"fmt"
)

// VulkanError represents a structured Vulkan error with additional context
type VulkanError struct {
	Result    Result
	Operation string
	Details   string
}

// Error implements the error interface
func (e *VulkanError) Error() string {
	if e.Details != "" {
		return e.Operation + " failed: " + e.Result.Error() + " (" + e.Details + ")"
	}
	return e.Operation + " failed: " + e.Result.Error()
}

// Unwrap returns the underlying Result as an error for error unwrapping
func (e *VulkanError) Unwrap() error {
	return e.Result
}

// IsVulkanError checks if an error is a VulkanError
func IsVulkanError(err error) bool {
	_, ok := err.(*VulkanError)
	return ok
}

// NewVulkanError creates a new VulkanError
func NewVulkanError(result Result, operation string, details string) *VulkanError {
	return &VulkanError{
		Result:    result,
		Operation: operation,
		Details:   details,
	}
}

// IsErrorDeviceLost checks if an error indicates that the Vulkan device has been lost
// (VK_ERROR_DEVICE_LOST). It correctly unwraps nested errors.
func IsErrorDeviceLost(err error) bool {
	var vErr *VulkanError
	return errors.As(err, &vErr) && vErr.Result == ErrorDeviceLost
}

// IsErrorSurfaceLost checks if an error indicates that the Vulkan surface has been lost
// (VK_ERROR_SURFACE_LOST_KHR). It correctly unwraps nested errors.
func IsErrorSurfaceLost(err error) bool {
	var vErr *VulkanError
	return errors.As(err, &vErr) && vErr.Result == ErrorSurfaceLostKHR
}

// IsErrorOutOfDate checks if an error indicates that the Vulkan swapchain is out of date
// (VK_ERROR_OUT_OF_DATE_KHR). It correctly unwraps nested errors.
func IsErrorOutOfDate(err error) bool {
	var vErr *VulkanError
	return errors.As(err, &vErr) && vErr.Result == ErrorOutOfDateKHR
}

// ValidationError represents input validation errors
type ValidationError struct {
	Field  string
	Reason string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("vulkan validation error: %s %s", e.Field, e.Reason)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, reason string) *ValidationError {
	return &ValidationError{
		Field:  field,
		Reason: reason,
	}
}
