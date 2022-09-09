package env

import (
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
)

// Executable is the name of the 'wdcli' executable.
// It should be located inside the deployment directory.
const Executable = "wdcli"

// ExecutablePath returns the path to the executable of this distillery.
func (dis *Distillery) ExecutablePath() string {
	return filepath.Join(dis.Config.DeployRoot, Executable)
}

// UsingDistilleryExecutable checks if the current process
func (dis *Distillery) UsingDistilleryExecutable() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return fsx.SameFile(exe, dis.ExecutablePath())
}

// Config file is the name of the config file.
// It should be located inside the deployment directory.
const ConfigFile = ".env"

// ConfigFilePath returns the path to the configuration file of this distillery.
// TODO: This should be moved to the Config struct.
func (dis *Distillery) ConfigFilePath() string {
	return filepath.Join(dis.Config.DeployRoot, ConfigFile)
}
