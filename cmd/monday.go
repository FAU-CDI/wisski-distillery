package cmd

import (
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
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
		Requirements: env.Requirements{
			NeedsDistillery: true,
		},
		Command:     "monday",
		Description: "Runs regular monday tasks",
	}
}

func (monday monday) AfterParse() error {
	_, err := os.Stat(monday.Positionals.GraphdbZip)
	if os.IsNotExist(err) {
		return errNoGraphDBZip.WithMessageF(monday.Positionals.GraphdbZip)
	}
	if err != nil {
		return err
	}
	return nil
}

func (monday monday) Run(context wisski_distillery.Context) error {
	if err := logging.LogOperation(func() error {
		return context.Exec("backup")
	}, context.IOStream, "Running backup"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("system_update", monday.Positionals.GraphdbZip)
	}, context.IOStream, "Running system_update"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		return context.Exec("rebuild")
	}, context.IOStream, "Running rebuild"); err != nil {
		return err
	}

	if monday.UpdateInstances {
		if err := logging.LogOperation(func() error {
			return context.Exec("blind_update")
		}, context.IOStream, "Running blind_update"); err != nil {
			return err
		}
	}

	logging.LogMessage(context.IOStream, "Done, have a great week!")
	return nil
}
