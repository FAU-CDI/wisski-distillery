package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Config is the configuration command
var Config wisski_distillery.Command = cfg{}

type cfg struct{}

func (c cfg) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "config",
		Description: "prints information about configuration",
	}
}

var errMarshalConfig = exit.Error{
	Message:  "unable to marshal config",
	ExitCode: exit.ExitGeneric,
}

func (cfg) Run(context wisski_distillery.Context) error {
	if err := context.Environment.Config.Marshal(context.Stdout); err != nil {
		return errMarshalConfig.Wrap(err)
	}
	return nil
}
