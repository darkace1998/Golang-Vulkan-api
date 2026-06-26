# 🧪 [Testing Improvement] Add TestCreateBufferValidation

## Description

🎯 **What:**
Added unit tests to cover input validation for the `CreateBuffer` function. Specifically, it verifies that the function correctly returns a `ValidationError` when provided with invalid inputs.

📊 **Coverage:**
The following scenarios are now tested:
*   **nil device:** Verifies that passing a `nil` device handle returns an error for the `"device"` field.
*   **nil createInfo:** Verifies that passing a `nil` `BufferCreateInfo` pointer returns an error for the `"createInfo"` field.
*   **zero size:** Verifies that attempting to create a buffer with a size of 0 returns an error for the `"Size"` field.
*   **exceeds max size:** Verifies that attempting to create a buffer exceeding the 1GB safety limit returns an error for the `"Size"` field.
*   **zero usage:** Verifies that attempting to create a buffer with 0 usage flags returns an error for the `"Usage"` field.

✨ **Result:**
This testing improvement increases the test coverage of `memory.go` by properly testing all of the preliminary validation checks within the `CreateBuffer` function. This provides a safety net ensuring that future refactoring will not accidentally remove or break these critical validation checks.
