package cmd

import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/status"
	"github.com/tkw1536/goprogram/stream"
)

// Cron is the 'cron' command
var UpdatePrefixConfig wisski_distillery.Command = updateprefixconfig{}

type updateprefixconfig struct {
	Parallel int `short:"p" long:"parallel" description:"run on (at most) this many instances in parallel. 0 for no limit." default:"1"`
}

func (updateprefixconfig) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "update_prefix_config",
		Description: "Updates the prefix configuration",
	}
}

var errPrefixUpdateFailed = exit.Error{
	Message:  "Failed to update the prefix configuration",
	ExitCode: exit.ExitGeneric,
}

func (upc updateprefixconfig) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	wissKIs, err := dis.Instances().All()
	if err != nil {
		return errPrefixUpdateFailed.Wrap(err)
	}

	return status.StreamGroup(context.IOStream, upc.Parallel, func(instance *wisski.WissKI, io stream.IOStream) error {
		io.Println("reading prefixes")
		err := instance.Prefixes().UpdatePrefixes()
		if err != nil {
			return errPrefixUpdateFailed.Wrap(err)
		}
		return nil
	}, wissKIs, status.SmartMessage(func(item *wisski.WissKI) string {
		return fmt.Sprintf("update_prefix %q", item.Slug)
	}))
}
