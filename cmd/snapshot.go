package cmd

//spellchecker:words github wisski distillery internal component exporter goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/exporter"
	"github.com/tkw1536/goprogram/exit"
)

// Snapshot creates a snapshot of an instance.
var Snapshot wisski_distillery.Command = snapshot{}

type snapshot struct {
	Keepalive   bool `description:"keep instance running while taking a backup. might lead to inconsistent state" long:"keepalive"    short:"k"`
	StagingOnly bool `description:"do not package into a snapshot archive, but only create a staging directory"   long:"staging-only" short:"s"`

	Parts []string `description:"parts to include in snapshots. defaults to all parts, use l to list all available parts" long:"parts"      short:"p"`
	List  bool     `description:"list available parts"                                                                    long:"list-parts" short:"l"`

	Positionals struct {
		Slug string `description:"slug of instance to take a snapshot of"                                                         positional-arg-name:"SLUG" required:"1-1"`
		Dest string `description:"destination path to write snapshot archive to. defaults to the 'snapshots/archives/' directory" positional-arg-name:"DEST"`
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

var (
	errSnapshotFailed = exit.NewErrorWithCode("failed to make a snapshot", exit.ExitGeneric)
	errSnapshotWissKI = exit.NewErrorWithCode("unable to find WissKI", exit.ExitGeneric)
)

func (sn snapshot) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// list available parts
	if sn.List {
		for _, part := range dis.Exporter().Parts() {
			_, _ = context.Println(part)
		}
		return nil
	}

	// find the instance!
	instance, err := dis.Instances().WissKI(context.Context, sn.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errSnapshotWissKI, err)
	}

	// do a snapshot of it!
	err = dis.Exporter().MakeExport(context.Context, context.Stderr, exporter.ExportTask{
		Dest:        sn.Positionals.Dest,
		StagingOnly: sn.StagingOnly,

		SnapshotDescription: exporter.SnapshotDescription{
			Parts: sn.Parts,
		},
		Instance: instance,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errSnapshotFailed, err)
	}
	return nil
}
