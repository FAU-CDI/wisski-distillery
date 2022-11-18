package info

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"golang.org/x/sync/errgroup"
)

type Info struct {
	ingredient.Base

	PHP      *php.PHP
	Fetchers []ingredient.WissKIFetcher

	Analytics *lazy.PoolAnalytics
}

// Information fetches information about this WissKI.
// TODO: Rework this to be able to determine what kind of information is available.
func (wisski *Info) Information(quick bool) (info status.Information, err error) {
	// setup flags
	flags := ingredient.FetcherFlags{
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

func (wisski *Info) Fetch(flags ingredient.FetcherFlags, info *status.Information) error {
	info.Time = time.Now().UTC()
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()
	return nil
}
