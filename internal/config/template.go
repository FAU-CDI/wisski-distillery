//spellchecker:words config
package config

//spellchecker:words crypto rand path filepath time github wisski distillery internal bootstrap passwordx pkglib password
import (
	"crypto/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/passwordx"
	"github.com/tkw1536/pkglib/password"
)

// Template is used to generate a configuration file.
//
//nolint:recvcheck
type Template struct {
	RootPath      string
	DefaultDomain string

	TSAdminUser     string
	TSAdminPassword string

	SQLAdminUsername string
	SQLAdminPassword string

	DockerNetworkPrefix string
	SessionSecret       string
}

// SetDefaults sets defaults on the template.
func (tpl *Template) SetDefaults() (err error) {
	if tpl.RootPath == "" {
		tpl.RootPath = bootstrap.BaseDirectoryDefault
	}

	if tpl.DefaultDomain == "" {
		tpl.DefaultDomain, err = os.Hostname()
		if err != nil {
			return err
		}
	}

	if tpl.TSAdminUser == "" {
		tpl.TSAdminUser = "admin"
	}

	if tpl.TSAdminPassword == "" {
		tpl.TSAdminPassword, err = password.Generate(rand.Reader, 64, passwordx.Safe)
		if err != nil {
			return err
		}
	}

	if tpl.SQLAdminUsername == "" {
		tpl.SQLAdminUsername = "admin"
	}

	if tpl.SQLAdminPassword == "" {
		tpl.SQLAdminPassword, err = password.Generate(rand.Reader, 64, passwordx.Safe)
		if err != nil {
			return err
		}
	}

	if tpl.DockerNetworkPrefix == "" {
		tpl.DockerNetworkPrefix, err = password.Generate(rand.Reader, 10, passwordx.Identifier)
		if err != nil {
			return err
		}
		tpl.DockerNetworkPrefix = `distillery-` + tpl.DockerNetworkPrefix
	}

	if tpl.SessionSecret == "" {
		tpl.SessionSecret, err = password.Generate(rand.Reader, 100, passwordx.Printable)
		if err != nil {
			return err
		}
	}

	return nil
}

// Generate generates a configuration file for this configuration.
func (tpl Template) Generate() Config {
	return Config{
		Listen: ListenConfig{
			Ports:   []uint16{80},
			SSHPort: 80,
		},
		Paths: PathsConfig{
			Root:           tpl.RootPath,
			OverridesJSON:  filepath.Join(tpl.RootPath, bootstrap.OverridesJSON),
			ResolverBlocks: filepath.Join(tpl.RootPath, bootstrap.ResolverBlockedTXT),
		},
		HTTP: HTTPConfig{
			PrimaryDomain: tpl.DefaultDomain,
			ExtraDomains:  []string{},
		},
		Docker: DockerConfig{
			NetworkPrefix: tpl.DockerNetworkPrefix,
		},
		SQL: SQLConfig{
			DatabaseConfig: DatabaseConfig{
				AdminUsername: tpl.SQLAdminUsername,
				AdminPassword: tpl.SQLAdminPassword,

				UserPrefix: "mysql-factory-",
				DataPrefix: "mysql-factory-",
			},

			Database: "distillery",
		},
		TS: TSConfig{
			DatabaseConfig: DatabaseConfig{
				AdminUsername: tpl.TSAdminUser,
				AdminPassword: tpl.TSAdminPassword,

				UserPrefix: "graphdb-factory-",
				DataPrefix: "graphdb-factory-",
			},
		},
		MaxBackupAge:   30 * 24 * time.Hour, // 1 month
		PasswordLength: 64,

		SessionSecret: tpl.SessionSecret,
		CronInterval:  10 * time.Minute,
	}
}
