package cmd

//spellchecker:words github wisski distillery internal goprogram exit pkglib collection status
import (
	"fmt"
	"io"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/status"
)

// BlindUpdate is the 'blind_update' command.
var BlindUpdate wisski_distillery.Command = blindUpdate{}

type blindUpdate struct {
	Parallel    int  `default:"1"                                                                        description:"run on (at most) this many instances in parallel. 0 for no limit" long:"parallel" short:"p"`
	Force       bool `description:"force running blind-update even if 'AutoBlindUpdate' is set to false" long:"force"                                                                   short:"f"`
	Positionals struct {
		Slug []string `description:"slug of instances to update" positional-arg-name:"SLUG" required:"0"`
	} `positional-args:"true"`
}

func (blindUpdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "blind_update",
		Description: "runs the blind update in the provided instances",
	}
}

var errBlindUpdateFailed = exit.NewErrorWithCode("failed to run blind update", exit.ExitGeneric)

func (bu blindUpdate) Run(context wisski_distillery.Context) (err error) {
	// find all the instances!
	wissKIs, err := context.Environment.Instances().Load(context.Context, bu.Positionals.Slug...)
	if err != nil {
		return fmt.Errorf("%w: %w", errBlindUpdateFailed, err)
	}
	if !bu.Force {
		wissKIs = collection.KeepFunc(wissKIs, func(instance *wisski.WissKI) bool {
			return bool(instance.AutoBlindUpdateEnabled)
		})
	}

	// and do the actual blind_update!
	if err := status.WriterGroup(context.Stderr, bu.Parallel, func(instance *wisski.WissKI, writer io.Writer) error {
		return instance.Composer().Update(context.Context, writer)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("blind_update %q", item.Slug)
	})); err != nil {
		return fmt.Errorf("%w: failed to blind_update: %w", errBlindUpdateFailed, err)
	}
	return nil
}
