package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"go.tkw01536.de/goprogram/exit"
)

// Ls is the 'ls' command.
var Ls wisski_distillery.Command = ls{}

type ls struct {
	Positionals struct {
		Slug []string `description:"slugs of instances to list. if empty, list all instances" positional-arg-name:"SLUG" required:"0"`
	} `positional-args:"true"`
}

func (ls) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "ls",
		Description: "lists instances",
	}
}

var errLsWissKI = exit.NewErrorWithCode("unable to get WissKIs", exit.ExitGeneric)

func (l ls) Run(context wisski_distillery.Context) error {
	instances, err := context.Environment.Instances().Load(context.Context, l.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: %w", errLsWissKI, err)
	}

	for _, instance := range instances {
		_, _ = context.Println(instance.Slug)
	}

	return nil
}
