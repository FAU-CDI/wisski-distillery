package wisski_distillery

import (
	"os/user"

	"github.com/FAU-CDI/wisski-distillery/core"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/tkw1536/goprogram"
	"github.com/tkw1536/goprogram/exit"
)

// these define the ggman-specific program types
// none of these are strictly needed, they're just around for convenience
type wdcliEnv = *env.Distillery
type wdcliParameters = core.Params
type wdcliRequirements = core.Requirements
type wdCliFlags = core.Flags

type Program = goprogram.Program[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]
type Command = goprogram.Command[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]
type Context = goprogram.Context[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]
type Arguments = goprogram.Arguments[wdCliFlags]
type Description = goprogram.Description[wdCliFlags, wdcliRequirements]

// an error when nor arguments are provided.
var errUserIsNotRoot = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "This command has to be executed as root. The current user is not root.",
}

var warnNoDeployWdcli = "Warning: Not using %q executable at %q. This might leave the distillery in an inconsistent state. \n"

func NewProgram() Program {
	return Program{
		BeforeCommand: func(context Context, command Command) error {
			// make sure that we are root!
			usr, err := user.Current()
			if err != nil || usr.Uid != "0" || usr.Gid != "0" {
				return errUserIsNotRoot
			}

			// warn when not using the distillery excutable
			if context.Description.Requirements.NeedsDistillery {
				dis := context.Environment
				if !dis.UsingDistilleryExecutable() {
					context.EPrintf(warnNoDeployWdcli, core.Executable, dis.ExecutablePath())
				}
			}

			return nil
		},

		NewEnvironment: func(params wdcliParameters, context Context) (e wdcliEnv, err error) {
			return env.NewDistillery(params, context.Args.Flags, context.Description.Requirements)
		},
	}
}
