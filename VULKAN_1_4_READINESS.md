# Vulkan 1.4 Readiness

This document outlines the current state of Vulkan 1.4 readiness in the `golang-vulkan-api` library.

## Current State

As of the current release, the library provides fundamental support for Vulkan 1.4 API versioning:

*   **Version Definition:** The constant `vulkan.Version14` is defined in `types.go`. This corresponds to `VK_MAKE_API_VERSION(0, 1, 4, 0)`.
*   **Version Querying:** You can use `vulkan.Version14` to check if the underlying Vulkan implementation supports version 1.4 (e.g., using `EnumerateInstanceVersion()`).
*   **Instance Creation:** You can request a Vulkan 1.4 instance by setting `APIVersion: vulkan.Version14` in your `ApplicationInfo` struct when calling `CreateInstance()`.

**Important Note:** Requesting `vulkan.Version14` during instance creation will only succeed if the user's system has a Vulkan driver and runtime installed that supports Vulkan 1.4. Otherwise, instance creation will fail with an error like `VK_ERROR_INCOMPATIBLE_DRIVER`.

### Example Usage

```go
package main

import (
    "fmt"
    "log"
    vulkan "github.com/darkace1998/golang-vulkan-api"
)

func main() {
    // 1. Check supported instance version (optional but recommended)
    supportedVersion, err := vulkan.EnumerateInstanceVersion()
    if err == nil {
        fmt.Printf("Supported Instance Version: %d.%d.%d\n",
            supportedVersion.Major(), supportedVersion.Minor(), supportedVersion.Patch())

        if supportedVersion >= vulkan.Version14 {
            fmt.Println("Vulkan 1.4 is supported by the system!")
        } else {
            fmt.Println("Vulkan 1.4 is NOT supported by the system.")
        }
    }

    // 2. Request Vulkan 1.4 during instance creation
    instanceCreateInfo := &vulkan.InstanceCreateInfo{
        ApplicationInfo: &vulkan.ApplicationInfo{
            ApplicationName:    "My Vulkan 1.4 App",
            ApplicationVersion: vulkan.MakeVersion(1, 0, 0),
            EngineName:         "My Engine",
            EngineVersion:      vulkan.MakeVersion(1, 0, 0),
            // Request Vulkan 1.4
            APIVersion:         vulkan.Version14,
        },
    }

    instance, err := vulkan.CreateInstance(instanceCreateInfo)
    if err != nil {
        log.Fatalf("Failed to create Vulkan instance: %v (Driver might not support 1.4)", err)
    }
    defer vulkan.DestroyInstance(instance)

    fmt.Println("Successfully created a Vulkan instance (requested 1.4).")
}
```

## Planned Features

Full support for Vulkan 1.4 is a planned long-term goal. As documented in our roadmap (`todo.md`):

*   **Vulkan 1.4 Full Support:** We plan to implement new 1.4 core features as they are promoted from extensions into the core API.
*   **Implementation Strategy:** New types, constants, and functions specific to Vulkan 1.4 will likely be added in dedicated files (e.g., similar to the existing `vulkan13.go`) to maintain organization and backwards compatibility where possible.

We will continually monitor the Vulkan specification and update this repository to expose new 1.4 capabilities to Go developers.
