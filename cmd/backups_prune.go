package cmd

//spellchecker:words slog github wisski distillery internal component exporter wdlog logging goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// backup is the 'backups_prune' command.
var BackupsPrune wisski_distillery.Command = backupsPrune{}

type backupsPrune struct{}

func (backupsPrune) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "backups_prune",
		Description: "prunes old backup archives",
	}
}

var errPruneFailed = exit.NewErrorWithCode("failed to prune backups", exit.ExitGeneric)

func (backupsPrune) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	if err := dis.Exporter().PruneExports(context.Context, context.Stderr); err != nil {
		return fmt.Errorf("%w: %w", errPruneFailed, err)
	}
	return nil
}
