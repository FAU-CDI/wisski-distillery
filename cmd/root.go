package cmd

//spellchecker:words context user github wisski distillery internal wdlog cobra pkglib exit stream
import (
	"context"
	"fmt"
	"os/user"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/cgo"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/stream"
)

//spellchecker:words contextcheck unsynced pflags shellrc GGROOT
var (
	errInvalidFlags        = exit.NewErrorWithCode("unknown flags passed", cli.ExitGeneralArguments)
	errNoArgumentsProvided = exit.NewErrorWithCode("need at least one argument. use `wdcli license` to view licensing information", cli.ExitGeneralArguments)
	errUserIsNotRoot       = exit.NewErrorWithCode("this command has to be executed as root. the current user is not root", cli.ExitGeneralArguments)
)

const warnCGoEnabled = "Warning: This executable has been built with cgo enabled. This means certain commands may not work. \n"

// Command returns the main wdcli command
//
//nolint:contextcheck // don't need to pass down the context
func NewCommand(ctx context.Context, parameters cli.Params) *cobra.Command {
	var flags cli.Flags

	root := &cobra.Command{
		Use:   "wdcli",
		Short: "A command line tool for the wisski-distillery",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger := wdlog.New(cmd.ErrOrStderr(), flags.LogLevel.Level())
			cmd.SetContext(wdlog.Set(cmd.Context(), logger))

			// make sure that we are root!
			usr, err := user.Current()
			if err != nil || usr.Uid != "0" || usr.Gid != "0" {
				return errUserIsNotRoot
			}

			// warn about cgo!
			if cgo.Enabled {
				if _, err := fmt.Fprint(cmd.ErrOrStderr(), warnCGoEnabled); err != nil {
					return fmt.Errorf("failed to print error: %w", err)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return errNoArgumentsProvided
		},
	}

	// setup flags
	{
		pflags := root.PersistentFlags()

		pflags.StringVarP((*string)(&flags.LogLevel), "loglevel", "l", "info", "log level")
		pflags.StringVarP(&flags.ConfigPath, "config", "c", "", "path to distillery configuration file")
		pflags.BoolVar(&flags.InternalInDocker, "internal-in-docker", false, "internal flag to signal the shell that it is running inside a docker stack belonging to the distillery")
	}

	root.SetContext(ctx)
	root.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%w: %w", errInvalidFlags, err)
	})

	root.SilenceErrors = true
	root.SilenceUsage = true

	cli.SetFlags(root, &flags)
	cli.SetParameters(root, &parameters)

	// add all the commands
	root.AddCommand(
		// self commands
		NewConfigCommand(),

		// setup commands
		NewBootstrapCommand(),
		NewSystemUpdateCommand(),
		NewSystemPauseCommand(),

		// sql commands
		NewMysqlCommand(),
		NewMakeMysqlAccountCommand(),

		// instance setup and teardown
		NewProvisionCommand(),
		NewPurgeCommand(),
		NewReserveCommand(),
		NewRebuildCommand(),

		// instance management
		NewLsCommand(),
		NewInfoCommand(),
		NewInstanceLockCommand(),
		NewInstancePauseCommand(),
		NewInstanceLogCommand(),

		// instance tasks
		NewShellCommand(),
		NewBlindUpdateCommand(),
		NewUpdatePrefixConfigCommand(), // TODO: Move into post-instance configuration

		NewPathbuildersCommand(),
		NewPrefixesCommand(),
		NewDrupalSettingCommand(),
		NewDrupalUserCommand(),

		// distillery auth
		NewDisUserCommand(),
		NewDisGrantCommand(),
		NewDisSSHCommand(),

		// backup & cron
		NewSnapshotCommand(),
		NewRebuildTSCommand(),
		NewBackupCommand(),
		NewSnapshotRestoreCommand(),
		NewBackupsPruneCommand(),
		NewCronCommand(),
		NewMondayCommand(),

		// servers
		NewServerCommand(),
		NewSSHCommand(),

		// status
		NewStatusCommand(),

		NewMakeBlockCommand(),

		// self commands
		NewLicenseCommand(),
	)

	// wrap all the argument errors
	var wrapAllArgs func(cmd *cobra.Command)
	wrapAllArgs = func(cmd *cobra.Command) {
		cmd.Args = wrapArgs(cmd.Args)
		for _, child := range cmd.Commands() {
			wrapAllArgs(child)
		}
	}
	wrapAllArgs(root)

	// setup more flags

	return root
}

var errInvalidArguments = exit.NewErrorWithCode("invalid arguments passed", cli.ExitCommandArguments)

// wrapArgs wraps a [cobra.PositionalArgs] error with the given error.
// The wrapping occurs by calling [fmt.Errorf] with a string of "%w: %w" and [errInvalidArguments].
// If pos is nil, it is passed through as-is.
func wrapArgs(pos cobra.PositionalArgs) cobra.PositionalArgs {
	if pos == nil {
		return pos
	}

	return func(cmd *cobra.Command, args []string) error {
		err := pos(cmd, args)
		if err == nil {
			return nil
		}
		return fmt.Errorf("%w: %w", errInvalidArguments, err)
	}
}

// streamFromCommand returns a stream.IOStream from the given command.
func streamFromCommand(cmd *cobra.Command) stream.IOStream {
	return stream.IOStream{
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.ErrOrStderr(),
		Stdin:  cmd.InOrStdin(),
	}
}
