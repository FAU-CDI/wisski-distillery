package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
)

// Monday is the 'monday' command
var Monday wisski_distillery.Command = monday{}

type monday struct {
	UpdateInstances bool `short:"u" long:"update-instances" description:"Fully update instances. May take a long time, and is potentially breaking. "`
	Positionals     struct {
		GraphdbZip string `positional-arg-name:"PATH_TO_GRAPHDB_ZIP" required:"1-1" description:"path to the graphdb.zip file"`
	} `positional-args:"true"`
}

func (monday) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "monday",
		Description: "Runs regular monday tasks",
	}
}

func (monday monday) AfterParse() error {
	// TODO: Use a generic environment here!
	if !fsx.IsFile(new(environment.Native), monday.Positionals.GraphdbZip) {
		return errNoGraphDBZip.WithMessageF(monday.Positionals.GraphdbZip)
	}
	return nil
}

func (monday monday) Run(context wisski_distillery.Context) error {
	if err := logging.LogOperation(func() error {
		return context.Exec("backup")
	}, context.Stderr, context.Context, "Running backup"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("system_update", monday.Positionals.GraphdbZip)
	}, context.Stderr, context.Context, "Running system_update"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("rebuild")
	}, context.Stderr, context.Context, "Running rebuild"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("update_prefix_config")
	}, context.Stderr, context.Context, "Running update_prefix_config"); err != nil {
		return err
	}

	if monday.UpdateInstances {
		if err := logging.LogOperation(func() error {
			return context.Exec("blind_update")
		}, context.Stderr, context.Context, "Running blind_update"); err != nil {
			return err
		}
	}

	logging.LogMessage(context.Stderr, context.Context, "Done, have a great week!")
	return nil
}
