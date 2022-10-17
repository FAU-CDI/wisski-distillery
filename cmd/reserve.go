package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Reserve is the 'reserve' command
var Reserve wisski_distillery.Command = reserve{}

type reserve struct {
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"name of WissKI Instance to reserve"`
	} `positional-args:"true"`
}

func (reserve) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "reserve",
		Description: "Reserves a new WissKI Instance",
	}
}

// TODO: AfterParse to check instance!

var errReserveAlreadyExists = exit.Error{
	Message:  "Instance %q already exists",
	ExitCode: exit.ExitGeneric,
}

var errReserveGeneric = exit.Error{
	Message:  "Unable to provision instance %s: %s",
	ExitCode: exit.ExitGeneric,
}

func (r reserve) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	slug := r.Positionals.Slug

	// check that it doesn't already exist
	logging.LogMessage(context.IOStream, "Reserving new WissKI instance %s", slug)
	if exists, err := dis.Instances().Has(slug); err != nil || exists {
		return errProvisionAlreadyExists.WithMessageF(slug)
	}

	// make it in-memory
	instance, err := dis.Instances().Create(slug)
	if err != nil {
		return errProvisionGeneric.WithMessageF(slug, err)
	}

	// check that the base directory does not exist
	logging.LogMessage(context.IOStream, "Checking that base directory %s does not exist", instance.FilesystemBase)
	if fsx.IsDirectory(dis.Environment, instance.FilesystemBase) {
		return errProvisionAlreadyExists.WithMessageF(slug)
	}

	// setup docker stack
	s := instance.Reserve()
	{
		if err := logging.LogOperation(func() error {
			return s.Install(context.IOStream, component.InstallationContext{})
		}, context.IOStream, "Installing docker stack"); err != nil {
			return err
		}

		if err := logging.LogOperation(func() error {
			return s.Update(context.IOStream, true)
		}, context.IOStream, "Updating docker stack"); err != nil {
			return err
		}
	}

	// and we're done!
	logging.LogMessage(context.IOStream, "Instance has been reserved")
	context.Printf("URL:      %s\n", instance.URL().String())

	return nil
}
