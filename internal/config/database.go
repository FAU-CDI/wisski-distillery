package config

type DatabaseConfig struct {
	// Credentials for the admin user.
	// Is automatically created if it does not exist.
	AdminUsername string `yaml:"username" default:"admin" validate:"nonempty"`
	AdminPassword string `yaml:"password" validate:"nonempty"`

	// Prefix for new users and data setss
	UserPrefix string `yaml:"user_prefix" default:"wisski-distillery-" validate:"slug"`
	DataPrefix string `yaml:"fragment_prefix" default:"wisski-distillery-" validate:"slug"`
}

type SQLConfig struct {
	DatabaseConfig `yaml:",inline" recurse:"true"`

	// Database to use to store distillery datastructures
	Database string `yaml:"database" default:"distillery" validate:"slug"`
}

type TSConfig struct {
	DatabaseConfig `yaml:",inline" recurse:"true"`
}
