package ssh

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/component"
)

type SSH struct {
	component.ComponentBase
}

func (SSH) Name() string {
	return "ssh"
}

//go:embed all:stack
var resources embed.FS

func (ssh SSH) Stack() component.Installable {
	return ssh.ComponentBase.MakeStack(component.Installable{
		Resources:   resources,
		ContextPath: "stack",
	})
}
