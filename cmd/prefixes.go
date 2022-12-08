package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Prefixes is then 'prefixes' command
var Prefixes wisski_distillery.Command = prefixes{}

type prefixes struct {
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to show prefixes for"`
	} `positional-args:"true"`
}

func (prefixes) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "prefixes",
		Description: "List all Prefixes for a specific WissKI",
	}
}

var errPrefixesGeneric = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to load prefixes",
}

func (p prefixes) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, p.Positionals.Slug)
	if err != nil {
		return err
	}

	prefixes, err := instance.Prefixes().All(context.Context, nil)
	if err != nil {
		return errPrefixesGeneric.Wrap(err)
	}

	for _, p := range prefixes {
		context.Println(p)
	}

	return nil
}
