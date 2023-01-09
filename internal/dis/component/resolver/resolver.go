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
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/rs/zerolog"

	_ "embed"
)

type Resolver struct {
	component.Base
	Dependencies struct {
		Instances *instances.Instances
		Custom    *custom.Custom
	}

	prefixes        lazy.Lazy[map[string]string] // cached prefixes (from the server)
	RefreshInterval time.Duration
}

var (
	_ component.Routeable = (*Resolver)(nil)
	_ component.Cronable  = (*Resolver)(nil)
)

func (resolver *Resolver) Routes() component.Routes {
	return component.Routes{
		Prefix:  "/wisski/get/",
		Aliases: []string{"/go/"},
		CSRF:    false,
	}
}

//go:embed "resolver.html"
var resolverHTMLStr string
var resolverTemplate = static.AssetsHome.MustParseShared("resolver.html", resolverHTMLStr)

type resolverContext struct {
	custom.BaseContext
	wdresolve.IndexContext
}

func (resolver *Resolver) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	resolverTemplate := resolver.Dependencies.Custom.Template(resolverTemplate)

	logger := zerolog.Ctx(ctx)

	var p wdresolve.ResolveHandler
	var err error

	p.HandleIndex = func(context wdresolve.IndexContext, w http.ResponseWriter, r *http.Request) {
		ctx := resolverContext{
			IndexContext: context,
		}
		resolver.Dependencies.Custom.Update(&ctx, r)

		httpx.WriteHTML(ctx, nil, resolverTemplate, "", w, r)
	}
	p.TrustXForwardedProto = true

	fallback := &resolvers.Regexp{
		Data: map[string]string{},
	}

	// handle the default domain name!
	domainName := resolver.Config.DefaultDomain
	if domainName != "" {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
		logger.Info().Str("name", domainName).Msg("registering default domain")
	}

	// handle the extra domains!
	for _, domain := range resolver.Config.SelfExtraDomains {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domain))] = fmt.Sprintf("https://$1.%s", domainName)
		logger.Info().Str("name", domainName).Msg("registering legacy domain")
	}

	// resolve the prefixes
	p.Resolver = resolvers.InOrder{
		resolver,
		fallback,
	}
	return p, err
}

func (resolver *Resolver) Target(uri string) string {
	return wdresolve.PrefixTarget(resolver, uri)
}

// Prefixes returns a cached list of prefixes
func (resolver *Resolver) Prefixes() (prefixes map[string]string) {
	return resolver.prefixes.Get(nil) // by precondition there always is a cached value
}
