package vulkan

import (
	"testing"
)

// TestRenderingInfoStruct tests the RenderingInfo struct defaults
func TestRenderingInfoStruct(t *testing.T) {
	info := RenderingInfo{
		Flags:      RenderingContentsSecondaryCommandBuffers,
		LayerCount: 1,
	}

	if info.Flags != RenderingContentsSecondaryCommandBuffers {
		t.Errorf("Expected Flags %v, got %v", RenderingContentsSecondaryCommandBuffers, info.Flags)
	}
	if info.LayerCount != 1 {
		t.Errorf("Expected LayerCount 1, got %d", info.LayerCount)
	}
	if info.ViewMask != 0 {
		t.Errorf("Expected ViewMask 0, got %d", info.ViewMask)
	}
}

// TestSubmitInfo2Struct tests the SubmitInfo2 struct defaults
func TestSubmitInfo2Struct(t *testing.T) {
	info := SubmitInfo2{
		Flags: SubmitProtected,
	}

	if info.Flags != SubmitProtected {
		t.Errorf("Expected Flags %v, got %v", SubmitProtected, info.Flags)
	}
	if len(info.CommandBufferInfos) != 0 {
		t.Errorf("Expected 0 CommandBufferInfos, got %d", len(info.CommandBufferInfos))
	}
}

// TestExtendedDynamicState tests the extended dynamic state constants and functions
// Since the functions wrap C code, we will just test that the constants are defined correctly
func TestExtendedDynamicState(t *testing.T) {
	if CullModeFrontAndBack != CullModeFront|CullModeBack {
		t.Errorf("Expected CullModeFrontAndBack to be combination of front and back")
	}
	if FrontFaceCounterClockwise != 0 {
		t.Errorf("Expected FrontFaceCounterClockwise to be 0")
	}
	if PrimitiveTopologyTriangleList != 3 {
		t.Errorf("Expected PrimitiveTopologyTriangleList to be 3 (VK_PRIMITIVE_TOPOLOGY_TRIANGLE_LIST)")
	}
}
