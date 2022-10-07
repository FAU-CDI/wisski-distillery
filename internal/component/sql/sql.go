package sql

import (
	"context"
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type SQL struct {
	component.ComponentBase

	ServerURL string // upstream server url

	PollContext  context.Context // context to abort polling with
	PollInterval time.Duration   // duration to wait for during wait

	lazyNetwork lazy.Lazy[string]
}

func (*SQL) Name() string {
	return "sql"
}

func (sql *SQL) Path() string {
	return filepath.Join(sql.Core.Config.DeployRoot, "core", sql.Name())
}

func (*SQL) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:sql
//go:embed sql.env
var resources embed.FS

func (sql *SQL) Stack(env environment.Environment) component.StackWithResources {
	return component.MakeStack(sql, env, component.StackWithResources{
		Resources:   resources,
		ContextPath: "sql",

		EnvPath: "sql.env",
		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": sql.Config.DockerNetworkName,
			"HTTPS_ENABLED":       sql.Config.HTTPSEnabledEnv(),
		},

		MakeDirsPerm: environment.DefaultDirPerm,
		MakeDirs: []string{
			"data",
		},
	})
}
