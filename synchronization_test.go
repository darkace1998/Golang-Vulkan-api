package vulkan

import (
	"errors"
	"testing"
)

// ============================================================================
// Semaphore Validation Tests
// ============================================================================

func TestCreateTimelineSemaphoreValidation(t *testing.T) {
	_, err := CreateTimelineSemaphore(nil, 0)
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else {
		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if valErr.Field != testDeviceParameter {
			t.Errorf("Expected field %q, got %q", testDeviceParameter, valErr.Field)
		}
	}
}

func TestWaitSemaphoresValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		waitInfo      *SemaphoreWaitInfo
		expectedField string
	}{
		{
			name:          testNilDevice,
			device:        nil,
			waitInfo:      &SemaphoreWaitInfo{},
			expectedField: testDeviceParameter,
		},
		{
			name:          "nil waitInfo",
			device:        fakeDevice(),
			waitInfo:      nil,
			expectedField: "waitInfo",
		},
		{
			name:          "empty semaphores",
			device:        fakeDevice(),
			waitInfo:      &SemaphoreWaitInfo{},
			expectedField: "Semaphores",
		},
		{
			name:   "mismatched lengths",
			device: fakeDevice(),
			waitInfo: &SemaphoreWaitInfo{
				Semaphores: []Semaphore{fakeSemaphore()},
				Values:     []uint64{1, 2},
			},
			expectedField: "Values",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := WaitSemaphores(tc.device, tc.waitInfo, 0)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			} else {
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				} else if valErr.Field != tc.expectedField {
					t.Errorf("Expected field %q, got %q", tc.expectedField, valErr.Field)
				}
			}
		})
	}
}

func TestSignalSemaphoreValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		signalInfo    *SemaphoreSignalInfo
		expectedField string
	}{
		{
			name:          testNilDevice,
			device:        nil,
			signalInfo:    &SemaphoreSignalInfo{},
			expectedField: testDeviceParameter,
		},
		{
			name:          "nil signalInfo",
			device:        fakeDevice(),
			signalInfo:    nil,
			expectedField: "signalInfo",
		},
		{
			name:          "nil semaphore",
			device:        fakeDevice(),
			signalInfo:    &SemaphoreSignalInfo{},
			expectedField: "Semaphore",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SignalSemaphore(tc.device, tc.signalInfo)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			} else {
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				} else if valErr.Field != tc.expectedField {
					t.Errorf("Expected field %q, got %q", tc.expectedField, valErr.Field)
				}
			}
		})
	}
}

func TestGetSemaphoreCounterValueValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		semaphore     Semaphore
		expectedField string
	}{
		{
			name:          testNilDevice,
			device:        nil,
			semaphore:     fakeSemaphore(),
			expectedField: testDeviceParameter,
		},
		{
			name:          "nil semaphore",
			device:        fakeDevice(),
			semaphore:     nil,
			expectedField: "semaphore",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetSemaphoreCounterValue(tc.device, tc.semaphore)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			} else {
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				} else if valErr.Field != tc.expectedField {
					t.Errorf("Expected field %q, got %q", tc.expectedField, valErr.Field)
				}
			}
		})
	}
}

// ============================================================================
// Event Validation Tests
// ============================================================================

func TestCreateEventValidation(t *testing.T) {
	_, err := CreateEvent(nil, nil)
	if err == nil {
		t.Errorf("Expected error when device is nil")
	} else {
		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if valErr.Field != testDeviceParameter {
			t.Errorf("Expected field %q, got %q", testDeviceParameter, valErr.Field)
		}
	}
}

func TestSetEventValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		event         Event
		expectedField string
	}{
		{
			name:          testNilDevice,
			device:        nil,
			event:         fakeEvent(),
			expectedField: testDeviceParameter,
		},
		{
			name:          testNilEvent,
			device:        fakeDevice(),
			event:         nil,
			expectedField: testEventParameter,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SetEvent(tc.device, tc.event)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			} else {
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				} else if valErr.Field != tc.expectedField {
					t.Errorf("Expected field %q, got %q", tc.expectedField, valErr.Field)
				}
			}
		})
	}
}

func TestResetEventValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		event         Event
		expectedField string
	}{
		{
			name:          testNilDevice,
			device:        nil,
			event:         fakeEvent(),
			expectedField: testDeviceParameter,
		},
		{
			name:          testNilEvent,
			device:        fakeDevice(),
			event:         nil,
			expectedField: testEventParameter,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ResetEvent(tc.device, tc.event)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			} else {
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				} else if valErr.Field != tc.expectedField {
					t.Errorf("Expected field %q, got %q", tc.expectedField, valErr.Field)
				}
			}
		})
	}
}

func TestGetEventStatusValidation(t *testing.T) {
	tests := []struct {
		name          string
		device        Device
		event         Event
		expectedField string
	}{
		{
			name:          testNilDevice,
			device:        nil,
			event:         fakeEvent(),
			expectedField: testDeviceParameter,
		},
		{
			name:          testNilEvent,
			device:        fakeDevice(),
			event:         nil,
			expectedField: testEventParameter,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetEventStatus(tc.device, tc.event)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			} else {
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					t.Errorf("Expected ValidationError, got %T", err)
				} else if valErr.Field != tc.expectedField {
					t.Errorf("Expected field %q, got %q", tc.expectedField, valErr.Field)
				}
			}
		})
	}
}

// ============================================================================
// Nil Check Tests for Destroy and Command Functions
// ============================================================================

func TestDestroyEventNilArgs(t *testing.T) {
	DestroyEvent(nil, nil)
	DestroyEvent(nil, fakeEvent())
	DestroyEvent(fakeDevice(), nil)
}

func TestCmdSetEventNilArgs(t *testing.T) {
	CmdSetEvent(nil, nil, 0)
	CmdSetEvent(nil, fakeEvent(), 0)
	CmdSetEvent(fakeCommandBuffer(), nil, 0)
}

func TestCmdResetEventNilArgs(t *testing.T) {
	CmdResetEvent(nil, nil, 0)
	CmdResetEvent(nil, fakeEvent(), 0)
	CmdResetEvent(fakeCommandBuffer(), nil, 0)
}

func TestCmdPipelineBarrierFullNilArgs(t *testing.T) {
	CmdPipelineBarrierFull(nil, 0, 0, 0, nil, nil, nil)
}

func TestCmdWaitEventsNilArgs(t *testing.T) {
	CmdWaitEvents(nil, nil, 0, 0, nil, nil, nil)
	CmdWaitEvents(nil, []Event{fakeEvent()}, 0, 0, nil, nil, nil)
	CmdWaitEvents(fakeCommandBuffer(), nil, 0, 0, nil, nil, nil)
	CmdWaitEvents(fakeCommandBuffer(), []Event{}, 0, 0, nil, nil, nil)
}
