package env

import (
	"context"
	"os"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/tkw1536/goprogram/exit"
)

// Distillery represents a running instance for the distillery
type Distillery struct {
	Config *config.Config
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
func NewDistillery(params Params, req Requirements) (env *Distillery, err error) {
	env = &Distillery{}

	// if we don't need to load the config, there is nothing to do
	if !req.NeedsConfig {
		return
	}

	// if there is no no config file, return
	cfg := params.ConfigFilePath()
	if cfg == "" {
		return nil, errNoConfigFile
	}

	f, err := os.Open(params.ConfigFilePath())
	if err != nil {
		return nil, errOpenConfig.WithMessageF(err)
	}
	defer f.Close()

	// unmarshal the config
	env.Config = &config.Config{}
	err = env.Config.Unmarshal(f)
	return
}
