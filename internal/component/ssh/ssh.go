package ssh

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

type SSH struct {
	component.ComponentBase
}

func (SSH) Name() string {
	return "ssh"
}

func (ssh SSH) Path() string {
	return filepath.Join(ssh.Core.Config.DeployRoot, "core", ssh.Name())
}

func (SSH) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:ssh
var resources embed.FS

func (ssh *SSH) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(ssh, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "ssh",
	})
}
