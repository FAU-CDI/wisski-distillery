// Command wdcli implement the entry point for the wisski-distillery
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/cmd"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
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
	wdcli.Register(cmd.Backup)
	wdcli.Register(cmd.Cron)
	wdcli.Register(cmd.Monday)

	// servers
	wdcli.Register(cmd.Server)
	wdcli.Register(cmd.SSH)

	// status
	wdcli.Register(cmd.Status)
}

// an error when no arguments are provided.
var errNoArgumentsProvided = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "need at least one argument. use `wdcli license` to view licensing information",
}

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
	// just immediatly return a custom error message.
	if len(os.Args) == 1 {
		streams.Die(errNoArgumentsProvided)
		errNoArgumentsProvided.Return()
		return
	}

	// creat a new set of parameters
	// and then use them to execute the main command
	err := func() error {
		params, err := cli.ParamsFromEnv()
		if err != nil {
			return streams.Die(err)
		}

		return wdcli.Main(streams, params, os.Args[1:])
	}()

	// return the error to the user

	exit.AsError(err).Return()
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
