package resolver

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/tkw1536/goprogram/stream"
)

type Resolver struct {
	component.Base

	Instances *instances.Instances

	prefixes        lazy.Lazy[map[string]string] // cached prefixes (from the server)
	RefreshInterval time.Duration

	handler lazy.Lazy[wdresolve.ResolveHandler] // handler
}

var (
	_ component.Servable = (*Resolver)(nil)
)

func (resolver *Resolver) Routes() []string { return []string{"/go/", "/wisski/get/"} }

func (resolver *Resolver) Handler(ctx context.Context, route string, io stream.IOStream) (http.Handler, error) {
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

		// start updating prefixes
		resolver.updatePrefixes(ctx, io)

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

// Prefixes returns a cached list of prefixes
func (resolver *Resolver) Prefixes() (prefixes map[string]string) {
	return resolver.prefixes.Get(nil) // by precondition there always is a cached value
}
