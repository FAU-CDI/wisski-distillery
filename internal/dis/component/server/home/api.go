package home

import (
	"context"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/tkw1536/pkglib/httpx"
)

type API struct {
	component.Base

	Dependencies struct {
		Home *Home
	}
}

var (
	_ component.Routeable = (*API)(nil)
)

func (api *API) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/api/v1/systems",
		Exact:     true,
		Decorator: api.Config.HTTP.APIDecorator("GET"),
	}
}

type APISystemInfo struct {
	Slug    string
	URL     string
	Tagline string

	EntityCount int
	BundleCount int
	LastEdit    time.Time
}

func (api *API) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	return httpx.JSON(func(r *http.Request) ([]APISystemInfo, error) {
		var statuses []status.WissKI
		if api.Dependencies.Home.ShouldShowList(r) {
			statuses = api.Dependencies.Home.homeInstances.Get(nil)
		}

		if len(statuses) == 0 {
			return []APISystemInfo{}, nil
		}

		infos := make([]APISystemInfo, len(statuses))
		for i, status := range statuses {
			infos[i].Slug = status.Slug
			infos[i].URL = status.URL
			infos[i].EntityCount = status.Statistics.Bundles.TotalCount()
			infos[i].BundleCount = status.Statistics.Bundles.TotalBundles
			infos[i].LastEdit = status.Statistics.Bundles.LastEdit().Time
		}
		return infos, nil
	}), nil
}
