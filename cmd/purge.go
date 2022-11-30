package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Provision is the 'provision' command
var Purge wisski_distillery.Command = purge{}

type purge struct {
	Yes         bool `short:"y" long:"yes" description:"Skip asking for confirmation"`
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"name of WissKI Instance to purge"`
	} `positional-args:"true"`
}

func (purge) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "purge",
		Description: "Purges a WissKI Instance",
	}
}

var errPurgeNoDetails = exit.Error{
	Message:  "Unable to find instance details for purge: %s",
	ExitCode: exit.ExitGeneric,
}

var errPurgeNoConfirmation = exit.Error{
	Message:  "Aborting after request was not confirmed. Either type `yes` or pass `--yes` on the command line",
	ExitCode: exit.ExitGeneric,
}

var errPurgeGeneric = exit.Error{
	Message:  "Unable to purge instance %s: %s",
	ExitCode: exit.ExitGeneric,
}

func (p purge) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	slug := p.Positionals.Slug

	// check the confirmation from the user
	if !p.Yes {
		context.Printf("About to remove repository %s. This cannot be undone.\n", slug)
		context.Printf("Type 'yes' to continue: ")
		line, err := context.ReadLine()
		if err != nil || line != "yes" {
			return errPurgeNoConfirmation
		}
	}

	// load the instance (first via bookkeeping, then via defaults)
	logging.LogMessage(context.Stderr, "Checking bookkeeping table")
	instance, err := dis.Instances().WissKI(context.Context, slug)
	if err == instances.ErrWissKINotFound {
		context.Println("Not found in bookkeeping table, assuming defaults")
		instance, err = dis.Instances().Create(slug)
	}
	if err != nil {
		return errPurgeNoDetails.WithMessageF(err)
	}

	// remove docker stack
	logging.LogMessage(context.Stderr, "Stopping and removing docker container")
	if err := instance.Barrel().Stack().Down(context.Context, context.Stderr); err != nil {
		context.EPrintln(err)
	}

	// remove the filesystem
	logging.LogMessage(context.Stderr, "Removing from filesystem %s", instance.FilesystemBase)
	if err := dis.Still.Environment.RemoveAll(instance.FilesystemBase); err != nil {
		context.EPrintln(err)
	}

	// purge all the instance specific resources
	if err := logging.LogOperation(func() error {
		domain := instance.Domain()
		for _, pc := range dis.Provisionable() {
			logging.LogMessage(context.Stderr, "Purging %s resources", pc.Name())
			err := pc.Purge(context.Context, instance.Instance, domain)
			if err != nil {
				return err
			}
		}

		return nil
	}, context.Stderr, "Purging instance-specific resources"); err != nil {
		return errPurgeGeneric.WithMessageF(slug, err)
	}

	// remove from bookkeeping
	logging.LogMessage(context.Stderr, "Removing instance from bookkeeping")
	if err := instance.Bookkeeping().Delete(context.Context); err != nil {
		context.EPrintln(err)
	}

	// remove the filesystem
	logging.LogMessage(context.Stderr, "Remove lock data")
	if instance.Locker().TryUnlock(context.Context) {
		context.EPrintln("instance was not locked")
	}

	logging.LogMessage(context.Stderr, "Instance %s has been purged", slug)
	return nil
}
