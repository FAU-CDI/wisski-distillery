package cli

//spellchecker:words github wisski distillery internal wdlog
import "github.com/FAU-CDI/wisski-distillery/internal/wdlog"

// Flags are global flags for the wdcli executable.
type Flags struct {
	//lint:ignore SA5008 required by the argument framework
	//nolint:staticcheck
	LogLevel   wdlog.Flag `choice:"trace"                                      choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"fatal" choice:"panic" default:"info" description:"log level" long:"loglevel" short:"l"`
	ConfigPath string     `description:"path to distillery configuration file" long:"config"  short:"c"`

	InternalInDocker bool `description:"internal flag to signal the shell that it is running inside a docker stack belonging to the distillery" long:"internal-in-docker"`
}
