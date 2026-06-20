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

func TestResultHelpers(t *testing.T) {
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

func TestBoolConversion(t *testing.T) {
	t.Run("ToBool", func(t *testing.T) {
		if !Bool32(1).ToBool() {
			t.Error("Expected Bool32(1).ToBool() to be true")
		}
		if Bool32(0).ToBool() {
			t.Error("Expected Bool32(0).ToBool() to be false")
		}
	})
	t.Run("FromBool", func(t *testing.T) {
		if FromBool(true) != 1 {
			t.Error("Expected FromBool(true) to be 1")
		}
		if FromBool(false) != 0 {
			t.Error("Expected FromBool(false) to be 0")
		}
	})
}

func TestVersion(t *testing.T) {
    v := MakeVersion(1, 2, 3)
    if v.Major() != 1 {
        t.Errorf("Expected Major() 1, got %d", v.Major())
    }
    if v.Minor() != 2 {
        t.Errorf("Expected Minor() 2, got %d", v.Minor())
    }
    if v.Patch() != 3 {
        t.Errorf("Expected Patch() 3, got %d", v.Patch())
    }
}
