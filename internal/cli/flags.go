package cli

import "github.com/FAU-CDI/wisski-distillery/internal/wdlog"

// Flags are global flags for the wdcli executable
type Flags struct {
	//lint:ignore SA5008 required by the argument framework
	LogLevel   wdlog.Flag `short:"l" long:"loglevel" description:"log level" default:"info" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" choice:"panic"`
	ConfigPath string     `short:"c" long:"config" description:"path to distillery configuration file"`

	InternalInDocker bool `long:"internal-in-docker" description:"internal flag to signal the shell that it is running inside a docker stack belonging to the distillery"`
}
