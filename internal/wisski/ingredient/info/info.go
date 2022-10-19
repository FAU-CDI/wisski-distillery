package info

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/drush"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"golang.org/x/sync/errgroup"
)

type Info struct {
	ingredient.Base

	PHP      *php.PHP
	Fetchers []ingredient.Fetcher

	Barrel      *barrel.Barrel
	Locker      *locker.Locker
	Drush       *drush.Drush
	Prefixes    *extras.Prefixes
	Pathbuilder *extras.Pathbuilder
}

// TODO: Use the information struct globally
type WissKIInfo = ingredient.Information

// Information fetches information about this WissKI.
// TODO: Rework this to be able to determine what kind of information is available.
func (wisski *Info) Information(quick bool) (info WissKIInfo, err error) {
	// setup flags
	flags := ingredient.FetchFlags{
		Quick: quick,
	}

	// potentially setup a new server
	if !flags.Quick {
		flags.Server = wisski.PHP.NewServer()
		if err == nil {
			defer flags.Server.Close()
		}
	}

	// run all the fetchers!
	var group errgroup.Group
	for _, fetcher := range wisski.Fetchers {
		fetcher := fetcher
		group.Go(func() error {
			return fetcher.Fetch(flags, &info)
		})
	}

	err = group.Wait()
	return
}

func (wisski *Info) Fetch(flags ingredient.FetchFlags, info *ingredient.Information) error {
	info.Time = time.Now().UTC()
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()
	return nil
}
