package home

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
)

type Home struct {
	component.Base
	Dependencies struct {
		Templating *templating.Templating
		Instances  *instances.Instances
	}

	instanceNames lazy.Lazy[map[string]struct{}] // instance names
	homeInstances lazy.Lazy[[]status.WissKI]     // list of home instances (updated via cron)
}

var (
	_ component.Routeable = (*Home)(nil)
)

func (*Home) Routes() component.Routes {
	return component.Routes{
		Prefix:          "/",
		MatchAllDomains: true,
		CSRF:            false,

		MenuTitle:    "WissKI Distillery",
		MenuPriority: component.MenuHome,
	}
}

var (
	menuHome = component.MenuItem{Title: "WissKI Distillery", Path: "/"}
)

func (home *Home) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	// generate a default handler
	dflt, err := home.loadRedirect(ctx)
	if err != nil {
		return nil, err
	}
	dflt.Fallback = home.publicHandler(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug, ok := home.Config.HTTP.SlugFromHost(r.Host)
		switch {
		case !ok:
			http.NotFound(w, r)
		case slug != "":
			home.serveWissKI(w, slug, r)
		default:
			dflt.ServeHTTP(w, r)
		}
	}), nil
}

func (home *Home) instanceMap(ctx context.Context) (map[string]struct{}, error) {
	wissKIs, err := home.Dependencies.Instances.All(ctx)
	if err != nil {
		return nil, err
	}

	names := make(map[string]struct{}, len(wissKIs))
	for _, w := range wissKIs {
		names[w.Slug] = struct{}{}
	}
	return names, nil
}

func (home *Home) serveWissKI(w http.ResponseWriter, slug string, r *http.Request) {
	if _, ok := home.instanceNames.Get(nil)[slug]; !ok {
		// Get(nil) guaranteed to work by precondition
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "WissKI %q not found\n", slug)
		return
	}

	w.WriteHeader(http.StatusBadGateway)
	fmt.Fprintf(w, "WissKI %q is currently offline\n", slug)
}
