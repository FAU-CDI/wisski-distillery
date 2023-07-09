package composer

import (
	"context"
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/tkw1536/pkglib/stream"
)

// Drush implements commands related to drush
type Composer struct {
	ingredient.Base
	Dependencies struct {
		Barrel *barrel.Barrel
		//		PHP    *php.PHP
	}
}

// Exec executes a composer command
func (composer *Composer) Exec(ctx context.Context, progress io.Writer, command ...string) error {
	if err := composer.Dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), append([]string{"composer", "--no-interaction", "--working-dir", barrel.ComposerDirectory}, command...)...); err != nil {
		return err
	}
	return nil
}

// FixPermissions fixes the permissions of the sites directory.
// This needs to be run after every installation of a composer module.
func (composer *Composer) FixPermission(ctx context.Context, progress io.Writer) error {
	composer.Dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), "chmod", "-R", "u+w", barrel.SitesDirectory)
	return nil
}

// Install attempts runs 'composer require' with the given arguments
// Spec is like a specification on the command line.
func (composer *Composer) Install(ctx context.Context, progress io.Writer, args ...string) error {
	composer.FixPermission(ctx, progress)

	requires := append([]string{"require"}, args...)
	if err := composer.Exec(ctx, progress, requires...); err != nil {
		return err
	}
	return nil
}

var ErrNotInstalled = errors.New("Composer: Not installed")

// TryInstall attempts to install the given package.
// If it cannot be installed, returns ErrNotInstalled.
func (composer *Composer) TryInstall(ctx context.Context, progress io.Writer, spec string) error {
	if err := composer.Exec(ctx, io.Discard, "require", "--dry-run", spec); err != nil {
		return ErrNotInstalled
	}

	return composer.Install(ctx, progress, spec)
}
