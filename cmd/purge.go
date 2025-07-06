package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"go.tkw01536.de/goprogram/exit"
)

// Provision is the 'provision' command.
var Purge wisski_distillery.Command = purge{}

type purge struct {
	Yes         bool `description:"do not ask for confirmation" long:"yes" short:"y"`
	Positionals struct {
		Slug string `description:"name of instance to purge" positional-arg-name:"slug" required:"1-1"`
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

var (
	errPurgeNoConfirmation = exit.NewErrorWithCode("aborting after request was not confirmed. either type `yes` or pass `--yes` on the command line", exit.ExitGeneric)
	errPurgeFailed         = exit.NewErrorWithCode("failed to run purge", exit.ExitGeneric)
)

func (p purge) Run(context wisski_distillery.Context) error {
	dis := context.Environment
	slug := p.Positionals.Slug

	// check the confirmation from the user
	if !p.Yes {
		_, _ = context.Printf("About to remove instance %q. This cannot be undone.\n", slug)
		_, _ = context.Printf("Type 'yes' to continue: ")
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
