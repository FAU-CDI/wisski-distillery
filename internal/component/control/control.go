package control

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Control represents the running control server.
type Control struct {
	component.ComponentBase

	Servables []component.Servable
}

func (control Control) Name() string {
	return "dis" // TODO: Rename this to control!
}

func (control Control) Path() string {
	return filepath.Join(control.Core.Config.DeployRoot, "core", control.Name())
}

//go:embed all:control control.env
var resources embed.FS

func (control *Control) Stack(env environment.Environment) component.StackWithResources {
	stt := component.MakeStack(control, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "control",
		EnvPath:     "control.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": control.Config.DockerNetworkName,
			"HOST_RULE":           control.Config.DefaultHostRule(),
			"HTTPS_ENABLED":       control.Config.HTTPSEnabledEnv(),

			"CONFIG_PATH": control.Config.ConfigPath,
			"DEPLOY_ROOT": control.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": control.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         control.Config.SelfOverridesFile,
			"SELF_RESOLVER_BLOCK_FILE":    control.Config.SelfResolverBlockFile,
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
	return stt
}

func (control Control) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: control.Config.CurrentExecutable(control.Environment), // TODO: Does this make sense?
	}
}
