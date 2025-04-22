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

// Prefixes implements reading and writing prefix.
type Theme struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

//go:embed theme.php
var themePHP string

// Get returns the currently active theme.
func (t *Theme) Get(ctx context.Context, server *phpx.Server) (theme string, err error) {
	err = t.dependencies.PHP.ExecScript(
		ctx, server, &theme, themePHP,
		"get_active_theme",
	)
	return
}

func (t *Theme) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Theme, _ = t.Get(flags.Context, flags.Server)
	return
}
