package ipmi

import (
	"context"
	"fmt"
	"strings"

	"github.com/bancey/ipmitool-api/internal/config"
)

type PowerStatus struct {
	Status string `json:"status"`
}

func GetPowerStatus(ctx context.Context, exec Executor, server *config.Server) (*PowerStatus, error) {
	output, err := exec.Execute(ctx, server, "power", "status")
	if err != nil {
		return nil, err
	}

	status := "unknown"
	lower := strings.ToLower(output)
	if strings.Contains(lower, "is on") {
		status = "on"
	} else if strings.Contains(lower, "is off") {
		status = "off"
	}

	return &PowerStatus{Status: status}, nil
}

func SetPowerState(ctx context.Context, exec Executor, server *config.Server, action string) error {
	validActions := map[string]bool{
		"on":    true,
		"off":   true,
		"reset": true,
		"cycle": true,
		"soft":  true,
	}

	if !validActions[action] {
		return fmt.Errorf("invalid power action: %q (valid: on, off, reset, cycle, soft)", action)
	}

	_, err := exec.Execute(ctx, server, "power", action)
	return err
}
