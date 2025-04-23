package cmd

//spellchecker:words sync github wisski distillery internal component execx logging goprogram exit parser pkglib umaskfree status
import (
	"fmt"
	"io"
	"sync"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/errwrap"
	"github.com/FAU-CDI/wisski-distillery/pkg/execx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/fsx"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
	"github.com/tkw1536/pkglib/status"
)

// SystemUpdate is the 'system_update' command.
var SystemUpdate wisski_distillery.Command = systemupdate{}

type systemupdate struct {
	InstallDocker bool "short:\"a\" long:\"install-docker\" description:\"try to automatically install docker. assumes `apt-get` as a package manager\""
	Positionals   struct {
		GraphdbZip string `positional-arg-name:"PATH_TO_GRAPHDB_ZIP" required:"1-1" description:"path to the graphdb.zip file"`
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

var errNoGraphDBZip = exit.Error{
	Message:  "%q does not exist",
	ExitCode: exit.ExitCommandArguments,
}

func (s systemupdate) AfterParse() error {
	isFile, err := fsx.IsRegular(s.Positionals.GraphdbZip, true)
	if err != nil {
		return err
	}

	if !isFile {
		return errNoGraphDBZip.WithMessageF(s.Positionals.GraphdbZip)
	}
	return nil
}

var errBoostrapFailedToCreateDirectory = exit.Error{
	Message:  "failed to create directory %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapComponent = exit.Error{
	Message:  "unable to bootstrap %s",
	ExitCode: exit.ExitGeneric,
}

var errDockerUnreachable = exit.Error{
	Message:  "unable to reach docker api",
	ExitCode: exit.ExitGeneric,
}

var errNetworkCreateFailed = exit.Error{
	Message:  "unable to create docker network",
	ExitCode: exit.ExitGeneric,
}

var errSystemUpdateGeneric = exit.Error{
	Message:  "generic system update error",
	ExitCode: exit.ExitGeneric,
}

func (si systemupdate) Run(context wisski_distillery.Context) (err error) {
	defer errwrap.DeferWrap(errSystemUpdateGeneric, &err)

	dis := context.Environment

	// create all the other directories
	if _, err := logging.LogMessage(context.Stderr, "Ensuring distillery installation directories exist"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	for _, d := range []string{
		dis.Config.Paths.Root,
		dis.Instances().Path(),
		dis.Exporter().StagingPath(),
		dis.Exporter().ArchivePath(),
		dis.Templating().CustomAssetsPath(),
	} {
		context.Println(d)
		if err := umaskfree.MkdirAll(d, umaskfree.DefaultDirPerm); err != nil {
			return errBoostrapFailedToCreateDirectory.WithMessageF(d).WrapError(err)
		}
	}

	if si.InstallDocker {
		// install system updates
		if _, err := logging.LogMessage(context.Stderr, "Updating Operating System Packages"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if err := si.mustExec(context, "", "apt-get", "update"); err != nil {
			return err
		}
		if err := si.mustExec(context, "", "apt-get", "upgrade", "-y"); err != nil {
			return err
		}

		// install docker
		if _, err := logging.LogMessage(context.Stderr, "Installing / Updating Docker"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if err := si.mustExec(context, "", "apt-get", "install", "curl"); err != nil {
			return err
		}
		// TODO: Download directly
		if err := si.mustExec(context, "", "/bin/sh", "-c", "curl -fsSL https://get.docker.com -o - | /bin/sh"); err != nil {
			return err
		}
	}

	// check that the docker api is available
	{
		if _, err := logging.LogMessage(context.Stderr, "Checking that the 'docker' api is reachable"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		ping, err := dis.Docker().Ping(context.Context)
		if err != nil {
			return errDockerUnreachable.WrapError(err)
		}
		context.Printf("API Version:     %s (experimental: %t)\nBuilder Version: %s\n", ping.APIVersion, ping.Experimental, ping.BuilderVersion)
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
			id, existed, err := dis.Docker().CreateNetwork(context.Context, name)
			if err != nil {
				return errNetworkCreateFailed.WrapError(err)
			}
			if existed {
				context.Printf("Network %s (id %s) already existed\n", name, id)
			} else {
				context.Printf("Network %s (id %s) created\n", name, id)
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

			Handler: func(item component.Installable, index int, writer io.Writer) error {
				stack := item.Stack()

				if err := stack.Install(context.Context, writer, item.Context(ctx)); err != nil {
					return err
				}

				if err := stack.Update(context.Context, writer, true); err != nil {
					return err
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
		return err
	}

	if err := logging.LogOperation(func() error {
		for _, item := range dis.Updatable() {
			name := item.Name()
			if err := logging.LogOperation(func() error {
				_, ok := updated[item.ID()]
				if ok {
					context.Println("Already updated")
					return nil
				}
				return item.Update(context.Context, context.Stderr)
			}, context.Stderr, "Updating Component: %s", name); err != nil {
				return errBootstrapComponent.WithMessageF(name).WrapError(err)
			}
		}
		return nil
	}, context.Stderr, "Performing Component Updates"); err != nil {
		return err
	}

	if _, err := logging.LogMessage(context.Stderr, "System has been updated"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	return nil
}

var errMustExecFailed = exit.Error{
	Message: "process exited with code %d",
}

// If it does not, returns errMustExecFailed.
//
//nolint:unparam
func (si systemupdate) mustExec(context wisski_distillery.Context, workdir string, exe string, argv ...string) error {
	dis := context.Environment
	if workdir == "" {
		workdir = dis.Config.Paths.Root
	}
	code := execx.Exec(context.Context, context.IOStream, workdir, exe, argv...)()

	if code := exit.Code(code); code != 0 {
		err := errMustExecFailed.WithMessageF(code)
		err.ExitCode = code
		return err
	}
	return nil
}
