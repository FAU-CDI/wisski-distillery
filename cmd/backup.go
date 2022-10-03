package cmd

import (
	"io"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/FAU-CDI/wisski-distillery/pkg/targz"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/status"
)

// Backup is the 'backup' command
var Backup wisski_distillery.Command = backupC{}

type backupC struct {
	NoPrune             bool `short:"n" long:"no-prune" description:"Do not prune older backup archives"`
	StagingOnly         bool `short:"s" long:"staging-only" description:"Do not package into a backup archive, but only create a staging directory"`
	ConcurrentSnapshots int  `short:"c" long:"concurrent-snapshots" description:"Maximum number of concurrent snapshots" default:"2"`
	Positionals         struct {
		Dest string `positional-arg-name:"DEST" description:"Destination path to write backup archive to. Defaults to the snapshots/archives/ directory"`
	} `positional-args:"true"`
}

func (backupC) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "backup",
		Description: "Makes a backup of the entire distillery",
	}
}

var errBackupFailed = exit.Error{
	Message:  "Failed to make a backup",
	ExitCode: exit.ExitGeneric,
}

func (bk backupC) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// prune old backups
	if !bk.NoPrune {
		defer logging.LogOperation(func() error {
			return dis.SnapshotManager().PruneBackups(context.IOStream)
		}, context.IOStream, "Pruning old backups")
	}

	// do the handling
	err := handleSnapshotLike(context, SnapshotFlags{
		Dest:        bk.Positionals.Dest,
		Slug:        "",
		Title:       "Backup",
		StagingOnly: bk.StagingOnly,

		Do: func(dest string) SnapshotLike {
			backup := dis.SnapshotManager().NewBackup(context.IOStream, snapshots.BackupDescription{
				Dest:                dest,
				ConcurrentSnapshots: bk.ConcurrentSnapshots,
			})
			return &backup
		},
	})

	if err != nil {
		return errBackupFailed.Wrap(err)
	}
	return nil
}

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

func handleSnapshotLike(context wisski_distillery.Context, flags SnapshotFlags) (err error) {
	dis := context.Environment

	// determine target paths
	logging.LogMessage(context.IOStream, "Determining target paths")
	var stagingDir, archivePath string
	if flags.StagingOnly {
		stagingDir = flags.Dest
	} else {
		archivePath = flags.Dest
	}
	if stagingDir == "" {
		stagingDir, err = dis.SnapshotManager().NewStagingDir(flags.Slug)
		if err != nil {
			return err
		}
	}
	if !flags.StagingOnly && archivePath == "" {
		archivePath = dis.SnapshotManager().NewArchivePath(flags.Slug)
	}
	context.Printf("Staging Directory: %s\n", stagingDir)
	context.Printf("Archive Path:      %s\n", archivePath)

	// create the staging directory
	logging.LogMessage(context.IOStream, "Creating staging directory")
	err = dis.Core.Environment.Mkdir(stagingDir, environment.DefaultDirPerm)
	if !environment.IsExist(err) && err != nil {
		return err
	}

	// if it was requested to not do staging only
	// we need the staging directory to be deleted at the end
	if !flags.StagingOnly {
		defer func() {
			logging.LogMessage(context.IOStream, "Removing staging directory")
			dis.Environment.RemoveAll(stagingDir)
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
		report, err := dis.Environment.Create(reportPath, environment.DefaultFilePerm)
		if err != nil {
			return err
		}

		// and write out the report
		{
			_, err := sl.Report(report)
			return err
		}
	}, context.IOStream, "Generating %s", flags.Title)

	// if we only requested staging
	// all that is left is to write the log entry
	if flags.StagingOnly {
		logging.LogMessage(context.IOStream, "Writing Log Entry")

		// write out the log entry
		entry.Path = stagingDir
		entry.Packed = false
		dis.Instances().AddSnapshotLog(entry)

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

		count, err = targz.Package(dis.Core.Environment, archivePath, stagingDir, func(dst, src string) {
			st.Set(0, dst)
		})

		return err
	}, context.IOStream, "Writing archive"); err != nil {
		return err
	}

	// write out the log entry
	logging.LogMessage(context.IOStream, "Writing Log Entry")
	entry.Path = archivePath
	entry.Packed = true
	dis.Instances().AddSnapshotLog(entry)

	// and we're done!
	return nil
}
