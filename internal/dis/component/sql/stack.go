package sql

//spellchecker:words embed path filepath time github wisski distillery internal config package component docker pkglib umaskfree yamlx gopkg yaml
import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	config_package "github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/docker"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/yamlx"
	"gopkg.in/yaml.v3"
)

type SQL struct {
	component.Base
	dependencies struct {
		Tables []component.Table
		Docker *docker.Docker
	}

	m  sync.Mutex // m protects db
	db *sql.DB

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
	return filepath.Join(component.GetStill(sql).Config.Paths.Root, "core", "sql")
}

func (*SQL) Context(parent component.InstallationContext) component.InstallationContext {
	return parent
}

const (
	queryExecutable = "mariadb"      // "mysql"-compatible executable for raw sql queries.
	dumpExecutable  = "mariadb-dump" // "mysqldump"-compatible executable for dumping an entire database.
)

//go:embed all:sql
var resources embed.FS

func (sql *SQL) OpenStack() (component.StackWithResources, error) {
	config := component.GetStill(sql).Config
	//nolint:wrapcheck
	return component.OpenStack(sql, sql.dependencies.Docker, component.StackWithResources{
		Resources:   resources,
		ContextPath: "sql",

		EnvContext: map[string]string{
			"DOCKER_NETWORK_NAME": config.Docker.Network(),
			"HTTPS_ENABLED":       config.HTTP.HTTPSEnabledEnv(),
			"HOST_RULE":           config.HTTP.HostRule(config_package.PHPMyAdminDomain.Domain()),
		},

		ComposerYML: func(root *yaml.Node) (*yaml.Node, error) {
			// phpmyadmin is exposed => everything is fine
			if config.HTTP.PhpMyAdmin.Set && config.HTTP.PhpMyAdmin.Value {
				return root, nil
			}

			// not exposed => remove the appropriate labels
			if err := yamlx.ReplaceWith(root, []string{
				"eu.wiss-ki.barrel.distillery=${DOCKER_NETWORK_NAME}",
			}, "services", "phpmyadmin", "labels"); err != nil {
				return nil, fmt.Errorf("failed to replace yaml network name: %w", err)
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
