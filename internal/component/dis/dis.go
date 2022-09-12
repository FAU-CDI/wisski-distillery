package dis

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
)

type Dis struct {
	component.ComponentBase

	Executable string // path to the current executable
}

func (dis Dis) Name() string {
	return "dis"
}

//go:embed all:stack dis.env
var resources embed.FS

func (dis Dis) Stack() component.Installable {
	return dis.ComponentBase.MakeStack(component.Installable{
		Resources:   resources,
		ContextPath: "stack",
		EnvPath:     "dis.env",

		EnvContext: map[string]string{
			"VIRTUAL_HOST":      dis.Config.DefaultHost(),
			"LETSENCRYPT_HOST":  dis.Config.DefaultSSLHost(),
			"LETSENCRYPT_EMAIL": dis.Config.CertbotEmail,

			"CONFIG_PATH": dis.Config.ConfigPath,
			"DEPLOY_ROOT": dis.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": dis.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         dis.Config.SelfOverridesFile,
		},
		CopyContextFiles: []string{core.Executable},
	})
}

func (dis Dis) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		core.Executable: dis.Executable,
	}
}
