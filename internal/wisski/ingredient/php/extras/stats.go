package extras

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

	PHP *php.PHP
}

//go:embed stats.php
var statsPHP string

// Get fetches all statistics from the server
func (stats *Stats) Get(ctx context.Context, server *phpx.Server) (data status.Statistics, err error) {
	err = stats.PHP.ExecScript(ctx, server, &data, statsPHP, "export_statistics")
	return
}

func (stats *Stats) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) (err error) {
	if flags.Quick {
		return
	}

	info.Statistics, _ = stats.Get(flags.Context, flags.Server)
	return
}
