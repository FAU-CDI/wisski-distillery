//spellchecker:words extras
package extras

//spellchecker:words context github wisski distillery internal phpx status ingredient embed
import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"

	_ "embed"
)

// Version implements reading the current drupal version.
type Version struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

// Get returns the currently active theme.
func (v *Version) Get(ctx context.Context, server *phpx.Server) (version string, err error) {
	err = v.dependencies.PHP.EvalCode(
		ctx, server, &version, "return Drupal::VERSION; ",
	)
	return
}

func (v *Version) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.DrupalVersion, _ = v.Get(flags.Context, flags.Server)
	return
}
