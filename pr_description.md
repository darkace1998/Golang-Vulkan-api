🧪 Add test for DestroySwapchain nil validation

🎯 **What:** The testing gap addressed
The `DestroySwapchain` function in `swapchain.go` was missing a test to ensure it handles nil arguments correctly without panicking. This aligns it with other destroy functions like `DestroyBuffer` or `DestroyCommandPool` which already have similar nil validation tests.

📊 **Coverage:** What scenarios are now tested
The `TestDestroySwapchainNilArgs` test case covers the scenario of calling `DestroySwapchain` with:
- `nil` for both `device` and `swapchain`.
- `nil` for `device` and a valid `swapchain`.
- A valid `device` and `nil` for `swapchain`.

✨ **Result:** The improvement in test coverage
Increased test coverage for `swapchain.go` and improved the reliability by verifying that `DestroySwapchain` behaves safely even when provided with invalid/nil handles.
