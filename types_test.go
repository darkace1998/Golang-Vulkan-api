package vulkan

import "testing"

func TestResult_Error(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		expected string
	}{
		{"Success", Success, "VK_SUCCESS"},
		{"NotReady", NotReady, "VK_NOT_READY"},
		{"Timeout", Timeout, "VK_TIMEOUT"},
		{"EventSet", EventSet, "VK_EVENT_SET"},
		{"EventReset", EventReset, "VK_EVENT_RESET"},
		{"Incomplete", Incomplete, "VK_INCOMPLETE"},
		{"ErrorOutOfHostMemory", ErrorOutOfHostMemory, "VK_ERROR_OUT_OF_HOST_MEMORY"},
		{"ErrorOutOfDeviceMemory", ErrorOutOfDeviceMemory, "VK_ERROR_OUT_OF_DEVICE_MEMORY"},
		{"ErrorInitializationFailed", ErrorInitializationFailed, "VK_ERROR_INITIALIZATION_FAILED"},
		{"ErrorDeviceLost", ErrorDeviceLost, "VK_ERROR_DEVICE_LOST"},
		{"ErrorMemoryMapFailed", ErrorMemoryMapFailed, "VK_ERROR_MEMORY_MAP_FAILED"},
		{"ErrorLayerNotPresent", ErrorLayerNotPresent, "VK_ERROR_LAYER_NOT_PRESENT"},
		{"ErrorExtensionNotPresent", ErrorExtensionNotPresent, "VK_ERROR_EXTENSION_NOT_PRESENT"},
		{"ErrorFeatureNotPresent", ErrorFeatureNotPresent, "VK_ERROR_FEATURE_NOT_PRESENT"},
		{"ErrorIncompatibleDriver", ErrorIncompatibleDriver, "VK_ERROR_INCOMPATIBLE_DRIVER"},
		{"ErrorTooManyObjects", ErrorTooManyObjects, "VK_ERROR_TOO_MANY_OBJECTS"},
		{"ErrorFormatNotSupported", ErrorFormatNotSupported, "VK_ERROR_FORMAT_NOT_SUPPORTED"},
		{"ErrorFragmentedPool", ErrorFragmentedPool, "VK_ERROR_FRAGMENTED_POOL"},
		{"ErrorUnknown", ErrorUnknown, "VK_ERROR_UNKNOWN"},
		{"ErrorOutOfPoolMemory", ErrorOutOfPoolMemory, "VK_ERROR_OUT_OF_POOL_MEMORY"},
		{"ErrorInvalidExternalHandle", ErrorInvalidExternalHandle, "VK_ERROR_INVALID_EXTERNAL_HANDLE"},
		{"ErrorFragmentation", ErrorFragmentation, "VK_ERROR_FRAGMENTATION"},
		{"ErrorInvalidOpaqueCaptureAddress", ErrorInvalidOpaqueCaptureAddress, "VK_ERROR_INVALID_OPAQUE_CAPTURE_ADDRESS"},
		{"ErrorSurfaceLostKHR", ErrorSurfaceLostKHR, "VK_ERROR_SURFACE_LOST_KHR"},
		{"ErrorNativeWindowInUseKHR", ErrorNativeWindowInUseKHR, "VK_ERROR_NATIVE_WINDOW_IN_USE_KHR"},
		{"SuboptimalKHR", SuboptimalKHR, "VK_SUBOPTIMAL_KHR"},
		{"ErrorOutOfDateKHR", ErrorOutOfDateKHR, "VK_ERROR_OUT_OF_DATE_KHR"},
		{"ErrorIncompatibleDisplayKHR", ErrorIncompatibleDisplayKHR, "VK_ERROR_INCOMPATIBLE_DISPLAY_KHR"},
		{"ErrorValidationFailedEXT", ErrorValidationFailedEXT, "VK_ERROR_VALIDATION_FAILED_EXT"},
		{"ErrorInvalidShaderNV", ErrorInvalidShaderNV, "VK_ERROR_INVALID_SHADER_NV"},
		{"ErrorInvalidDrmFormatModifierPlaneLayoutEXT", ErrorInvalidDrmFormatModifierPlaneLayoutEXT, "VK_ERROR_INVALID_DRM_FORMAT_MODIFIER_PLANE_LAYOUT_EXT"},
		{"ErrorNotPermittedEXT", ErrorNotPermittedEXT, "VK_ERROR_NOT_PERMITTED_EXT"},
		{"ErrorFullScreenExclusiveModeLostEXT", ErrorFullScreenExclusiveModeLostEXT, "VK_ERROR_FULL_SCREEN_EXCLUSIVE_MODE_LOST_EXT"},
		{"ThreadIdleKHR", ThreadIdleKHR, "VK_THREAD_IDLE_KHR"},
		{"ThreadDoneKHR", ThreadDoneKHR, "VK_THREAD_DONE_KHR"},
		{"OperationDeferredKHR", OperationDeferredKHR, "VK_OPERATION_DEFERRED_KHR"},
		{"OperationNotDeferredKHR", OperationNotDeferredKHR, "VK_OPERATION_NOT_DEFERRED_KHR"},
		{"PipelineCompileRequiredEXT", PipelineCompileRequiredEXT, "VK_PIPELINE_COMPILE_REQUIRED_EXT"},
		{"Unknown Code", Result(9999), "Unknown Vulkan error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.result.Error(); got != tc.expected {
				t.Errorf("Result.Error() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestResultHelpersTypes(t *testing.T) {
	t.Run("IsError", func(t *testing.T) {
		tests := []struct {
			name     string
			result   Result
			expected bool
		}{
			{"Negative Result (Error)", ErrorOutOfHostMemory, true},
			{"Zero Result (Success)", Success, false},
			{"Positive Result (Incomplete)", Incomplete, false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.result.IsError() != tc.expected {
					t.Errorf("Expected IsError() to be %v for %v, got %v", tc.expected, tc.result, tc.result.IsError())
				}
			})
		}
	})

	t.Run("IsSuccess", func(t *testing.T) {
		tests := []struct {
			name     string
			result   Result
			expected bool
		}{
			{"Zero Result (Success)", Success, true},
			{"Positive Result (Incomplete)", Incomplete, true},
			{"Negative Result (Error)", ErrorOutOfHostMemory, false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.result.IsSuccess() != tc.expected {
					t.Errorf("Expected IsSuccess() to be %v for %v, got %v", tc.expected, tc.result, tc.result.IsSuccess())
				}
			})
		}
	})
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name  string
		major uint32
		minor uint32
		patch uint32
	}{
		{
			name:  "Standard 1.2.3",
			major: 1,
			minor: 2,
			patch: 3,
		},
		{
			name:  "Zero values",
			major: 0,
			minor: 0,
			patch: 0,
		},
		{
			name:  "Max values (127, 1023, 4095)",
			major: 127,
			minor: 1023,
			patch: 4095,
		},
		{
			name:  "Vulkan 1.0.0",
			major: 1,
			minor: 0,
			patch: 0,
		},
		{
			name:  "Vulkan 1.1.0",
			major: 1,
			minor: 1,
			patch: 0,
		},
		{
			name:  "Vulkan 1.2.0",
			major: 1,
			minor: 2,
			patch: 0,
		},
		{
			name:  "Vulkan 1.3.0",
			major: 1,
			minor: 3,
			patch: 0,
		},
		{
			name:  "Mixed values",
			major: 42,
			minor: 512,
			patch: 2048,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := MakeVersion(tc.major, tc.minor, tc.patch)
			if v.Major() != tc.major {
				t.Errorf("Expected Major() %d, got %d", tc.major, v.Major())
			}
			if v.Minor() != tc.minor {
				t.Errorf("Expected Minor() %d, got %d", tc.minor, v.Minor())
			}
			if v.Patch() != tc.patch {
				t.Errorf("Expected Patch() %d, got %d", tc.patch, v.Patch())
			}
		})
	}
}

func TestBoolConversion(t *testing.T) {
	t.Run("ToBool", func(t *testing.T) {
		tests := []struct {
			name     string
			input    Bool32
			expected bool
		}{
			{"True", True, true},
			{"False", False, false},
			{"Other Non-zero", Bool32(2), false},
			{"Other Non-zero 2", Bool32(100), false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.input.ToBool(); got != tt.expected {
					t.Errorf("Bool32.ToBool() = %v, want %v", got, tt.expected)
				}
			})
		}
	})

	t.Run("FromBool", func(t *testing.T) {
		tests := []struct {
			name     string
			input    bool
			expected Bool32
		}{
			{"true", true, True},
			{"false", false, False},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := FromBool(tt.input); got != tt.expected {
					t.Errorf("FromBool() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
}
