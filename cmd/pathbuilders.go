package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
)

// Pathbuilders is the 'pathbuilders' command.
var Pathbuilders wisski_distillery.Command = pathbuilders{}

type pathbuilders struct {
	Positionals struct {
		Slug string `description:"slug of instance to export pathbuilders of"                              positional-arg-name:"SLUG" required:"1-1"`
		Name string `description:"name of pathbuilder to get. if omitted, show a list of all pathbuilders" positional-arg-name:"NAME"`
	} `positional-args:"true"`
}

func (pathbuilders) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "pathbuilder",
		Description: "list pathbuilders of a specific instance",
	}
}

var (
	errPathbuildersExport  = exit.NewErrorWithCode("unable to export pathbuilder", exit.ExitGeneric)
	errPathbuildersNoExist = exit.NewErrorWithCode("pathbuilder does not exist", exit.ExitGeneric)
	errPathbuildersWissKI  = exit.NewErrorWithCode("unable to find WissKI", exit.ExitGeneric)
)

func (pb pathbuilders) Run(context wisski_distillery.Context) error {
	// get the wisski
	instance, err := context.Environment.Instances().WissKI(context.Context, pb.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersWissKI, err)
	}

	// get all of the pathbuilders
	if pb.Positionals.Name == "" {
		names, err := instance.Pathbuilder().All(context.Context, nil)
		if err != nil {
			return fmt.Errorf("%w: %w", errPathbuildersExport, err)
		}
		for _, name := range names {
			_, _ = context.Println(name)
		}
		return nil
	}

	// get all the pathbuilders
	xml, err := instance.Pathbuilder().Get(context.Context, nil, pb.Positionals.Name)
	if xml == "" {
		return fmt.Errorf("%q: %w", pb.Positionals.Name, errPathbuildersNoExist)
	}
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersExport, err)
	}
	_, _ = context.Printf("%s", xml)

	return nil
}
