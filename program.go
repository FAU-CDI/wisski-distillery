//spellchecker:words wisski distillery
package wisski_distillery

//spellchecker:words context signal user github wisski distillery internal bootstrap wdlog goprogram exit pkglib
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"os/user"

	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"go.tkw01536.de/goprogram"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/cgo"
)

// none of these are strictly needed, they're just around for convenience.
type wdcliEnv = *dis.Distillery
type wdcliParameters = cli.Params
type wdcliRequirements = cli.Requirements
type wdCliFlags = cli.Flags

type Program = goprogram.Program[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]
type Command = goprogram.Command[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]

// Context holds the context passed to any wdcli command.
//
// The context contains a reference to a "context.Context" as well as an IOStream.
//
// The context.Context holds a global context.
// It is initialized in the NewContext function below.
//
// It is cancelled if the user sends SIGINT or SIGKILL.
// Despite the context being a pseudo-global, it is passed to (almost) every function using the variable name "ctx".
//
// The IOStream is typically used in three ways:
//
// - Standard output is used to log events
// - Standard error is used to interactively display progress
// - Standard input is passed to a (few) interactive programs
//
// The standard output writer is passed directly into the context, see the "pkg/logging" package for conventions.
// Other parts are passed around (standard error using the variable name "progress") as required.
//
// The IOStream as a whole is only passed to functions that exist directly under cmd/.
//
// TODO(twiesing): The logging on the standard logger is still to be done.
type Context = goprogram.Context[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]
type Arguments = goprogram.Arguments[wdCliFlags]
type ContextCleanupFunc = goprogram.ContextCleanupFunc[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]
type Description = goprogram.Description[wdCliFlags, wdcliRequirements]

var GetContext = goprogram.GetContext[wdcliEnv, wdcliParameters, wdCliFlags, wdcliRequirements]

// an error when nor arguments are provided.
var errUserIsNotRoot = exit.NewErrorWithCode("this command has to be executed as root. the current user is not root", exit.ExitGeneralArguments)

const warnCGoEnabled = "Warning: This executable has been built with cgo enabled. This means certain commands may not work. \n"
const warnNoDeployWdcli = "Warning: Not using %q executable at %q. This might leave the distillery in an inconsistent state. \n"

func NewProgram() Program {
	return Program{
		BeforeCommand: func(context Context, command Command) error {
			// make sure that we are root!
			usr, err := user.Current()
			if err != nil || usr.Uid != "0" || usr.Gid != "0" {
				return errUserIsNotRoot
			}

			// warn about cgo!
			if cgo.Enabled {
				if _, err := context.Printf(warnCGoEnabled); err != nil {
					return fmt.Errorf("failed to print error: %w", err)
				}
			}

			// when not running inside docker and we need a distillery
			// then we should warn if we are not using the distillery executable.
			if dis := context.Environment; !context.Args.Flags.InternalInDocker && context.Description.Requirements.NeedsDistillery && !dis.Config.Paths.UsingDistilleryExecutable() {
				if _, err := context.EPrintf(warnNoDeployWdcli, bootstrap.Executable, dis.Config.Paths.ExecutablePath()); err != nil {
					return fmt.Errorf("failed to log error: %w", err)
				}
			}

			return nil
		},

		NewContext: func(params *wdcliParameters, parent context.Context) (context.Context, ContextCleanupFunc, error) {
			if params == nil {
				return parent, nil, nil
			}

			{
				context := GetContext(parent)
				logger := wdlog.New(os.Stderr, context.Args.Flags.LogLevel.Level())
				parent = wdlog.Set(parent, logger)
			}

			ctx, stop := signal.NotifyContext(parent, os.Interrupt, os.Kill)
			return ctx, func(context *Context) { stop() }, nil
		},

		NewEnvironment: func(params wdcliParameters, context Context) (e wdcliEnv, err error) {
			return dis.NewDistillery(params, context.Args.Flags, context.Description.Requirements)
		},
	}
}

// TODO: Check if / when we need umaskfree?
// TODO: Rework calls to WrapError
