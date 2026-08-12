package vulkan

import (
	"errors"
	"testing"
)

// TestH264EncodeSessionProfilePropagation verifies that the profile requested
// in H264EncodeSessionCreateInfo is carried into the video profile chain
// instead of being silently replaced by the default.
func TestH264EncodeSessionProfilePropagation(t *testing.T) {
	createInfo := DefaultH264EncodeSessionCreateInfo(1920, 1080)
	createInfo.Profile = H264ProfileBaseline

	config := createInfo.helperConfig()
	if config.EncodeH264Profile == nil {
		t.Fatal("Expected EncodeH264Profile to be set from createInfo.Profile")
	}
	if config.EncodeH264Profile.StdProfileIdc != H264ProfileBaseline {
		t.Errorf("Expected profile %d (Baseline), got %d", H264ProfileBaseline, config.EncodeH264Profile.StdProfileIdc)
	}

	// Zero profile keeps the documented default behavior (nil -> defaults).
	createInfo.Profile = 0
	if createInfo.helperConfig().EncodeH264Profile != nil {
		t.Error("Expected nil EncodeH264Profile when Profile is unset")
	}
}

// TestH265EncodeSessionProfilePropagation is the H.265 counterpart.
func TestH265EncodeSessionProfilePropagation(t *testing.T) {
	createInfo := DefaultH265EncodeSessionCreateInfo(1920, 1080)
	createInfo.Profile = H265ProfileMain10

	config := createInfo.helperConfig()
	if config.EncodeH265Profile == nil {
		t.Fatal("Expected EncodeH265Profile to be set from createInfo.Profile")
	}
	if config.EncodeH265Profile.StdProfileIdc != H265ProfileMain10 {
		t.Errorf("Expected profile %d (Main10), got %d", H265ProfileMain10, config.EncodeH265Profile.StdProfileIdc)
	}

	createInfo.Profile = 0
	if createInfo.helperConfig().EncodeH265Profile != nil {
		t.Error("Expected nil EncodeH265Profile when Profile is unset")
	}
}

// TestGetVideoSessionMemoryBindRequirementsValidation tests input validation
// of the bind-index-aware memory requirements query.
func TestGetVideoSessionMemoryBindRequirementsValidation(t *testing.T) {
	tests := []struct {
		name         string
		device       Device
		videoSession VideoSession
		errorParam   string
	}{
		{testNilDevice, nil, fakeVideoSession(), testDeviceParameter},
		{"null videoSession", fakeDevice(), VideoSession(NullHandle), "videoSession"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetVideoSessionMemoryBindRequirements(tt.device, tt.videoSession)
			if err == nil {
				t.Fatal("Expected error but got nil")
			}
			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if validationErr.Field != tt.errorParam {
				t.Errorf("Expected error for parameter '%s', got '%s'", tt.errorParam, validationErr.Field)
			}
		})
	}
}

// TestCmdDecodeVideoRejectsReferenceSlots verifies that unimplemented
// reference slots are rejected instead of silently dropped.
func TestCmdDecodeVideoRejectsReferenceSlots(t *testing.T) {
	decodeInfo := &VideoDecodeInfo{}
	decodeInfo.ReferenceSlots = make([]struct {
		SlotIndex   int32
		ImageView   ImageView
		ImageLayout ImageLayout
	}, 1)

	err := CmdDecodeVideo(fakeCommandBuffer(), decodeInfo)
	if err == nil {
		t.Fatal("Expected error for unimplemented reference slots")
	}
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if validationErr.Field != "decodeInfo.ReferenceSlots" {
		t.Errorf("Expected error for 'decodeInfo.ReferenceSlots', got '%s'", validationErr.Field)
	}
}

// TestCmdEncodeVideoRejectsReferenceSlots is the encode counterpart.
func TestCmdEncodeVideoRejectsReferenceSlots(t *testing.T) {
	encodeInfo := &VideoEncodeInfo{}
	encodeInfo.ReferenceSlots = make([]struct {
		SlotIndex   int32
		ImageView   ImageView
		ImageLayout ImageLayout
	}, 1)

	err := CmdEncodeVideo(fakeCommandBuffer(), encodeInfo)
	if err == nil {
		t.Fatal("Expected error for unimplemented reference slots")
	}
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %T: %v", err, err)
	}
	if validationErr.Field != "encodeInfo.ReferenceSlots" {
		t.Errorf("Expected error for 'encodeInfo.ReferenceSlots', got '%s'", validationErr.Field)
	}
}

// TestVideoEncodeCapabilityConstants tests the encode capability and feedback
// flag values against the Vulkan headers.
func TestVideoEncodeCapabilityConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    uint32
		expected uint32
	}{
		{"PrecedingExternallyEncodedBytes", uint32(VideoEncodeCapabilityPrecedingExternallyEncodedBytesBit), 0x00000001},
		{"InsufficientBitstreamBufferRangeDetection", uint32(VideoEncodeCapabilityInsufficientBitstreamBufferRangeDetectionBit), 0x00000002},
		{"FeedbackBitstreamBufferOffset", uint32(VideoEncodeFeedbackBitstreamBufferOffsetBit), 0x00000001},
		{"FeedbackBitstreamBytesWritten", uint32(VideoEncodeFeedbackBitstreamBytesWrittenBit), 0x00000002},
		{"FeedbackBitstreamHasOverrides", uint32(VideoEncodeFeedbackBitstreamHasOverridesBit), 0x00000004},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Expected %s to be 0x%08X, got 0x%08X", tt.name, tt.expected, tt.value)
			}
		})
	}
}

// TestVideoCodingControlRateControlBit verifies the rate control bit value and
// its deprecated alias stay in sync.
func TestVideoCodingControlRateControlBit(t *testing.T) {
	if uint32(VideoCodingControlEncodeRateControlBit) != 0x00000002 {
		t.Errorf("Expected VideoCodingControlEncodeRateControlBit to be 0x2, got 0x%X", uint32(VideoCodingControlEncodeRateControlBit))
	}
	if VideoCodingControlEncodeBit != VideoCodingControlEncodeRateControlBit {
		t.Error("Deprecated alias VideoCodingControlEncodeBit must equal VideoCodingControlEncodeRateControlBit")
	}
}

// TestLoadVideoFormatFunctionsNilInstance verifies nil-instance handling and
// that the loader can be reset and retried.
func TestLoadVideoFormatFunctionsNilInstance(t *testing.T) {
	ResetVideoFormatFunctions()
	if LoadVideoFormatFunctions(nil) {
		t.Error("Expected LoadVideoFormatFunctions(nil) to return false")
	}
	// The result is cached; reset must allow a fresh attempt.
	ResetVideoFormatFunctions()
	if LoadVideoFormatFunctions(nil) {
		t.Error("Expected LoadVideoFormatFunctions(nil) to return false after reset")
	}
	ResetVideoFormatFunctions()
}
