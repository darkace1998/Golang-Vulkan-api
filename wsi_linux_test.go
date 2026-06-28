//go:build linux || freebsd || openbsd || netbsd

package vulkan

import (
	"errors"
	"testing"
	"unsafe"
)

func TestCreateXlibSurfaceKHRValidation(t *testing.T) {
	instance := Instance(fakeHandle())
	var dpy unsafe.Pointer = unsafe.Pointer(uintptr(1))

	// Test nil instance
	_, err := CreateXlibSurfaceKHR(nil, dpy, 0)
	if err == nil {
		t.Errorf("Expected validation error for nil instance")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) || valErr.Field != "instance" {
		t.Errorf("Expected ValidationError for instance, got %v", err)
	}

	// Test nil dpy
	_, err = CreateXlibSurfaceKHR(instance, nil, 0)
	if err == nil {
		t.Errorf("Expected validation error for nil dpy")
	}
	if !errors.As(err, &valErr) || valErr.Field != "dpy" {
		t.Errorf("Expected ValidationError for dpy, got %v", err)
	}
}

func TestCreateXcbSurfaceKHRValidation(t *testing.T) {
	instance := Instance(fakeHandle())
	var connection unsafe.Pointer = unsafe.Pointer(uintptr(1))

	// Test nil instance
	_, err := CreateXcbSurfaceKHR(nil, connection, 0)
	if err == nil {
		t.Errorf("Expected validation error for nil instance")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) || valErr.Field != "instance" {
		t.Errorf("Expected ValidationError for instance, got %v", err)
	}

	// Test nil connection
	_, err = CreateXcbSurfaceKHR(instance, nil, 0)
	if err == nil {
		t.Errorf("Expected validation error for nil connection")
	}
	if !errors.As(err, &valErr) || valErr.Field != "connection" {
		t.Errorf("Expected ValidationError for connection, got %v", err)
	}
}

func TestCreateWaylandSurfaceKHRValidation(t *testing.T) {
	instance := Instance(fakeHandle())
	var display unsafe.Pointer = unsafe.Pointer(uintptr(1))
	var surface unsafe.Pointer = unsafe.Pointer(uintptr(2))

	// Test nil instance
	_, err := CreateWaylandSurfaceKHR(nil, display, surface)
	if err == nil {
		t.Errorf("Expected validation error for nil instance")
	}
	var valErr *ValidationError
	if !errors.As(err, &valErr) || valErr.Field != "instance" {
		t.Errorf("Expected ValidationError for instance, got %v", err)
	}

	// Test nil display
	_, err = CreateWaylandSurfaceKHR(instance, nil, surface)
	if err == nil {
		t.Errorf("Expected validation error for nil display")
	}
	if !errors.As(err, &valErr) || valErr.Field != "display" {
		t.Errorf("Expected ValidationError for display, got %v", err)
	}

	// Test nil surface
	_, err = CreateWaylandSurfaceKHR(instance, display, nil)
	if err == nil {
		t.Errorf("Expected validation error for nil surface")
	}
	if !errors.As(err, &valErr) || valErr.Field != "surface" {
		t.Errorf("Expected ValidationError for surface, got %v", err)
	}
}
