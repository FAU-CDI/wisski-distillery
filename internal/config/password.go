package config

import "github.com/FAU-CDI/wisski-distillery/internal/password"

// NewPassword returns a new password using the password settings from this configuration
func (cfg Config) NewPassword() (string, error) {
	return password.Password(cfg.PasswordLength)
}
