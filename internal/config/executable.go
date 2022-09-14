package config

import (
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
)

// ExecutablePath returns the path to the executable of this distillery.
func (cfg Config) ExecutablePath() string {
	return filepath.Join(cfg.DeployRoot, core.Executable)
}

// UsingDistilleryExecutable checks if the current process is using the distillery executable
func (cfg Config) UsingDistilleryExecutable() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return fsx.SameFile(exe, cfg.ExecutablePath())
}

// CurrentExecutable returns the path to the current executable being used.
// When it does not exist, falls back to the default executable.
func (cfg Config) CurrentExecutable() string {
	exe, err := os.Executable()
	if err != nil || !fsx.IsFile(exe) {
		return cfg.ExecutablePath()
	}
	return exe
}
