package env

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

func (dis Distillery) BackupDir() string {
	return filepath.Join(dis.Config.DeployRoot, "backups")
}

func (dis Distillery) InprogressBackupPath() string {
	return filepath.Join(dis.BackupDir(), "inprogress")
}

func (dis Distillery) FinalBackupPath() string {
	return filepath.Join(dis.BackupDir(), "final")
}

// NewFinalBackupFile returns the path to a new final backup file.
func (dis Distillery) FinalBackupArchive(prefix string) string {
	counter := atomic.AddUint64(&globalBackupCounter, 1)

	// generate a new name with the current time, a global counter, and the prefix
	name := fmt.Sprintf("%s-%d-%d.tar.gz", prefix, time.Now().Unix(), counter)
	path := filepath.Join(dis.FinalBackupPath(), name)

	return path
}

var globalBackupCounter uint64

// NewInprogressBackupPath returns the path to a new inprogress backup directory.
// The directory is guaranteed to have been freshly created.
func (dis Distillery) NewInprogressBackupPath(prefix string) (string, error) {
	counter := atomic.AddUint64(&globalBackupCounter, 1)

	// generate a new name with the current time, a global counter, and the prefix
	name := fmt.Sprintf("%s-%d-%d", prefix, time.Now().Unix(), counter)
	path := filepath.Join(dis.InprogressBackupPath(), name)

	// create the directory
	if err := os.Mkdir(path, os.ModeDir); err != nil {
		return "", err
	}

	// and it is here!
	return path, nil
}
