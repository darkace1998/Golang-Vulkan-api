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
🎯 **What:** Extracted manual C array cleanup code from `CreateInstance` and placed it into a single `defer` statement at the top of the function to reduce repetitive memory cleanups.
💡 **Why:** `CreateInstance` contained a lot of repetitive manual cleanup if early allocations failed, making it error-prone and overly complex. Moving cleanup logic into a unified `defer` block ensures reliable and DRY resource teardown on error exits.
✅ **Verification:** Verified changes through inspection of `git diff` output, ensuring redundant `C.free` code segments were removed. Verified code functionality by ensuring no regressions by running tests via `make test` and specific matching via `go test -v ./... | grep -C 5 "CreateInstance"`.
✨ **Result:** Improved overall codebase maintainability by reducing function length and decreasing probability of memory leaks on function failures.
