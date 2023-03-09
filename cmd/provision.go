package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Provision is the 'provision' command
var Provision wisski_distillery.Command = pv{}

type pv struct {
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"slug of instance to create"`
	} `positional-args:"true"`
}

func (pv) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "provision",
		Description: "creates a new instance",
	}
}

var errProvisionGeneric = exit.Error{
	Message:  "unable to provision instance %s",
	ExitCode: exit.ExitGeneric,
}

// TODO: AfterParse to check instance!

func (p pv) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Provision().Provision(context.Stderr, context.Context, provision.ProvisionFlags{
		Slug: p.Positionals.Slug,
	})
	if err != nil {
		return errProvisionGeneric.WithMessageF(p.Positionals.Slug).Wrap(err)
	}

	// and we're done!
	logging.LogMessage(context.Stderr, context.Context, "Instance has been provisioned")
	context.Printf("URL:      %s\n", instance.URL().String())
	context.Printf("Username: %s\n", instance.DrupalUsername)
	context.Printf("Password: %s\n", instance.DrupalPassword)

	return nil
}
