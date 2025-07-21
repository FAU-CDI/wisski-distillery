//spellchecker:words manager
package manager

//spellchecker:words context errors time github wisski distillery internal models ingredient barrel composer dockerx logging pkglib contextx errorsx stream
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"go.tkw01536.de/pkglib/contextx"
	"go.tkw01536.de/pkglib/errorsx"
	"go.tkw01536.de/pkglib/stream"
)

// Provision applies defaults to flags, to ensure some values are set.
func (manager *Manager) Provision(ctx context.Context, progress io.Writer, system models.System, flags Profile) (e error) {
	// Force building and applying the system!
	if err := manager.dependencies.SystemManager.ApplyInitial(ctx, progress, system); err != nil {
		return fmt.Errorf("failed to apply initial configuration: %w", err)
	}

	stack, err := manager.dependencies.Barrel.OpenStack()
	if err != nil {
		return fmt.Errorf("failed to open stack: %w", err)
	}
	defer errorsx.Close(stack, &e, "stack")

	// Create the composer directory!
	if _, err := logging.LogMessage(progress, "Creating required directories"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		code, err := stack.Run(ctx, stream.FromNil(), dockerx.RunFlags{Detach: true, AutoRemove: true}, "barrel", "sudo", "-u", "www-data", "mkdir", "-p", barrel.ComposerDirectory)
		if code != 0 {
			err = barrel.ExitError(code)
		}
		if err != nil {
			return fmt.Errorf("failed to create composer directory in barrel: %w", err)
		}
	}

	// start the container, and have it do nothing!
	code, err := stack.Run(ctx, stream.FromNil(), dockerx.RunFlags{Detach: true, AutoRemove: true}, "barrel", "tail", "-f", "/dev/null")
	if code != 0 {
		err = barrel.ExitError(code)
	}
	if err != nil {
		return fmt.Errorf("failed to start barrel in dummy mode: %w", err)
	}

	// when we are done, shut it down!
	defer func() {
		anyways, cancel := contextx.Anyways(ctx, time.Minute)
		defer cancel()

		// stop the container (even if the context was cancelled)
		if err := stack.DownAll(anyways, progress); err != nil {
			err = fmt.Errorf("unable to down stack: %w", err)
			e = errorsx.Combine(e, err)
		}
	}()

	return manager.bootstrap(ctx, progress, flags)
}

// TODO: Move this to the flags.
var drushVariants = []string{
	"drush/drush", "drush/drush:^13", "drush/drush:^12", "drush/drush:^11",
}

// bootstrap applies the initial flags induced by flags.
// Applies defaults to the flags.
func (provision *Manager) bootstrap(ctx context.Context, progress io.Writer, flags Profile) error {
	// TODO: Check if we can remove the easyrdf patch!
	flags.ApplyDefaults()

	if _, err := logging.LogMessage(progress, "Creating Composer Project"); err != nil {
		return fmt.Errorf("failed to log progress: %w", err)
	}
	{
		drupal := "drupal/recommended-project"
		if flags.Drupal != "" {
			drupal += ":" + flags.Drupal
		}
		err := provision.dependencies.Composer.Exec(ctx, progress, "create-project", drupal, ".")
		if err != nil {
			return fmt.Errorf("failed to create drupal project structure: %w", err)
		}
	}

	if _, err := logging.LogMessage(progress, "Configuring Composer"); err != nil {
		return fmt.Errorf("failed to log progress: %w", err)
	}
	{
		// needed for composer > 2.2
		err := provision.dependencies.Composer.Exec(ctx, progress, "config", "allow-plugins", "true")
		if err != nil {
			return fmt.Errorf("failed to configure composer: %w", err)
		}
	}

	if _, err := logging.LogMessage(progress, "Installing drush"); err != nil {
		return fmt.Errorf("failed to log progress: %w", err)
	}
	{
		for _, v := range drushVariants {
			err := provision.dependencies.Composer.TryInstall(ctx, progress, v)
			if errors.Is(err, composer.ErrNotInstalled) {
				continue
			}
			if err != nil {
				return fmt.Errorf("failed to install drush %q: %w", v, err)
			}
			break
		}
	}

	liquid := ingredient.GetLiquid(provision)

	var sqlDBURL = "mysql://" + liquid.SqlUsername + ":" + liquid.SqlPassword + "@sql/" + liquid.SqlDatabase

	// Use 'drush' to run the site-installation.
	// Here we need to use the username, password and database creds we made above.
	if _, err := logging.LogMessage(progress, "Running Drupal installation scripts"); err != nil {
		return fmt.Errorf("failed to log progress: %w", err)
	}
	{
		if err := provision.dependencies.Drush.Exec(
			ctx, progress,
			"site-install",
			"standard", "--yes", "--site-name="+liquid.Domain(),
			"--account-name="+liquid.DrupalUsername, "--account-pass="+liquid.DrupalPassword,
			"--db-url="+sqlDBURL,
			"-vvv",
		); err != nil {
			return fmt.Errorf("failed to execute drush site-install command: %w", err)
		}

		if err := provision.dependencies.Composer.FixPermission(ctx, progress); err != nil {
			return fmt.Errorf("failed to fix permissions: %w", err)
		}
	}

	// Rebuild the settings file
	if _, err := logging.LogMessage(progress, "Rebuilding Settings"); err != nil {
		return fmt.Errorf("failed to log progress: %w", err)
	}
	{
		if err := provision.dependencies.SystemManager.BuildSettings(ctx, progress); err != nil {
			return fmt.Errorf("failed to build settings: %w", err)
		}
	}

	// Create directory for ontologies
	if _, err := logging.LogMessage(progress, fmt.Sprintf("Creating %q", barrel.OntologyDirectory)); err != nil {
		return fmt.Errorf("failed to log progress: %w", err)
	}
	{
		if err := provision.dependencies.Barrel.BashScript(ctx, stream.NonInteractive(progress), "mkdir", "-p", barrel.OntologyDirectory); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
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
	if _, err := logging.LogMessage(progress, "Creating default adapter"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		if _, err := provision.dependencies.Adapters.SetAdapter(ctx, nil, provision.dependencies.Adapters.DefaultAdapter()); err != nil {
			return fmt.Errorf("failed to create default adapter: %w", err)
		}
	}

	if _, err := logging.LogMessage(progress, "Running initial cron"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		if err := provision.dependencies.Drush.Exec(ctx, progress, "core-cron"); err != nil {
			return fmt.Errorf("failed to run initial cron: %w", err)
		}
	}

	if _, err := logging.LogMessage(progress, "Provisioning is now complete"); err != nil {
		return fmt.Errorf("failed to log message: %w", err)
	}
	{
		if _, err := fmt.Fprintf(progress, "URL:                  %s\n", liquid.URL()); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if _, err := fmt.Fprintf(progress, "Username:             %s\n", liquid.DrupalUsername); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
		if _, err := fmt.Fprintf(progress, "Password:             %s\n", liquid.DrupalPassword); err != nil {
			return fmt.Errorf("failed to log message: %w", err)
		}
	}

	return nil
}
