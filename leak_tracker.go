package vulkan

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

// LeakTracker is a utility to track Vulkan resource allocations and detect leaks.
type LeakTracker struct {
	mu        sync.Mutex
	resources map[string]string
	enabled   bool
}

var globalTracker = &LeakTracker{
	resources: make(map[string]string),
	enabled:   false,
}

// EnableLeakTracker turns on tracking of Vulkan object allocations.
func EnableLeakTracker() {
	globalTracker.mu.Lock()
	defer globalTracker.mu.Unlock()
	globalTracker.enabled = true
}

// DisableLeakTracker turns off tracking of Vulkan object allocations.
func DisableLeakTracker() {
	globalTracker.mu.Lock()
	defer globalTracker.mu.Unlock()
	globalTracker.enabled = false
}

// ClearLeaks resets the current list of tracked allocations.
func ClearLeaks() {
	globalTracker.mu.Lock()
	defer globalTracker.mu.Unlock()
	globalTracker.resources = make(map[string]string)
}

func trackResource(resourceType string, handle unsafe.Pointer) {
	if !globalTracker.enabled || handle == nil {
		return
	}

	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	globalTracker.mu.Lock()
	defer globalTracker.mu.Unlock()

	key := fmt.Sprintf("%s:%p", resourceType, handle)
	globalTracker.resources[key] = stack
}

func untrackResource(resourceType string, handle unsafe.Pointer) {
	if !globalTracker.enabled || handle == nil {
		return
	}

	globalTracker.mu.Lock()
	defer globalTracker.mu.Unlock()

	key := fmt.Sprintf("%s:%p", resourceType, handle)
	delete(globalTracker.resources, key)
}

// ReportLeaks returns a formatted string containing information about any un-freed resources.
func ReportLeaks() string {
	globalTracker.mu.Lock()
	defer globalTracker.mu.Unlock()

	if len(globalTracker.resources) == 0 {
		return "No resource leaks detected."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Detected %d potential resource leaks:\n", len(globalTracker.resources)))

	for key, stack := range globalTracker.resources {
		sb.WriteString(fmt.Sprintf("\n--- Leak: %s ---\n", key))
		sb.WriteString(stack)
	}

	return sb.String()
}
