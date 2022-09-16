package instances

import (
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

// Info represents some info about this WissKI
type Info struct {
	Slug string // The slug of the instance
	URL  string // The public URL of this instance

	Running      bool     // is the instance running?
	Pathbuilders []string // list of pathbuilders
}

// Info returns information about this WissKI instance.
func (wisski *WissKI) Info(quick bool) (info Info, err error) {
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
	// these execute php code
	if !quick {
		group.Go(func() (err error) {
			info.Pathbuilders, err = wisski.Pathbuilders()
			return
		})
	}

	err = group.Wait()
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
