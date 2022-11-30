package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Backup is the 'backup' command
var Backup wisski_distillery.Command = backup{}

type backup struct {
	NoPrune             bool `short:"n" long:"no-prune" description:"Do not prune older backup archives"`
	StagingOnly         bool `short:"s" long:"staging-only" description:"Do not package into a backup archive, but only create a staging directory"`
	ConcurrentSnapshots int  `short:"c" long:"concurrent-snapshots" description:"Maximum number of concurrent snapshots" default:"2"`
	Positionals         struct {
		Dest string `positional-arg-name:"DEST" description:"Destination path to write backup archive to. Defaults to the snapshots/archives/ directory"`
	} `positional-args:"true"`
}

func (backup) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
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

func (bk backup) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// prune old backups
	if !bk.NoPrune {
		defer logging.LogOperation(func() error {
			return dis.Exporter().PruneExports(context.Context, context.Stderr)
		}, context.Stderr, "Pruning old backups")
	}

	// do the handling
	err := dis.Exporter().MakeExport(context.Context, context.Stderr, exporter.ExportTask{
		Dest:        bk.Positionals.Dest,
		StagingOnly: bk.StagingOnly,

		Instance: nil,

		BackupDescription: exporter.BackupDescription{
			ConcurrentSnapshots: bk.ConcurrentSnapshots,
		},
	})

	if err != nil {
		return errBackupFailed.Wrap(err)
	}
	return nil
}
