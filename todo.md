# TODO — Golang-Vulkan-api

> Codebase analysis and improvement roadmap for the Go Vulkan 1.3+ binding library.
> Verified on the current `origin/main` checkout with `libvulkan-dev` installed; `go build ./...`, `go test ./...`, and `go test -race ./...` all pass.

---

## Table of Contents

- [High Priority](#high-priority)
- [Current Verified State](#current-verified-state)
- [Testing](#testing)
- [Documentation](#documentation)
- [Examples](#examples)
- [Build & CI/CD](#build--cicd)
- [API Coverage](#api-coverage)
- [Performance](#performance)
- [Security & Robustness](#security--robustness)
- [Code Quality](#code-quality)
- [Long-Term / Future](#long-term--future)

---

## Current Verified State

- [x] `libvulkan-dev` is installed
- [x] `go build ./...` passes
- [x] `go test ./...` passes
- [x] `go test -race ./...` passes
- [x] Repository is synced to `origin/main`

---

## High Priority

- [x] **Implement benchmark CSV export** — The backup benchmark (`examples/benchmark_backup/`) now uses the `-csv` flag to call `exportToCSV`, matching the main benchmark.
- [x] **Complete swapchain usage example** — `examples/swapchain/main.go` demonstrates all swapchain types, constants, input validation, synchronization objects, and the full present loop workflow.
- [x] **Add a full graphics pipeline example** — `examples/graphics_pipeline/main.go` demonstrates offscreen rendering with vertex buffers, shader modules, render pass, framebuffer, graphics pipeline creation, and draw commands.

---

## Testing

- [ ] **Expand unit test coverage beyond validation** — Current tests (`instance_test.go`, `validation_test.go`, `video_test.go`) focus on input validation and error types; add tests for struct field defaults, constant correctness, and helper function edge cases.
- [ ] **Add tests for `resources.go`** — `TransitionImageLayout()` has non-trivial logic for selecting access masks and pipeline stages based on layout transitions; unit test all layout transition paths.
- [ ] **Add tests for `synchronization.go`** — Timeline semaphore creation helpers and fence status logic should be tested in pure Go (validation-only, no GPU required).
- [x] **Add tests for `vulkan13.go`** — `RenderingInfo`, `SubmitInfo2`, and extended dynamic state types need struct and validation tests.
- [x] **Add tests for `video_helpers.go`** — `IsExtensionSupported()` and `IsLayerSupported()` are pure Go functions that can be easily unit tested.
- [ ] **Add integration tests with a mock or software Vulkan driver** — Consider using `lavapipe` (Mesa software renderer) in CI to run tests that require a Vulkan instance without a physical GPU.
- [x] **Increase benchmark coverage** — Add benchmarks for memory allocation helpers, command buffer recording, and descriptor set updates.

---

## Documentation

- [ ] **Add a "Getting Started" tutorial** — A step-by-step guide for new users that walks through instance creation → device selection → memory allocation → command buffer → first draw/compute dispatch.
- [x] **Add inline doc comments on all exported types and functions** — Some exported symbols (e.g., in `queries.go`, `swapchain.go`, `misc.go`) lack Go doc comments.
- [x] **Document thread safety guarantees** — Video codec functions require single-threaded loading (`LoadVideoDeviceFunctions`); document which API functions are safe for concurrent use.
- [ ] **Document error recovery patterns** — Show idiomatic Go patterns for handling `VulkanError` vs. `ValidationError`, including retry logic for transient failures like `VK_ERROR_DEVICE_LOST`.
- [ ] **Add a performance tuning guide** — Tips for AI/ML compute workloads: batch size, memory alignment, queue family selection, timeline semaphore usage.
- [ ] **Document Vulkan 1.4 readiness** — `types.go` defines `Version14`; document which 1.4 features are ready and which are planned.

---

## Examples

- [x] **Add a render-to-texture offscreen example** — Demonstrate framebuffer creation, render pass, and reading back pixels without a window/surface.
- [x] **Add a multi-queue example** — Show transfer + graphics queue usage for async resource uploads.
- [ ] **Add a descriptor set update example** — Demonstrate uniform buffer and combined image sampler descriptor binding.
- [ ] **Add a push constants example** — Push constants are defined in the pipeline layout types but not demonstrated.
- [ ] **Expand the `vulkan13/` example** — Show private data slots and maintenance4 (`GetDeviceBufferMemoryRequirements` / `GetDeviceImageMemoryRequirements`) in action.
- [x] **Clean up `benchmark_backup/`** — Either merge useful parts into the main `benchmark/` example or remove the backup directory to reduce duplication.

---

## Build & CI/CD

- [x] **Add CI test job** — The current GitHub Actions workflow (`golangci-lint.yml`) only lints; add a job that runs `go test ./...` on Ubuntu.
- [ ] **Add multi-platform CI** — Test builds on Windows and macOS runners in addition to Ubuntu to catch platform-specific CGO issues early.
- [ ] **Add a CI job with `lavapipe`** — Install Mesa's software Vulkan driver in CI to run integration tests that need a Vulkan instance.
- [ ] **Pin golangci-lint version in CI** — Currently uses `latest`; pin to a specific version for reproducible builds.
- [ ] **Add a build matrix for Go versions** — Test against Go 1.22 and the latest stable Go release.
- [ ] **Add `go vet` to the CI pipeline** — Runs additional static analysis beyond `golangci-lint`.

---

## API Coverage

- [ ] **Add sparse memory binding support** — `SparseMemoryBind` and related types are partially defined; complete the binding functions (`QueueBindSparse`, sparse image/buffer support).
- [ ] **Add pipeline cache support** — `VkPipelineCache` creation, retrieval, and merging functions are not yet wrapped.
- [x] **Add secondary command buffer support** — `CommandBufferInheritanceInfo` is defined but secondary command buffer recording and execution (`CmdExecuteCommands`) need examples and tests.
- [ ] **Add subpass dependency support** — `SubpassDependency` is defined but subpass dependency chains and self-dependencies are not demonstrated.
- [ ] **Complete surface/WSI integration** — Surface creation functions (`vkCreateXlibSurfaceKHR`, `vkCreateWin32SurfaceKHR`, `vkCreateMetalSurfaceKHR`) are not yet wrapped; these are needed for on-screen rendering.
- [ ] **Add ray tracing extension support** — `VK_KHR_ray_tracing_pipeline` and acceleration structure extensions for modern GPU ray tracing.
- [ ] **Add `VK_EXT_mesh_shader` support** — Mesh and task shader pipeline stages for advanced geometry processing.

---

## Performance

- [ ] **Profile CGO overhead** — Measure and document the cost of CGO calls per Vulkan function; identify hot paths where batching or caching would help.
- [ ] **Add object handle caching** — Frequently retrieved handles (queue, device properties) could be cached on the Go side to avoid repeated CGO roundtrips.
- [x] **Pool allocator for descriptor sets** — Provide a higher-level descriptor set pool manager that reduces allocation overhead.

---

## Security & Robustness

- [ ] **Audit CGO pointer passing** — Review all `unsafe.Pointer` usage for correctness under Go's CGO pointer passing rules; ensure no Go pointers are stored on the C side.
- [ ] **Add resource leak detection** — Provide a debug mode or utility that tracks `Create*` / `Destroy*` call pairs and reports leaks.
- [ ] **Validate memory mapping bounds** — `MapMemory()` should validate offset + size against allocation size where possible.
- [ ] **Handle GPU device loss gracefully** — Add helper utilities for detecting and recovering from `VK_ERROR_DEVICE_LOST`.

---

## Code Quality

- [ ] **Reduce duplication in video.go** — The C function pointer loading pattern repeats for each video function; consider a code generator or table-driven approach.
- [ ] **Standardize error wrapping** — Some functions return raw `fmt.Errorf` while others return typed `VulkanError`; standardize on typed errors for all Vulkan calls.
- [ ] **Add `go generate` for constants** — Vulkan constants (formats, flags, result codes) could be auto-generated from `vk.xml` to stay in sync with the Vulkan spec.
- [ ] **Enable additional linters** — Consider enabling `errcheck`, `exhaustive` (enum switch coverage), and `nilnil` in `.golangci.yml`.

---

## Long-Term / Future

- [ ] **Vulkan 1.4 full support** — Implement new 1.4 core features as they are promoted from extensions.
- [ ] **Shader reflection utilities** — Parse SPIR-V to auto-generate descriptor set layouts and pipeline layouts.
- [ ] **Higher-level rendering abstractions** — Optional layer providing render graph, material system, or scene graph on top of the raw Vulkan bindings.
- [ ] **WebGPU compatibility layer** — Thin adapter to map common WebGPU patterns onto Vulkan for cross-API portability.
- [ ] **Automated Vulkan spec sync** — Script or CI job that diffs the latest `vk.xml` against the current bindings and reports missing functions/types.
