package exporter

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

// ShouldPrune determines if a file with the provided modification time should be
// removed from the export log.
func (exporter *Exporter) ShouldPrune(modtime time.Time) bool {
	return time.Since(modtime) > time.Duration(exporter.Config.MaxBackupAge)*24*time.Hour
}

// Prune prunes all old exports
func (exporter *Exporter) PruneExports(ctx context.Context, progress io.Writer) error {
	sPath := exporter.ArchivePath()

	// list all the files
	entries, err := exporter.Still.Environment.ReadDir(sPath)
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
		if !exporter.ShouldPrune(info.ModTime()) {
			continue
		}

		// assemble path, and then remove the file!
		path := filepath.Join(sPath, entry.Name())
		logging.ProgressF(progress, ctx, "Removing %s cause it is older than %d days\n", path, exporter.Config.MaxBackupAge)

		if err := exporter.Still.Environment.Remove(path); err != nil {
			return err
		}
	}

	// prune the snapshot log!
	_, err = exporter.ExporterLogger.Log(ctx)
	return err
}
