package vulkan

import (
	"strings"
	"testing"
	"unsafe"
)

func TestLeakTracker(t *testing.T) {
	EnableLeakTracker()
	defer DisableLeakTracker()
	ClearLeaks()

	var dummy1, dummy2 int
	handle1 := unsafe.Pointer(&dummy1)
	handle2 := unsafe.Pointer(&dummy2)

	trackResource("TestResource", handle1)
	trackResource("TestResource", handle2)

	untrackResource("TestResource", handle1)

	report := ReportLeaks()

	if !strings.Contains(report, "Detected 1 potential resource leaks") {
		t.Errorf("Expected 1 leak detected, got: %s", report)
	}
	if !strings.Contains(report, "TestResource") {
		t.Errorf("Expected report to contain resource name, got: %s", report)
	}

	ClearLeaks()
	report = ReportLeaks()
	if report != "No resource leaks detected." {
		t.Errorf("Expected no leaks, got: %s", report)
	}
}
