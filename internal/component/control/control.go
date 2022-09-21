package control

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

// Control represents the running control server.
type Control struct {
	component.ComponentBase

	Instances *instances.Instances

	ResolverFile string
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
	return component.MakeStack(control, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "control",
		EnvPath:     "control.env",

		EnvContext: map[string]string{
			"VIRTUAL_HOST":      control.Config.DefaultHost(),
			"LETSENCRYPT_HOST":  control.Config.DefaultSSLHost(),
			"LETSENCRYPT_EMAIL": control.Config.CertbotEmail,

			"CONFIG_PATH": control.Config.ConfigPath,
			"DEPLOY_ROOT": control.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": control.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         control.Config.SelfOverridesFile,
		},

		TouchFiles:       []string{control.ResolverFile},
		CopyContextFiles: []string{core.Executable},
	})
}

func (control Control) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		core.Executable: control.Config.CurrentExecutable(control.Environment), // TODO: Does this make sense?
	}
}
