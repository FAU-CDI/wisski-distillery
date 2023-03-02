package exporter

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/targz"
	"github.com/tkw1536/goprogram/status"
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
	Report(w io.Writer) (int, error)
}

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
	logging.LogMessage(progress, ctx, "Determining target paths")
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
	logging.ProgressF(progress, ctx, "Staging Directory: %s\n", stagingDir)
	logging.ProgressF(progress, ctx, "Archive Path:      %s\n", archivePath)

	// create the staging directory
	logging.LogMessage(progress, ctx, "Creating staging directory")
	err = fsx.Mkdir(stagingDir, fsx.DefaultDirPerm)
	if !os.IsExist(err) && err != nil {
		return err
	}

	// if it was requested to not do staging only
	// we need the staging directory to be deleted at the end
	if !task.StagingOnly {
		defer func() {
			logging.LogMessage(progress, ctx, "Removing staging directory")
			os.RemoveAll(stagingDir)
		}()
	}

	// create the actual snapshot or backup
	// write out the report
	// and retain a log entry
	var entry models.Export
	logging.LogOperation(func() error {
		var sl export
		if task.Instance == nil {
			task.BackupDescription.Dest = stagingDir
			backup := exporter.NewBackup(ctx, progress, task.BackupDescription)
			sl = &backup
		} else {
			task.SnapshotDescription.Dest = stagingDir
			snapshot := exporter.NewSnapshot(ctx, task.Instance, progress, task.SnapshotDescription)
			sl = &snapshot
		}

		// create a log entry
		entry = sl.LogEntry()

		// find the report path
		reportPath := filepath.Join(stagingDir, "report.txt")
		logging.ProgressF(progress, ctx, reportPath)

		// create the path
		report, err := fsx.Create(reportPath, fsx.DefaultFilePerm)
		if err != nil {
			return err
		}

		// and write out the report
		{
			_, err := sl.Report(report)
			return err
		}
	}, progress, ctx, "Generating %s", Title)

	// if we only requested staging
	// all that is left is to write the log entry
	if task.StagingOnly {
		logging.LogMessage(progress, ctx, "Writing Log Entry")

		// write out the log entry
		entry.Path = stagingDir
		entry.Packed = false
		exporter.Dependencies.ExporterLogger.Add(ctx, entry)

		logging.ProgressF(progress, ctx, "Wrote %s\n", stagingDir)
		return nil
	}

	// package everything up as an archive!
	if err := logging.LogOperation(func() error {
		var count int64
		defer func() { logging.ProgressF(progress, ctx, "Wrote %d byte(s) to %s\n", count, archivePath) }()

		st := status.NewWithCompat(progress, 1)
		st.Start()
		defer st.Stop()

		count, err = targz.Package(exporter.Environment, archivePath, stagingDir, func(dst, src string) {
			st.Set(0, dst)
		})

		return err
	}, progress, ctx, "Writing archive"); err != nil {
		return err
	}

	// write out the log entry
	logging.LogMessage(progress, ctx, "Writing Log Entry")
	entry.Path = archivePath
	entry.Packed = true
	exporter.Dependencies.ExporterLogger.Add(ctx, entry)

	// and we're done!
	return nil
}
