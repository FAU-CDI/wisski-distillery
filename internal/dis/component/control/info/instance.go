package info

import (
	_ "embed"
	"net/http"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/info"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
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
	Info     info.WissKIInfo
}

func (info *Info) instance(r *http.Request) (is instanceContext, err error) {
	// find the slug as the last component of path!
	slug := strings.TrimSuffix(r.URL.Path, "/")
	slug = slug[strings.LastIndex(slug, "/")+1:]

	// find the instance itself!
	instance, err := info.Instances.WissKI(slug)
	if err == instances.ErrWissKINotFound {
		return is, httpx.ErrNotFound
	}
	if err != nil {
		return is, err
	}
	is.Instance = instance.Instance

	// get some more info about the wisski
	is.Info, err = instance.Info().Information(false)
	if err != nil {
		return is, err
	}

	// current time
	is.Time = time.Now().UTC()

	return
}
