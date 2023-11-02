package server

import (
	"context"
	"embed"
	"io"
	"path/filepath"
	"syscall"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
)

func (control Server) Path() string {
	return filepath.Join(control.Still.Config.Paths.Root, "core", "dis")
}

//go:embed all:server
var resources embed.FS

func (server *Server) Stack() component.StackWithResources {
	return component.MakeStack(server, component.StackWithResources{
		Resources:   resources,
		ContextPath: "server",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": server.Config.Docker.Network(),
			"HOST_RULE":           server.Config.HTTP.PanelHostRule(),
			"HTTPS_ENABLED":       server.Config.HTTP.HTTPSEnabledEnv(),

			"CONFIG_PATH": server.Config.ConfigPath,
			"DEPLOY_ROOT": server.Config.Paths.Root,

			"SELF_OVERRIDES_FILE":      server.Config.Paths.OverridesJSON,
			"SELF_RESOLVER_BLOCK_FILE": server.Config.Paths.ResolverBlocks,

			"CUSTOM_ASSETS_PATH": server.dependencies.Templating.CustomAssetsPath(),
		},

		CopyContextFiles: []string{bootstrap.Executable},
	})
}

// Trigger triggers the active cron run to immediatly invoke cron.
func (server *Server) Trigger(ctx context.Context) error {
	return server.Stack().Kill(ctx, io.Discard, "control", syscall.SIGHUP)
}

func (server *Server) Context(parent component.InstallationContext) component.InstallationContext {
	return component.InstallationContext{
		bootstrap.Executable: server.Config.Paths.CurrentExecutable(), // TODO: Does this make sense?
	}
}
