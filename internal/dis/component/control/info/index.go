package info

import (
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"golang.org/x/sync/errgroup"
)

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = static.AssetsControlIndex.MustParseShared(
	"index.html",
	indexTemplateStr,
)

type indexContext struct {
	Time time.Time

	Config *config.Config

	Instances []ingredient.Information

	TotalCount   int
	RunningCount int
	StoppedCount int

	Backups []models.Export
}

func (nfo *Info) index(r *http.Request) (idx indexContext, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := nfo.Instances.All()
		if err != nil {
			return err
		}

		// get all of their info!
		idx.Instances = make([]ingredient.Information, len(all))
		for i, instance := range all {
			{
				i := i
				instance := instance

				// store the info for this group!
				group.Go(func() (err error) {
					idx.Instances[i], err = instance.Info().Information(true)
					return err
				})
			}
		}

		return nil
	})

	// get the log entries
	group.Go(func() (err error) {
		idx.Backups, err = nfo.SnapshotsLog.For("")
		return
	})

	// get the static properties
	idx.Config = nfo.Config
	idx.Time = time.Now().UTC()

	group.Wait()

	// count how many are running and how many are stopped
	for _, i := range idx.Instances {
		if i.Running {
			idx.RunningCount++
		} else {
			idx.StoppedCount++
		}
	}
	idx.TotalCount = len(idx.Instances)

	return
}
