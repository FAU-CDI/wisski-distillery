package home

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/list"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
)

type Home struct {
	component.Base
	dependencies struct {
		ListInstances *list.ListInstances
		Templating    *templating.Templating
		Handling      *handling.Handling
	}
}

var (
	_ component.Routeable = (*Home)(nil)
)

func (home *Home) Routes() component.Routes {
	return component.Routes{
		Prefix:          "/",
		MatchAllDomains: true,
		CSRF:            false,

		MenuTitle:    home.Config.Home.Title,
		MenuSticky:   true,
		MenuPriority: component.MenuHome,
	}
}

func (home *Home) HandleRoute(ctx context.Context, route string) (http.Handler, error) {
	// generate a default handler
	dflt, err := home.loadRedirect(ctx)
	if err != nil {
		return nil, err
	}
	dflt.Fallback = home.publicHandler(ctx)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug, ok := home.Config.HTTP.NormSlugFromHost(r.Host)
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

func (home *Home) serveWissKI(w http.ResponseWriter, slug string, r *http.Request) {
	if _, ok := home.dependencies.ListInstances.Names()[slug]; !ok {
		// Get(nil) guaranteed to work by precondition
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "WissKI %q not found\n", slug)
		return
	}

	w.WriteHeader(http.StatusBadGateway)
	fmt.Fprintf(w, "WissKI %q is currently offline\n", slug)
}
