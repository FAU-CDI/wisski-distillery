package home

import (
	"bytes"
	"context"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/static"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

func (home *Home) updateInstances(ctx context.Context, io stream.IOStream) {
	go func() {
		for t := range timex.TickContext(ctx, home.RefreshInterval) {
			io.Printf("[%s]: reloading instance list\n", t.Format(time.Stamp))

			names, _ := home.instanceMap()
			home.instanceNames.Set(names)
		}
	}()
}

func (home *Home) instanceMap() (map[string]struct{}, error) {
	wissKIs, err := home.Instances.All()
	if err != nil {
		return nil, err
	}

	names := make(map[string]struct{}, len(wissKIs))
	for _, w := range wissKIs {
		names[w.Slug] = struct{}{}
	}
	return names, nil
}

func (home *Home) updateRender(ctx context.Context, io stream.IOStream) {
	go func() {
		for t := range timex.TickContext(ctx, home.RefreshInterval) {
			io.Printf("[%s]: reloading home render\n", t.Format(time.Stamp))

			bytes, _ := home.homeRender()
			home.homeBytes.Set(bytes)
		}
	}()
}

//go:embed "home.html"
var homeHTMLStr string
var homeTemplate = static.AssetsHomeHome.MustParse(homeHTMLStr)

func (home *Home) homeRender() ([]byte, error) {
	var context HomeContext

	// setup a couple of static things
	context.Time = time.Now().UTC()
	context.SelfRedirect = home.Config.SelfRedirect.String()

	// find all the WissKIs
	wissKIs, err := home.Instances.All()
	if err != nil {
		return nil, err
	}
	context.Instances = make([]wisski.WissKIInfo, len(wissKIs))

	// determine their infos
	var eg errgroup.Group
	for i, instance := range wissKIs {
		i := i
		wissKI := instance
		eg.Go(func() (err error) {
			context.Instances[i], err = wissKI.Info(true)
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
	Instances []wisski.WissKIInfo

	Time time.Time

	SelfRedirect string
}
