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
	"github.com/tkw1536/pkglib/collection"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/status"
)

// ExportTask describes a task that makes either a [Backup] or a [Snapshot].
// See [Exporter.MakeExport]
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

// export is implemented by [Backup] and [Snapshot]
type export interface {
	LogEntry() models.Export
	// ReportPlain writes a plaintext report summary into w
	ReportPlain(w io.Writer) error
	// ReportMachine writes a machine readable report summary into w
	ReportMachine(w io.Writer) error
}

// Parts lists all available snapshot parts
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
	logging.LogMessage(progress, "Determining target paths")
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
	fmt.Fprintf(progress, "Staging Directory: %s\n", stagingDir)
	fmt.Fprintf(progress, "Archive Path:      %s\n", archivePath)

	// create the staging directory
	logging.LogMessage(progress, "Creating staging directory")
	err = umaskfree.Mkdir(stagingDir, umaskfree.DefaultDirPerm)
	if !errors.Is(err, fs.ErrExist) && err != nil {
		return err
	}

	// if it was requested to not do staging only
	// we need the staging directory to be deleted at the end
	if !task.StagingOnly {
		defer func() {
			logging.LogMessage(progress, "Removing staging directory")
			os.RemoveAll(stagingDir)
		}()
	}

	// create the actual snapshot or backup
	// write out the report
	// and retain a log entry
	var entry models.Export
	logging.LogOperation(func() error {
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
			fmt.Fprintln(progress, reportPath)

			report, err := umaskfree.Create(reportPath, umaskfree.DefaultFilePerm)
			if err != nil {
				return err
			}

			if err := export.ReportMachine(report); err != nil {
				return err
			}
		}

		// write the plaintext report
		{
			reportPath := filepath.Join(stagingDir, ReportPlainPath)
			fmt.Fprintln(progress, reportPath)

			report, err := umaskfree.Create(reportPath, umaskfree.DefaultFilePerm)
			if err != nil {
				return err
			}

			if err := export.ReportPlain(report); err != nil {
				return err
			}
		}

		return nil
	}, progress, "Generating %s", Title)

	// if we only requested staging
	// all that is left is to write the log entry
	if task.StagingOnly {
		fmt.Fprintln(progress, "Writing Log Entry")

		// write out the log entry
		entry.Path = stagingDir
		entry.Packed = false
		if err := exporter.dependencies.ExporterLogger.Add(ctx, entry); err != nil {
			return err
		}

		fmt.Fprintf(progress, "Wrote %s\n", stagingDir)
		return nil
	}

	// package everything up as an archive!
	if err := logging.LogOperation(func() error {
		var count int64
		defer func() { fmt.Fprintf(progress, "Wrote %d byte(s) to %s\n", count, archivePath) }()

		st := status.NewWithCompat(progress, 1)
		st.Start()
		defer st.Stop()

		count, err = targz.Package(archivePath, stagingDir, func(dst, src string) {
			st.Set(0, dst)
		})

		return err
	}, progress, "Writing archive"); err != nil {
		return err
	}

	// write out the log entry
	_, _ = logging.LogMessage(progress, "Writing Log Entry") // shouldn't fail because of log
	entry.Path = archivePath
	entry.Packed = true
	return exporter.dependencies.ExporterLogger.Add(ctx, entry)
}
