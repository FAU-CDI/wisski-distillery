package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/core"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Bootstrap is the 'bootstrap' command
var Bootstrap wisski_distillery.Command = bootstrap{}

type bootstrap struct {
	Directory string `short:"r" long:"root-directory" description:"path to the root deployment directory" default:"/var/www/deploy"`
	Hostname  string `short:"h" long:"hostname" description:"default hostname of the distillery (default: system hostname)"`
}

func (bootstrap) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: core.Requirements{
			NeedsDistillery: false,
		},
		Command:     "bootstrap",
		Description: "Bootstraps the installation of a Distillery System",
	}
}

var errBootstrapDifferent = exit.Error{
	Message:  "refusing to bootstrap: base directory is already set to %s.",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapFailedToCreateDirectory = exit.Error{
	Message:  "failed to create directory %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapFailedToSaveDirectory = exit.Error{
	Message:  "failed to register base directory: %s",
	ExitCode: exit.ExitGeneric,
}

var errBoostrapFailedToCopyExe = exit.Error{
	Message:  "failed to copy wdcli executable: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapWriteConfig = exit.Error{
	Message:  "failed to write configuration file: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapOpenConfig = exit.Error{
	Message:  "failed to open configuration file: %s",
	ExitCode: exit.ExitGeneric,
}

var errBootstrapCreateFile = exit.Error{
	Message:  "failed to touch configuration file: %s",
	ExitCode: exit.ExitGeneric,
}

func (bs bootstrap) Run(context wisski_distillery.Context) error {
	root := bs.Directory

	// check that we didn't get a different base directory
	{
		got, err := core.ReadBaseDirectory()
		if err == nil && got != "" && got != root {
			return errBootstrapDifferent.WithMessageF(got)
		}
	}

	{
		logging.LogMessage(context.IOStream, "Creating root deployment directory")
		if err := os.MkdirAll(root, fs.ModeDir); err != nil {
			return errBootstrapFailedToCreateDirectory.WithMessageF(root)
		}
		if err := core.WriteBaseDirectory(root); err != nil {
			return errBootstrapFailedToSaveDirectory.WithMessageF(root)
		}
		context.Println(root)
	}

	// TODO: Should we read an existing configuration file?
	wdcliPath := filepath.Join(root, core.Executable)
	envPath := filepath.Join(root, core.ConfigFile)

	// setup a new template for the configuration file!
	var tpl config.Template
	tpl.DeployRoot = bs.Directory
	tpl.DefaultDomain = bs.Hostname

	// and use thge defaults
	if err := tpl.SetDefaults(); err != nil {
		return errBootstrapWriteConfig.WithMessageF(err)
	}

	{
		logging.LogMessage(context.IOStream, "Copying over wdcli executable")
		exe, err := os.Executable()
		if err != nil {
			return errBoostrapFailedToCopyExe.WithMessageF(err)
		}

		err = fsx.CopyFile(wdcliPath, exe)
		if err != nil && err != fsx.ErrCopySameFile {
			return errBoostrapFailedToCopyExe.WithMessageF(err)
		}
		context.Println(wdcliPath)
	}

	{
		if !fsx.IsFile(envPath) {
			if err := logging.LogOperation(func() error {
				env, err := os.Create(envPath)
				if err != nil {
					return err
				}
				defer env.Close()

				return tpl.MarshalTo(env)
			}, context.IOStream, "Installing configuration file"); err != nil {
				return errBootstrapWriteConfig.WithMessageF(err)
			}

			if err := logging.LogOperation(func() error {

				context.Println(tpl.SelfOverridesFile)
				if err := os.WriteFile(
					tpl.SelfOverridesFile,
					core.DefaultOverridesJSON,
					fs.ModePerm,
				); err != nil {
					return errBootstrapCreateFile.WithMessageF(err)
				}

				context.Println(tpl.AuthorizedKeys)
				if err := os.WriteFile(
					tpl.AuthorizedKeys,
					core.DefaultAuthorizedKeys,
					fs.ModePerm,
				); err != nil {
					return errBootstrapCreateFile.WithMessageF(err)
				}

				return nil
			}, context.IOStream, "Creating additional config files"); err != nil {
				return err
			}
		}

	}

	// re-read the configuration and print it!
	logging.LogMessage(context.IOStream, "Configuration is now complete")
	f, err := os.Open(envPath)
	if err != nil {
		return errBootstrapOpenConfig.WithMessageF(err)
	}
	defer f.Close()

	var cfg config.Config
	if err := cfg.Unmarshal(f); err != nil {
		return errBootstrapOpenConfig.WithMessageF(err)
	}
	context.Println(cfg)

	// Tell the user how to proceed
	logging.LogMessage(context.IOStream, "Bootstrap is complete")
	context.Printf("Adjust the configuration file at %s\n", envPath)
	context.Printf("Then grab a GraphDB zipped source file and run:\n")
	context.Printf("%s system_update /path/to/graphdb.zip\n", wdcliPath)

	return nil
}
