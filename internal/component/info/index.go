package info

import (
	"html/template"
	"net/http"
	"time"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"golang.org/x/sync/errgroup"
)

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = template.Must(template.New("index.html").Parse(indexTemplateStr))

type indexPageContext struct {
	Time time.Time

	Config *config.Config

	Instances []instances.WissKIInfo

	TotalCount   int
	RunningCount int
	StoppedCount int

	Backups []models.Export
}

func (info *Info) indexPageAPI(r *http.Request) (idx indexPageContext, err error) {
	var group errgroup.Group

	group.Go(func() error {
		// list all the instances
		all, err := info.Instances.All()
		if err != nil {
			return err
		}

		// get all of their info!
		idx.Instances = make([]instances.WissKIInfo, len(all))
		for i, instance := range all {
			{
				i := i
				instance := instance

				// store the info for this group!
				group.Go(func() (err error) {
					idx.Instances[i], err = instance.Info(true)
					return err
				})
			}
		}

		return nil
	})

	// get the log entries
	group.Go(func() (err error) {
		idx.Backups, err = info.Instances.ExportLogFor("")
		return
	})

	// get the static properties
	idx.Config = info.Config
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
