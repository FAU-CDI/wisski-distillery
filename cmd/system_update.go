package cmd

import (
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/distillery"
	"github.com/FAU-CDI/wisski-distillery/env"
	"github.com/FAU-CDI/wisski-distillery/internal/execx"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
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
		Requirements: env.Requirements{
			NeedsConfig: true,
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
	_, err := os.Stat(s.Positionals.GraphdbZip)
	if os.IsNotExist(err) {
		return errNoGraphDBZip.WithMessageF(s.Positionals.GraphdbZip)
	}
	if err != nil {
		return err
	}
	return nil
}

var errFailedToCreateDirectory = exit.Error{
	Message:  "failed to create directory %s: %s",
	ExitCode: exit.ExitGeneric,
}

var errFailedRuntime = exit.Error{
	Message:  "failed to update runtime: %s",
	ExitCode: exit.ExitGeneric,
}

func (si systemupdate) Run(context wisski_distillery.Context) error {
	dis := context.Environment

	// create all the other directories
	logging.LogMessage(context.IOStream, "Ensuring distillery installation directories exist")
	for _, d := range []string{
		dis.Config.DeployRoot,
		dis.InstancesDir(),
		dis.InprogressBackupPath(),
		dis.FinalBackupPath(),
	} {
		context.Println(d)
		if err := os.MkdirAll(d, os.ModeDir); err != nil {
			return errFailedToCreateDirectory.WithMessageF(d, err)
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
		if err := si.mustExec(context, "", "/bin/sh", "-c", "curl -fsSL https://get.docker.com -o - | /bin/sh"); err != nil {
			return err
		}
	}

	// create the docker network
	// TODO: Use docker API for this
	logging.LogMessage(context.IOStream, "Updating Docker Configuration")
	si.mustExec(context, "", "docker", "network", "create", "distillery")

	// install and update the various stacks!
	ctx := stack.InstallationContext{
		"graphdb.zip": si.Positionals.GraphdbZip,
	}

	if err := logging.LogOperation(func() error {
		for _, stack := range dis.Stacks() {
			if err := logging.LogOperation(func() error {
				return stack.Install(context.IOStream, ctx)
			}, context.IOStream, "Installing docker stack %q", stack.Dir); err != nil {
				return err
			}

			if err := logging.LogOperation(func() error {
				return stack.Update(context.IOStream, true)
			}, context.IOStream, "Updating docker stack %q", stack.Dir); err != nil {
				return err
			}
		}
		return nil
	}, context.IOStream, "Updating Components"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		if err := distillery.InstallResource(dis.RuntimeDir(), filepath.Join("resources", "runtime"), func(dst, src string) {
			context.Printf("[copy]  %s\n", dst)
		}); err != nil {
			return errFailedRuntime.WithMessageF(err)
		}
		return nil
	}, context.IOStream, "Unpacking Runtime Components"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		if err := dis.SQLBootstrap(context.IOStream); err != nil {
			return err
		}
		return nil
	}, context.IOStream, "Bootstraping SQL database"); err != nil {
		return err
	}

	if err := logging.LogOperation(func() error {
		if err := dis.TriplestoreBootstrap(context.IOStream); err != nil {
			return err
		}
		return nil
	}, context.IOStream, "Bootstraping Triplestore"); err != nil {
		return err
	}

	logging.LogMessage(context.IOStream, "System has been updated")
	return nil
}

var errMustExecFailed = exit.Error{
	Message: "process exited with code %d",
}

// mustExec indicates that the given executable process must complete successfully.
// If it does not, returns errMustExecFailed
func (si systemupdate) mustExec(context wisski_distillery.Context, workdir string, exe string, argv ...string) error {
	if workdir == "" {
		workdir = context.Environment.Config.DeployRoot
	}
	code := execx.Exec(context.IOStream, workdir, exe, argv...)
	if code != 0 {
		err := errMustExecFailed.WithMessageF(code)
		err.ExitCode = exit.ExitCode(code)
		return err
	}
	return nil
}
