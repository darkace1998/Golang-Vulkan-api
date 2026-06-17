package vulkan

import (
	"errors"
	"testing"
)

// TestCreateTimelineSemaphoreValidation tests validation for CreateTimelineSemaphore
func TestCreateTimelineSemaphoreValidation(t *testing.T) {
	_, err := CreateTimelineSemaphore(nil, 0)
	if err == nil {
		t.Errorf("Expected error for nil device, got nil")
	} else {
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T: %v", err, err)
		}
	}
}
