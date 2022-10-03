package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/smartp"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

// Cron is the 'cron' command
var Cron wisski_distillery.Command = cron{}

type cron struct {
	Parallel int `short:"p" long:"parallel" description:"run on (at most) this many instances in parallel. 0 for no limit." default:"1"`

	Positionals struct {
		Slug []string `positional-arg-name:"SLUG" required:"0" description:"slug of instance(s) to run cron in"`
	} `positional-args:"true"`
}

func (cron) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		Command:     "cron",
		Description: "Runs the cron script for several instances",
	}
}

var errCronFailed = exit.Error{
	Message:  "Failed to run cron script for instance %q: exited with code %s",
	ExitCode: exit.ExitGeneric,
}

func (cr cron) Run(context wisski_distillery.Context) error {
	// find all the instances!
	wissKIs, err := context.Environment.Instances().Load(cr.Positionals.Slug...)
	if err != nil {
		return err
	}

	// and do the actual blind_update!
	return smartp.Run(context.IOStream, cr.Parallel, func(instance instances.WissKI, io stream.IOStream) error {
		code, err := instance.Shell(io, "/runtime/cron.sh")
		if err != nil {
			io.EPrintln(err)
		}
		if code != 0 {
			// keep going, because we want to run as many crons as possible
			err = errBlindUpdateFailed.WithMessageF(instance.Slug, code)
			io.EPrintln(err)
		}

		return nil
	}, wissKIs, smartp.SmartMessage(func(item instances.WissKI) string {
		return fmt.Sprintf("cron %q", item.Slug)
	}))
}
