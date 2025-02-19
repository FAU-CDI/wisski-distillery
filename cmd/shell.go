package cmd

//spellchecker:words github wisski distillery internal ingredient barrel goprogram exit parser
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Shell is the 'shell' command
var Shell wisski_distillery.Command = shell{}

type shell struct {
	Positionals struct {
		Slug string   `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to run shell in"`
		Args []string `positional-arg-name:"ARGS" description:"arguments to pass to the shell"`
	} `positional-args:"true"`
}

func (shell) Description() wisski_distillery.Description {
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

var errShellWissKI = exit.Error{
	Message:  "unable to find WissKI",
	ExitCode: exit.ExitGeneric,
}

func (sh shell) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, sh.Positionals.Slug)
	if err != nil {
		return errShellWissKI.WrapError(err)
	}

	{
		err := instance.Barrel().Shell(context.Context, context.IOStream, sh.Positionals.Args...)
		if err != nil {
			code := err.(barrel.ExitError).Code()
			return exit.Error{
				ExitCode: code,
				Message:  fmt.Sprintf("Exit code %d", code),
			}
		}
	}

	return nil
}
