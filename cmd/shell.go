package cmd

//spellchecker:words errors github wisski distillery internal ingredient barrel goprogram exit parser
import (
	"errors"
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/spf13/cobra"
	"go.tkw01536.de/goprogram/parser"
	"go.tkw01536.de/pkglib/exit"
)

func NewShellCommand() *cobra.Command {
	impl := new(shell)

	cmd := &cobra.Command{
		Use:     "shell SLUG [ARGS...]",
		Short:   "open a shell in the provided instance",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type shell struct {
	Positionals struct {
		Slug string
		Args []string
	}
}

func (sh *shell) ParseArgs(cmd *cobra.Command, args []string) error {
	sh.Positionals.Slug = args[0]
	if len(args) >= 2 {
		sh.Positionals.Args = args[1:]
	}
	return nil
}

func (*shell) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "shell",
		Description: "open a shell in the provided instance",
	}
}

var errShellWissKI = exit.NewErrorWithCode("unable to find WissKI", exit.ExitGeneric)

func (sh *shell) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errShellWissKI, err)
	}

	instance, err := dis.Instances().WissKI(cmd.Context(), sh.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errShellWissKI, err)
	}

	{
		args := append([]string{"/bin/bash"}, sh.Positionals.Args...)
		err := instance.Barrel().BashScript(cmd.Context(), streamFromCommand(cmd), args...)
		if err != nil {
			var ee barrel.ExitError
			if !(errors.As(err, &ee)) {
				return fmt.Errorf("barrel.Shell returned unexpected error: %w", err)
			}
			code := ee.Code()

			return exit.NewErrorWithCode(fmt.Sprintf("exit code %d", code), code)
		}
	}

	return nil
}
