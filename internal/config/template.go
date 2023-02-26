package config

import (
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/password"
	"github.com/tkw1536/pkglib/hostname"
)

// Template is a template for the configuration file
type Template struct {
	DeployRoot               string `env:"DEPLOY_ROOT"`
	DefaultDomain            string `env:"DEFAULT_DOMAIN"`
	SelfOverridesFile        string `env:"SELF_OVERRIDES_FILE"`
	SelfResolverBlockFile    string `env:"SELF_RESOLVER_BLOCK_FILE"`
	TriplestoreAdminUser     string `env:"GRAPHDB_ADMIN_USER"`
	TriplestoreAdminPassword string `env:"GRAPHDB_ADMIN_PASSWORD"`
	MysqlAdminUsername       string `env:"MYSQL_ADMIN_USER"`
	MysqlAdminPassword       string `env:"MYSQL_ADMIN_PASSWORD"`
	DockerNetworkName        string `env:"DOCKER_NETWORK_NAME"`
	SessionSecret            string `env:"SESSION_SECRET"`
}

// SetDefaults sets defaults on the template
func (tpl *Template) SetDefaults(env environment.Environment) (err error) {
	if tpl.DeployRoot == "" {
		tpl.DeployRoot = bootstrap.BaseDirectoryDefault
	}

	if tpl.DefaultDomain == "" {
		tpl.DefaultDomain = hostname.FQDN() // TODO: Make this environment specific
	}

	if tpl.SelfOverridesFile == "" {
		tpl.SelfOverridesFile = filepath.Join(tpl.DeployRoot, bootstrap.OverridesJSON)
	}

	if tpl.SelfResolverBlockFile == "" {
		tpl.SelfResolverBlockFile = filepath.Join(tpl.DeployRoot, bootstrap.ResolverBlockedTXT)
	}

	if tpl.TriplestoreAdminUser == "" {
		tpl.TriplestoreAdminUser = "admin"
	}

	if tpl.TriplestoreAdminPassword == "" {
		tpl.TriplestoreAdminPassword, err = password.Password(64)
		if err != nil {
			return err
		}
	}

	if tpl.MysqlAdminUsername == "" {
		tpl.MysqlAdminUsername = "admin"
	}

	if tpl.MysqlAdminPassword == "" {
		tpl.MysqlAdminPassword, err = password.Password(64)
		if err != nil {
			return err
		}
	}

	if tpl.DockerNetworkName == "" {
		tpl.DockerNetworkName, err = password.Password(10)
		if err != nil {
			return err
		}
		tpl.DockerNetworkName = `distillery-` + tpl.DockerNetworkName
	}

	if tpl.SessionSecret == "" {
		tpl.SessionSecret, err = password.Password(100)
		if err != nil {
			return err
		}
	}

	return nil
}

// Generate generates a configuration file for this configuration
func (tpl Template) Generate() Config {
	return Config{
		Paths: PathsConfig{
			Root:           tpl.DeployRoot,
			OverridesJSON:  tpl.SelfOverridesFile,
			ResolverBlocks: tpl.SelfResolverBlockFile,
		},
		HTTP: HTTPConfig{
			PrimaryDomain: tpl.DefaultDomain,
			ExtraDomains:  []string{},
		},
		Docker: DockerConfig{
			tpl.DockerNetworkName,
		},
		SQL: SQLConfig{
			DatabaseConfig: DatabaseConfig{
				AdminUsername: tpl.MysqlAdminUsername,
				AdminPassword: tpl.MysqlAdminPassword,

				UserPrefix: "mysql-factory-",
				DataPrefix: "mysql-factory-",
			},

			Database: "distillery",
		},
		TS: TSConfig{
			DatabaseConfig: DatabaseConfig{
				AdminUsername: tpl.TriplestoreAdminUser,
				AdminPassword: tpl.TriplestoreAdminPassword,

				UserPrefix: "graphdb-factory-",
				DataPrefix: "graphdb-factory-",
			},
		},
		MaxBackupAge:   30 * 24 * time.Hour, // 1 month
		PasswordLength: 64,

		PublicSSHPort: 2222,

		SessionSecret: tpl.SessionSecret,
		CronInterval:  10 * time.Minute,
	}
}
