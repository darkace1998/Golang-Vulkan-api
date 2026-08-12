package vulkan

import (
	"errors"
	"testing"
)

func TestCreateRayTracingPipelinesKHR_Validation(t *testing.T) {
	t.Run("nil device", func(t *testing.T) {
		_, err := CreateRayTracingPipelinesKHR(nil, nil, []RayTracingPipelineCreateInfoKHR{{}})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var valErr *ValidationError
		if !errors.As(err, &valErr) || valErr.Field != "device" {
			t.Errorf("expected ValidationError for device, got %v", err)
		}
	})

	t.Run("empty create infos", func(t *testing.T) {
		device := fakeDevice()
		_, err := CreateRayTracingPipelinesKHR(device, nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var valErr *ValidationError
		if !errors.As(err, &valErr) || valErr.Field != "createInfos" {
			t.Errorf("expected ValidationError for createInfos, got %v", err)
		}
	})
}
