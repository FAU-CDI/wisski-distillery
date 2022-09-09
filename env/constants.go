package env

import (
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
)

// ExecutablePath returns the path to the executable of this distillery.
func (dis *Distillery) ExecutablePath() string {
	return filepath.Join(dis.Config.DeployRoot, core.Executable)
}

// UsingDistilleryExecutable checks if the current process
func (dis *Distillery) UsingDistilleryExecutable() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return fsx.SameFile(exe, dis.ExecutablePath())
}

// CurrentExecutable returns the path to the current executable being used.
// When it does not exist, falls back to the default executable.
func (dis *Distillery) CurrentExecutable() string {
	exe, err := os.Executable()
	if err != nil || !fsx.IsFile(exe) {
		return dis.ExecutablePath()
	}
	return exe
}
