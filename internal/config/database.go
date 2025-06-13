//spellchecker:words config
package config

//spellchecker:words github wisski distillery internal config validators
import "github.com/FAU-CDI/wisski-distillery/internal/config/validators"

type DatabaseConfig struct {
	// Credentials for the admin user.
	AdminUsername string `default:"admin"  validate:"nonempty" yaml:"username"`
	AdminPassword string `sensitive:"****" validate:"nonempty" yaml:"password"`

	// Prefix for new users and data setss
	UserPrefix string `default:"wisski-distillery-" validate:"slug" yaml:"user_prefix"`
	DataPrefix string `default:"wisski-distillery-" validate:"slug" yaml:"data_prefix"`
}

type SQLConfig struct {
	DatabaseConfig `recurse:"true" yaml:",inline"`

	// Database to use to store distillery datastructures
	Database string `default:"distillery" validate:"slug" yaml:"database"`
}

type TSConfig struct {
	DatabaseConfig `recurse:"true" yaml:",inline"`

	// DangerouslyUseAdapterPrefixes inidicates if scanning for prefixes should just use prefixes declared in all adapters.
	// This may not reflect what is actually in the database.
	DangerouslyUseAdapterPrefixes validators.NullableBool `default:"false" validate:"bool" yaml:"dangerously_use_adapter_prefixes"`
}
