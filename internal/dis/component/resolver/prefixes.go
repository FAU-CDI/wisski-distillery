package resolver

import (
	"context"
)

func (resolver *Resolver) TaskName() string {
	return "reloading prefixes"
}

func (resolver *Resolver) Cron(ctx context.Context) error {
	prefixes, err := resolver.AllPrefixes(ctx)
	if err != nil {
		return err
	}

	resolver.prefixes.Set(prefixes)
	return nil
}

// AllPrefixes returns a list of all prefixes from the server.
// Prefixes may be cached on the server
func (resolver *Resolver) AllPrefixes(ctx context.Context) (map[string]string, error) {
	instances, err := resolver.Dependencies.Instances.All(ctx)
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
