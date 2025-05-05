package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

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

var (
	errPrefixesGeneric = exit.NewErrorWithCode("unable to load prefixes", exit.ExitGeneric)
	errPrefixesWissKI  = exit.NewErrorWithCode("unable to find WissKI", exit.ExitGeneric)
)

func (p prefixes) Run(context wisski_distillery.Context) error {
	instance, err := context.Environment.Instances().WissKI(context.Context, p.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPrefixesWissKI, err)
	}

	prefixes, err := instance.Prefixes().All(context.Context, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", errPrefixesGeneric, err)
	}

	for _, p := range prefixes {
		_, _ = context.Println(p)
	}

	return nil
}
