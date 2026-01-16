# TODO - Golang-Vulkan-api

This document outlines areas for improvement, planned features, and known issues in the Go Vulkan API binding.

## Table of Contents

- [High Priority](#high-priority)
- [Medium Priority](#medium-priority)
- [Low Priority](#low-priority)
- [Testing](#testing)
- [Documentation](#documentation)
- [Performance](#performance)
- [Platform Support](#platform-support)

---

## High Priority

### API Completeness

- [ ] **Graphics Pipeline Creation**: Complete the graphics pipeline API with all vertex input state options
- [ ] **Swapchain Management**: Implement complete swapchain creation, image acquisition, and presentation
- [ ] **Buffer Views**: Add `CreateBufferView` and `DestroyBufferView` functions
- [x] ~~**Descriptor Pool Management**: Add `ResetDescriptorPool` function~~ *(Implemented in descriptors.go)*
- [x] ~~**Descriptor Set Allocation**: Implement `AllocateDescriptorSets` and `FreeDescriptorSets`~~ *(Implemented in descriptors.go)*
- [x] ~~**Descriptor Set Updates**: Add `UpdateDescriptorSets` and `WriteDescriptorSet` support~~ *(Implemented in descriptors.go)*

### Video Codec Enhancements

- [ ] **Reference Slots Implementation**: Complete reference slot handling in `CmdDecodeVideo` and `CmdEncodeVideo`
- [ ] **VP9 Codec Support**: Add VP9 decode/encode when Vulkan extensions become available
- [ ] **Video Encode Quality Control**: Add rate control and quality parameter configuration

### Memory Management

- [ ] **Dedicated Allocation**: Implement `VK_KHR_dedicated_allocation` support
- [ ] **External Memory**: Add support for `VK_KHR_external_memory` extensions
- [ ] **Memory Budget Tracking**: Add utilities to track and manage GPU memory budget

---

## Medium Priority

### Missing Core Features

- [x] ~~**Events**: Implement `CreateEvent`, `DestroyEvent`, `SetEvent`, `ResetEvent`, `GetEventStatus`~~ *(Implemented in synchronization.go)*
- [x] ~~**Query Pools**: Complete query pool implementation with `CreateQueryPool`, `DestroyQueryPool`, `BeginQuery`, `EndQuery`, `GetQueryPoolResults`~~ *(Implemented in queries.go)*
- [ ] **Pipeline Cache**: Add `CreatePipelineCache`, `DestroyPipelineCache`, `GetPipelineCacheData`, `MergePipelineCaches`
- [x] ~~**Secondary Command Buffers**: Improve secondary command buffer support with `CmdExecuteCommands`~~ *(Implemented in commands.go)*
- [x] ~~**Push Constants**: Add `CmdPushConstants` function~~ *(Implemented in commands.go)*

### Vulkan 1.2 Features

- [x] ~~**Timeline Semaphores**: Expose timeline semaphore creation with `VK_SEMAPHORE_TYPE_TIMELINE`~~ *(Implemented in synchronization.go)*
- [ ] **Buffer Device Address**: Add helpers for buffer device address usage
- [ ] **Descriptor Indexing**: Support descriptor binding flags for non-uniform indexing

### Vulkan 1.3 Features (Remaining)

- [ ] **Copy Commands 2**: Implement `CmdCopyBuffer2`, `CmdCopyImage2`, `CmdBlitImage2`, `CmdResolveImage2`
- [ ] **Format Feature Flags 2**: Add extended format feature query functions
- [ ] **Shader Module Create Flags**: Support shader module create flags for maintenance

### Compute Shader Support

- [ ] **Compute Dispatch Base**: Add `CmdDispatchBase` for multi-device dispatch
- [ ] **Dispatch Indirect Count**: Implement count-based indirect dispatch variants

---

## Low Priority

### Extensions

- [ ] **Ray Tracing**: Add `VK_KHR_ray_tracing_pipeline` and `VK_KHR_acceleration_structure` support
- [ ] **Mesh Shaders**: Implement `VK_EXT_mesh_shader` extension bindings
- [ ] **Fragment Shading Rate**: Add `VK_KHR_fragment_shading_rate` support
- [ ] **Debug Utils**: Complete debug messenger implementation with callback support
- [ ] **Validation Features**: Add `VK_EXT_validation_features` for enhanced debugging

### Platform-Specific

- [ ] **Wayland Surface**: Add Wayland surface creation for Linux
- [ ] **XCB/Xlib Surface**: Add X11 surface creation support
- [ ] **Win32 Surface**: Complete Windows surface creation
- [ ] **Metal/MoltenVK**: Test and validate macOS support via MoltenVK

### Advanced Features

- [ ] **Sparse Resources**: Implement sparse buffer and image binding
- [ ] **Multi-View Rendering**: Add multi-view/multi-GPU rendering support
- [ ] **Variable Rate Shading**: Implement `VK_KHR_fragment_shading_rate` variable rate support

---

## Testing

### Unit Tests

- [ ] **Instance Creation Tests**: Add comprehensive instance creation/destruction tests
- [ ] **Device Selection Tests**: Test physical device selection and feature querying
- [ ] **Memory Allocation Tests**: Validate memory type selection and allocation
- [ ] **Command Buffer Tests**: Test command buffer recording and submission
- [ ] **Synchronization Tests**: Verify fence and semaphore behavior

### Integration Tests

- [ ] **Full Render Pipeline**: Create end-to-end render pipeline test
- [ ] **Compute Shader Tests**: Test complete compute shader workflow
- [ ] **Video Encode/Decode Tests**: Integration tests for video codec operations (hardware dependent)

### Test Infrastructure

- [ ] **Mock Vulkan Loader**: Create mock for testing without GPU
- [x] ~~**CI/CD Integration**: Add GitHub Actions workflow for automated testing~~ *(golangci-lint workflow exists)*
- [ ] **Build and Test Workflow**: Add GitHub Actions workflow for building and running tests
- [ ] **Coverage Reporting**: Implement test coverage tracking

---

## Documentation

### API Documentation

- [ ] **GoDoc Comments**: Add comprehensive godoc comments to all exported functions
- [ ] **Parameter Validation Docs**: Document expected parameter ranges and constraints
- [ ] **Error Code Reference**: Create reference for all possible error return values
- [ ] **Thread Safety Documentation**: Document which functions are thread-safe

### Guides and Tutorials

- [ ] **Getting Started Guide**: Write beginner-friendly setup and usage guide
- [ ] **Compute Shader Tutorial**: Step-by-step compute shader example
- [ ] **Video Processing Guide**: Documentation for video encode/decode workflows
- [ ] **Memory Management Best Practices**: Guide for efficient GPU memory usage
- [ ] **Multi-Threading Guide**: Document parallel command buffer recording patterns

### Examples

- [ ] **Triangle Example**: Complete minimal triangle rendering example
- [ ] **Texture Mapping**: Add textured quad example with sampler setup
- [ ] **Deferred Rendering**: Example of multiple render targets
- [ ] **Post-Processing**: Screen-space effects example
- [ ] **Video Transcoding**: Complete video encode/decode example

---

## Performance

### Optimization Opportunities

- [ ] **CGO Overhead Reduction**: Batch Vulkan calls where possible to reduce CGO overhead
- [ ] **Memory Pool Improvements**: Add free-list or buddy allocator for memory pool
- [ ] **Command Buffer Caching**: Implement command buffer reuse patterns
- [ ] **Pipeline State Object Caching**: Add PSO cache for pipeline recreation

### Profiling

- [ ] **GPU Timestamp Queries**: Expose timestamp query functionality
- [ ] **Pipeline Statistics**: Add pipeline statistics query support
- [ ] **Memory Statistics**: Implement memory heap usage tracking

### Benchmarks

- [ ] **API Call Benchmarks**: Benchmark CGO overhead for common operations
- [ ] **Allocation Benchmarks**: Compare memory allocation strategies
- [ ] **Command Buffer Benchmarks**: Measure command buffer recording performance

---

## Platform Support

### Linux

- [x] Basic Vulkan support via pkg-config
- [ ] XCB surface creation
- [ ] Xlib surface creation
- [ ] Wayland surface creation

### Windows

- [x] Basic Vulkan support
- [ ] Win32 surface creation
- [ ] DXGI interop

### macOS

- [x] MoltenVK support (via Vulkan SDK)
- [ ] Metal surface creation
- [ ] macOS-specific optimizations

### Other

- [ ] Android NDK support
- [ ] FreeBSD testing
- [ ] ARM architecture testing (Raspberry Pi, Apple Silicon)

---

## Code Quality

### Refactoring

- [ ] **Error Handling Consistency**: Ensure all functions use consistent error patterns
- [ ] **Resource Cleanup**: Audit all resource creation for proper cleanup paths
- [ ] **Nil Safety**: Add comprehensive nil checks for all handle parameters
- [ ] **Const Correctness**: Review and fix constant definitions

### Static Analysis

- [ ] **Address `gosec` Warnings**: Review and address any security warnings
- [ ] **Fix `staticcheck` Issues**: Resolve all staticcheck warnings
- [ ] **Reduce Cyclomatic Complexity**: Simplify complex functions

### Code Organization

- [ ] **Separate Extension Code**: Move extension-specific code to separate files
- [ ] **Group Related Types**: Organize type definitions by Vulkan object category
- [ ] **Reduce File Sizes**: Split large files into logical components

---

## Known Issues

### Current Limitations

1. **Single Device Support for Video**: Video extension function pointers only support one device at a time
2. **Reference Slot Handling**: Video reference slots in decode/encode are not yet implemented
3. **Graphics Pipeline State**: Some advanced graphics pipeline states are not fully exposed
4. **Surface Creation**: No windowing system integration (requires external window library)

### Workarounds

1. For windowing, use with GLFW, SDL2, or other cross-platform window libraries
2. For advanced features not yet implemented, use the raw CGO bridge

---

## Contributing

When contributing to this project, please:

1. Follow existing code style and patterns
2. Add tests for new functionality
3. Update documentation for API changes
4. Test on multiple platforms when possible
5. Reference this TODO when picking up tasks

---

*Last updated: 2026-01-16*
