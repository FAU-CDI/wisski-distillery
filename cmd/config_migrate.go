package cmd

import (
	"fmt"
	"os"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/config/legacy"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// ConfigMigrate is the config-migrate command
var ConfigMigrate wisski_distillery.Command = cfgMigrate{}

type cfgMigrate struct {
	Positionals struct {
		Input string `positional-arg-name:"input" required:"1-1" description:"old config to migrate"`
	} `positional-args:"true"`
}

func (cfgMigrate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: false,
		},
		Command:     "config_migrate",
		Description: "migrate legacy configuration",
	}
}

func (c cfgMigrate) Run(context wisski_distillery.Context) error {
	// migration environment is the native environment!
	var env environment.Environment

	// open the legacy file
	file, err := os.Open(c.Positionals.Input)
	if err != nil {
		return err
	}
	defer file.Close()

	// migrate from a legacy configuration
	// then marshal, and re-read

	var cfg config.Config
	// migrate the legacy config
	if err := legacy.Migrate(&cfg, env, file); err != nil {
		return err
	}

	// validate it!
	if err := cfg.Validate(env); err != nil {
		return err
	}

	// marshal the config
	bytes, err := config.Marshal(&cfg, nil)
	if err != nil {
		return err
	}

	// and print it!
	fmt.Println(string(bytes))
	return nil
}
