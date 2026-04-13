package mihomo

import (
	"context"
	"os"
	"path/filepath"
)

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

func (m *Manager) Start(context.Context, string) error {
	return os.MkdirAll(m.opts.RuntimeDir, 0o755)
}

func (m *Manager) Stop() {}

func (m *Manager) ProbeMixedPort() int {
	return m.opts.ProbeMixedPort
}

func (m *Manager) ProdControllerAddr() string {
	return m.opts.ProdControllerAddr
}

func (m *Manager) ProbeControllerAddr() string {
	return m.opts.ProbeControllerAddr
}

func (m *Manager) ApplyProdConfig(payload []byte) error {
	return writeFileAtomic(m.opts.ProdConfigPath, payload)
}

func (m *Manager) ApplyProbeConfig(payload []byte) error {
	return writeFileAtomic(m.opts.ProbeConfigPath, payload)
}

func writeFileAtomic(path string, payload []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, payload, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
