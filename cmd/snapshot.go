package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/tkw1536/goprogram/exit"
)

// Snapshot creates a snapshot of an instance
var Snapshot wisski_distillery.Command = snapshot{}

type snapshot struct {
	Keepalive   bool `short:"k" long:"keepalive" description:"keep instance running while taking a backup. might lead to inconsistent state"`
	StagingOnly bool `short:"s" long:"staging-only" description:"do not package into a snapshot archive, but only create a staging directory"`

	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to take a snapshot of"`
		Dest string "positional-arg-name:\"DEST\" description:\"destination path to write snapshot archive to. defaults to the `snapshots/archives/` directory\""
	} `positional-args:"true"`
}

func (snapshot) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "snapshot",
		Description: "generates a snapshot archive for the provided instance",
	}
}

var errSnapshotFailed = exit.Error{
	Message:  "failed to make a snapshot",
	ExitCode: exit.ExitGeneric,
}

var errSnapshotWissKI = exit.Error{
	Message:  "unable to find WissKI",
	ExitCode: exit.ExitGeneric,
}

func (sn snapshot) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// find the instance!
	instance, err := dis.Instances().WissKI(context.Context, sn.Positionals.Slug)
	if err != nil {
		return errSnapshotWissKI.Wrap(err)
	}

	// do a snapshot of it!
	err = dis.Exporter().MakeExport(context.Context, context.Stderr, exporter.ExportTask{
		Dest:        sn.Positionals.Dest,
		StagingOnly: sn.StagingOnly,

		Instance: instance,
	})

	if err != nil {
		return errSnapshotFailed.Wrap(err)
	}
	return nil
}
