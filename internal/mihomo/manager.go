package mihomo

import "context"

type Options struct {
	BinaryPath          string
	RuntimeDir          string
	ProdConfigPath      string
	ProbeConfigPath     string
	ProdControllerAddr  string
	ProbeControllerAddr string
	ProbeMixedPort      int
}

type Manager struct {
	opts Options
}

func NewManager(opts Options) *Manager {
	return &Manager{opts: opts}
}

func (m *Manager) Start(context.Context, string) error { return nil }
func (m *Manager) Stop()                               {}
func (m *Manager) ProbeMixedPort() int                 { return m.opts.ProbeMixedPort }
