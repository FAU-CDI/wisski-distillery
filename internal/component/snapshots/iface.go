package snapshots

import (
	"io"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/targz"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

type SnapshotFlags struct {
	Dest        string
	Slug        string
	Title       string // "Backup" or "Snapshot"
	StagingOnly bool

	Do func(dest string) SnapshotLike
}

type SnapshotLike interface {
	LogEntry() models.Snapshot
	Report(w io.Writer) (int, error)
}

func (manager *Manager) HandleSnapshotLike(context stream.IOStream, flags SnapshotFlags) (err error) {
	// determine target paths
	logging.LogMessage(context, "Determining target paths")
	var stagingDir, archivePath string
	if flags.StagingOnly {
		stagingDir = flags.Dest
	} else {
		archivePath = flags.Dest
	}
	if stagingDir == "" {
		stagingDir, err = manager.NewStagingDir(flags.Slug)
		if err != nil {
			return err
		}
	}
	if !flags.StagingOnly && archivePath == "" {
		archivePath = manager.NewArchivePath(flags.Slug)
	}
	context.Printf("Staging Directory: %s\n", stagingDir)
	context.Printf("Archive Path:      %s\n", archivePath)

	// create the staging directory
	logging.LogMessage(context, "Creating staging directory")
	err = manager.Environment.Mkdir(stagingDir, environment.DefaultDirPerm)
	if !environment.IsExist(err) && err != nil {
		return err
	}

	// if it was requested to not do staging only
	// we need the staging directory to be deleted at the end
	if !flags.StagingOnly {
		defer func() {
			logging.LogMessage(context, "Removing staging directory")
			manager.Environment.RemoveAll(stagingDir)
		}()
	}

	// create the actual snapshot or backup
	// write out the report
	// and retain a log entry
	var entry models.Snapshot
	logging.LogOperation(func() error {
		// do the snapshot!
		sl := flags.Do(stagingDir)

		// create a log entry
		entry = sl.LogEntry()

		// find the report path
		reportPath := filepath.Join(stagingDir, "report.txt")
		context.Println(reportPath)

		// create the path
		report, err := manager.Environment.Create(reportPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}

		// and write out the report
		{
			_, err := sl.Report(report)
			return err
		}
	}, context, "Generating %s", flags.Title)

	// if we only requested staging
	// all that is left is to write the log entry
	if flags.StagingOnly {
		logging.LogMessage(context, "Writing Log Entry")

		// write out the log entry
		entry.Path = stagingDir
		entry.Packed = false
		manager.Instances.AddSnapshotLog(entry)

		context.Printf("Wrote %s\n", stagingDir)
		return nil
	}

	// package everything up as an archive!
	if err := logging.LogOperation(func() error {
		var count int64
		defer func() { context.Printf("Wrote %d byte(s) to %s\n", count, archivePath) }()

		st := status.NewWithCompat(context.Stdout, 1)
		st.Start()
		defer st.Stop()

		count, err = targz.Package(manager.Environment, archivePath, stagingDir, func(dst, src string) {
			st.Set(0, dst)
		})

		return err
	}, context, "Writing archive"); err != nil {
		return err
	}

	// write out the log entry
	logging.LogMessage(context, "Writing Log Entry")
	entry.Path = archivePath
	entry.Packed = true
	manager.Instances.AddSnapshotLog(entry)

	// and we're done!
	return nil
}
