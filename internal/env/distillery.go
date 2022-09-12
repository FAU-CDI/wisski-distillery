package env

import (
	"context"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
)

// Distillery represents an interface to the running distillery.
type Distillery struct {
	// Config holds the configuration of the distillery.
	// It is read directly from a configuration file.
	Config *config.Config

	// Upstream holds information to connect to the various running
	// distillery components.
	//
	// NOTE(twiesing): This is intended to eventually allow full remote management of the distillery.
	// But for now this will just hold upstream configuration.
	Upstream Upstream

	// components hold references to the various components of the distillery.
	components
}

// Upstream are the upstream urls connecting to the various external components.
type Upstream struct {
	SQL         string
	Triplestore string
}

// Context returns a new Context belonging to this distillery
func (dis *Distillery) Context() context.Context {
	return context.Background()
}

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
