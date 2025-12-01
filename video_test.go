package vulkan

import (
	"errors"
	"testing"
)

// TestVideoCodecOperationFlags tests video codec operation flag constants
func TestVideoCodecOperationFlags(t *testing.T) {
	tests := []struct {
		name     string
		flag     VideoCodecOperationFlags
		expected uint32
	}{
		{"None", VideoCodecOperationNone, 0},
		{"DecodeH264", VideoCodecOperationDecodeH264Bit, 0x00000001},
		{"DecodeH265", VideoCodecOperationDecodeH265Bit, 0x00000002},
		{"DecodeAV1", VideoCodecOperationDecodeAV1Bit, 0x00000004},
		{"EncodeH264", VideoCodecOperationEncodeH264Bit, 0x00010000},
		{"EncodeH265", VideoCodecOperationEncodeH265Bit, 0x00020000},
		{"EncodeAV1", VideoCodecOperationEncodeAV1Bit, 0x00040000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if uint32(tt.flag) != tt.expected {
				t.Errorf("Expected %s to be 0x%08X, got 0x%08X", tt.name, tt.expected, uint32(tt.flag))
			}
		})
	}
}

// TestVideoChromaSubsampling tests video chroma subsampling constants
func TestVideoChromaSubsampling(t *testing.T) {
	tests := []struct {
		name     string
		value    VideoChromaSubsampling
		expected uint32
	}{
		{"Invalid", VideoChromaSubsamplingInvalid, 0},
		{"Monochrome", VideoChromaSubsamplingMonochrome, 0x00000001},
		{"420", VideoChromaSubsampling420, 0x00000002},
		{"422", VideoChromaSubsampling422, 0x00000004},
		{"444", VideoChromaSubsampling444, 0x00000008},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if uint32(tt.value) != tt.expected {
				t.Errorf("Expected %s to be 0x%08X, got 0x%08X", tt.name, tt.expected, uint32(tt.value))
			}
		})
	}
}

// TestVideoComponentBitDepth tests video component bit depth constants
func TestVideoComponentBitDepth(t *testing.T) {
	tests := []struct {
		name     string
		value    VideoComponentBitDepth
		expected uint32
	}{
		{"Invalid", VideoComponentBitDepthInvalid, 0},
		{"8bit", VideoComponentBitDepth8, 0x00000001},
		{"10bit", VideoComponentBitDepth10, 0x00000004},
		{"12bit", VideoComponentBitDepth12, 0x00000010},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if uint32(tt.value) != tt.expected {
				t.Errorf("Expected %s to be 0x%08X, got 0x%08X", tt.name, tt.expected, uint32(tt.value))
			}
		})
	}
}

// TestVideoExtensionConstants tests video extension name constants
func TestVideoExtensionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"VideoDecodeH264", ExtensionNameVideoDecodeH264, "VK_KHR_video_decode_h264"},
		{"VideoEncodeH264", ExtensionNameVideoEncodeH264, "VK_KHR_video_encode_h264"},
		{"VideoDecodeH265", ExtensionNameVideoDecodeH265, "VK_KHR_video_decode_h265"},
		{"VideoEncodeH265", ExtensionNameVideoEncodeH265, "VK_KHR_video_encode_h265"},
		{"VideoDecodeAV1", ExtensionNameVideoDecodeAV1, "VK_KHR_video_decode_av1"},
		{"VideoEncodeAV1", ExtensionNameVideoEncodeAV1, "VK_KHR_video_encode_av1"},
		{"VideoQueue", ExtensionNameVideoQueue, "VK_KHR_video_queue"},
		{"VideoDecodeQueue", ExtensionNameVideoDecodeQueue, "VK_KHR_video_decode_queue"},
		{"VideoEncodeQueue", ExtensionNameVideoEncodeQueue, "VK_KHR_video_encode_queue"},
		{"VideoMaintenance1", ExtensionNameVideoMaintenance1, "VK_KHR_video_maintenance1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be '%s', got '%s'", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// TestGetVideoCapabilitiesValidation tests input validation for GetVideoCapabilities
func TestGetVideoCapabilitiesValidation(t *testing.T) {
	tests := []struct {
		name           string
		physicalDevice PhysicalDevice
		videoProfile   *VideoProfileInfo
		expectError    bool
		errorParam     string
	}{
		{
			name:           "nil physicalDevice",
			physicalDevice: nil,
			videoProfile: &VideoProfileInfo{
				VideoCodecOperation: VideoCodecOperationDecodeH264Bit,
				ChromaSubsampling:   VideoChromaSubsampling420,
				LumaBitDepth:        VideoComponentBitDepth8,
				ChromaBitDepth:      VideoComponentBitDepth8,
			},
			expectError: true,
			errorParam:  "physicalDevice",
		},
		{
			name:           "nil videoProfile",
			physicalDevice: PhysicalDevice(uintptr(0x1234)), // Non-nil fake handle
			videoProfile:   nil,
			expectError:    true,
			errorParam:     "videoProfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetVideoCapabilities(tt.physicalDevice, tt.videoProfile)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
					return
				}

				if validationErr.Parameter != tt.errorParam {
					t.Errorf("Expected error for parameter '%s', got '%s'", tt.errorParam, validationErr.Parameter)
				}
			}
		})
	}
}

// TestCreateVideoSessionValidation tests input validation for CreateVideoSession
func TestCreateVideoSessionValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *VideoSessionCreateInfo
		expectError bool
		errorParam  string
	}{
		{
			name:   "nil device",
			device: nil,
			createInfo: &VideoSessionCreateInfo{
				VideoProfile: &VideoProfileInfo{
					VideoCodecOperation: VideoCodecOperationDecodeH264Bit,
				},
			},
			expectError: true,
			errorParam:  "device",
		},
		{
			name:        "nil createInfo",
			device:      Device(uintptr(0x1234)), // Non-nil fake handle
			createInfo:  nil,
			expectError: true,
			errorParam:  "createInfo",
		},
		{
			name:   "nil videoProfile in createInfo",
			device: Device(uintptr(0x1234)),
			createInfo: &VideoSessionCreateInfo{
				VideoProfile: nil,
			},
			expectError: true,
			errorParam:  "createInfo.VideoProfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateVideoSession(tt.device, tt.createInfo)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
					return
				}

				if validationErr.Parameter != tt.errorParam {
					t.Errorf("Expected error for parameter '%s', got '%s'", tt.errorParam, validationErr.Parameter)
				}
			}
		})
	}
}

// TestGetVideoSessionMemoryRequirementsValidation tests input validation
func TestGetVideoSessionMemoryRequirementsValidation(t *testing.T) {
	tests := []struct {
		name         string
		device       Device
		videoSession VideoSession
		expectError  bool
		errorParam   string
	}{
		{
			name:         "nil device",
			device:       nil,
			videoSession: VideoSession(uintptr(0x1234)),
			expectError:  true,
			errorParam:   "device",
		},
		{
			name:         "null videoSession",
			device:       Device(uintptr(0x1234)),
			videoSession: VideoSession(NullHandle),
			expectError:  true,
			errorParam:   "videoSession",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetVideoSessionMemoryRequirements(tt.device, tt.videoSession)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
					return
				}

				if validationErr.Parameter != tt.errorParam {
					t.Errorf("Expected error for parameter '%s', got '%s'", tt.errorParam, validationErr.Parameter)
				}
			}
		})
	}
}

// TestBindVideoSessionMemoryValidation tests input validation
func TestBindVideoSessionMemoryValidation(t *testing.T) {
	tests := []struct {
		name         string
		device       Device
		videoSession VideoSession
		bindInfos    []VideoBindMemoryInfo
		expectError  bool
		errorParam   string
	}{
		{
			name:         "nil device",
			device:       nil,
			videoSession: VideoSession(uintptr(0x1234)),
			bindInfos:    []VideoBindMemoryInfo{{MemoryBindIndex: 0}},
			expectError:  true,
			errorParam:   "device",
		},
		{
			name:         "null videoSession",
			device:       Device(uintptr(0x1234)),
			videoSession: VideoSession(NullHandle),
			bindInfos:    []VideoBindMemoryInfo{{MemoryBindIndex: 0}},
			expectError:  true,
			errorParam:   "videoSession",
		},
		{
			name:         "empty bindInfos",
			device:       Device(uintptr(0x1234)),
			videoSession: VideoSession(uintptr(0x5678)),
			bindInfos:    []VideoBindMemoryInfo{},
			expectError:  true,
			errorParam:   "bindInfos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindVideoSessionMemory(tt.device, tt.videoSession, tt.bindInfos)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
					return
				}

				if validationErr.Parameter != tt.errorParam {
					t.Errorf("Expected error for parameter '%s', got '%s'", tt.errorParam, validationErr.Parameter)
				}
			}
		})
	}
}

// TestCreateVideoSessionParametersValidation tests input validation
func TestCreateVideoSessionParametersValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *VideoSessionParametersCreateInfo
		expectError bool
		errorParam  string
	}{
		{
			name:        "nil device",
			device:      nil,
			createInfo:  &VideoSessionParametersCreateInfo{},
			expectError: true,
			errorParam:  "device",
		},
		{
			name:        "nil createInfo",
			device:      Device(uintptr(0x1234)),
			createInfo:  nil,
			expectError: true,
			errorParam:  "createInfo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateVideoSessionParameters(tt.device, tt.createInfo)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
					return
				}

				if validationErr.Parameter != tt.errorParam {
					t.Errorf("Expected error for parameter '%s', got '%s'", tt.errorParam, validationErr.Parameter)
				}
			}
		})
	}
}

// TestVideoProfileInfo tests VideoProfileInfo structure creation
func TestVideoProfileInfo(t *testing.T) {
	profile := &VideoProfileInfo{
		VideoCodecOperation: VideoCodecOperationDecodeH264Bit,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
	}

	if profile.VideoCodecOperation != VideoCodecOperationDecodeH264Bit {
		t.Errorf("Expected VideoCodecOperation to be DecodeH264Bit")
	}
	if profile.ChromaSubsampling != VideoChromaSubsampling420 {
		t.Errorf("Expected ChromaSubsampling to be 420")
	}
	if profile.LumaBitDepth != VideoComponentBitDepth8 {
		t.Errorf("Expected LumaBitDepth to be 8")
	}
	if profile.ChromaBitDepth != VideoComponentBitDepth8 {
		t.Errorf("Expected ChromaBitDepth to be 8")
	}
}

// TestVideoCapabilities tests VideoCapabilities structure
func TestVideoCapabilities(t *testing.T) {
	caps := &VideoCapabilities{
		Flags:                         0x01,
		MinBitstreamBufferOffsetAlign: 256,
		MinBitstreamBufferSizeAlign:   4096,
		PictureAccessGranularity:      Extent2D{Width: 16, Height: 16},
		MinCodedExtent:                Extent2D{Width: 64, Height: 64},
		MaxCodedExtent:                Extent2D{Width: 4096, Height: 4096},
		MaxDpbSlots:                   16,
		MaxActiveReferencePictures:    8,
	}

	if caps.Flags != 0x01 {
		t.Errorf("Expected Flags to be 0x01, got 0x%X", caps.Flags)
	}
	if caps.MinBitstreamBufferOffsetAlign != 256 {
		t.Errorf("Expected MinBitstreamBufferOffsetAlign to be 256")
	}
	if caps.MinBitstreamBufferSizeAlign != 4096 {
		t.Errorf("Expected MinBitstreamBufferSizeAlign to be 4096")
	}
	if caps.MaxDpbSlots != 16 {
		t.Errorf("Expected MaxDpbSlots to be 16")
	}
	if caps.MaxActiveReferencePictures != 8 {
		t.Errorf("Expected MaxActiveReferencePictures to be 8")
	}
}

// TestVideoSessionCreateInfo tests VideoSessionCreateInfo structure
func TestVideoSessionCreateInfo(t *testing.T) {
	profile := &VideoProfileInfo{
		VideoCodecOperation: VideoCodecOperationDecodeH265Bit,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth10,
		ChromaBitDepth:      VideoComponentBitDepth10,
	}

	createInfo := &VideoSessionCreateInfo{
		QueueFamilyIndex:       0,
		VideoProfile:           profile,
		PictureFormat:          FormatR8G8B8A8Unorm,
		MaxCodedExtent:         Extent2D{Width: 1920, Height: 1080},
		ReferencePictureFormat: FormatR8G8B8A8Unorm,
		MaxDpbSlots:            16,
		MaxActiveReferences:    8,
	}

	if createInfo.QueueFamilyIndex != 0 {
		t.Errorf("Expected QueueFamilyIndex to be 0")
	}
	if createInfo.VideoProfile != profile {
		t.Errorf("Expected VideoProfile to match")
	}
	if createInfo.MaxCodedExtent.Width != 1920 || createInfo.MaxCodedExtent.Height != 1080 {
		t.Errorf("Expected MaxCodedExtent to be 1920x1080")
	}
	if createInfo.MaxDpbSlots != 16 {
		t.Errorf("Expected MaxDpbSlots to be 16")
	}
}

// TestVideoPictureResource tests VideoPictureResource structure
func TestVideoPictureResource(t *testing.T) {
	resource := VideoPictureResource{
		ImageView:      ImageView(uintptr(0x1234)),
		ImageLayout:    ImageLayoutGeneral,
		CodedOffset:    Offset2D{X: 0, Y: 0},
		CodedExtent:    Extent2D{Width: 1920, Height: 1080},
		BaseArrayLayer: 0,
	}

	if resource.CodedExtent.Width != 1920 {
		t.Errorf("Expected CodedExtent.Width to be 1920")
	}
	if resource.CodedExtent.Height != 1080 {
		t.Errorf("Expected CodedExtent.Height to be 1080")
	}
	if resource.BaseArrayLayer != 0 {
		t.Errorf("Expected BaseArrayLayer to be 0")
	}
}

// TestVideoDecodeInfo tests VideoDecodeInfo structure
func TestVideoDecodeInfo(t *testing.T) {
	decodeInfo := &VideoDecodeInfo{
		SrcBuffer:       Buffer(uintptr(0x1234)),
		SrcBufferOffset: 0,
		SrcBufferRange:  1024,
		DstPictureResource: VideoPictureResource{
			ImageView:   ImageView(uintptr(0x5678)),
			ImageLayout: ImageLayoutGeneral,
			CodedExtent: Extent2D{Width: 1920, Height: 1080},
		},
	}

	if decodeInfo.SrcBufferRange != 1024 {
		t.Errorf("Expected SrcBufferRange to be 1024")
	}
	if decodeInfo.DstPictureResource.CodedExtent.Width != 1920 {
		t.Errorf("Expected DstPictureResource.CodedExtent.Width to be 1920")
	}
}

// TestVideoEncodeInfo tests VideoEncodeInfo structure
func TestVideoEncodeInfo(t *testing.T) {
	encodeInfo := &VideoEncodeInfo{
		SrcPictureResource: VideoPictureResource{
			ImageView:   ImageView(uintptr(0x1234)),
			ImageLayout: ImageLayoutGeneral,
			CodedExtent: Extent2D{Width: 1920, Height: 1080},
		},
		DstBuffer:       Buffer(uintptr(0x5678)),
		DstBufferOffset: 0,
		DstBufferRange:  4096,
	}

	if encodeInfo.DstBufferRange != 4096 {
		t.Errorf("Expected DstBufferRange to be 4096")
	}
	if encodeInfo.SrcPictureResource.CodedExtent.Width != 1920 {
		t.Errorf("Expected SrcPictureResource.CodedExtent.Width to be 1920")
	}
}

// TestVideoBindMemoryInfo tests VideoBindMemoryInfo structure
func TestVideoBindMemoryInfo(t *testing.T) {
	bindInfo := VideoBindMemoryInfo{
		MemoryBindIndex: 0,
		Memory:          DeviceMemory(uintptr(0x1234)),
		MemoryOffset:    0,
		MemorySize:      1024 * 1024, // 1MB
	}

	if bindInfo.MemoryBindIndex != 0 {
		t.Errorf("Expected MemoryBindIndex to be 0")
	}
	if bindInfo.MemorySize != 1024*1024 {
		t.Errorf("Expected MemorySize to be 1MB")
	}
}

// TestCombinedCodecFlags tests combining video codec operation flags
func TestCombinedCodecFlags(t *testing.T) {
	// Test combining multiple decode flags
	decodeFlags := VideoCodecOperationDecodeH264Bit | VideoCodecOperationDecodeH265Bit
	expected := VideoCodecOperationFlags(0x00000003)
	if decodeFlags != expected {
		t.Errorf("Expected combined decode flags to be 0x%08X, got 0x%08X", expected, decodeFlags)
	}

	// Test combining encode flags
	encodeFlags := VideoCodecOperationEncodeH264Bit | VideoCodecOperationEncodeH265Bit
	expectedEncode := VideoCodecOperationFlags(0x00030000)
	if encodeFlags != expectedEncode {
		t.Errorf("Expected combined encode flags to be 0x%08X, got 0x%08X", expectedEncode, encodeFlags)
	}

	// Test all flags combined
	allFlags := VideoCodecOperationDecodeH264Bit | VideoCodecOperationDecodeH265Bit |
		VideoCodecOperationDecodeAV1Bit | VideoCodecOperationEncodeH264Bit |
		VideoCodecOperationEncodeH265Bit | VideoCodecOperationEncodeAV1Bit
	expectedAll := VideoCodecOperationFlags(0x00070007)
	if allFlags != expectedAll {
		t.Errorf("Expected all flags to be 0x%08X, got 0x%08X", expectedAll, allFlags)
	}
}

// BenchmarkVideoProfileCreation benchmarks video profile creation
func BenchmarkVideoProfileCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = &VideoProfileInfo{
			VideoCodecOperation: VideoCodecOperationDecodeH264Bit,
			ChromaSubsampling:   VideoChromaSubsampling420,
			LumaBitDepth:        VideoComponentBitDepth8,
			ChromaBitDepth:      VideoComponentBitDepth8,
		}
	}
}

// BenchmarkVideoSessionCreateInfoCreation benchmarks VideoSessionCreateInfo creation
func BenchmarkVideoSessionCreateInfoCreation(b *testing.B) {
	profile := &VideoProfileInfo{
		VideoCodecOperation: VideoCodecOperationDecodeH264Bit,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &VideoSessionCreateInfo{
			QueueFamilyIndex:       0,
			VideoProfile:           profile,
			PictureFormat:          FormatR8G8B8A8Unorm,
			MaxCodedExtent:         Extent2D{Width: 1920, Height: 1080},
			ReferencePictureFormat: FormatR8G8B8A8Unorm,
			MaxDpbSlots:            16,
			MaxActiveReferences:    8,
		}
	}
}
