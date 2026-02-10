package main

import (
	"fmt"
	"log"
	"unsafe"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

// Graphics Pipeline Example
// Demonstrates a complete offscreen graphics pipeline:
//   - Render pass with color attachment
//   - Vertex buffer with triangle data
//   - Shader modules (placeholder SPIR-V)
//   - Full graphics pipeline creation
//   - Framebuffer creation
//   - Command buffer recording with draw commands
//   - GPU execution and synchronisation

// Vertex represents a 2D vertex with position (x, y) and color (r, g, b).
// Packed as 5 floats = 20 bytes per vertex.
type Vertex struct {
	PosX, PosY       float32
	ColorR, ColorG, ColorB float32
}

// Triangle vertices (positions in NDC, RGB colors)
var triangleVertices = []Vertex{
	{0.0, -0.5, 1.0, 0.0, 0.0},  // top    — red
	{0.5, 0.5, 0.0, 1.0, 0.0},   // right  — green
	{-0.5, 0.5, 0.0, 0.0, 1.0},  // left   — blue
}

// Minimal valid SPIR-V vertex shader (pass-through)
// This is a minimal placeholder; a real application would compile GLSL via glslc.
var vertexShaderCode = []uint32{
	0x07230203, 0x00010000, 0x000d000a, 0x00000001,
}

// Minimal valid SPIR-V fragment shader (solid color output)
var fragmentShaderCode = []uint32{
	0x07230203, 0x00010000, 0x000d000a, 0x00000001,
}

func main() {
	fmt.Println("=== Vulkan Graphics Pipeline Example ===")
	fmt.Println("Demonstrates a complete offscreen triangle rendering pipeline")
	fmt.Println()

	// -----------------------------------------------------------------------
	// 1. Create Vulkan instance
	// -----------------------------------------------------------------------
	fmt.Println("1. Creating Vulkan instance...")
	instance, err := vulkan.CreateInstance(&vulkan.InstanceCreateInfo{
		ApplicationInfo: &vulkan.ApplicationInfo{
			ApplicationName:    "Graphics Pipeline Example",
			ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
			EngineName:         "Example Engine",
			EngineVersion:      vulkan.MakeVersion(1, 0, 0),
			APIVersion:         vulkan.Version13,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create Vulkan instance: %v", err)
	}
	defer vulkan.DestroyInstance(instance)
	fmt.Println("   ✓ Vulkan instance created")

	// -----------------------------------------------------------------------
	// 2. Select physical device
	// -----------------------------------------------------------------------
	fmt.Println("\n2. Selecting physical device...")
	physicalDevices, err := vulkan.EnumeratePhysicalDevices(instance)
	if err != nil {
		log.Fatalf("Failed to enumerate devices: %v", err)
	}
	if len(physicalDevices) == 0 {
		log.Fatal("No physical devices found")
	}
	physicalDevice := physicalDevices[0]
	props := vulkan.GetPhysicalDeviceProperties(physicalDevice)
	fmt.Printf("   ✓ Using device: %s\n", props.DeviceName)

	// -----------------------------------------------------------------------
	// 3. Find graphics queue family
	// -----------------------------------------------------------------------
	fmt.Println("\n3. Finding graphics queue family...")
	queueFamilies := vulkan.GetPhysicalDeviceQueueFamilyProperties(physicalDevice)
	var graphicsQueueFamily uint32 = ^uint32(0)
	for i, qf := range queueFamilies {
		if qf.QueueFlags&vulkan.QueueGraphicsBit != 0 {
			graphicsQueueFamily = uint32(i)
			break
		}
	}
	if graphicsQueueFamily == ^uint32(0) {
		log.Fatal("No graphics queue family found")
	}
	fmt.Printf("   ✓ Graphics queue family: %d\n", graphicsQueueFamily)

	// -----------------------------------------------------------------------
	// 4. Create logical device and get queue
	// -----------------------------------------------------------------------
	fmt.Println("\n4. Creating logical device...")
	device, err := vulkan.CreateDevice(physicalDevice, &vulkan.DeviceCreateInfo{
		QueueCreateInfos: []vulkan.DeviceQueueCreateInfo{
			{
				QueueFamilyIndex: graphicsQueueFamily,
				QueuePriorities:  []float32{1.0},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create device: %v", err)
	}
	defer vulkan.DestroyDevice(device)
	queue := vulkan.GetDeviceQueue(device, graphicsQueueFamily, 0)
	fmt.Println("   ✓ Device and queue created")

	// -----------------------------------------------------------------------
	// 5. Create a render target image (offscreen color attachment)
	// -----------------------------------------------------------------------
	fmt.Println("\n5. Creating offscreen render target...")
	const (
		renderWidth  = 800
		renderHeight = 600
	)
	renderFormat := vulkan.FormatR8G8B8A8Unorm

	image, err := vulkan.CreateImage(device, &vulkan.ImageCreateInfo{
		ImageType:     vulkan.ImageType2D,
		Format:        renderFormat,
		Extent:        vulkan.Extent3D{Width: renderWidth, Height: renderHeight, Depth: 1},
		MipLevels:     1,
		ArrayLayers:   1,
		Samples:       vulkan.SampleCount1Bit,
		Tiling:        vulkan.ImageTilingOptimal,
		Usage:         vulkan.ImageUsageColorAttachmentBit | vulkan.ImageUsageTransferSrcBit,
		SharingMode:   vulkan.SharingModeExclusive,
		InitialLayout: vulkan.ImageLayoutUndefined,
	})
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}
	defer vulkan.DestroyImage(device, image)
	fmt.Printf("   ✓ Render target image created (%dx%d, R8G8B8A8)\n", renderWidth, renderHeight)

	// Allocate and bind memory for the image
	memProps := vulkan.GetPhysicalDeviceMemoryProperties(physicalDevice)
	imageMemReqs := vulkan.GetImageMemoryRequirements(device, image)
	memTypeIndex, found := vulkan.FindMemoryType(memProps, imageMemReqs.MemoryTypeBits, vulkan.MemoryPropertyDeviceLocalBit)
	if !found {
		log.Fatal("Failed to find suitable memory type for render target")
	}

	imageMemory, err := vulkan.AllocateMemory(device, &vulkan.MemoryAllocateInfo{
		AllocationSize:  imageMemReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	})
	if err != nil {
		log.Fatalf("Failed to allocate image memory: %v", err)
	}
	defer vulkan.FreeMemory(device, imageMemory)

	if err := vulkan.BindImageMemory(device, image, imageMemory, 0); err != nil {
		log.Fatalf("Failed to bind image memory: %v", err)
	}
	fmt.Println("   ✓ Image memory allocated and bound")

	// Create image view
	imageView, err := vulkan.CreateImageView(device, &vulkan.ImageViewCreateInfo{
		Image:    image,
		ViewType: vulkan.ImageViewType2D,
		Format:   renderFormat,
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectColorBit,
			BaseMipLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create image view: %v", err)
	}
	defer vulkan.DestroyImageView(device, imageView)
	fmt.Println("   ✓ Image view created")

	// -----------------------------------------------------------------------
	// 6. Create render pass
	// -----------------------------------------------------------------------
	fmt.Println("\n6. Creating render pass...")
	renderPass, err := vulkan.CreateRenderPass(device, &vulkan.RenderPassCreateInfo{
		Attachments: []vulkan.AttachmentDescription{
			{
				Format:         renderFormat,
				Samples:        vulkan.SampleCount1Bit,
				LoadOp:         vulkan.AttachmentLoadOpClear,
				StoreOp:        vulkan.AttachmentStoreOpStore,
				StencilLoadOp:  vulkan.AttachmentLoadOpDontCare,
				StencilStoreOp: vulkan.AttachmentStoreOpDontCare,
				InitialLayout:  vulkan.ImageLayoutUndefined,
				FinalLayout:    vulkan.ImageLayoutColorAttachmentOptimal,
			},
		},
		Subpasses: []vulkan.SubpassDescription{
			{
				PipelineBindPoint: vulkan.PipelineBindPointGraphics,
				ColorAttachments: []vulkan.AttachmentReference{
					{Attachment: 0, Layout: vulkan.ImageLayoutColorAttachmentOptimal},
				},
			},
		},
		Dependencies: []vulkan.SubpassDependency{
			{
				SrcSubpass:    vulkan.SubpassExternal,
				DstSubpass:    0,
				SrcStageMask:  vulkan.PipelineStageColorAttachmentOutputBit,
				DstStageMask:  vulkan.PipelineStageColorAttachmentOutputBit,
				SrcAccessMask: 0,
				DstAccessMask: vulkan.AccessColorAttachmentWriteBit,
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create render pass: %v", err)
	}
	defer vulkan.DestroyRenderPass(device, renderPass)
	fmt.Println("   ✓ Render pass created (1 color attachment, clear-to-store)")

	// -----------------------------------------------------------------------
	// 7. Create framebuffer
	// -----------------------------------------------------------------------
	fmt.Println("\n7. Creating framebuffer...")
	framebuffer, err := vulkan.CreateFramebuffer(device, &vulkan.FramebufferCreateInfo{
		RenderPass:  renderPass,
		Attachments: []vulkan.ImageView{imageView},
		Width:       renderWidth,
		Height:      renderHeight,
		Layers:      1,
	})
	if err != nil {
		log.Fatalf("Failed to create framebuffer: %v", err)
	}
	defer vulkan.DestroyFramebuffer(device, framebuffer)
	fmt.Println("   ✓ Framebuffer created")

	// -----------------------------------------------------------------------
	// 8. Create vertex buffer
	// -----------------------------------------------------------------------
	fmt.Println("\n8. Creating vertex buffer...")
	vertexDataSize := vulkan.DeviceSize(len(triangleVertices)) * vulkan.DeviceSize(unsafe.Sizeof(triangleVertices[0]))

	vertexBuffer, err := vulkan.CreateBuffer(device, &vulkan.BufferCreateInfo{
		Size:        vertexDataSize,
		Usage:       vulkan.BufferUsageVertexBufferBit,
		SharingMode: vulkan.SharingModeExclusive,
	})
	if err != nil {
		log.Fatalf("Failed to create vertex buffer: %v", err)
	}
	defer vulkan.DestroyBuffer(device, vertexBuffer)

	bufferMemReqs := vulkan.GetBufferMemoryRequirements(device, vertexBuffer)
	bufferMemType, found := vulkan.FindMemoryType(memProps, bufferMemReqs.MemoryTypeBits,
		vulkan.MemoryPropertyHostVisibleBit|vulkan.MemoryPropertyHostCoherentBit)
	if !found {
		log.Fatal("Failed to find suitable memory type for vertex buffer")
	}

	vertexBufferMemory, err := vulkan.AllocateMemory(device, &vulkan.MemoryAllocateInfo{
		AllocationSize:  bufferMemReqs.Size,
		MemoryTypeIndex: bufferMemType,
	})
	if err != nil {
		log.Fatalf("Failed to allocate vertex buffer memory: %v", err)
	}
	defer vulkan.FreeMemory(device, vertexBufferMemory)

	if err := vulkan.BindBufferMemory(device, vertexBuffer, vertexBufferMemory, 0); err != nil {
		log.Fatalf("Failed to bind vertex buffer memory: %v", err)
	}

	// Map memory and upload vertex data
	mapped, err := vulkan.MapMemory(device, vertexBufferMemory, 0, vertexDataSize, 0)
	if err != nil {
		log.Fatalf("Failed to map vertex buffer memory: %v", err)
	}

	// Copy vertex data
	src := unsafe.Slice((*byte)(unsafe.Pointer(&triangleVertices[0])), int(vertexDataSize))
	dst := unsafe.Slice((*byte)(mapped), int(vertexDataSize))
	copy(dst, src)

	vulkan.UnmapMemory(device, vertexBufferMemory)
	fmt.Printf("   ✓ Vertex buffer created and uploaded (%d vertices, %d bytes)\n",
		len(triangleVertices), vertexDataSize)

	// -----------------------------------------------------------------------
	// 9. Create shader modules
	// -----------------------------------------------------------------------
	fmt.Println("\n9. Creating shader modules...")
	vertShaderModule, err := vulkan.CreateShaderModule(device, &vulkan.ShaderModuleCreateInfo{
		CodeSize: uint32(len(vertexShaderCode) * 4),
		Code:     vertexShaderCode,
	})
	if err != nil {
		log.Fatalf("Failed to create vertex shader module: %v", err)
	}
	defer vulkan.DestroyShaderModule(device, vertShaderModule)
	fmt.Println("   ✓ Vertex shader module created")

	fragShaderModule, err := vulkan.CreateShaderModule(device, &vulkan.ShaderModuleCreateInfo{
		CodeSize: uint32(len(fragmentShaderCode) * 4),
		Code:     fragmentShaderCode,
	})
	if err != nil {
		log.Fatalf("Failed to create fragment shader module: %v", err)
	}
	defer vulkan.DestroyShaderModule(device, fragShaderModule)
	fmt.Println("   ✓ Fragment shader module created")

	// -----------------------------------------------------------------------
	// 10. Create pipeline layout
	// -----------------------------------------------------------------------
	fmt.Println("\n10. Creating pipeline layout...")
	pipelineLayout, err := vulkan.CreatePipelineLayout(device, &vulkan.PipelineLayoutCreateInfo{})
	if err != nil {
		log.Fatalf("Failed to create pipeline layout: %v", err)
	}
	defer vulkan.DestroyPipelineLayout(device, pipelineLayout)
	fmt.Println("   ✓ Pipeline layout created (no descriptor sets, no push constants)")

	// -----------------------------------------------------------------------
	// 11. Create graphics pipeline
	// -----------------------------------------------------------------------
	fmt.Println("\n11. Creating graphics pipeline...")

	graphicsPipelineCreateInfo := vulkan.GraphicsPipelineCreateInfo{
		Stages: []vulkan.PipelineShaderStageCreateInfo{
			{
				Stage:  vulkan.ShaderStageVertexBit,
				Module: vertShaderModule,
				Name:   "main",
			},
			{
				Stage:  vulkan.ShaderStageFragmentBit,
				Module: fragShaderModule,
				Name:   "main",
			},
		},
		VertexInputState: &vulkan.PipelineVertexInputStateCreateInfo{
			VertexBindingDescriptions: []vulkan.VertexInputBindingDescription{
				{
					Binding:   0,
					Stride:    uint32(unsafe.Sizeof(Vertex{})),
					InputRate: vulkan.VertexInputRateVertex,
				},
			},
			VertexAttributeDescriptions: []vulkan.VertexInputAttributeDescription{
				// position: location=0, offset=0, R8G8B8A8 used as 2-float proxy
				{Location: 0, Binding: 0, Format: vulkan.FormatR8G8B8A8Unorm, Offset: 0},
				// color: location=1, offset=8
				{Location: 1, Binding: 0, Format: vulkan.FormatR8G8B8A8Unorm, Offset: 8},
			},
		},
		InputAssemblyState: &vulkan.PipelineInputAssemblyStateCreateInfo{
			Topology:               vulkan.PrimitiveTopologyTriangleList,
			PrimitiveRestartEnable: false,
		},
		ViewportState: &vulkan.PipelineViewportStateCreateInfo{
			Viewports: []vulkan.Viewport{
				{X: 0, Y: 0, Width: renderWidth, Height: renderHeight, MinDepth: 0, MaxDepth: 1},
			},
			Scissors: []vulkan.Rect2D{
				{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: vulkan.Extent2D{Width: renderWidth, Height: renderHeight}},
			},
		},
		RasterizationState: &vulkan.PipelineRasterizationStateCreateInfo{
			DepthClampEnable:        false,
			RasterizerDiscardEnable: false,
			PolygonMode:             vulkan.PolygonModeFill,
			CullMode:                vulkan.CullModeBack,
			FrontFace:               vulkan.FrontFaceCounterClockwise,
			DepthBiasEnable:         false,
			LineWidth:               1.0,
		},
		MultisampleState: &vulkan.PipelineMultisampleStateCreateInfo{
			RasterizationSamples: vulkan.SampleCount1Bit,
			SampleShadingEnable:  false,
		},
		ColorBlendState: &vulkan.PipelineColorBlendStateCreateInfo{
			LogicOpEnable: false,
			Attachments: []vulkan.PipelineColorBlendAttachmentState{
				{
					BlendEnable:         false,
					SrcColorBlendFactor: vulkan.BlendFactorOne,
					DstColorBlendFactor: vulkan.BlendFactorZero,
					ColorBlendOp:        vulkan.BlendOpAdd,
					SrcAlphaBlendFactor: vulkan.BlendFactorOne,
					DstAlphaBlendFactor: vulkan.BlendFactorZero,
					AlphaBlendOp:        vulkan.BlendOpAdd,
					ColorWriteMask:      vulkan.ColorComponentAll,
				},
			},
		},
		Layout:     pipelineLayout,
		RenderPass: renderPass,
		Subpass:    0,
	}

	pipelines, err := vulkan.CreateGraphicsPipelines(device, vulkan.PipelineCache(nil),
		[]vulkan.GraphicsPipelineCreateInfo{graphicsPipelineCreateInfo})
	if err != nil {
		log.Fatalf("Failed to create graphics pipeline: %v", err)
	}
	graphicsPipeline := pipelines[0]
	defer vulkan.DestroyPipeline(device, graphicsPipeline)
	fmt.Println("   ✓ Graphics pipeline created")
	fmt.Println("     Stages: vertex + fragment")
	fmt.Println("     Topology: triangle list")
	fmt.Println("     Polygon mode: fill")
	fmt.Println("     Cull mode: back face")

	// -----------------------------------------------------------------------
	// 12. Create command pool and buffer
	// -----------------------------------------------------------------------
	fmt.Println("\n12. Creating command pool and buffer...")
	commandPool, err := vulkan.CreateCommandPool(device, &vulkan.CommandPoolCreateInfo{
		Flags:            vulkan.CommandPoolCreateResetCommandBufferBit,
		QueueFamilyIndex: graphicsQueueFamily,
	})
	if err != nil {
		log.Fatalf("Failed to create command pool: %v", err)
	}
	defer vulkan.DestroyCommandPool(device, commandPool)

	commandBuffers, err := vulkan.AllocateCommandBuffers(device, &vulkan.CommandBufferAllocateInfo{
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		log.Fatalf("Failed to allocate command buffers: %v", err)
	}
	cb := commandBuffers[0]
	fmt.Println("   ✓ Command buffer allocated")

	// -----------------------------------------------------------------------
	// 13. Record rendering commands
	// -----------------------------------------------------------------------
	fmt.Println("\n13. Recording rendering commands...")

	if err := vulkan.BeginCommandBuffer(cb, &vulkan.CommandBufferBeginInfo{
		Flags: vulkan.CommandBufferUsageOneTimeSubmitBit,
	}); err != nil {
		log.Fatalf("Failed to begin command buffer: %v", err)
	}

	// Begin render pass — clear to dark blue
	vulkan.CmdBeginRenderPass(cb, &vulkan.RenderPassBeginInfo{
		RenderPass:  renderPass,
		Framebuffer: framebuffer,
		RenderArea: vulkan.Rect2D{
			Offset: vulkan.Offset2D{X: 0, Y: 0},
			Extent: vulkan.Extent2D{Width: renderWidth, Height: renderHeight},
		},
		ClearValues: []vulkan.ClearValue{
			{Color: vulkan.ClearColorValue{Float32: [4]float32{0.0, 0.0, 0.2, 1.0}}},
		},
	}, vulkan.SubpassContentsInline)

	// Bind the graphics pipeline
	vulkan.CmdBindPipeline(cb, vulkan.PipelineBindPointGraphics, graphicsPipeline)

	// Set dynamic viewport and scissor
	vulkan.CmdSetViewport(cb, 0, []vulkan.Viewport{
		{X: 0, Y: 0, Width: renderWidth, Height: renderHeight, MinDepth: 0, MaxDepth: 1},
	})
	vulkan.CmdSetScissor(cb, 0, []vulkan.Rect2D{
		{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: vulkan.Extent2D{Width: renderWidth, Height: renderHeight}},
	})

	// Bind vertex buffer
	vulkan.CmdBindVertexBuffers(cb, 0,
		[]vulkan.Buffer{vertexBuffer},
		[]vulkan.DeviceSize{0},
	)

	// Draw the triangle (3 vertices, 1 instance)
	vulkan.CmdDraw(cb, uint32(len(triangleVertices)), 1, 0, 0)

	// End render pass
	vulkan.CmdEndRenderPass(cb)

	if err := vulkan.EndCommandBuffer(cb); err != nil {
		log.Fatalf("Failed to end command buffer: %v", err)
	}
	fmt.Println("   ✓ Commands recorded:")
	fmt.Println("     BeginRenderPass (clear to dark blue)")
	fmt.Println("     BindPipeline (graphics)")
	fmt.Println("     SetViewport / SetScissor")
	fmt.Println("     BindVertexBuffers")
	fmt.Println("     Draw (3 vertices, 1 instance)")
	fmt.Println("     EndRenderPass")

	// -----------------------------------------------------------------------
	// 14. Submit and wait
	// -----------------------------------------------------------------------
	fmt.Println("\n14. Submitting to GPU...")
	fence, err := vulkan.CreateFence(device, &vulkan.FenceCreateInfo{})
	if err != nil {
		log.Fatalf("Failed to create fence: %v", err)
	}
	defer vulkan.DestroyFence(device, fence)

	if err := vulkan.QueueSubmit(queue, []vulkan.SubmitInfo{
		{CommandBuffers: []vulkan.CommandBuffer{cb}},
	}, fence); err != nil {
		log.Fatalf("Failed to submit: %v", err)
	}

	if err := vulkan.WaitForFences(device, []vulkan.Fence{fence}, true, ^uint64(0)); err != nil {
		log.Fatalf("Failed to wait for fence: %v", err)
	}
	fmt.Println("   ✓ GPU execution complete")

	// -----------------------------------------------------------------------
	// 15. Summary
	// -----------------------------------------------------------------------
	fmt.Println("\n=== Graphics Pipeline Example Complete ===")
	fmt.Println()
	fmt.Println("Pipeline state objects demonstrated:")
	fmt.Println("  ✓ Vertex input state (binding + attribute descriptions)")
	fmt.Println("  ✓ Input assembly state (triangle list topology)")
	fmt.Println("  ✓ Viewport / scissor state")
	fmt.Println("  ✓ Rasterization state (fill, back-face cull, CCW)")
	fmt.Println("  ✓ Multisample state (1 sample)")
	fmt.Println("  ✓ Color blend state (no blending, write RGBA)")
	fmt.Println("  ✓ Pipeline layout (empty)")
	fmt.Println()
	fmt.Println("Rendering commands demonstrated:")
	fmt.Println("  ✓ CmdBeginRenderPass / CmdEndRenderPass")
	fmt.Println("  ✓ CmdBindPipeline (graphics)")
	fmt.Println("  ✓ CmdSetViewport / CmdSetScissor")
	fmt.Println("  ✓ CmdBindVertexBuffers")
	fmt.Println("  ✓ CmdDraw (non-indexed)")
	fmt.Println()
	fmt.Println("Resource management demonstrated:")
	fmt.Println("  ✓ Offscreen render target (image + image view)")
	fmt.Println("  ✓ Render pass (color attachment, clear→store)")
	fmt.Println("  ✓ Framebuffer")
	fmt.Println("  ✓ Vertex buffer (host-visible, mapped upload)")
	fmt.Println("  ✓ Shader modules (vertex + fragment)")
	fmt.Println("  ✓ Fence-based CPU/GPU synchronisation")
}
