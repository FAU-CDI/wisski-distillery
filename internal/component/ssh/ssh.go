package ssh

import (
	"embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

type SSH struct {
	component.ComponentBase
}

func (SSH) Name() string {
	return "ssh"
}

//go:embed all:stack
var resources embed.FS

func (ssh SSH) Stack(env environment.Environment) component.StackWithResources {
	return ssh.ComponentBase.MakeStack(env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "stack",
	})
}
