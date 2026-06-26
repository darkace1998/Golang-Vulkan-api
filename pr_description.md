🎯 **What:**
- Added missing `nil` validation for the `device` and `buffer` parameters in `GetBufferMemoryRequirements`.
- Extracted the CGO wrapper logic into a package-level variable `getBufferMemoryRequirementsFunc` to enable headless testing.
- Added a `TestGetBufferMemoryRequirements` function to verify `GetBufferMemoryRequirements`.

📊 **Coverage:**
- Now covers testing parameter validation (for a nil `device` and a nil `buffer`) where an empty `MemoryRequirements` object should be returned without causing a segment fault.
- Tests the happy path success scenario where `GetBufferMemoryRequirements` correctly delegates to the C wrapper and correctly populates and translates the underlying structure fields.

✨ **Result:**
- Enhanced overall safety of the memory operations API by guarding against driver crashes that occur when accessing missing references in parameters.
- Improved explicit memory test coverage to include `GetBufferMemoryRequirements`.
