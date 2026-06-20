package vulkan

/*
#include <vulkan/vulkan.h>
#include <stdlib.h>
#include <string.h>

// Function pointer for UpdateVideoSessionParameters
static PFN_vkUpdateVideoSessionParametersKHR pfn_vkUpdateVideoSessionParametersKHR = NULL;
static PFN_vkGetPhysicalDeviceVideoFormatPropertiesKHR pfn_vkGetPhysicalDeviceVideoFormatPropertiesKHR = NULL;

// Load additional video device functions
static int loadAdditionalVideoDeviceFunctions(VkDevice device) {
    if (device == VK_NULL_HANDLE) {
        return 0;
    }
    pfn_vkUpdateVideoSessionParametersKHR = (PFN_vkUpdateVideoSessionParametersKHR)
        vkGetDeviceProcAddr(device, "vkUpdateVideoSessionParametersKHR");
    return pfn_vkUpdateVideoSessionParametersKHR != NULL;
}

// Load video format query functions
static int loadVideoFormatFunctions(VkInstance instance) {
    if (instance == VK_NULL_HANDLE) {
        return 0;
    }
    pfn_vkGetPhysicalDeviceVideoFormatPropertiesKHR = (PFN_vkGetPhysicalDeviceVideoFormatPropertiesKHR)
        vkGetInstanceProcAddr(instance, "vkGetPhysicalDeviceVideoFormatPropertiesKHR");
    return pfn_vkGetPhysicalDeviceVideoFormatPropertiesKHR != NULL;
}

// Wrapper for UpdateVideoSessionParameters
static VkResult call_vkUpdateVideoSessionParametersKHR(
    VkDevice device,
    VkVideoSessionParametersKHR videoSessionParameters,
    const VkVideoSessionParametersUpdateInfoKHR* pUpdateInfo) {
    if (pfn_vkUpdateVideoSessionParametersKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkUpdateVideoSessionParametersKHR(device, videoSessionParameters, pUpdateInfo);
}

// Wrapper for GetPhysicalDeviceVideoFormatProperties
static VkResult call_vkGetPhysicalDeviceVideoFormatPropertiesKHR(
    VkPhysicalDevice physicalDevice,
    const VkPhysicalDeviceVideoFormatInfoKHR* pVideoFormatInfo,
    uint32_t* pVideoFormatPropertyCount,
    VkVideoFormatPropertiesKHR* pVideoFormatProperties) {
    if (pfn_vkGetPhysicalDeviceVideoFormatPropertiesKHR == NULL) {
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }
    return pfn_vkGetPhysicalDeviceVideoFormatPropertiesKHR(physicalDevice, pVideoFormatInfo, pVideoFormatPropertyCount, pVideoFormatProperties);
}
*/
import "C"

import (
	"sync"
	"unsafe"
)

// VideoDeviceFunctions defines the VideoDeviceFunctions type
// VideoDeviceFunctions holds per-device video function pointers
// This provides thread-safe access to video extension functions
type VideoDeviceFunctions struct {
	device    Device
	loaded    bool
	loadMutex sync.RWMutex
}

// videoDeviceFunctionsMap provides per-device function pointer storage
var (
	videoDeviceFunctionsMap     = make(map[Device]*VideoDeviceFunctions)
	videoDeviceFunctionsMapLock sync.RWMutex
)

// GetVideoDeviceFunctions returns the video functions for a device
func GetVideoDeviceFunctions(device Device) *VideoDeviceFunctions {
	videoDeviceFunctionsMapLock.RLock()
	funcs, exists := videoDeviceFunctionsMap[device]
	videoDeviceFunctionsMapLock.RUnlock()

	if !exists {
		return nil
	}
	return funcs
}

// CreateVideoDeviceFunctions performs the operation
// CreateVideoDeviceFunctions creates and loads video functions for a device
// This provides per-device function pointer storage with thread-safe loading
func CreateVideoDeviceFunctions(device Device) (*VideoDeviceFunctions, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}

	videoDeviceFunctionsMapLock.Lock()
	defer videoDeviceFunctionsMapLock.Unlock()

	// Check if already exists
	if funcs, exists := videoDeviceFunctionsMap[device]; exists {
		return funcs, nil
	}

	// Create new entry
	funcs := &VideoDeviceFunctions{
		device: device,
		loaded: false,
	}

	// Load functions - this also uses the global loading for now
	// but tracks per-device state
	if C.loadAdditionalVideoDeviceFunctions(C.VkDevice(device)) != 0 {
		funcs.loaded = true
	}

	videoDeviceFunctionsMap[device] = funcs
	return funcs, nil
}

// IsLoaded returns whether the video functions are loaded
func (vdf *VideoDeviceFunctions) IsLoaded() bool {
	vdf.loadMutex.RLock()
	defer vdf.loadMutex.RUnlock()
	return vdf.loaded
}

// LoadVideoFormatFunctions performs the operation
// LoadVideoFormatFunctions loads video format query functions
// This must be called after creating a Vulkan instance
func LoadVideoFormatFunctions(instance Instance) bool {
	return C.loadVideoFormatFunctions(C.VkInstance(instance)) != 0
}

// --------------------------------
// Video Encoding Rate Control
// --------------------------------

// VideoEncodeRateControlMode represents video encode rate control modes
type VideoEncodeRateControlMode uint32

const (
	VideoEncodeRateControlModeDefault  VideoEncodeRateControlMode = 0
	VideoEncodeRateControlModeDisabled VideoEncodeRateControlMode = 1
	VideoEncodeRateControlModeCBR      VideoEncodeRateControlMode = 2
	VideoEncodeRateControlModeVBR      VideoEncodeRateControlMode = 3
)

// VideoEncodeRateControlInfo contains rate control configuration
type VideoEncodeRateControlInfo struct {
	Mode                 VideoEncodeRateControlMode
	LayerCount           uint32
	AverageBitrate       uint64
	MaxBitrate           uint64
	FrameRateNumerator   uint32
	FrameRateDenominator uint32
	VirtualBufferSize    uint64
	InitialBufferFill    uint64
}

// --------------------------------
// Shared Session Creation Helpers
// --------------------------------

type videoSessionHelperConfig struct {
	Width               uint32
	Height              uint32
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

func createVideoSessionHelper(device Device, config *videoSessionHelperConfig, codecOp VideoCodecOperationFlags) (VideoSession, error) {
	if device == nil {
		return VideoSession(NullHandle), NewValidationError("device", "cannot be nil")
	}

	videoProfile := &VideoProfileInfo{
		VideoCodecOperation: codecOp,
		ChromaSubsampling:   config.ChromaSubsampling,
		LumaBitDepth:        config.LumaBitDepth,
		ChromaBitDepth:      config.ChromaBitDepth,
	}

	sessionCreateInfo := &VideoSessionCreateInfo{
		QueueFamilyIndex:       config.QueueFamilyIndex,
		VideoProfile:           videoProfile,
		PictureFormat:          config.PictureFormat,
		MaxCodedExtent:         Extent2D{Width: config.Width, Height: config.Height},
		ReferencePictureFormat: config.ReferenceFormat,
		MaxDpbSlots:            config.MaxDpbSlots,
		MaxActiveReferences:    config.MaxActiveReferences,
	}

	return CreateVideoSession(device, sessionCreateInfo)
}

// --------------------------------
// H.264 Encoding Session Helpers
// --------------------------------

// H264Profile represents H.264/AVC profile identifiers
type H264Profile uint32

const (
	H264ProfileBaseline H264Profile = 66
	H264ProfileMain     H264Profile = 77
	H264ProfileHigh     H264Profile = 100
	H264ProfileHigh10   H264Profile = 110
	H264ProfileHigh422  H264Profile = 122
	H264ProfileHigh444  H264Profile = 244
)

// H264Level represents H.264/AVC levels
type H264Level uint32

const (
	H264Level1_0 H264Level = 10
	H264Level1_1 H264Level = 11
	H264Level1_2 H264Level = 12
	H264Level1_3 H264Level = 13
	H264Level2_0 H264Level = 20
	H264Level2_1 H264Level = 21
	H264Level2_2 H264Level = 22
	H264Level3_0 H264Level = 30
	H264Level3_1 H264Level = 31
	H264Level3_2 H264Level = 32
	H264Level4_0 H264Level = 40
	H264Level4_1 H264Level = 41
	H264Level4_2 H264Level = 42
	H264Level5_0 H264Level = 50
	H264Level5_1 H264Level = 51
	H264Level5_2 H264Level = 52
)

// H264EncodeSessionCreateInfo contains configuration for H.264 encode session
type H264EncodeSessionCreateInfo struct {
	Width               uint32
	Height              uint32
	Profile             H264Profile
	Level               H264Level
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	RateControl         *VideoEncodeRateControlInfo
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

// DefaultH264EncodeSessionCreateInfo returns a default H.264 encode session configuration
func DefaultH264EncodeSessionCreateInfo(width, height uint32) *H264EncodeSessionCreateInfo {
	return &H264EncodeSessionCreateInfo{
		Width:               width,
		Height:              height,
		Profile:             H264ProfileHigh,
		Level:               H264Level4_1,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
		MaxDpbSlots:         5,
		MaxActiveReferences: 2,
		PictureFormat:       FormatG8B8R82Plane420Unorm,
		ReferenceFormat:     FormatG8B8R82Plane420Unorm,
	}
}

// CreateH264EncodeSession creates an H.264 encode session with the given configuration
func CreateH264EncodeSession(device Device, createInfo *H264EncodeSessionCreateInfo) (VideoSession, error) {
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	config := &videoSessionHelperConfig{
		Width:               createInfo.Width,
		Height:              createInfo.Height,
		ChromaSubsampling:   createInfo.ChromaSubsampling,
		LumaBitDepth:        createInfo.LumaBitDepth,
		ChromaBitDepth:      createInfo.ChromaBitDepth,
		MaxDpbSlots:         createInfo.MaxDpbSlots,
		MaxActiveReferences: createInfo.MaxActiveReferences,
		QueueFamilyIndex:    createInfo.QueueFamilyIndex,
		PictureFormat:       createInfo.PictureFormat,
		ReferenceFormat:     createInfo.ReferenceFormat,
	}

	return createVideoSessionHelper(device, config, VideoCodecOperationEncodeH264Bit)
}

// --------------------------------
// H.265 Encoding Session Helpers
// --------------------------------

// H265Profile represents H.265/HEVC profile identifiers
type H265Profile uint32

const (
	H265ProfileMain             H265Profile = 1
	H265ProfileMain10           H265Profile = 2
	H265ProfileMainStillPicture H265Profile = 3
	H265ProfileRext             H265Profile = 4
	H265ProfileSCC              H265Profile = 9
)

// H265Level represents H.265/HEVC levels
type H265Level uint32

const (
	H265Level1_0 H265Level = 30
	H265Level2_0 H265Level = 60
	H265Level2_1 H265Level = 63
	H265Level3_0 H265Level = 90
	H265Level3_1 H265Level = 93
	H265Level4_0 H265Level = 120
	H265Level4_1 H265Level = 123
	H265Level5_0 H265Level = 150
	H265Level5_1 H265Level = 153
	H265Level5_2 H265Level = 156
	H265Level6_0 H265Level = 180
	H265Level6_1 H265Level = 183
	H265Level6_2 H265Level = 186
)

// H265EncodeSessionCreateInfo contains configuration for H.265 encode session
type H265EncodeSessionCreateInfo struct {
	Width               uint32
	Height              uint32
	Profile             H265Profile
	Level               H265Level
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	RateControl         *VideoEncodeRateControlInfo
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

// DefaultH265EncodeSessionCreateInfo returns a default H.265 encode session configuration
func DefaultH265EncodeSessionCreateInfo(width, height uint32) *H265EncodeSessionCreateInfo {
	return &H265EncodeSessionCreateInfo{
		Width:               width,
		Height:              height,
		Profile:             H265ProfileMain,
		Level:               H265Level5_1,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
		MaxDpbSlots:         5,
		MaxActiveReferences: 2,
		PictureFormat:       FormatG8B8R82Plane420Unorm,
		ReferenceFormat:     FormatG8B8R82Plane420Unorm,
	}
}

// CreateH265EncodeSession creates an H.265 encode session with the given configuration
func CreateH265EncodeSession(device Device, createInfo *H265EncodeSessionCreateInfo) (VideoSession, error) {
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	config := &videoSessionHelperConfig{
		Width:               createInfo.Width,
		Height:              createInfo.Height,
		ChromaSubsampling:   createInfo.ChromaSubsampling,
		LumaBitDepth:        createInfo.LumaBitDepth,
		ChromaBitDepth:      createInfo.ChromaBitDepth,
		MaxDpbSlots:         createInfo.MaxDpbSlots,
		MaxActiveReferences: createInfo.MaxActiveReferences,
		QueueFamilyIndex:    createInfo.QueueFamilyIndex,
		PictureFormat:       createInfo.PictureFormat,
		ReferenceFormat:     createInfo.ReferenceFormat,
	}

	return createVideoSessionHelper(device, config, VideoCodecOperationEncodeH265Bit)
}

// --------------------------------
// AV1 Encoding Session Helpers
// --------------------------------

// AV1Profile represents AV1 profile identifiers
type AV1Profile uint32

const (
	AV1ProfileMain         AV1Profile = 0
	AV1ProfileHigh         AV1Profile = 1
	AV1ProfileProfessional AV1Profile = 2
)

// AV1Level represents AV1 levels
type AV1Level uint32

const (
	AV1Level2_0 AV1Level = 0
	AV1Level2_1 AV1Level = 1
	AV1Level3_0 AV1Level = 4
	AV1Level3_1 AV1Level = 5
	AV1Level4_0 AV1Level = 8
	AV1Level4_1 AV1Level = 9
	AV1Level5_0 AV1Level = 12
	AV1Level5_1 AV1Level = 13
	AV1Level5_2 AV1Level = 14
	AV1Level5_3 AV1Level = 15
	AV1Level6_0 AV1Level = 16
	AV1Level6_1 AV1Level = 17
	AV1Level6_2 AV1Level = 18
	AV1Level6_3 AV1Level = 19
)

// AV1EncodeSessionCreateInfo contains configuration for AV1 encode session
type AV1EncodeSessionCreateInfo struct {
	Width               uint32
	Height              uint32
	Profile             AV1Profile
	Level               AV1Level
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	RateControl         *VideoEncodeRateControlInfo
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

// DefaultAV1EncodeSessionCreateInfo returns a default AV1 encode session configuration
func DefaultAV1EncodeSessionCreateInfo(width, height uint32) *AV1EncodeSessionCreateInfo {
	return &AV1EncodeSessionCreateInfo{
		Width:               width,
		Height:              height,
		Profile:             AV1ProfileMain,
		Level:               AV1Level5_0,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
		MaxDpbSlots:         8,
		MaxActiveReferences: 3,
		PictureFormat:       FormatG8B8R82Plane420Unorm,
		ReferenceFormat:     FormatG8B8R82Plane420Unorm,
	}
}

// CreateAV1EncodeSession creates an AV1 encode session with the given configuration
func CreateAV1EncodeSession(device Device, createInfo *AV1EncodeSessionCreateInfo) (VideoSession, error) {
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	config := &videoSessionHelperConfig{
		Width:               createInfo.Width,
		Height:              createInfo.Height,
		ChromaSubsampling:   createInfo.ChromaSubsampling,
		LumaBitDepth:        createInfo.LumaBitDepth,
		ChromaBitDepth:      createInfo.ChromaBitDepth,
		MaxDpbSlots:         createInfo.MaxDpbSlots,
		MaxActiveReferences: createInfo.MaxActiveReferences,
		QueueFamilyIndex:    createInfo.QueueFamilyIndex,
		PictureFormat:       createInfo.PictureFormat,
		ReferenceFormat:     createInfo.ReferenceFormat,
	}

	return createVideoSessionHelper(device, config, VideoCodecOperationEncodeAV1Bit)
}

// --------------------------------
// Decoding Session Helpers
// --------------------------------

// H264DecodeSessionCreateInfo contains configuration for H.264 decode session
type H264DecodeSessionCreateInfo struct {
	Width               uint32
	Height              uint32
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

// DefaultH264DecodeSessionCreateInfo returns a default H.264 decode session configuration
func DefaultH264DecodeSessionCreateInfo(width, height uint32) *H264DecodeSessionCreateInfo {
	return &H264DecodeSessionCreateInfo{
		Width:               width,
		Height:              height,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
		MaxDpbSlots:         17,
		MaxActiveReferences: 16,
		PictureFormat:       FormatG8B8R82Plane420Unorm,
		ReferenceFormat:     FormatG8B8R82Plane420Unorm,
	}
}

// CreateH264DecodeSession creates an H.264 decode session with the given configuration
func CreateH264DecodeSession(device Device, createInfo *H264DecodeSessionCreateInfo) (VideoSession, error) {
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	config := &videoSessionHelperConfig{
		Width:               createInfo.Width,
		Height:              createInfo.Height,
		ChromaSubsampling:   createInfo.ChromaSubsampling,
		LumaBitDepth:        createInfo.LumaBitDepth,
		ChromaBitDepth:      createInfo.ChromaBitDepth,
		MaxDpbSlots:         createInfo.MaxDpbSlots,
		MaxActiveReferences: createInfo.MaxActiveReferences,
		QueueFamilyIndex:    createInfo.QueueFamilyIndex,
		PictureFormat:       createInfo.PictureFormat,
		ReferenceFormat:     createInfo.ReferenceFormat,
	}

	return createVideoSessionHelper(device, config, VideoCodecOperationDecodeH264Bit)
}

// H265DecodeSessionCreateInfo contains configuration for H.265 decode session
type H265DecodeSessionCreateInfo struct {
	Width               uint32
	Height              uint32
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

// DefaultH265DecodeSessionCreateInfo returns a default H.265 decode session configuration
func DefaultH265DecodeSessionCreateInfo(width, height uint32) *H265DecodeSessionCreateInfo {
	return &H265DecodeSessionCreateInfo{
		Width:               width,
		Height:              height,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
		MaxDpbSlots:         17,
		MaxActiveReferences: 16,
		PictureFormat:       FormatG8B8R82Plane420Unorm,
		ReferenceFormat:     FormatG8B8R82Plane420Unorm,
	}
}

// CreateH265DecodeSession creates an H.265 decode session with the given configuration
func CreateH265DecodeSession(device Device, createInfo *H265DecodeSessionCreateInfo) (VideoSession, error) {
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	config := &videoSessionHelperConfig{
		Width:               createInfo.Width,
		Height:              createInfo.Height,
		ChromaSubsampling:   createInfo.ChromaSubsampling,
		LumaBitDepth:        createInfo.LumaBitDepth,
		ChromaBitDepth:      createInfo.ChromaBitDepth,
		MaxDpbSlots:         createInfo.MaxDpbSlots,
		MaxActiveReferences: createInfo.MaxActiveReferences,
		QueueFamilyIndex:    createInfo.QueueFamilyIndex,
		PictureFormat:       createInfo.PictureFormat,
		ReferenceFormat:     createInfo.ReferenceFormat,
	}

	return createVideoSessionHelper(device, config, VideoCodecOperationDecodeH265Bit)
}

// AV1DecodeSessionCreateInfo contains configuration for AV1 decode session
type AV1DecodeSessionCreateInfo struct {
	Width               uint32
	Height              uint32
	ChromaSubsampling   VideoChromaSubsampling
	LumaBitDepth        VideoComponentBitDepth
	ChromaBitDepth      VideoComponentBitDepth
	MaxDpbSlots         uint32
	MaxActiveReferences uint32
	QueueFamilyIndex    uint32
	PictureFormat       Format
	ReferenceFormat     Format
}

// DefaultAV1DecodeSessionCreateInfo returns a default AV1 decode session configuration
func DefaultAV1DecodeSessionCreateInfo(width, height uint32) *AV1DecodeSessionCreateInfo {
	return &AV1DecodeSessionCreateInfo{
		Width:               width,
		Height:              height,
		ChromaSubsampling:   VideoChromaSubsampling420,
		LumaBitDepth:        VideoComponentBitDepth8,
		ChromaBitDepth:      VideoComponentBitDepth8,
		MaxDpbSlots:         8,
		MaxActiveReferences: 7,
		PictureFormat:       FormatG8B8R82Plane420Unorm,
		ReferenceFormat:     FormatG8B8R82Plane420Unorm,
	}
}

// CreateAV1DecodeSession creates an AV1 decode session with the given configuration
func CreateAV1DecodeSession(device Device, createInfo *AV1DecodeSessionCreateInfo) (VideoSession, error) {
	if createInfo == nil {
		return VideoSession(NullHandle), NewValidationError("createInfo", "cannot be nil")
	}

	config := &videoSessionHelperConfig{
		Width:               createInfo.Width,
		Height:              createInfo.Height,
		ChromaSubsampling:   createInfo.ChromaSubsampling,
		LumaBitDepth:        createInfo.LumaBitDepth,
		ChromaBitDepth:      createInfo.ChromaBitDepth,
		MaxDpbSlots:         createInfo.MaxDpbSlots,
		MaxActiveReferences: createInfo.MaxActiveReferences,
		QueueFamilyIndex:    createInfo.QueueFamilyIndex,
		PictureFormat:       createInfo.PictureFormat,
		ReferenceFormat:     createInfo.ReferenceFormat,
	}

	return createVideoSessionHelper(device, config, VideoCodecOperationDecodeAV1Bit)
}

// --------------------------------
// DPB (Decoded Picture Buffer) Management
// --------------------------------

// DPBSlot represents a slot in the decoded picture buffer
type DPBSlot struct {
	SlotIndex         int32
	ImageView         ImageView
	ImageLayout       ImageLayout
	IsReference       bool
	PictureOrderCount int32
	FrameNum          int32
	IsLongTerm        bool
}

// DPBManager manages the decoded picture buffer for video decode/encode
type DPBManager struct {
	slots      []DPBSlot
	maxSlots   uint32
	currentPOC int32
	frameNum   int32
}

// CreateDPBManager creates a new DPB manager with the specified number of slots
func CreateDPBManager(maxSlots uint32) *DPBManager {
	return &DPBManager{
		slots:    make([]DPBSlot, 0, maxSlots),
		maxSlots: maxSlots,
	}
}

// AddSlot adds a picture to the DPB
func (dpb *DPBManager) AddSlot(imageView ImageView, imageLayout ImageLayout, poc int32) (*DPBSlot, error) {
	if uint32(len(dpb.slots)) >= dpb.maxSlots {
		// Need to remove oldest reference
		dpb.RemoveOldestReference()
	}

	slot := DPBSlot{
		SlotIndex:         int32(len(dpb.slots)),
		ImageView:         imageView,
		ImageLayout:       imageLayout,
		IsReference:       true,
		PictureOrderCount: poc,
		FrameNum:          dpb.frameNum,
		IsLongTerm:        false,
	}

	dpb.slots = append(dpb.slots, slot)
	dpb.frameNum++

	return &dpb.slots[len(dpb.slots)-1], nil
}

// RemoveOldestReference removes the oldest short-term reference from the DPB
func (dpb *DPBManager) RemoveOldestReference() {
	if len(dpb.slots) == 0 {
		return
	}

	// Find oldest short-term reference
	oldestIdx := -1
	oldestPOC := int32(0x7FFFFFFF)

	for i, slot := range dpb.slots {
		if slot.IsReference && !slot.IsLongTerm && slot.PictureOrderCount < oldestPOC {
			oldestPOC = slot.PictureOrderCount
			oldestIdx = i
		}
	}

	if oldestIdx >= 0 {
		// Remove by setting IsReference to false
		dpb.slots[oldestIdx].IsReference = false
	}
}

// GetReferenceSlots returns all current reference slots
func (dpb *DPBManager) GetReferenceSlots() []DPBSlot {
	refs := make([]DPBSlot, 0, len(dpb.slots))
	for _, slot := range dpb.slots {
		if slot.IsReference {
			refs = append(refs, slot)
		}
	}
	return refs
}

// MarkAsLongTerm marks a slot as a long-term reference
func (dpb *DPBManager) MarkAsLongTerm(slotIndex int32) {
	for i := range dpb.slots {
		if dpb.slots[i].SlotIndex == slotIndex {
			dpb.slots[i].IsLongTerm = true
			break
		}
	}
}

// Reset clears all slots from the DPB
func (dpb *DPBManager) Reset() {
	dpb.slots = dpb.slots[:0]
	dpb.currentPOC = 0
	dpb.frameNum = 0
}

// CalculatePOC performs the operation
// CalculatePOC calculates the Picture Order Count for the next frame
// This is a simplified implementation for H.264/H.265
func (dpb *DPBManager) CalculatePOC() int32 {
	poc := dpb.currentPOC
	dpb.currentPOC += 2 // Increment by 2 for progressive, non-field coding
	return poc
}

// --------------------------------
// Video Queue Family Detection
// --------------------------------

// FindVideoDecodeQueueFamily finds a queue family that supports video decode
func FindVideoDecodeQueueFamily(physicalDevice PhysicalDevice) (uint32, bool) {
	var queueFamilyCount C.uint32_t
	C.vkGetPhysicalDeviceQueueFamilyProperties(C.VkPhysicalDevice(physicalDevice), &queueFamilyCount, nil)

	if queueFamilyCount == 0 {
		return 0, false
	}

	queueFamilies := make([]C.VkQueueFamilyProperties, queueFamilyCount)
	C.vkGetPhysicalDeviceQueueFamilyProperties(C.VkPhysicalDevice(physicalDevice), &queueFamilyCount, &queueFamilies[0])

	for i := uint32(0); i < uint32(queueFamilyCount); i++ {
		// VK_QUEUE_VIDEO_DECODE_BIT_KHR = 0x00000020
		if queueFamilies[i].queueFlags&0x00000020 != 0 {
			return i, true
		}
	}

	return 0, false
}

// FindVideoEncodeQueueFamily finds a queue family that supports video encode
func FindVideoEncodeQueueFamily(physicalDevice PhysicalDevice) (uint32, bool) {
	var queueFamilyCount C.uint32_t
	C.vkGetPhysicalDeviceQueueFamilyProperties(C.VkPhysicalDevice(physicalDevice), &queueFamilyCount, nil)

	if queueFamilyCount == 0 {
		return 0, false
	}

	queueFamilies := make([]C.VkQueueFamilyProperties, queueFamilyCount)
	C.vkGetPhysicalDeviceQueueFamilyProperties(C.VkPhysicalDevice(physicalDevice), &queueFamilyCount, &queueFamilies[0])

	for i := uint32(0); i < uint32(queueFamilyCount); i++ {
		// VK_QUEUE_VIDEO_ENCODE_BIT_KHR = 0x00000040
		if queueFamilies[i].queueFlags&0x00000040 != 0 {
			return i, true
		}
	}

	return 0, false
}

// --------------------------------
// Video Coding Control Helpers
// --------------------------------

// VideoCodingControlFlags represents video coding control flags
type VideoCodingControlFlags uint32

const (
	VideoCodingControlResetBit              VideoCodingControlFlags = 0x00000001
	VideoCodingControlEncodeBit             VideoCodingControlFlags = 0x00000002
	VideoCodingControlEncodeQualityLevelBit VideoCodingControlFlags = 0x00000004
)

// CmdControlVideoCodingReset issues a reset control command for video coding
func CmdControlVideoCodingReset(commandBuffer CommandBuffer) error {
	return CmdControlVideoCoding(commandBuffer, &VideoCodingControlInfo{
		Flags: uint32(VideoCodingControlResetBit),
	})
}

// --------------------------------
// Video Format Properties
// --------------------------------

// VideoFormatProperties contains video format properties information
type VideoFormatProperties struct {
	Format           Format
	ImageCreateFlags uint32
	ImageType        ImageType
	ImageTiling      ImageTiling
	ImageUsageFlags  ImageUsageFlags
}

// GetVideoFormatProperties queries the video format properties for a physical device
func GetVideoFormatProperties(physicalDevice PhysicalDevice, videoProfile *VideoProfileInfo, imageUsage ImageUsageFlags) ([]VideoFormatProperties, error) {
	if physicalDevice == nil {
		return nil, NewValidationError("physicalDevice", "cannot be nil")
	}
	if videoProfile == nil {
		return nil, NewValidationError("videoProfile", "cannot be nil")
	}

	// Create profile info
	var cVideoProfile C.VkVideoProfileInfoKHR
	cVideoProfile.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_INFO_KHR
	cVideoProfile.pNext = nil
	cVideoProfile.videoCodecOperation = C.VkVideoCodecOperationFlagBitsKHR(videoProfile.VideoCodecOperation)
	cVideoProfile.chromaSubsampling = C.VkVideoChromaSubsamplingFlagsKHR(videoProfile.ChromaSubsampling)
	cVideoProfile.lumaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(videoProfile.LumaBitDepth)
	cVideoProfile.chromaBitDepth = C.VkVideoComponentBitDepthFlagsKHR(videoProfile.ChromaBitDepth)

	// Create profile list
	var cProfileList C.VkVideoProfileListInfoKHR
	cProfileList.sType = C.VK_STRUCTURE_TYPE_VIDEO_PROFILE_LIST_INFO_KHR
	cProfileList.pNext = nil
	cProfileList.profileCount = 1
	cProfileList.pProfiles = &cVideoProfile

	// Create format info
	var cFormatInfo C.VkPhysicalDeviceVideoFormatInfoKHR
	cFormatInfo.sType = C.VK_STRUCTURE_TYPE_PHYSICAL_DEVICE_VIDEO_FORMAT_INFO_KHR
	cFormatInfo.pNext = unsafe.Pointer(&cProfileList)
	cFormatInfo.imageUsage = C.VkImageUsageFlags(imageUsage)

	// Get count first
	var count C.uint32_t
	result := Result(C.call_vkGetPhysicalDeviceVideoFormatPropertiesKHR(
		C.VkPhysicalDevice(physicalDevice),
		&cFormatInfo,
		&count,
		nil,
	))

	if result != Success {
		return nil, NewVulkanError(result, "GetVideoFormatProperties", "failed to get format count")
	}

	if count == 0 {
		return []VideoFormatProperties{}, nil
	}

	// Get properties
	cProps := make([]C.VkVideoFormatPropertiesKHR, count)
	for i := range cProps {
		cProps[i].sType = C.VK_STRUCTURE_TYPE_VIDEO_FORMAT_PROPERTIES_KHR
		cProps[i].pNext = nil
	}

	result = Result(C.call_vkGetPhysicalDeviceVideoFormatPropertiesKHR(
		C.VkPhysicalDevice(physicalDevice),
		&cFormatInfo,
		&count,
		&cProps[0],
	))

	if result != Success {
		return nil, NewVulkanError(result, "GetVideoFormatProperties", "failed to get format properties")
	}

	// Convert to Go types
	props := make([]VideoFormatProperties, count)
	for i := range props {
		props[i] = VideoFormatProperties{
			Format:           Format(cProps[i].format),
			ImageCreateFlags: uint32(cProps[i].imageCreateFlags),
			ImageType:        ImageType(cProps[i].imageType),
			ImageTiling:      ImageTiling(cProps[i].imageTiling),
			ImageUsageFlags:  ImageUsageFlags(cProps[i].imageUsageFlags),
		}
	}

	return props, nil
}

// --------------------------------
// Video Session Parameter Update
// --------------------------------

// VideoSessionParametersUpdateInfo contains update information for video session parameters
type VideoSessionParametersUpdateInfo struct {
	UpdateSequenceCount uint32
}

// UpdateVideoSessionParameters updates video session parameters
func UpdateVideoSessionParameters(device Device, videoSessionParameters VideoSessionParameters, updateInfo *VideoSessionParametersUpdateInfo) error {
	if device == nil {
		return NewValidationError("device", "cannot be nil")
	}
	if videoSessionParameters == VideoSessionParameters(NullHandle) {
		return NewValidationError("videoSessionParameters", "cannot be null")
	}
	if updateInfo == nil {
		return NewValidationError("updateInfo", "cannot be nil")
	}

	var cUpdateInfo C.VkVideoSessionParametersUpdateInfoKHR
	cUpdateInfo.sType = C.VK_STRUCTURE_TYPE_VIDEO_SESSION_PARAMETERS_UPDATE_INFO_KHR
	cUpdateInfo.pNext = nil
	cUpdateInfo.updateSequenceCount = C.uint32_t(updateInfo.UpdateSequenceCount)

	result := Result(C.call_vkUpdateVideoSessionParametersKHR(
		C.VkDevice(device),
		C.VkVideoSessionParametersKHR(videoSessionParameters),
		&cUpdateInfo,
	))

	if result != Success {
		return NewVulkanError(result, "UpdateVideoSessionParameters", "failed to update video session parameters")
	}

	return nil
}

// --------------------------------
// YUV Format Helpers
// --------------------------------

// YUVFormat represents common YUV video formats
type YUVFormat uint32

const (
	YUVFormatNV12 YUVFormat = 0 // 4:2:0, 8-bit, semi-planar
	YUVFormatP010 YUVFormat = 1 // 4:2:0, 10-bit, semi-planar
	YUVFormatP016 YUVFormat = 2 // 4:2:0, 16-bit, semi-planar
	YUVFormatYUY2 YUVFormat = 3 // 4:2:2, 8-bit, packed
	YUVFormatY210 YUVFormat = 4 // 4:2:2, 10-bit, packed
	YUVFormatY410 YUVFormat = 5 // 4:4:4, 10-bit, packed
	YUVFormatAYUV YUVFormat = 6 // 4:4:4, 8-bit, packed
)

// YUVFormatToVulkanFormat converts a YUV format to the corresponding Vulkan format
func YUVFormatToVulkanFormat(yuvFormat YUVFormat) Format {
	switch yuvFormat {
	case YUVFormatNV12:
		return FormatG8B8R82Plane420Unorm
	case YUVFormatP010:
		return FormatG10X6B10X6R10X62Plane420Unorm3Pack16
	case YUVFormatP016:
		return FormatG16B16R162Plane420Unorm
	case YUVFormatYUY2:
		return FormatG8B8G8R8422Unorm
	case YUVFormatY210:
		return FormatG10X6B10X6G10X6R10X6422Unorm4Pack16
	case YUVFormatY410:
		return FormatA2R10G10B10UnormPack32 // Note: Closest available
	case YUVFormatAYUV:
		return FormatR8G8B8A8Unorm // Note: Requires format conversion
	default:
		return FormatG8B8R82Plane420Unorm
	}
}

// GetChromaSubsamplingForYUVFormat returns the chroma subsampling for a YUV format
func GetChromaSubsamplingForYUVFormat(yuvFormat YUVFormat) VideoChromaSubsampling {
	switch yuvFormat {
	case YUVFormatNV12, YUVFormatP010, YUVFormatP016:
		return VideoChromaSubsampling420
	case YUVFormatYUY2, YUVFormatY210:
		return VideoChromaSubsampling422
	case YUVFormatY410, YUVFormatAYUV:
		return VideoChromaSubsampling444
	default:
		return VideoChromaSubsampling420
	}
}

// GetBitDepthForYUVFormat returns the luma bit depth for a YUV format
func GetBitDepthForYUVFormat(yuvFormat YUVFormat) VideoComponentBitDepth {
	switch yuvFormat {
	case YUVFormatNV12, YUVFormatYUY2, YUVFormatAYUV:
		return VideoComponentBitDepth8
	case YUVFormatP010, YUVFormatY210, YUVFormatY410:
		return VideoComponentBitDepth10
	case YUVFormatP016:
		return VideoComponentBitDepth12 // or 16
	default:
		return VideoComponentBitDepth8
	}
}

// --------------------------------
// Video Picture Resource Helpers
// --------------------------------

// CreateVideoPictureResource creates a VideoPictureResource from an image view
func CreateVideoPictureResource(imageView ImageView, imageLayout ImageLayout, codedExtent Extent2D) VideoPictureResource {
	return VideoPictureResource{
		ImageView:      imageView,
		ImageLayout:    imageLayout,
		CodedOffset:    Offset2D{X: 0, Y: 0},
		CodedExtent:    codedExtent,
		BaseArrayLayer: 0,
	}
}

// CreateVideoPictureResourceWithOffset creates a VideoPictureResource with a specific offset
func CreateVideoPictureResourceWithOffset(imageView ImageView, imageLayout ImageLayout, codedOffset Offset2D, codedExtent Extent2D, baseArrayLayer uint32) VideoPictureResource {
	return VideoPictureResource{
		ImageView:      imageView,
		ImageLayout:    imageLayout,
		CodedOffset:    codedOffset,
		CodedExtent:    codedExtent,
		BaseArrayLayer: baseArrayLayer,
	}
}
