package cmd

import (
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
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
		Requirements: env.Requirements{
			NeedsConfig: true,
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
	logging.LogMessage(context.IOStream, "Checking bookkeeping table")
	instance, err := dis.Instance(slug)
	if err == env.ErrInstanceNotFound {
		context.Println("Not found in bookkeeping table, assuming defaults")
		instance, err = dis.NewInstance(slug)
	}
	if err != nil {
		return errPurgeNoDetails.WithMessageF(err)
	}

	// remove docker stack
	logging.LogMessage(context.IOStream, "Stopping and removing docker container")
	if err := instance.Stack().Down(context.IOStream); err != nil {
		context.EPrintln(err)
	}

	// remove the filesystem
	logging.LogMessage(context.IOStream, "Removing from filesystem %s", instance.FilesystemBase)
	if err := os.RemoveAll(instance.FilesystemBase); err != nil {
		context.EPrintln(err)
	}

	// remove the triplestore
	logging.LogOperation(func() error {
		logging.LogMessage(context.IOStream, "Removing user %s", instance.GraphDBUser)
		if err := dis.TriplestorePurgeUser(instance.GraphDBUser); err != nil {
			context.EPrintln(err)
		}

		logging.LogMessage(context.IOStream, "Removing repository %s", instance.GraphDBRepository)
		if err := dis.TriplestorePurgeRepo(instance.GraphDBRepository); err != nil {
			context.EPrintln(err)
		}

		return nil
	}, context.IOStream, "Removing from Triplestore")

	// remove the sql
	logging.LogOperation(func() error {
		logging.LogMessage(context.IOStream, "Removing user %s", instance.SqlUser)
		if err := dis.SQLPurgeUser(instance.SqlUser); err != nil {
			context.EPrintln(err)
		}

		logging.LogMessage(context.IOStream, "Removing database %s", instance.SqlDatabase)
		if err := dis.SQLPurgeDatabase(instance.SqlDatabase); err != nil {
			context.EPrintln(err)
		}

		return nil
	}, context.IOStream, "Removing from SQL")

	// remove from bookkeeping
	logging.LogMessage(context.IOStream, "Removing instance from bookkeeping")
	if err := instance.Delete(); err != nil {
		context.EPrintln(err)
	}

	logging.LogMessage(context.IOStream, "Instance %s has been purged", slug)
	return nil
}
