package config

import (
	"bytes"
	"io"
	"path/filepath"
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/hostname"
	"github.com/FAU-CDI/wisski-distillery/pkg/password"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"

	_ "embed"
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
		tpl.DefaultDomain = hostname.FQDN(env)
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

//go:embed config_template
var templateBytes []byte

// MarshalTo marshals this template into dst
func (tpl Template) MarshalTo(dst io.Writer) error {
	tplVal := reflect.ValueOf(tpl)
	tplType := reflect.TypeOf(tpl)

	context := make(map[string]string, tplType.NumField())
	for i := 0; i < tplType.NumField(); i++ {
		field := tplType.Field(i)

		key := field.Tag.Get("env")
		value := tplVal.FieldByName(field.Name).String()

		context[key] = value
	}

	return unpack.WriteTemplate(dst, context, bytes.NewReader(templateBytes))
}
