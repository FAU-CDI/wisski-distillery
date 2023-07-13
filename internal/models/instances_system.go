package models

// System represents system information.
// It is embedded into the instances struct by gorm.
type System struct {
	// NOTE(twiesing): Any changes here should be reflected in instance_{provision,rebuild}.html and remote/api.ts.
	PHP                string `gorm:"column:php;not null"`           // php version to use
	OpCacheDevelopment bool   `gorm:"column:opcache_devel;not null"` // opcache development

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
	phpVersions   = []string{"8.0", "8.1", "8.2"}
	phpVersionMap = (func() map[string]struct{} {
		m := make(map[string]struct{}, len(phpVersions))
		for _, v := range phpVersions {
			m[v] = struct{}{}
		}
		return m
	})()
)

// DefaultPHPVersion is the default php version
const DefaultPHPVersion = "8.1"

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

const (
	// Content Security Policy used by the internal server
	ContentSecurityPolicyNothing = "base-uri 'self'; default-src 'none';"

	// Content Security policy used by the distillery admin server
	ContentSecurityPolicyDistilery = "base-uri 'self'; default-src 'self'; img-src 'self' data:; media-src 'none'; worker-src 'none'; frame-src 'none'; frame-ancestors 'none';"
)

func ContentSecurityPolicyExamples() []string {
	return []string{
		ContentSecurityPolicyDistilery,
	}
}
