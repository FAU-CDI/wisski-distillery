package config

import (
	"embed"
	"path/filepath"
)

// Runtime contains runtime resources to be installed into any instance
//go:embed all:runtime
var Runtime embed.FS

// RuntimeDir returns the path to the runtime directory
func (cfg Config) RuntimeDir() string {
	return filepath.Join(cfg.DeployRoot, "runtime")
}
