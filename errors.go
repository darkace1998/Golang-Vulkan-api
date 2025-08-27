package vulkan

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

// ValidationError represents input validation errors
type ValidationError struct {
	Parameter string
	Message   string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return "validation error for parameter '" + e.Parameter + "': " + e.Message
}

// NewValidationError creates a new ValidationError
func NewValidationError(parameter, message string) *ValidationError {
	return &ValidationError{
		Parameter: parameter,
		Message:   message,
	}
}