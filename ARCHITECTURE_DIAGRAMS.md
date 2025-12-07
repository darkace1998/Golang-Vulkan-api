# Architecture Diagrams - Golang-Vulkan-api Video Extension Loading

## Overall System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Your Go Application                       │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ func main() {                                         │   │
│  │   instance := CreateInstance()                        │   │
│  │   LoadVideoInstanceFunctions(instance)  ← CRITICAL   │   │
│  │                                                       │   │
│  │   device := CreateDevice(physicalDevice)             │   │
│  │   LoadVideoDeviceFunctions(device)      ← CRITICAL   │   │
│  │                                                       │   │
│  │   videoSession := CreateVideoSession()  ← Safe now   │   │
│  │   CmdDecodeVideo(...)                   ← Safe now   │   │
│  │ }                                                     │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
        ↓                          ↓                      ↓
     video.go                  video.go               video.go
   (Go wrappers)            (C loaders)            (C wrappers)
        ↓                          ↓                      ↓
┌─────────────────────────────────────────────────────────────┐
│                      CGO Layer (video.go)                    │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ // Public Go APIs                                      │  │
│  │ func LoadVideoInstanceFunctions(instance) bool { ... } │  │
│  │ func LoadVideoDeviceFunctions(device) bool { ... }     │  │
│  │ func CreateVideoSession(...) (VideoSession, error)     │  │
│  │ func CmdBeginVideoCoding(...) error { ... }            │  │
│  └────────────────────────────────────────────────────────┘  │
│                              ↓                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ // C Function Pointers                                 │  │
│  │ static PFN_vkCreateVideoSessionKHR pfn_... = NULL;     │  │
│  │ static PFN_vkCmdBeginVideoCodingKHR pfn_... = NULL;    │  │
│  │ // ... 10 more function pointers                       │  │
│  └────────────────────────────────────────────────────────┘  │
│                              ↓                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ // Loader Functions                                    │  │
│  │ static int loadVideoDeviceFunctions(VkDevice device) { │  │
│  │   pfn_vkCreateVideoSessionKHR = (...)                  │  │
│  │     vkGetDeviceProcAddr(device, "...");                │  │
│  │   // Check all pointers loaded successfully            │  │
│  │   return all_not_null;                                 │  │
│  │ }                                                      │  │
│  └────────────────────────────────────────────────────────┘  │
│                              ↓                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ // Safe Wrapper Functions                              │  │
│  │ static VkResult call_vkCreateVideoSessionKHR(...) {    │  │
│  │   if (pfn_vkCreateVideoSessionKHR == NULL) {          │  │
│  │     return VK_ERROR_EXTENSION_NOT_PRESENT;            │  │
│  │   }                                                    │  │
│  │   return pfn_vkCreateVideoSessionKHR(...);            │  │
│  │ }                                                      │  │
│  └────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────────────────────────┐
│              Vulkan Loader & Implementation                  │
│                                                               │
│  vkGetDeviceProcAddr/vkGetInstanceProcAddr                  │
│      ↓                                                       │
│  ┌────────────────────────────────────────────────────────┐  │
│  │        Vulkan Video Extension Functions                │  │
│  │  (if available on this GPU/driver)                     │  │
│  │                                                        │  │
│  │  vkCreateVideoSessionKHR()                             │  │
│  │  vkCmdBeginVideoCodingKHR()                            │  │
│  │  vkCmdDecodeVideoKHR()                                 │  │
│  │  etc.                                                  │  │
│  └────────────────────────────────────────────────────────┘  │
│      ↓                                                       │
│  ┌────────────────────────────────────────────────────────┐  │
│  │           GPU Video Hardware                           │  │
│  │  (H.264, H.265, AV1 encode/decode)                    │  │
│  └────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Function Loading Timeline

```
                    Application Startup
                           │
                           ↓
            ┌──────────────────────────────┐
            │  Create Vulkan Instance      │
            │  instance = CreateInstance() │
            └──────────────────────────────┘
                           │
                           ↓
            ┌──────────────────────────────────────────┐
            │  Load Instance Functions (CRITICAL!)      │
            │  LoadVideoInstanceFunctions(instance)     │
            │                                           │
            │  vkGetInstanceProcAddr() called for:       │
            │  - vkGetPhysicalDeviceVideoCapabilitiesKHR│
            └──────────────────────────────────────────┘
                           │
                    ┌──────┴──────┐
                    ↓             ↓
            ┌─────────────┐  ┌────────────────┐
            │  SUCCESS    │  │   FAILURE      │
            │  Can use    │  │   Not supported│
            │  instance   │  │   Return false │
            │  functions  │  │   Log error    │
            └─────────────┘  └────────────────┘
                    │             │
                    │         (STOP HERE)
                    ↓
            ┌──────────────────────────────┐
            │  Create Logical Device       │
            │  device = CreateDevice()     │
            └──────────────────────────────┘
                    │
                    ↓
            ┌──────────────────────────────────────────┐
            │  Load Device Functions (CRITICAL!)        │
            │  LoadVideoDeviceFunctions(device)         │
            │                                           │
            │  vkGetDeviceProcAddr() called for:        │
            │  - vkCreateVideoSessionKHR                │
            │  - vkCmdBeginVideoCodingKHR               │
            │  - vkCmdDecodeVideoKHR                    │
            │  - vkCmdEncodeVideoKHR                    │
            │  + 7 more functions                       │
            │                                           │
            │  Validates ALL 11 functions loaded        │
            └──────────────────────────────────────────┘
                    │
                    ↓
            ┌──────────────────────────────┐
            │  Video APIs Now Safe         │
            │                              │
            │  CreateVideoSession()  OK    │
            │  CmdBeginVideoCoding() OK    │
            │  CmdDecodeVideo()      OK    │
            │  CmdEncodeVideo()      OK    │
            │                              │
            │  All wrapper functions       │
            │  have valid pointers         │
            └──────────────────────────────┘
                    │
                    ↓
            ┌──────────────────────────────┐
            │  Application Uses Video      │
            │  (Normal Vulkan operations)  │
            └──────────────────────────────┘
```

## Error Handling Paths

```
Calling a Video Function
        │
        ↓
    ┌─────────────────────────┐
    │  Wrapper function?      │
    │  (call_vkXxxKHR)        │
    └─────────────────────────┘
        │
        ├─ Yes ──→ Check if pointer is NULL
        │               │
        │               ├─ Is NULL → Return error
        │               │            "Extension not loaded"
        │               │            "Call LoadVideoDeviceFunctions"
        │               │
        │               └─ Not NULL → Call actual function
        │                             Return result
        │
        └─ No → Error
                "Function not wrapped!"
```

## Memory Layout of Function Pointers

```
        .bss section (Global Static Memory)
        ┌──────────────────────────────────────────────────┐
        │ Global Function Pointers                         │
        │ (Initialized to NULL at startup)                 │
        │                                                   │
        │ PFN_vkGetPhysicalDeviceVideoCapabilitiesKHR     │ ← 8 bytes (ptr)
        │   pfn_vkGetPhysicalDeviceVideoCapabilitiesKHR    │
        │   [NULL] → [0x7f1234567890] after loading        │
        │                                                   │
        │ PFN_vkCreateVideoSessionKHR                      │ ← 8 bytes (ptr)
        │   pfn_vkCreateVideoSessionKHR                    │
        │   [NULL] → [0x7f1234567891] after loading        │
        │                                                   │
        │ PFN_vkCmdBeginVideoCodingKHR                     │ ← 8 bytes (ptr)
        │   pfn_vkCmdBeginVideoCodingKHR                   │
        │   [NULL] → [0x7f1234567892] after loading        │
        │                                                   │
        │ ... 9 more function pointers ...                 │
        │                                                   │
        │ Total: 12 × 8 = 96 bytes                         │
        └──────────────────────────────────────────────────┘
                         ↓
                   vkGetDeviceProcAddr
                   gets address of actual
                   Vulkan function from driver
                         ↓
        ┌──────────────────────────────────────────────────┐
        │ Vulkan Driver Library (libvulkan.so)             │
        │                                                   │
        │ vkGetPhysicalDeviceVideoCapabilitiesKHR  ────────┐
        │ [actual code]                                    │
        │                                                  │
        │ vkCreateVideoSessionKHR ──────────────────────┐  │
        │ [actual code]                                │  │
        │                                              │  │
        │ vkCmdBeginVideoCodingKHR ────────────────┐   │  │
        │ [actual code]                            │   │  │
        │                                          │   │  │
        └──────────────────────────────────────────┼───┼──┘
                                                   │   │
        After loading, pointers point here ────────┘   │
                                                       │
        (Addresses stored in our function pointers)────┘
```

## Comparison: Before vs After

### BEFORE (Direct Linking - Broken)
```
Compile Time                  Runtime
─────────────                 ───────
Source Code                   Vulkan Driver
  │                             │
  ├─ #include                   │
  │  <vulkan/vulkan.h>          │
  │                             │
  └─ vkCreateVideoSessionKHR    │
              │                 │
              ├─ Linker tries   │
              │ to find symbol  │
              │                 │
              └─ ERROR!         │
                 "Undefined     │
                  reference"    │
                                │
                          (Never gets here)

❌ RESULT: Build failure even with Vulkan headers installed
```

### AFTER (Dynamic Loading - Works!)
```
Compile Time                  Runtime
─────────────                 ───────
Source Code                   Vulkan Driver
  │                             │
  ├─ #include                   │
  │  <vulkan/vulkan.h>          │
  │  (Standard functions only)  │
  │                             │
  └─ vkGetDeviceProcAddr        │
              │                 │
              └─ Link OK        ├─ Returns pointer to
                 ✓              │  vkCreateVideoSessionKHR
                                │  
                          ✓ Function loaded!
                          ✓ Works!

✅ RESULT: Builds and runs on all systems
```

## Thread Safety Model

```
                     Application Threads
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
    Thread 1            Thread 2            Thread 3
        │                   │                   │
        ↓                   ↓                   ↓
    ┌────────────────────────────────────────────┐
    │ INITIALIZATION PHASE (Must be single-threaded)
    │
    │  LoadVideoInstanceFunctions(instance)
    │        ↓
    │   ┌─────────────────────────────────┐
    │   │ Global function pointers        │
    │   │ being modified:                 │
    │   │ pfn_vkCreateVideoSessionKHR     │
    │   │ pfn_vkCmdBeginVideoCodingKHR    │
    │   │ ... etc                         │
    │   └─────────────────────────────────┘
    │  LoadVideoDeviceFunctions(device)
    │        ↓
    │   ┌─────────────────────────────────┐
    │   │ Function pointers updated       │
    │   │ All 11 functions now ready      │
    │   └─────────────────────────────────┘
    └────────────────────────────────────────────┘
                      ↓
    ┌────────────────────────────────────────────┐
    │ USAGE PHASE (Can be multi-threaded)
    │
    │  Multiple threads can now safely:
    │  - Read function pointers
    │  - Call wrapper functions
    │  - Use video APIs
    │
    │ IMPORTANT: Don't call Load* functions again!
    └────────────────────────────────────────────┘
            ↓
    ┌───────────────────────────────────────────────┐
    │ Thread-Safe Usage (after initialization):
    │
    │ Thread 1: CmdDecodeVideo(...)    ✓
    │ Thread 2: CmdEncodeVideo(...)    ✓
    │ Thread 3: CreateVideoSession(..) ✓
    │
    │ All read from same function pointers
    │ All function calls are safe
    │ No data races because pointers
    │ are read-only during this phase
    └───────────────────────────────────────────────┘
```

## Extension Availability Detection

```
System Initialization
        │
        ├─ Check: Does GPU driver support video?
        │         (via vkGetDeviceProcAddr)
        │
        ├─ YES ──→ LoadVideoDeviceFunctions returns true
        │           ✓ vkCreateVideoSessionKHR available
        │           ✓ vkCmdDecodeVideoKHR available
        │           ✓ vkCmdEncodeVideoKHR available
        │
        └─ NO  ──→ LoadVideoDeviceFunctions returns false
                    ✗ Video extensions not available
                    ✓ But application doesn't crash!
                    ✓ Can fall back to software codec
                    
         Safe degradation based on hardware capability
```

---

These diagrams illustrate:
1. Overall system architecture with layering
2. Initialization timeline and critical steps
3. Error handling paths
4. Memory layout of function pointers
5. Before/after comparison
6. Thread safety model
7. Hardware capability detection

Understanding these diagrams will help you debug issues and understand the design decisions.
