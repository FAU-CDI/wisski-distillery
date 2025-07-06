// Command wdcli implement the entry point for the wisski-distillery
//
//spellchecker:words main
package main

//spellchecker:words runtime debug github wisski distillery internal goprogram exit pkglib stream
import (
	"fmt"
	"os"
	"runtime/debug"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/cmd"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/stream"
)

var wdcli = wisski_distillery.NewProgram()

func init() {
	// self commands
	wdcli.Register(cmd.Config)
	wdcli.Register(cmd.License)

	// setup commands
	wdcli.Register(cmd.Bootstrap)
	wdcli.Register(cmd.SystemUpdate)
	wdcli.Register(cmd.SystemPause)

	// sql commands
	wdcli.Register(cmd.Mysql)
	wdcli.Register(cmd.MakeMysqlAccount)

	// instance setup and teardown
	wdcli.Register(cmd.Provision)
	wdcli.Register(cmd.Purge)
	wdcli.Register(cmd.Reserve)
	wdcli.Register(cmd.Rebuild)

	// instance management
	wdcli.Register(cmd.Ls)
	wdcli.Register(cmd.Info)
	wdcli.Register(cmd.InstanceLock)
	wdcli.Register(cmd.InstancePause)
	wdcli.Register(cmd.InstanceLog)

	// instance tasks
	wdcli.Register(cmd.Shell)
	wdcli.Register(cmd.BlindUpdate)
	wdcli.Register(cmd.UpdatePrefixConfig) // TODO: Move into post-instance configuration

	wdcli.Register(cmd.Pathbuilders)
	wdcli.Register(cmd.Prefixes)
	wdcli.Register(cmd.DrupalSetting)
	wdcli.Register(cmd.DrupalUser)

	// distillery auth
	wdcli.Register(cmd.DisUser)
	wdcli.Register(cmd.DisGrant)
	wdcli.Register(cmd.DisSSH)

	// backup & cron
	wdcli.Register(cmd.Snapshot)
	wdcli.Register(cmd.RebuildTS)
	wdcli.Register(cmd.Backup)
	wdcli.Register(cmd.BackupsPrune)
	wdcli.Register(cmd.Cron)
	wdcli.Register(cmd.Monday)

	// servers
	wdcli.Register(cmd.Server)
	wdcli.Register(cmd.SSH)

	// status
	wdcli.Register(cmd.Status)

	wdcli.Register(cmd.MakeBlock)
}

// an error when no arguments are provided.
var errNoArgumentsProvided = exit.NewErrorWithCode("need at least one argument. use `wdcli license` to view licensing information", exit.ExitGeneralArguments)

func main() {
	// recover from calls to panic(), and exit the program appropriatly.
	// This has to be in the main() function because any of the library functions might be broken.
	// For this reason, as few ggman functions as possible are used here; just stuff from the top-level ggman package.
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, fatalPanicMessage, err)
			debug.PrintStack()
			exit.ExitPanic.Return()
		}
	}()

	streams := stream.FromEnv()

	// when there are no arguments then parsing argument *will* fail
	//
	// we don't need to even bother with the rest of the program
	// just immediately return a custom error message.
	if len(os.Args) == 1 {
		_ = exit.Die(streams, errNoArgumentsProvided) // returned below anyways
		code, _ := exit.CodeFromError(errNoArgumentsProvided)
		code.Return()
		return
	}

	// creat a new set of parameters
	// and then use them to execute the main command
	code, _ := exit.CodeFromError(func() error {
		params, err := cli.ParamsFromEnv()
		if err != nil {
			return exit.Die(streams, err)
		}

		return wdcli.Main(streams, params, os.Args[1:])
	}(),
	)
	code.Return()
}

const fatalPanicMessage = `Fatal Error: Panic

The wdcli program panicked and had to abort execution. This is usually
indicative of a bug. If this occurs repeatedly you might want to consider
filing an issue in the issue tracker at:

https://github.com/FAU-CDI/wisski-distillery/issues

Below is debug information that might help the developers track down what
happened.

panic: %v
`
