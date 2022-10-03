package snapshots

import (
	"path/filepath"
	"time"

	"github.com/tkw1536/goprogram/stream"
)

// ShouldPrune determines if a file with the provided modtime
func (manager *Manager) ShouldPrune(modtime time.Time) bool {
	return time.Since(modtime) > time.Duration(manager.Config.MaxBackupAge)*24*time.Hour
}

// Prune prunes all backups and snapshots older than the maximum backup age
func (manager *Manager) PruneBackups(io stream.IOStream) error {
	sPath := manager.ArchivePath()

	// list all the files
	entries, err := manager.Core.Environment.ReadDir(sPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// skip directories
		if entry.IsDir() {
			continue
		}

		// grab info about the file
		info, err := entry.Info()
		if err != nil {
			return err
		}

		// check if it should be pruned!
		if !manager.ShouldPrune(info.ModTime()) {
			continue
		}

		// assemble path, and then remove the file!
		path := filepath.Join(sPath, entry.Name())
		io.Printf("Removing %s cause it is older than %d days", path, manager.Config.MaxBackupAge)

		if err := manager.Core.Environment.Remove(path); err != nil {
			return err
		}
	}

	// prune the snapshot log!
	_, err = manager.Instances.SnapshotLog()
	return err
}
