package core

import _ "embed"

// DefaultOverridesJSON contains a template for a new 'overrides.json' file
//go:embed bootstrap/overrides.json
var DefaultOverridesJSON []byte

// DefaultAuthorizedKeys contains a template for a new 'global_authorized_keys' file
//go:embed bootstrap/global_authorized_keys
var DefaultAuthorizedKeys []byte

// ConfigFileTemplate contains a template for a new configuration file
//go:embed bootstrap/env
var ConfigFileTemplate []byte
