package vulkan

import (
	"errors"
	"sync"
	"testing"
)

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyInstanceNilArgs tests that DestroyInstance handles nil gracefully
func TestDestroyInstanceNilArgs(t *testing.T) {
	DestroyInstance(nil) // Should not panic
}

// TestEnumeratePhysicalDevicesValidation tests nil parameter validation
func TestEnumeratePhysicalDevicesValidation(t *testing.T) {
	_, err := EnumeratePhysicalDevices(nil)
	if err == nil {
		t.Fatal("Expected error for nil instance")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Parameter != "instance" {
		t.Errorf("Expected param 'instance', got '%s'", valErr.Parameter)
	}
}

// ============================================================================
// Video Function Loader Thread-Safety Tests
// ============================================================================

// TestVideoFunctionLoaderReset tests that Reset functions work
func TestVideoFunctionLoaderReset(t *testing.T) {
	// Reset both so we start clean
	ResetVideoInstanceFunctions()
	ResetVideoDeviceFunctions()

	// Verify that the loaded flags are reset
	// After reset, calling Load with nil instance will try to load but should fail
	// (no actual Vulkan instance)
	// We just verify no panics occur
	ResetVideoInstanceFunctions()
	ResetVideoDeviceFunctions()
}

// TestVideoFunctionLoaderConcurrent tests thread-safety of LoadVideoInstanceFunctions
func TestVideoFunctionLoaderConcurrent(t *testing.T) {
	ResetVideoInstanceFunctions()

	const goroutines = 50
	var wg sync.WaitGroup
	results := make([]bool, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Pass nil instance - will fail but should not race
			results[idx] = LoadVideoInstanceFunctions(nil)
		}(i)
	}

	wg.Wait()

	// All goroutines should get the same result (sync.Once ensures single execution)
	first := results[0]
	for i, r := range results {
		if r != first {
			t.Errorf("Goroutine %d got result %v, expected %v (inconsistent due to race?)", i, r, first)
		}
	}

	ResetVideoInstanceFunctions() // cleanup
}

// TestVideoDeviceFunctionLoaderConcurrent tests thread-safety of LoadVideoDeviceFunctions
func TestVideoDeviceFunctionLoaderConcurrent(t *testing.T) {
	ResetVideoDeviceFunctions()

	const goroutines = 50
	var wg sync.WaitGroup
	results := make([]bool, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = LoadVideoDeviceFunctions(nil)
		}(i)
	}

	wg.Wait()

	first := results[0]
	for i, r := range results {
		if r != first {
			t.Errorf("Goroutine %d got result %v, expected %v", i, r, first)
		}
	}

	ResetVideoDeviceFunctions() // cleanup
}
