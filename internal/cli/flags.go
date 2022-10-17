package cli

// Flags are global flags for the wdcli executable
type Flags struct {
	ConfigPath string `short:"c" long:"config" description:"Path to distillery configuration file"`

	InternalInDocker bool `long:"internal-in-docker" description:"Internal Flag to signal the shell that it is running inside a docker stack belonging to the distillery"`
}
