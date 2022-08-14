package env

import "github.com/FAU-CDI/wisski-distillery/internal/stack"

func (dis *Distillery) SelfStack() stack.Installable {
	TARGET := "https://github.com/FAU-CDI/wisski-distillery"
	if dis.Config.SelfRedirect != nil {
		TARGET = dis.Config.SelfRedirect.String()
	}

	return dis.asCoreStack(stack.Installable{
		Stack: stack.Stack{
			Name: "self",
		},

		EnvFileContext: map[string]string{
			"VIRTUAL_HOST":      dis.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  dis.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": dis.Config.CertbotEmail,
			"TARGET":            TARGET,
			"OVERRIDES_FILE":    dis.Config.SelfOverridesFile,
		},
	})
}

func (dis *Distillery) SelfStackPath() string {
	return dis.SelfStack().Dir
}
