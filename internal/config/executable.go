package config

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
)

// ExecutablePath returns the path to the executable of this distillery.
func (cfg Config) ExecutablePath() string {
	return filepath.Join(cfg.DeployRoot, bootstrap.Executable)
}

// UsingDistilleryExecutable checks if the current process is using the distillery executable
func (cfg Config) UsingDistilleryExecutable(env environment.Environment) bool {
	exe, err := env.Executable()
	if err != nil {
		return false
	}
	return fsx.SameFile(env, exe, cfg.ExecutablePath())
}

// CurrentExecutable returns the path to the current executable being used.
// When it does not exist, falls back to the default executable.
func (cfg Config) CurrentExecutable(env environment.Environment) string {
	exe, err := env.Executable()
	if err != nil || !fsx.IsFile(env, exe) {
		return cfg.ExecutablePath()
	}
	return exe
}
