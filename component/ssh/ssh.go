package ssh

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/component"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
)

type SSH struct {
	component.ComponentBase
}

func (SSH) Name() string {
	return "ssh"
}

//go:embed all:stack
var resources embed.FS

func (ssh SSH) Stack() stack.Installable {
	return ssh.ComponentBase.MakeStack(stack.Installable{
		Resources:   resources,
		ContextPath: "stack",
	})
}
