//spellchecker:words info
package info

//spellchecker:words context reflect sync atomic time github wisski distillery internal phpx status wdlog ingredient pkglib sema golang errgroup
import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/tkw1536/pkglib/sema"
	"golang.org/x/sync/errgroup"
)

type Info struct {
	ingredient.Base
	dependencies struct {
		PHP      *php.PHP
		Fetchers []ingredient.WissKIFetcher
	}
}

var (
	_ ingredient.WissKIFetcher = (*Info)(nil)
)

// Information fetches information about this WissKI.
func (nfo *Info) Information(ctx context.Context, quick bool) (info status.WissKI, err error) {
	// setup flags
	flags := ingredient.FetcherFlags{
		Quick:   quick,
		Context: ctx,
	}

	var serversUsed uint64
	pool := sema.Pool[*phpx.Server]{
		// limit the number of processes running in this container
		// to avoid long overheads
		Limit: 5,
		New: func() *phpx.Server {
			atomic.AddUint64(&serversUsed, 1)
			return nfo.dependencies.PHP.NewServer()
		},
		Discard: func(s *phpx.Server) {
			s.Close()
		},
	}
	defer pool.Close()

	// setup a dictionary to record data about how long each operation took.
	// we use a slice as opposed to a map to avoid having to mutex!
	fetcherTimes := make([]time.Duration, len(nfo.dependencies.Fetchers))
	recordTime := func(i int) func() {
		start := time.Now()
		return func() {
			fetcherTimes[i] = time.Since(start)
		}
	}

	start := time.Now()
	{
		var group errgroup.Group
		for i, fetcher := range nfo.dependencies.Fetchers {
			fetcher, flags, i := fetcher, flags, i
			group.Go(func() error {
				// quick: don't need to create servers
				if flags.Quick {
					defer recordTime(i)()

					err := fetcher.Fetch(flags, &info)
					if err != nil {
						return fmt.Errorf("fetcher %s (quick): %w", reflect.TypeOf(fetcher), err)
					}
					return nil
				}

				// complete: need to use a server from the pool
				return pool.Use(func(s *phpx.Server) error {
					defer recordTime(i)()
					flags.Server = s

					err := fetcher.Fetch(flags, &info)
					if err != nil {
						return fmt.Errorf("fetcher %s (pool): %w", reflect.TypeOf(fetcher), err)
					}
					return nil
				})
			})
		}

		// wait for all the results
		err = group.Wait()
	}
	took := time.Since(start)

	var tookSum time.Duration

	// get a map of how long each fetcher took
	times := make(map[string]time.Duration, len(nfo.dependencies.Fetchers))
	for i, fetcher := range nfo.dependencies.Fetchers {
		tookSum += fetcherTimes[i]
		times[fetcher.Name()] = fetcherTimes[i]
	}

	// compute the ratio taken
	tookRatio := float64(took) / float64(tookSum)

	// and send it to debugging output
	wdlog.Of(ctx).Debug(
		"ran information fetchers",

		"servers", serversUsed,
		"fetchers_took_ms", times,
		"took_ms", took,
		"took_sum_ms", tookSum,
		"took_ratio", tookRatio,
		"quick", quick,
	)

	return
}

func (nfo *Info) Fetch(flags ingredient.FetcherFlags, info *status.WissKI) error {
	liquid := ingredient.GetLiquid(nfo)

	info.Time = time.Now().UTC()
	info.Slug = liquid.Slug
	info.URL = liquid.URL().String()
	return nil
}
