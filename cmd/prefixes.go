package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Prefixes is then 'prefixes' command.
var Prefixes wisski_distillery.Command = prefixes{}

type prefixes struct {
	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to show prefixes for"`
	} `positional-args:"true"`
}

func (prefixes) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "prefixes",
		Description: "list all prefixes for a specific instance",
	}
}

var errPrefixesGeneric = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to load prefixes",
}

var errPrefixesWissKI = exit.Error{
	Message:  "unable to find WissKI",
	ExitCode: exit.ExitGeneric,
}

func (p prefixes) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, p.Positionals.Slug)
	if err != nil {
		return errPrefixesWissKI.WrapError(err)
	}

	prefixes, err := instance.Prefixes().All(context.Context, nil)
	if err != nil {
		return errPrefixesGeneric.WrapError(err)
	}

	for _, p := range prefixes {
		context.Println(p)
	}

	return nil
}
