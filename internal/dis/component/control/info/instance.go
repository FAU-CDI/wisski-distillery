package info

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/gorilla/mux"
)

//go:embed "html/instance.html"
var instanceTemplateString string
var instanceTemplate = static.AssetsControlInstance.MustParseShared(
	"instance.html",
	instanceTemplateString,
)

type instanceContext struct {
	Time time.Time

	Instance models.Instance
	Info     status.WissKI
}

func (info *Info) instance(r *http.Request) (is instanceContext, err error) {
	// find the instance itself!
	instance, err := info.Instances.WissKI(r.Context(), mux.Vars(r)["slug"])
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

	// current time
	is.Time = time.Now().UTC()

	return
}
