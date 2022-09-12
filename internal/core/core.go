// Package core implements the core of the WissKI Distillery and the wdcli executable.
// It does not depend on any other packages.
package core

// BaseDirectoryDefault is the default deploy directory to load the distillery from.
const BaseDirectoryDefault = "/var/www/deploy"

// Executable is the name of the 'wdcli' executable.
// It should be located inside the deployment directory.
const Executable = "wdcli"

// Config file is the name of the config file.
// It should be located inside the deployment directory.
const ConfigFile = ".env"
