package env

import "path/filepath"

// RuntimeDir returns the path to the runtime directory
func (dis Distillery) RuntimeDir() string {
	return filepath.Join(dis.Config.DeployRoot, "runtime")
}
