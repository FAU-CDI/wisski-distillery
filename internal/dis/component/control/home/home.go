package home

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/pkg/lazy"
	"github.com/tkw1536/goprogram/stream"
)

type Home struct {
	component.Base

	Instances *instances.Instances

	RefreshInterval time.Duration

	redirect      lazy.Lazy[*Redirect]
	instanceNames lazy.Lazy[map[string]struct{}]
	homeBytes     lazy.Lazy[[]byte]
}

func (*Home) Routes() []string { return []string{"/"} }

func (home *Home) Handler(ctx context.Context, route string, io stream.IOStream) (http.Handler, error) {
	home.updateRedirect(ctx, io)
	home.updateInstances(ctx, io)
	home.updateRender(ctx, io)
	return home, nil
}

func (home *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug, ok := home.Config.SlugFromHost(r.Host)
	switch {
	case !ok:
		http.NotFound(w, r)
	case slug != "":
		home.serveWissKI(w, slug, r)
	default:
		home.serveRoot(w, r)
	}
}

func (home *Home) serveRoot(w http.ResponseWriter, r *http.Request) {
	// not the root url => server the fallback
	if !(r.URL.Path == "" || r.URL.Path == "/") {
		home.redirect.Get(nil).ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusAccepted)
	w.Write(home.homeBytes.Get(nil))
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
