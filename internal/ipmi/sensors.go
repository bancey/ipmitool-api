package ipmi

import (
	"bufio"
	"context"
	"strings"

	"github.com/bancey/ipmitool-api/internal/config"
)

type Sensor struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Units  string `json:"units"`
	Status string `json:"status"`
}

func GetSensors(ctx context.Context, exec Executor, server *config.Server) ([]Sensor, error) {
	output, err := exec.Execute(ctx, server, "sensor", "list")
	if err != nil {
		return nil, err
	}

	return parseSensorOutput(output), nil
}

func parseSensorOutput(output string) []Sensor {
	var sensors []Sensor
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}

		sensor := Sensor{
			Name:   strings.TrimSpace(parts[0]),
			Value:  strings.TrimSpace(parts[1]),
			Units:  strings.TrimSpace(parts[2]),
			Status: strings.TrimSpace(parts[3]),
		}
		sensors = append(sensors, sensor)
	}

	return sensors
}
