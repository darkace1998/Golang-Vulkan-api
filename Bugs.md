# Bug Registry

This document tracks identified bugs and potential issues in the codebase discovered through deep code analysis.

---

## Bug #1: Unsafe Union Field Handling in CmdClearAttachments

**Channel Operation JIT**: Memory layout operations for clear values

**Channel Read/Write**: Unsafe pointer writes to VkClearValue union fields
**Blocking Semantics**: Pointer arithmetic using Go type sizes for C struct offsets
**Callback Integration**: Could cause memory corruption affecting render pass operations
**All 32 Channels**: Affects all color and depth/stencil clear operations
**Location**: misc.go (CmdClearAttachments function, lines 46-49)

### Description
The code performs unsafe pointer casting to write ClearValue union fields. When writing DepthStencil values, the code manually computes offsets into the VkClearValue union:
```go
*(*float32)(unsafe.Pointer(&cAttachments[i].clearValue)) = att.ClearValue.DepthStencil.Depth
*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&cAttachments[i].clearValue)) + unsafe.Sizeof(att.ClearValue.DepthStencil.Depth))) = att.ClearValue.DepthStencil.Stencil
```
The use of `unsafe.Sizeof(att.ClearValue.DepthStencil.Depth)` (a Go float32) to compute the offset into a C struct may not match the actual VkClearDepthStencilValue layout if there is padding between fields in the C struct.

### Impact
- **Severity**: High
- **Type**: Potential Memory Corruption
- If C struct padding differs from Go's float32 size, the stencil value would be written to the wrong offset

---

## Bug #2: Memory Leak on Allocation Failure in CreateDevice

**Channel Operation JIT**: C memory allocation for queue priorities

**Channel Read/Write**: Sequential malloc calls with partial cleanup on failure
**Blocking Semantics**: Early return paths may leave allocated memory unreleased
**Callback Integration**: Memory pressure could cascade through device creation
**All 32 Channels**: Affects device creation for all queue families
**Location**: device.go (CreateDevice function, lines 226-232)

### Description
When allocating memory for queue priorities fails, the code attempts to clean up previously allocated priorities. However, the cleanup loop only frees priorities already in `cPrioritiesToFree`, and there's no defer mechanism to ensure cleanup in all error paths throughout the function:
```go
if cPrioritiesPtr == nil {
    for _, ptr := range cPrioritiesToFree {
        C.free(unsafe.Pointer(ptr))
    }
    return nil, NewVulkanError(...)
}
```

### Impact
- **Severity**: Medium
- **Type**: Memory Leak
- C memory may not be freed in certain failure scenarios

---

## Bug #3: Unsafe Pointer Arithmetic for Array Indexing

**Channel Operation JIT**: Queue creation info array indexing

**Channel Read/Write**: Manual offset calculation into C struct arrays
**Blocking Semantics**: Misaligned access could cause performance issues or crashes on strict-alignment platforms
**Callback Integration**: Affects device queue setup for all Vulkan operations
**All 32 Channels**: Every device queue family uses this code path
**Location**: device.go (CreateDevice function, lines 217-218)

### Description
The code performs manual pointer arithmetic to access array elements:
```go
offset := uintptr(i) * uintptr(C.sizeof_VkDeviceQueueCreateInfo)
cQueueInfo := (*C.VkDeviceQueueCreateInfo)(unsafe.Pointer(uintptr(...) + offset))
```
While this works in practice for properly aligned structures on most platforms, on architectures with strict alignment requirements (e.g., some ARM platforms), misaligned access could cause performance degradation or faults.

### Impact
- **Severity**: Medium
- **Type**: Potential Misalignment
- Could cause issues on platforms with strict alignment requirements

---

## Summary

| Bug # | Location | Severity | Type |
|-------|----------|----------|------|
| 1 | misc.go:46-49 | High | Memory Corruption |
| 2 | device.go:226-232 | Medium | Memory Leak |
| 3 | device.go:217-218 | Medium | Unsafe Access |
