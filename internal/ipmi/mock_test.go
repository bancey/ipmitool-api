package ipmi

import (
	"context"

	"github.com/bancey/ipmitool-api/internal/config"
)

type mockExecutor struct {
	output string
	err    error
	called []string
}

func (m *mockExecutor) Execute(_ context.Context, _ *config.Server, args ...string) (string, error) {
	m.called = args
	return m.output, m.err
}
