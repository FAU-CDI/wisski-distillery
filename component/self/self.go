package self

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/component"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

type Self struct {
	component.ComponentBase
}

func (Self) Name() string {
	return "self"
}

//go:embed all:stack
//go:embed self.env
var resources embed.FS

func (self Self) Stack() stack.Installable {
	// TODO: Move me into config!
	TARGET := "https://github.com/FAU-CDI/wisski-distillery"
	if self.Config.SelfRedirect != nil { // TODO: move to config!
		TARGET = self.Config.SelfRedirect.String()
	}

	return self.ComponentBase.MakeStack(stack.Installable{
		Resources: resources,

		ContextPath: "stack",
		EnvPath:     "self.env",

		EnvContext: map[string]string{
			"VIRTUAL_HOST":      self.Config.DefaultVirtualHost(),
			"LETSENCRYPT_HOST":  self.Config.DefaultLetsencryptHost(),
			"LETSENCRYPT_EMAIL": self.Config.CertbotEmail,
			"TARGET":            TARGET,
			"OVERRIDES_FILE":    self.Config.SelfOverridesFile,
		},
	})
}
