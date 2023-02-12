// Package bootstrap implements the core of the WissKI Distillery and the wdcli executable.
// It does not depend on any other packages.
package bootstrap

import _ "embed"

// TODO: This package should be split up into a true bootstrap component, and something else.

// BaseDirectoryDefault is the default deploy directory to load the distillery from.
const BaseDirectoryDefault = "/var/www/deploy"

// Executable is the name of the 'wdcli' executable.
// It should be located inside the deployment directory.
const Executable = "wdcli"

// ConfigFile is the name of the config file.
// It should be located inside the deployment directory.
const ConfigFile = "distillery.yaml"

// OverridesJSON is the name of the json overrides file.
// It should be located inside the deployment directory.
const OverridesJSON = "overrides.json"

// DefaultOverridesJSON contains a template for a new 'overrides.json' file
//
//go:embed overrides.json
var DefaultOverridesJSON []byte

// ResolverBlockTXT is the name of the resolver blocked prefix file.
// It should be located inside the deployment directory.
const ResolverBlockedTXT = "resolver-blocked.txt"

// ResolverBlockTXT contains a template for 'resolver-blocked' file
//
//go:embed resolver-blocked.txt
var DefaultResolverBlockedTXT []byte
