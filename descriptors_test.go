package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyImageViewNilArgs tests that DestroyImageView handles nil gracefully
func TestDestroyImageViewNilArgs(t *testing.T) {
	DestroyImageView(nil, nil)
	DestroyImageView(nil, fakeImageView())
	DestroyImageView(fakeDevice(), nil)
}

// TestDestroySamplerNilArgs tests that DestroySampler handles nil gracefully
func TestDestroySamplerNilArgs(t *testing.T) {
	DestroySampler(nil, nil)
	DestroySampler(nil, fakeSampler())
	DestroySampler(fakeDevice(), nil)
}

// TestDestroyDescriptorSetLayoutNilArgs tests that DestroyDescriptorSetLayout handles nil gracefully
func TestDestroyDescriptorSetLayoutNilArgs(t *testing.T) {
	DestroyDescriptorSetLayout(nil, nil)
	DestroyDescriptorSetLayout(nil, fakeDescriptorSetLayout())
	DestroyDescriptorSetLayout(fakeDevice(), nil)
}

// TestDestroyDescriptorPoolNilArgs tests that DestroyDescriptorPool handles nil gracefully
func TestDestroyDescriptorPoolNilArgs(t *testing.T) {
	DestroyDescriptorPool(nil, nil)
	DestroyDescriptorPool(nil, fakeDescriptorPool())
	DestroyDescriptorPool(fakeDevice(), nil)
}

// ============================================================================
// Validation Tests for Create Functions
// ============================================================================

// TestCreateSamplerValidation tests nil parameter validation
func TestCreateSamplerValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *SamplerCreateInfo
		expectParam string
	}{
		{testNilDevice, nil, &SamplerCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateSampler(tt.device, tt.createInfo)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
			}
		})
	}
}

// TestCreateDescriptorSetLayoutValidation tests nil parameter validation
func TestCreateDescriptorSetLayoutValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *DescriptorSetLayoutCreateInfo
		expectParam string
	}{
		{testNilDevice, nil, &DescriptorSetLayoutCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateDescriptorSetLayout(tt.device, tt.createInfo)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
			}
		})
	}
}

// TestCreateDescriptorPoolValidation tests nil parameter validation
func TestCreateDescriptorPoolValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *DescriptorPoolCreateInfo
		expectParam string
	}{
		{testNilDevice, nil, &DescriptorPoolCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateDescriptorPool(tt.device, tt.createInfo)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if valErr.Field != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
			}
		})
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkUpdateDescriptorSets(b *testing.B) {
	device := Device(nil)

	writes := make([]WriteDescriptorSet, 10)
	for i := 0; i < 10; i++ {
		writes[i] = WriteDescriptorSet{
			DstSet:          fakeDescriptorSet(),
			DstBinding:      uint32(i), //nolint:gosec // benchmark safe
			DstArrayElement: 0,
			DescriptorType:  DescriptorTypeUniformBuffer,
			BufferInfo: []DescriptorBufferInfo{
				{
					Buffer: fakeBuffer(),
					Offset: 0,
					Range:  256,
				},
			},
		}
	}

	copies := make([]CopyDescriptorSet, 5)
	for i := 0; i < 5; i++ {
		copies[i] = CopyDescriptorSet{
			SrcSet:          fakeDescriptorSet(),
			SrcBinding:      0,
			SrcArrayElement: 0,
			DstSet:          fakeDescriptorSet(),
			DstBinding:      1,
			DstArrayElement: 0,
			DescriptorCount: 1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UpdateDescriptorSets(device, writes, copies)
	}
}
