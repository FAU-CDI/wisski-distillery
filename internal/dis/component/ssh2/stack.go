package ssh2

import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (ssh *SSH2) Path() string {
	return filepath.Join(ssh.Still.Config.Paths.Root, "core", "ssh2")
}

//go:embed all:ssh2
var resources embed.FS

func (ssh *SSH2) Stack() component.StackWithResources {
	return component.MakeStack(ssh, component.StackWithResources{
		Resources:   resources,
		ContextPath: "ssh2",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": ssh.Config.Docker.Network(),

			"CONFIG_PATH": ssh.Config.ConfigPath,
			"DEPLOY_ROOT": ssh.Config.Paths.Root,

			"SELF_OVERRIDES_FILE":      ssh.Config.Paths.OverridesJSON,
			"SELF_RESOLVER_BLOCK_FILE": ssh.Config.Paths.ResolverBlocks,
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
}

func (ssh *SSH2) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: ssh.Config.Paths.CurrentExecutable(), // TODO: Does this make sense?
	}
}
