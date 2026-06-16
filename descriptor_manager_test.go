package vulkan

import (
	"errors"
	"testing"
)

func TestNewDescriptorPoolManagerValidation(t *testing.T) {
	tests := []struct {
		name           string
		device         Device
		maxSetsPerPool uint32
		poolSizes      []DescriptorPoolSize
		expectParam    string
	}{
		{"nil device", nil, 10, []DescriptorPoolSize{{Type: DescriptorTypeUniformBuffer, DescriptorCount: 10}}, "device"},
		{"zero maxSets", fakeDevice(), 0, []DescriptorPoolSize{{Type: DescriptorTypeUniformBuffer, DescriptorCount: 10}}, "maxSetsPerPool"},
		{"empty poolSizes", fakeDevice(), 10, []DescriptorPoolSize{}, "poolSizes"},
		{"nil poolSizes", fakeDevice(), 10, nil, "poolSizes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDescriptorPoolManager(tt.device, tt.maxSetsPerPool, tt.poolSizes, 0)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("Expected ValidationError, got %T: %v", err, err)
			}
			if valErr.Parameter != tt.expectParam {
				t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Parameter)
			}
		})
	}
}

func TestDescriptorPoolManagerAllocationValidation(t *testing.T) {
	manager, err := NewDescriptorPoolManager(fakeDevice(), 10, []DescriptorPoolSize{
		{Type: DescriptorTypeUniformBuffer, DescriptorCount: 10},
	}, 0)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	sets, err := manager.AllocateDescriptorSets(nil)
	if err != nil {
		t.Fatalf("Expected no error for nil layouts, got %v", err)
	}
	if sets != nil {
		t.Fatalf("Expected nil sets for nil layouts, got %v", sets)
	}

	sets, err = manager.AllocateDescriptorSets([]DescriptorSetLayout{})
	if err != nil {
		t.Fatalf("Expected no error for empty layouts, got %v", err)
	}
	if sets != nil {
		t.Fatalf("Expected nil sets for empty layouts, got %v", sets)
	}
}

func TestDescriptorPoolManagerDestroy(t *testing.T) {
	manager, err := NewDescriptorPoolManager(fakeDevice(), 10, []DescriptorPoolSize{
		{Type: DescriptorTypeUniformBuffer, DescriptorCount: 10},
	}, 0)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Should not panic or error
	manager.Destroy()
}
