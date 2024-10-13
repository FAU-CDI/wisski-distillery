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
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/tkw1536/pkglib/lazy"

	_ "embed"
)

type Resolver struct {
	component.Base
	dependencies struct {
		Instances  *instances.Instances
		Templating *templating.Templating
		Handling   *handling.Handling
		Auth       *auth.Auth
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

		MenuTitle:    "Resolver",
		MenuPriority: component.MenuResolver,
	}
}

//go:embed "resolver.html"
var resolverHTML []byte
var resolverTemplate = templating.Parse[resolverContext](
	"resolver.html", resolverHTML, nil,

	templating.Title("Resolver"),
	templating.Assets(assets.AssetsDefault),
)

type resolverContext struct {
	templating.RuntimeFlags
	wdresolve.IndexContext
}

var (
	menuResolver = component.MenuItem{Title: "Resolver", Path: "/wisski/get/"}
)

func (resolver *Resolver) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	// get the resolver template
	tpl := resolverTemplate.Prepare(
		resolver.dependencies.Templating,
		templating.Crumbs(
			menuResolver,
		),
	)
	t := tpl.Template()

	// extract a logger and the fallback
	logger := wdlog.Of(ctx)
	fallback := &resolvers.Regexp{
		Data: map[string]string{},
	}

	config := component.GetStill(resolver).Config

	// handle the default domain name!
	domainName := config.HTTP.PrimaryDomain
	if domainName != "" {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
		logger.Info().Str("name", domainName).Msg("registering default domain")
	}

	// handle the extra domains!
	for _, domain := range config.HTTP.ExtraDomains {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domain))] = fmt.Sprintf("https://$1.%s", domainName)
		logger.Info().Str("name", domainName).Msg("registering legacy domain")
	}

	p := wdresolve.ResolveHandler{
		HandleIndex: func(context wdresolve.IndexContext, w http.ResponseWriter, r *http.Request) {
			ctx := resolverContext{
				IndexContext: context,
			}

			if resolver.dependencies.Auth.CheckScope("", scopes.ScopeUserValid, r) != nil {
				ctx.IndexContext.Prefixes = nil
			}
			resolver.dependencies.Handling.WriteHTML(tpl.Context(r, ctx), nil, t, w, r)
		},

		Resolver: resolvers.InOrder{
			resolver,
			fallback,
		},

		TrustXForwardedProto: true,
	}

	return p, nil
}

func (resolver *Resolver) Target(uri string) string {
	return wdresolve.PrefixTarget(resolver, uri)
}

// Prefixes returns a cached list of prefixes
func (resolver *Resolver) Prefixes() (prefixes map[string]string) {
	return resolver.prefixes.Get(nil) // by precondition there always is a cached value
}
