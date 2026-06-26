import sys

def patch_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    target = "// TestFreeMemoryNilArgs tests that FreeMemory handles nil gracefully"

    new_test = """// TestGetBufferMemoryRequirements tests the GetBufferMemoryRequirements function
func TestGetBufferMemoryRequirements(t *testing.T) {
	t.Run("nil device", func(t *testing.T) {
		reqs := GetBufferMemoryRequirements(nil, fakeBuffer())
		if reqs.Size != 0 || reqs.Alignment != 0 || reqs.MemoryTypeBits != 0 {
			t.Errorf("Expected empty memory requirements for nil device, got: %+v", reqs)
		}
	})

	t.Run("nil buffer", func(t *testing.T) {
		reqs := GetBufferMemoryRequirements(fakeDevice(), nil)
		if reqs.Size != 0 || reqs.Alignment != 0 || reqs.MemoryTypeBits != 0 {
			t.Errorf("Expected empty memory requirements for nil buffer, got: %+v", reqs)
		}
	})

	t.Run("success", func(t *testing.T) {
		// Mock the CGO function
		originalFunc := getBufferMemoryRequirementsFunc
		defer func() { getBufferMemoryRequirementsFunc = originalFunc }()

		expectedReqs := MemoryRequirements{
			Size:           1024,
			Alignment:      256,
			MemoryTypeBits: 7,
		}

		getBufferMemoryRequirementsFunc = func(device Device, buffer Buffer) MemoryRequirements {
			return expectedReqs
		}

		reqs := GetBufferMemoryRequirements(fakeDevice(), fakeBuffer())

		if reqs.Size != expectedReqs.Size {
			t.Errorf("Expected size %d, got %d", expectedReqs.Size, reqs.Size)
		}
		if reqs.Alignment != expectedReqs.Alignment {
			t.Errorf("Expected alignment %d, got %d", expectedReqs.Alignment, reqs.Alignment)
		}
		if reqs.MemoryTypeBits != expectedReqs.MemoryTypeBits {
			t.Errorf("Expected memoryTypeBits %d, got %d", expectedReqs.MemoryTypeBits, reqs.MemoryTypeBits)
		}
	})
}

// TestFreeMemoryNilArgs tests that FreeMemory handles nil gracefully"""

    if target in content:
        content = content.replace(target, new_test)
        with open(filepath, 'w') as f:
            f.write(content)
        print("Patched successfully")
    else:
        print("Target string not found in the file")

if __name__ == "__main__":
    patch_file('memory_test.go')
