package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Ls is the 'ls' command
var Ls wisski_distillery.Command = ls{}

type ls struct {
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slugs of instances to list. if empty, list all instances"`
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

var errLsWissKI = exit.Error{
	Message:  "unable to get WissKIs",
	ExitCode: exit.ExitGeneric,
}

func (l ls) Run(context wisski_distillery.Context) error {
	instances, err := context.Environment.Instances().Load(context.Context, l.Positionals.Slug...)
	if err != nil {
		return errLsWissKI.WrapError(err)
	}

	for _, instance := range instances {
		context.Println(instance.Slug)
	}

	return nil
}
