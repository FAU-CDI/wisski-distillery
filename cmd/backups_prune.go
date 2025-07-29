package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewBackupsPruneCommand() *cobra.Command {
	impl := new(backupsPrune)

	cmd := &cobra.Command{
		Use:     "backups_prune",
		Short:   "prunes old backup archives",
		Args:    cobra.NoArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type backupsPrune struct{}

func (bp *backupsPrune) ParseArgs(cmd *cobra.Command, args []string) error {
	return nil
}

var errPruneFailed = exit.NewErrorWithCode("failed to prune backups", exit.ExitGeneric)

func (bp *backupsPrune) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errPruneFailed, err)
	}

	if err := dis.Exporter().PruneExports(cmd.Context(), cmd.ErrOrStderr()); err != nil {
		return fmt.Errorf("%w: %w", errPruneFailed, err)
	}
	return nil
}
