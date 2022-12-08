package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Pathbuilders is the 'pathbuilders' command
var Pathbuilders wisski_distillery.Command = pathbuilders{}

type pathbuilders struct {
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to export pathbuilders of"`
		Name string `positional-arg-name:"NAME" description:"name of pathbuilder to get. if omitted, show a list of all pathbuilders"`
	} `positional-args:"true"`
}

func (pathbuilders) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "pathbuilder",
		Description: "list pathbuilders of a specific instance",
	}
}

var errPathbuilders = exit.Error{
	Message:  "unable to export pathbuilder: %s",
	ExitCode: exit.ExitGeneric,
}

var errNoPathbuilder = exit.Error{
	Message:  "pathbuilder %q does not exist",
	ExitCode: exit.ExitGeneric,
}

func (pb pathbuilders) Run(context wisski_distillery.Context) error {

	// get the wisski
	instance, err := context.Environment.Instances().WissKI(context.Context, pb.Positionals.Slug)
	if err != nil {
		return err
	}

	// get all of the pathbuilders
	if pb.Positionals.Name == "" {
		names, err := instance.Pathbuilder().All(context.Context, nil)
		if err != nil {
			return errPathbuilders.WithMessageF(err)
		}
		for _, name := range names {
			context.Println(name)
		}
		return nil
	}

	// get all the pathbuilders
	xml, err := instance.Pathbuilder().Get(context.Context, nil, pb.Positionals.Name)
	if xml == "" {
		return errNoPathbuilder.WithMessageF(pb.Positionals.Name)
	}
	if err != nil {
		return errPathbuilders.WithMessageF(err)
	}
	context.Printf("%s", xml)

	return nil
}
