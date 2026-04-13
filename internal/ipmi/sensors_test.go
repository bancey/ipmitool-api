package ipmi

import (
	"context"
	"fmt"
	"testing"
)

func TestGetSensors(t *testing.T) {
	output := `CPU Temp         | 45.000     | degrees C  | ok
Fan1             | 3200.000   | RPM        | ok
Vcore            | 1.020      | Volts      | ok`

	mock := &mockExecutor{output: output}
	sensors, err := GetSensors(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sensors) != 3 {
		t.Fatalf("expected 3 sensors, got %d", len(sensors))
	}

	if sensors[0].Name != "CPU Temp" {
		t.Errorf("expected name 'CPU Temp', got %q", sensors[0].Name)
	}
	if sensors[0].Value != "45.000" {
		t.Errorf("expected value '45.000', got %q", sensors[0].Value)
	}
	if sensors[0].Units != "degrees C" {
		t.Errorf("expected units 'degrees C', got %q", sensors[0].Units)
	}
	if sensors[0].Status != "ok" {
		t.Errorf("expected status 'ok', got %q", sensors[0].Status)
	}

	if sensors[1].Name != "Fan1" {
		t.Errorf("expected name 'Fan1', got %q", sensors[1].Name)
	}
}

func TestGetSensorsEmpty(t *testing.T) {
	mock := &mockExecutor{output: ""}
	sensors, err := GetSensors(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sensors) != 0 {
		t.Errorf("expected 0 sensors, got %d", len(sensors))
	}
}

func TestGetSensorsMalformedLines(t *testing.T) {
	output := `CPU Temp         | 45.000     | degrees C  | ok
bad line without pipes
Fan1             | 3200.000   | RPM        | ok`

	mock := &mockExecutor{output: output}
	sensors, err := GetSensors(context.Background(), mock, testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sensors) != 2 {
		t.Errorf("expected 2 sensors (skipping malformed), got %d", len(sensors))
	}
}

func TestGetSensorsError(t *testing.T) {
	mock := &mockExecutor{err: fmt.Errorf("connection refused")}
	_, err := GetSensors(context.Background(), mock, testServer)
	if err == nil {
		t.Error("expected error from executor")
	}
}

func TestParseSensorOutput(t *testing.T) {
	output := `Inlet Temp       | 22.000     | degrees C  | ok
Exhaust Temp     | 35.000     | degrees C  | ok
Fan1 RPM         | na         | RPM        | na`

	sensors := parseSensorOutput(output)
	if len(sensors) != 3 {
		t.Fatalf("expected 3 sensors, got %d", len(sensors))
	}

	if sensors[2].Value != "na" {
		t.Errorf("expected 'na' value, got %q", sensors[2].Value)
	}
	if sensors[2].Status != "na" {
		t.Errorf("expected 'na' status, got %q", sensors[2].Status)
	}
}
