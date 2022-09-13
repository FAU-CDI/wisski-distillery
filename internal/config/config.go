// Package config implements reading and validating a WissKIDistillery configuration file.
package config

import (
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// Config represents the configuration of a distillery instance
type Config struct {
	// Several docker-compose files are created to manage global services and the system itself.
	// On top of this all real-system space will be created under this directory.
	DeployRoot string `env:"DEPLOY_ROOT" default:"/var/www/deploy" validator:"is_valid_abspath"`

	// Each created Drupal Instance corresponds to a single domain name.
	// These domain names should either be a complete domain name or a sub-domain of a default domain.
	// This setting configures the default domain-name to create subdomains of.
	DefaultDomain string `env:"DEFAULT_DOMAIN" default:"localhost.kwarc.info" validator:"is_valid_domain"`

	// By default, the default domain redirects to the distillery repository.
	// If you want to change this, set an alternate domain name here.
	SelfRedirect *url.URL `env:"SELF_REDIRECT" default:"" validator:"is_valid_https_url"`

	// By default, only the 'self' domain above is caught.
	// To catch additional domains, add them here (comma seperated)
	SelfExtraDomains []string `env:"SELF_EXTRA_DOMAINS" default:"" validator:"is_valid_domains"`

	// You can override individual URLS in the homepage
	// Do this by adding URLs (without trailing '/'s) into a JSON file
	SelfOverridesFile string `env:"SELF_OVERRIDES_FILE" default:"" validator:"is_valid_file"`

	// The system can support setting up certificate(s) automatically.
	// It can be enabled by setting an email for certbot certificates.
	// This email address can be configured here.
	CertbotEmail string `env:"CERTBOT_EMAIL" default:"" validator:"is_valid_email"`

	// Maximum age for backup in days
	MaxBackupAge int `env:"MAX_BACKUP_AGE" default:"" validator:"is_valid_number"`

	// Each Drupal instance requires a corresponding system user, database users and databases.
	// These are also set by the appropriate domain name.
	// To differentiate them from other users of the system, these names can be prefixed.
	// The prefix to use can be configured here.
	// When changing these please consider that no system user may exist that has the same name as a mysql user.
	// This is a MariaDB restriction.
	MysqlUserPrefix     string `env:"MYSQL_USER_PREFIX" default:"mysql-factory-" validator:"is_valid_slug"`
	MysqlDatabasePrefix string `env:"MYSQL_DATABASE_PREFIX" default:"mysql-factory-" validator:"is_valid_slug"`
	GraphDBUserPrefix   string `env:"GRAPHDB_USER_PREFIX" default:"mysql-factory-" validator:"is_valid_slug"`
	GraphDBRepoPrefix   string `env:"GRAPHDB_REPO_PREFIX" default:"mysql-factory-" validator:"is_valid_slug"`

	// In addition to the filesystem the WissKI distillery requires a single SQL table.
	// It uses this database to store a list of installed things.
	DistilleryBookkeepingDatabase string `env:"DISTILLERY_BOOKKEEPING_DATABASE" default:"distillery" validator:"is_valid_slug"`
	DistilleryBookkeepingTable    string `env:"DISTILLERY_BOOKKEEPING_TABLE" default:"distillery" validator:"is_valid_slug"`

	// Various components use password-based-authentication.
	// These passwords are generated automatically.
	// This variable can be used to determine their length.
	PasswordLength int `env:"PASSWORD_LENGTH" default:"64" validator:"is_valid_number"`

	// A file to be used for global authorized_keys for the ssh server.
	GlobalAuthorizedKeysFile string `env:"GLOBAL_AUTHORIZED_KEYS_FILE" default:"/var/www/deploy/authorized_keys" validator:"is_valid_file"`

	// admin credentials for graphdb
	TriplestoreAdminUser     string `env:"GRAPHDB_ADMIN_USER" default:"admin" validator:"is_nonempty"`
	TriplestoreAdminPassword string `env:"GRAPHDB_ADMIN_PASSWORD" default:"" validator:"is_nonempty"`

	// admin credentials for the Mysql database
	MysqlAdminUser     string `env:"MYSQL_ADMIN_USER" default:"admin" validator:"is_nonempty"`
	MysqlAdminPassword string `env:"MYSQL_ADMIN_PASSWORD" default:"admin" validator:"is_nonempty"`

	// ConfigPath is the path this configuration was loaded from (if any)
	ConfigPath string
}

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

func (config *Config) Unmarshal(src io.Reader) error {
	// read all the values!
	values, err := ReadAll(src)
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

		env := tField.Tag.Get("env")
		dflt := tField.Tag.Get("default")
		validator := tField.Tag.Get("validator")

		// skip it if it isn't loaded!
		if env == "" {
			continue
		}

		// read the value with a default
		value, ok := values[env]
		if !ok || value == "" {
			if dflt == "" {
				continue
			}
			value = dflt
		}

		// use the validator
		vFunc, ok := knownValidators[validator]
		if vFunc == nil || !ok {
			return errors.Errorf("Unable to read %q refers to unknown validator %s", env, validator)
		}

		// get the parsed value
		checked, err := vFunc(value)
		if err != nil {
			return errors.Wrapf(err, "Unable to read %q: Validator %s", env, validator)
		}

		// set the value of the field
		var errSet interface{}
		func() {
			defer func() {
				errSet = recover()
			}()
			vField.Set(reflect.ValueOf(checked))
		}()

		// capture any error
		if errSet != nil {
			return errors.Errorf("Unable to parse %q: validator %s returned %q", tField.Name, validator, errSet)
		}
	}

	return nil
}
