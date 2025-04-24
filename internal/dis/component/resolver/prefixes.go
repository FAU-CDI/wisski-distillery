//spellchecker:words resolver
package resolver

//spellchecker:words context
import (
	"context"
	"fmt"
)

func (resolver *Resolver) TaskName() string {
	return "reloading prefixes"
}

func (resolver *Resolver) Cron(ctx context.Context) error {
	prefixes, err := resolver.AllPrefixes(ctx)
	if err != nil {
		return fmt.Errorf("failed resolve all prefixes: %w", err)
	}

	resolver.prefixes.Set(prefixes)
	return nil
}

// Prefixes may be cached on the server.
func (resolver *Resolver) AllPrefixes(ctx context.Context) (map[string]string, error) {
	instances, err := resolver.dependencies.Instances.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all prefixes: %w", err)
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
