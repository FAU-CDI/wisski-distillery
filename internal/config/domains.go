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

// SlugFromHost returns the slug belonging to the appropriate host.'
//
// When host is a top-level domain, returns "", true.
// When no slug is found, returns "", false.
func (cfg Config) SlugFromHost(host string) (slug string, ok bool) {
	// extract an ':port' that happens to be in the host.
	domain, _, _ := strings.Cut(host, ":")
	domainL := strings.ToLower(domain)

	// check all the possible domain endings
	for _, suffix := range append([]string{cfg.DefaultDomain}, cfg.SelfExtraDomains...) {
		suffixL := strings.ToLower(suffix)
		if domainL == suffixL {
			return "", true
		}
		if strings.HasSuffix(domainL, "."+suffixL) {
			return domain[:len(domain)-len(suffix)-1], true
		}
	}

	// no domain found!
	return "", ok
}
