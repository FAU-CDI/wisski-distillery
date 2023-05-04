package list

import (
	"context"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/api"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/auth/scopes"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
)

// API implements an API to list all instances
type API struct {
	component.Base
	Dependencies struct {
		ListInstances *ListInstances
		Auth          *auth.Auth
	}
}

func (lia *API) Routes() component.Routes {
	return component.Routes{
		Prefix: "/api/v1/instances/directory",
		Exact:  true,
	}
}

// APISystem represents a system returned by the api
type APISystem struct {
	Slug    string
	URL     string
	Tagline string

	EntityCount int
	BundleCount int
	LastEdit    time.Time
}

func (a *API) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return &api.Handler[[]APISystem]{
		Config: a.Config,
		Auth:   a.Dependencies.Auth,

		Methods: []string{"GET"},
		Scope:   scopes.ScopeInstanceDirectory,

		Handler: func(s string, r *http.Request) ([]APISystem, error) {
			var statuses []status.WissKI
			if a.Dependencies.ListInstances.ShouldShowList(r) {
				statuses = a.Dependencies.ListInstances.infos.Get(nil)
			}

			if len(statuses) == 0 {
				return []APISystem{}, nil
			}

			infos := make([]APISystem, len(statuses))
			for i, status := range statuses {
				infos[i].Slug = status.Slug
				infos[i].URL = status.URL
				infos[i].EntityCount = status.Statistics.Bundles.TotalCount()
				infos[i].BundleCount = status.Statistics.Bundles.TotalBundles
				infos[i].LastEdit = status.Statistics.Bundles.LastEdit().Time
			}
			return infos, nil
		},
	}, nil
}
