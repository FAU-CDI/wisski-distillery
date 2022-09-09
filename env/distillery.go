package env

import (
	"context"
	"os"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/tkw1536/goprogram/exit"
)

// Distillery represents a running instance for the distillery
type Distillery struct {
	Config   *config.Config
	Upstream Upstream
}

// Upstream are the upstream urls connecting to the various external components.
type Upstream struct {
	SQL         string
	Triplestore string
}

func (dis Distillery) HTTPSEnabled() bool {
	return dis.Config.CertbotEmail != ""
}

// Returns the default virtual host
func (dis Distillery) DefaultVirtualHost() string {
	VIRTUAL_HOST := dis.Config.DefaultDomain
	if len(dis.Config.SelfExtraDomains) > 0 {
		VIRTUAL_HOST += "," + strings.Join(dis.Config.SelfExtraDomains, ",")
	}
	return VIRTUAL_HOST
}

func (dis Distillery) DefaultLetsencryptHost() string {
	if !dis.HTTPSEnabled() {
		return ""
	}
	return dis.DefaultVirtualHost()
}

// Context returns a new Context belonging to this distillery
func (dis Distillery) Context() context.Context {
	return context.Background()
}

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
