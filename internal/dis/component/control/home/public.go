package home

import (
	"bytes"
	"context"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"golang.org/x/sync/errgroup"
)

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

//go:embed "home.html"
var homeHTMLStr string
var homeTemplate = static.AssetsHome.MustParseShared("home.html", homeHTMLStr)

func (home *Home) homeRender(ctx context.Context) ([]byte, error) {
	var context HomeContext

	// setup a couple of static things
	context.Time = time.Now().UTC()
	context.SelfRedirect = home.Config.SelfRedirect.String()

	// find all the WissKIs
	wissKIs, err := home.Dependencies.Instances.All(ctx)
	if err != nil {
		return nil, err
	}
	context.Instances = make([]status.WissKI, len(wissKIs))

	// determine their infos
	var eg errgroup.Group
	for i, instance := range wissKIs {
		i := i
		wissKI := instance
		eg.Go(func() (err error) {
			context.Instances[i], err = wissKI.Info().Information(ctx, false)
			return
		})
	}
	eg.Wait()

	// render the template
	var buffer bytes.Buffer
	homeTemplate.Execute(&buffer, context)
	return buffer.Bytes(), nil
}

type HomeContext struct {
	Instances []status.WissKI

	Time time.Time

	SelfRedirect string
}
