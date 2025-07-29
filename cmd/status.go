package cmd

//spellchecker:words encoding json github wisski distillery internal goprogram exit
import (
	"encoding/json"
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewStatusCommand() *cobra.Command {
	impl := new(cStatus)

	cmd := &cobra.Command{
		Use:     "status",
		Short:   "provide information about the distillery as a whole",
		Args:    cobra.NoArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.JSON, "json", false, "print status as JSON instead of as string")

	return cmd
}

type cStatus struct {
	JSON bool
}

func (s *cStatus) ParseArgs(cmd *cobra.Command, args []string) error {
	return nil
}

var errStatusGeneric = exit.NewErrorWithCode("unable to get status", exit.ExitGeneric)

func (s *cStatus) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errStatusGeneric, err)
	}

	status, _, err := dis.Info().Status(cmd.Context(), true)
	if err != nil {
		return fmt.Errorf("%w: %w", errStatusGeneric, err)
	}

	if s.JSON {
		err := json.NewEncoder(cmd.OutOrStdout()).Encode(status)
		if err != nil {
			return fmt.Errorf("%w: %w", errStatusGeneric, err)
		}
		return nil
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Total Instances:      %v\n", status.TotalCount)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "      (running):      %v\n", status.RunningCount)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "      (stopped):      %v\n", status.StoppedCount)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Backups: (count %d)\n", len(status.Backups))
	for _, s := range status.Backups {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s (slug %q, taken %s, packed %v)\n", s.Path, s.Slug, s.Created.String(), s.Packed)
	}

	return nil
}
