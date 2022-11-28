package resolver

import (
	"context"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
	"github.com/tkw1536/goprogram/stream"
)

// updatePrefixes starts updating prefixes
func (resolver *Resolver) updatePrefixes(ctx context.Context, io stream.IOStream) {
	go func() {
		for t := range timex.TickContext(ctx, resolver.RefreshInterval) {
			io.Printf("[%s]: reloading prefixes\n", t.Format(time.Stamp))

			err := (func() (err error) {
				ctx, cancel := context.WithTimeout(ctx, resolver.RefreshInterval)
				defer cancel()

				prefixes, err := resolver.AllPrefixes(ctx)
				if err != nil {
					return err
				}

				resolver.prefixes.Set(prefixes)
				return nil
			})()
			if err != nil {
				io.EPrintf("error reloading prefixes: ", err.Error())
			}
		}
	}()
}

// AllPrefixes returns a list of all prefixes from the server.
// Prefixes may be cached on the server
func (resolver *Resolver) AllPrefixes(ctx context.Context) (map[string]string, error) {
	instances, err := resolver.Instances.All(ctx)
	if err != nil {
		return nil, err
	}

	gPrefixes := make(map[string]string)
	var lastErr error
	for _, instance := range instances {
		if instance.Prefixes().NoPrefix() {
			continue
		}
		url := instance.URL().String()

		// failed to fetch prefixes for this particular instance
		// => skip it!
		prefixes, err := instance.Prefixes().AllCached(ctx)
		if err != nil {
			lastErr = err
			continue
		}

		for _, p := range prefixes {
			gPrefixes[p] = url
		}
	}

	return gPrefixes, lastErr
}
