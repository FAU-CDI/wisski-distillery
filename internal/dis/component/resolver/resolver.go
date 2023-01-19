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
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templates"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/rs/zerolog"

	_ "embed"
)

type Resolver struct {
	component.Base
	Dependencies struct {
		Instances  *instances.Instances
		Templating *templates.Templating
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
var resolverTemplate = templates.Parse[resolverContext]("resolver.html", resolverHTML, assets.AssetsDefault)

type resolverContext struct {
	templates.BaseContext
	wdresolve.IndexContext
}

func (resolver *Resolver) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	tpl := resolverTemplate.Prepare(resolver.Dependencies.Templating, templates.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "Resolver", Path: "/wisski/get/"},
		},
	})
	logger := zerolog.Ctx(ctx)

	var p wdresolve.ResolveHandler
	var err error

	p.HandleIndex = func(context wdresolve.IndexContext, w http.ResponseWriter, r *http.Request) {
		ctx := resolverContext{
			IndexContext: context,
		}
		if !resolver.Dependencies.Auth.Has(auth.User, r) {
			ctx.IndexContext.Prefixes = nil
		}

		tpl.Execute(w, r, ctx)
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
