// Package config contains distillery configuration
package config

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Config represents the configuration of a WissKI Distillery.
//
// Config is read from a byte stream using [Unmarshal].
//
// Config contains many methods that do not require any interaction with any running components.
// Methods that require running components are instead store inside the [Distillery] or an appropriate [Component].
type Config struct {
	// Several docker-compose files are created to manage global services and the system itself.
	// On top of this all real-system space will be created under this directory.
	DeployRoot string `env:"DEPLOY_ROOT" default:"/var/www/deploy" parser:"abspath"`

	// Each created Drupal Instance corresponds to a single domain name.
	// These domain names should either be a complete domain name or a sub-domain of a default domain.
	// This setting configures the default domain-name to create subdomains of.
	DefaultDomain string `env:"DEFAULT_DOMAIN" default:"localhost.kwarc.info" parser:"domain"`

	// By default, the default domain redirects to the distillery repository.
	// If you want to change this, set an alternate domain name here.
	SelfRedirect *url.URL `env:"SELF_REDIRECT" default:"https://github.com/FAU-CDI/wisski-distillery" parser:"https_url"`

	// By default, only the 'self' domain above is caught.
	// To catch additional domains, add them here (comma seperated)
	SelfExtraDomains []string `env:"SELF_EXTRA_DOMAINS" default:"" parser:"domains"`

	// You can override individual URLS in the homepage
	// Do this by adding URLs (without trailing '/'s) into a JSON file
	SelfOverridesFile string `env:"SELF_OVERRIDES_FILE" default:"" parser:"file"`

	// You can block specific prefixes from being picked up by the resolver.
	// Do this by adding one prefix per file.
	SelfResolverBlockFile string `env:"SELF_RESOLVER_BLOCK_FILE" default:"" parser:"file"`

	// The system can support setting up certificate(s) automatically.
	// It can be enabled by setting an email for certbot certificates.
	// This email address can be configured here.
	CertbotEmail string `env:"CERTBOT_EMAIL" default:"" parser:"email"`

	// Maximum age for backup in days
	MaxBackupAge int `env:"MAX_BACKUP_AGE" default:"" parser:"number"`

	// Each Drupal instance requires a corresponding system user, database users and databases.
	// These are also set by the appropriate domain name.
	// To differentiate them from other users of the system, these names can be prefixed.
	// The prefix to use can be configured here.
	// When changing these please consider that no system user may exist that has the same name as a mysql user.
	// This is a MariaDB restriction.
	MysqlUserPrefix     string `env:"MYSQL_USER_PREFIX" default:"mysql-factory-" parser:"slug"`
	MysqlDatabasePrefix string `env:"MYSQL_DATABASE_PREFIX" default:"mysql-factory-" parser:"slug"`
	GraphDBUserPrefix   string `env:"GRAPHDB_USER_PREFIX" default:"mysql-factory-" parser:"slug"`
	GraphDBRepoPrefix   string `env:"GRAPHDB_REPO_PREFIX" default:"mysql-factory-" parser:"slug"`

	// In addition to the filesystem the WissKI distillery requires a single SQL table.
	// It uses this database to store a list of installed things.
	DistilleryDatabase string `env:"DISTILLERY_BOOKKEEPING_DATABASE" default:"distillery" parser:"slug"`

	// Various components use password-based-authentication.
	// These passwords are generated automatically.
	// This variable can be used to determine their length.
	PasswordLength int `env:"PASSWORD_LENGTH" default:"64" parser:"number"`

	// Public port to use for the ssh server
	PublicSSHPort uint16 `env:"SSH_PORT" default:"2222" parser:"port"`

	// A file to be used for global authorized_keys for the ssh server.
	GlobalAuthorizedKeysFile string `env:"GLOBAL_AUTHORIZED_KEYS_FILE" default:"/var/www/deploy/authorized_keys" parser:"file"`

	// admin credentials for graphdb
	TriplestoreAdminUser     string `env:"GRAPHDB_ADMIN_USER" default:"admin" parser:"nonempty"`
	TriplestoreAdminPassword string `env:"GRAPHDB_ADMIN_PASSWORD" default:"" parser:"nonempty"`

	// admin credentials for the Mysql database
	MysqlAdminUser     string `env:"MYSQL_ADMIN_USER" default:"admin" parser:"nonempty"`
	MysqlAdminPassword string `env:"MYSQL_ADMIN_PASSWORD" default:"" parser:"nonempty"`

	// admin credentials for the keycloak server
	KeycloakAdminUser     string `env:"KEYCLOAK_ADMIN_USER" default:"admin" parser:"nonempty"`
	KeycloakAdminPassword string `env:"KEYCLOAK_ADMIN_PASSWORD" default:"" parser:"nonempty"`

	// admin credentials for the dis server
	DisAdminUser     string `env:"DIS_ADMIN_USER" default:"admin" parser:"nonempty"`
	DisAdminPassword string `env:"DIS_ADMIN_PASSWORD" default:"" parser:"nonempty"`

	// name of docker network to use
	DockerNetworkName string `env:"DOCKER_NETWORK_NAME" default:"distillery" parser:"nonempty"`

	// ConfigPath is the path this configuration was loaded from (if any)
	ConfigPath string
}

// String serializes this configuration into a string
func (config Config) String() string {
	values := &strings.Builder{}

	vConfig := reflect.ValueOf(config)
	tConfig := vConfig.Type()

	// iterate over the types
	numValues := tConfig.NumField()
	for i := 0; i < numValues; i++ {
		tField := tConfig.Field(i)
		vField := vConfig.FieldByName(tField.Name)

		env := tField.Tag.Get("env")
		if env == "" {
			continue
		}

		fmt.Fprintf(values, "%s=%v\n", env, vField.Interface())
	}

	return values.String()
}
