package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/lib/collection"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

// BlindUpdate is the 'blind_update' command
var BlindUpdate wisski_distillery.Command = blindUpdate{}

type blindUpdate struct {
	Parallel    int  `short:"p" long:"parallel" description:"run on (at most) this many instances in parallel. 0 for no limit." default:"1"`
	Force       bool `short:"f" long:"force" description:"force running blind-update even if AutoBlindUpdate is set to false"`
	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run blind-update in"`
	} `positional-args:"true"`
}

func (blindUpdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "blind_update",
		Description: "Runs the blind update in the provided instances",
	}
}

var errBlindUpdateFailed = exit.Error{
	Message:  "Failed to run blind update script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (bu blindUpdate) Run(context wisski_distillery.Context) error {
	// find all the instances!
	wissKIs, err := context.Environment.Instances().Load(bu.Positionals.Slug...)
	if err != nil {
		return err
	}
	if !bu.Force {
		wissKIs = collection.Filter(wissKIs, func(instance *wisski.WissKI) bool {
			return bool(instance.AutoBlindUpdateEnabled)
		})
	}

	// and do the actual blind_update!
	return status.StreamGroup(context.IOStream, bu.Parallel, func(instance *wisski.WissKI, str stream.IOStream) error {
		return instance.BlindUpdate(str)
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("blind_update %q", item.Slug)
	}))
}
