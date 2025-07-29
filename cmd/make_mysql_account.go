package cmd

//spellchecker:words strings github wisski distillery internal cobra pkglib exit nobufio
import (
	"fmt"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/nobufio"
)

func NewMakeMysqlAccountCommand() *cobra.Command {
	impl := new(makeMysqlAccount)

	cmd := &cobra.Command{
		Use:   "make_mysql_account",
		Short: "creates a MySQL account",
		Args:  cobra.NoArgs,
		RunE:  impl.Exec,
	}

	return cmd
}

type makeMysqlAccount struct{}

var (
	errUnableToReadUsername = exit.NewErrorWithCode("unable to read username", exit.ExitGeneric)
	errUnableToReadPassword = exit.NewErrorWithCode("unable to read password", exit.ExitGeneric)
	errUnableToMakeAccount  = exit.NewErrorWithCode("unable to create account", exit.ExitGeneric)
)

func (mma *makeMysqlAccount) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToMakeAccount, err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Username>")
	username, err := nobufio.ReadLine(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToReadUsername, err)
	}
	username = strings.TrimSpace(username)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Password>")
	password, err := nobufio.ReadPassword(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToReadPassword, err)
	}

	if err := dis.SQL().CreateSuperuser(cmd.Context(), username, password, false); err != nil {
		return fmt.Errorf("%w: %w", errUnableToMakeAccount, err)
	}

	return nil
}
