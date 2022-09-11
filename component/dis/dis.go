package dis

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/component"
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

type Dis struct {
	component.ComponentBase

	// TODO: SQL Component

	Executable string // path to the current executable
}

func (dis Dis) Name() string {
	return "dis"
}

//go:embed all:stack dis.env
var resources embed.FS

func (dis Dis) Stack() stack.Installable {
	return dis.ComponentBase.MakeStack(stack.Installable{
		Resources:   resources,
		ContextPath: "stack",
		EnvPath:     "dis.env",

		EnvContext: map[string]string{
			"VIRTUAL_HOST":      dis.Config.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  dis.Config.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": dis.Config.CertbotEmail,

			"CONFIG_PATH": dis.Config.ConfigPath,
			"DEPLOY_ROOT": dis.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": dis.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         dis.Config.SelfOverridesFile,
		},
		CopyContextFiles: []string{core.Executable},
	})
}

func (dis Dis) Context(parent stack.InstallationContext) stack.InstallationContext {
	return stack.InstallationContext{
		core.Executable: dis.Executable,
	}
}
