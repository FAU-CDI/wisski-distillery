package resolver

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/tkw1536/goprogram/stream"
)

type Resolver struct {
	component.ComponentBase

	Instances *instances.Instances

	prefixes        lazy.Lazy[map[string]string] // cached prefixes (from the server)
	RefreshInterval time.Duration

	handler lazy.Lazy[wdresolve.ResolveHandler] // handler
}

func (*Resolver) Name() string { return "resolver" }

func (resolver *Resolver) Routes() []string { return []string{"/go/", "/wisski/get/"} }

func (resolver *Resolver) Handler(route string, context context.Context, io stream.IOStream) (http.Handler, error) {
	var err error
	return resolver.handler.Get(func() (p wdresolve.ResolveHandler) {
		p.TrustXForwardedProto = true

		fallback := &resolvers.Regexp{
			Data: map[string]string{},
		}

		// handle the default domain name!
		domainName := resolver.Config.DefaultDomain
		if domainName != "" {
			fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
			io.Printf("registering default domain %s\n", domainName)
		}

		// handle the extra domains!
		for _, domain := range resolver.Config.SelfExtraDomains {
			fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domain))] = fmt.Sprintf("https://$1.%s", domainName)
			io.Printf("registering legacy domain %s\n", domain)
		}

		go resolver.updatePrefixes(io, context)

		// resolve the prefixes
		p.Resolver = resolvers.InOrder{
			resolver,
			fallback,
		}
		return p
	}), err
}

func (resolver *Resolver) updatePrefixes(io stream.IOStream, ctx context.Context) {
	t := time.NewTicker(resolver.RefreshInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			io.Println("resolver: Reloading prefixes from database")
			prefixes, _ := resolver.AllPrefixes()
			resolver.prefixes.Set(prefixes)
		case <-ctx.Done():
			return
		}
	}
}

func (resolver *Resolver) Target(uri string) string {
	return wdresolve.PrefixTarget(resolver, uri)
}

// Prefixes returns a cached list of prefixes
func (resolver *Resolver) Prefixes() (prefixes map[string]string) {
	return resolver.prefixes.Get(func() map[string]string {
		prefixes, _ := resolver.AllPrefixes()
		return prefixes
	})
}

// AllPrefixes returns a list of all prefixes from the server.
// Prefixes may be cached on the server
func (resolver *Resolver) AllPrefixes() (map[string]string, error) {
	instances, err := resolver.Instances.All()
	if err != nil {
		return nil, err
	}

	gPrefixes := make(map[string]string)
	var lastErr error
	for _, instance := range instances {
		if instance.NoPrefix() {
			continue
		}
		url := instance.URL().String()

		// failed to fetch prefixes for this particular instance
		// => skip it!
		prefixes, err := instance.PrefixesCached()
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
