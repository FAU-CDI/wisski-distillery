//spellchecker:words manager
package manager

//spellchecker:words context github wisski distillery internal ingredient barrel composer logging pkglib stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/stream"
)

// The instance must be running.
func (manager *Manager) Apply(ctx context.Context, progress io.Writer, flags Profile) error {
	// Update drupal
	if flags.Drupal != "" {
		err := manager.applyDrupal(ctx, progress, flags.Drupal)
		if err != nil {
			return err
		}
	}

	// Update WissKI
	if flags.WissKI != "" {
		err := manager.applyWissKI(ctx, progress, flags.WissKI)
		if err != nil {
			return err
		}
	}

	// install custom modules
	if len(flags.InstallModules) > 0 {
		err := manager.installModules(ctx, progress, flags.InstallModules, false)
		if err != nil {
			return err
		}
	}

	// install + enable custom modules
	if len(flags.EnableModules) > 0 {
		err := manager.installModules(ctx, progress, flags.EnableModules, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (manager *Manager) installModules(ctx context.Context, progress io.Writer, modules []string, enable bool) error {
	message := ""
	if enable {
		message = "Installing and enabling modules"
	} else {
		message = "Installing modules"
	}

	// enable the module
	return logging.LogOperation(func() error {
		for _, spec := range modules {
			logging.LogMessage(progress, fmt.Sprintf("Installing %q", spec))
			err := manager.dependencies.Composer.Install(ctx, progress, spec)
			if err != nil {
				return err
			}

			if enable {
				name := composer.ModuleName(spec)
				logging.LogMessage(progress, fmt.Sprintf("Enabling %q (from spec %q)", name, spec))
				err := manager.dependencies.Drush.Enable(ctx, progress, name)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}, progress, "%s", message)
}

// applyDrupal applies a specific drupal version.
// Assumes that drupal != "".
func (manager *Manager) applyDrupal(ctx context.Context, progress io.Writer, drupal string) error {
	return logging.LogOperation(func() error {
		logging.LogMessage(progress, "Clearing up permissions for update")
		{
			for _, script := range [][]string{
				{"chmod", "777", "web/sites/default"},
				{"chmod", "666", "web/sites/default/*settings.php"},
				{"chmod", "666", "web/sites/default/*services.php"},
			} {
				err := manager.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), script...)
				if err != nil {
					return err
				}
			}
		}

		defer func() {
			logging.LogMessage(progress, "Resetting permissions")
			{
				for _, script := range [][]string{
					{"chmod", "755", "web/sites/default"},
					{"chmod", "644", "web/sites/default/*settings.php"},
					{"chmod", "644", "web/sites/default/*services.php"},
				} {
					manager.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), script...)
				}
			}
		}()

		// write out a specific Drupal version
		logging.LogMessage(progress, "Performing Drupal update")
		{
			args := []string{
				"drupal/internal/core-recommended:", "drupal/internal/core-composer-scaffold:", "drupal/internal/core-project-message:",
			}
			for i, cm := range args {
				args[i] = cm + drupal
			}
			args = append(args, "--update-with-dependencies", "--no-update")

			if err := manager.dependencies.Composer.Install(ctx, progress, args...); err != nil {
				return err
			}
		}

		logging.LogMessage(progress, "Running composer update")
		{
			if err := manager.dependencies.Composer.Exec(ctx, progress, "update"); err != nil {
				return err
			}
		}

		logging.LogMessage(progress, "Performing database updates (if any)")
		{
			if err := manager.dependencies.Drush.Exec(ctx, progress, "updatedb", "--yes"); err != nil {
				return err
			}
		}

		return nil
	}, progress, "%s", "Updating to Drupal %q", drupal)
}

// applyWissKI applies the WissKI version.
func (manager *Manager) applyWissKI(ctx context.Context, progress io.Writer, wisski string) error {
	return logging.LogOperation(func() error {
		logging.LogMessage(progress, "Installing WissKI Module")
		{
			spec := "drupal/wisski"
			if wisski != "" {
				spec += ":" + wisski
			}

			err := manager.dependencies.Composer.Install(ctx, progress, spec)
			if err != nil {
				return err
			}
		}

		// install dependencies in the WissKI directory
		logging.LogMessage(progress, "Installing WissKI Dependencies")
		{
			if err := manager.dependencies.Composer.ExecWissKI(ctx, progress, "install"); err != nil {
				return err
			}
		}

		logging.LogMessage(progress, "Enable Wisski modules")
		{
			if err := manager.dependencies.Drush.Enable(ctx, progress,
				"wisski_core", "wisski_linkblock", "wisski_pathbuilder", "wisski_adapter_sparql11_pb", "wisski_salz",
			); err != nil {
				return err
			}

			if err := manager.dependencies.Composer.FixPermission(ctx, progress); err != nil {
				return err
			}
		}

		logging.LogMessage(progress, "Performing database updates (if any)")
		{
			if err := manager.dependencies.Drush.Exec(ctx, progress, "updatedb", "--yes"); err != nil {
				return err
			}
		}

		return nil
	}, progress, "Installing WissKI version %q", wisski)
}
