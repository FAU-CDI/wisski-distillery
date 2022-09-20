package cmd

import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/environment"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// SystemUpdate is the 'system_update' command
var SystemUpdate wisski_distillery.Command = systemupdate{}

type systemupdate struct {
	SkipCoreUpdates bool `short:"s" long:"skip-core-updates" description:"Skip applying operating system and other core system updates"`
	Positionals     struct {
		GraphdbZip string `positional-arg-name:"PATH_TO_GRAPHDB_ZIP" required:"1-1" description:"path to the graphdb.zip file"`
	} `positional-args:"true"`
}

func (systemupdate) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: true,
		},
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},
		Command:     "system_update",
		Description: "Installs and Update Components of the WissKI Distillery System",
	}
}

var errNoGraphDBZip = exit.Error{
	Message:  "%s does not exist",
	ExitCode: exit.ExitCommandArguments,
}

func (s systemupdate) AfterParse() error {
	// TODO: Use a generic environment here!
	if !fsx.IsFile(environment.Native{}, s.Positionals.GraphdbZip) {
		return errNoGraphDBZip.WithMessageF(s.Positionals.GraphdbZip)
	}
	return nil
}

var errBoostrapFailedToCreateDirectory = exit.Error{
	Message:  "failed to create directory %s: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapComponent = exit.Error{
	Message:  "Unable to bootstrap %s: %s",
	ExitCode: exit.ExitGeneric,
}

func (si systemupdate) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// create all the other directories
	logging.LogMessage(context.IOStream, "Ensuring distillery installation directories exist")
	for _, d := range []string{
		dis.Config.DeployRoot,
		dis.Instances().Path(),
		dis.SnapshotsStagingPath(),
		dis.SnapshotsArchivePath(),
	} {
		context.Println(d)
		if err := dis.Core.Environment.MkdirAll(d, environment.DefaultDirPerm); err != nil {
			return errBoostrapFailedToCreateDirectory.WithMessageF(d, err)
		}
	}

	if !si.SkipCoreUpdates {
		// install system updates
		logging.LogMessage(context.IOStream, "Updating Operating System Packages")
		if err := si.mustExec(context, "", "apt-get", "update"); err != nil {
			return err
		}
		if err := si.mustExec(context, "", "apt-get", "upgrade", "-y"); err != nil {
			return err
		}

		// install docker
		logging.LogMessage(context.IOStream, "Installing / Updating Docker")
		if err := si.mustExec(context, "", "apt-get", "install", "curl"); err != nil {
			return err
		}
		// TODO: Download directly
		if err := si.mustExec(context, "", "/bin/sh", "-c", "curl -fsSL https://get.docker.com -o - | /bin/sh"); err != nil {
			return err
		}
	}

	// create the docker network
	// TODO: Use docker API for this
	logging.LogMessage(context.IOStream, "Updating Docker Configuration")
	si.mustExec(context, "", "docker", "network", "create", "distillery")

	// install and update the various stacks!
	ctx := component.InstallationContext{
		"graphdb.zip": si.Positionals.GraphdbZip,
	}

	if err := logging.LogOperation(func() error {
		for _, component := range dis.Installables() {
			name := component.Name()
			stack := component.Stack(dis.Core.Environment)
			ctx := component.Context(ctx)

			if err := logging.LogOperation(func() error {
				return stack.Install(context.IOStream, ctx)
			}, context.IOStream, "Installing Docker Stack %q", name); err != nil {
				return err
			}

			if err := logging.LogOperation(func() error {
				return stack.Update(context.IOStream, true)
			}, context.IOStream, "Updating Docker Stack: %q", name); err != nil {
				return err
			}
		}
		return nil
	}, context.IOStream, "Performing Stack Updates"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		for _, component := range dis.Updateable() {
			name := component.Name()
			if err := logging.LogOperation(func() error {
				return component.Update(context.IOStream)
			}, context.IOStream, "Updating Component: %s", name); err != nil {
				return errBootstrapComponent.WithMessageF(name, err)
			}
		}
		return nil
	}, context.IOStream, "Performing Component Updates"); err != nil {
		return err
	}
	// TODO: Register cronjob in /etc/cron.d!

	logging.LogMessage(context.IOStream, "System has been updated")
	return nil
}

var errMustExecFailed = exit.Error{
	Message: "process exited with code %d",
}

// mustExec indicates that the given executable process must complete successfully.
// If it does not, returns errMustExecFailed
func (si systemupdate) mustExec(context wisski_distillery.Context, workdir string, exe string, argv ...string) error {
	dis := context.Environment
	if workdir == "" {
		workdir = context.Environment.Config.DeployRoot
	}
	code := dis.Core.Environment.Exec(context.IOStream, workdir, exe, argv...)
	if code != 0 {
		err := errMustExecFailed.WithMessageF(code)
		err.ExitCode = exit.ExitCode(code)
		return err
	}
	return nil
}
