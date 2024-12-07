//spellchecker:words extras
package extras

//spellchecker:words context embed github wisski distillery internal phpx status ingredient
import (
	"context"
	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
)

type Stats struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

var (
	_ ingredient.WissKIFetcher = (*Stats)(nil)
)

//go:embed stats.php
var statsPHP string

// Get fetches all statistics from the server
func (stats *Stats) Get(ctx context.Context, server *phpx.Server) (data status.Statistics, err error) {
	err = stats.dependencies.PHP.ExecScript(ctx, server, &data, statsPHP, "export_statistics")
	return
}

func (stats *Stats) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Statistics, _ = stats.Get(flags.Context, flags.Server)
	return
}
