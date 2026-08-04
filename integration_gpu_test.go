package vulkan

import (
	"os"
	"testing"
	"unsafe"
)

// This file contains integration tests that exercise real GPU work end-to-end
// against a software Vulkan implementation (Lavapipe in CI): a compute
// dispatch whose output is read back and verified, a render pass clear whose
// pixels are read back and verified, and timestamp queries. These paths cover
// pipeline creation, descriptor updates, memory mapping, command recording,
// barriers, and queue submission — the code paths a unit test with fake
// handles can never reach.

// doubleCompSpirv is the SPIR-V for the following GLSL compute shader,
// compiled with glslangValidator -V:
//
//	#version 450
//	layout(local_size_x = 64) in;
//	layout(set = 0, binding = 0) buffer InBuf  { uint inData[]; };
//	layout(set = 0, binding = 1) buffer OutBuf { uint outData[]; };
//	void main() {
//	    uint i = gl_GlobalInvocationID.x;
//	    outData[i] = inData[i] * 2u + 1u;
//	}
var doubleCompSpirv = []uint32{
	0x07230203, 0x00010000, 0x0008000b, 0x00000026, 0x00000000, 0x00020011, 0x00000001, 0x0006000b,
	0x00000001, 0x4c534c47, 0x6474732e, 0x3035342e, 0x00000000, 0x0003000e, 0x00000000, 0x00000001,
	0x0006000f, 0x00000005, 0x00000004, 0x6e69616d, 0x00000000, 0x0000000b, 0x00060010, 0x00000004,
	0x00000011, 0x00000040, 0x00000001, 0x00000001, 0x00030003, 0x00000002, 0x000001c2, 0x00040005,
	0x00000004, 0x6e69616d, 0x00000000, 0x00030005, 0x00000008, 0x00000069, 0x00080005, 0x0000000b,
	0x475f6c67, 0x61626f6c, 0x766e496c, 0x7461636f, 0x496e6f69, 0x00000044, 0x00040005, 0x00000011,
	0x4274754f, 0x00006675, 0x00050006, 0x00000011, 0x00000000, 0x4474756f, 0x00617461, 0x00030005,
	0x00000013, 0x00000000, 0x00040005, 0x00000018, 0x75426e49, 0x00000066, 0x00050006, 0x00000018,
	0x00000000, 0x61446e69, 0x00006174, 0x00030005, 0x0000001a, 0x00000000, 0x00040047, 0x0000000b,
	0x0000000b, 0x0000001c, 0x00040047, 0x00000010, 0x00000006, 0x00000004, 0x00030047, 0x00000011,
	0x00000003, 0x00050048, 0x00000011, 0x00000000, 0x00000023, 0x00000000, 0x00040047, 0x00000013,
	0x00000021, 0x00000001, 0x00040047, 0x00000013, 0x00000022, 0x00000000, 0x00040047, 0x00000017,
	0x00000006, 0x00000004, 0x00030047, 0x00000018, 0x00000003, 0x00050048, 0x00000018, 0x00000000,
	0x00000023, 0x00000000, 0x00040047, 0x0000001a, 0x00000021, 0x00000000, 0x00040047, 0x0000001a,
	0x00000022, 0x00000000, 0x00040047, 0x00000025, 0x0000000b, 0x00000019, 0x00020013, 0x00000002,
	0x00030021, 0x00000003, 0x00000002, 0x00040015, 0x00000006, 0x00000020, 0x00000000, 0x00040020,
	0x00000007, 0x00000007, 0x00000006, 0x00040017, 0x00000009, 0x00000006, 0x00000003, 0x00040020,
	0x0000000a, 0x00000001, 0x00000009, 0x0004003b, 0x0000000a, 0x0000000b, 0x00000001, 0x0004002b,
	0x00000006, 0x0000000c, 0x00000000, 0x00040020, 0x0000000d, 0x00000001, 0x00000006, 0x0003001d,
	0x00000010, 0x00000006, 0x0003001e, 0x00000011, 0x00000010, 0x00040020, 0x00000012, 0x00000002,
	0x00000011, 0x0004003b, 0x00000012, 0x00000013, 0x00000002, 0x00040015, 0x00000014, 0x00000020,
	0x00000001, 0x0004002b, 0x00000014, 0x00000015, 0x00000000, 0x0003001d, 0x00000017, 0x00000006,
	0x0003001e, 0x00000018, 0x00000017, 0x00040020, 0x00000019, 0x00000002, 0x00000018, 0x0004003b,
	0x00000019, 0x0000001a, 0x00000002, 0x00040020, 0x0000001c, 0x00000002, 0x00000006, 0x0004002b,
	0x00000006, 0x0000001f, 0x00000002, 0x0004002b, 0x00000006, 0x00000021, 0x00000001, 0x0004002b,
	0x00000006, 0x00000024, 0x00000040, 0x0006002c, 0x00000009, 0x00000025, 0x00000024, 0x00000021,
	0x00000021, 0x00050036, 0x00000002, 0x00000004, 0x00000000, 0x00000003, 0x000200f8, 0x00000005,
	0x0004003b, 0x00000007, 0x00000008, 0x00000007, 0x00050041, 0x0000000d, 0x0000000e, 0x0000000b,
	0x0000000c, 0x0004003d, 0x00000006, 0x0000000f, 0x0000000e, 0x0003003e, 0x00000008, 0x0000000f,
	0x0004003d, 0x00000006, 0x00000016, 0x00000008, 0x0004003d, 0x00000006, 0x0000001b, 0x00000008,
	0x00060041, 0x0000001c, 0x0000001d, 0x0000001a, 0x00000015, 0x0000001b, 0x0004003d, 0x00000006,
	0x0000001e, 0x0000001d, 0x00050084, 0x00000006, 0x00000020, 0x0000001e, 0x0000001f, 0x00050080,
	0x00000006, 0x00000022, 0x00000020, 0x00000021, 0x00060041, 0x0000001c, 0x00000023, 0x00000013,
	0x00000015, 0x00000016, 0x0003003e, 0x00000023, 0x00000022, 0x000100fd, 0x00010038,
}

// integrationDevice bundles the objects every GPU integration test needs.
type integrationDevice struct {
	instance       Instance
	physicalDevice PhysicalDevice
	device         Device
	queue          Queue
	memProps       PhysicalDeviceMemoryProperties
}

// setupIntegrationDevice creates an instance, device, and queue, registering
// cleanup on t. Tests are skipped unless RUN_INTEGRATION_TESTS=1.
func setupIntegrationDevice(t *testing.T) *integrationDevice {
	t.Helper()
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration tests; set RUN_INTEGRATION_TESTS=1 to enable")
	}

	instance, err := CreateInstance(&InstanceCreateInfo{
		ApplicationInfo: &ApplicationInfo{
			ApplicationName: "GPUIntegrationTest",
			EngineName:      "NoEngine",
			APIVersion:      Version13,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	t.Cleanup(func() { DestroyInstance(instance) })

	physicalDevices, err := EnumeratePhysicalDevices(instance)
	if err != nil || len(physicalDevices) == 0 {
		t.Fatalf("Failed to enumerate physical devices: %v", err)
	}
	physicalDevice := physicalDevices[0]

	device, err := CreateDevice(physicalDevice, &DeviceCreateInfo{
		QueueCreateInfos: []DeviceQueueCreateInfo{
			{QueueFamilyIndex: 0, QueuePriorities: []float32{1.0}},
		},
		EnableTimelineSemaphores: true,
	})
	if err != nil {
		t.Fatalf("Failed to create logical device: %v", err)
	}
	t.Cleanup(func() { DestroyDevice(device) })

	return &integrationDevice{
		instance:       instance,
		physicalDevice: physicalDevice,
		device:         device,
		queue:          GetDeviceQueue(device, 0, 0),
		memProps:       GetPhysicalDeviceMemoryProperties(physicalDevice),
	}
}

// createHostBuffer creates a buffer bound to host-visible, host-coherent
// memory and registers cleanup on t.
func (env *integrationDevice) createHostBuffer(t *testing.T, size DeviceSize, usage BufferUsageFlags) (Buffer, DeviceMemory) {
	t.Helper()

	buffer, err := CreateBuffer(env.device, &BufferCreateInfo{Size: size, Usage: usage})
	if err != nil {
		t.Fatalf("CreateBuffer failed: %v", err)
	}
	t.Cleanup(func() { DestroyBuffer(env.device, buffer) })

	memReqs := GetBufferMemoryRequirements(env.device, buffer)
	memTypeIndex, found := FindMemoryType(env.memProps, memReqs.MemoryTypeBits,
		MemoryPropertyHostVisibleBit|MemoryPropertyHostCoherentBit)
	if !found {
		t.Fatal("No host-visible, host-coherent memory type found")
	}

	memory, err := AllocateMemory(env.device, &MemoryAllocateInfo{
		AllocationSize:  memReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	})
	if err != nil {
		t.Fatalf("AllocateMemory failed: %v", err)
	}
	t.Cleanup(func() { FreeMemory(env.device, memory) })

	if err := BindBufferMemory(env.device, buffer, memory, 0); err != nil {
		t.Fatalf("BindBufferMemory failed: %v", err)
	}
	return buffer, memory
}

// recordAndSubmit allocates a primary command buffer from a fresh pool,
// records commands via the callback, submits, and waits for completion.
func (env *integrationDevice) recordAndSubmit(t *testing.T, record func(cb CommandBuffer)) {
	t.Helper()

	commandPool, err := CreateCommandPool(env.device, &CommandPoolCreateInfo{QueueFamilyIndex: 0})
	if err != nil {
		t.Fatalf("CreateCommandPool failed: %v", err)
	}
	t.Cleanup(func() { DestroyCommandPool(env.device, commandPool) })

	commandBuffers, err := AllocateCommandBuffers(env.device, &CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		t.Fatalf("AllocateCommandBuffers failed: %v", err)
	}
	cb := commandBuffers[0]

	if err := BeginCommandBuffer(cb, &CommandBufferBeginInfo{Flags: CommandBufferUsageOneTimeSubmitBit}); err != nil {
		t.Fatalf("BeginCommandBuffer failed: %v", err)
	}
	record(cb)
	if err := EndCommandBuffer(cb); err != nil {
		t.Fatalf("EndCommandBuffer failed: %v", err)
	}

	fence, err := CreateFence(env.device, &FenceCreateInfo{})
	if err != nil {
		t.Fatalf("CreateFence failed: %v", err)
	}
	t.Cleanup(func() { DestroyFence(env.device, fence) })

	err = QueueSubmit(env.queue, []SubmitInfo{{CommandBuffers: []CommandBuffer{cb}}}, fence)
	if err != nil {
		t.Fatalf("QueueSubmit failed: %v", err)
	}
	result, err := WaitForFences(env.device, []Fence{fence}, true, 5_000_000_000)
	if err != nil {
		t.Fatalf("WaitForFences failed: %v", err)
	}
	if result != Success {
		t.Fatalf("WaitForFences returned %v, want Success", result)
	}
}

// TestIntegrationComputeDispatch runs a real compute shader that computes
// out[i] = in[i]*2+1 over a storage buffer and verifies every output value.
// It covers shader module creation, descriptor set layout/pool/allocation,
// descriptor writes, compute pipeline creation, dispatch, and memory mapping.
func TestIntegrationComputeDispatch(t *testing.T) {
	env := setupIntegrationDevice(t)

	const elemCount = 64
	const bufSize = DeviceSize(elemCount * 4)

	inBuffer, inMemory := env.createHostBuffer(t, bufSize, BufferUsageStorageBufferBit)
	outBuffer, outMemory := env.createHostBuffer(t, bufSize, BufferUsageStorageBufferBit)

	// Fill the input buffer with 0..63.
	ptr, err := MapMemory(env.device, inMemory, 0, bufSize, 0)
	if err != nil {
		t.Fatalf("MapMemory(in) failed: %v", err)
	}
	inData := unsafe.Slice((*uint32)(ptr), elemCount)
	for i := uint32(0); i < elemCount; i++ {
		inData[i] = i
	}
	UnmapMemory(env.device, inMemory)

	// Shader module from the embedded SPIR-V.
	const spirvWordBytes = 4
	shaderModule, err := CreateShaderModule(env.device, &ShaderModuleCreateInfo{
		CodeSize: uint32(len(doubleCompSpirv)) * spirvWordBytes, //nolint:gosec // embedded shader is ~1 KiB
		Code:     doubleCompSpirv,
	})
	if err != nil {
		t.Fatalf("CreateShaderModule failed: %v", err)
	}
	defer DestroyShaderModule(env.device, shaderModule)

	// Descriptor set layout: two storage buffers.
	setLayout, err := CreateDescriptorSetLayout(env.device, &DescriptorSetLayoutCreateInfo{
		Bindings: []DescriptorSetLayoutBinding{
			{Binding: 0, DescriptorType: DescriptorTypeStorageBuffer, DescriptorCount: 1, StageFlags: ShaderStageComputeBit},
			{Binding: 1, DescriptorType: DescriptorTypeStorageBuffer, DescriptorCount: 1, StageFlags: ShaderStageComputeBit},
		},
	})
	if err != nil {
		t.Fatalf("CreateDescriptorSetLayout failed: %v", err)
	}
	defer DestroyDescriptorSetLayout(env.device, setLayout)

	pipelineLayout, err := CreatePipelineLayout(env.device, &PipelineLayoutCreateInfo{
		SetLayouts: []DescriptorSetLayout{setLayout},
	})
	if err != nil {
		t.Fatalf("CreatePipelineLayout failed: %v", err)
	}
	defer DestroyPipelineLayout(env.device, pipelineLayout)

	pipelines, err := CreateComputePipelines(env.device, nil, []ComputePipelineCreateInfo{
		{
			Stage: PipelineShaderStageCreateInfo{
				Stage:  ShaderStageComputeBit,
				Module: shaderModule,
				Name:   "main",
			},
			Layout: pipelineLayout,
		},
	})
	if err != nil {
		t.Fatalf("CreateComputePipelines failed: %v", err)
	}
	defer DestroyPipeline(env.device, pipelines[0])

	// Descriptor pool + set, and point the bindings at the buffers.
	descriptorPool, err := CreateDescriptorPool(env.device, &DescriptorPoolCreateInfo{
		MaxSets: 1,
		PoolSizes: []DescriptorPoolSize{
			{Type: DescriptorTypeStorageBuffer, DescriptorCount: 2},
		},
	})
	if err != nil {
		t.Fatalf("CreateDescriptorPool failed: %v", err)
	}
	defer DestroyDescriptorPool(env.device, descriptorPool)

	descriptorSets, err := AllocateDescriptorSets(env.device, &DescriptorSetAllocateInfo{
		DescriptorPool: descriptorPool,
		SetLayouts:     []DescriptorSetLayout{setLayout},
	})
	if err != nil {
		t.Fatalf("AllocateDescriptorSets failed: %v", err)
	}

	UpdateDescriptorSets(env.device, []WriteDescriptorSet{
		{
			DstSet:          descriptorSets[0],
			DstBinding:      0,
			DescriptorCount: 1,
			DescriptorType:  DescriptorTypeStorageBuffer,
			BufferInfo:      []DescriptorBufferInfo{{Buffer: inBuffer, Offset: 0, Range: bufSize}},
		},
		{
			DstSet:          descriptorSets[0],
			DstBinding:      1,
			DescriptorCount: 1,
			DescriptorType:  DescriptorTypeStorageBuffer,
			BufferInfo:      []DescriptorBufferInfo{{Buffer: outBuffer, Offset: 0, Range: bufSize}},
		},
	}, nil)

	env.recordAndSubmit(t, func(cb CommandBuffer) {
		CmdBindPipeline(cb, PipelineBindPointCompute, pipelines[0])
		CmdBindDescriptorSets(cb, PipelineBindPointCompute, pipelineLayout, 0,
			[]DescriptorSet{descriptorSets[0]}, nil)
		CmdDispatch(cb, 1, 1, 1)
		CmdPipelineBarrierFull(cb, PipelineStageComputeShaderBit, PipelineStageHostBit, 0,
			nil,
			[]BufferMemoryBarrier{{
				SrcAccessMask:       AccessShaderWriteBit,
				DstAccessMask:       AccessHostReadBit,
				SrcQueueFamilyIndex: QueueFamilyIgnored,
				DstQueueFamilyIndex: QueueFamilyIgnored,
				Buffer:              outBuffer,
				Offset:              0,
				Size:                uint64(bufSize),
			}},
			nil)
	})

	// Read back and verify every element.
	ptr, err = MapMemory(env.device, outMemory, 0, bufSize, 0)
	if err != nil {
		t.Fatalf("MapMemory(out) failed: %v", err)
	}
	outData := unsafe.Slice((*uint32)(ptr), elemCount)
	for i := uint32(0); i < elemCount; i++ {
		want := i*2 + 1
		if outData[i] != want {
			t.Errorf("outData[%d] = %d, want %d", i, outData[i], want)
		}
	}
	UnmapMemory(env.device, outMemory)
}

// TestIntegrationRenderPassClear creates an offscreen color image, clears it
// through a render pass (CmdBeginRenderPass with clear values), copies the
// image to a buffer, and verifies the resulting pixels. It covers render pass
// and framebuffer creation, image views, image barriers, and image-to-buffer
// copies.
func TestIntegrationRenderPassClear(t *testing.T) {
	env := setupIntegrationDevice(t)

	const width, height = 4, 4
	const pixelBytes = 4
	const bufSize = DeviceSize(width * height * pixelBytes)

	// Offscreen color image.
	image, err := CreateImage(env.device, &ImageCreateInfo{
		ImageType:     ImageType2D,
		Format:        FormatR8G8B8A8Unorm,
		Extent:        Extent3D{Width: width, Height: height, Depth: 1},
		MipLevels:     1,
		ArrayLayers:   1,
		Samples:       SampleCount1Bit,
		Tiling:        ImageTilingOptimal,
		Usage:         ImageUsageColorAttachmentBit | ImageUsageTransferSrcBit,
		InitialLayout: ImageLayoutUndefined,
	})
	if err != nil {
		t.Fatalf("CreateImage failed: %v", err)
	}
	defer DestroyImage(env.device, image)

	imgReqs := GetImageMemoryRequirements(env.device, image)
	memTypeIndex, found := FindMemoryType(env.memProps, imgReqs.MemoryTypeBits, MemoryPropertyDeviceLocalBit)
	if !found {
		// Software renderers may not advertise DEVICE_LOCAL; take any type.
		memTypeIndex, found = FindMemoryType(env.memProps, imgReqs.MemoryTypeBits, 0)
		if !found {
			t.Fatal("No suitable memory type for image")
		}
	}
	imageMemory, err := AllocateMemory(env.device, &MemoryAllocateInfo{
		AllocationSize:  imgReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	})
	if err != nil {
		t.Fatalf("AllocateMemory(image) failed: %v", err)
	}
	defer FreeMemory(env.device, imageMemory)
	if err := BindImageMemory(env.device, image, imageMemory, 0); err != nil {
		t.Fatalf("BindImageMemory failed: %v", err)
	}

	imageView, err := CreateImageView(env.device, &ImageViewCreateInfo{
		Image:    image,
		ViewType: ImageViewType2D,
		Format:   FormatR8G8B8A8Unorm,
		SubresourceRange: ImageSubresourceRange{
			AspectMask: ImageAspectColorBit,
			LevelCount: 1,
			LayerCount: 1,
		},
	})
	if err != nil {
		t.Fatalf("CreateImageView failed: %v", err)
	}
	defer DestroyImageView(env.device, imageView)

	// Render pass that clears the attachment and leaves it TRANSFER_SRC.
	renderPass, err := CreateRenderPass(env.device, &RenderPassCreateInfo{
		Attachments: []AttachmentDescription{
			{
				Format:         FormatR8G8B8A8Unorm,
				Samples:        SampleCount1Bit,
				LoadOp:         AttachmentLoadOpClear,
				StoreOp:        AttachmentStoreOpStore,
				StencilLoadOp:  AttachmentLoadOpDontCare,
				StencilStoreOp: AttachmentStoreOpDontCare,
				InitialLayout:  ImageLayoutUndefined,
				FinalLayout:    ImageLayoutTransferSrcOptimal,
			},
		},
		Subpasses: []SubpassDescription{
			{
				PipelineBindPoint: PipelineBindPointGraphics,
				ColorAttachments: []AttachmentReference{
					{Attachment: 0, Layout: ImageLayoutColorAttachmentOptimal},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateRenderPass failed: %v", err)
	}
	defer DestroyRenderPass(env.device, renderPass)

	framebuffer, err := CreateFramebuffer(env.device, &FramebufferCreateInfo{
		RenderPass:  renderPass,
		Attachments: []ImageView{imageView},
		Width:       width,
		Height:      height,
		Layers:      1,
	})
	if err != nil {
		t.Fatalf("CreateFramebuffer failed: %v", err)
	}
	defer DestroyFramebuffer(env.device, framebuffer)

	readbackBuffer, readbackMemory := env.createHostBuffer(t, bufSize, BufferUsageTransferDstBit)

	// Clear to opaque red via the render pass, then copy to the buffer.
	env.recordAndSubmit(t, func(cb CommandBuffer) {
		CmdBeginRenderPass(cb, &RenderPassBeginInfo{
			RenderPass:  renderPass,
			Framebuffer: framebuffer,
			RenderArea:  Rect2D{Extent: Extent2D{Width: width, Height: height}},
			ClearValues: []ClearValue{
				{Color: ClearColorValue{Float32: [4]float32{1.0, 0.0, 0.0, 1.0}}},
			},
		}, SubpassContentsInline)
		CmdEndRenderPass(cb)

		// Make the attachment write visible to the transfer read; the implicit
		// end-of-render-pass dependency alone does not cover the copy below.
		CmdPipelineBarrierFull(cb, PipelineStageColorAttachmentOutputBit, PipelineStageTransferBit, 0,
			nil, nil,
			[]ImageMemoryBarrier{{
				SrcAccessMask:       AccessColorAttachmentWriteBit,
				DstAccessMask:       AccessTransferReadBit,
				OldLayout:           ImageLayoutTransferSrcOptimal,
				NewLayout:           ImageLayoutTransferSrcOptimal,
				SrcQueueFamilyIndex: QueueFamilyIgnored,
				DstQueueFamilyIndex: QueueFamilyIgnored,
				Image:               image,
				SubresourceRange: ImageSubresourceRange{
					AspectMask: ImageAspectColorBit,
					LevelCount: 1,
					LayerCount: 1,
				},
			}})

		CmdCopyImageToBuffer(cb, image, ImageLayoutTransferSrcOptimal, readbackBuffer,
			[]BufferImageCopy{
				{
					ImageSubresource: ImageSubresourceLayers{
						AspectMask: ImageAspectColorBit,
						LayerCount: 1,
					},
					ImageExtent: Extent3D{Width: width, Height: height, Depth: 1},
				},
			})

		CmdPipelineBarrierFull(cb, PipelineStageTransferBit, PipelineStageHostBit, 0,
			nil,
			[]BufferMemoryBarrier{{
				SrcAccessMask:       AccessTransferWriteBit,
				DstAccessMask:       AccessHostReadBit,
				SrcQueueFamilyIndex: QueueFamilyIgnored,
				DstQueueFamilyIndex: QueueFamilyIgnored,
				Buffer:              readbackBuffer,
				Offset:              0,
				Size:                uint64(bufSize),
			}},
			nil)
	})

	// Every pixel must be opaque red: R=255, G=0, B=0, A=255.
	ptr, err := MapMemory(env.device, readbackMemory, 0, bufSize, 0)
	if err != nil {
		t.Fatalf("MapMemory(readback) failed: %v", err)
	}
	pixels := unsafe.Slice((*byte)(ptr), int(bufSize))
	for i := 0; i < len(pixels); i += pixelBytes {
		r, g, b, a := pixels[i], pixels[i+1], pixels[i+2], pixels[i+3]
		if r != 255 || g != 0 || b != 0 || a != 255 {
			t.Fatalf("pixel %d = (%d,%d,%d,%d), want (255,0,0,255)", i/pixelBytes, r, g, b, a)
		}
	}
	UnmapMemory(env.device, readbackMemory)
}

// TestIntegrationClearColorImageUint verifies that integer clear colors are
// honored: it clears an R8G8B8A8_UINT image via CmdClearColorImage with the
// Uint32 member of the clear-color union and reads the values back. Before
// clearColorBits existed, the Uint32 member was silently ignored and the
// image would have been cleared to zero.
func TestIntegrationClearColorImageUint(t *testing.T) {
	env := setupIntegrationDevice(t)

	const width, height = 4, 4
	const pixelBytes = 4
	const bufSize = DeviceSize(width * height * pixelBytes)

	image, err := CreateImage(env.device, &ImageCreateInfo{
		ImageType:     ImageType2D,
		Format:        FormatR8G8B8A8Uint,
		Extent:        Extent3D{Width: width, Height: height, Depth: 1},
		MipLevels:     1,
		ArrayLayers:   1,
		Samples:       SampleCount1Bit,
		Tiling:        ImageTilingOptimal,
		Usage:         ImageUsageTransferDstBit | ImageUsageTransferSrcBit,
		InitialLayout: ImageLayoutUndefined,
	})
	if err != nil {
		t.Fatalf("CreateImage failed: %v", err)
	}
	defer DestroyImage(env.device, image)

	imgReqs := GetImageMemoryRequirements(env.device, image)
	memTypeIndex, found := FindMemoryType(env.memProps, imgReqs.MemoryTypeBits, 0)
	if !found {
		t.Fatal("No suitable memory type for image")
	}
	imageMemory, err := AllocateMemory(env.device, &MemoryAllocateInfo{
		AllocationSize:  imgReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	})
	if err != nil {
		t.Fatalf("AllocateMemory(image) failed: %v", err)
	}
	defer FreeMemory(env.device, imageMemory)
	if err := BindImageMemory(env.device, image, imageMemory, 0); err != nil {
		t.Fatalf("BindImageMemory failed: %v", err)
	}

	readbackBuffer, readbackMemory := env.createHostBuffer(t, bufSize, BufferUsageTransferDstBit)

	fullRange := ImageSubresourceRange{
		AspectMask: ImageAspectColorBit,
		LevelCount: 1,
		LayerCount: 1,
	}
	clearColor := ClearColorValue{Uint32: [4]uint32{7, 42, 99, 255}}

	env.recordAndSubmit(t, func(cb CommandBuffer) {
		TransitionImageLayout(cb, image, FormatR8G8B8A8Uint,
			ImageLayoutUndefined, ImageLayoutTransferDstOptimal, fullRange)

		CmdClearColorImage(cb, image, ImageLayoutTransferDstOptimal, &clearColor,
			[]ImageSubresourceRange{fullRange})

		// TRANSFER_DST -> TRANSFER_SRC for the readback copy.
		CmdPipelineBarrierFull(cb, PipelineStageTransferBit, PipelineStageTransferBit, 0,
			nil, nil,
			[]ImageMemoryBarrier{{
				SrcAccessMask:       AccessTransferWriteBit,
				DstAccessMask:       AccessTransferReadBit,
				OldLayout:           ImageLayoutTransferDstOptimal,
				NewLayout:           ImageLayoutTransferSrcOptimal,
				SrcQueueFamilyIndex: QueueFamilyIgnored,
				DstQueueFamilyIndex: QueueFamilyIgnored,
				Image:               image,
				SubresourceRange:    fullRange,
			}})

		CmdCopyImageToBuffer(cb, image, ImageLayoutTransferSrcOptimal, readbackBuffer,
			[]BufferImageCopy{
				{
					ImageSubresource: ImageSubresourceLayers{
						AspectMask: ImageAspectColorBit,
						LayerCount: 1,
					},
					ImageExtent: Extent3D{Width: width, Height: height, Depth: 1},
				},
			})

		CmdPipelineBarrierFull(cb, PipelineStageTransferBit, PipelineStageHostBit, 0,
			nil,
			[]BufferMemoryBarrier{{
				SrcAccessMask:       AccessTransferWriteBit,
				DstAccessMask:       AccessHostReadBit,
				SrcQueueFamilyIndex: QueueFamilyIgnored,
				DstQueueFamilyIndex: QueueFamilyIgnored,
				Buffer:              readbackBuffer,
				Offset:              0,
				Size:                uint64(bufSize),
			}},
			nil)
	})

	ptr, err := MapMemory(env.device, readbackMemory, 0, bufSize, 0)
	if err != nil {
		t.Fatalf("MapMemory(readback) failed: %v", err)
	}
	pixels := unsafe.Slice((*byte)(ptr), int(bufSize))
	for i := 0; i < len(pixels); i += pixelBytes {
		r, g, b, a := pixels[i], pixels[i+1], pixels[i+2], pixels[i+3]
		if r != 7 || g != 42 || b != 99 || a != 255 {
			t.Fatalf("pixel %d = (%d,%d,%d,%d), want (7,42,99,255)", i/pixelBytes, r, g, b, a)
		}
	}
	UnmapMemory(env.device, readbackMemory)
}

// TestIntegrationWaitPolling verifies the non-blocking polling idiom: waiting
// on an unsignaled fence or an unreached timeline value with timeout=0 must
// return Timeout with a nil error (VK_TIMEOUT is a success code), and must
// return Success once the condition is met.
func TestIntegrationWaitPolling(t *testing.T) {
	env := setupIntegrationDevice(t)

	// Unsignaled fence: poll must report Timeout without an error.
	fence, err := CreateFence(env.device, &FenceCreateInfo{})
	if err != nil {
		t.Fatalf("CreateFence failed: %v", err)
	}
	defer DestroyFence(env.device, fence)

	result, err := WaitForFences(env.device, []Fence{fence}, true, 0)
	if err != nil {
		t.Fatalf("Polling an unsignaled fence returned an error: %v", err)
	}
	if result != Timeout {
		t.Fatalf("Polling an unsignaled fence returned %v, want Timeout", result)
	}

	// Pre-signaled fence: poll must report Success.
	signaled, err := CreateFence(env.device, &FenceCreateInfo{Flags: FenceCreateSignaledBit})
	if err != nil {
		t.Fatalf("CreateFence(signaled) failed: %v", err)
	}
	defer DestroyFence(env.device, signaled)

	result, err = WaitForFences(env.device, []Fence{signaled}, true, 0)
	if err != nil {
		t.Fatalf("Polling a signaled fence returned an error: %v", err)
	}
	if result != Success {
		t.Fatalf("Polling a signaled fence returned %v, want Success", result)
	}

	// Timeline semaphore at value 5: waiting for 10 with timeout=0 must
	// report Timeout without an error; waiting for 5 must report Success.
	timeline, err := CreateTimelineSemaphore(env.device, 5)
	if err != nil {
		t.Fatalf("CreateTimelineSemaphore failed: %v", err)
	}
	defer DestroySemaphore(env.device, timeline)

	result, err = WaitSemaphores(env.device, &SemaphoreWaitInfo{
		Semaphores: []Semaphore{timeline},
		Values:     []uint64{10},
	}, 0)
	if err != nil {
		t.Fatalf("Polling an unreached timeline value returned an error: %v", err)
	}
	if result != Timeout {
		t.Fatalf("Polling an unreached timeline value returned %v, want Timeout", result)
	}

	result, err = WaitSemaphores(env.device, &SemaphoreWaitInfo{
		Semaphores: []Semaphore{timeline},
		Values:     []uint64{5},
	}, 0)
	if err != nil {
		t.Fatalf("Polling a reached timeline value returned an error: %v", err)
	}
	if result != Success {
		t.Fatalf("Polling a reached timeline value returned %v, want Success", result)
	}
}

// TestIntegrationTimestampQueries writes two timestamps around GPU work and
// reads them back, covering query pool creation, CmdResetQueryPool,
// CmdWriteTimestamp, and GetQueryPoolResultsUint64.
func TestIntegrationTimestampQueries(t *testing.T) {
	env := setupIntegrationDevice(t)

	props := GetPhysicalDeviceProperties(env.physicalDevice)
	if props.Limits.TimestampPeriod == 0 {
		t.Skip("Device does not support timestamps")
	}

	queryPool, err := CreateQueryPool(env.device, &QueryPoolCreateInfo{
		QueryType:  QueryTypeTimestamp,
		QueryCount: 2,
	})
	if err != nil {
		t.Fatalf("CreateQueryPool failed: %v", err)
	}
	defer DestroyQueryPool(env.device, queryPool)

	env.recordAndSubmit(t, func(cb CommandBuffer) {
		CmdResetQueryPool(cb, queryPool, 0, 2)
		CmdWriteTimestamp(cb, PipelineStageTopOfPipeBit, queryPool, 0)
		CmdPipelineBarrier(cb, PipelineStageTopOfPipeBit, PipelineStageBottomOfPipeBit, 0)
		CmdWriteTimestamp(cb, PipelineStageBottomOfPipeBit, queryPool, 1)
	})

	results, queryResult, err := GetQueryPoolResultsUint64(env.device, queryPool, 0, 2, QueryResultWait)
	if err != nil {
		t.Fatalf("GetQueryPoolResultsUint64 failed: %v", err)
	}
	if queryResult != Success {
		t.Fatalf("GetQueryPoolResultsUint64 returned %v, want Success", queryResult)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	if results[1] < results[0] {
		t.Errorf("Timestamps out of order: t0=%d t1=%d", results[0], results[1])
	}
	t.Logf("Timestamps: t0=%d t1=%d (period %.2f ns)", results[0], results[1], props.Limits.TimestampPeriod)
}
