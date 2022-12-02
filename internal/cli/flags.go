package cli

import "github.com/rs/zerolog"

// Flags are global flags for the wdcli executable
type Flags struct {
	LogLevel   LogLevelString `short:"l" long:"loglevel" description:"Log level to use for logger" default:"info" choice:"trace" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" choice:"panic"`
	ConfigPath string         `short:"c" long:"config" description:"Path to distillery configuration file"`

	InternalInDocker bool `long:"internal-in-docker" description:"Internal Flag to signal the shell that it is running inside a docker stack belonging to the distillery"`
}

type LogLevelString string

func (ls LogLevelString) Level() zerolog.Level {
	level, err := zerolog.ParseLevel(string(ls))
	if err != nil {
		return zerolog.InfoLevel
	}
	return level
}
