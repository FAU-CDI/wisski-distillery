package cmd

//spellchecker:words github wisski distillery internal goprogram exit parser
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewMysqlCommand() *cobra.Command {
	impl := new(mysql)

	cmd := &cobra.Command{
		Use:     "mysql [ARGS...]",
		Short:   "opens a mysql shell",
		Args:    cobra.ArbitraryArgs,
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type mysql struct {
	Positionals struct {
		Args []string
	}
}

func (ms *mysql) ParseArgs(cmd *cobra.Command, args []string) error {
	ms.Positionals.Args = args
	return nil
}

func (ms *mysql) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("failed to get distillery: %w", err)
	}

	code := dis.SQL().Shell(cmd.Context(), streamFromCommand(cmd), ms.Positionals.Args...)

	if code := exit.Code(code); code != 0 {
		return exit.NewErrorWithCode(fmt.Sprintf("exit code %d", code), code)
	}
	return nil
}
