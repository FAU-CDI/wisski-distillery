//spellchecker:words config
package config

//spellchecker:words github wisski distillery internal config validators
import "github.com/FAU-CDI/wisski-distillery/internal/config/validators"

type DatabaseConfig struct {
	// Credentials for the admin user.
	// Is automatically created if it does not exist.
	AdminUsername string `yaml:"username" default:"admin" validate:"nonempty"`
	AdminPassword string `yaml:"password" validate:"nonempty"  sensitive:"****"`

	// Prefix for new users and data setss
	UserPrefix string `yaml:"user_prefix" default:"wisski-distillery-" validate:"slug"`
	DataPrefix string `yaml:"data_prefix" default:"wisski-distillery-" validate:"slug"`
}

type SQLConfig struct {
	DatabaseConfig `yaml:",inline" recurse:"true"`

	// Database to use to store distillery datastructures
	Database string `yaml:"database" default:"distillery" validate:"slug"`
}

type TSConfig struct {
	DatabaseConfig `yaml:",inline" recurse:"true"`

	// DangerouslyUseAdapterPrefixes inidicates if scanning for prefixes should just use prefixes declared in all adapters.
	// This may not reflect what is actually in the database.
	DangerouslyUseAdapterPrefixes validators.NullableBool `yaml:"dangerously_use_adapter_prefixes" default:"false" validate:"bool"`
}
