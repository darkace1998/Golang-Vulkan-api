package vulkan

/*
#include <vulkan/vulkan.h>
*/
import "C"
import "unsafe"

// ColorSpace represents color space values
type ColorSpace uint32

const (
	ColorSpaceSRGBNonlinear         ColorSpace = C.VK_COLOR_SPACE_SRGB_NONLINEAR_KHR
	ColorSpaceDisplayP3Nonlinear    ColorSpace = C.VK_COLOR_SPACE_DISPLAY_P3_NONLINEAR_EXT
	ColorSpaceExtendedSRGBLinear    ColorSpace = C.VK_COLOR_SPACE_EXTENDED_SRGB_LINEAR_EXT
	ColorSpaceDisplayP3Linear       ColorSpace = C.VK_COLOR_SPACE_DISPLAY_P3_LINEAR_EXT
	ColorSpaceDCIP3Nonlinear        ColorSpace = C.VK_COLOR_SPACE_DCI_P3_NONLINEAR_EXT
	ColorSpaceBT709Linear           ColorSpace = C.VK_COLOR_SPACE_BT709_LINEAR_EXT
	ColorSpaceBT709Nonlinear        ColorSpace = C.VK_COLOR_SPACE_BT709_NONLINEAR_EXT
	ColorSpaceBT2020Linear          ColorSpace = C.VK_COLOR_SPACE_BT2020_LINEAR_EXT
	ColorSpaceHDR10ST2084           ColorSpace = C.VK_COLOR_SPACE_HDR10_ST2084_EXT
	ColorSpaceDolbyVision           ColorSpace = C.VK_COLOR_SPACE_DOLBYVISION_EXT
	ColorSpaceHDR10HLG              ColorSpace = C.VK_COLOR_SPACE_HDR10_HLG_EXT
	ColorSpaceAdobeRGBLinear        ColorSpace = C.VK_COLOR_SPACE_ADOBERGB_LINEAR_EXT
	ColorSpaceAdobeRGBNonlinear     ColorSpace = C.VK_COLOR_SPACE_ADOBERGB_NONLINEAR_EXT
	ColorSpacePassThrough           ColorSpace = C.VK_COLOR_SPACE_PASS_THROUGH_EXT
	ColorSpaceExtendedSRGBNonlinear ColorSpace = C.VK_COLOR_SPACE_EXTENDED_SRGB_NONLINEAR_EXT
)

// SurfaceTransformFlags represents surface transform flags
type SurfaceTransformFlags uint32

const (
	SurfaceTransformIdentity                  SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_IDENTITY_BIT_KHR
	SurfaceTransformRotate90                  SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_ROTATE_90_BIT_KHR
	SurfaceTransformRotate180                 SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_ROTATE_180_BIT_KHR
	SurfaceTransformRotate270                 SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_ROTATE_270_BIT_KHR
	SurfaceTransformHorizontalMirror          SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_HORIZONTAL_MIRROR_BIT_KHR
	SurfaceTransformHorizontalMirrorRotate90  SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_HORIZONTAL_MIRROR_ROTATE_90_BIT_KHR
	SurfaceTransformHorizontalMirrorRotate180 SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_HORIZONTAL_MIRROR_ROTATE_180_BIT_KHR
	SurfaceTransformHorizontalMirrorRotate270 SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_HORIZONTAL_MIRROR_ROTATE_270_BIT_KHR
	SurfaceTransformInherit                   SurfaceTransformFlags = C.VK_SURFACE_TRANSFORM_INHERIT_BIT_KHR
)

// CompositeAlphaFlags represents composite alpha flags
type CompositeAlphaFlags uint32

const (
	CompositeAlphaOpaque         CompositeAlphaFlags = C.VK_COMPOSITE_ALPHA_OPAQUE_BIT_KHR
	CompositeAlphaPreMultiplied  CompositeAlphaFlags = C.VK_COMPOSITE_ALPHA_PRE_MULTIPLIED_BIT_KHR
	CompositeAlphaPostMultiplied CompositeAlphaFlags = C.VK_COMPOSITE_ALPHA_POST_MULTIPLIED_BIT_KHR
	CompositeAlphaInherit        CompositeAlphaFlags = C.VK_COMPOSITE_ALPHA_INHERIT_BIT_KHR
)

// SwapchainCreateFlags represents swapchain creation flags
type SwapchainCreateFlags uint32

const (
	SwapchainCreateSplitInstanceBindRegions SwapchainCreateFlags = C.VK_SWAPCHAIN_CREATE_SPLIT_INSTANCE_BIND_REGIONS_BIT_KHR
	SwapchainCreateProtected                SwapchainCreateFlags = C.VK_SWAPCHAIN_CREATE_PROTECTED_BIT_KHR
	SwapchainCreateMutableFormat            SwapchainCreateFlags = C.VK_SWAPCHAIN_CREATE_MUTABLE_FORMAT_BIT_KHR
)

// SwapchainCreateInfo contains swapchain creation information
type SwapchainCreateInfo struct {
	Flags              SwapchainCreateFlags
	Surface            Surface
	MinImageCount      uint32
	ImageFormat        Format
	ImageColorSpace    ColorSpace
	ImageExtent        Extent2D
	ImageArrayLayers   uint32
	ImageUsage         ImageUsageFlags
	ImageSharingMode   SharingMode
	QueueFamilyIndices []uint32
	PreTransform       SurfaceTransformFlags
	CompositeAlpha     CompositeAlphaFlags
	PresentMode        PresentMode
	Clipped            bool
	OldSwapchain       Swapchain
}

// CreateSwapchain creates a swapchain
func CreateSwapchain(device Device, createInfo *SwapchainCreateInfo) (Swapchain, error) {
	// Input validation
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if createInfo == nil {
		return nil, NewValidationError("createInfo", "cannot be nil")
	}
	if createInfo.Surface == nil {
		return nil, NewValidationError("Surface", "cannot be nil")
	}
	if createInfo.MinImageCount == 0 {
		return nil, NewValidationError("MinImageCount", "cannot be zero")
	}
	if createInfo.ImageExtent.Width == 0 || createInfo.ImageExtent.Height == 0 {
		return nil, NewValidationError("ImageExtent", "width and height cannot be zero")
	}
	if createInfo.ImageArrayLayers == 0 {
		return nil, NewValidationError("ImageArrayLayers", "cannot be zero")
	}

	var cCreateInfo C.VkSwapchainCreateInfoKHR
	cCreateInfo.sType = C.VK_STRUCTURE_TYPE_SWAPCHAIN_CREATE_INFO_KHR
	cCreateInfo.pNext = nil
	cCreateInfo.flags = C.VkSwapchainCreateFlagsKHR(createInfo.Flags)
	cCreateInfo.surface = C.VkSurfaceKHR(createInfo.Surface)
	cCreateInfo.minImageCount = C.uint32_t(createInfo.MinImageCount)
	cCreateInfo.imageFormat = C.VkFormat(createInfo.ImageFormat)
	cCreateInfo.imageColorSpace = C.VkColorSpaceKHR(createInfo.ImageColorSpace)
	cCreateInfo.imageExtent.width = C.uint32_t(createInfo.ImageExtent.Width)
	cCreateInfo.imageExtent.height = C.uint32_t(createInfo.ImageExtent.Height)
	cCreateInfo.imageArrayLayers = C.uint32_t(createInfo.ImageArrayLayers)
	cCreateInfo.imageUsage = C.VkImageUsageFlags(createInfo.ImageUsage)
	cCreateInfo.imageSharingMode = C.VkSharingMode(createInfo.ImageSharingMode)

	// Queue family indices
	if len(createInfo.QueueFamilyIndices) > 0 {
		cCreateInfo.queueFamilyIndexCount = C.uint32_t(len(createInfo.QueueFamilyIndices))
		cCreateInfo.pQueueFamilyIndices = (*C.uint32_t)(&createInfo.QueueFamilyIndices[0])
	} else {
		cCreateInfo.queueFamilyIndexCount = 0
		cCreateInfo.pQueueFamilyIndices = nil
	}

	cCreateInfo.preTransform = C.VkSurfaceTransformFlagBitsKHR(createInfo.PreTransform)
	cCreateInfo.compositeAlpha = C.VkCompositeAlphaFlagBitsKHR(createInfo.CompositeAlpha)
	cCreateInfo.presentMode = C.VkPresentModeKHR(createInfo.PresentMode)

	if createInfo.Clipped {
		cCreateInfo.clipped = C.VK_TRUE
	} else {
		cCreateInfo.clipped = C.VK_FALSE
	}

	if createInfo.OldSwapchain != nil {
		cCreateInfo.oldSwapchain = C.VkSwapchainKHR(createInfo.OldSwapchain)
	} else {
		cCreateInfo.oldSwapchain = nil
	}

	var swapchain C.VkSwapchainKHR
	result := Result(C.vkCreateSwapchainKHR(C.VkDevice(device), &cCreateInfo, nil, &swapchain))
	if result != Success {
		return nil, NewVulkanError(result, "CreateSwapchain", "Vulkan swapchain creation failed")
	}

	trackResource("Swapchain", unsafe.Pointer(swapchain))
	return Swapchain(swapchain), nil
}

// DestroySwapchain destroys a swapchain
func DestroySwapchain(device Device, swapchain Swapchain) {
	if device == nil || swapchain == nil {
		return
	}
	C.vkDestroySwapchainKHR(C.VkDevice(device), C.VkSwapchainKHR(swapchain), nil)
	untrackResource("Swapchain", unsafe.Pointer(swapchain))
	untrackResource("SwapchainKHR", unsafe.Pointer(swapchain))
}

// GetSwapchainImages gets the swapchain images
func GetSwapchainImages(device Device, swapchain Swapchain) ([]Image, error) {
	if device == nil {
		return nil, NewValidationError("device", "cannot be nil")
	}
	if swapchain == nil {
		return nil, NewValidationError("swapchain", "cannot be nil")
	}

	var imageCount C.uint32_t
	result := Result(C.vkGetSwapchainImagesKHR(C.VkDevice(device), C.VkSwapchainKHR(swapchain), &imageCount, nil))
	if result != Success {
		return nil, NewVulkanError(result, "GetSwapchainImages", "Vulkan swapchain images query failed")
	}

	if imageCount == 0 {
		return []Image{}, nil
	}

	cImages := make([]C.VkImage, imageCount)
	result = Result(C.vkGetSwapchainImagesKHR(C.VkDevice(device), C.VkSwapchainKHR(swapchain), &imageCount, &cImages[0]))
	if result != Success {
		return nil, NewVulkanError(result, "GetSwapchainImages", "Vulkan swapchain images query failed")
	}

	images := make([]Image, imageCount)
	for i := range images {
		images[i] = Image(cImages[i])
	}

	return images, nil
}

// AcquireNextImage acquires the next presentable image from a swapchain. Returns the index of the next image to use, and whether the swapchain is suboptimal.
func AcquireNextImage(device Device, swapchain Swapchain, timeout uint64, semaphore Semaphore, fence Fence) (uint32, bool, error) {
	if device == nil {
		return 0, false, NewValidationError("device", "cannot be nil")
	}
	if swapchain == nil {
		return 0, false, NewValidationError("swapchain", "cannot be nil")
	}

	var imageIndex C.uint32_t
	var cSemaphore C.VkSemaphore
	var cFence C.VkFence

	if semaphore != nil {
		cSemaphore = C.VkSemaphore(semaphore)
	}
	if fence != nil {
		cFence = C.VkFence(fence)
	}

	result := Result(C.vkAcquireNextImageKHR(
		C.VkDevice(device),
		C.VkSwapchainKHR(swapchain),
		C.uint64_t(timeout),
		cSemaphore,
		cFence,
		&imageIndex,
	))

	// VK_SUBOPTIMAL_KHR is a success code that indicates the swapchain is suboptimal
	if result == Result(C.VK_SUBOPTIMAL_KHR) {
		return uint32(imageIndex), true, nil
	}

	if result != Success {
		return 0, false, NewVulkanError(result, "AcquireNextImage", "Vulkan acquire next image failed")
	}

	return uint32(imageIndex), false, nil
}

// PresentInfo contains presentation information
type PresentInfo struct {
	WaitSemaphores []Semaphore
	Swapchains     []Swapchain
	ImageIndices   []uint32
}

// QueuePresent queues an image for presentation. Returns true if the swapchain is suboptimal.
func QueuePresent(queue Queue, presentInfo *PresentInfo) (bool, error) {
	if queue == nil {
		return false, NewValidationError("queue", "cannot be nil")
	}
	if presentInfo == nil {
		return false, NewValidationError("presentInfo", "cannot be nil")
	}
	if len(presentInfo.Swapchains) == 0 {
		return false, NewValidationError("Swapchains", "cannot be empty")
	}
	if len(presentInfo.ImageIndices) != len(presentInfo.Swapchains) {
		return false, NewValidationError("ImageIndices", "must match Swapchains length")
	}

	var cPresentInfo C.VkPresentInfoKHR
	cPresentInfo.sType = C.VK_STRUCTURE_TYPE_PRESENT_INFO_KHR
	cPresentInfo.pNext = nil

	// Wait semaphores
	var cWaitSemaphores []C.VkSemaphore
	if len(presentInfo.WaitSemaphores) > 0 {
		cWaitSemaphores = make([]C.VkSemaphore, len(presentInfo.WaitSemaphores))
		for i, sem := range presentInfo.WaitSemaphores {
			cWaitSemaphores[i] = C.VkSemaphore(sem)
		}
		cPresentInfo.waitSemaphoreCount = C.uint32_t(len(cWaitSemaphores))
		cPresentInfo.pWaitSemaphores = &cWaitSemaphores[0]
	} else {
		cPresentInfo.waitSemaphoreCount = 0
		cPresentInfo.pWaitSemaphores = nil
	}

	// Swapchains
	cSwapchains := make([]C.VkSwapchainKHR, len(presentInfo.Swapchains))
	for i, sc := range presentInfo.Swapchains {
		cSwapchains[i] = C.VkSwapchainKHR(sc)
	}
	cPresentInfo.swapchainCount = C.uint32_t(len(cSwapchains))
	cPresentInfo.pSwapchains = &cSwapchains[0]

	// Image indices
	cImageIndices := make([]C.uint32_t, len(presentInfo.ImageIndices))
	for i, idx := range presentInfo.ImageIndices {
		cImageIndices[i] = C.uint32_t(idx)
	}
	cPresentInfo.pImageIndices = &cImageIndices[0]

	// We don't support pResults for now
	cPresentInfo.pResults = nil

	result := Result(C.vkQueuePresentKHR(C.VkQueue(queue), &cPresentInfo))

	// VK_SUBOPTIMAL_KHR is a success code that indicates the swapchain is suboptimal
	if result == Result(C.VK_SUBOPTIMAL_KHR) {
		return true, nil
	}

	if result != Success {
		return false, NewVulkanError(result, "QueuePresent", "Vulkan queue present failed")
	}

	return false, nil
}
