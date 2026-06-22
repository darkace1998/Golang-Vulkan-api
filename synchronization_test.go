package vulkan

import (
	"testing"
)

func TestWaitSemaphoresValidation(t *testing.T) {
	err := WaitSemaphores(nil, &SemaphoreWaitInfo{}, 0)
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	err = WaitSemaphores(fakeDevice(), nil, 0)
	if err == nil {
		t.Errorf("Expected error when waitInfo is nil")
	} else if err.Error() != "vulkan validation error: waitInfo cannot be nil" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSignalSemaphoreValidation(t *testing.T) {
	err := SignalSemaphore(nil, &SemaphoreSignalInfo{})
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	err = SignalSemaphore(fakeDevice(), nil)
	if err == nil {
		t.Errorf("Expected error when signalInfo is nil")
	} else if err.Error() != "vulkan validation error: signalInfo cannot be nil" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestGetSemaphoreCounterValueValidation(t *testing.T) {
	_, err := GetSemaphoreCounterValue(nil, fakeSemaphore())
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	_, err = GetSemaphoreCounterValue(fakeDevice(), nil)
	if err == nil {
		t.Errorf("Expected error when semaphore is nil")
	} else if err.Error() != "vulkan validation error: semaphore cannot be nil" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestCreateEventValidation(t *testing.T) {
	_, err := CreateEvent(nil, &EventCreateInfo{})
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	// CreateEvent uses a default EventCreateInfo if nil, so we don't expect a validation error for it
}

func TestDestroyEventNilArgs(t *testing.T) {
	// Should not panic
	DestroyEvent(nil, nil)
	DestroyEvent(fakeDevice(), nil)
	DestroyEvent(nil, fakeEvent())
}

func TestSetEventValidation(t *testing.T) {
	err := SetEvent(nil, fakeEvent())
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	err = SetEvent(fakeDevice(), nil)
	if err == nil {
		t.Errorf("Expected error when event is nil")
	} else if err.Error() != testNilEventError {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestResetEventValidation(t *testing.T) {
	err := ResetEvent(nil, fakeEvent())
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	err = ResetEvent(fakeDevice(), nil)
	if err == nil {
		t.Errorf("Expected error when event is nil")
	} else if err.Error() != testNilEventError {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestGetEventStatusValidation(t *testing.T) {
	_, err := GetEventStatus(nil, fakeEvent())
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else if err.Error() != testNilDeviceError {
		t.Errorf("Unexpected error message: %v", err)
	}

	_, err = GetEventStatus(fakeDevice(), nil)
	if err == nil {
		t.Errorf("Expected error when event is nil")
	} else if err.Error() != testNilEventError {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestCmdSetEventNilArgs(t *testing.T) {
	// Should not panic
	CmdSetEvent(nil, nil, 0)
	CmdSetEvent(fakeCommandBuffer(), nil, 0)
	CmdSetEvent(nil, fakeEvent(), 0)
}

func TestCmdResetEventNilArgs(t *testing.T) {
	// Should not panic
	CmdResetEvent(nil, nil, 0)
	CmdResetEvent(fakeCommandBuffer(), nil, 0)
	CmdResetEvent(nil, fakeEvent(), 0)
}

func TestCmdWaitEventsNilArgs(t *testing.T) {
	// Should not panic
	CmdWaitEvents(nil, nil, 0, 0, nil, nil, nil)
	CmdWaitEvents(fakeCommandBuffer(), nil, 0, 0, nil, nil, nil)
	CmdWaitEvents(nil, []Event{fakeEvent()}, 0, 0, nil, nil, nil)
}
