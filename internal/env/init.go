package env

import (
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/tkw1536/goprogram/exit"
)

var errNoConfigFile = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Configuration File does not exist",
}

var errOpenConfig = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "error loading configuration file: %s",
}

// NewDistillery creates a new distillery object from a set of parameters and requirements
func NewDistillery(params core.Params, flags core.Flags, req core.Requirements) (env *Distillery, err error) {
	env = &Distillery{
		Upstream: Upstream{
			SQL:         "127.0.0.1:3306",
			Triplestore: "127.0.0.1:7200",
		},
	}

	if flags.InternalInDocker {
		env.Upstream.SQL = "sql:3306"
		env.Upstream.Triplestore = "triplestore:7200"
	}

	// if we don't need to load the config, there is nothing to do
	if !req.NeedsDistillery {
		return
	}

	// try to find the configuration file
	cfg := flags.ConfigPath // command line flags first
	if cfg == "" {
		cfg = params.ConfigPath // then globally provided files
	}
	if cfg == "" {
		return nil, errNoConfigFile
	}

	// open the config file!
	f, err := os.Open(params.ConfigPath)
	if err != nil {
		return nil, errOpenConfig.WithMessageF(err)
	}
	defer f.Close()

	// unmarshal the config
	env.Config = &config.Config{
		ConfigPath: cfg,
	}
	err = env.Config.Unmarshal(f)
	return
}
