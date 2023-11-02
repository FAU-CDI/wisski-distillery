package config

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config/validators"
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

	// Also serve the panel on the toplevel domain.
	// Note that the panel is *always* servered under the "panel" domain.
	// Disabling this is not recommended.
	Panel validators.NullableBool `yaml:"panel" validate:"bool" default:"true"`

	// API determines if the API is enabled.
	// In a future version of the distillery, it will be enabled by default.
	API validators.NullableBool `yaml:"api" validate:"bool" default:"false"`

	// TS determintes if the special Triplestore domain is enabled.
	TS validators.NullableBool `yaml:"ts" validate:"bool" default:"false"`

	// PhpMyAdmin determines if the special PhpMyAdmin domain is enabled.
	PhpMyAdmin validators.NullableBool `yaml:"phpmyadmin" validate:"bool" default:"false"`
}

// PanelDomain is the domain name where the control panel runs.
func (hcfg HTTPConfig) PanelDomain() string {
	// if we have panel domain enabled, then return it
	if hcfg.Panel.Set && hcfg.Panel.Value {
		return hcfg.PrimaryDomain
	}

	// else use the domain itself
	return hcfg.Domains(PanelDomain.Domain())[0]
}

// TSDomain returns the full url to the triplestore, if any
func (hcfg HTTPConfig) TSURL() template.URL {
	return hcfg.optionalURL(TriplestoreDomain.Domain(), hcfg.TS)
}

func (hcfg HTTPConfig) PhpMyAdminURL() template.URL {
	return hcfg.optionalURL(PHPMyAdminDomain.Domain(), hcfg.PhpMyAdmin)
}

// optionalURL returns the public-facing url to domain if enabled is true.
func (hcfg HTTPConfig) optionalURL(domain string, enabled validators.NullableBool) template.URL {
	if !(enabled.Set && enabled.Value) {
		return ""
	}

	u := url.URL{
		Scheme: "http",
		Host:   hcfg.Domains(domain)[0],
		Path:   "/",
	}
	if hcfg.HTTPSEnabled() {
		u.Scheme = "https"
	}
	return template.URL(u.String())
}

// JoinPath returns the root public url joined with the provided parts.
func (hcfg HTTPConfig) JoinPath(elem ...string) *url.URL {
	u := url.URL{
		Scheme: "http",
		Host:   hcfg.PanelDomain(),
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

// SpecialDomain represents a reserved domain
type SpecialDomain string

var (
	PanelDomain       SpecialDomain = "panel"
	TriplestoreDomain SpecialDomain = "ts"
	PHPMyAdminDomain  SpecialDomain = "phpmyadmin"
)

var RestrictedSlugs = []string{
	"www",
	"admin",
	PanelDomain.Domain(),
	TriplestoreDomain.Domain(),
	PHPMyAdminDomain.Domain(),
}

func (sd SpecialDomain) Domain() string {
	return string(sd)
}

// Domains adds the given subdomain to the primary and alias domains.
// If sub is empty, returns only the domains.
//
// sub is not otherwise validated, and should be normalized by the caller.
//
// It is guaranteed that the first domain returned will always be the primary domain.
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

// NormSlugFromHost is like SlugFromHost, but normalizes the panel host
func (cfg HTTPConfig) NormSlugFromHost(host string) (string, bool) {
	// if we didn't get a domain, don't do anything
	slug, ok := cfg.SlugFromHost(host)
	if !ok {
		return "", false
	}

	// always serve the panel domain
	if slug == PanelDomain.Domain() {
		return "", true
	}

	// if we don't serve the toplevel domain then the toplevel domain is an error.
	if slug == "" && !(cfg.Panel.Set && cfg.Panel.Value) {
		return "", false
	}

	return slug, true
}

func TrimSuffixFold(s string, suffix string) string {
	if len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// DefaultHostRule returns the host rule for the control panel of this distillery.
func (cfg HTTPConfig) PanelHostRule() string {
	all := cfg.Domains(PanelDomain.Domain())
	if cfg.Panel.Set && cfg.Panel.Value {
		all = append(all, cfg.Domains("")...)
	}
	return MakeHostRule(all...)
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
