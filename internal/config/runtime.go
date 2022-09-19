package config

import (
	"path/filepath"
)

// RuntimeDir returns the path to the runtime directory
func (cfg Config) RuntimeDir() string {
	return filepath.Join(cfg.DeployRoot, "runtime")
}
