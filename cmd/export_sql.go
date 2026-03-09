package cmd

//spellchecker:words github wisski distillery internal cobra
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
)

func NewExportSQLCommand() *cobra.Command {
	impl := new(exportSQL)

	cmd := &cobra.Command{
		Use:     "export_sql SLUG",
		Short:   "export the SQL database for a specific instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type exportSQL struct {
	Positionals struct {
		Slug string
	}
}

func (es *exportSQL) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		es.Positionals.Slug = args[0]
	}
	return nil
}

func (ets *exportSQL) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get WissKI: %w", err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), ets.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("failed to get WissKI: %w", err)
	}

	if err := instance.BoundSQL().Snapshot(cmd.Context(), cmd.ErrOrStderr(), cmd.OutOrStdout()); err != nil {
		return fmt.Errorf("failed to snapshot SQL: %w", err)
	}
	return nil
}
