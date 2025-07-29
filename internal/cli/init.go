package cli

//spellchecker:words github wisski distillery internal config component pkglib errorsx exit
import (
	"fmt"
	"os"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/cgo"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
)

var (
	errNoConfigFile = exit.NewErrorWithCode("configuration file does not exist", exit.ExitGeneralArguments)
	errOpenConfig   = exit.NewErrorWithCode("error loading configuration file", exit.ExitGeneralArguments)
	errCGoEnabled   = exit.NewErrorWithCode("this functionality is only available when cgo support is disabled. Set `CGO_ENABLED=0' at build time and try again", exit.ExitGeneralArguments)
)

// GetDistillery gets the distillery for the currently running command.
// [SetFlags] and [SetParameters] must have been called.
func GetDistillery(cmd *cobra.Command, req Requirements) (d *dis.Distillery, e error) {
	flags := get[Flags](cmd, flagsKey)
	params := get[Params](cmd, parametersKey)

	// check cgo support to prevent weird error messages
	// this has to happen either when we are inside docker, or when explicity requested by the command.
	if cgo.Enabled && (flags.InternalInDocker || req.FailOnCgo) {
		return nil, errCGoEnabled
	}

	d = new(dis.Distillery)
	d.Upstream = component.Upstream{
		SQL:         component.HostPort{Host: "127.0.0.1", Port: 3306},
		Triplestore: component.HostPort{Host: "127.0.0.1", Port: 7200},
		Solr:        component.HostPort{Host: "127.0.0.1", Port: 8983},
	}

	// we are within the docker
	//
	// so setup the ports to connect everything to properly.
	// also override some of the parameters for the environment.
	if flags.InternalInDocker {
		d.Upstream.SQL = component.HostPort{Host: "sql", Port: 3306}
		d.Upstream.Triplestore = component.HostPort{Host: "triplestore", Port: 7200}
		d.Upstream.Solr = component.HostPort{Host: "solr", Port: 8983}
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
		return nil, fmt.Errorf("%w: %w", errOpenConfig, err)
	}
	defer errorsx.Close(f, &e, "config file")

	// unmarshal the config
	d.Config = &config.Config{
		ConfigPath: cfg,
	}
	e = d.Config.Unmarshal(f)
	return
}
