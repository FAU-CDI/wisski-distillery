//spellchecker:words server
package server

//spellchecker:words context embed path filepath syscall github wisski distillery internal bootstrap component
import (
	"context"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"syscall"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (server *Server) Path() string {
	return filepath.Join(component.GetStill(server).Config.Paths.Root, "core", "dis")
}

//go:embed all:server
var resources embed.FS

func (server *Server) Stack() component.StackWithResources {
	config := component.GetStill(server).Config

	return component.MakeStack(server, component.StackWithResources{
		Resources:   resources,
		ContextPath: "server",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),
			"HOST_RULE":           config.HTTP.PanelHostRule(),
			"HTTPS_ENABLED":       config.HTTP.HTTPSEnabledEnv(),

			"CONFIG_PATH": config.ConfigPath,
			"DEPLOY_ROOT": config.Paths.Root,

			"SELF_OVERRIDES_FILE":      config.Paths.OverridesJSON,
			"SELF_RESOLVER_BLOCK_FILE": config.Paths.ResolverBlocks,

			"CUSTOM_ASSETS_PATH": server.dependencies.Templating.CustomAssetsPath(),
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
}

// Trigger triggers the active cron run to immediatly invoke cron.
func (server *Server) Trigger(ctx context.Context) error {
	if err := server.Stack().Kill(ctx, io.Discard, "control", syscall.SIGHUP); err != nil {
		return fmt.Errorf("failed to trigger 'control' service: %w", err)
	}
	return nil
}

func (server *Server) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: component.GetStill(server).Config.Paths.CurrentExecutable(), // TODO: Does this make sense?
	}
}
