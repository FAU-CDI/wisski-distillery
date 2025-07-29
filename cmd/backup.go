package cmd

//spellchecker:words slog github wisski distillery internal component exporter wdlog logging goprogram exit
import (
	"fmt"
	"log/slog"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewBackupCommand() *cobra.Command {
	impl := new(backup)

	cmd := &cobra.Command{
		Use:     "backup [DEST]",
		Short:   "makes a backup of the entire distillery",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Prune, "prune", false, "prune older backup archives")
	flags.BoolVar(&impl.StagingOnly, "staging-only", false, "do not package into a backup archive, but only create a staging directory")
	flags.IntVar(&impl.ConcurrentSnapshots, "concurrent-snapshots", 2, "maximum number of concurrent snapshots")
	flags.StringVar(&impl.Positionals.Dest, "dest", "", "destination path to write backup archive to. defaults to the 'snapshots/archives/' directory")

	return cmd
}

type backup struct {
	Prune               bool
	StagingOnly         bool
	ConcurrentSnapshots int
	Positionals         struct {
		Dest string
	}
}

func (bk *backup) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		bk.Positionals.Dest = args[0]
	}
	return nil
}

var errBackupFailed = exit.NewErrorWithCode("failed to make a backup", exit.ExitGeneric)

func (bk *backup) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errBackupFailed, err)
	}

	// prune old backups
	if bk.Prune {
		defer func() {
			err := logging.LogOperation(func() error {
				return dis.Exporter().PruneExports(cmd.Context(), cmd.ErrOrStderr())
			}, cmd.ErrOrStderr(), "Pruning old backups")
			if err != nil {
				wdlog.Of(cmd.Context()).Error("failed to prune backups", slog.Any("error", err))
			}
		}()
	}

	// do the handling
	{
		err := dis.Exporter().MakeExport(cmd.Context(), cmd.ErrOrStderr(), exporter.ExportTask{
			Dest:        bk.Positionals.Dest,
			StagingOnly: bk.StagingOnly,

			Instance: nil,

			BackupDescription: exporter.BackupDescription{
				ConcurrentSnapshots: bk.ConcurrentSnapshots,
			},
		})
		if err != nil {
			return fmt.Errorf("%w: %w", errBackupFailed, err)
		}
	}

	return nil
}
