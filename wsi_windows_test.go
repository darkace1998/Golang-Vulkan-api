//go:build windows

package vulkan

import (
	"errors"
	"testing"
	"unsafe"
)

func TestCreateWin32SurfaceKHRValidation(t *testing.T) {
	instance := Instance(fakeHandle())
	var hinstance unsafe.Pointer = unsafe.Pointer(uintptr(1))
	var hwnd unsafe.Pointer = unsafe.Pointer(uintptr(2))

	// Test nil instance
	_, err := CreateWin32SurfaceKHR(nil, hinstance, hwnd)
	if err == nil {
		t.Errorf("Expected validation error for nil instance")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) || valErr.Field != "instance" {
		t.Errorf("Expected ValidationError for instance, got %v", err)
	}

	// Test nil hinstance
	_, err = CreateWin32SurfaceKHR(instance, nil, hwnd)
	if err == nil {
		t.Errorf("Expected validation error for nil hinstance")
	}
	if !errors.As(err, &valErr) || valErr.Field != "hinstance" {
		t.Errorf("Expected ValidationError for hinstance, got %v", err)
	}

	// Test nil hwnd
	_, err = CreateWin32SurfaceKHR(instance, hinstance, nil)
	if err == nil {
		t.Errorf("Expected validation error for nil hwnd")
	}
	if !errors.As(err, &valErr) || valErr.Field != "hwnd" {
		t.Errorf("Expected ValidationError for hwnd, got %v", err)
	}
}
