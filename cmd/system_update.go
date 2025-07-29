package cmd

//spellchecker:words sync github wisski distillery internal component execx logging goprogram exit parser pkglib errorsx umaskfree status
import (
	"fmt"
	"io"
	"sync"

	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
	"go.tkw01536.de/pkglib/fsx/umaskfree"
	"go.tkw01536.de/pkglib/status"
)

func NewSystemUpdateCommand() *cobra.Command {
	impl := new(systemupdate)

	cmd := &cobra.Command{
		Use:     "system_update GRAPHDB_ZIP",
		Short:   "installs and updates components of the distillery system",
		Args:    cobra.ExactArgs(1),
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVar(&impl.InstallDocker, "install-docker", false, "try to automatically install docker. assumes 'apt-get' as a package manager")

	return cmd
}

type systemupdate struct {
	InstallDocker bool
	Positionals   struct {
		GraphdbZip string
	}
}

func (s *systemupdate) ParseArgs(cmd *cobra.Command, args []string) error {
	s.Positionals.GraphdbZip = args[0]

	isFile, err := fsx.IsRegular(s.Positionals.GraphdbZip, true)
	if err != nil {
		return fmt.Errorf("failed to check for regular file: %w", err)
	}

	if !isFile {
		return fmt.Errorf("%q: %w", s.Positionals.GraphdbZip, exit.NewErrorWithCode("does not exist", exit.ExitCommandArguments))
	}
	return nil
}

var (
	errSystemUpdateFailedToLog           = exit.NewErrorWithCode("failed to log message", exit.ExitGeneric)
	errBoostrapFailedToCreateDirectory   = exit.NewErrorWithCode("failed to create directory", exit.ExitGeneric)
	errBootstrapComponent                = exit.NewErrorWithCode("unable to bootstrap", exit.ExitGeneric)
	errNetworkCreateFailed               = exit.NewErrorWithCode("unable to create docker network", exit.ExitGeneric)
	errSystemUpdateFailedToPing          = exit.NewErrorWithCode("failed to ping docker client", exit.ExitGeneric)
	errSystemUpdateDockerClient          = exit.NewErrorWithCode("failed to create docker client", exit.ExitGeneric)
	errSystemUpdateFailedStackUpdate     = exit.NewErrorWithCode("failed to perform stack updates", exit.ExitGeneric)
	errSystemUpdateFailedComponentUpdate = exit.NewErrorWithCode("failed to perform component updates", exit.ExitGeneric)
)

func (s *systemupdate) Exec(cmd *cobra.Command, args []string) (e error) {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
		FailOnCgo:       true,
	})
	if err != nil {
		return fmt.Errorf("failed to get distillery: %w", err)
	}

	// create all the other directories
	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Ensuring distillery installation directories exist"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
	}
	for _, d := range []string{
		dis.Config.Paths.Root,
		dis.Instances().Path(),
		dis.Exporter().StagingPath(),
		dis.Exporter().ArchivePath(),
		dis.Templating().CustomAssetsPath(),
	} {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), d)
		if err := umaskfree.MkdirAll(d, umaskfree.DefaultDirPerm); err != nil {
			return fmt.Errorf("%q: %w: %w", d, errBoostrapFailedToCreateDirectory, err)
		}
	}

	if s.InstallDocker {
		// install system updates
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Updating Operating System Packages"); err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
		}
		if err := s.mustExec(cmd, dis, "", "apt-get", "update"); err != nil {
			return err
		}
		if err := s.mustExec(cmd, dis, "", "apt-get", "upgrade", "-y"); err != nil {
			return err
		}

		// install docker
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Installing / Updating Docker"); err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
		}
		if err := s.mustExec(cmd, dis, "", "apt-get", "install", "curl"); err != nil {
			return err
		}
		// TODO: Download directly
		if err := s.mustExec(cmd, dis, "", "/bin/sh", "-c", "curl -fsSL https://get.docker.com -o - | /bin/sh"); err != nil {
			return err
		}
	}

	// create the docker client!
	client, err := dis.Docker().NewClient()
	if err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateDockerClient, err)
	}
	defer errorsx.Close(client, &e, "client")

	// check that the docker api is available
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Checking that the 'docker' api is reachable"); err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
		}

		ping, err := client.Ping(cmd.Context())
		if err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToPing, err)
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "API Version:     %s (experimental: %t)\nBuilder Version: %s\n", ping.APIVersion, ping.Experimental, ping.BuilderVersion)
	}

	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Checking that 'docker compose' is available"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if err := s.mustExec(cmd, dis, "", "docker", "compose", "version"); err != nil {
			return err
		}
	}

	// create the docker networks
	{
		if _, err := logging.LogMessage(cmd.ErrOrStderr(), "Configuring docker networks"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		for _, name := range dis.Config.Docker.Networks() {
			id, existed, err := client.NetworkCreate(cmd.Context(), name)
			if err != nil {
				return fmt.Errorf("%w: %w", errNetworkCreateFailed, err)
			}
			if existed {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Network %s (id %s) already existed\n", name, id)
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Network %s (id %s) created\n", name, id)
			}
		}
	}

	// install and update the various stacks!
	ctx := component.InstallationContext{
		"graphdb.zip": s.Positionals.GraphdbZip,
	}

	var updated = make(map[string]struct{})
	var updateMutex sync.Mutex

	if err := logging.LogOperation(func() error {
		return status.RunErrorGroup(cmd.ErrOrStderr(), status.Group[component.Installable, error]{
			PrefixString: func(item component.Installable, index int) string {
				return fmt.Sprintf("[update %q]: ", item.Name())
			},
			PrefixAlign: true,

			Handler: func(item component.Installable, index int, writer io.Writer) (e error) {
				stack, err := item.OpenStack()
				if err != nil {
					return fmt.Errorf("failed open stack: %w", err)
				}
				defer errorsx.Close(stack, &e, "stack")

				if err := stack.Install(cmd.Context(), writer, item.Context(ctx)); err != nil {
					return fmt.Errorf("failed to install stack: %w", err)
				}

				if err := stack.Update(cmd.Context(), writer, true); err != nil {
					return fmt.Errorf("failed to update stack: %w", err)
				}

				ud, ok := item.(component.Updatable)
				if !ok {
					return nil
				}

				defer func() {
					updateMutex.Lock()
					defer updateMutex.Unlock()
					updated[item.ID()] = struct{}{}
				}()

				return ud.Update(cmd.Context(), writer)
			},
		}, dis.Installable())
	}, cmd.ErrOrStderr(), "Performing Stack Updates"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedStackUpdate, err)
	}

	if err := logging.LogOperation(func() error {
		for _, item := range dis.Updatable() {
			name := item.Name()
			if err := logging.LogOperation(func() error {
				_, ok := updated[item.ID()]
				if ok {
					_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Already updated")
					return nil
				}
				return item.Update(cmd.Context(), cmd.ErrOrStderr())
			}, cmd.ErrOrStderr(), "Updating Component: %s", name); err != nil {
				return fmt.Errorf("%w: %q: %w", errBootstrapComponent, name, err)
			}
		}
		return nil
	}, cmd.ErrOrStderr(), "Performing Component Updates"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedComponentUpdate, err)
	}

	if _, err := logging.LogMessage(cmd.ErrOrStderr(), "System has been updated"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
	}
	return nil
}

//nolint:unparam
func (s *systemupdate) mustExec(cmd *cobra.Command, dis *dis.Distillery, workdir string, exe string, argv ...string) error {
	if workdir == "" {
		workdir = dis.Config.Paths.Root
	}
	code := execx.Exec(cmd.Context(), streamFromCommand(cmd), workdir, exe, argv...)()

	if code := exit.Code(code); code != 0 {
		return exit.NewErrorWithCode(fmt.Sprintf("process exited with code %d", code), code)
	}
	return nil
}
