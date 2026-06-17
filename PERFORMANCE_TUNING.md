# Performance Tuning Guide

This guide provides strategies for optimizing Vulkan compute workloads using `golang-vulkan-api`, with a specific focus on AI/ML and general parallel processing tasks.

## AI/ML Compute Workloads

When implementing AI/ML operations (like matrix multiplication, convolution, or tensor transformations) via Vulkan compute shaders, consider the following optimization strategies.

### 1. Batch Size and Workgroup Dimensions

The choice of workgroup dimensions is critical to achieving high occupancy on the GPU.

- **Match Hardware Waves**: Adjust your workgroup dimensions to be a multiple of the underlying hardware's wavefront or warp size (typically 32 for NVIDIA, 64 for AMD). This minimizes idle threads within a workgroup.
- **Batch Processing**: Instead of dispatching single operations, try to batch multiple AI operations together in a single `CmdDispatch` call. This reduces host-to-device API overhead.
- **Tiling**: For matrix operations, implement shared memory tiling in your compute shaders to reduce global memory accesses.
- **Avoid Underutilization**: If your dataset is small, a single dispatch may not provide enough workgroups to saturate the GPU's compute units. In such cases, executing multiple parallel passes or grouping small tasks is more efficient.

### 2. Memory Alignment and Layout

Memory access patterns heavily dictate the performance of compute shaders.

- **Optimal Layouts**: Store tensors and weights in a linear layout that matches how your shader will read them. Prefer Struct of Arrays (SoA) or Array of Structs (AoS) based on spatial locality in the shader.
- **Buffer Alignment**: Comply with `minStorageBufferOffsetAlignment` when allocating sub-regions of a larger Vulkan buffer to ensure that the memory offset fits hardware constraints.
- **Memory Pooling**: Do not allocate and free Vulkan `DeviceMemory` per operation. Use a memory allocator to sub-allocate large memory blocks for intermediate tensor data.
- **Host-Visible Memory**: Use `VK_MEMORY_PROPERTY_HOST_VISIBLE_BIT | VK_MEMORY_PROPERTY_HOST_COHERENT_BIT` for buffers that the CPU needs to read frequently, but consider utilizing an intermediate device-local staging buffer if the GPU needs to read the buffer multiple times.

### 3. Queue Family Selection

Using the right queues can enhance concurrency and avoid stalling rendering or presentation work.

- **Dedicated Compute Queues**: When creating a logical device, check the properties of available queue families using `GetPhysicalDeviceQueueFamilyProperties`. Select a queue family that supports `VK_QUEUE_COMPUTE_BIT` but ideally *not* `VK_QUEUE_GRAPHICS_BIT`. Dedicated asynchronous compute queues often have fewer synchronization overheads with graphics tasks.
- **Concurrent Execution**: Submit compute workloads to a separate queue from your graphics/transfer queues if you are interleaving rendering and AI tasks. This allows the GPU's command processor to overlap workloads.

### 4. Timeline Semaphore Usage (Vulkan 1.2+)

Timeline semaphores (`VK_KHR_timeline_semaphore`) offer a much more flexible synchronization primitive than binary semaphores, particularly suited for chained ML operations.

- **Monotonic Progress**: Use a single timeline semaphore to track the progress of a sequence of compute dispatches. The CPU can easily query the current value to know exactly which stages of the pipeline have completed.
- **Wait-Before-Signal**: You can submit a compute workload that waits on a timeline value that *has not yet been signaled* by the host or another queue. This allows you to build an entire command graph upfront and submit it, reducing CPU bottlenecks.
- **Host Integration**: Use `WaitSemaphores` on the CPU to block until the GPU finishes a specific batch of ML processing, eliminating the need for coarse-grained fences (`vkWaitForFences`) per task.
