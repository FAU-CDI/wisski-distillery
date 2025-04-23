package dis

//spellchecker:words github wisski distillery internal config component goprogram exit pkglib
import (
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/cgo"
)

var errNoConfigFile = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "configuration file does not exist",
}

var errOpenConfig = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "error loading configuration file: %q",
}

// An error to be returned when cgo is enabled unexpectedly.
var ErrCGoEnabled = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "this functionality is only available when cgo support is disabled. Set `CGO_ENABLED=0' at build time and try again",
}

// NewDistillery creates a new distillery from the provided flags.
func NewDistillery(params cli.Params, flags cli.Flags, req cli.Requirements) (dis *Distillery, err error) {
	// check cgo support to prevent weird error messages
	// this has to happen either when we are inside docker, or when explicity requested by the command.
	if cgo.Enabled && (flags.InternalInDocker || req.FailOnCgo) {
		return nil, ErrCGoEnabled
	}

	dis = new(Distillery)
	dis.Upstream = component.Upstream{
		SQL:         component.HostPort{Host: "127.0.0.1", Port: 3306},
		Triplestore: component.HostPort{Host: "127.0.0.1", Port: 7200},
		Solr:        component.HostPort{Host: "127.0.0.1", Port: 8983},
	}

	// we are within the docker
	//
	// so setup the ports to connect everything to properly.
	// also override some of the parameters for the environment.
	if flags.InternalInDocker {
		dis.Upstream.SQL = component.HostPort{Host: "sql", Port: 3306}
		dis.Upstream.Triplestore = component.HostPort{Host: "triplestore", Port: 7200}
		dis.Upstream.Solr = component.HostPort{Host: "solr", Port: 8983}
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
