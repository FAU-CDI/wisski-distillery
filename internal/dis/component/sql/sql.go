package sql

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type SQL struct {
	component.Base
	Dependencies struct {
		Tables []component.Table
	}

	ServerURL string // upstream server url

	PollInterval time.Duration // duration to wait for during wait

	lazyNetwork lazy.Lazy[string]
}

var (
	_ component.Backupable    = (*SQL)(nil)
	_ component.Snapshotable  = (*SQL)(nil)
	_ component.Installable   = (*SQL)(nil)
	_ component.Provisionable = (*SQL)(nil)
)

func (sql *SQL) Path() string {
	return filepath.Join(sql.Still.Config.Paths.Root, "core", "sql")
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
			"DOCKER_NETWORK_NAME": sql.Config.Docker.Network,
			"HTTPS_ENABLED":       sql.Config.HTTP.HTTPSEnabledEnv(),
		},

		MakeDirsPerm: environment.DefaultDirPerm,
		MakeDirs: []string{
			"data",
			"imports",
		},
	})
}
