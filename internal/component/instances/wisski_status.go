package instances

import (
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

// Info represents some info about this WissKI
type Info struct {
	Slug string // The slug of the instance

	Running bool // is the instance running?

	DrupalVersion interface{} // version of drupal being used
}

// Info returns information about this WissKI instance.
func (wisski *WissKI) Info() (info Info, err error) {
	// static properties
	info.Slug = wisski.Slug

	// dynamic properties, TODO: Add more properties here!
	var group errgroup.Group

	group.Go(func() (err error) {
		info.Running, err = wisski.Alive()
		return
	})

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
