package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"bufio"
	"fmt"
	"strings"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewPurgeCommand() *cobra.Command {
	impl := new(purge)

	cmd := &cobra.Command{
		Use:     "purge",
		Short:   "purges an instance",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.Yes, "yes", false, "do not ask for confirmation")

	return cmd
}

type purge struct {
	Yes         bool
	Positionals struct {
		Slug string
	}
}

func (p *purge) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		p.Positionals.Slug = args[0]
	}
	return nil
}

func (*purge) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "purge",
		Description: "purges an instance",
	}
}

var (
	errPurgeNoConfirmation = exit.NewErrorWithCode("aborting after request was not confirmed. either type `yes` or pass `--yes` on the command line", exit.ExitGeneric)
	errPurgeFailed         = exit.NewErrorWithCode("failed to run purge", exit.ExitGeneric)
)

func (p *purge) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errPurgeFailed, err)
	}

	slug := p.Positionals.Slug

	// check the confirmation from the user
	if !p.Yes {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "About to remove instance %q. This cannot be undone.\n", slug)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Type 'yes' to continue: ")
		reader := bufio.NewReader(cmd.InOrStdin())
		line, err := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(line) != "yes" {
			return errPurgeNoConfirmation
		}
	}

	// do the purge!
	if err := dis.Purger().Purge(cmd.Context(), cmd.OutOrStdout(), slug); err != nil {
		return fmt.Errorf("%w: %w", errPurgeFailed, err)
	}
	return nil
}
