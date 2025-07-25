//spellchecker:words config
package config

//spellchecker:words path filepath github wisski distillery internal bootstrap pkglib
import (
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"go.tkw01536.de/pkglib/fsx"
)

type PathsConfig struct {
	// Several docker-compose files are created to manage global services and the system itself.
	// On top of this all real-system space will be created under this directory.
	Root string `default:"/var/www/deploy" validate:"directory" yaml:"root"`

	// You can override individual URLS in the homepage
	// Do this by adding URLs (without trailing '/'s) into a JSON file
	OverridesJSON string `validate:"file" yaml:"overrides"`

	// You can block specific prefixes from being picked up by the resolver.
	// Do this by adding one prefix per file.
	ResolverBlocks string `validate:"file" yaml:"blocks"`
}

// RuntimeDir returns the path to the runtime directory.
func (pcfg PathsConfig) RuntimeDir() string {
	return filepath.Join(pcfg.Root, "runtime")
}

// ExecutablePath returns the path to the executable of this distillery.
func (pcfg PathsConfig) ExecutablePath() string {
	return filepath.Join(pcfg.Root, bootstrap.Executable)
}

// UsingDistilleryExecutable checks if the current process is using the distillery executable.
func (pcfg PathsConfig) UsingDistilleryExecutable() bool {
	// TODO: Log
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return fsx.Same(exe, pcfg.ExecutablePath())
}

// CurrentExecutable returns the path to the current executable being used.
// When it does not exist, falls back to the default executable.
func (pcfg PathsConfig) CurrentExecutable() string {
	exe, err := os.Executable()
	if err == nil {
		isFile, err := fsx.IsRegular(exe, true)
		if err == nil && isFile {
			return exe
		}
	}

	return pcfg.ExecutablePath()
}
