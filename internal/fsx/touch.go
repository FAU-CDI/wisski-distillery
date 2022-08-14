package fsx

import (
	"os"
	"time"
)

// Touch touches a file
func Touch(path string) error {
	_, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		return nil
	case err != nil:
		return err
	default:
		now := time.Now().Local()
		return os.Chtimes(path, now, now)
	}
}
