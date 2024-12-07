//spellchecker:words composer
package composer

//spellchecker:words context errors github wisski distillery internal ingredient barrel drush mstore pkglib stream
import (
	"context"
	"errors"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/drush"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/mstore"
	"github.com/tkw1536/pkglib/stream"
)

// Drush implements commands related to drush
type Composer struct {
	ingredient.Base
	dependencies struct {
		Barrel *barrel.Barrel
		MStore *mstore.MStore
		Drush  *drush.Drush
	}
}

// Exec executes a composer command for the main composer package.
// Returns an error iff composer does not exit with 0.
func (composer *Composer) Exec(ctx context.Context, progress io.Writer, command ...string) error {
	return composer.exec(ctx, progress, append([]string{"--working-dir", barrel.ComposerDirectory}, command...)...)
}

// Exec executes a composer command for the wisski directory.
// Returns an error iff composer does not exit with 0.
func (composer *Composer) ExecWissKI(ctx context.Context, progress io.Writer, command ...string) error {
	return composer.exec(ctx, progress, append([]string{"--working-dir", barrel.WissKIDirectory}, command...)...)
}

func (composer *Composer) exec(ctx context.Context, progress io.Writer, command ...string) error {
	if err := composer.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), append([]string{"composer", "--no-interaction"}, command...)...); err != nil {
		return err
	}
	return nil
}

// FixPermissions fixes the permissions of the sites directory.
// This needs to be run after every installation of a composer module.
func (composer *Composer) FixPermission(ctx context.Context, progress io.Writer) error {
	composer.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), "chmod", "-R", "u+w", barrel.SitesDirectory)
	return nil
}

// Install attempts runs 'composer require' with the given arguments
// Spec is like a specification on the command line.
func (composer *Composer) Install(ctx context.Context, progress io.Writer, args ...string) error {
	if err := composer.FixPermission(ctx, progress); err != nil {
		return err
	}

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
