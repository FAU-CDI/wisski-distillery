package cli

import "go.tkw01536.de/pkglib/exit"

const (
	// ExitZero indicates that no error occurred.
	// It is the zero value of type ExitCode.
	ExitZero exit.ExitCode = 0

	// ExitGeneric indicates a generic error occurred within this invocation.
	// This typically implies a subcommand-specific behavior wants to return failure to the caller.
	ExitGeneric exit.ExitCode = 1

	// ExitUnknownCommand indicates that the user attempted to call a subcommand that is not defined.
	ExitUnknownCommand exit.ExitCode = 2

	// ExitGeneralArguments indicates that the user attempted to pass invalid general arguments to the program.
	ExitGeneralArguments exit.ExitCode = 3
	// ExitCommandArguments indicates that the user attempted to pass invalid command-specific arguments to a subcommand.
	ExitCommandArguments exit.ExitCode = 4

	// ExitContext indicates an error with the underlying command context.
	ExitContext exit.ExitCode = 254

	// ExitPanic indicates that the go code called panic() inside the execution of the current program.
	// This typically implies a bug inside a program.
	ExitPanic exit.ExitCode = 255
)
