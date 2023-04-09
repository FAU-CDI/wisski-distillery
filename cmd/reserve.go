package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/fsx"
)

// Reserve is the 'reserve' command
var Reserve wisski_distillery.Command = reserve{}

type reserve struct {
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"name of instance to reserve"`
	} `positional-args:"true"`
}

func (reserve) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "reserve",
		Description: "reserves a new instance",
	}
}

// TODO: AfterParse to check instance!

var errReserveAlreadyExists = exit.Error{
	Message:  "instance %q already exists",
	ExitCode: exit.ExitGeneric,
}

var errReserveGeneric = exit.Error{
	Message:  "unable to provision instance",
	ExitCode: exit.ExitGeneric,
}

func (r reserve) Run(context wisski_distillery.Context) (err error) {
	defer errReserveGeneric.DeferWrap(&err)

	dis := context.Environment
	slug := r.Positionals.Slug

	// check that it doesn't already exist
	logging.LogMessage(context.Stderr, "Reserving new WissKI instance %s", slug)
	if exists, err := dis.Instances().Has(context.Context, slug); err != nil || exists {
		return errReserveAlreadyExists.WithMessageF(slug)
	}

	// make it in-memory
	instance, err := dis.Instances().Create(slug)
	if err != nil {
		return errProvisionGeneric.WithMessageF(slug, err)
	}

	// check that the base directory does not exist
	{
		logging.LogMessage(context.Stderr, "Checking that base directory %s does not exist", instance.FilesystemBase)
		exists, err := fsx.Exists(instance.FilesystemBase)
		if err != nil {
			return errProvisionGeneric.Wrap(err)
		}
		if exists {
			return errReserveAlreadyExists.WithMessageF(slug)
		}
	}

	// setup docker stack
	s := instance.Reserve().Stack()
	{
		if err := logging.LogOperation(func() error {
			return s.Install(context.Context, context.Stderr, component.InstallationContext{})
		}, context.Stderr, "Installing docker stack"); err != nil {
			return err
		}

		if err := logging.LogOperation(func() error {
			return s.Update(context.Context, context.Stderr, true)
		}, context.Stderr, "Updating docker stack"); err != nil {
			return err
		}
	}

	// and we're done!
	logging.LogMessage(context.Stderr, "Instance has been reserved")
	context.Printf("URL:      %s\n", instance.URL().String())

	return nil
}
