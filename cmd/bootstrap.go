package cmd

//spellchecker:words errors path filepath github wisski distillery internal bootstrap config logging goprogram exit pkglib errorsx umaskfree
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/bootstrap"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/config"

	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/errorsx"
	"github.com/tkw1536/pkglib/fsx"
	"github.com/tkw1536/pkglib/fsx/umaskfree"
)

// Bootstrap is the 'bootstrap' command.
var Bootstrap wisski_distillery.Command = cBootstrap{}

type cBootstrap struct {
	Directory string `default:"/var/www/deploy"                                                   description:"path to the root deployment directory" long:"root-directory" short:"r"`
	Hostname  string `description:"default hostname of the distillery (default: system hostname)" long:"hostname"                                     short:"h"`
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

var (
	errBootstrapDifferent               = exit.NewErrorWithCode("refusing to bootstrap: base directory is already set to", exit.ExitGeneric)
	errBootstrapFailedToCreateDirectory = exit.NewErrorWithCode("failed to create directory", exit.ExitGeneric)
	errBootstrapFailedToSaveDirectory   = exit.NewErrorWithCode("failed to register base directory", exit.ExitGeneric)
	errBoostrapFailedToCopyExe          = exit.NewErrorWithCode("failed to copy wdcli executable", exit.ExitGeneric)
	errBootstrapWriteConfig             = exit.NewErrorWithCode("failed to write configuration file", exit.ExitGeneric)
	errBootstrapOpenConfig              = exit.NewErrorWithCode("failed to open configuration file", exit.ExitGeneric)
	errBootstrapCreateFile              = exit.NewErrorWithCode("failed to touch configuration file", exit.ExitGeneric)
)

func (bs cBootstrap) Run(context wisski_distillery.Context) (e error) {
	root := bs.Directory

	// check that we didn't get a different base directory
	{
		got, err := cli.ReadBaseDirectory()
		if err == nil && got != "" && got != root {
			return fmt.Errorf("%w %q", errBootstrapDifferent, got)
		}
	}

	{
		if _, err := logging.LogMessage(context.Stderr, "Creating root deployment directory"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if err := umaskfree.MkdirAll(root, umaskfree.DefaultDirPerm); err != nil {
			return fmt.Errorf("%q: %w: %w", root, errBootstrapFailedToCreateDirectory, err)
		}
		if err := cli.WriteBaseDirectory(root); err != nil {
			return fmt.Errorf("%q: %w: %w", root, errBootstrapFailedToSaveDirectory, err)
		}
		_, _ = context.Println(root)
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
		return fmt.Errorf("%w: %w", errBootstrapWriteConfig, err)
	}

	{
		if _, err := logging.LogMessage(context.Stderr, "Copying over wdcli executable"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("%w: %w", errBoostrapFailedToCopyExe, err)
		}

		err = umaskfree.CopyFile(context.Context, wdcliPath, exe)
		if err != nil && !errors.Is(err, umaskfree.ErrCopySameFile) {
			return fmt.Errorf("%w: %w", errBoostrapFailedToCopyExe, err)
		}
		_, _ = context.Println(wdcliPath)
	}

	{
		isFile, err := fsx.IsRegular(cfgPath, false)
		if err != nil {
			return fmt.Errorf("%w: %w", errBootstrapWriteConfig, err)
		}
		if !isFile {
			// generate the configuration from the template
			cfg := tpl.Generate()

			// write out all the extra config files
			if err := logging.LogOperation(func() error {
				if _, err := context.Println(cfg.Paths.OverridesJSON); err != nil {
					return fmt.Errorf("failed to write text: %w", err)
				}
				if err := umaskfree.WriteFile(
					cfg.Paths.OverridesJSON,
					bootstrap.DefaultOverridesJSON,
					fs.ModePerm,
				); err != nil {
					return fmt.Errorf("failed to write overrides file: %w", err)
				}

				_, _ = context.Println(cfg.Paths.ResolverBlocks)
				if err := umaskfree.WriteFile(
					cfg.Paths.ResolverBlocks,
					bootstrap.DefaultResolverBlockedTXT,
					fs.ModePerm,
				); err != nil {
					return fmt.Errorf("failed to write resolver blocks file: %w", err)
				}

				return nil
			}, context.Stderr, "Creating custom config files"); err != nil {
				return fmt.Errorf("%w: %w", errBootstrapCreateFile, err)
			}

			// Validate configuration file!
			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("failed to validate configuration: %w", err)
			}

			// and marshal it out!
			if err := logging.LogOperation(func() (e error) {
				configYML, err := umaskfree.Create(cfgPath, umaskfree.DefaultFilePerm)
				if err != nil {
					return fmt.Errorf("failed to create configuration path: %w", err)
				}
				defer errorsx.Close(configYML, &e, "configuration file")

				bytes, err := config.Marshal(&cfg, nil)
				if err != nil {
					return fmt.Errorf("failed to marshal configuration file: %w", err)
				}

				{
					_, err := configYML.Write(bytes)
					return fmt.Errorf("failed to write config yml: %w", err)
				}
			}, context.Stderr, "Installing primary configuration file"); err != nil {
				return fmt.Errorf("%w: %w", err, errBootstrapWriteConfig)
			}
		}
	}

	// re-read the configuration and print it!
	if _, err := logging.LogMessage(context.Stderr, "Configuration is now complete"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	f, err := os.Open(cfgPath) // #nosec G304 -- intended
	if err != nil {
		return fmt.Errorf("%w: %w", errBootstrapOpenConfig, err)
	}
	defer errorsx.Close(f, &e, "configuration file")

	var cfg config.Config
	if err := cfg.Unmarshal(f); err != nil {
		return fmt.Errorf("%w: %w", errBootstrapOpenConfig, err)
	}
	_, _ = context.Println(cfg)

	// Tell the user how to proceed
	if _, err := logging.LogMessage(context.Stderr, "Bootstrap is complete"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	if _, err := context.Printf("Adjust the configuration file at %s\n", cfgPath); err != nil {
		return fmt.Errorf("failed to report progress: %w", err)
	}
	if _, err := context.Printf("Then make sure 'docker compose' is installed.\n"); err != nil {
		return fmt.Errorf("failed to report progress: %w", err)
	}
	if _, err := context.Printf("Finally grab a GraphDB zipped source file and run:\n"); err != nil {
		return fmt.Errorf("failed to report progress: %w", err)
	}
	if _, err := context.Printf("%s system_update /path/to/graphdb.zip\n", wdcliPath); err != nil {
		return fmt.Errorf("failed to report progress: %w", err)
	}

	return nil
}
