package models

// System represents system information.
// It is embedded into the instances struct by gorm.
type System struct {
	// NOTE(twiesing): Any changes here should be reflected in instance_{provision,rebuild}.html and remote/api.ts.
	PHP                string `gorm:"column:php;not null"`
	OpCacheDevelopment bool   `gorm:"column:opcache_devel;not null"`
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
