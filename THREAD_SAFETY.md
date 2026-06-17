# Thread Safety and Host Synchronization Guarantees

This document outlines the thread safety guarantees and host synchronization requirements for the Golang-Vulkan-API library. These rules are aligned with the official Vulkan specification but also include specific details about this Go binding's behavior.

## General Guarantees

In general, Vulkan handles thread safety internally for most operations, *provided* that you adhere to host synchronization rules. This Go wrapper does not introduce locking mechanisms around core API calls. If the Vulkan spec says a function requires external synchronization (e.g., modifying the same queue from multiple threads without synchronization), you must enforce that synchronization in your Go code using channels, `sync.Mutex`, or other primitives.

*   **Handle Creation/Destruction:** Functions that create or destroy Vulkan objects (e.g., `CreateBuffer`, `DestroyImage`) are generally thread-safe with respect to the `Device` or `Instance` they belong to, as per the Vulkan specification.
*   **Object Modification:** Modifying the *same* Vulkan object concurrently from multiple goroutines (e.g., recording to the same `CommandBuffer` simultaneously) requires explicit external synchronization, unless specified otherwise by Vulkan.

## Function Loaders (Video Extensions)

A critical area of thread safety in this binding concerns the loading of C function pointers for extensions, specifically the video codec extensions (`VK_KHR_video_queue`, etc.).

### Single-Threaded Loading Requirement

*   **NOT Thread-Safe:** Loading video codec functions requires **single-threaded loading**. You **must** call `LoadVideoDeviceFunctions(Device)` and `LoadVideoInstanceFunctions(Instance)` from a single thread before making any video-related API calls.
*   **Execution:** These loaders set global C function pointers internally. If multiple goroutines attempt to load these concurrently, or if one goroutine loads them while another is making a video API call, it will lead to data races and crashes.
*   **Limitation:** Because the function pointers are stored globally, only one `Instance` and one `Device` are fully supported for video operations at a time.

### `ResetVideoInstanceFunctions` and `ResetVideoDeviceFunctions`

*   **NOT Thread-Safe:** These functions are also **NOT thread-safe**.
*   **Usage:** They clear the global function pointers, allowing you to load them again for a *different* `Instance` or `Device`.
*   **Host Synchronization:** You **must** ensure that no goroutine is calling the loaders or executing any video-related API calls while a reset function is being called.

## API Functions Safe for Concurrent Use

Unless the Vulkan specification dictates that a specific parameter requires "externally synchronized" access, the Go API functions can be safely called from multiple goroutines.

Examples of inherently thread-safe operations (assuming different target objects):
*   Allocating different memory chunks (`AllocateMemory`)
*   Recording different command buffers (`CmdBegin`, `CmdDraw`, `CmdEnd`) in parallel.
*   Creating pipelines (`CreateComputePipelines`, `CreateGraphicsPipelines`).

Examples of operations requiring explicit synchronization if sharing the target object:
*   Submitting to the *same* `Queue` (`QueueSubmit`). You must use a mutex if multiple goroutines submit to a single queue concurrently.
*   Updating the *same* Descriptor Set.
*   Recording to the *same* Command Buffer.

## CGO and Go Pointer Rules

When passing data to Vulkan, this binding strictly adheres to Go's CGO pointer passing rules. When using functions that require slices or memory buffers, ensure that the underlying Go memory is not modified concurrently by another goroutine while the Vulkan C function is reading it, as this can lead to undefined behavior or SIGABRT crashes.
