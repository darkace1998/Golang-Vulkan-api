package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyDeviceNilArgs tests that DestroyDevice handles nil gracefully
func TestDestroyDeviceNilArgs(t *testing.T) {
	DestroyDevice(nil) // Should not panic
}

// TestQueueWaitIdleValidation tests nil parameter validation
func TestQueueWaitIdleValidation(t *testing.T) {
	err := QueueWaitIdle(nil)
	if err == nil {
		t.Fatal("Expected error for nil queue")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Parameter != "queue" {
		t.Errorf("Expected param 'queue', got '%s'", valErr.Parameter)
	}
}

// TestDeviceWaitIdleValidation tests nil parameter validation
func TestDeviceWaitIdleValidation(t *testing.T) {
	err := DeviceWaitIdle(nil)
	if err == nil {
		t.Fatalf("Expected error for nil %s", testDeviceParameter)
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if valErr.Parameter != testDeviceParameter {
		t.Errorf("Expected param '%s', got '%s'", testDeviceParameter, valErr.Parameter)
	}
}
