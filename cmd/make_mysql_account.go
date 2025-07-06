package cmd

//spellchecker:words github wisski distillery internal goprogram exit parser
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/goprogram/parser"
)

// Shell is the 'shell' command.
var MakeMysqlAccount wisski_distillery.Command = makeMysqlAccount{}

type makeMysqlAccount struct{}

func (makeMysqlAccount) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "make_mysql_account",
		Description: "opens a shell in the provided instance",
	}
}

var (
	errUnableToReadUsername = exit.NewErrorWithCode("unable to read username", exit.ExitGeneric)
	errUnableToReadPassword = exit.NewErrorWithCode("unable to read password", exit.ExitGeneric)
	errUnableToMakeAccount  = exit.NewErrorWithCode("unable to create account", exit.ExitGeneric)
)

func (mma makeMysqlAccount) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	_, _ = context.Printf("Username>")
	username, err := context.ReadLine()
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToReadUsername, err)
	}

	_, _ = context.Printf("Password>")
	password, err := context.ReadPassword()
	if err != nil {
		return fmt.Errorf("%w: %w", errUnableToReadPassword, err)
	}

	if err := dis.SQL().CreateSuperuser(context.Context, username, password, false); err != nil {
		return fmt.Errorf("%w: %w", errUnableToMakeAccount, err)
	}

	return nil
}
