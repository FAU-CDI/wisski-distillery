package wisski

import (
	"path/filepath"
	"time"

	"github.com/tkw1536/goprogram/stream"
)

// ShouldPrune determines if a file with the provided modtime
func (dis *Distillery) ShouldPrune(modtime time.Time) bool {
	return time.Since(modtime) > time.Duration(dis.Config.MaxBackupAge)*24*time.Hour
}

// PruneBackups prunes all backups older than the maximum backup age
func (dis *Distillery) PruneBackups(io stream.IOStream) error {
	sPath := dis.SnapshotsArchivePath()

	// list all the files
	entries, err := dis.Core.Environment.ReadDir(sPath)
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
		if !dis.ShouldPrune(info.ModTime()) {
			continue
		}

		// assemble path, and then remove the file!
		path := filepath.Join(sPath, entry.Name())
		io.Printf("Removing %s cause it is older than %d days", path, dis.Config.MaxBackupAge)

		if err := dis.Core.Environment.Remove(path); err != nil {
			return err
		}
	}
	return nil
}
