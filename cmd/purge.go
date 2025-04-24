package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Provision is the 'provision' command.
var Purge wisski_distillery.Command = purge{}

type purge struct {
	Yes         bool `short:"y" long:"yes" description:"do not ask for confirmation"`
	Positionals struct {
		Slug string `positional-arg-name:"slug" required:"1-1" description:"name of instance to purge"`
	} `positional-args:"true"`
}

func (purge) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "purge",
		Description: "purges an instance",
	}
}

var errPurgeNoConfirmation = exit.Error{
	Message:  "aborting after request was not confirmed. either type `yes` or pass `--yes` on the command line",
	ExitCode: exit.ExitGeneric,
}

var errPurgeFailed = exit.Error{
	Message:  "failed to run purge",
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

	// do the purge!
	if err := dis.Purger().Purge(context.Context, context.Stdout, slug); err != nil {
		return fmt.Errorf("%w: %w", errPurgeFailed, err)
	}
	return nil
}
