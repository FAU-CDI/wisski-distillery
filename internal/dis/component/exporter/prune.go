//spellchecker:words exporter
package exporter

//spellchecker:words context path filepath time github wisski distillery internal component
import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

// ShouldPrune determines if a file with the provided modification time should be
// removed from the export log.
func (exporter *Exporter) ShouldPrune(modtime time.Time) bool {
	return time.Since(modtime) > component.GetStill(exporter).Config.MaxBackupAge
}

// Prune prunes all old exports
func (exporter *Exporter) PruneExports(ctx context.Context, progress io.Writer) error {
	sPath := exporter.ArchivePath()

	// list all the files
	entries, err := os.ReadDir(sPath)
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
		fmt.Fprintf(progress, "Removing %s cause it is older than %d days\n", path, component.GetStill(exporter).Config.MaxBackupAge)

		if err := os.Remove(path); err != nil {
			return err
		}
	}

	// prune the snapshot log!
	_, err = exporter.dependencies.ExporterLogger.Log(ctx)
	return err
}
