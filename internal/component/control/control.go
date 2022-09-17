package control

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
)

// Control represents the control server
type Control struct {
	component.ComponentBase

	Instances *instances.Instances

	ResolverFile string
}

func (control Control) Name() string {
	return "dis" // TODO: Rename this to control!
}

//go:embed all:control control.env
var resources embed.FS

func (control Control) Stack() component.StackWithResources {
	return control.ComponentBase.MakeStack(component.StackWithResources{
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
		core.Executable: control.Config.CurrentExecutable(),
	}
}
