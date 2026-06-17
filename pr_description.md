🎯 **What:** Added comprehensive validation and nil-check tests for functions in `synchronization.go` to ensure robust error handling.

📊 **Coverage:**
- Added test coverage for `CreateTimelineSemaphore`, `WaitSemaphores`, `SignalSemaphore`, and `GetSemaphoreCounterValue`.
- Added test coverage for `CreateEvent`, `SetEvent`, `ResetEvent`, and `GetEventStatus`.
- Added nil argument check coverage for `DestroyEvent`, `CmdSetEvent`, `CmdResetEvent`, `CmdPipelineBarrierFull`, and `CmdWaitEvents`.
- Updated `test_helpers_test.go` to support `fakeEvent()` for testing synchronization events.

✨ **Result:** Improved overall test coverage for synchronization primitives, ensuring `vulkan` package validation logic triggers correctly before making CGO calls.
