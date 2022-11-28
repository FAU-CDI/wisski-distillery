package home

import (
	"bytes"
	"context"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/status"
	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

func (home *Home) updateInstances(ctx context.Context, io stream.IOStream) {
	go func() {
		for t := range timex.TickContext(ctx, home.RefreshInterval) {
			io.Printf("[%s]: reloading instance list\n", t.Format(time.Stamp))

			err := (func() error {
				ctx, cancel := context.WithTimeout(ctx, home.RefreshInterval)
				defer cancel()

				names, err := home.instanceMap(ctx)
				if err != nil {
					return err
				}

				home.instanceNames.Set(names)
				return nil
			})()
			if err != nil {
				io.EPrintf("error reloading instances: ", err.Error())
			}
		}
	}()
}

func (home *Home) instanceMap(ctx context.Context) (map[string]struct{}, error) {
	wissKIs, err := home.Instances.All(ctx)
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
			io.Printf("[%s]: reloading home render list\n", t.Format(time.Stamp))

			err := (func() error {
				ctx, cancel := context.WithTimeout(ctx, home.RefreshInterval)
				defer cancel()

				bytes, err := home.homeRender(ctx)
				if err != nil {
					return err
				}

				home.homeBytes.Set(bytes)
				return nil
			})()
			if err != nil {
				io.EPrintf("error reloading instances: ", err.Error())
			}
		}
	}()
}

//go:embed "home.html"
var homeHTMLStr string
var homeTemplate = static.AssetsHomeHome.MustParseShared("home.html", homeHTMLStr)

func (home *Home) homeRender(ctx context.Context) ([]byte, error) {
	var context HomeContext

	// setup a couple of static things
	context.Time = time.Now().UTC()
	context.SelfRedirect = home.Config.SelfRedirect.String()

	// find all the WissKIs
	wissKIs, err := home.Instances.All(ctx)
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
