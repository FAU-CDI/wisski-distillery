package execx

import (
	"os/exec"
	"path/filepath"
)

// LookPathAbs is like [exec.LookPath], but always returns an absolute path
func LookPathAbs(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}
