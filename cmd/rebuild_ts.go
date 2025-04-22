package cmd

//spellchecker:words github wisski distillery internal
import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
)

// RebuildTS is the 'rebuild_ts' setting.
var RebuildTS wisski_distillery.Command = rebuildTS{}

type rebuildTS struct {
	AllowEmptyRepository bool `short:"a" long:"allow-empty" description:"don't abort if repository is empty"`
	Positionals          struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to rebuild triplestore for"`
	} `positional-args:"true"`
}

func (rebuildTS) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "rebuild_ts",
		Description: "rebuild the triplestore for a specific instance",
	}
}

func (rts rebuildTS) Run(context wisski_distillery.Context) (err error) {
	instance, err := context.Environment.Instances().WissKI(context.Context, rts.Positionals.Slug)
	if err != nil {
		return err
	}

	_, err = instance.TRB().RebuildTriplestore(context.Context, context.Stdout, rts.AllowEmptyRepository)
	return err
}
