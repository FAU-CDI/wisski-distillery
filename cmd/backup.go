package cmd

//spellchecker:words slog github wisski distillery internal component exporter wdlog logging goprogram exit
import (
	"log/slog"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Backup is the 'backup' command.
var Backup wisski_distillery.Command = backup{}

type backup struct {
	NoPrune             bool `short:"n" long:"no-prune" description:"do not prune older backup archives"`
	StagingOnly         bool `short:"s" long:"staging-only" description:"do not package into a backup archive, but only create a staging directory"`
	ConcurrentSnapshots int  `short:"c" long:"concurrent-snapshots" description:"maximum number of concurrent snapshots" default:"2"`
	Positionals         struct {
		Dest string "positional-arg-name:\"DEST\" description:\"destination path to write backup archive to. defaults to the `snapshots/archives/` directory\""
	} `positional-args:"true"`
}

func (backup) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "backup",
		Description: "makes a backup of the entire distillery",
	}
}

var errBackupFailed = exit.Error{
	Message:  "failed to make a backup",
	ExitCode: exit.ExitGeneric,
}

func (bk backup) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// prune old backups
	if !bk.NoPrune {
		defer func() {
			err := logging.LogOperation(func() error {
				return dis.Exporter().PruneExports(context.Context, context.Stderr)
			}, context.Stderr, "Pruning old backups")
			if err != nil {
				wdlog.Of(context.Context).Error("failed to prune backups", slog.Any("error", err))
			}
		}()
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
		return errBackupFailed.WrapError(err)
	}
	return nil
}
