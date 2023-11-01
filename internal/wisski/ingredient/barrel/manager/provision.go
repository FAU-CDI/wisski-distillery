package manager

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/composer"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/pkglib/contextx"
	"github.com/tkw1536/pkglib/stream"
)

// Provision provisions this instance with the given flags.
//
// Provision assumes that the instance does not yet exist, and may fail with an existing instance.
//
// Provision applies defaults to flags, to ensure some values are set
func (manager *Manager) Provision(ctx context.Context, progress io.Writer, system models.System, flags Profile) error {
	// Force building and applying the system!
	if err := manager.dependencies.SystemManager.Apply(ctx, progress, system, false); err != nil {
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

	// Apply the defaults to the flags
	flags.ApplyDefaults()
	return manager.bootstrap(ctx, progress, flags)
}

// TODO: Move this to the flags
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
			if err == composer.ErrNotInstalled {
				continue
			}
			if err != nil {
				return err
			}
			break
		}
	}

	var sqlDBURL = "mysql://" + provision.SqlUsername + ":" + provision.SqlPassword + "@sql/" + provision.SqlDatabase

	// Use 'drush' to run the site-installation.
	// Here we need to use the username, password and database creds we made above.
	logging.LogMessage(progress, "Running Drupal installation scripts")
	{
		if err := provision.dependencies.Drush.Exec(
			ctx, progress,
			"site-install",
			"standard", "--yes", "--site-name="+provision.Domain(),
			"--account-name="+provision.DrupalUsername, "--account-pass="+provision.DrupalPassword,
			"--db-url="+sqlDBURL,
		); err != nil {
			return err
		}

		if err := provision.dependencies.Composer.FixPermission(ctx, progress); err != nil {
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
		if err := provision.dependencies.Adapters.CreateDistilleryAdapter(ctx, nil, extras.DistilleryAdapter{
			Label:             "Default WissKI Distillery Adapter",
			MachineName:       "default",
			Description:       "Default Adapter for " + provision.Domain(),
			InstanceDomain:    provision.Domain(),
			GraphDBRepository: provision.GraphDBRepository,
			GraphDBUsername:   provision.GraphDBUsername,
			GraphDBPassword:   provision.GraphDBPassword,
		}); err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Updating TRUSTED_HOST_PATTERNS in settings.php")
	{
		if err := provision.dependencies.Settings.SetTrustedDomain(ctx, nil, provision.Domain()); err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Running initial cron")
	{
		if err := provision.dependencies.Drush.Exec(ctx, progress, "core-cron"); err != nil {
			return err
		}
	}

	logging.LogMessage(progress, "Provisioning is now complete")
	{
		fmt.Fprintf(progress, "URL:                  %s\n", provision.URL())
		fmt.Fprintf(progress, "Username:             %s\n", provision.DrupalUsername)
		fmt.Fprintf(progress, "Password:             %s\n", provision.DrupalPassword)
	}

	return nil
}
