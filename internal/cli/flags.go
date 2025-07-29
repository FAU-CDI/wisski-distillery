package cli

//spellchecker:words github wisski distillery internal wdlog
import "github.com/FAU-CDI/wisski-distillery/internal/wdlog"

// Flags are global flags for the wdcli executable.
type Flags struct {
	LogLevel         wdlog.Flag
	ConfigPath       string
	InternalInDocker bool
}
