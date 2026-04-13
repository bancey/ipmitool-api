package ipmi

import (
	"context"
	"fmt"
	"testing"
)

func TestGetChassisStatus(t *testing.T) {
	output := `System Power         : on
Power Overload       : false
Power Interlock      : inactive
Main Power Fault     : false
Power Control Fault  : false
Power Restore Policy : always-off
Last Power Event     : command
Chassis Intrusion    : inactive
Front-Panel Lockout  : inactive
Drive Fault          : false
Cooling/Fan Fault    : false`

	mock := &mockExecutor{output: output}
	status, err := GetChassisStatus(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !status.PowerOn {
		t.Error("expected PowerOn to be true")
	}
	if status.PowerOverload {
		t.Error("expected PowerOverload to be false")
	}
	if status.MainPowerFault {
		t.Error("expected MainPowerFault to be false")
	}
	if status.LastPowerEvent != "command" {
		t.Errorf("expected LastPowerEvent 'command', got %q", status.LastPowerEvent)
	}
	if status.DrivesFault {
		t.Error("expected DrivesFault to be false")
	}
	if status.CoolingFault {
		t.Error("expected CoolingFault to be false")
	}
}

func TestGetChassisStatusPowerOff(t *testing.T) {
	output := `System Power         : off
Power Overload       : false
Main Power Fault     : false
Last Power Event     : 
Drive Fault          : false
Cooling/Fan Fault    : false`

	mock := &mockExecutor{output: output}
	status, err := GetChassisStatus(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.PowerOn {
		t.Error("expected PowerOn to be false")
	}
}

func TestGetChassisStatusFaults(t *testing.T) {
	output := `System Power         : on
Power Overload       : true
Power Fault          : true
Main Power Fault     : true
Last Power Event     : ac-failed
Drive Fault          : true
Cooling/Fan Fault    : true`

	mock := &mockExecutor{output: output}
	status, err := GetChassisStatus(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !status.PowerOverload {
		t.Error("expected PowerOverload to be true")
	}
	if !status.PowerFault {
		t.Error("expected PowerFault to be true")
	}
	if !status.MainPowerFault {
		t.Error("expected MainPowerFault to be true")
	}
	if !status.DrivesFault {
		t.Error("expected DrivesFault to be true")
	}
	if !status.CoolingFault {
		t.Error("expected CoolingFault to be true")
	}
	if status.LastPowerEvent != "ac-failed" {
		t.Errorf("expected LastPowerEvent 'ac-failed', got %q", status.LastPowerEvent)
	}
}

func TestGetChassisStatusError(t *testing.T) {
	mock := &mockExecutor{err: fmt.Errorf("connection refused")}
	_, err := GetChassisStatus(context.Background(), mock, testServer)
	if err == nil {
		t.Error("expected error from executor")
	}
}

func TestParseChassisStatusEmpty(t *testing.T) {
	status := parseChassisStatus("")
	if status.PowerOn {
		t.Error("expected PowerOn to be false for empty output")
	}
}
