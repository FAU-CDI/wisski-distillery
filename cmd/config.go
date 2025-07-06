package cmd

//spellchecker:words github wisski distillery internal goprogram exit
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"go.tkw01536.de/goprogram/exit"
)

// Config is the configuration command.
var Config wisski_distillery.Command = cfg{}

type cfg struct {
	Human bool `description:"Print configuration in human-readable format" long:"human"`
}

func (c cfg) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "config",
		Description: "prints information about configuration",
	}
}

var errMarshalConfig = exit.NewErrorWithCode("unable to marshal config", exit.ExitGeneric)

func (cfg cfg) Run(context wisski_distillery.Context) error {
	if cfg.Human {
		human := context.Environment.Config.MarshalSensitive()
		_, _ = context.Println(human)
		return nil
	}
	if err := context.Environment.Config.Marshal(context.Stdout); err != nil {
		return fmt.Errorf("%w: %w", errMarshalConfig, err)
	}
	return nil
}
