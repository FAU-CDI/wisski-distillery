package cmd

//spellchecker:words sync github wisski distillery internal component execx logging goprogram exit parser pkglib errorsx umaskfree status
import (
	"fmt"
	"io"
	"sync"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/fsx"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/status"
)

// SystemUpdate is the 'system_update' command.
var SystemUpdate wisski_distillery.Command = systemupdate{}

type systemupdate struct {
	InstallDocker bool `description:"try to automatically install docker. assumes 'apt-get' as a package manager" long:"install-docker" short:"a"`
	Positionals   struct {
		GraphdbZip string `description:"path to the graphdb.zip file" positional-arg-name:"PATH_TO_GRAPHDB_ZIP" required:"1-1"`
	} `positional-args:"true"`
}

func (systemupdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
			FailOnCgo:       true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "system_update",
		Description: "installs and updates components of the distillery system",
	}
}

var errNoGraphDBZip = exit.NewErrorWithCode("does not exist", exit.ExitCommandArguments)

func (s systemupdate) AfterParse() error {
	isFile, err := fsx.IsRegular(s.Positionals.GraphdbZip, true)
	if err != nil {
		return fmt.Errorf("failed to check for regular file: %w", err)
	}

	if !isFile {
		return fmt.Errorf("%q: %w", s.Positionals.GraphdbZip, errNoGraphDBZip)
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

func (si systemupdate) Run(context wisski_distillery.Context) (e error) {
	dis := context.Environment

	// create all the other directories
	if _, err := logging.LogMessage(context.Stderr, "Ensuring distillery installation directories exist"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
	}
	for _, d := range []string{
		dis.Config.Paths.Root,
		dis.Instances().Path(),
		dis.Exporter().StagingPath(),
		dis.Exporter().ArchivePath(),
		dis.Templating().CustomAssetsPath(),
	} {
		_, _ = context.Println(d)
		if err := umaskfree.MkdirAll(d, umaskfree.DefaultDirPerm); err != nil {
			return fmt.Errorf("%q: %w: %w", d, errBoostrapFailedToCreateDirectory, err)
		}
	}

	if si.InstallDocker {
		// install system updates
		if _, err := logging.LogMessage(context.Stderr, "Updating Operating System Packages"); err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
		}
		if err := si.mustExec(context, "", "apt-get", "update"); err != nil {
			return err
		}
		if err := si.mustExec(context, "", "apt-get", "upgrade", "-y"); err != nil {
			return err
		}

		// install docker
		if _, err := logging.LogMessage(context.Stderr, "Installing / Updating Docker"); err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
		}
		if err := si.mustExec(context, "", "apt-get", "install", "curl"); err != nil {
			return err
		}
		// TODO: Download directly
		if err := si.mustExec(context, "", "/bin/sh", "-c", "curl -fsSL https://get.docker.com -o - | /bin/sh"); err != nil {
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
		if _, err := logging.LogMessage(context.Stderr, "Checking that the 'docker' api is reachable"); err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
		}

		ping, err := client.Ping(context.Context)
		if err != nil {
			return fmt.Errorf("%w: %w", errSystemUpdateFailedToPing, err)
		}
		_, _ = context.Printf("API Version:     %s (experimental: %t)\nBuilder Version: %s\n", ping.APIVersion, ping.Experimental, ping.BuilderVersion)
	}

	{
		if _, err := logging.LogMessage(context.Stderr, "Checking that 'docker compose' is available"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if err := si.mustExec(context, "", "docker", "compose", "version"); err != nil {
			return err
		}
	}

	// create the docker networks
	{
		if _, err := logging.LogMessage(context.Stderr, "Configuring docker networks"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}

		for _, name := range dis.Config.Docker.Networks() {
			id, existed, err := client.NetworkCreate(context.Context, name)
			if err != nil {
				return fmt.Errorf("%w: %w", errNetworkCreateFailed, err)
			}
			if existed {
				_, _ = context.Printf("Network %s (id %s) already existed\n", name, id)
			} else {
				_, _ = context.Printf("Network %s (id %s) created\n", name, id)
			}
		}
	}

	// install and update the various stacks!
	ctx := component.InstallationContext{
		"graphdb.zip": si.Positionals.GraphdbZip,
	}

	var updated = make(map[string]struct{})
	var updateMutex sync.Mutex

	if err := logging.LogOperation(func() error {
		return status.RunErrorGroup(context.Stderr, status.Group[component.Installable, error]{
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

				if err := stack.Install(context.Context, writer, item.Context(ctx)); err != nil {
					return fmt.Errorf("failed to install stack: %w", err)
				}

				if err := stack.Update(context.Context, writer, true); err != nil {
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

				return ud.Update(context.Context, writer)
			},
		}, dis.Installable())
	}, context.Stderr, "Performing Stack Updates"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedStackUpdate, err)
	}

	if err := logging.LogOperation(func() error {
		for _, item := range dis.Updatable() {
			name := item.Name()
			if err := logging.LogOperation(func() error {
				_, ok := updated[item.ID()]
				if ok {
					_, _ = context.Println("Already updated")
					return nil
				}
				return item.Update(context.Context, context.Stderr)
			}, context.Stderr, "Updating Component: %s", name); err != nil {
				return fmt.Errorf("%w: %q: %w", errBootstrapComponent, name, err)
			}
		}
		return nil
	}, context.Stderr, "Performing Component Updates"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedComponentUpdate, err)
	}

	if _, err := logging.LogMessage(context.Stderr, "System has been updated"); err != nil {
		return fmt.Errorf("%w: %w", errSystemUpdateFailedToLog, err)
	}
	return nil
}

//nolint:unparam
func (si systemupdate) mustExec(context wisski_distillery.Context, workdir string, exe string, argv ...string) error {
	dis := context.Environment
	if workdir == "" {
		workdir = dis.Config.Paths.Root
	}
	code := execx.Exec(context.Context, context.IOStream, workdir, exe, argv...)()

	if code := exit.Code(code); code != 0 {
		return exit.NewErrorWithCode(fmt.Sprintf("process exited with code %d", code), code)
	}
	return nil
}
