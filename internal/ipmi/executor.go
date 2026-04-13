package ipmi

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/bancey/ipmitool-api/internal/config"
)

// Executor runs ipmitool commands and returns their output.
type Executor interface {
	Execute(ctx context.Context, server *config.Server, args ...string) (string, error)
}

// CommandExecutor runs actual ipmitool commands on the system.
type CommandExecutor struct{}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

func (e *CommandExecutor) Execute(ctx context.Context, server *config.Server, args ...string) (string, error) {
	cmdArgs := []string{
		"-I", "lanplus",
		"-H", server.Host,
		"-p", fmt.Sprintf("%d", server.Port),
		"-U", server.Username,
		"-P", server.Password,
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.CommandContext(ctx, "ipmitool", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ipmitool error: %s: %w", strings.TrimSpace(string(output)), err)
	}

	return strings.TrimSpace(string(output)), nil
}
