//spellchecker:words models
package models

// System represents system information.
// It is embedded into the instances struct by gorm.
type System struct {
	// NOTE(twiesing): Any changes here should be reflected in instance_{provision,rebuild}.html and remote/api.ts.
	PHP                string `gorm:"column:php;not null"`                    // php version to use
	IIPServer          bool   `gorm:"column:iipimage;not null;default:false"` // should we enable the IIPServer?
	OpCacheDevelopment bool   `gorm:"column:opcache_devel;not null"`          // opcache development

	ContentSecurityPolicy string `gorm:"column:csp;not null"` // content security policy for the system
}

const (
	imagePrefix = "docker.io/library/php:"
	imageSuffix = "-apache-bullseye"
)

// OpCacheMode returns the name of the `opcache-*.ini` configuration being included in the docker image
func (system System) OpCacheMode() string {
	if system.OpCacheDevelopment {
		return "devel"
	}
	return "prod"
}

var (
	phpVersions   = []string{"8.1", "8.2", "8.3", "8.4"}
	phpVersionMap = (func() map[string]struct{} {
		m := make(map[string]struct{}, len(phpVersions))
		for _, v := range phpVersions {
			m[v] = struct{}{}
		}
		return m
	})()
)

// DefaultPHPVersion is the default php version
const DefaultPHPVersion = "8.3"

// KnownPHPVersions returns a slice of php versions.
func KnownPHPVersions() []string {
	return append([]string(nil), phpVersions...)
}

// GetDockerBaseImage returns the docker base image used by the given system.
func (system System) GetDockerBaseImage() string {
	version := DefaultPHPVersion
	if _, ok := phpVersionMap[system.PHP]; ok {
		version = system.PHP
	}
	return imagePrefix + version + imageSuffix
}

// GetIIPServerEnabled returns if the IIPServer was enabled
func (system System) GetIIPServerEnabled() string {
	if !system.IIPServer {
		return ""
	}
	return "1"
}

const (
	// Content Security Policy used by the internal server
	ContentSecurityPolicyNothing = "base-uri 'self'; default-src 'none';"

	// Content Security policy used by the distillery admin server
	ContentSecurityPolicyPanel = "base-uri 'self'; default-src 'self'; img-src 'self' data:; media-src 'none'; worker-src 'none'; frame-src 'none'; frame-ancestors 'none';"

	ContentSecurityPolicyPanelUnsafeScripts = ContentSecurityPolicyPanel + " script-src 'self' 'unsafe-inline';"
	ContentSecurityPolicyPanelUnsafeStyles  = ContentSecurityPolicyPanel + " style-src 'self' 'unsafe-inline';"
)

func ContentSecurityPolicyExamples() []string {
	return []string{
		ContentSecurityPolicyPanel,
	}
}
