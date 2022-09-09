package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/internal/sqle"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Shell is the 'shell' command
var MakeMysqlAccount wisski_distillery.Command = makeMysqlAccount{}

type makeMysqlAccount struct{}

func (makeMysqlAccount) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "make_mysql_account",
		Description: "Open a shell in the provided instance",
	}
}

var errUnableToReadUsername = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to read username: %s",
}

var errUnableToReadPassword = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to read password: %s",
}

func (mma makeMysqlAccount) Run(context wisski_distillery.Context) error {
	context.Printf("Username>")
	username, err := context.ReadLine()
	if err != nil {
		return errUnableToReadUsername.WithMessageF(err)
	}

	context.Printf("Password>")
	password, err := context.ReadPassword()
	if err != nil {
		return errUnableToReadPassword.WithMessageF(err)
	}

	query := sqle.Format("CREATE USER ?@'%' IDENTIFIED BY ?; GRANT ALL PRIVILEGES ON *.* TO ?@`%` WITH GRANT OPTION; FLUSH PRIVILEGES;", username, password, username)
	if err != nil {
		return err
	}
	code, err := context.Environment.SQL().OpenShell(context.IOStream, "-e", query)
	if err != nil {
		return err
	}

	if code != 0 {
		return exit.Error{
			ExitCode: exit.ExitCode(uint8(code)),
			Message:  fmt.Sprintf("Exit code %d", code),
		}
	}
	return nil
}
