package dis

import (
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
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
	dis = new(Distillery)
	dis.Still.Upstream = component.Upstream{
		SQL:         component.HostPort{Host: "127.0.0.1", Port: 3306},
		Triplestore: component.HostPort{Host: "127.0.0.1", Port: 7200},
		Solr:        component.HostPort{Host: "127.0.0.1", Port: 8983},
	}

	// we are within the docker
	//
	// so setup the ports to connect everything to properly.
	// also override some of the parameters for the environment.
	if flags.InternalInDocker {
		dis.Still.Upstream.SQL = component.HostPort{Host: "sql", Port: 3306}
		dis.Still.Upstream.Triplestore = component.HostPort{Host: "triplestore", Port: 7200}
		dis.Still.Upstream.Solr = component.HostPort{Host: "solr", Port: 8983}
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
	err = dis.Config.Unmarshal(f)
	return
}
