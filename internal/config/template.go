package config

import (
	"bytes"
	"io"
	"path/filepath"
	"reflect"

	"github.com/FAU-CDI/wisski-distillery/internal/core"
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
	AuthorizedKeys           string `env:"AUTHORIZED_KEYS_FILE"`
	TriplestoreAdminUser     string `env:"GRAPHDB_ADMIN_USER"`
	TriplestoreAdminPassword string `env:"GRAPHDB_ADMIN_PASSWORD"`
	MysqlAdminUsername       string `env:"MYSQL_ADMIN_USER"`
	MysqlAdminPassword       string `env:"MYSQL_ADMIN_PASSWORD"`
	DisAdminUsername         string `env:"DIS_ADMIN_USER"`
	DisAdminPassword         string `env:"DIS_ADMIN_PASSWORD"`
}

// SetDefaults sets defaults on the template
func (tpl *Template) SetDefaults() (err error) {
	if tpl.DeployRoot == "" {
		tpl.DeployRoot = core.BaseDirectoryDefault
	}

	if tpl.DefaultDomain == "" {
		tpl.DefaultDomain = hostname.FQDN()
	}

	if tpl.SelfOverridesFile == "" {
		tpl.SelfOverridesFile = filepath.Join(tpl.DeployRoot, core.OverridesJSON)
	}

	if tpl.AuthorizedKeys == "" {
		tpl.AuthorizedKeys = filepath.Join(tpl.DeployRoot, core.AuthorizedKeys)
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

	if tpl.DisAdminUsername == "" {
		tpl.DisAdminUsername = "admin"
	}

	if tpl.DisAdminPassword == "" {
		tpl.DisAdminPassword, err = password.Password(64)
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
