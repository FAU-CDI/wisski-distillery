package config

import (
	"strings"
)

// This file contains domain related derived configuration values.

// HTTPSEnabled returns if the distillery has HTTPS enabled, and false otherwise.
func (cfg Config) HTTPSEnabled() bool {
	return cfg.CertbotEmail != ""
}

// IfHttps returns value when the distillery has https enabled, and the empty string otherwise.
func (cfg Config) IfHttps(value string) string {
	if !cfg.HTTPSEnabled() {
		return ""
	}
	return value
}

// DefaultHost returns the default hostname for the distillery.
//
// This consists of the [DefaultDomain] as well as [ExtraDomains].
// Domain names are concatinated with commas.
func (cfg Config) DefaultHost() string {
	var builder strings.Builder

	builder.WriteString(cfg.DefaultDomain)
	for _, domain := range cfg.SelfExtraDomains {
		builder.WriteRune(',')
		builder.WriteString(domain)
	}

	return builder.String()
}

// DefaultSSLHost returns the default hostname for the ssl version of the distillery.
//
// This is exactly [DefaultHost] when SSL is enabled, and the empty string otherwise.
func (cfg Config) DefaultSSLHost() string {
	return cfg.IfHttps(cfg.DefaultHost())
}

// SlugFromHost returns the slug belonging to the appropriate host.
func (cfg Config) SlugFromHost(host string) (slug string) {
	// extract an ':port' that happens to be in the host.
	domain, _, _ := strings.Cut(host, ":")

	// check all the possible domain endings
	for _, suffix := range append([]string{cfg.DefaultDomain}, cfg.SelfExtraDomains...) {
		if strings.HasSuffix(domain, "."+suffix) {
			return domain[:len(domain)-len(suffix)-1]
		}
	}

	// no domain found!
	return ""
}
