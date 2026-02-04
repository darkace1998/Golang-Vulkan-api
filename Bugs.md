# Bug Registry

This document tracks identified bugs and potential issues in the codebase discovered through deep code analysis.

---

## Bug #1: Unsafe Union Field Handling in CmdClearAttachments - **FIXED**

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

### Fix Applied
Changed to use direct C struct field access via `C.VkClearDepthStencilValue` pointer cast:
```go
cDepthStencil := (*C.VkClearDepthStencilValue)(unsafe.Pointer(&cAttachments[i].clearValue))
cDepthStencil.depth = C.float(att.ClearValue.DepthStencil.Depth)
cDepthStencil.stencil = C.uint32_t(att.ClearValue.DepthStencil.Stencil)
```
This ensures proper C struct layout is used for field access.

---

## Bug #2: Memory Leak on Allocation Failure in CreateDevice - **FIXED**

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

### Fix Applied
The issue was actually a **double-free bug** - the manual cleanup at line 287-289 (in the features allocation error path) would free memory that the defer at lines 251-255 would also try to free. Removed the redundant manual cleanup since the defer already handles it properly.

---

## Bug #3: Unsafe Pointer Arithmetic for Array Indexing - **FIXED**

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

### Fix Applied
Changed to use Go's `unsafe.Slice` for safe slice-based indexing:
```go
cQueueInfos := unsafe.Slice(cQueueCreateInfosPtr, len(createInfo.QueueCreateInfos))
// ... then use cQueueInfos[i] for array access
```
Also applied the same fix to the priority array access using `unsafe.Slice(cPrioritiesPtr, len(qci.QueuePriorities))`.

---

## Summary

| Bug # | Location | Severity | Type | Status |
|-------|----------|----------|------|--------|
| 1 | misc.go:46-49 | High | Memory Corruption | **FIXED** |
| 2 | device.go:226-232 | Medium | Memory Leak | **FIXED** |
| 3 | device.go:217-218 | Medium | Unsafe Access | **FIXED** |
