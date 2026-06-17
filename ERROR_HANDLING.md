# Error Handling in Vulkan Go API

This document provides idiomatic Go patterns for handling errors when working with the `golang-vulkan-api` package. It covers differentiating between validation errors and Vulkan execution errors, as well as strategies for recovering from common failure scenarios.

## Error Types

The library primarily uses two structured error types:

1.  **`ValidationError`**: Returned when the API detects invalid input *before* calling the underlying Vulkan C API. This includes checking for `nil` pointers, invalid array lengths, or invalid enum values.
2.  **`VulkanError`**: Returned when a Vulkan C API call fails. It wraps the raw Vulkan `Result` code and provides context about the operation that failed.

### Handling Errors

You should generally check for errors after every Vulkan API call. The `IsVulkanError` and standard `errors.As` or `errors.Is` functions can be used to inspect the error.

```go
package main

import (
	"errors"
	"fmt"
	"log"

	vulkan "github.com/darkace1998/golang-vulkan-api"
)

func exampleErrorHandling(device vulkan.Device) {
	// Attempt to create a buffer with potentially invalid parameters
	buffer, err := vulkan.CreateBuffer(device, &vulkan.BufferCreateInfo{
		// Intentional error: size 0 is often invalid depending on usage
		Size: 0,
		Usage: vulkan.BufferUsageUniformBufferBit,
	})

	if err != nil {
		var valErr *vulkan.ValidationError
		var vulkanErr *vulkan.VulkanError

		if errors.As(err, &valErr) {
			// Handle validation errors (typically programming mistakes)
			log.Fatalf("Validation error: %s - %s", valErr.Field, valErr.Reason)
		} else if errors.As(err, &vulkanErr) {
			// Handle Vulkan execution errors
			fmt.Printf("Vulkan operation '%s' failed with result: %s\n", vulkanErr.Operation, vulkanErr.Result.Error())

            // Further inspection of specific Vulkan errors
            if vulkanErr.Result == vulkan.ErrorOutOfDeviceMemory || vulkanErr.Result == vulkan.ErrorOutOfHostMemory {
                log.Println("Fatal: Out of memory!")
                // Trigger memory cleanup or application shutdown
            }
		} else {
			// Handle unexpected errors
			log.Fatalf("Unknown error: %v", err)
		}
		return
	}
	defer vulkan.DestroyBuffer(device, buffer, nil)
}
```

## Recovering from Transient Errors

Certain Vulkan errors indicate transient failures where the operation might succeed if retried later, or where the application can recover by recreating resources.

### Handling `VK_ERROR_DEVICE_LOST`

`VK_ERROR_DEVICE_LOST` (represented as `vulkan.ErrorDeviceLost` in Go) is one of the most severe but potentially recoverable errors. It indicates that the logical or physical device has been lost (e.g., due to a driver crash, hardware reset, or GPU hanging).

**Recovery Strategy:**

1.  **Detect the Error:** Any command that interacts with the device (queue submission, waiting on fences, memory mapping) might return this error.
2.  **Clean Up Current State:** You cannot use the current `vulkan.Device` object anymore. You must destroy all resources associated with it (buffers, images, pipelines, swapchains) and finally destroy the device itself. *Note: Destroying resources on a lost device is generally safe and expected in Vulkan.*
3.  **Re-enumerate and Recreate:** You must re-enumerate physical devices (the physical device handle might still be valid, but querying properties might fail or return different results, or you might need to find a new physical device entirely).
4.  **Recreate Everything:** Create a new logical device and recreate all necessary resources.

```go
func submitWork(queue vulkan.Queue, submitInfos []vulkan.SubmitInfo, fence vulkan.Fence) error {
	err := vulkan.QueueSubmit(queue, submitInfos, fence)
	if err != nil {
		var vErr *vulkan.VulkanError
		if errors.As(err, &vErr) && vErr.Result == vulkan.ErrorDeviceLost {
			log.Println("WARNING: GPU Device Lost detected during QueueSubmit!")
			// Trigger a recovery routine in your application loop
			// return errDeviceLostToAppLoop // Return a specific signal
		}
		return err
	}
	return nil
}

// Conceptual Recovery Loop in your main application:
func runApplication() {
    for {
        // Initialize Vulkan, Device, Swapchain, Pipelines...
        // appState := initVulkan()

        // Main rendering loop
        // err := appState.renderLoop()

        // if err == errDeviceLostToAppLoop {
        //     log.Println("Attempting Device Lost recovery...")
        //     appState.cleanupAllResources()
        //     continue // Restart initialization
        // } else if err != nil {
        //     log.Fatalf("Fatal error: %v", err)
        // }
        // break // Normal exit
    }
}
```

### Handling `VK_ERROR_OUT_OF_DATE_KHR` and `VK_SUBOPTIMAL_KHR`

These errors are common when working with swapchains (e.g., when resizing the window).

*   `vulkan.ErrorOutOfDateKhr`: The swapchain is no longer compatible with the surface (e.g., the window size changed). You *must* recreate the swapchain before rendering again.
*   `vulkan.SuboptimalKhr`: The swapchain can still be used, but the surface properties no longer match exactly. You *should* recreate the swapchain for optimal performance, but it's not strictly required immediately.

These are typically returned by `AcquireNextImageKHR` or `QueuePresentKHR`.

```go
// Example pseudo-code for swapchain handling
func drawFrame(...) {
    imageIndex, err := vulkan.AcquireNextImageKHR(device, swapchain, timeout, imageAvailableSemaphore, nil)

    if err != nil {
        var vErr *vulkan.VulkanError
        if errors.As(err, &vErr) {
            if vErr.Result == vulkan.ErrorOutOfDateKhr {
                // Must recreate swapchain
                // recreateSwapchain()
                return // Skip drawing this frame
            } else if vErr.Result == vulkan.SuboptimalKhr {
                // Can still render, but consider recreating soon
                // markedForSwapchainRecreation = true
            } else {
                log.Fatalf("Failed to acquire image: %v", err)
            }
        }
    }

    // ... proceed with recording and submitting command buffers ...

    presentErr := vulkan.QueuePresentKHR(presentQueue, &vulkan.PresentInfoKHR{...})
    if presentErr != nil {
         var vErr *vulkan.VulkanError
         if errors.As(presentErr, &vErr) && (vErr.Result == vulkan.ErrorOutOfDateKhr || vErr.Result == vulkan.SuboptimalKhr) {
             // recreateSwapchain()
         }
    }
}
```
