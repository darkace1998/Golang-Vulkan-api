//go:build darwin

package vulkan

import (
	"errors"
	"testing"
	"unsafe"
)

func TestCreateMetalSurfaceEXTValidation(t *testing.T) {
	instance := Instance(fakeHandle())
	var layer unsafe.Pointer = unsafe.Pointer(uintptr(1))

	// Test nil instance
	_, err := CreateMetalSurfaceEXT(nil, layer)
	if err == nil {
		t.Errorf("Expected validation error for nil instance")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) || valErr.Field != "instance" {
		t.Errorf("Expected ValidationError for instance, got %v", err)
	}

	// Test nil layer
	_, err = CreateMetalSurfaceEXT(instance, nil)
	if err == nil {
		t.Errorf("Expected validation error for nil layer")
	}
	if !errors.As(err, &valErr) || valErr.Field != "layer" {
		t.Errorf("Expected ValidationError for layer, got %v", err)
	}
}
