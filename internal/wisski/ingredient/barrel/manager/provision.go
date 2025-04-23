//spellchecker:words manager
package manager

//spellchecker:words context time github wisski distillery internal component models ingredient barrel composer extras logging pkglib contextx stream
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/contextx"
	"github.com/tkw1536/pkglib/stream"
)

// Provision applies defaults to flags, to ensure some values are set.
func (manager *Manager) Provision(ctx context.Context, progress io.Writer, system models.System, flags Profile) error {
	// Force building and applying the system!
	if err := manager.dependencies.SystemManager.ApplyInitial(ctx, progress, system); err != nil {
		return err
	}

	// Create the composer directory!
	logging.LogMessage(progress, "Creating required directories")
	{
		code, err := manager.dependencies.Barrel.Stack().Run(ctx, stream.FromNil(), component.RunFlags{Detach: true, AutoRemove: true}, "barrel", "sudo", "-u", "www-data", "mkdir", "-p", barrel.ComposerDirectory)
		if code != 0 {
			err = barrel.ExitError(code)
		}
		if err != nil {
			return err
		}
	}

	// start the container, and have it do nothing!
	code, err := manager.dependencies.Barrel.Stack().Run(ctx, stream.FromNil(), component.RunFlags{Detach: true, AutoRemove: true}, "barrel", "tail", "-f", "/dev/null")
	if code != 0 {
		err = barrel.ExitError(code)
	}
	if err != nil {
		return err
	}

	// when we are done, shut it down!
	defer func() {
		anyways, cancel := contextx.Anyways(ctx, time.Minute)
		defer cancel()

		// stop the container (even if the context was cancelled)
		manager.dependencies.Barrel.Stack().DownAll(anyways, progress)
	}()

	return manager.bootstrap(ctx, progress, flags)
}

// TODO: Move this to the flags.
var drushVariants = []string{
	"drush/drush", "drush/drush:^12", "drush/drush:^11",
}

// bootstrap applies the initial flags induced by flags.
// Applies defaults to the flags.
func (provision *Manager) bootstrap(ctx context.Context, progress io.Writer, flags Profile) error {
	// TODO: Check if we can remove the easyrdf patch!
	flags.ApplyDefaults()

	logging.LogMessage(progress, "Creating Composer Project")
	{
		drupal := "drupal/recommended-project"
		if flags.Drupal != "" {
			drupal += ":" + flags.Drupal
		}
		err := provision.dependencies.Composer.Exec(ctx, progress, "create-project", drupal, ".")
		if err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Configuring Composer")
	{
		// needed for composer > 2.2
		err := provision.dependencies.Composer.Exec(ctx, progress, "config", "allow-plugins", "true")
		if err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Installing drush")
	{
		for _, v := range drushVariants {
			err := provision.dependencies.Composer.TryInstall(ctx, progress, v)
			if errors.Is(err, composer.ErrNotInstalled) {
				continue
			}
			if err != nil {
				return err
			}
			break
		}
	}

	liquid := ingredient.GetLiquid(provision)

	var sqlDBURL = "mysql://" + liquid.SqlUsername + ":" + liquid.SqlPassword + "@sql/" + liquid.SqlDatabase

	// Use 'drush' to run the site-installation.
	// Here we need to use the username, password and database creds we made above.
	logging.LogMessage(progress, "Running Drupal installation scripts")
	{
		if err := provision.dependencies.Drush.Exec(
			ctx, progress,
			"site-install",
			"standard", "--yes", "--site-name="+liquid.Domain(),
			"--account-name="+liquid.DrupalUsername, "--account-pass="+liquid.DrupalPassword,
			"--db-url="+sqlDBURL,
		); err != nil {
			return err
		}

		if err := provision.dependencies.Composer.FixPermission(ctx, progress); err != nil {
			return err
		}
	}

	// Rebuild the settings file
	logging.LogMessage(progress, "Rebuilding Settings")
	{
		if err := provision.dependencies.SystemManager.BuildSettings(ctx, progress); err != nil {
			return err
		}
	}

	// Create directory for ontologies
	logging.LogMessage(progress, fmt.Sprintf("Creating %q", barrel.OntologyDirectory))
	{
		if err := provision.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), "mkdir", "-p", barrel.OntologyDirectory); err != nil {
			return err
		}
	}

	{
		// make a set of flags to apply to the given instance
		flags := flags
		flags.Drupal = "" // Do not upgrade Drupal
		flags.WissKI = "" // Do not upgrade WissKI

		// apply the rest of the flags
		if err := provision.Apply(ctx, progress, flags); err != nil {
			return err
		}
	}

	// install WissKI
	if err := provision.applyWissKI(ctx, progress, flags.WissKI); err != nil {
		return err
	}

	// create the default adapter
	logging.LogMessage(progress, "Creating default adapter")
	{
		if _, err := provision.dependencies.Adapters.SetAdapter(ctx, nil, provision.dependencies.Adapters.DefaultAdapter()); err != nil {
			return fmt.Errorf("failed to create default adapter: %w", err)
		}
	}

	logging.LogMessage(progress, "Running initial cron")
	{
		if err := provision.dependencies.Drush.Exec(ctx, progress, "core-cron"); err != nil {
			return fmt.Errorf("failed to run initial cron: %w", err)
		}
	}

	logging.LogMessage(progress, "Provisioning is now complete")
	{
		fmt.Fprintf(progress, "URL:                  %s\n", liquid.URL())
		fmt.Fprintf(progress, "Username:             %s\n", liquid.DrupalUsername)
		fmt.Fprintf(progress, "Password:             %s\n", liquid.DrupalPassword)
	}

	return nil
}
