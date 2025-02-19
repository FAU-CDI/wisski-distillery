//spellchecker:words drush
package drush

//spellchecker:words context github wisski distillery internal ingredient barrel pkglib stream
import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/tkw1536/pkglib/stream"
)

// Drush implements commands related to drush
type Drush struct {
	ingredient.Base
	dependencies struct {
		Barrel *barrel.Barrel
		PHP    *php.PHP
	}
}

// Enable enables the given drush modules
func (drush *Drush) Enable(ctx context.Context, progress io.Writer, modules ...string) error {
	return drush.Exec(ctx, progress, append([]string{"pm-enable", "--yes"}, modules...)...)
}

func (drush *Drush) Exec(ctx context.Context, progress io.Writer, command ...string) error {
	script := append([]string{"drush"}, command...)
	return drush.dependencies.Barrel.ShellScript(ctx, stream.NonInteractive(progress), script...)
}
