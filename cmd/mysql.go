package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Mysql is the 'mysql' command
var Mysql wisski_distillery.Command = mysql{}

type mysql struct {
	Positionals struct {
		Args []string `positional-arg-name:"ARGS" description:"arguments to pass to the mysql command"`
	} `positional-args:"true"`
}

func (mysql) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "mysql",
		Description: "Opens a mysql shell",
	}
}

func (ms mysql) Run(context wisski_distillery.Context) error {
	code, err := context.Environment.SQL().Shell(context.IOStream, ms.Positionals.Args...)
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
