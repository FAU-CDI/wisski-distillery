package instances

import (
	"fmt"
	"time"

	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

// Info represents some info about this WissKI
type Info struct {
	Slug string // The slug of the instance
	URL  string // The public URL of this instance

	LastRebuild time.Time

	Running      bool              // is the instance running?
	Pathbuilders map[string]string // list of pathbuilders
}

// Info returns information about this WissKI instance.
func (wisski *WissKI) Info(quick bool) (info Info, err error) {
	fmt.Println("call to info")
	// static properties
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()

	// dynamic properties, TODO: Add more properties here!
	var group errgroup.Group

	// quick check if this wisski is running
	group.Go(func() (err error) {
		info.Running, err = wisski.Alive()
		return
	})

	// slower checks for extra properties.
	// these might execute php code or require additional database queries.
	if !quick {
		group.Go(func() error {
			info.Pathbuilders, _ = wisski.AllPathbuilders()
			return nil
		})
		group.Go(func() (err error) {
			info.LastRebuild, _ = wisski.LastRebuild()
			return nil
		})
	}

	err = group.Wait()
	fmt.Println(err)
	return
}

// Alive checks if this WissKI is currently running.
func (wisski *WissKI) Alive() (bool, error) {
	ps, err := wisski.Barrel().Ps(stream.FromNil())
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}
