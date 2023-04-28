package news

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/tkw1536/pkglib/httpx"
)

type API struct {
	component.Base
}

var (
	_ component.Routeable = (*API)(nil)
)

func (api *API) Routes() component.Routes {
	return component.Routes{
		Prefix:    "/api/v1/news/",
		Exact:     true,
		Decorator: api.Config.HTTP.APIDecorator("GET"),
	}
}

func (api *API) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	items, err := Items()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return httpx.Response{
		ContentType: "application/json",
		Body:        data,
	}, nil
}
