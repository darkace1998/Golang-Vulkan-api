package vulkan

import (
	"strings"
	"testing"
)

func TestGetSwapchainImages_NilValidation(t *testing.T) {
	device := fakeDevice()
	swapchain := fakeSwapchain()

	// Test nil device
	_, err := GetSwapchainImages(nil, swapchain)
	if err == nil {
		t.Error("expected error for nil device, got nil")
	} else if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	} else if !strings.Contains(err.Error(), "device") {
		t.Errorf("expected error message to contain 'device', got %q", err.Error())
	}

	// Test nil swapchain
	_, err = GetSwapchainImages(device, nil)
	if err == nil {
		t.Error("expected error for nil swapchain, got nil")
	} else if _, ok := err.(*ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	} else if !strings.Contains(err.Error(), "swapchain") {
		t.Errorf("expected error message to contain 'swapchain', got %q", err.Error())
	}
}
