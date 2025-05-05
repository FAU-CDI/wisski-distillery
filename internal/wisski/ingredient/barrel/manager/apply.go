//spellchecker:words manager
package manager

//spellchecker:words context github wisski distillery internal ingredient barrel composer logging pkglib stream
import (
	"context"
	"fmt"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/errorsx"
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
	var message string
	if enable {
		message = "Installing and enabling modules"
	} else {
		message = "Installing modules"
	}

	if err := logging.LogOperation(func() error {
		for _, spec := range modules {
			if _, err := logging.LogMessage(progress, fmt.Sprintf("Installing %q", spec)); err != nil {
				return fmt.Errorf("failed to log message: %w", err)
			}
			err := manager.dependencies.Composer.Install(ctx, progress, spec)
			if err != nil {
				return fmt.Errorf("failed to install module %q: %w", spec, err)
			}

			if enable {
				name := composer.ModuleName(spec)
				if _, err := logging.LogMessage(progress, fmt.Sprintf("Enabling %q (from spec %q)", name, spec)); err != nil {
					return fmt.Errorf("failed to log message: %w", err)
				}
				err := manager.dependencies.Drush.Enable(ctx, progress, name)
				if err != nil {
					return fmt.Errorf("failed to enable module %q: %w", name, err)
				}
			}
		}
		return nil
	}, progress, "%s", message); err != nil {
		return fmt.Errorf("failed to install modules: %w", err)
	}
	return nil
}

// applyDrupal applies a specific drupal version.
// Assumes that drupal != "".
func (manager *Manager) applyDrupal(ctx context.Context, progress io.Writer, drupal string) (e error) {
	if err := logging.LogOperation(func() error {
		if _, err := logging.LogMessage(progress, "Clearing up permissions for update"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			for _, script := range [][]string{
				{"chmod", "777", "web/sites/default"},
				{"chmod", "666", "web/sites/default/*settings.php"},
				{"chmod", "666", "web/sites/default/*services.php"},
			} {
				err := manager.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), script...)
				if err != nil {
					return fmt.Errorf("failed to change permissions before update: %w", err)
				}
			}
		}

		defer func() {
			if _, err := logging.LogMessage(progress, "Resetting permissions"); err != nil {
				err = fmt.Errorf("failed to log message: %w", err)
				e = errorsx.Combine(e, err)
				return
			}

			{
				for _, script := range [][]string{
					{"chmod", "755", "web/sites/default"},
					{"chmod", "644", "web/sites/default/*settings.php"},
					{"chmod", "644", "web/sites/default/*services.php"},
				} {
					if err := manager.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), script...); err != nil {
						err = fmt.Errorf("failed to reset permissions after update: %w", err)
						e = errorsx.Combine(e, err)
					}
				}
			}
		}()

		// write out a specific Drupal version
		if _, err := logging.LogMessage(progress, "Performing Drupal update"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			args := []string{
				"drupal/internal/core-recommended:", "drupal/internal/core-composer-scaffold:", "drupal/internal/core-project-message:",
			}
			for i, cm := range args {
				args[i] = cm + drupal
			}
			args = append(args, "--update-with-dependencies", "--no-update")

			if err := manager.dependencies.Composer.Install(ctx, progress, args...); err != nil {
				return fmt.Errorf("failed to install drupal core: %w", err)
			}
		}

		if _, err := logging.LogMessage(progress, "Running composer update"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			if err := manager.dependencies.Composer.Exec(ctx, progress, "update"); err != nil {
				return fmt.Errorf("failed to update: %w", err)
			}
		}

		if _, err := logging.LogMessage(progress, "Performing database updates (if any)"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			if err := manager.dependencies.Drush.Exec(ctx, progress, "updatedb", "--yes"); err != nil {
				return fmt.Errorf("failed to update database: %w", err)
			}
		}

		return nil
	}, progress, "%s", "Updating to Drupal %q", drupal); err != nil {
		return fmt.Errorf("failed to update drupal: %w", err)
	}
	return nil
}

// applyWissKI applies the WissKI version.
func (manager *Manager) applyWissKI(ctx context.Context, progress io.Writer, wisski string) error {
	if err := logging.LogOperation(func() error {
		if _, err := logging.LogMessage(progress, "Installing WissKI Module"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			spec := "drupal/wisski"
			if wisski != "" {
				spec += ":" + wisski
			}

			err := manager.dependencies.Composer.Install(ctx, progress, spec)
			if err != nil {
				return fmt.Errorf("failed to install WissKI: %w", err)
			}
		}

		// install dependencies in the WissKI directory
		if _, err := logging.LogMessage(progress, "Installing WissKI Dependencies"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			if err := manager.dependencies.Composer.ExecWissKI(ctx, progress, "install"); err != nil {
				return fmt.Errorf("failed to install wisski dependencies: %w", err)
			}
		}

		if _, err := logging.LogMessage(progress, "Enable Wisski modules"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			if err := manager.dependencies.Drush.Enable(ctx, progress,
				"wisski_core", "wisski_linkblock", "wisski_pathbuilder", "wisski_adapter_sparql11_pb", "wisski_salz",
			); err != nil {
				return fmt.Errorf("failed to enable wisski modules: %w", err)
			}

			if err := manager.dependencies.Composer.FixPermission(ctx, progress); err != nil {
				return fmt.Errorf("failed to fix permissions: %w", err)
			}
		}

		if _, err := logging.LogMessage(progress, "Performing database updates (if any)"); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		{
			if err := manager.dependencies.Drush.Exec(ctx, progress, "updatedb", "--yes"); err != nil {
				return fmt.Errorf("failed to update database with drush: %w", err)
			}
		}

		return nil
	}, progress, "Installing WissKI version %q", wisski); err != nil {
		return fmt.Errorf("failed to install drupal version: %w", err)
	}
	return nil
}
