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
		Slug string `positional-arg-name:"slug" required:"1-1" description:"slug of instance to create"`
	} `positional-args:"true"`
}

func (provision) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "provision",
		Description: "creates a new instance",
	}
}

// TODO: AfterParse to check instance!

var errProvisionAlreadyExists = exit.Error{
	Message:  "instance %q already exists",
	ExitCode: exit.ExitGeneric,
}

var errProvisionGeneric = exit.Error{
	Message:  "unable to provision instance %s: %s",
	ExitCode: exit.ExitGeneric,
}

func (p provision) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	slug := p.Positionals.Slug

	// check that it doesn't already exist
	logging.LogMessage(context.Stderr, context.Context, "Provisioning new WissKI instance %s", slug)
	if exists, err := dis.Instances().Has(context.Context, slug); err != nil || exists {
		return errProvisionAlreadyExists.WithMessageF(slug)
	}

	// make it in-memory
	instance, err := dis.Instances().Create(slug)
	if err != nil {
		return errProvisionGeneric.WithMessageF(slug, err)
	}

	// check that the base directory does not exist
	logging.LogMessage(context.Stderr, context.Context, "Checking that base directory %s does not exist", instance.FilesystemBase)
	if fsx.IsDirectory(dis.Environment, instance.FilesystemBase) {
		return errProvisionAlreadyExists.WithMessageF(slug)
	}

	// Store in the instances table!
	if err := logging.LogOperation(func() error {
		if err := instance.Bookkeeping().Save(context.Context); err != nil {
			return errProvisionGeneric.WithMessageF(slug, err)
		}

		return nil
	}, context.Stderr, context.Context, "Updating bookkeeping database"); err != nil {
		return err
	}

	// create all the resources!
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range dis.Provisionable() {
			logging.LogMessage(context.Stderr, context.Context, "Provisioning %s resources", pc.Name())
			err := pc.Provision(context.Context, instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, context.Stderr, context.Context, "Provisioning instance-specific resources"); err != nil {
		return errProvisionGeneric.WithMessageF(slug, err)
	}

	// run the provision script
	if err := logging.LogOperation(func() error {
		if err := instance.Provisioner().Provision(context.Context, context.Stderr); err != nil {
			return errProvisionGeneric.WithMessageF(slug, err)
		}

		return nil
	}, context.Stderr, context.Context, "Running setup scripts"); err != nil {
		return err
	}

	// start the container!
	logging.LogMessage(context.Stderr, context.Context, "Starting Container")
	if err := instance.Barrel().Stack().Up(context.Context, context.Stderr); err != nil {
		return err
	}

	// and we're done!
	logging.LogMessage(context.Stderr, context.Context, "Instance has been provisioned")
	context.Printf("URL:      %s\n", instance.URL().String())
	context.Printf("Username: %s\n", instance.DrupalUsername)
	context.Printf("Password: %s\n", instance.DrupalPassword)

	return nil
}
