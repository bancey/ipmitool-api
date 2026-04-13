package ipmi

import (
	"bufio"
	"context"
	"strings"

	"github.com/bancey/ipmitool-api/internal/config"
)

type ChassisStatus struct {
	PowerOn        bool   `json:"power_on"`
	PowerOverload  bool   `json:"power_overload"`
	PowerFault     bool   `json:"power_fault"`
	MainPowerFault bool   `json:"main_power_fault"`
	LastPowerEvent string `json:"last_power_event"`
	DrivesFault    bool   `json:"drives_fault"`
	CoolingFault   bool   `json:"cooling_fault"`
}

func GetChassisStatus(ctx context.Context, exec Executor, server *config.Server) (*ChassisStatus, error) {
	output, err := exec.Execute(ctx, server, "chassis", "status")
	if err != nil {
		return nil, err
	}

	return parseChassisStatus(output), nil
}

func parseChassisStatus(output string) *ChassisStatus {
	status := &ChassisStatus{}
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(strings.ToLower(parts[1]))

		switch key {
		case "system power":
			status.PowerOn = value == "on"
		case "power overload":
			status.PowerOverload = value == "true"
		case "power fault":
			status.PowerFault = value == "true"
		case "main power fault":
			status.MainPowerFault = value == "true"
		case "last power event":
			status.LastPowerEvent = strings.TrimSpace(parts[1])
		case "drive fault":
			status.DrivesFault = value == "true"
		case "cooling/fan fault":
			status.CoolingFault = value == "true"
		}
	}

	return status
}
