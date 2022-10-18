package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Provision is the 'provision' command
var Provision wisski_distillery.Command = provision{}

type provision struct {
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"name of WissKI Instance to create"`
	} `positional-args:"true"`
}

func (provision) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "provision",
		Description: "Creates a new WissKI Instance",
	}
}

// TODO: AfterParse to check instance!

var errProvisionAlreadyExists = exit.Error{
	Message:  "Instance %q already exists",
	ExitCode: exit.ExitGeneric,
}

var errProvisionGeneric = exit.Error{
	Message:  "Unable to provision instance %s: %s",
	ExitCode: exit.ExitGeneric,
}

func (p provision) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	slug := p.Positionals.Slug

	// check that it doesn't already exist
	logging.LogMessage(context.IOStream, "Provisioning new WissKI instance %s", slug)
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

	// Store in the instances table!
	if err := logging.LogOperation(func() error {
		if err := instance.Bookkeeping().Save(); err != nil {
			return errProvisionGeneric.WithMessageF(slug, err)
		}

		return nil
	}, context.IOStream, "Updating bookkeeping database"); err != nil {
		return err
	}

	// create all the resources!
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range dis.Provisionable() {
			logging.LogMessage(context.IOStream, "Provisioning %s resources", pc.Name())
			err := pc.Provision(instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, context.IOStream, "Provisioning instance-specific resources"); err != nil {
		return errProvisionGeneric.WithMessageF(slug, err)
	}

	// run the provision script
	if err := logging.LogOperation(func() error {
		if err := instance.Provisioner().Provision(context.IOStream); err != nil {
			return errProvisionGeneric.WithMessageF(slug, err)
		}

		return nil
	}, context.IOStream, "Running setup scripts"); err != nil {
		return err
	}

	// start the container!
	logging.LogMessage(context.IOStream, "Starting Container")
	if err := instance.Barrel().Stack().Up(context.IOStream); err != nil {
		return err
	}

	// and we're done!
	logging.LogMessage(context.IOStream, "Instance has been provisioned")
	context.Printf("URL:      %s\n", instance.URL().String())
	context.Printf("Username: %s\n", instance.DrupalUsername)
	context.Printf("Password: %s\n", instance.DrupalPassword)

	return nil
}
