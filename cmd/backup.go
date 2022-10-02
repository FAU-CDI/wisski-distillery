package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
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
	var err error

	if !bk.NoPrune {
		defer logging.LogOperation(func() error {
			return dis.PruneBackups(context.IOStream)
		}, context.IOStream, "Pruning old backups")
	}

	// determine the target path for the archive
	var sPath string
	if !bk.StagingOnly {
		// regular mode: create a temporary staging directory
		logging.LogMessage(context.IOStream, "Creating new backup staging directory")
		sPath, err = dis.SnapshotManager().NewStagingDir("")
		if err != nil {
			return errSnapshotFailed.Wrap(err)
		}
		defer func() {
			logging.LogMessage(context.IOStream, "Removing snapshot staging directory")
			dis.Environment.RemoveAll(sPath)
		}()
	} else {
		// staging mode: use dest as a destination
		sPath = bk.Positionals.Dest
		if sPath == "" {
			sPath, err = dis.SnapshotManager().NewStagingDir("")
			if err != nil {
				return errSnapshotFailed.Wrap(err)
			}
		}

		// create the directory (if it doesn't already exist)
		logging.LogMessage(context.IOStream, "Creating staging directory")
		err = dis.Core.Environment.Mkdir(sPath, environment.DefaultDirPerm)
		if !environment.IsExist(err) && err != nil {
			return errSnapshotFailed.WithMessageF(err)
		}
		err = nil
	}
	context.Println(sPath)

	logging.LogOperation(func() error {
		backup := dis.SnapshotManager().NewBackup(context.IOStream, snapshots.BackupDescription{
			Dest:                sPath,
			Auto:                bk.Positionals.Dest == "",
			ConcurrentSnapshots: bk.ConcurrentSnapshots,
		})
		backup.WriteReport(dis.Core.Environment, context.IOStream)
		return nil
	}, context.IOStream, "Generating Backup")

	// if we requested to only have a staging area, then we are done
	if bk.StagingOnly {
		context.Printf("Wrote %s\n", sPath)
		return nil
	}

	// create the archive path
	archivePath := bk.Positionals.Dest
	if archivePath == "" {
		archivePath = dis.SnapshotManager().NewArchivePath("")
	}

	// and write everything into it!
	var count int64
	if err := logging.LogOperation(func() error {
		context.IOStream.Println(archivePath)

		st := status.NewWithCompat(context.Stdout, 1)
		st.Start()
		defer st.Stop()

		count, err = targz.Package(dis.Core.Environment, archivePath, sPath, func(dst, src string) {
			st.Set(0, dst)
		})
		return err
	}, context.IOStream, "Writing backup archive"); err != nil {
		return errSnapshotFailed.Wrap(err)
	}
	context.Printf("Wrote %d byte(s) to %s\n", count, archivePath)

	return nil
}
