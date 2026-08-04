//go:build !(amd64 || arm64 || riscv64 || loong64 || ppc64le || s390x)

package vulkan

// This package models Vulkan non-dispatchable handles as unsafe.Pointer,
// which matches VK_DEFINE_NON_DISPATCHABLE_HANDLE only on 64-bit targets.
// On 32-bit architectures Vulkan defines those handles as uint64_t, so the
// package cannot compile there. See MULTIPLATFORM.md.
//
// The undefined identifier below produces a clear compile-time error on
// unsupported architectures instead of hundreds of cgo type mismatches.
var _ = golang_vulkan_api_requires_a_64bit_architecture_see_MULTIPLATFORM_md
