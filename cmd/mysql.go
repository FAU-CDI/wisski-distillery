package cmd

//spellchecker:words github wisski distillery internal cobra pkglib exit
import (
	"fmt"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
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

	flags := cmd.Flags()
	flags.StringVar(&impl.Slug, "slug", "", "Optional slug to open instance shell for. If not provided, gives a shell in the global sql database.")

	return cmd
}

type mysql struct {
	Positionals struct {
		Args []string
	}
	Slug string
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

	if ms.Slug == "" {
		return ms.globalShell(cmd, dis)
	}

	return ms.instanceShell(cmd, dis)
}

func (ms *mysql) globalShell(cmd *cobra.Command, dis *dis.Distillery) error {
	code := dis.SQL().DeprecatedShell(cmd.Context(), streamFromCommand(cmd), ms.Positionals.Args...)

	if code := exit.Code(code); code != 0 {
		return exit.NewErrorWithCode(fmt.Sprintf("exit code %d", code), code)
	}
	return nil
}

func (ms *mysql) instanceShell(cmd *cobra.Command, dis *dis.Distillery) error {
	instance, err := dis.Instances().WissKI(cmd.Context(), ms.Slug)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	code := instance.DelegatedSQL().Shell(cmd.Context(), streamFromCommand(cmd), ms.Positionals.Args...)
	if code := exit.Code(code); code != 0 {
		return exit.NewErrorWithCode(fmt.Sprintf("exit code %d", code), code)
	}
	return nil
}
