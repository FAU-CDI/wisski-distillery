package config

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

type HTTPConfig struct {
	// Each created Drupal Instance corresponds to a single domain name.
	// These domain names should either be a complete domain name or a sub-domain of a default domain.
	// This setting configures the default domain-name to create subdomains of.
	PrimaryDomain string `yaml:"domain" default:"localhost.kwarc.info" validate:"domain"`

	// By default, only the 'self' domain above is caught.
	// To catch additional domains, add them here (comma separated)
	ExtraDomains []string `yaml:"domains" validate:"domains"`

	// The system can support setting up certificate(s) automatically.
	// It can be enabled by setting an email for certbot certificates.
	// This email address can be configured here.
	CertbotEmail string `yaml:"certbot_email" validate:"email"`
}

// JoinPath returns the root public url joined with the provided parts.
func (hcfg HTTPConfig) JoinPath(elem ...string) *url.URL {
	u := url.URL{
		Scheme: "http",
		Host:   hcfg.PrimaryDomain,
		Path:   "/",
	}
	if hcfg.HTTPSEnabled() {
		u.Scheme = "https"
	}

	return u.JoinPath(elem...)
}

// TCPMuxCommand generates a command line for the sslh executable.
func (hcfg HTTPConfig) TCPMuxCommand(addr string, http string, https string, ssh string) string {
	if hcfg.HTTPSEnabled() {
		return fmt.Sprintf("-bind %s -http %s -tls %s -rest %s", addr, http, https, ssh)
	}
	return fmt.Sprintf("-bind %s -http %s -rest %s", addr, http, ssh)
}

// HTTPSEnabled returns if the distillery has HTTPS enabled, and false otherwise.
func (hcfg HTTPConfig) HTTPSEnabled() bool {
	return hcfg.CertbotEmail != ""
}

type SpecialDomain string

var (
	TriplestoreDomain SpecialDomain = "ts"
)

func (hcfg HTTPConfig) SpecialDomain(domain SpecialDomain) string {
	return fmt.Sprintf("%s.%s", string(domain)+"_", hcfg.PrimaryDomain)
}

// Domains adds the given subdomain to the primary and alias domains.
// If sub is empty, returns only the domains.
//
// sub is not otherwise validated, and should be normalized by the caller.
func (hcfg HTTPConfig) Domains(sub string) []string {
	domains := append([]string{hcfg.PrimaryDomain}, hcfg.ExtraDomains...)
	if sub == "" {
		return domains
	}

	for i, d := range domains {
		domains[i] = sub + "." + d
	}
	return domains
}

// HostRule returns a HostRule for the provided subdomain.
// See Domains() for usage of sub.
func (hcfg HTTPConfig) HostRule(sub string) string {
	return MakeHostRule(hcfg.Domains(sub)...)
}

// HTTPSEnabledEnv returns "true" if https is enabled, and "false" otherwise.
func (hcfg HTTPConfig) HTTPSEnabledEnv() string {
	if hcfg.HTTPSEnabled() {
		return "true"
	}
	return "false"
}

// HostFromSlug returns the hostname belonging to a given slug.
// When the slug is empty, returns the default (top-level) domain.
func (cfg HTTPConfig) HostFromSlug(slug string) string {
	if slug == "" {
		return cfg.PrimaryDomain
	}
	return fmt.Sprintf("%s.%s", slug, cfg.PrimaryDomain)
}

// SlugFromHost returns the slug belonging to the appropriate host.'
//
// When host is a top-level domain, returns "", true.
// When no slug is found, returns "", false.
func (cfg HTTPConfig) SlugFromHost(host string) (slug string, ok bool) {
	// extract an ':port' that happens to be in the host.
	domain, _, _ := strings.Cut(host, ":")
	domain = TrimSuffixFold(domain, ".wisski") // remove optional ".wisski" ending that is used inside docker

	domainL := strings.ToLower(domain)

	// check all the possible domain endings
	for _, suffix := range append([]string{cfg.PrimaryDomain}, cfg.ExtraDomains...) {
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

func TrimSuffixFold(s string, suffix string) string {
	if len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// DefaultHostRule returns the default traefik hostname rule for this distillery.
// This consists of the [DefaultDomain] as well as [ExtraDomains].
func (cfg HTTPConfig) DefaultHostRule() string {
	return cfg.HostRule("")
}

// MakeHostRule builds a new Host() rule string to be used by traefik.
func MakeHostRule(hosts ...string) string {
	var builder strings.Builder

	first := true
	for _, host := range hosts {
		// HACK HACK HACK: Very minimal domain validation to prevent validation.
		// Just skip everything that isn't a domain.
		if strings.Contains(host, "`") {
			continue
		}

		if first {
			builder.WriteString("Host(`")
		} else {
			builder.WriteString("||Host(`")
		}

		// domain should be punycode!
		domain, err := idna.ToASCII(host)
		if err != nil {
			domain = host
		}

		builder.WriteString(domain)
		builder.WriteString("`)")

		first = false
	}

	return builder.String()
}
