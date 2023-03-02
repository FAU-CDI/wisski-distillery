package ssh2

import (
	"embed"
	"path/filepath"
	"strconv"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (ssh SSH2) Path() string {
	return filepath.Join(ssh.Still.Config.Paths.Root, "core", "ssh2")
}

//go:embed all:ssh2 ssh2.env
var resources embed.FS

func (ssh *SSH2) Stack() component.StackWithResources {
	stt := component.MakeStack(ssh, component.StackWithResources{
		Resources:   resources,
		ContextPath: "ssh2",
		EnvPath:     "ssh2.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": ssh.Config.Docker.Network,
			"HOST_RULE":           ssh.Config.HTTP.DefaultHostRule(),
			"HTTPS_ENABLED":       ssh.Config.HTTP.HTTPSEnabledEnv(),

			"CONFIG_PATH": ssh.Config.ConfigPath,
			"DEPLOY_ROOT": ssh.Config.Paths.Root,

			"SELF_OVERRIDES_FILE":      ssh.Config.Paths.OverridesJSON,
			"SELF_RESOLVER_BLOCK_FILE": ssh.Config.Paths.ResolverBlocks,

			"SSH_PORT": strconv.FormatUint(uint64(ssh.Config.PublicSSHPort), 10),
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
	return stt
}

func (ssh SSH2) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: ssh.Config.Paths.CurrentExecutable(), // TODO: Does this make sense?
	}
}
