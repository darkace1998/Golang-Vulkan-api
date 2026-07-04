# Examples

This directory contains various examples demonstrating how to use the `golang-vulkan-api` library.

## Available Examples

- **`basic`**: Basic Vulkan setup and physical device enumeration.
- **`benchmark`**: A GPU stress testing and benchmarking tool. See `benchmark/README.md` for more details.
- **`compute`**: Demonstrates how to run a compute shader, including buffer creation and memory binding.
- **`descriptor_manager`**: Example demonstrating the usage of the high-level `DescriptorPoolManager` to easily allocate descriptor sets.
- **`descriptor_update`**: Demonstrates how to bind uniform buffers and combined image samplers by updating descriptor sets.
- **`graphics_pipeline`**: A comprehensive graphics pipeline example showing offscreen rendering with vertex buffers, shader modules, render pass, framebuffer, graphics pipeline creation, and draw commands.
- **`multi_queue`**: Shows how to discover, create, and use multiple Vulkan queues, specifically focusing on parallel transfer and graphics operations.
- **`pipeline_cache`**: Example demonstrating pipeline cache creation, retrieval, merging, and loading data.
- **`push_constants`**: Demonstrates how to use push constants to pass small amounts of data to shaders efficiently.
- **`render_to_texture`**: Demonstrates framebuffer creation, render pass, and reading back pixels without a window/surface.
- **`secondary_command_buffer`**: Shows how to record and execute secondary command buffers.
- **`simple`**: Minimal example for Vulkan instance creation.
- **`subpass_dependencies`**: Demonstrates creating a multi-subpass render pass with a dependency chain and self-dependency.
- **`swapchain`**: Demonstrates all swapchain types, constants, input validation, synchronization objects, and the full present loop workflow.
- **`type`**: Type system and constant validation example.
- **`video`**: Demonstrates video codec support detection.
- **`vulkan13`**: Vulkan 1.3 feature demonstration, including dynamic state and rendering info.

## Running an Example

To run an example, use the `go run` command from the root directory or within the specific example directory.

```bash
# Run from the root directory
go run ./examples/basic

# Or cd into the directory
cd examples/compute
go run main.go
```
