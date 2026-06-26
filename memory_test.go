package vulkan

import (
	"errors"
	"sync"
	"testing"
	"unsafe"
)

// ============================================================================
// MemoryPool Thread-Safety Tests
// ============================================================================

// TestMemoryPoolAllocateBasic tests basic allocation from a memory pool
func TestMemoryPoolAllocateBasic(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024,
		Alignment: 16,
	}

	offset, err := pool.Allocate(64, 16)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}

	offset2, err := pool.Allocate(64, 16)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if offset2 != 64 {
		t.Errorf("Expected offset 64, got %d", offset2)
	}
}

// TestMemoryPoolAllocateAlignment tests alignment handling
func TestMemoryPoolAllocateAlignment(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024,
		Alignment: 4,
	}

	// Allocate 10 bytes with alignment 4
	_, err := pool.Allocate(10, 4)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Next allocation with alignment 256 should align to 256
	offset, err := pool.Allocate(32, 256)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if offset%256 != 0 {
		t.Errorf("Expected offset aligned to 256, got %d", offset)
	}
}

// TestMemoryPoolAllocateExhaustion tests pool exhaustion
func TestMemoryPoolAllocateExhaustion(t *testing.T) {
	pool := &MemoryPool{
		Size:      128,
		Alignment: 1,
	}

	_, err := pool.Allocate(100, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// This should fail - not enough space
	_, err = pool.Allocate(100, 1)
	if err == nil {
		t.Error("Expected error for pool exhaustion, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

// TestMemoryPoolAllocateNilPool tests nil pool handling
func TestMemoryPoolAllocateNilPool(t *testing.T) {
	var pool *MemoryPool
	_, err := pool.Allocate(64, 16)
	if err == nil {
		t.Error("Expected error for nil pool")
	}
}

// TestMemoryPoolAllocateZeroSize tests zero size allocation
func TestMemoryPoolAllocateZeroSize(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024,
		Alignment: 16,
	}

	_, err := pool.Allocate(0, 16)
	if err == nil {
		t.Error("Expected error for zero size")
	}
}

// TestMemoryPoolAllocateInvalidAlignment tests non-power-of-two alignment
func TestMemoryPoolAllocateInvalidAlignment(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024,
		Alignment: 16,
	}

	_, err := pool.Allocate(64, 3) // 3 is not a power of two
	if err == nil {
		t.Error("Expected error for invalid alignment")
	}
}

// TestMemoryPoolReset tests pool reset
func TestMemoryPoolReset(t *testing.T) {
	pool := &MemoryPool{
		Size:      128,
		Alignment: 1,
	}

	_, err := pool.Allocate(100, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	pool.Reset()

	// After reset, should be able to allocate again from the start
	offset, err := pool.Allocate(100, 1)
	if err != nil {
		t.Fatalf("Expected no error after reset, got: %v", err)
	}
	if offset != 0 {
		t.Errorf("Expected offset 0 after reset, got %d", offset)
	}
}

// TestMemoryPoolResetNil tests nil pool reset (should not panic)
func TestMemoryPoolResetNil(t *testing.T) {
	var pool *MemoryPool
	pool.Reset() // should not panic
}

// TestMemoryPoolConcurrentAllocate tests thread-safety of Allocate
func TestMemoryPoolConcurrentAllocate(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024 * 1024, // 1MB
		Alignment: 8,
	}

	const goroutines = 100
	const allocsPerGoroutine = 10
	allocSize := DeviceSize(64)

	var wg sync.WaitGroup
	errCh := make(chan error, goroutines*allocsPerGoroutine)
	offsets := make(chan DeviceSize, goroutines*allocsPerGoroutine)

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < allocsPerGoroutine; i++ {
				offset, err := pool.Allocate(allocSize, 8)
				if err != nil {
					errCh <- err
					return
				}
				offsets <- offset
			}
		}()
	}

	wg.Wait()
	close(errCh)
	close(offsets)

	// Check for errors
	for err := range errCh {
		t.Errorf("Concurrent allocation error: %v", err)
	}

	// Check that all offsets are unique (no data race on Offset)
	seen := make(map[DeviceSize]bool)
	for offset := range offsets {
		if seen[offset] {
			t.Errorf("Duplicate offset %d detected - indicates data race", offset)
		}
		seen[offset] = true
	}
}

// TestMemoryPoolConcurrentAllocateAndReset tests concurrent allocate and reset
func TestMemoryPoolConcurrentAllocateAndReset(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024 * 1024,
		Alignment: 8,
	}

	const iterations = 1000
	var wg sync.WaitGroup

	// Allocator goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if _, err := pool.Allocate(64, 8); err != nil {
				return
			}
		}
	}()

	// Resetter goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations/10; i++ {
			pool.Reset()
		}
	}()

	wg.Wait()
	// If we get here without panic or race detector complaints, the test passes
}

// ============================================================================
// Nil Check Tests for Destroy Functions
// ============================================================================

// TestDestroyBufferNilArgs tests that DestroyBuffer handles nil gracefully
func TestDestroyBufferNilArgs(t *testing.T) {
	// Should not panic
	DestroyBuffer(nil, nil)
	DestroyBuffer(nil, fakeBuffer())
	DestroyBuffer(fakeDevice(), nil)
}

// TestFreeMemoryNilArgs tests that FreeMemory handles nil gracefully
func TestFreeMemoryNilArgs(t *testing.T) {
	FreeMemory(nil, nil)
	FreeMemory(nil, fakeDeviceMemory())
	FreeMemory(fakeDevice(), nil)
}

// TestDestroyImageNilArgs tests that DestroyImage handles nil gracefully
func TestDestroyImageNilArgs(t *testing.T) {
	DestroyImage(nil, nil)
	DestroyImage(nil, fakeImage())
	DestroyImage(fakeDevice(), nil)
}

// ============================================================================
// Validation Tests for Create/Bind Functions
// ============================================================================

// TestAllocateMemoryValidation tests nil parameter validation
func TestAllocateMemoryValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		allocInfo   *MemoryAllocateInfo
		expectParam string
	}{
		{testNilDevice, nil, &MemoryAllocateInfo{}, testDeviceParameter},
		{"nil allocateInfo", fakeDevice(), nil, "allocateInfo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AllocateMemory(tt.device, tt.allocInfo)
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

// TestBindBufferMemoryValidation tests nil parameter validation
func TestBindBufferMemoryValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		buffer      Buffer
		memory      DeviceMemory
		expectParam string
	}{
		{testNilDevice, nil, fakeBuffer(), fakeDeviceMemory(), testDeviceParameter},
		{"nil buffer", fakeDevice(), nil, fakeDeviceMemory(), "buffer"},
		{testNilMemory, fakeDevice(), fakeBuffer(), nil, testMemoryParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindBufferMemory(tt.device, tt.buffer, tt.memory, 0)
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

// TestCreateImageValidation tests nil parameter validation
func TestCreateImageValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		createInfo  *ImageCreateInfo
		expectParam string
	}{
		{testNilDevice, nil, &ImageCreateInfo{}, testDeviceParameter},
		{testNilCreateInfo, fakeDevice(), nil, testCreateInfoParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateImage(tt.device, tt.createInfo)
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

// TestBindImageMemoryValidation tests nil parameter validation
func TestBindImageMemoryValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		image       Image
		memory      DeviceMemory
		expectParam string
	}{
		{testNilDevice, nil, fakeImage(), fakeDeviceMemory(), testDeviceParameter},
		{"nil image", fakeDevice(), nil, fakeDeviceMemory(), "image"},
		{testNilMemory, fakeDevice(), fakeImage(), nil, testMemoryParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindImageMemory(tt.device, tt.image, tt.memory, 0)
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

// TestGetDeviceMemoryCommitmentValidation tests nil parameter validation
func TestGetDeviceMemoryCommitmentValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		memory      DeviceMemory
		expectParam string
	}{
		{testNilDevice, nil, fakeDeviceMemory(), testDeviceParameter},
		{testNilMemory, fakeDevice(), nil, testMemoryParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDeviceMemoryCommitment(tt.device, tt.memory)
			// GetDeviceMemoryCommitment returns 0 on error/nil params,
			// it does not return an error itself to maintain API compatibility
			// with similar query functions, but we can visually verify it handles nil
			// safely without crashing during tests.
		})
	}
}

// TestMapMemoryValidation tests nil parameter validation
func TestMapMemoryValidation(t *testing.T) {
	tests := []struct {
		name        string
		device      Device
		memory      DeviceMemory
		expectParam string
	}{
		{testNilDevice, nil, fakeDeviceMemory(), testDeviceParameter},
		{testNilMemory, fakeDevice(), nil, testMemoryParameter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MapMemory(tt.device, tt.memory, 0, 1024, 0)
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

// TestMapMemoryBoundsValidation tests offset and size validation bounds.
func TestMapMemoryBoundsValidation(t *testing.T) {
	device := fakeDevice()
	memory := fakeDeviceMemory()

	// Add a dummy entry to the tracker directly.
	// Normally AllocateMemory does this, but since we are mocking, we just inject it.
	allocationTrackerMu.Lock()
	allocationTracker[memory] = 1024
	allocationTrackerMu.Unlock()

	defer func() {
		allocationTrackerMu.Lock()
		delete(allocationTracker, memory)
		allocationTrackerMu.Unlock()
	}()

	tests := []struct {
		name        string
		offset      DeviceSize
		size        DeviceSize
		expectError bool
		expectParam string
	}{
		{"Valid Bounds", 0, 1024, false, ""},
		{"Valid Sub-Bounds", 128, 512, false, ""},
		{"Offset Exceeds Size", 1024, 1, true, "offset"},
		{"Offset Out of Bounds", 2048, 1, true, "offset"},
		{"Zero Size", 0, 0, true, "size"},
		{"Whole Size Valid", 0, DeviceSize(WholeSize), false, ""},
		{"Whole Size Valid with Offset", 512, DeviceSize(WholeSize), false, ""},
		{"Size Exceeds Bounds", 0, 2048, true, "size"},
		{"Offset Plus Size Exceeds Bounds", 512, 1024, true, "size"},
		{"Offset Plus Size Overflow", 512, DeviceSize(0xFFFFFFFFFFFFFFF0), true, "size"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock vkMapMemoryFunc
			origMapMemoryFunc := vkMapMemoryFunc
			vkMapMemoryFunc = func(device Device, memory DeviceMemory, offset, size DeviceSize, flags uint32, data *unsafe.Pointer) Result {
				return Success
			}
			defer func() { vkMapMemoryFunc = origMapMemoryFunc }()

			_, err := MapMemory(device, memory, tt.offset, tt.size, 0)
			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				var valErr *ValidationError
				if errors.As(err, &valErr) {
					if valErr.Field != tt.expectParam {
						t.Errorf("Expected error param '%s', got '%s'", tt.expectParam, valErr.Field)
					}
				} else {
					t.Errorf("Expected ValidationError, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestMemoryPoolDefaultAlignment tests zero alignment uses pool default
func TestMemoryPoolDefaultAlignment(t *testing.T) {
	pool := &MemoryPool{
		Size:      1024,
		Alignment: 64,
	}

	// First alloc: 10 bytes
	_, err := pool.Allocate(10, 0) // 0 should use pool.Alignment=64
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Second alloc should be aligned to 64
	offset, err := pool.Allocate(10, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if offset%64 != 0 {
		t.Errorf("Expected offset aligned to 64, got %d", offset)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkMemoryPoolAllocate(b *testing.B) {
	pool := &MemoryPool{
		Size:      DeviceSize(uint64(b.N)) * 64, //nolint:gosec // benchmark safe // Ensure enough size
		Alignment: 64,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pool.Allocate(64, 64)
	}
}

func TestFindMemoryType(t *testing.T) {
	memProperties := PhysicalDeviceMemoryProperties{
		MemoryTypeCount: 4,
	}
	memProperties.MemoryTypes[0] = MemoryType{PropertyFlags: MemoryPropertyDeviceLocalBit}
	memProperties.MemoryTypes[1] = MemoryType{PropertyFlags: MemoryPropertyHostVisibleBit | MemoryPropertyHostCoherentBit}
	memProperties.MemoryTypes[2] = MemoryType{PropertyFlags: MemoryPropertyHostVisibleBit | MemoryPropertyHostCachedBit}
	memProperties.MemoryTypes[3] = MemoryType{PropertyFlags: MemoryPropertyDeviceLocalBit | MemoryPropertyHostVisibleBit}

	tests := []struct {
		name          string
		typeFilter    uint32
		properties    MemoryPropertyFlags
		expectedIndex uint32
		expectedFound bool
	}{
		{
			name:          "exact match single property",
			typeFilter:    0b1111,
			properties:    MemoryPropertyDeviceLocalBit,
			expectedIndex: 0,
			expectedFound: true,
		},
		{
			name:          "exact match multiple properties",
			typeFilter:    0b1111,
			properties:    MemoryPropertyHostVisibleBit | MemoryPropertyHostCoherentBit,
			expectedIndex: 1,
			expectedFound: true,
		},
		{
			name:          "type filter excludes match",
			typeFilter:    0b0101, // Only allows index 0 and 2
			properties:    MemoryPropertyHostVisibleBit | MemoryPropertyHostCoherentBit, // This is at index 1
			expectedIndex: 0,
			expectedFound: false,
		},
		{
			name:          "property mismatch",
			typeFilter:    0b1111,
			properties:    MemoryPropertyProtectedBit,
			expectedIndex: 0,
			expectedFound: false,
		},
		{
			name:          "multiple valid matches returns lowest index",
			typeFilter:    0b1111,
			properties:    MemoryPropertyHostVisibleBit, // Index 1, 2, and 3 have this, returns 1
			expectedIndex: 1,
			expectedFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, found := FindMemoryType(memProperties, tt.typeFilter, tt.properties)
			if found != tt.expectedFound {
				t.Errorf("expected found %v, got %v", tt.expectedFound, found)
			}
			if found && idx != tt.expectedIndex {
				t.Errorf("expected index %d, got %d", tt.expectedIndex, idx)
			}
		})
	}
}
