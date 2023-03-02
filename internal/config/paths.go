package config

import (
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
)

type PathsConfig struct {
	// Several docker-compose files are created to manage global services and the system itself.
	// On top of this all real-system space will be created under this directory.
	Root string `yaml:"root" default:"/var/www/deploy" validate:"directory"`

	// You can override individual URLS in the homepage
	// Do this by adding URLs (without trailing '/'s) into a JSON file
	OverridesJSON string `yaml:"overrides" validate:"file"`

	// You can block specific prefixes from being picked up by the resolver.
	// Do this by adding one prefix per file.
	ResolverBlocks string `yaml:"blocks" validate:"file"`
}

// RuntimeDir returns the path to the runtime directory
func (pcfg PathsConfig) RuntimeDir() string {
	return filepath.Join(pcfg.Root, "runtime")
}

// ExecutablePath returns the path to the executable of this distillery.
func (pcfg PathsConfig) ExecutablePath() string {
	return filepath.Join(pcfg.Root, bootstrap.Executable)
}

// UsingDistilleryExecutable checks if the current process is using the distillery executable
func (pcfg PathsConfig) UsingDistilleryExecutable(env environment.Environment) bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return fsx.SameFile(exe, pcfg.ExecutablePath())
}

// CurrentExecutable returns the path to the current executable being used.
// When it does not exist, falls back to the default executable.
func (pcfg PathsConfig) CurrentExecutable(env environment.Environment) string {
	exe, err := os.Executable()
	if err != nil || !fsx.IsFile(exe) {
		return pcfg.ExecutablePath()
	}
	return exe
}
