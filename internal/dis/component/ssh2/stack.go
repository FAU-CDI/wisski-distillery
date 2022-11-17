package ssh2

import (
	"embed"
	"path/filepath"
	"strconv"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
)

func (ssh SSH2) Path() string {
	return filepath.Join(ssh.Still.Config.DeployRoot, "core", "ssh2")
}

//go:embed all:ssh2 ssh2.env
var resources embed.FS

func (ssh *SSH2) Stack(env environment.Environment) component.StackWithResources {
	stt := component.MakeStack(ssh, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "ssh2",
		EnvPath:     "ssh2.env",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": ssh.Config.DockerNetworkName,
			"HOST_RULE":           ssh.Config.DefaultHostRule(),
			"HTTPS_ENABLED":       ssh.Config.HTTPSEnabledEnv(),

			"CONFIG_PATH": ssh.Config.ConfigPath,
			"DEPLOY_ROOT": ssh.Config.DeployRoot,

			"GLOBAL_AUTHORIZED_KEYS_FILE": ssh.Config.GlobalAuthorizedKeysFile,
			"SELF_OVERRIDES_FILE":         ssh.Config.SelfOverridesFile,
			"SELF_RESOLVER_BLOCK_FILE":    ssh.Config.SelfResolverBlockFile,

			"SSH_PORT": strconv.FormatUint(uint64(ssh.Config.PublicSSHPort), 10),
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
	return stt
}

func (ssh SSH2) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: ssh.Config.CurrentExecutable(ssh.Environment), // TODO: Does this make sense?
	}
}
