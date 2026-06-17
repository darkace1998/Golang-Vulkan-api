package vulkan

import "unsafe"

const (
	testDeviceParameter     = "device"
	testCreateInfoParameter = "createInfo"
	testMemoryParameter     = "memory"
	testNilCreateInfo       = "nil createInfo"
	testNilDevice           = "nil device"
	testNilMemory           = "nil memory"
	testNilEvent            = "nil event"
	testEventParameter      = "event"
	testValidationErrorType = "ValidationError"
)

// fakeHandle creates a non-nil handle backed by real Go memory, safe for use
// with -race checkptr. It must never be passed to actual Vulkan C functions.
func fakeHandle() unsafe.Pointer {
	x := new(uint64)
	return unsafe.Pointer(x)
}

func fakeDevice() Device                           { return Device(fakeHandle()) }
func fakeBuffer() Buffer                           { return Buffer(fakeHandle()) }
func fakeImage() Image                             { return Image(fakeHandle()) }
func fakeDeviceMemory() DeviceMemory               { return DeviceMemory(fakeHandle()) }
func fakeCommandPool() CommandPool                 { return CommandPool(fakeHandle()) }
func fakeCommandBuffer() CommandBuffer             { return CommandBuffer(fakeHandle()) }
func fakeFence() Fence                             { return Fence(fakeHandle()) }
func fakeSemaphore() Semaphore                     { return Semaphore(fakeHandle()) }
func fakeShaderModule() ShaderModule               { return ShaderModule(fakeHandle()) }
func fakePipelineLayout() PipelineLayout           { return PipelineLayout(fakeHandle()) }
func fakeRenderPass() RenderPass                   { return RenderPass(fakeHandle()) }
func fakePipeline() Pipeline                       { return Pipeline(fakeHandle()) }
func fakeImageView() ImageView                     { return ImageView(fakeHandle()) }
func fakeSampler() Sampler                         { return Sampler(fakeHandle()) }
func fakeDescriptorSetLayout() DescriptorSetLayout { return DescriptorSetLayout(fakeHandle()) }
func fakeDescriptorPool() DescriptorPool           { return DescriptorPool(fakeHandle()) }
func fakePhysicalDevice() PhysicalDevice           { return PhysicalDevice(fakeHandle()) }
func fakeVideoSession() VideoSession               { return VideoSession(fakeHandle()) }
func fakeSwapchain() Swapchain                     { return Swapchain(fakeHandle()) }
func fakeEvent() Event                             { return Event(fakeHandle()) }
