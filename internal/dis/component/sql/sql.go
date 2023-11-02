package sql

import (
	"embed"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/yamlx"
	"gopkg.in/yaml.v3"
)

type SQL struct {
	component.Base
	dependencies struct {
		Tables []component.Table
	}

	ServerURL string // upstream server url

	PollInterval time.Duration // duration to wait for during wait
}

var (
	_ component.Backupable    = (*SQL)(nil)
	_ component.Snapshotable  = (*SQL)(nil)
	_ component.Installable   = (*SQL)(nil)
	_ component.Provisionable = (*SQL)(nil)
	_ component.Updatable     = (*SQL)(nil)
)

func (sql *SQL) Path() string {
	return filepath.Join(sql.Still.Config.Paths.Root, "core", "sql")
}

func (*SQL) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

//go:embed all:sql
var resources embed.FS

func (sql *SQL) Stack() component.StackWithResources {
	return component.MakeStack(sql, component.StackWithResources{
		Resources:   resources,
		ContextPath: "sql",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": sql.Config.Docker.Network(),
			"HTTPS_ENABLED":       sql.Config.HTTP.HTTPSEnabledEnv(),
			"HOST_RULE":           sql.Config.HTTP.HostRule(config.PHPMyAdminDomain.Domain()),
		},

		ComposerYML: func(root *yaml.Node) (*yaml.Node, error) {
			// phpmyadmin is exposed => everything is fine
			if sql.Config.HTTP.PhpMyAdmin.Set && sql.Config.HTTP.PhpMyAdmin.Value {
				return root, nil
			}

			// not exposed => remove the appropriate labels
			if err := yamlx.ReplaceWith(root, []string{
				"eu.wiss-ki.barrel.distillery=${DOCKER_NETWORK_NAME}",
			}, "services", "phpmyadmin", "labels"); err != nil {
				return nil, err
			}

			return root, nil
		},

		MakeDirsPerm: umaskfree.DefaultDirPerm,
		MakeDirs: []string{
			"data",
			"imports",
		},
	})
}

const (
	// "mysql"-compatible executable for raw sql queries
	SQLQueryExecutable = "mariadb"

	// "mysqldump"-compatible executable for dumping an entire database
	SQlDumpExecutable = "mariadb-dump"
)
