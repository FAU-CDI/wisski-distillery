package admin

import (
	_ "embed"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/julienschmidt/httprouter"
)

//go:embed "html/instance.html"
var instanceTemplateString string
var instanceTemplate = static.AssetsAdmin.MustParseShared(
	"instance.html",
	instanceTemplateString,
)

type instanceContext struct {
	custom.BaseContext

	Instance models.Instance
	Info     status.WissKI
}

func (admin *Admin) instance(r *http.Request) (is instanceContext, err error) {
	admin.Dependencies.Custom.Update(&is, r)

	// find the instance itself!
	slug := httprouter.ParamsFromContext(r.Context()).ByName("slug")
	instance, err := admin.Dependencies.Instances.WissKI(r.Context(), slug)
	if err == instances.ErrWissKINotFound {
		return is, httpx.ErrNotFound
	}
	if err != nil {
		return is, err
	}
	is.Instance = instance.Instance

	// get some more info about the wisski
	is.Info, err = instance.Info().Information(r.Context(), false)
	if err != nil {
		return is, err
	}

	return
}
