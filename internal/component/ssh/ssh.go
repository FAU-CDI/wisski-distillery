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

func (ssh *SSH) Path() string {
	return filepath.Join(ssh.Still.Config.DeployRoot, "core", "ssh")
}

func (*SSH) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:ssh
//go:embed ssh.env
var resources embed.FS

func (ssh *SSH) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(ssh, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "ssh",

		EnvPath: "ssh.env",
		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": ssh.Config.DockerNetworkName,
		},
	})
}
