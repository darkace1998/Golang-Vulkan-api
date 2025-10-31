package vulkan

import (
	"errors"
	"strings"
	"testing"
)

// TestStringSliceToCharArrayValidation tests the stringSliceToCharArray function
func TestStringSliceToCharArrayValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected bool // true if function should succeed
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: true, // returns nil for empty slice
		},
		{
			name:     "normal strings",
			input:    []string{"test1", "test2", "test3"},
			expected: true,
		},
		{
			name:     "single string",
			input:    []string{"test"},
			expected: true,
		},
		{
			name:     "large valid array",
			input:    make([]string, 100), // 100 empty strings
			expected: true,
		},
		{
			name:     "too large array",
			input:    make([]string, 20000), // exceeds maxAllowedSize
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize strings for non-empty test cases
			if len(tt.input) > 0 && tt.input[0] == "" {
				for i := range tt.input {
					tt.input[i] = "test"
				}
			}

			result := stringSliceToCharArray(tt.input)

			if tt.expected {
				if len(tt.input) == 0 {
					if result != nil {
						t.Errorf("Expected nil for empty slice, got non-nil")
					}
				} else {
					if result == nil {
						t.Errorf("Expected non-nil result for valid input, got nil")
					} else {
						// Clean up
						freeStringArray(result, len(tt.input))
					}
				}
			} else {
				if result != nil {
					t.Errorf("Expected nil for invalid input, got non-nil")
					freeStringArray(result, len(tt.input))
				}
			}
		})
	}
}

// TestCreateInstanceValidation tests input validation for CreateInstance
func TestCreateInstanceValidation(t *testing.T) {
	tests := []struct {
		name        string
		createInfo  *InstanceCreateInfo
		expectError bool
		errorType   string
	}{
		{
			name:        "nil createInfo",
			createInfo:  nil,
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name: "valid minimal createInfo",
			createInfo: &InstanceCreateInfo{
				ApplicationInfo:       nil,
				EnabledLayerNames:     []string{},
				EnabledExtensionNames: []string{},
			},
			expectError: false,
		},
		{
			name: "valid createInfo with application info",
			createInfo: &InstanceCreateInfo{
				ApplicationInfo: &ApplicationInfo{
					ApplicationName:    "TestApp",
					ApplicationVersion: MakeVersion(1, 0, 0),
					EngineName:         "TestEngine",
					EngineVersion:      MakeVersion(1, 0, 0),
					APIVersion:         Version13,
				},
				EnabledLayerNames:     []string{"VK_LAYER_KHRONOS_validation"},
				EnabledExtensionNames: []string{"VK_KHR_surface"},
			},
			expectError: false,
		},
		{
			name: "application name too long",
			createInfo: &InstanceCreateInfo{
				ApplicationInfo: &ApplicationInfo{
					ApplicationName: strings.Repeat("a", 300), // exceeds 256 chars
				},
			},
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name: "engine name too long",
			createInfo: &InstanceCreateInfo{
				ApplicationInfo: &ApplicationInfo{
					ApplicationName: "TestApp",
					EngineName:      strings.Repeat("a", 300), // exceeds 256 chars
				},
			},
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name: "too many layers",
			createInfo: &InstanceCreateInfo{
				EnabledLayerNames: make([]string, 100), // exceeds maxLayers
			},
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name: "too many extensions",
			createInfo: &InstanceCreateInfo{
				EnabledExtensionNames: make([]string, 300), // exceeds maxExtensions
			},
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name: "layer name too long",
			createInfo: &InstanceCreateInfo{
				EnabledLayerNames: []string{strings.Repeat("a", 300)}, // exceeds 256 chars
			},
			expectError: true,
			errorType:   "ValidationError",
		},
		{
			name: "extension name too long",
			createInfo: &InstanceCreateInfo{
				EnabledExtensionNames: []string{strings.Repeat("a", 300)}, // exceeds 256 chars
			},
			expectError: true,
			errorType:   "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize string slices with valid content for size tests
			if tt.createInfo != nil {
				if len(tt.createInfo.EnabledLayerNames) > 0 && tt.createInfo.EnabledLayerNames[0] == "" {
					for i := range tt.createInfo.EnabledLayerNames {
						tt.createInfo.EnabledLayerNames[i] = "test_layer"
					}
				}
				if len(tt.createInfo.EnabledExtensionNames) > 0 && tt.createInfo.EnabledExtensionNames[0] == "" {
					for i := range tt.createInfo.EnabledExtensionNames {
						tt.createInfo.EnabledExtensionNames[i] = "test_extension"
					}
				}
			}

			_, err := CreateInstance(tt.createInfo)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				// Check error type
				switch tt.errorType {
				case "ValidationError":
					var validationErr *ValidationError
					if !errors.As(err, &validationErr) {
						t.Errorf("Expected ValidationError, got %T: %v", err, err)
					}
				case "VulkanError":
					var vulkanErr *VulkanError
					if !errors.As(err, &vulkanErr) {
						t.Errorf("Expected VulkanError, got %T: %v", err, err)
					}
				}
			} else {
				// For valid inputs, we expect Vulkan-related errors since no Vulkan driver is installed
				// But we should not get validation errors
				if err != nil {
					var validationErr *ValidationError
					if errors.As(err, &validationErr) {
						t.Errorf("Got unexpected validation error for valid input: %v", err)
					}
					// VulkanError is expected since no Vulkan driver is available
				}
			}
		})
	}
}

// TestVulkanErrorType tests the VulkanError type
func TestVulkanErrorType(t *testing.T) {
	err := NewVulkanError(ErrorInitializationFailed, "TestOperation", "test details")

	if err.Result != ErrorInitializationFailed {
		t.Errorf("Expected Result %v, got %v", ErrorInitializationFailed, err.Result)
	}

	if err.Operation != "TestOperation" {
		t.Errorf("Expected Operation 'TestOperation', got '%s'", err.Operation)
	}

	if err.Details != "test details" {
		t.Errorf("Expected Details 'test details', got '%s'", err.Details)
	}

	expectedMsg := "TestOperation failed: VK_ERROR_INITIALIZATION_FAILED (test details)"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test Unwrap
	unwrapped := err.Unwrap()
	if unwrapped != ErrorInitializationFailed {
		t.Errorf("Expected unwrapped error %v, got %v", ErrorInitializationFailed, unwrapped)
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

// TestVersionHelpers tests version utility functions
func TestVersionHelpers(t *testing.T) {
	version := MakeVersion(1, 2, 3)

	if version.Major() != 1 {
		t.Errorf("Expected major version 1, got %d", version.Major())
	}

	if version.Minor() != 2 {
		t.Errorf("Expected minor version 2, got %d", version.Minor())
	}

	if version.Patch() != 3 {
		t.Errorf("Expected patch version 3, got %d", version.Patch())
	}
}

// TestResultHelpers tests Result helper functions
func TestResultHelpers(t *testing.T) {
	// Test success result
	if !Success.IsSuccess() {
		t.Errorf("Success should return true for IsSuccess()")
	}

	// Test error result
	if ErrorInitializationFailed.IsSuccess() {
		t.Errorf("Error result should return false for IsSuccess()")
	}

	// Test error message
	expected := "VK_ERROR_INITIALIZATION_FAILED"
	if ErrorInitializationFailed.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, ErrorInitializationFailed.Error())
	}
}

// BenchmarkStringSliceToCharArray benchmarks the string slice conversion
func BenchmarkStringSliceToCharArray(b *testing.B) {
	testSlice := []string{"layer1", "layer2", "layer3", "layer4", "layer5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := stringSliceToCharArray(testSlice)
		if result != nil {
			freeStringArray(result, len(testSlice))
		}
	}
}
