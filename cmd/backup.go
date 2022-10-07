package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component/snapshots"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
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
	err := dis.SnapshotManager().HandleSnapshotLike(context.IOStream, snapshots.SnapshotFlags{
		Dest:        bk.Positionals.Dest,
		Slug:        "",
		Title:       "Backup",
		StagingOnly: bk.StagingOnly,

		Do: func(dest string) snapshots.SnapshotLike {
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
