package cmd

//spellchecker:words github wisski distillery internal goprogram exit parser
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Mysql is the 'mysql' command.
var Mysql wisski_distillery.Command = mysql{}

type mysql struct {
	Positionals struct {
		Args []string `positional-arg-name:"ARGS" description:"arguments to pass to the mysql command"`
	} `positional-args:"true"`
}

func (mysql) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "mysql",
		Description: "opens a mysql shell",
	}
}

func (ms mysql) Run(context wisski_distillery.Context) error {
	code := context.Environment.SQL().Shell(context.Context, context.IOStream, ms.Positionals.Args...)

	if code := exit.Code(code); code != 0 {
		return exit.Error{
			ExitCode: code,
			Message:  fmt.Sprintf("Exit code %d", code),
		}
	}
	return nil
}
