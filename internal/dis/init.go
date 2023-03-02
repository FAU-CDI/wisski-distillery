package dis

import (
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/tkw1536/goprogram/exit"
)

var errNoConfigFile = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "configuration file does not exist",
}

var errOpenConfig = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "error loading configuration file: %q",
}

// NewDistillery creates a new distillery from the provided flags
func NewDistillery(params cli.Params, flags cli.Flags, req cli.Requirements) (dis *Distillery, err error) {
	dis = &Distillery{
		Still: component.Still{
			Environment: new(environment.Native),
		},
		Upstream: Upstream{
			SQL:         "127.0.0.1:3306",
			Triplestore: "127.0.0.1:7200",
			Solr:        "127.0.0.1:8983",
		},
	}

	// we are within the docker
	//
	// so setup the ports to connect everything to peroperly.
	// also override some of the parameters for the environment.
	if flags.InternalInDocker {
		dis.Upstream.SQL = "sql:3306"
		dis.Upstream.Triplestore = "triplestore:7200"
		dis.Upstream.Solr = "solr:8983"
		params.ConfigPath = os.Getenv("CONFIG_PATH")
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
	dis.Config = &config.Config{
		ConfigPath: cfg,
	}
	err = dis.Config.Unmarshal(dis.Still.Environment, f)
	return
}
