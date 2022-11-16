package extras

import (
	_ "embed"
	"log"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
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
func (stats *Stats) Get(server *phpx.Server) (data ingredient.Statistics, err error) {
	err = stats.PHP.ExecScript(server, &data, statsPHP, "export_statistics")
	if err != nil {
		log.Println(err)
	}
	return
}

func (stats *Stats) Fetch(flags ingredient.FetchFlags, info *ingredient.Information) (err error) {
	if flags.Quick {
		return
	}

	info.Statistics, _ = stats.Get(flags.Server)
	return
}
