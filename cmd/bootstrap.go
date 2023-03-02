package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/config"

	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Bootstrap is the 'bootstrap' command
var Bootstrap wisski_distillery.Command = cBootstrap{}

type cBootstrap struct {
	Directory string `short:"r" long:"root-directory" description:"path to the root deployment directory" default:"/var/www/deploy"`
	Hostname  string `short:"h" long:"hostname" description:"default hostname of the distillery (default: system hostname)"`
}

func (cBootstrap) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: false,
		},
		Command:     "bootstrap",
		Description: "bootstraps the installation of a distillery system",
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

func (bs cBootstrap) Run(context wisski_distillery.Context) error {
	root := bs.Directory

	// check that we didn't get a different base directory
	{
		got, err := cli.ReadBaseDirectory()
		if err == nil && got != "" && got != root {
			return errBootstrapDifferent.WithMessageF(got)
		}
	}

	{
		logging.LogMessage(context.Stderr, context.Context, "Creating root deployment directory")
		if err := fsx.MkdirAll(root, fsx.DefaultDirPerm); err != nil {
			return errBootstrapFailedToCreateDirectory.WithMessageF(root)
		}
		if err := cli.WriteBaseDirectory(root); err != nil {
			return errBootstrapFailedToSaveDirectory.WithMessageF(root)
		}
		context.Println(root)
	}

	// TODO: Should we read an existing configuration file?
	wdcliPath := filepath.Join(root, bootstrap.Executable)
	cfgPath := filepath.Join(root, bootstrap.ConfigFile)

	// setup a new template for the configuration file!
	var tpl config.Template
	tpl.RootPath = bs.Directory
	tpl.DefaultDomain = bs.Hostname

	// and use thge defaults
	if err := tpl.SetDefaults(); err != nil {
		return errBootstrapWriteConfig.WithMessageF(err)
	}

	{
		logging.LogMessage(context.Stderr, context.Context, "Copying over wdcli executable")
		exe, err := os.Executable()
		if err != nil {
			return errBoostrapFailedToCopyExe.WithMessageF(err)
		}

		err = fsx.CopyFile(context.Context, wdcliPath, exe)
		if err != nil && err != fsx.ErrCopySameFile {
			return errBoostrapFailedToCopyExe.WithMessageF(err)
		}
		context.Println(wdcliPath)
	}

	{
		if !fsx.IsFile(cfgPath) {
			// generate the configuration from the template
			cfg := tpl.Generate()

			// write out all the extra config files
			if err := logging.LogOperation(func() error {
				context.Println(cfg.Paths.OverridesJSON)
				if err := fsx.WriteFile(
					cfg.Paths.OverridesJSON,
					bootstrap.DefaultOverridesJSON,
					fs.ModePerm,
				); err != nil {
					return err
				}

				context.Println(cfg.Paths.ResolverBlocks)
				if err := fsx.WriteFile(
					cfg.Paths.ResolverBlocks,
					bootstrap.DefaultResolverBlockedTXT,
					fs.ModePerm,
				); err != nil {
					return err
				}

				return nil
			}, context.Stderr, context.Context, "Creating custom config files"); err != nil {
				return errBootstrapCreateFile.WithMessageF(err)
			}

			// Validate configuration file!
			if err := cfg.Validate(); err != nil {
				return err
			}

			// and marshal it out!
			if err := logging.LogOperation(func() error {
				configYML, err := fsx.Create(cfgPath, fsx.DefaultFilePerm)
				if err != nil {
					return err
				}
				defer configYML.Close()

				bytes, err := config.Marshal(&cfg, nil)
				if err != nil {
					return err
				}

				{
					_, err := configYML.Write(bytes)
					return err
				}
			}, context.Stderr, context.Context, "Installing primary configuration file"); err != nil {
				return errBootstrapWriteConfig.WithMessageF(err)
			}
		}

	}

	// re-read the configuration and print it!
	logging.LogMessage(context.Stderr, context.Context, "Configuration is now complete")
	f, err := os.Open(cfgPath)
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
	logging.LogMessage(context.Stderr, context.Context, "Bootstrap is complete")
	context.Printf("Adjust the configuration file at %s\n", cfgPath)
	context.Printf("Then make sure 'docker compose' is installed.\n")
	context.Printf("Finally grab a GraphDB zipped source file and run:\n")
	context.Printf("%s system_update /path/to/graphdb.zip\n", wdcliPath)

	return nil
}
