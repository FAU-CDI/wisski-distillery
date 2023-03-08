// Package legacy provides support for reading legacy configuration.
// It is deprecated and will be removed in a future release.
package legacy

import (
	"io"
	"net/url"
	"reflect"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/config/legacy/envreader"
	"github.com/FAU-CDI/wisski-distillery/internal/config/legacy/stringparser"
	"github.com/FAU-CDI/wisski-distillery/internal/config/validators"
	"github.com/pkg/errors"
)

// Migrate parses a configuration from an old configuration.
func Migrate(config *config.Config, src io.Reader) error {
	var legacy Legacy
	if err := legacy.Unmarshal(src); err != nil {
		return nil
	}
	return legacy.Migrate(config)
}

// Legacy represents a legacy configuration file.
//
// NOTE(twiesing): This will be deprecated soon.
type Legacy struct {
	DeployRoot string `env:"DEPLOY_ROOT" default:"/var/www/deploy" parser:"abspath"`

	DefaultDomain string `env:"DEFAULT_DOMAIN" default:"localhost.kwarc.info" parser:"domain"`

	SelfRedirect *url.URL `env:"SELF_REDIRECT" default:"https://github.com/FAU-CDI/wisski-distillery" parser:"https_url"`

	SelfExtraDomains []string `env:"SELF_EXTRA_DOMAINS" default:"" parser:"domains"`

	SelfOverridesFile string `env:"SELF_OVERRIDES_FILE" default:"" parser:"file"`

	SelfResolverBlockFile string `env:"SELF_RESOLVER_BLOCK_FILE" default:"" parser:"file"`

	CertbotEmail string `env:"CERTBOT_EMAIL" default:"" parser:"email"`

	MaxBackupAge int `env:"MAX_BACKUP_AGE" default:"" parser:"number"`

	MysqlUserPrefix     string `env:"MYSQL_USER_PREFIX" default:"mysql-factory-" parser:"slug"`
	MysqlDatabasePrefix string `env:"MYSQL_DATABASE_PREFIX" default:"mysql-factory-" parser:"slug"`
	GraphDBUserPrefix   string `env:"GRAPHDB_USER_PREFIX" default:"mysql-factory-" parser:"slug"`
	GraphDBRepoPrefix   string `env:"GRAPHDB_REPO_PREFIX" default:"mysql-factory-" parser:"slug"`

	DistilleryDatabase string `env:"DISTILLERY_BOOKKEEPING_DATABASE" default:"distillery" parser:"slug"`

	PasswordLength int `env:"PASSWORD_LENGTH" default:"64" parser:"number"`

	PublicSSHPort uint16 `env:"SSH_PORT" default:"2222" parser:"port"`

	TriplestoreAdminUser     string `env:"GRAPHDB_ADMIN_USER" default:"admin" parser:"nonempty"`
	TriplestoreAdminPassword string `env:"GRAPHDB_ADMIN_PASSWORD" default:"" parser:"nonempty"`

	MysqlAdminUser     string `env:"MYSQL_ADMIN_USER" default:"admin" parser:"nonempty"`
	MysqlAdminPassword string `env:"MYSQL_ADMIN_PASSWORD" default:"" parser:"nonempty"`

	SessionSecret string `env:"SESSION_SECRET" default:"" parser:"nonempty"`

	// name of docker network to use
	DockerNetworkName string        `env:"DOCKER_NETWORK_NAME" default:"distillery" parser:"nonempty"`
	CronInterval      time.Duration `env:"CRON_INTERVAL" default:"10m" parser:"duration"`
}

// Migrate migrates this LegacyConfig into a new configuration.
func (legacy *Legacy) Migrate(cfg *config.Config) error {
	cfg.Paths.Root = legacy.DeployRoot
	cfg.HTTP.PrimaryDomain = legacy.DefaultDomain
	cfg.Theme.SelfRedirect = (*validators.URL)(legacy.SelfRedirect)
	cfg.HTTP.ExtraDomains = legacy.SelfExtraDomains
	cfg.Paths.OverridesJSON = legacy.SelfOverridesFile
	cfg.Paths.ResolverBlocks = legacy.SelfResolverBlockFile
	cfg.HTTP.CertbotEmail = legacy.CertbotEmail
	cfg.MaxBackupAge = time.Duration(legacy.MaxBackupAge) * 24 * time.Hour
	cfg.SQL.UserPrefix = legacy.MysqlUserPrefix
	cfg.SQL.DataPrefix = legacy.MysqlDatabasePrefix
	cfg.TS.UserPrefix = legacy.GraphDBUserPrefix
	cfg.TS.DataPrefix = legacy.GraphDBRepoPrefix
	cfg.SQL.Database = legacy.DistilleryDatabase
	cfg.PasswordLength = legacy.PasswordLength
	cfg.Listen.Ports = []uint16{80, legacy.PublicSSHPort}
	if legacy.CertbotEmail != "" {
		cfg.Listen.Ports = append(cfg.Listen.Ports, 443)
	}
	cfg.Listen.AdvertisedSSHPort = legacy.PublicSSHPort
	cfg.TS.AdminUsername = legacy.TriplestoreAdminUser
	cfg.TS.AdminPassword = legacy.TriplestoreAdminPassword
	cfg.SQL.AdminUsername = legacy.MysqlAdminUser
	cfg.SQL.AdminPassword = legacy.MysqlAdminPassword
	cfg.SessionSecret = legacy.SessionSecret
	cfg.Docker.Network = legacy.DockerNetworkName
	cfg.CronInterval = legacy.CronInterval
	return nil
}

// Unmarshal opens a legacy configuration file.
//
// Data is read using the [envreader.ReadAll] method, see the appropriate documentation for the file format.
//
// The `env` and `parser` reflect tags of the [Config] struct determine the keys to read from, and the types to expect.
// When a key is missing, it is set to the default value.
//
// See also [stringparser.Parse].
func (config *Legacy) Unmarshal(src io.Reader) error {
	// read all the values!
	values, err := envreader.ReadAll(src)
	if err != nil {
		return err
	}

	vConfig := reflect.ValueOf(config).Elem()
	tConfig := vConfig.Type()

	// iterate over the types
	numValues := tConfig.NumField()
	for i := 0; i < numValues; i++ {
		tField := tConfig.Field(i)
		vField := vConfig.FieldByName(tField.Name)

		tEnv := tField.Tag.Get("env")
		tDefault := tField.Tag.Get("default")
		tParser := tField.Tag.Get("parser")

		// skip it if it isn't loaded!
		if tEnv == "" {
			continue
		}

		// read the value with a default
		value, ok := values[tEnv]
		if !ok || value == "" {
			if tDefault != "" {
				value = tDefault
			}
		}

		// parse the value!
		if err := stringparser.Parse(tParser, value, vField); err != nil {
			return errors.Errorf("Legacy.Unmarshal: Setting %q, Parser %q: %s", tEnv, tParser, err)
		}
	}

	return nil
}
