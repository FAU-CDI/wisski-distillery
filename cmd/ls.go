package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewLsCommand() *cobra.Command {
	impl := new(ls)

	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "lists instances",
		Args:    cobra.ArbitraryArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type ls struct {
	Positionals struct {
		Slug []string
	}
}

func (l *ls) ParseArgs(cmd *cobra.Command, args []string) error {
	l.Positionals.Slug = args
	return nil
}

var errLsWissKI = exit.NewErrorWithCode("unable to get WissKIs", exit.ExitGeneric)

func (l *ls) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errLsWissKI, err)
	}

	instances, err := dis.Instances().Load(cmd.Context(), l.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: %w", errLsWissKI, err)
	}

	for _, instance := range instances {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), instance.Slug)
	}

	return nil
}
