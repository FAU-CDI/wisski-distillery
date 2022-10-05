package resolver

import (
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

	prefixes lazy.Lazy[map[string]string]        // cached prefixes (from the server)
	handler  lazy.Lazy[wdresolve.ResolveHandler] // handler
}

func (*Resolver) Name() string { return "resolver" }

func (resolver *Resolver) Routes() []string { return []string{"/go/", "/wisski/get/"} }

func (resolver *Resolver) Handler(route string, io stream.IOStream) (http.Handler, error) {
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

		// resolve the prefixes
		p.Resolver = resolvers.InOrder{
			resolver,
			fallback,
		}
		return p
	}), err
}

func (resolver *Resolver) Target(uri string) string {
	return wdresolve.PrefixTarget(resolver, uri)
}

// allow reloading prefixes from the server every minute
const prefixesRefresh = time.Minute

func (resolver *Resolver) Prefixes() (prefixes map[string]string) {
	// reset the prefixes after a specific time, but only if requested
	resolver.prefixes.ResetAfter(prefixesRefresh)
	return resolver.prefixes.Get(resolver.freshPrefixes)
}

func (resolver *Resolver) freshPrefixes() map[string]string {
	instances, err := resolver.Instances.All()
	if err != nil {
		return nil
	}

	gPrefixes := make(map[string]string)
	for _, instance := range instances {
		url := instance.URL().String()

		// failed to fetch prefixes for this particular instance
		// => skip it!
		prefixes, err := instance.PrefixesCached()
		if err != nil {
			continue
		}

		for _, p := range prefixes {
			gPrefixes[url] = p
		}
	}
	return gPrefixes
}
