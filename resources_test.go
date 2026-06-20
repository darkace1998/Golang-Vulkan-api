package vulkan

import (
	"testing"
)

// TestTransitionImageLayoutValidation tests early exit conditions for TransitionImageLayout
func TestTransitionImageLayoutValidation(t *testing.T) {
	// This should not crash if early exits work
	TransitionImageLayout(nil, nil, FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
	TransitionImageLayout(fakeCommandBuffer(), nil, FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
	TransitionImageLayout(nil, fakeImage(), FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})

	// Also test TransitionImageLayout
	TransitionImageLayout(nil, nil, FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
	TransitionImageLayout(fakeCommandBuffer(), nil, FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
	TransitionImageLayout(nil, fakeImage(), FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
}

func TestGetImageSubresourceLayoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		image       Image
		subresource *ImageSubresource
		expectParam string
	}{
		{testNilDevice, nil, fakeImage(), &ImageSubresource{}, testDeviceParameter},
		{"nil image", fakeDevice(), nil, &ImageSubresource{}, "image"},
		{"nil subresource", fakeDevice(), fakeImage(), nil, "subresource"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetImageSubresourceLayout(tt.device, tt.image, tt.subresource)
			// Verify it handles nil safely without crashing.
		})
	}
}
