package env

import (
	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

// DisComponent represents the 'dis' layer belonging to a distillery
type DisComponent struct {
	dis *Distillery
}

// Dis returns the DisComponent belonging to this distillery
func (dis *Distillery) Dis() DisComponent {
	return DisComponent{dis: dis}
}

func (DisComponent) Name() string {
	return "dis"
}

func (dis DisComponent) Stack() stack.Installable {
	return dis.dis.makeComponentStack(dis, stack.Installable{
		EnvContext: map[string]string{
			"VIRTUAL_HOST":      dis.dis.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  dis.dis.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": dis.dis.Config.CertbotEmail,

			"CONFIG_PATH": dis.dis.Config.ConfigPath,
			"DEPLOY_ROOT": dis.dis.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": dis.dis.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         dis.dis.Config.SelfOverridesFile,
		},
		CopyContextFiles: []string{core.Executable},
	})
}

func (dis DisComponent) Context(parent stack.InstallationContext) stack.InstallationContext {
	return stack.InstallationContext{
		core.Executable: dis.dis.CurrentExecutable(),
	}
}

func (dis DisComponent) Path() string {
	return dis.Stack().Dir
}
