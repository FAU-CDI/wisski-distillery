package cmd

//spellchecker:words github wisski distillery internal goprogram exit parser
import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Shell is the 'shell' command
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

var errUnableToReadUsername = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to read username",
}

var errUnableToReadPassword = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to read password",
}

var errUnableToMakeAccount = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to create account",
}

func (mma makeMysqlAccount) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	context.Printf("Username>")
	username, err := context.ReadLine()
	if err != nil {
		return errUnableToReadUsername.WrapError(err)
	}

	context.Printf("Password>")
	password, err := context.ReadPassword()
	if err != nil {
		return errUnableToReadPassword.WrapError(err)
	}

	if err := dis.SQL().CreateSuperuser(context.Context, username, password, false); err != nil {
		return errUnableToMakeAccount.WrapError(err)
	}

	return nil
}
