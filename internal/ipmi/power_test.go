package ipmi

import (
	"context"
	"fmt"
	"testing"

	"github.com/bancey/ipmitool-api/internal/config"
)

var testServer = &config.Server{Host: "10.0.0.1", Port: 623, Username: "admin", Password: "pass"}

func TestGetPowerStatusOn(t *testing.T) {
	mock := &mockExecutor{output: "Chassis Power is on"}
	status, err := GetPowerStatus(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Status != "on" {
		t.Errorf("expected 'on', got %q", status.Status)
	}
}

func TestGetPowerStatusOff(t *testing.T) {
	mock := &mockExecutor{output: "Chassis Power is off"}
	status, err := GetPowerStatus(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Status != "off" {
		t.Errorf("expected 'off', got %q", status.Status)
	}
}

func TestGetPowerStatusUnknown(t *testing.T) {
	mock := &mockExecutor{output: "something unexpected"}
	status, err := GetPowerStatus(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Status != "unknown" {
		t.Errorf("expected 'unknown', got %q", status.Status)
	}
}

func TestGetPowerStatusError(t *testing.T) {
	mock := &mockExecutor{err: fmt.Errorf("connection refused")}
	_, err := GetPowerStatus(context.Background(), mock, testServer)
	if err == nil {
		t.Error("expected error from executor")
	}
}

func TestSetPowerStateValid(t *testing.T) {
	for _, action := range []string{"on", "off", "reset", "cycle", "soft"} {
		t.Run(action, func(t *testing.T) {
			mock := &mockExecutor{output: ""}
			err := SetPowerState(context.Background(), mock, testServer, action)
			if err != nil {
				t.Errorf("unexpected error for action %q: %v", action, err)
			}
			if len(mock.called) != 2 || mock.called[0] != "power" || mock.called[1] != action {
				t.Errorf("expected args [power %s], got %v", action, mock.called)
			}
		})
	}
}

func TestSetPowerStateInvalid(t *testing.T) {
	mock := &mockExecutor{output: ""}
	err := SetPowerState(context.Background(), mock, testServer, "destroy")
	if err == nil {
		t.Error("expected error for invalid action")
	}
}

func TestSetPowerStateExecutorError(t *testing.T) {
	mock := &mockExecutor{err: fmt.Errorf("timeout")}
	err := SetPowerState(context.Background(), mock, testServer, "on")
	if err == nil {
		t.Error("expected error from executor")
	}
}
