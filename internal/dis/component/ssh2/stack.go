package ssh2

//spellchecker:words embed path filepath github wisski distillery internal bootstrap component
import (
	"embed"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (ssh *SSH2) Path() string {
	return filepath.Join(component.GetStill(ssh).Config.Paths.Root, "core", "ssh2")
}

//go:embed all:ssh2
var resources embed.FS

func (ssh *SSH2) OpenStack() (component.StackWithResources, error) {
	config := component.GetStill(ssh).Config
	//nolint:wrapcheck
	return component.OpenStack(ssh, ssh.dependencies.Docker, component.StackWithResources{
		Resources:   resources,
		ContextPath: "ssh2",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),

			"CONFIG_PATH": config.ConfigPath,
			"DEPLOY_ROOT": config.Paths.Root,

			"SELF_OVERRIDES_FILE":      config.Paths.OverridesJSON,
			"SELF_RESOLVER_BLOCK_FILE": config.Paths.ResolverBlocks,
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
}

func (ssh *SSH2) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: component.GetStill(ssh).Config.Paths.CurrentExecutable(), // TODO: Does this make sense?
	}
}
