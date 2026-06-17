package vulkan

import (
	"testing"
)

// TestGetLayoutTransitionAccessAndStages verifies that layout transitions map to the correct
// access masks and pipeline stages based on the internal switch logic.
func TestGetLayoutTransitionAccessAndStages(t *testing.T) {
	testCases := []struct {
		name          string
		oldLayout     ImageLayout
		newLayout     ImageLayout
		srcAccessMask AccessFlags
		dstAccessMask AccessFlags
		srcStage      PipelineStageFlags
		dstStage      PipelineStageFlags
	}{
		{
			name:          "Undefined to TransferDst",
			oldLayout:     ImageLayoutUndefined,
			newLayout:     ImageLayoutTransferDstOptimal,
			srcAccessMask: 0,
			dstAccessMask: AccessTransferWriteBit,
			srcStage:      PipelineStageTopOfPipeBit,
			dstStage:      PipelineStageTransferBit,
		},
		{
			name:          "TransferDst to ShaderReadOnly",
			oldLayout:     ImageLayoutTransferDstOptimal,
			newLayout:     ImageLayoutShaderReadOnlyOptimal,
			srcAccessMask: AccessTransferWriteBit,
			dstAccessMask: AccessShaderReadBit,
			srcStage:      PipelineStageTransferBit,
			dstStage:      PipelineStageFragmentShaderBit,
		},
		{
			name:          "Undefined to DepthStencilAttachment",
			oldLayout:     ImageLayoutUndefined,
			newLayout:     ImageLayoutDepthStencilAttachmentOptimal,
			srcAccessMask: 0,
			dstAccessMask: AccessDepthStencilAttachmentReadBit | AccessDepthStencilAttachmentWriteBit,
			srcStage:      PipelineStageTopOfPipeBit,
			dstStage:      PipelineStageEarlyFragmentTestsBit,
		},
		{
			name:          "Undefined to ColorAttachment",
			oldLayout:     ImageLayoutUndefined,
			newLayout:     ImageLayoutColorAttachmentOptimal,
			srcAccessMask: 0,
			dstAccessMask: AccessColorAttachmentReadBit | AccessColorAttachmentWriteBit,
			srcStage:      PipelineStageTopOfPipeBit,
			dstStage:      PipelineStageColorAttachmentOutputBit,
		},
		{
			name:          "ColorAttachment to PresentSrcKHR",
			oldLayout:     ImageLayoutColorAttachmentOptimal,
			newLayout:     ImageLayoutPresentSrcKHR,
			srcAccessMask: AccessColorAttachmentWriteBit,
			dstAccessMask: 0,
			srcStage:      PipelineStageColorAttachmentOutputBit,
			dstStage:      PipelineStageBottomOfPipeBit,
		},
		{
			name:          "TransferSrc to ShaderReadOnly",
			oldLayout:     ImageLayoutTransferSrcOptimal,
			newLayout:     ImageLayoutShaderReadOnlyOptimal,
			srcAccessMask: AccessTransferReadBit,
			dstAccessMask: AccessShaderReadBit,
			srcStage:      PipelineStageTransferBit,
			dstStage:      PipelineStageFragmentShaderBit,
		},
		{
			name:          "ShaderReadOnly to TransferSrc",
			oldLayout:     ImageLayoutShaderReadOnlyOptimal,
			newLayout:     ImageLayoutTransferSrcOptimal,
			srcAccessMask: AccessShaderReadBit,
			dstAccessMask: AccessTransferReadBit,
			srcStage:      PipelineStageFragmentShaderBit,
			dstStage:      PipelineStageTransferBit,
		},
		{
			name:          "ShaderReadOnly to TransferDst",
			oldLayout:     ImageLayoutShaderReadOnlyOptimal,
			newLayout:     ImageLayoutTransferDstOptimal,
			srcAccessMask: AccessShaderReadBit,
			dstAccessMask: AccessTransferWriteBit,
			srcStage:      PipelineStageFragmentShaderBit,
			dstStage:      PipelineStageTransferBit,
		},
		{
			name:          "General to TransferSrc",
			oldLayout:     ImageLayoutGeneral,
			newLayout:     ImageLayoutTransferSrcOptimal,
			srcAccessMask: AccessMemoryReadBit | AccessMemoryWriteBit,
			dstAccessMask: AccessTransferReadBit,
			srcStage:      PipelineStageAllCommandsBit,
			dstStage:      PipelineStageTransferBit,
		},
		{
			name:          "General to TransferDst",
			oldLayout:     ImageLayoutGeneral,
			newLayout:     ImageLayoutTransferDstOptimal,
			srcAccessMask: AccessMemoryReadBit | AccessMemoryWriteBit,
			dstAccessMask: AccessTransferWriteBit,
			srcStage:      PipelineStageAllCommandsBit,
			dstStage:      PipelineStageTransferBit,
		},
		{
			name:          "Fallback (e.g. Preinitialized to General)",
			oldLayout:     ImageLayoutPreinitialized,
			newLayout:     ImageLayoutGeneral,
			srcAccessMask: AccessMemoryReadBit | AccessMemoryWriteBit,
			dstAccessMask: AccessMemoryReadBit | AccessMemoryWriteBit,
			srcStage:      PipelineStageAllCommandsBit,
			dstStage:      PipelineStageAllCommandsBit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srcAccess, dstAccess, srcStage, dstStage := getLayoutTransitionAccessAndStages(tc.oldLayout, tc.newLayout)

			if srcAccess != tc.srcAccessMask {
				t.Errorf("Expected SrcAccessMask %v, got %v", tc.srcAccessMask, srcAccess)
			}
			if dstAccess != tc.dstAccessMask {
				t.Errorf("Expected DstAccessMask %v, got %v", tc.dstAccessMask, dstAccess)
			}
			if srcStage != tc.srcStage {
				t.Errorf("Expected SrcStage %v, got %v", tc.srcStage, srcStage)
			}
			if dstStage != tc.dstStage {
				t.Errorf("Expected DstStage %v, got %v", tc.dstStage, dstStage)
			}
		})
	}
}

// TestTransitionImageLayoutNilArgs ensures TransitionImageLayout gracefully handles nil inputs
func TestTransitionImageLayoutNilArgs(t *testing.T) {
	// Should not panic
	TransitionImageLayout(nil, nil, FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
	TransitionImageLayout(fakeCommandBuffer(), nil, FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
	TransitionImageLayout(nil, fakeImage(), FormatUndefined, ImageLayoutUndefined, ImageLayoutTransferDstOptimal, ImageSubresourceRange{})
}

// TestTransitionImageLayoutSimpleNilArgs ensures TransitionImageLayoutSimple gracefully handles nil inputs
func TestTransitionImageLayoutSimpleNilArgs(t *testing.T) {
	// Should not panic
	TransitionImageLayoutSimple(nil, nil, ImageLayoutUndefined, ImageLayoutTransferDstOptimal)
	TransitionImageLayoutSimple(fakeCommandBuffer(), nil, ImageLayoutUndefined, ImageLayoutTransferDstOptimal)
	TransitionImageLayoutSimple(nil, fakeImage(), ImageLayoutUndefined, ImageLayoutTransferDstOptimal)
}
