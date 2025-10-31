package vulkan

import (
	"errors"
	"strings"
	"testing"
)

// Pure Go tests that don't require CGO compilation

// TestVulkanErrorTypeBasic tests the VulkanError type with basic Result values
func TestVulkanErrorTypeBasic(t *testing.T) {
	// Use a simple Result value that doesn't require CGO constants
	testResult := Result(-1) // Generic error
	err := NewVulkanError(testResult, "TestOperation", "test details")

	if err.Result != testResult {
		t.Errorf("Expected Result %v, got %v", testResult, err.Result)
	}

	if err.Operation != "TestOperation" {
		t.Errorf("Expected Operation 'TestOperation', got '%s'", err.Operation)
	}

	if err.Details != "test details" {
		t.Errorf("Expected Details 'test details', got '%s'", err.Details)
	}

	// Test error message contains expected components
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "TestOperation failed") {
		t.Errorf("Error message should contain operation name, got '%s'", errorMsg)
	}
	if !strings.Contains(errorMsg, "test details") {
		t.Errorf("Error message should contain details, got '%s'", errorMsg)
	}

	// Test Unwrap
	unwrapped := err.Unwrap()
	if unwrapped != testResult {
		t.Errorf("Expected unwrapped error %v, got %v", testResult, unwrapped)
	}

	// Test IsVulkanError
	if !IsVulkanError(err) {
		t.Errorf("IsVulkanError should return true for VulkanError")
	}

	// Test with regular error
	regularErr := errors.New("regular error")
	if IsVulkanError(regularErr) {
		t.Errorf("IsVulkanError should return false for regular error")
	}
}

// TestVulkanErrorWithoutDetails tests VulkanError without details
func TestVulkanErrorWithoutDetails(t *testing.T) {
	testResult := Result(-2)
	err := NewVulkanError(testResult, "TestOperation", "")

	errorMsg := err.Error()
	expectedPrefix := "TestOperation failed:"
	if !strings.HasPrefix(errorMsg, expectedPrefix) {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedPrefix, errorMsg)
	}

	// Should not contain empty parentheses when no details
	if strings.Contains(errorMsg, "()") {
		t.Errorf("Error message should not contain empty parentheses, got '%s'", errorMsg)
	}
}

// TestValidationErrorType tests the ValidationError type
func TestValidationErrorType(t *testing.T) {
	err := NewValidationError("testParam", "test message")

	if err.Parameter != "testParam" {
		t.Errorf("Expected Parameter 'testParam', got '%s'", err.Parameter)
	}

	if err.Message != "test message" {
		t.Errorf("Expected Message 'test message', got '%s'", err.Message)
	}

	expectedMsg := "validation error for parameter 'testParam': test message"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestValidationPatterns tests common validation patterns
func TestValidationPatterns(t *testing.T) {
	tests := []struct {
		name        string
		validate    func() error
		expectError bool
	}{
		{
			name: "valid string length",
			validate: func() error {
				testStr := "valid"
				if len(testStr) > 256 {
					return NewValidationError("testStr", "exceeds maximum length")
				}
				return nil
			},
			expectError: false,
		},
		{
			name: "invalid string length",
			validate: func() error {
				testStr := strings.Repeat("a", 300)
				if len(testStr) > 256 {
					return NewValidationError("testStr", "exceeds maximum length")
				}
				return nil
			},
			expectError: true,
		},
		{
			name: "valid array size",
			validate: func() error {
				testArray := make([]string, 10)
				const maxSize = 64
				if len(testArray) > maxSize {
					return NewValidationError("testArray", "exceeds maximum size")
				}
				return nil
			},
			expectError: false,
		},
		{
			name: "invalid array size",
			validate: func() error {
				testArray := make([]string, 100)
				const maxSize = 64
				if len(testArray) > maxSize {
					return NewValidationError("testArray", "exceeds maximum size")
				}
				return nil
			},
			expectError: true,
		},
		{
			name: "nil parameter validation",
			validate: func() error {
				var testPtr *string = nil
				if testPtr == nil {
					return NewValidationError("testPtr", "cannot be nil")
				}
				return nil
			},
			expectError: true,
		},
		{
			name: "range validation",
			validate: func() error {
				testValue := 1.5
				if testValue < 0.0 || testValue > 1.0 {
					return NewValidationError("testValue", "must be between 0.0 and 1.0")
				}
				return nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validate()
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectError && err != nil {
				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				}
			}
		})
	}
}

// TestErrorWrapping tests error wrapping functionality
func TestErrorWrapping(t *testing.T) {
	testResult := Result(-3)
	vulkanErr := NewVulkanError(testResult, "TestOp", "test details")

	// Test errors.Is
	if !errors.Is(vulkanErr, testResult) {
		t.Errorf("errors.Is should return true for wrapped Result")
	}

	// Test errors.As
	var resultErr Result
	if !errors.As(vulkanErr, &resultErr) {
		t.Errorf("errors.As should extract Result from VulkanError")
	}
	if resultErr != testResult {
		t.Errorf("Expected Result %v, got %v", testResult, resultErr)
	}
}

// TestErrorTypeDistinction tests that different error types are distinct
func TestErrorTypeDistinction(t *testing.T) {
	vulkanErr := NewVulkanError(Result(-1), "VulkanOp", "details")
	validationErr := NewValidationError("param", "invalid")

	// Test that we can distinguish between error types
	if IsVulkanError(validationErr) {
		t.Errorf("ValidationError should not be identified as VulkanError")
	}

	// Test error.As with different types
	var vErr *VulkanError
	var valErr *ValidationError

	if !errors.As(vulkanErr, &vErr) {
		t.Errorf("Should be able to extract VulkanError")
	}
	if errors.As(vulkanErr, &valErr) {
		t.Errorf("Should not be able to extract ValidationError from VulkanError")
	}

	if errors.As(validationErr, &vErr) {
		t.Errorf("Should not be able to extract VulkanError from ValidationError")
	}
	if !errors.As(validationErr, &valErr) {
		t.Errorf("Should be able to extract ValidationError")
	}
}

// BenchmarkErrorCreation benchmarks error creation
func BenchmarkErrorCreation(b *testing.B) {
	b.Run("VulkanError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewVulkanError(Result(-1), "TestOp", "details")
		}
	})

	b.Run("ValidationError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewValidationError("param", "message")
		}
	})
}
