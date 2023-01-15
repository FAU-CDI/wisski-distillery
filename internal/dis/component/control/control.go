package control

import (
	"context"
	"embed"
	"io"
	"path/filepath"
	"syscall"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Control represents the running control server.
type Control struct {
	component.Base
	Dependencies struct {
		Routeables []component.Routeable
		Cronables  []component.Cronable

		Custom *custom.Custom
	}
}

var (
	_ component.Installable = (*Control)(nil)
)

func (control Control) Path() string {
	return filepath.Join(control.Still.Config.DeployRoot, "core", "dis")
}

//go:embed all:control control.env
var resources embed.FS

func (control *Control) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(control, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "control",
		EnvPath:     "control.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": control.Config.DockerNetworkName,
			"HOST_RULE":           control.Config.DefaultHostRule(),
			"HTTPS_ENABLED":       control.Config.HTTPSEnabledEnv(),

			"CONFIG_PATH": control.Config.ConfigPath,
			"DEPLOY_ROOT": control.Config.DeployRoot,

			"SELF_OVERRIDES_FILE":      control.Config.SelfOverridesFile,
			"SELF_RESOLVER_BLOCK_FILE": control.Config.SelfResolverBlockFile,

			"CUSTOM_ASSETS_PATH": control.Dependencies.Custom.CustomAssetsPath(),
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
}

// Trigger triggers the active cron run to immediatly invoke cron.
func (control *Control) Trigger(ctx context.Context, env environment.Environment) error {
	return control.Stack(env).Kill(ctx, io.Discard, "control", syscall.SIGHUP)
}

func (control Control) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: control.Config.CurrentExecutable(control.Environment), // TODO: Does this make sense?
	}
}
