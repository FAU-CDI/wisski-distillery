//spellchecker:words exporter
package exporter

//spellchecker:words context errors path filepath github wisski distillery internal component models logging targz pkglib collection umaskfree status
import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/targz"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/fsx/umaskfree"
	"go.tkw01536.de/pkglib/status"
)

// See [Exporter.MakeExport].
type ExportTask struct {
	// Dest is the destination path to write the backup to.
	// When empty, this is created automatically in the staging or archive directory.
	Dest string

	// By default, a .tar.gz file is generated.
	// To generated an unpacked directory, set [StagingOnly] to true.
	StagingOnly bool

	// Parts explicitly lists parts to include inside the snapshot.
	// If non-empty, only include parts with the specified names.
	// if empty, include all possible components.
	Parts []string

	// Instance is the instance to generate a snapshot of.
	// To generate a backup, leave this to be nil.
	Instance *wisski.WissKI

	// BackupDescriptions and SnapshotDescriptions further specitfy options for the export.
	// The Dest parameter is ignored, and updated automatically.
	BackupDescription   BackupDescription
	SnapshotDescription SnapshotDescription
}

// export is implemented by [Backup] and [Snapshot].
type export interface {
	LogEntry() models.Export
	// ReportPlain writes a plaintext report summary into w
	ReportPlain(w io.Writer) error
	// ReportMachine writes a machine readable report summary into w
	ReportMachine(w io.Writer) error
}

// Parts lists all available snapshot parts.
func (exporter *Exporter) Parts() []string {
	return collection.MapSlice(exporter.dependencies.Snapshotable, func(c component.Snapshotable) string { return c.SnapshotName() })
}

const (
	ReportPlainPath   = "README.txt"
	ReportMachinePath = "report.json"
)

// MakeExport performs an export task as described by flags.
// Output is directed to the provided io.
func (exporter *Exporter) MakeExport(ctx context.Context, progress io.Writer, task ExportTask) (err error) {
	// extract parameters
	Title := "Backup"
	Slug := ""
	if task.Instance != nil {
		Title = "Snapshot"
		Slug = task.Instance.Slug
	}

	// determine target paths
	if _, err := logging.LogMessage(progress, "Determining target paths"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	var stagingDir, archivePath string
	if task.StagingOnly {
		stagingDir = task.Dest
	} else {
		archivePath = task.Dest
	}
	if stagingDir == "" {
		stagingDir, err = exporter.NewStagingDir(Slug)
		if err != nil {
			return err
		}
	}
	if !task.StagingOnly && archivePath == "" {
		archivePath = exporter.NewArchivePath(Slug)
	}
	_, _ = fmt.Fprintf(progress, "Staging Directory: %s\n", stagingDir)
	_, _ = fmt.Fprintf(progress, "Archive Path:      %s\n", archivePath)

	// create the staging directory
	if _, err := logging.LogMessage(progress, "Creating staging directory"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	err = umaskfree.Mkdir(stagingDir, umaskfree.DefaultDirPerm)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("failed to create staging directory: %w", err)
	}

	// if it was requested to not do staging only
	// we need the staging directory to be deleted at the end
	if !task.StagingOnly {
		defer func() {
			// #nosec G104
			logging.LogMessage(progress, "Removing staging directory") //nolint:errcheck // no way to report error
			_ = os.RemoveAll(stagingDir)
		}()
	}

	// create the actual snapshot or backup
	// write out the report
	// and retain a log entry
	var entry models.Export
	_ = logging.LogOperation(func() error {
		var export export
		if task.Instance == nil {
			task.BackupDescription.Dest = stagingDir
			backup := exporter.NewBackup(ctx, progress, task.BackupDescription)
			export = &backup
		} else {
			task.SnapshotDescription.Dest = stagingDir
			snapshot := exporter.NewSnapshot(ctx, task.Instance, progress, task.SnapshotDescription)
			export = &snapshot
		}

		// create a log entry
		entry = export.LogEntry()

		// write the machine report
		{
			reportPath := filepath.Join(stagingDir, ReportMachinePath)
			_, _ = fmt.Fprintln(progress, reportPath)

			report, err := umaskfree.Create(reportPath, umaskfree.DefaultFilePerm)
			if err != nil {
				return fmt.Errorf("failed to create report file: %w", err)
			}

			if err := export.ReportMachine(report); err != nil {
				return fmt.Errorf("failed to generate report: %w", err)
			}
		}

		// write the plaintext report
		{
			reportPath := filepath.Join(stagingDir, ReportPlainPath)
			_, _ = fmt.Fprintln(progress, reportPath)

			report, err := umaskfree.Create(reportPath, umaskfree.DefaultFilePerm)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if err := export.ReportPlain(report); err != nil {
				return fmt.Errorf("failed to generate report: %w", err)
			}
		}

		return nil
	}, progress, "Generating %s", Title)

	// if we only requested staging
	// all that is left is to write the log entry
	if task.StagingOnly {
		_, _ = fmt.Fprintln(progress, "Writing Log Entry")

		// write out the log entry
		entry.Path = stagingDir
		entry.Packed = false
		if err := exporter.dependencies.ExporterLogger.Add(ctx, entry); err != nil {
			return fmt.Errorf("failed to add log entry entry: %w", err)
		}

		if _, err := fmt.Fprintf(progress, "Wrote %s\n", stagingDir); err != nil {
			return fmt.Errorf("failed to report progress: %w", err)
		}
		return nil
	}

	if err := logging.LogOperation(func() error {
		var count int64
		defer func() { _, _ = fmt.Fprintf(progress, "Wrote %d byte(s) to %s\n", count, archivePath) }()

		st := status.NewWithCompat(progress, 1)
		st.Start()
		defer st.Stop()

		count, err = targz.Package(archivePath, stagingDir, func(dst, src string) {
			st.Set(0, dst)
		})

		if err != nil {
			return fmt.Errorf("failed to package archive: %w", err)
		}
		return nil
	}, progress, "Writing archive"); err != nil {
		return fmt.Errorf("failed to write archive: %w", err)
	}

	if _, err := logging.LogMessage(progress, "Writing Log Entry"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}

	entry.Path = archivePath
	entry.Packed = true

	if err := exporter.dependencies.ExporterLogger.Add(ctx, entry); err != nil {
		return fmt.Errorf("failed to log backup: %w", err)
	}
	return nil
}
