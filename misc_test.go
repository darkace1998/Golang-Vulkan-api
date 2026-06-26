package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Pipeline Cache Tests
// ============================================================================

func TestCreatePipelineCacheValidation(t *testing.T) {
	tests := []struct {
		name       string
		device     Device
		createInfo *PipelineCacheCreateInfo
		wantErrMsg string
	}{
		{
			name:       "nil device",
			device:     nil,
			createInfo: &PipelineCacheCreateInfo{},
			wantErrMsg: testNilDeviceError,
		},
		{
			name:       "nil createInfo",
			device:     fakeDevice(),
			createInfo: nil,
			wantErrMsg: "vulkan validation error: createInfo cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreatePipelineCache(tt.device, tt.createInfo)
			if err == nil {
				t.Errorf("CreatePipelineCache() expected error, got nil")
				return
			}

			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Errorf("CreatePipelineCache() expected ValidationError, got %T: %v", err, err)
				return
			}

			if err.Error() != tt.wantErrMsg {
				t.Errorf("CreatePipelineCache() expected error message '%s', got '%s'", tt.wantErrMsg, err.Error())
			}
		})
	}
}

func TestDestroyPipelineCacheNilArgs(t *testing.T) {
	// Should not panic or return error
	DestroyPipelineCache(nil, nil)
	DestroyPipelineCache(nil, fakePipelineCache())
	DestroyPipelineCache(fakeDevice(), nil)
}

func TestGetPipelineCacheDataValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		pipelineCache PipelineCache
		wantErrMsg    string
	}{
		{
			name:          "nil device",
			device:        nil,
			pipelineCache: fakePipelineCache(),
			wantErrMsg:    testNilDeviceError,
		},
		{
			name:          "nil pipelineCache",
			device:        fakeDevice(),
			pipelineCache: nil,
			wantErrMsg:    "vulkan validation error: pipelineCache cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetPipelineCacheData(tt.device, tt.pipelineCache)
			if err == nil {
				t.Errorf("GetPipelineCacheData() expected error, got nil")
				return
			}

			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Errorf("GetPipelineCacheData() expected ValidationError, got %T: %v", err, err)
				return
			}

			if err.Error() != tt.wantErrMsg {
				t.Errorf("GetPipelineCacheData() expected error message '%s', got '%s'", tt.wantErrMsg, err.Error())
			}
		})
	}
}

func TestMergePipelineCachesValidation(t *testing.T) {
	tests := []struct {
		name       string
		device     Device
		dstCache   PipelineCache
		srcCaches  []PipelineCache
		wantErrMsg string
	}{
		{
			name:       "nil device",
			device:     nil,
			dstCache:   fakePipelineCache(),
			srcCaches:  []PipelineCache{fakePipelineCache()},
			wantErrMsg: testNilDeviceError,
		},
		{
			name:       "nil dstCache",
			device:     fakeDevice(),
			dstCache:   nil,
			srcCaches:  []PipelineCache{fakePipelineCache()},
			wantErrMsg: "vulkan validation error: dstCache cannot be nil",
		},
		{
			name:       "empty srcCaches",
			device:     fakeDevice(),
			dstCache:   fakePipelineCache(),
			srcCaches:  []PipelineCache{},
			wantErrMsg: "", // Should return nil error according to implementation
		},
		{
			name:       "nil srcCaches element",
			device:     fakeDevice(),
			dstCache:   fakePipelineCache(),
			srcCaches:  []PipelineCache{nil},
			wantErrMsg: "vulkan validation error: srcCaches contains nil cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MergePipelineCaches(tt.device, tt.dstCache, tt.srcCaches)

			if tt.wantErrMsg == "" {
				if err != nil {
					t.Errorf("MergePipelineCaches() expected nil error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("MergePipelineCaches() expected error, got nil")
				return
			}

			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Errorf("MergePipelineCaches() expected ValidationError, got %T: %v", err, err)
				return
			}

			if err.Error() != tt.wantErrMsg {
				t.Errorf("MergePipelineCaches() expected error message '%s', got '%s'", tt.wantErrMsg, err.Error())
			}
		})
	}
}
