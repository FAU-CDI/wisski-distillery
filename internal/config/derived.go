package config

import "strings"

// This file contains derived configuration values

func (cfg Config) HTTPSEnabled() bool {
	return cfg.CertbotEmail != ""
}

// Returns the default virtual host
func (cfg Config) DefaultVirtualHost() string {
	VIRTUAL_HOST := cfg.DefaultDomain
	if len(cfg.SelfExtraDomains) > 0 {
		VIRTUAL_HOST += "," + strings.Join(cfg.SelfExtraDomains, ",")
	}
	return VIRTUAL_HOST
}

func (cfg Config) DefaultLetsencryptHost() string {
	if !cfg.HTTPSEnabled() {
		return ""
	}
	return cfg.DefaultVirtualHost()
}
