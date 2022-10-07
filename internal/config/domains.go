package config

import (
	"fmt"
	"strings"

	"github.com/tkw1536/goprogram/lib/collection"
)

// This file contains domain related derived configuration values.

// HTTPSEnabled returns if the distillery has HTTPS enabled, and false otherwise.
func (cfg Config) HTTPSEnabled() bool {
	return cfg.CertbotEmail != ""
}

// HostRequirement returns a traefik rule for the given names
func (Config) HostRule(names ...string) string {
	quoted := collection.MapSlice(names, func(name string) string {
		return "`" + name + "`"
	})
	return fmt.Sprintf("Host(%s)", strings.Join(quoted, ","))
}

// HTTPSEnabledEnv returns "true" if https is enabled, and "false" otherwise.
func (cfg Config) HTTPSEnabledEnv() string {
	if cfg.HTTPSEnabled() {
		return "true"
	}
	return "false"
}

// IfHttps returns value when the distillery has https enabled, and the empty string otherwise.
func (cfg Config) IfHttps(value string) string {
	if !cfg.HTTPSEnabled() {
		return ""
	}
	return value
}

// DefaultHostRule returns the default traefik hostname rule for this distillery.
// This consists of the [DefaultDomain] as well as [ExtraDomains].
func (cfg Config) DefaultHostRule() string {
	return cfg.HostRule(append([]string{cfg.DefaultDomain}, cfg.SelfExtraDomains...)...)
}

// DefaultSSLHost returns the default hostname for the ssl version of the distillery.
//
// This is exactly [DefaultHost] when SSL is enabled, and the empty string otherwise.
func (cfg Config) xDefaultSSLHost() string {
	panic("not implemented")
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
