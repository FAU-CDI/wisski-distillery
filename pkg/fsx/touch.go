package fsx

import (
	"os"
	"time"
)

// Touch touches a file.
// It is similar to the unix 'touch' command.
//
// If the file does not exist exists, it is created using [os.Create].
// If the file does exist, it's access and modification times are updated to the current time.
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
