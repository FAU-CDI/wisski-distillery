package dis

import (
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
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

// NewDistillery creates a new distillery from the provided flags
func NewDistillery(params core.Params, flags core.Flags, req core.Requirements) (dis *Distillery, err error) {
	dis = &Distillery{
		context: params.Context,
		Core: component.Core{
			Environment: environment.Native{},
		},
		Upstream: Upstream{
			SQL:         "127.0.0.1:3306",
			Triplestore: "127.0.0.1:7200",
		},
	}

	// we are within the docker
	//
	// so setup the ports to connect everything to peroperly.
	// also override some of the parameters for the environment.
	if flags.InternalInDocker {
		dis.Upstream.SQL = "sql:3306"
		dis.Upstream.Triplestore = "triplestore:7200"
		params.ConfigPath = dis.Core.Environment.GetEnv("CONFIG_PATH")
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
	f, err := dis.Core.Environment.Open(params.ConfigPath)
	if err != nil {
		return nil, errOpenConfig.WithMessageF(err)
	}
	defer f.Close()

	// unmarshal the config
	dis.Config = &config.Config{
		ConfigPath: cfg,
	}
	err = dis.Config.Unmarshal(dis.Core.Environment, f)
	return
}
